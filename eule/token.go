package eule

import "fmt"

type tokenType string

const (
	tokenLParen      tokenType = "("
	tokenRParen      tokenType = ")"
	tokenLBrace      tokenType = "{"
	tokenRBrace      tokenType = "}"
	tokenLBrack      tokenType = "["
	tokenRBrack      tokenType = "]"
	tokenLAngle      tokenType = "<"
	tokenRAngle      tokenType = ">"
	tokenLAngleAngle tokenType = "<<"
	tokenRAngleAngle tokenType = ">>"

	tokenSemi        tokenType = ";"
	tokenColon       tokenType = ":"
	tokenComma       tokenType = ","
	tokenExcl        tokenType = "!"
	tokenDot         tokenType = "."
	tokenDotDotDot   tokenType = "..."
	tokenQuest       tokenType = "?"
	tokenQuestDot    tokenType = "?."
	tokenQuestLBrack tokenType = "?["
	tokenEq          tokenType = "="
	tokenArrow       tokenType = "=>"

	tokenPlus       tokenType = "+"
	tokenMinus      tokenType = "-"
	tokenStar       tokenType = "*"
	tokenSlash      tokenType = "/"
	tokenPercent    tokenType = "%"
	tokenTildeSlash tokenType = "~/"
	tokenPipe       tokenType = "|"
	tokenAmper      tokenType = "&"
	tokenCircum     tokenType = "^"
	tokenTilde      tokenType = "~"
	tokenPipePipe   tokenType = "||"
	tokenAmperAmper tokenType = "&&"
	tokenQuestQuest tokenType = "??"
	tokenPlusPlus   tokenType = "++"
	tokenMinusMinus tokenType = "--"

	tokenPlusEq        tokenType = "+="
	tokenMinusEq       tokenType = "-="
	tokenStarEq        tokenType = "*="
	tokenSlashEq       tokenType = "/="
	tokenPercentEq     tokenType = "%="
	tokenTildeSlashEq  tokenType = "~/="
	tokenPipeEq        tokenType = "|="
	tokenAmperEq       tokenType = "&="
	tokenCircumEq      tokenType = "^="
	tokenTildeEq       tokenType = "~="
	tokenLAngleAngleEq tokenType = "<<="
	tokenRAngleAngleEq tokenType = ">>="
	tokenPipePipeEq    tokenType = "||="
	tokenAmperAmperEq  tokenType = "&&="
	tokenQuestQuestEq  tokenType = "??="

	tokenExclEq   tokenType = "!="
	tokenEqEq     tokenType = "=="
	tokenLAngleEq tokenType = "<="
	tokenRAngleEq tokenType = ">="

	tokenIdentifier tokenType = "identifier"
	tokenString     tokenType = "string"
	tokenInteger    tokenType = "integer"
	tokenFloat      tokenType = "float"

	tokenVariable tokenType = "variable"
	tokenFunction tokenType = "function"

	tokenNihil tokenType = "nihil"
	tokenTrue  tokenType = "true"
	tokenFalse tokenType = "false"

	tokenIf       tokenType = "if"
	tokenElse     tokenType = "else"
	tokenFor      tokenType = "for"
	tokenForEach  tokenType = "for each"
	tokenIn       tokenType = "in"
	tokenWhile    tokenType = "while"
	tokenDo       tokenType = "do"
	tokenContinue tokenType = "continue"
	tokenBreak    tokenType = "break"
	tokenThrow    tokenType = "throw"
	tokenTry      tokenType = "try"
	tokenCatch    tokenType = "catch"
	tokenFinally  tokenType = "finally"
	tokenReturn   tokenType = "return"
	tokenYield    tokenType = "yield"
	tokenSwitch   tokenType = "switch"
	tokenCase     tokenType = "case"
	tokenDefault  tokenType = "default"
	tokenClass    tokenType = "class"
	tokenExtends  tokenType = "extends"
	tokenAsync    tokenType = "async"
	tokenAwait    tokenType = "await"

	tokenTypeOf tokenType = "typeof"

	tokenNewLine tokenType = "new line"

	tokenError tokenType = "__error"
	tokenEof   tokenType = "__eof"

	tokensCount = iota
)

type token struct {
	tokenType
	line    int
	literal string
}

func (t token) String() string {
	return fmt.Sprintf(
		"%04d: %-12s '%s'",
		t.line,
		t.tokenType,
		shortString(t.literal, 32),
	)
}
