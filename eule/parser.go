package eule

import (
	"fmt"
	"strconv"
)

type fnType int

const (
	fnScript fnType = iota

	fnSync
	fnSyncGen

	fnAsync
	fnAsyncGen
)

type fnCtx struct {
	fnType
	enclosing *fnCtx
	loopCtx   *loopCtx
}

type loopCtx struct {
	enclosing *loopCtx
}

type parser struct {
	scanner   scanner
	cur       token
	prev      token
	errors    []error
	fnCtx     *fnCtx
	isCrushed bool
}

func newParser(scanner scanner) *parser {
	return &parser{
		scanner:   scanner,
		errors:    make([]error, 0),
		fnCtx:     &fnCtx{fnScript, nil, nil},
		isCrushed: false,
	}
}

/* == utils ================================================================= */

func (p *parser) advance() {
	p.prev = p.cur
	p.cur = p.scanner.Scan()
}

func (p *parser) check(type_ tokenType) bool {
	return p.cur.tokenType == type_
}

func (p *parser) match(type_ tokenType) bool {
	if p.check(type_) {
		p.advance()
		return true
	}
	return false
}

func (p *parser) consume(type_ tokenType, message string) {
	if p.match(type_) {
		return
	}
	p.errorAt(p.cur, message)
}

func (p *parser) consumeSemi(message string) {
	if modeAutoSemicolons {
		if p.match(tokenNewLine) {
			return
		}
	}
	p.consume(tokenSemi, message)
}

func (p *parser) consumeIdentifier(message string) *identifierLit {
	p.consume(tokenIdentifier, message)
	return &identifierLit{p.prev.literal}
}

func (p *parser) ignoreNewLine() {
	if modeAutoSemicolons {
		p.match(tokenNewLine)
	}
}

type ParseError struct {
	token   token
	message string
}

func (pe ParseError) Error() string {
	return fmt.Sprintf(
		"line %d at '%s': %s",
		pe.token.line,
		pe.token.literal,
		pe.message,
	)
}

func (p *parser) errorAt(tk token, msg string) {
	panic(ParseError{tk, msg})
}

func (p *parser) fix() {
	defer func() { p.isCrushed = false }()

	for p.cur.tokenType != tokenEof {
		if p.prev.tokenType == tokenSemi { // TODO: what about new line?
			return
		}
		switch p.cur.tokenType {
		case tokenVariable, tokenFunction, tokenIf, tokenFor, tokenForEach,
			tokenWhile, tokenDo, tokenContinue, tokenBreak, tokenThrow,
			tokenTry, tokenReturn, tokenSwitch, tokenCase, tokenDefault:
			return
		}
		p.advance()
	}
}

type precedence int

const (
	precLowest precedence = iota

	precAssign // =
	precOr     // ||
	precAnd    // &&
	precEq     // == !=
	precComp   // < > <= >=
	precTerm   // + -
	precFact   // * / % ~/
	precUnary  // ! + - ~ typeof yield await ++ --
	precCall   // . () {} []

	precHighest
)

var precedences = map[tokenType]precedence{
	tokenEq: precAssign,

	tokenPipePipe: precOr,

	tokenAmperAmper: precAnd,

	tokenEqEq:   precEq,
	tokenExclEq: precEq,

	tokenLAngle:   precComp,
	tokenLAngleEq: precComp,
	tokenRAngle:   precComp,
	tokenRAngleEq: precComp,

	tokenPlus:  precTerm,
	tokenMinus: precTerm,

	tokenStar:    precFact,
	tokenSlash:   precFact,
	tokenPercent: precFact,

	tokenDot:    precCall,
	tokenLParen: precCall,
	tokenLBrack: precCall,
	tokenLBrace: precCall,
}

/* == parse ================================================================= */

func (p *parser) Parse() ([]astDecl, error) {
	script := []astDecl{}
	p.advance()

	for !p.match(tokenEof) {
		decl := p.decl()
		script = append(script, decl)
		if p.isCrushed {
			p.fix()
		}
	}

	if len(p.errors) != 0 {
		return nil, p.errors[0] // TODO: collect all errors
	}

	return script, nil
}

func (p *parser) decl() (decl astDecl) {
	defer catch(func(pe ParseError) {
		p.isCrushed = true
		p.errors = append(p.errors, pe)
	})

	switch {
	case p.match(tokenVariable):
		return p.variableDecl()
	case p.match(tokenFunction):
		return p.functionDecl(false)
	case p.match(tokenAsync):
		if p.match(tokenFunction) {
			return p.functionDecl(true)
		} else if p.match(tokenForEach) {
			return &stmtDecl{p.forEachStmt(true)}
		}
		p.errorAt(p.cur, "ERROR")
		return
	default:
		return &stmtDecl{p.stmt()}
	}
}

