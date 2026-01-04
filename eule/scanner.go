package eule

import (
	"fmt"
)

func init() {
	if modeArrowFunctions {
		dual[dualSymbol{'=', '>'}] = tokenArrow
	}
	if modeObjectOriented {
		keywords["class"] = tokenClass
		keywords["extends"] = tokenExtends
	}
}

const eofByte = 0

type scanner struct {
	source []byte
	cursor int
	start  int
	line   int
	inl    bool // Insert new line token.
}

func newScanner(source []byte) scanner {
	return scanner{
		source: source,
		cursor: 0,
		line:   1,
		inl:    false,
	}
}

func (s *scanner) Scan() token {
begin:
	line := s.line
	s.skipWhite()

	s.start = s.cursor

	if modeAutoSemicolons {
		if s.inl && line < s.line {
			return s.makeToken(tokenNewLine)
		}
	}

	if s.isAtEnd() {
		return s.makeToken(tokenEof)
	}

	char := s.advance()

	if char == '/' && s.current() == '/' {
		s.skipLineComment()
		goto begin
	}

	if char == '/' && s.current() == '*' {
		if errToken, hasErr := s.skipMultiLineComment(); hasErr {
			return errToken
		} else {
			goto begin
		}
	}

	switch {
	case isAlpha(char):
		return s.identifier()
	case isDigit(char, 10):
		return s.number()
	case char == '"':
		return s.string()
	}

	if t, ok := triple[tripleSymbol{char, s.current(), s.peek()}]; ok {
		s.advance()
		s.advance()
		return s.makeToken(t)
	} else if t, ok := dual[dualSymbol{char, s.current()}]; ok {
		s.advance()
		return s.makeToken(t)
	} else if t, ok := mono[char]; ok {
		return s.makeToken(t)
	}

	return s.errorToken("ERROR 1")
}

func (s *scanner) isAtEnd() bool {
	return s.current() == eofByte
}

func (s *scanner) current() byte {
	if s.cursor >= len(s.source) {
		return eofByte
	}
	return s.source[s.cursor]
}

func (s *scanner) peek() byte {
	if s.cursor+1 >= len(s.source) {
		return eofByte
	}
	return s.source[s.cursor+1]
}

func (s *scanner) advance() byte {
	char := s.current()
	s.cursor++
	return char
}

func (s *scanner) makeToken(t tokenType) token {
	_, s.inl = inlAfter[t]

	literal := string(s.source[s.start:s.cursor])
	tk := token{t, s.line, literal}

	if debugPrintTokens {
		fmt.Println(tk)
	}

	return tk
}

func (s *scanner) errorToken(message string) token {
	return token{
		tokenType: tokenError,
		line:      s.line,
		literal:   message,
	}
}

func (s *scanner) skipWhite() {
	for {
		switch char := s.current(); char {
		case '\n':
			s.line++
			fallthrough
		case ' ', '\r', '\t':
			s.advance()
		default:
			return
		}
	}
}

func (s *scanner) skipLineComment() {
	for s.current() != '\n' && !s.isAtEnd() {
		s.advance()
	}
}

func (s *scanner) skipMultiLineComment() (token, bool) {
	for (s.current() != '*' && s.peek() != '/') ||
		!s.isAtEnd() {
		if s.current() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		return s.errorToken("ERROR 2"), true
	}

	return token{}, false
}

func (s *scanner) identifierType() tokenType {
	if t, ok := keywords[string(s.source[s.start:s.cursor])]; ok {
		return t
	}
	return tokenIdentifier
}

func (s *scanner) identifier() token {
	for isAlpha(s.current()) || isDigit(s.current(), 10) {
		s.advance()
	}
	return s.makeToken(s.identifierType())
}

/*
 * 3.14, 3.1_4, 3., 3..method() - good
 * 3_.14, 3._14, 3.1__4, 3abc, 3.14abc - bad
 */
func (s *scanner) number() token {
	// Calculate base.
	var base uint8 = 10
	if b, ok := intBases[lowerChar(s.current())]; ok {
		base = b
		s.advance()
		if !isDigit(s.current(), base) {
			return s.errorToken("ERROR")
		}
	}

	// readDigits returns false if literal ends with underscore.
	readDigits := func() bool {
		allowUnderscore := false
		for isDigit(s.current(), base) ||
			(allowUnderscore && s.current() == '_') {
			allowUnderscore = s.current() != '_'
			s.advance()
		}
		return allowUnderscore
	}

	numberType := tokenInteger

	// Read integer.
	if !readDigits() { // Ends with underscore or double underscore.
		return s.errorToken("ERROR")
	} else if isAlpha(s.current()) { // '3abc' not allowed.
		return s.errorToken("ERROR")
	}

	// Read float.
	if base == 10 && s.current() == '.' {
		numberType = tokenFloat
		s.advance()
		if !readDigits() { // Ends with underscore or double underscore.
			return s.errorToken("ERROR")
		} else if isAlpha(s.current()) { // '3.14abc' not allowed.
			return s.errorToken("ERROR")
		}
	}

	return s.makeToken(numberType)
}

