package eule

import (
	"fmt"
	"math"
	"strconv"
)

/* == header ================================================================ */

// Interpreter version.
const Version = "0.0.0"

const (
	uint8Max  uint8  = math.MaxUint8
	uint16Max uint16 = math.MaxUint16

	uint8Count  int = math.MaxUint8 + 1
	uint16Count int = math.MaxUint16 + 1

	unreachable string = "unreachable"
)

// Debug information.
const (
	debugPrintTokens = true
	debugPrintAst    = true
)

// Interpreter modes.
const (
	// Automaticly inserts semicolons.
	modeAutoSemicolons = false
	// Allows class syntax.
	modeObjectOriented = false
	// Allows `function() => toReturn`` syntax.
	modeArrowFunctions = false
)

// Language constants.
const (
	stringThis  = "this"
	stringSuper = "super"

	stringNihil    = "void"
	stringVariable = "var"
	stringFunction = "function"

	typeOfNihil = stringNihil
	typeOfFloat = "number"
)

type empty = struct{}
type varName = string

/* == lib =================================================================== */

func panicf(format string, a ...any) {
	panic(fmt.Sprintf(format, a...))
}

func formatFloat(f Float) string {
	return strconv.FormatFloat(float64(f), 'g', -1, 64)
}

func coverString(str string, width int, char byte, space int) string {
	runes := []rune(str)
	if len(runes)+space*2 >= width {
		return str
	}

	left := (width - len(runes) - space*2) / 2
	right := left + (width-len(runes)-space*2)%2

	return fmt.Sprintf(
		"%*c%*c%s%*c%*c",
		left, char, space, ' ',
		str,
		right, char, space, ' ',
	)
}

func shortString(str string, length int) string {
	runes := []rune(str)
	if length < len(runes) {
		return string(runes[:length])
	}
	return str
}

func catch[E any](onCatch func(E)) {
	if p := recover(); p != nil {
		if pe, ok := p.(E); ok {
			onCatch(pe)
		} else {
			panic(p)
		}
	}
}