func (p *parser) stmt() astStmt {
	switch {
	case p.match(tokenLBrace):
		return p.blockStmt()
	case p.match(tokenIf):
		return p.ifStmt()
	case p.match(tokenFor):
		return p.forStmt()
	case p.match(tokenForEach):
		return p.forEachStmt(false)
	case p.match(tokenWhile):
		return p.whileStmt()
	case p.match(tokenDo):
		return p.doStmt()
	case p.match(tokenContinue):
		return p.continueStmt()
	case p.match(tokenBreak):
		return p.breakStmt()
	case p.match(tokenThrow):
		return p.throwStmt()
	case p.match(tokenTry):
		return p.tryStmt()
	case p.match(tokenReturn):
		return p.returnStmt()
	default:
		expr := &exprStmt{p.expr()}
		if modeAutoSemicolons {
			if !p.match(tokenSemi) {
				p.ignoreNewLine()
			}
		} else {
			p.consumeSemi("ERROR")
		}
		return expr
	}
}

func (p *parser) expr() astExpr {
	return p.precExpr(precAssign)
}

func (p *parser) precExpr(prec precedence) astExpr {
	canAssign := prec <= precAssign
	nud := p.nud(canAssign)

	for prec <= precedences[p.cur.tokenType] {
		nud = p.led(nud, canAssign)
	}

	if canAssign && p.match(tokenEq) {
		p.errorAt(p.prev, "Invalid assignment target.")
	}

	return nud
}

/* == declarations ========================================================== */

func (p *parser) variableDecl() *variableDecl {
	decl := &variableDecl{}

	for {
		vd := varDecl{}
		vd.name = p.consumeIdentifier("ERROR").varName
		if p.match(tokenEq) {
			vd.init = p.expr()
		} else {
			vd.init = &nihilLit{}
		}
		decl.vars = append(decl.vars, vd)

		if !p.match(tokenComma) {
			break
		}
	}

	p.consumeSemi("ERROR")
	return decl
}

func (p *parser) functionDecl(isAsync bool) *functionDecl {
	decl := &functionDecl{}

	isGen := false
	if p.match(tokenStar) {
		isGen = true
	}

	decl.name = p.consumeIdentifier("ERROR").varName
	var isArrow bool
	decl.function, isArrow = p.functionLit(isAsync, isGen)

	if isArrow {
		p.consumeSemi("ERROR")
	}

	return decl
}

/* == statements ============================================================ */

func (p *parser) blockStmt() *blockStmt {
	return &blockStmt{p.block()}
}

func (p *parser) block() block {
	block := make(block, 0)
	for !p.match(tokenRBrace) {
		if p.match(tokenEof) {
			p.errorAt(p.prev, "ERROR")
		}
		decl := p.decl()
		block = append(block, decl)
		if p.isCrushed {
			p.fix()
		}
	}
	return block
}

func (p *parser) ifStmt() *ifStmt {
	stmt := &ifStmt{}
	p.consume(tokenLParen, "ERROR")
	if p.match(tokenVariable) { // var a = b; a
		stmt.init = p.variableDecl()
		stmt.cond = p.expr()
	} else {
		stmt.init = &stmtDecl{&emptyStmt{}}
		stmt.cond = p.expr()
		if p.match(tokenSemi) { // a = b; a
			stmt.init = &stmtDecl{&exprStmt{stmt.cond}}
			stmt.cond = p.expr()
		}
	}
	p.consume(tokenRParen, "ERROR")
	p.ignoreNewLine()
	stmt.then = p.stmt()
	if p.match(tokenElse) {
		stmt.else_ = p.stmt()
	} else {
		stmt.else_ = &emptyStmt{}
	}
	return stmt
}

func (p *parser) forStmt() *forStmt {
	stmt := &forStmt{}
	p.consume(tokenLParen, "ERROR")
	if p.match(tokenVariable) {
		stmt.init = p.variableDecl()
	} else {
		stmt.init = &stmtDecl{&exprStmt{p.expr()}}
		p.consume(tokenSemi, "ERROR")
	}

	if p.match(tokenSemi) {
		stmt.cond = &emptyExpr{}
	} else {
		stmt.cond = p.expr()
		p.consume(tokenSemi, "ERROR")
	}

	if p.match(tokenRParen) {
		stmt.post = &emptyExpr{}
	} else {
		stmt.post = p.expr()
		p.consume(tokenRParen, "ERROR")
	}

	p.consume(tokenRParen, "ERROR")
	p.ignoreNewLine()
	stmt.loop = p.stmt()
	return stmt
}

