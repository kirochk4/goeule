package eule

import (
	"log"
	"math"
)

type env struct {
	encl *env
	vars map[string]Value
}

func (e *env) define(name string, value Value) {
	e.vars[name] = value
}

func (e *env) store(name string, value Value) {
	e.vars[name] = value
}

func (e *env) load(name string) Value {
	return e.vars[name]
}

type (
	continueSignal empty
	breakSignal    empty
	throwSignal    Value
	returnSignal   Value
)

type Interpreter struct {
	global *Table
	module *Table
	env
	callStack int
	callArgs  []Value // Using only for native functions.
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		global:    &Table{Proto: nil, Pairs: make(map[String]Value)},
		module:    &Table{Proto: nil, Pairs: make(map[String]Value)},
		env:       env{encl: nil, vars: make(map[string]Value)},
		callStack: 0,
		callArgs:  []Value{},
	}
}

func (it *Interpreter) GetArg(index int) Value {
	if index < 0 || index > len(it.callArgs)-1 {
		return Nihil{}
	}
	return it.callArgs[index]
}

func (it *Interpreter) Interpret(source []byte) {
	it.env.vars["print"] = &Native{fn: nativePrint}
	s := newScanner(source)
	p := newParser(s)
	tree, err := p.Parse()
	if err != nil {
		log.Fatal(err)
	}
	for _, node := range tree {
		it.eval(node)
	}
}

func (it *Interpreter) eval(node astNode) Value {
	switch node := node.(type) {
	/* == declarations ====================================================== */
	case *errorDecl:
		panic(unreachable)
	case *variableDecl:
		for _, decl := range node.vars {
			it.define(decl.name, it.eval(decl.init))
		}
		return nil
	case *functionDecl:
		it.define(node.name, it.eval(node.function))
		return nil
	case *stmtDecl:
		return it.eval(node.stmt)
	/* == statements ======================================================== */
	case *emptyStmt:
		return nil
	case *blockStmt:
		for _, decl := range node.block {
			it.eval(decl)
		}
		return nil
	case *ifStmt:
		return it.ifStmt(node)
	case *forStmt:
		return it.forStmt(node)
	case *forEachStmt:
		return it.forEachStmt(node)
	case *whileStmt:
		return it.whileStmt(node)
	case *doStmt:
		return it.doStmt(node)
	case *continueStmt:
		panic(continueSignal{})
	case *breakStmt:
		panic(breakSignal{})
	case *throwStmt:
		panic(throwSignal(it.eval(node.throw)))
	case *tryStmt:
		return it.tryStmt(node)
	case *returnStmt:
		panic(returnSignal(it.eval(node.value)))
	case *exprStmt:
		return it.eval(node.expr)
	/* == expressions ======================================================= */
	case *identifierLit:
		return it.load(node.varName)
	case *nihilLit:
		return Nihil{}
	case *booleanLit:
		return Boolean(node.value)
	case *integerLit:
		return Float(node.value)
	case *floatLit:
		return Float(node.value)
	case *stringLit:
		return String(node.value)
	case *tableLit:
		return it.tableLit(node)
	case *functionLit:
		return it.functionLit(node)

	case *emptyExpr:
		return nil
	case *assignExpr:
		return it.assignExpr(node)
	case *prefixExpr:
		return it.prefixExpr(node)
	case *infixExpr:
		return it.infixExpr(node)
	case *postfixExpr:
		return nil
	case *callExpr:
		callee := it.eval(node.left)
		switch callee := callee.(type) {
		case *Native:
			args := []Value{}
			for _, arg := range node.args {
				args = append(args, it.eval(arg))
			}
			return callee.fn(it, args)
		default:
			panic("ERROR")
		}
	case *indexExpr:
		return loadIndex(
			it.eval(node.left),
			it.eval(node.index),
		)
	case *protoTableExpr:
	default:
		panic(unreachable)
	}
	panic(unreachable)
}

func (it *Interpreter) beginScope() {
	it.env = env{encl: &it.env, vars: make(map[string]Value)}
}

