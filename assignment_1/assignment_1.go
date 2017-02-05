package main

import (
	"fmt"
	"io/ioutil"
)

/*

Goal: Scanning a file and correctly assign the token types
Possible tokens:

-> keyword		let, and, in, if, then, else, recur, loop, end

-> identifier	[a-zA-Z_]{[a-zA-Z0-9_]}* (letters, underscores, numbers. cannot start with number)

-> operator		(, ), =, &&, ||, !, <, ==, +, *, -

-> integer		Sequence of digits


*/

var keyword = []string{"and", "else", "end", "if", "in", "let", "loop", "recur", "then"}
var operator = []string{"(", ")", "=", "&&", "||", "!", "<", "==", "+", "*", "-"}

type StateFunc func(byte) int

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func arrContainsString(s string, arr []string) bool {
	for _, val := range arr {
		if s == val {
			return true
		}
	}
	return false
}

// IsKeyword checks, if the given string is a keyword
func IsKeyword(s string) bool {
	return arrContainsString(s, keyword)
}

// IsOperator checks, if the given string is an operator
func IsOperator(s string) bool {
	return arrContainsString(s, operator)
}

// IsLetter returns true, if byte value represents an upper- or lowercase letter from a-z
func IsLetter(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z')
}

// IsDigit returns true, if byte value represents a digit number 0-9
func IsDigit(b byte) bool {
	return (b >= '0' && b <= '9')
}

// IsWhitespace returns true, if byte value represents a whitespace, line break or tab
func IsWhitespace(b byte) bool {
	return (b == ' ' || b == '\n' || b == '\t')
}



func stateWhitespace(b byte) StateFunc {
	if IsWhitespace(b) {
		return stateWhitespace
	} else if IsDigit(b) {
		return stateDigit
	} else if IsLetter(b) {
		return stateLetter
	}

	fmt.Println("Unknown symbol:", b)
	return stateWhitespace
}

func stateDigit(b byte) StateFunc {
	
}

func stateLetter(b byte) StateFunc {

}


func main() {
	dat, err := ioutil.ReadFile("../src/github.com/compilers/sample.txt")
	str := string(dat)

	check(err)
	fmt.Println("Here's the whole text:")
	fmt.Println(str)

	start := 0

	StateEmpty, StateIdentifier, StateInteger := 0, 1, 2

	state := StateEmpty

	for index, value := range []byte(str) {
		if value == ' ' || value == '\n' {
			continue
		}

		switch state {
		case StateEmpty:
			if value == 
		}
	}
}