func (p *parser) forEachStmt(isAsync bool) *forEachStmt {
	stmt := &forEachStmt{}
	return stmt
}

func (p *parser) whileStmt() *whileStmt {
	stmt := &whileStmt{}
	p.consume(tokenLParen, "ERROR")
	stmt.cond = p.expr()
	p.consume(tokenRParen, "ERROR")
	p.fnCtx.loopCtx = &loopCtx{p.fnCtx.loopCtx}
	defer func() { p.fnCtx.loopCtx = p.fnCtx.loopCtx.enclosing }()
	stmt.loop = p.stmt()
	return stmt
}

func (p *parser) doStmt() *doStmt {
	stmt := &doStmt{}
	p.fnCtx.loopCtx = &loopCtx{p.fnCtx.loopCtx}
	defer func() { p.fnCtx.loopCtx = p.fnCtx.loopCtx.enclosing }()
	stmt.loop = p.stmt()
	p.consume(tokenWhile, "ERROR")
	p.consume(tokenLParen, "ERROR")
	stmt.cond = p.expr()
	p.consume(tokenRParen, "ERROR")
	return stmt
}

func (p *parser) continueStmt() *continueStmt {
	if p.fnCtx.loopCtx == nil {
		p.errorAt(p.prev, "ERROR")
	}
	stmt := &continueStmt{}
	p.consumeSemi("E")
	return stmt
}

func (p *parser) breakStmt() *breakStmt {
	if p.fnCtx.loopCtx == nil {
		p.errorAt(p.prev, "ERROR")
	}
	stmt := &breakStmt{}
	p.consumeSemi("expect new line")
	return stmt
}

func (p *parser) throwStmt() *throwStmt {
	stmt := &throwStmt{p.expr()}
	p.consumeSemi("expect new line")
	return stmt
}

func (p *parser) tryStmt() *tryStmt {
	stmt := &tryStmt{}
	p.consume(tokenLBrace, "ERROR")
	stmt.try = p.blockStmt()

	if p.match(tokenCatch) {
		if p.match(tokenLParen) {
			stmt.as = p.consumeIdentifier("ERROR").varName
			p.consume(tokenRParen, "ERROR")
		}
		p.consume(tokenLBrace, "ERROR")
		stmt.catch = p.blockStmt()
	}

	if p.match(tokenFinally) {
		p.consume(tokenLBrace, "ERROR")
		stmt.finally = p.blockStmt()
	}

	if stmt.catch == nil && stmt.finally == nil {
		p.errorAt(p.cur, "expect 'except' or 'finally'")
	}

	return stmt
}

func (p *parser) returnStmt() *returnStmt {
	if p.fnCtx.fnType == fnScript {
		p.errorAt(p.prev, "'return' outside function")
	}
	stmt := &returnStmt{p.expr()}
	p.consumeSemi("expect new line")
	return stmt
}

/* == expressions =========================================================== */

func (p *parser) nud(canAssign bool) astExpr {
	switch {
	case p.match(tokenIdentifier):
		ident := &identifierLit{p.prev.literal}
		if canAssign && p.match(tokenEq) {
			return &assignExpr{ident, p.expr()}
		}
		return ident

	case p.match(tokenNihil):
		return &nihilLit{}
	case p.match(tokenTrue):
		return &booleanLit{true}
	case p.match(tokenFalse):
		return &booleanLit{false}
	case p.match(tokenInteger):
		return parseInteger(p.prev.literal)
	case p.match(tokenFloat):
		return parseFloat(p.prev.literal)
	case p.match(tokenString):
		return &stringLit{
			p.prev.literal[1 : len(p.prev.literal)-2],
		}
	case p.match(tokenLBrace):
		return p.tableLit()
	case p.match(tokenLBrack):
		return p.arrayLit()
	case p.match(tokenFunction):
		var fl *functionLit
		if p.match(tokenStar) {
			fl, _ = p.functionLit(false, true)
		} else {
			fl, _ = p.functionLit(false, false)
		}
		return fl

	case p.match(tokenLParen):
		p.advance()
		group := p.expr()
		p.consume(tokenRParen, "ERROR")
		return group
	case p.match(tokenPlus), p.match(tokenMinus), p.match(tokenExcl),
		p.match(tokenTypeOf), p.match(tokenYield),
		p.match(tokenPlusPlus), p.match(tokenMinusMinus):
		op := p.prev
		right := p.precExpr(precUnary)
		return &prefixExpr{op, right}
	default:
		p.errorAt(p.cur, "ERROR")
		return nil
	}
}

