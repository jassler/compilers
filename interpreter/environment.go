package interpreter

import (
	"fmt"

	"os"

	"github.com/simplang/token"
)

type element struct {
	name  string
	value int64
}

type environment struct {
	elements []*element
}

func throwError(msg string, t *token.Token) {
	if t != nil {
		fmt.Printf("Error: %s (line %d.%d)\n", msg, t.Line, t.Column)
	} else {
		fmt.Printf("Error: %s\n", msg)
	}
	os.Exit(1)
}

func (e *element) String() string {
	return fmt.Sprintf("{%s = %d}", e.name, e.value)
}

func (e *environment) String() string {
	if len(e.elements) == 0 {
		return "[]"
	}

	s := "[" + e.elements[0].String()

	for i := 1; i < len(e.elements); i++ {
		s += ", " + e.elements[i].String()
	}

	s += "]"
	return s
}

func (e *environment) indexOfElement(name string) int {
	for i := len(e.elements) - 1; i >= 0; i-- {
		if e.elements[i].name == name {
			return i
		}
	}

	return -1
}

func (e *environment) setValue(name string, val int64) bool {
	if i := e.indexOfElement(name); i >= 0 {
		e.elements[i].value = val
		return true
	}

	return false
}

func (e *environment) getValue(t *token.Token) int64 {
	if i := e.indexOfElement(t.Literal); i >= 0 {
		return e.elements[i].value
	}

	throwError(fmt.Sprintf("Variable '%s' not defined (line %d.%d)", t.Literal, t.Line, t.Column), t)
	return 0
}

func (e *environment) appendElement(el *element) {
	e.elements = append(e.elements, el)
}

func (e *environment) removeElements(am int) {
	e.elements = e.elements[:len(e.elements)-am]
}