func (s *scanner) string() token {
	var previous byte = eofByte
	for !(s.current() == '"' && previous != '\\') && !s.isAtEnd() {
		if s.current() == '\n' {
			return s.errorToken("ERROR")
		}
		previous = s.advance()
	}
	if s.isAtEnd() {
		return s.errorToken("ERROR")
	}
	s.advance() // Read ending '"'.
	return s.makeToken(tokenString)
}

func isAlpha(char byte) bool {
	return 'a' <= char && char <= 'z' ||
		'A' <= char && char <= 'Z' ||
		char == '_'
}

func isDigit(char byte, base uint8) bool {
	if base <= 10 {
		return '0' <= char && char <= '0'+base-1
	}
	return ('0' <= char && char <= '9') ||
		('a' <= char && char <= 'a'+base-1) ||
		('A' <= char && char <= 'A'+base-1)
}

func lowerChar(char byte) byte {
	return ('a' - 'A') | char
}

var intBases = map[byte]uint8{
	'x': 16,
	'o': 8,
	'b': 2,
}

var inlAfter = map[tokenType]empty{
	tokenRParen: {},
	tokenRBrace: {},
	tokenRBrack: {},

	tokenIdentifier: {},
	tokenString:     {},
	tokenInteger:    {},
	tokenFloat:      {},
	tokenNihil:      {},
	tokenTrue:       {},
	tokenFalse:      {},

	tokenReturn: {},
	tokenYield:  {},
}

var mono = map[byte]tokenType{
	'(': tokenLParen,
	')': tokenRParen,
	'{': tokenLBrace,
	'}': tokenRBrace,
	'[': tokenLBrack,
	']': tokenRBrack,

	';': tokenSemi,
	':': tokenColon,
	',': tokenComma,
	'!': tokenExcl,
	'.': tokenDot,
	'?': tokenQuest,
	'=': tokenEq,

	'+': tokenPlus,
	'-': tokenMinus,
	'*': tokenStar,
	'/': tokenSlash,
	'%': tokenPercent,
	'|': tokenPipe,
	'&': tokenAmper,
	'^': tokenCircum,
	'~': tokenTilde,

	'<': tokenLAngle,
	'>': tokenRAngle,
}

type dualSymbol [2]byte

var dual = map[dualSymbol]tokenType{
	{'<', '<'}: tokenLAngleAngle,
	{'>', '>'}: tokenRAngleAngle,
	{'?', '.'}: tokenQuestDot,
	{'?', '['}: tokenQuestLBrack,
	{'~', '/'}: tokenTildeSlash,
	{'|', '|'}: tokenPipePipe,
	{'&', '&'}: tokenAmperAmper,
	{'?', '?'}: tokenQuestQuest,

	{'+', '+'}: tokenPlusPlus,
	{'-', '-'}: tokenMinusMinus,

	{'+', '='}: tokenPlusEq,
	{'-', '='}: tokenMinusEq,
	{'*', '='}: tokenStarEq,
	{'/', '='}: tokenSlashEq,
	{'%', '='}: tokenPercentEq,
	{'|', '='}: tokenPipeEq,
	{'&', '='}: tokenAmperEq,
	{'^', '='}: tokenCircumEq,
	{'~', '='}: tokenTildeEq,
	{'!', '='}: tokenExclEq,
	{'=', '='}: tokenEqEq,
	{'<', '='}: tokenLAngleEq,
	{'>', '='}: tokenRAngleEq,
}

type tripleSymbol [3]byte

var triple = map[tripleSymbol]tokenType{
	{'.', '.', '.'}: tokenDotDotDot,
	{'~', '/', '='}: tokenTildeSlashEq,
	{'<', '<', '='}: tokenLAngleAngleEq,
	{'>', '>', '='}: tokenRAngleAngleEq,
	{'|', '|', '='}: tokenPipePipeEq,
	{'&', '&', '='}: tokenAmperAmperEq,
	{'?', '?', '='}: tokenQuestQuestEq,
}

var keywords = map[string]tokenType{
	stringVariable: tokenVariable,
	stringFunction: tokenFunction,

	stringNihil: tokenNihil,
	"true":      tokenTrue,
	"false":     tokenFalse,

	"if":       tokenIf,
	"else":     tokenElse,
	"for":      tokenFor,
	"foreach":  tokenForEach,
	"in":       tokenIn,
	"while":    tokenWhile,
	"do":       tokenDo,
	"continue": tokenContinue,
	"break":    tokenBreak,
	"throw":    tokenThrow,
	"try":      tokenTry,
	"catch":    tokenCatch,
	"finally":  tokenFinally,
	"return":   tokenReturn,
	"yield":    tokenYield,
	"switch":   tokenSwitch,
	"case":     tokenCase,
	"default":  tokenDefault,
	"async":    tokenAsync,
	"await":    tokenAwait,

	"typeof": tokenTypeOf,
}