func (p *parser) led(nud astExpr, canAssign bool) astExpr {
	var to astExpr
	switch {
	case p.match(tokenDot):
		p.consume(tokenIdentifier, "ERROR")
		to = &indexExpr{
			left:  nud,
			index: &stringLit{p.prev.literal},
		}
		goto assign
	case p.match(tokenLParen):
		return &callExpr{
			left: nud,
			args: p.args(),
		}
	case p.match(tokenLBrace):
		return &protoTableExpr{
			proto: nud,
			table: p.tableLit(),
		}
	case p.match(tokenLBrack):
		index := p.expr()
		p.consume(tokenRBrack, "ERROR")
		to = &indexExpr{
			left:  nud,
			index: index,
		}
		goto assign

	case p.match(tokenPlus), p.match(tokenMinus),
		p.match(tokenStar), p.match(tokenSlash), p.match(tokenPercent),
		p.match(tokenPipe), p.match(tokenAmper), p.match(tokenCircum):
		op := p.prev
		return &infixExpr{
			left:  nud,
			op:    op,
			right: p.precExpr(precedences[op.tokenType]),
		}
	default:
		panic(unreachable)
	}

assign:
	if canAssign && p.match(tokenEq) {
		return &assignExpr{to, p.expr()}
	}
	return to
}

func parseInteger(literal string) astExpr {
	float, _ := strconv.ParseFloat(literal, 64)
	return &floatLit{float}
}

func parseFloat(literal string) astExpr {
	float, _ := strconv.ParseFloat(literal, 64)
	return &floatLit{float}
}

func (p *parser) tableLit() *tableLit {
	lit := &tableLit{
		pairs: make(map[astExpr]astExpr),
		array: make([]astExpr, 0),
	}
	if p.match(tokenRBrace) {
		return lit
	}
	for {
		var key, val astExpr
		switch {
		case p.match(tokenLBrack): // [key]: val
			key = p.expr()
			p.consume(tokenRBrack, "expect ']'")
			p.consume(tokenColon, "expect ':'")
			val = p.expr()
		case p.match(tokenColon): // :keyval
			val = p.consumeIdentifier("ERROR")
			key = &stringLit{p.prev.literal}
		default: // prop: val
			p.consume(tokenIdentifier, "ERROR")
			key = &stringLit{p.prev.literal}
			p.consume(tokenColon, "ERROR")
			val = p.expr()
		}
		lit.pairs[key] = val
		if !p.match(tokenComma) {
			break
		}
		if p.check(tokenRBrace) {
			break
		}
	}
	p.consume(tokenRBrace, "expect '}'")
	return lit
}

func (p *parser) arrayLit() *tableLit {
	lit := &tableLit{
		pairs: make(map[astExpr]astExpr, 0),
		array: make([]astExpr, 0),
	}
	if p.match(tokenRBrack) {
		return lit
	}
	for {
		lit.array = append(lit.array, p.expr())
		if !p.match(tokenComma) {
			break
		}
		if p.check(tokenRBrack) {
			break
		}
	}
	p.consume(tokenRBrack, "expect ']'")
	return lit
}

func (p *parser) functionLit(
	isAsync bool,
	isGen bool,
) (
	lit *functionLit,
	isArrow bool,
) {
	p.consume(tokenLParen, "ERROR")
	lit.params = p.params()
	p.ignoreNewLine()
	if modeArrowFunctions && p.match(tokenArrow) {
		isArrow = true
		lit.body = block{&stmtDecl{&returnStmt{p.expr()}}}
	} else {
		if isAsync {
			if isGen {
				p.fnCtx = &fnCtx{fnAsyncGen, p.fnCtx, nil}
			} else {
				p.fnCtx = &fnCtx{fnAsync, p.fnCtx, nil}
			}
		} else {
			if isGen {
				p.fnCtx = &fnCtx{fnSyncGen, p.fnCtx, nil}
			} else {
				p.fnCtx = &fnCtx{fnSync, p.fnCtx, nil}
			}
		}
		defer func() { p.fnCtx = p.fnCtx.enclosing }()
		p.consume(tokenLBrace, "ERROR")
		lit.body = p.block()
	}
	return
}

func (p *parser) params() []varName {
	params := []varName{}
	if p.match(tokenRParen) {
		return params
	}
	for {
		params = append(params, p.consumeIdentifier("ERROR").varName)
		if !p.match(tokenComma) {
			break
		}
		if p.check(tokenRParen) {
			break
		}
	}
	p.consume(tokenRParen, "ERROR")
	return params
}

func (p *parser) args() []astExpr {
	args := []astExpr{}
	if p.match(tokenRParen) {
		return args
	}
	for {
		args = append(args, p.expr())
		if !p.match(tokenComma) {
			break
		}
		if p.check(tokenRParen) {
			break
		}
	}
	p.consume(tokenRParen, "ERROR")
	return args
}
