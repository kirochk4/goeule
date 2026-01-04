package eule

import (
	"fmt"
	"strconv"
)

type Value interface {
	valueMark()
	typeOf() String
	fmt.Stringer
}

type Nihil empty
type Boolean bool
type Number float64
type String string
type Closure struct {
	closure env
	params  []varName
	block   block
}
type Native struct {
	fn func(it *Interpreter, args []Value) Value
}
type Table struct {
	Proto *Table
	Pairs map[String]Value
}
type Future empty

func nativePrint(it *Interpreter, args []Value) Value {
	for _, arg := range args {
		fmt.Print(arg)
		fmt.Print(" ")
	}
	fmt.Println()
	return Nihil{}
}

func testValue(v Value) bool {
	switch v := v.(type) {
	case Nihil:
		return false
	case Boolean:
		return bool(v)
	default:
		return true
	}
}

/* == interface ============================================================= */

func (v Nihil) typeOf() String    { return typeOfNihil }
func (v Boolean) typeOf() String  { return "boolean" }
func (v Number) typeOf() String   { return "number" }
func (v String) typeOf() String   { return "string" }
func (v *Closure) typeOf() String { return "function" }
func (v *Native) typeOf() String  { return "function" }
func (v *Table) typeOf() String   { return "table" }
func (v *Future) typeOf() String  { return "future" }

func (v Nihil) String() string    { return stringNihil }
func (v Boolean) String() string  { return strconv.FormatBool(bool(v)) }
func (v Number) String() string   { return formatFloat(v) }
func (v String) String() string   { return string(v) }
func (v *Closure) String() string { return fmt.Sprintf("<function %p>", v) }
func (v *Native) String() string  { return fmt.Sprintf("<function %p>", v) }
func (v *Table) String() string   { return fmt.Sprintf("<table %p>", v) }
func (v *Future) String() string  { return fmt.Sprintf("<future %p>", v) }

/* == marks ================================================================= */

func (v Nihil) valueMark()    {}
func (v Boolean) valueMark()  {}
func (v Number) valueMark()   {}
func (v String) valueMark()   {}
func (v *Closure) valueMark() {}
func (v *Native) valueMark()  {}
func (v *Table) valueMark()   {}
func (v *Future) valueMark()  {}