func (it *Interpreter) endScope() {
	it.env = *it.env.encl
}

func (it *Interpreter) ifStmt(node *ifStmt) Value {
	it.beginScope()
	defer it.endScope()
	it.eval(node.init)
	if testValue(it.eval(node.cond)) {
		it.eval(node.then)
	} else {
		it.eval(node.else_)
	}
	return nil
}

func (it *Interpreter) forStmt(node *forStmt) Value {
	defer catch(func(_ breakSignal) {})
	it.beginScope()
	defer it.endScope()
	it.eval(node.init)
	for testValue(it.eval(node.cond)) {
		func() {
			defer catch(func(_ continueSignal) {})
			it.eval(node.loop)
		}()
		it.eval(node.post)
	}
	return nil
}

func (it *Interpreter) forEachStmt(node *forEachStmt) Value {
	return nil
}

func (it *Interpreter) whileStmt(node *whileStmt) Value {
	defer catch(func(_ breakSignal) {})
	for testValue(it.eval(node.cond)) {
		func() {
			defer catch(func(_ continueSignal) {})
			it.eval(node.loop)
		}()
	}
	return nil
}

func (it *Interpreter) doStmt(node *doStmt) Value {
	defer catch(func(_ breakSignal) {})
	for {
		func() {
			defer catch(func(_ continueSignal) {})
			it.eval(node.loop)
		}()
		if !testValue(it.eval(node.cond)) {
			break
		}
	}
	return nil
}

func (it *Interpreter) tryStmt(node *tryStmt) Value {
	if node.finally != nil {
		defer func() { it.eval(node.finally) }()
	}

	if node.catch != nil {
		defer catch(func(throw throwSignal) {
			it.beginScope()
			defer it.endScope()
			it.define(node.as, throw)
			it.eval(node.catch)
		})
	}

	it.eval(node.try)
	return nil
}

func (it *Interpreter) tableLit(node *tableLit) Value {
	return nil
}

func (it *Interpreter) functionLit(node *functionLit) Value {
	return nil
}

func (it *Interpreter) assignExpr(node *assignExpr) Value {
	value := it.eval(node.right)
	switch left := node.left.(type) {
	case *identifierLit:
		it.store(left.varName, value)
	case *indexExpr:
		storeIndex(
			it.eval(left.left),
			it.eval(left.index),
			value,
		)
	default:
		panic(unreachable)
	}
	return value
}

func (it *Interpreter) prefixExpr(node *prefixExpr) Value {
	rVal := it.eval(node.right)

	switch node.op.tokenType {
	case tokenMinus:
		return -rVal.(Float)
	case tokenPlus:
		return +rVal.(Float)

	case tokenExcl:
		return Boolean(testValue(rVal))

	default:
		panic(unreachable)
	}
}

func (it *Interpreter) infixExpr(node *infixExpr) Value {
	lVal := it.eval(node.left)
	rVal := it.eval(node.right)

	switch node.op.tokenType {
	case tokenEq:
		return Boolean(rVal == lVal)
	case tokenExclEq:
		return Boolean(rVal != lVal)

	case tokenLAngle:
		return Boolean(rVal.(Float) < lVal.(Float))
	case tokenLAngleEq:
		return Boolean(rVal.(Float) <= lVal.(Float))
	case tokenRAngle:
		return Boolean(rVal.(Float) > lVal.(Float))
	case tokenRAngleEq:
		return Boolean(rVal.(Float) >= lVal.(Float))

	case tokenPlus:
		return rVal.(Float) + lVal.(Float)
	case tokenMinus:
		return rVal.(Float) - lVal.(Float)
	case tokenStar:
		return rVal.(Float) * lVal.(Float)
	case tokenSlash:
		return rVal.(Float) / lVal.(Float)
	case tokenPercent:
		return Float(math.Mod(float64(rVal.(Float)), float64(lVal.(Float))))

	default:
		panic(unreachable)
	}
}

func storeIndex(object Value, index Value, value Value) {}
func loadIndex(object Value, index Value) Value {
	return nil
}
