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
type Float float64
type String string
type Closure struct{}
type Table struct {
	Proto Value
	Pairs map[String]Value
}
type Future empty

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
func (v Float) typeOf() String    { return typeOfFloat }
func (v String) typeOf() String   { return "string" }
func (v *Closure) typeOf() String { return "function" }
func (v *Table) typeOf() String   { return "table" }
func (v *Future) typeOf() String  { return "future" }

func (v Nihil) String() string    { return stringNihil }
func (v Boolean) String() string  { return strconv.FormatBool(bool(v)) }
func (v Float) String() string    { return formatFloat(v) }
func (v String) String() string   { return string(v) }
func (v *Closure) String() string { return fmt.Sprintf("<function %p>", v) }
func (v *Table) String() string   { return fmt.Sprintf("<table %p>", v) }
func (v *Future) String() string  { return fmt.Sprintf("<future %p>", v) }

/* == marks ================================================================= */

func (v Nihil) valueMark()    {}
func (v Boolean) valueMark()  {}
func (v Float) valueMark()    {}
func (v String) valueMark()   {}
func (v *Closure) valueMark() {}
func (v *Table) valueMark()   {}
func (v *Future) valueMark()  {}
