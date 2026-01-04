package main

import (
	"os"

	"github.com/kirochk4/goeule/eule"
)

func main() {
	src, _ := os.ReadFile("script.eul")
	eule.NewInterpreter().Interpret(src)
}
