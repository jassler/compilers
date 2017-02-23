package scanner

import (
	"container/list"
	"fmt"
	"io/ioutil"
	"strconv"
	"unicode/utf8"
)

/*

Goal: Scanning a file and correctly assign the token types
Possible tokens:

-> keyword		let, and, in, if, then, else, recur, loop, end

-> identifier	[a-zA-Z_]{[a-zA-Z0-9_]}* (letters, underscores, numbers. cannot start with number)

-> operator		(, ), =, &&, ||, !, <, ==, +, *, -

-> integer		Sequence of digits
*/

// List of keywords. We get those keywords from token.go
var keyword []string

// List of operators. We get those operators from token.go
var operator []string

// List of single character operators. Used later on to check, if character is in operator more efficiently
var operatorStart = []rune{}

// lineNum for tokens
var lineNum int

// To simulate states I use function pointers.
// eg. if string starts with a letter, it'll probably be an identifier -> stateIdentifierID
type stateFunction func(rune, *list.List) (*stateFunction, bool)

// To avoid initialization loop, those variables will be initialized in main
// Note that those IDs exist since GO doesn't support direct function references
var stateWhitespaceID stateFunction
var stateIntegerID stateFunction
var stateIdentifierID stateFunction
var stateOperatorID stateFunction

func createError(msg string, line int) error {
	return fmt.Errorf("Syntax Error in line %d: %s", line, msg)
}

// arrContainsString checks, if given string appears in a string array
// Used for operators and keywords
func arrContainsString(s string, arr []string) bool {
	for _, val := range arr {
		if s == val {
			return true
		}
	}
	return false
}

// arrContainsRune checks, if given rune appears in a rune array
// Used for operators (operator starts with rune)
func arrContainsRune(r rune, arr []rune) bool {
	for _, val := range arr {
		if r == val {
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

// isPossibleOperator checks, if the given string could be a substring of an operator
// Used to distinguish between operators. eg. input = "++&&=" consists of "+", "+", "&&" and "=", so we kind of split them apart
func isPossibleOperator(s string) int {
	// -1: Definitely not an operator
	//  0: Could be one
	//  1: Definitely an operator
	returnVal := -1

	// Only works for operators of length 1 or 2!
	length := utf8.RuneCountInString(s)
	if length == 1 {
		if IsOperator(s) {
			returnVal = 1
		} else if isStartOfOperator([]rune(s)[0]) {
			returnVal = 0
		}
	} else if length == 2 {
		// Possible extension problem, when operators get longer than 2 characters
		if IsOperator(s) {
			returnVal = 1
		}
	}
	return returnVal
}

// IsStartOfOperator returns true, if given letter is the beginning of an operator
func isStartOfOperator(b rune) bool {
	return arrContainsRune(b, operatorStart)
}

// IsLetter returns true, if rune value represents an upper- or lowercase letter from a-z
func IsLetter(b rune) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b == '_')
}

// IsDigit returns true, if rune value represents a digit number 0-9
func IsDigit(b rune) bool {
	return (b >= '0' && b <= '9')
}

// IsWhitespace returns true, if rune value represents a whitespace, line break, tab or carriage return (13)
func IsWhitespace(b rune) bool {
	if b == '\n' {
		lineNum++
	}
	return (b == ' ' || b == '\n' || b == '\t' || b == 13)
}

/*
Here start the state functions that are called with the function pointers.
We can be in 4 states: Whitespace, Operator, Identifier or Integer.
Depending on what rune / character comes next, we switch states or stay in the same one.
If we switch states that means we completed reading a value type. eg. if we switch from stateInteger -> stateWhitespace, we just read an integer.

Each function takes the character as parameter and returns the new state it's in and whether or not it switched states
*/

// We are in stateWhitespace
func stateWhitespace(b rune, l *list.List) (*stateFunction, bool) {
	newState := &stateWhitespaceID

	if IsWhitespace(b) {
		// ignoring all other whitespaces
		newState = &stateWhitespaceID

	} else if IsDigit(b) {
		// start of an integer
		newState = &stateIntegerID

	} else if IsLetter(b) {
		// start of an identifier
		newState = &stateIdentifierID

	} else if isStartOfOperator(b) {
		// start of an operator
		newState = &stateOperatorID

	} else {
		// invalid character
		l.PushBack(createError("Invalid token '"+string(b)+"' (code: "+strconv.Itoa(int(b))+")", lineNum))
	}

	return newState, newState != &stateWhitespaceID
}

// We are in stateInteger
// Note that letters are not allowed to be in here
func stateInteger(b rune, l *list.List) (*stateFunction, bool) {
	newState := &stateWhitespaceID

	if IsWhitespace(b) {
		newState = &stateWhitespaceID

	} else if IsDigit(b) {
		newState = &stateIntegerID

	} else if IsLetter(b) {
		l.PushBack(createError("Identifiers can't start with numbers", lineNum))
		newState = &stateIdentifierID

	} else if isStartOfOperator(b) {
		newState = &stateOperatorID

	} else {
		l.PushBack(createError("Invalid token '"+string(b)+"' (code: "+strconv.Itoa(int(b))+")", lineNum))
	}

	return newState, newState != &stateIntegerID
}

// We are in stateIdentifier
func stateIdentifier(b rune, l *list.List) (*stateFunction, bool) {
	newState := &stateWhitespaceID

	if IsWhitespace(b) {
		newState = &stateWhitespaceID

	} else if IsDigit(b) {
		newState = &stateIdentifierID

	} else if IsLetter(b) {
		newState = &stateIdentifierID

	} else if isStartOfOperator(b) {
		newState = &stateOperatorID

	} else {
		l.PushBack(createError("Invalid token '"+string(b)+"' (code: "+strconv.Itoa(int(b))+")", lineNum))
	}

	return newState, newState != &stateIdentifierID
}

// We are in stateOperator
func stateOperator(b rune, l *list.List) (*stateFunction, bool) {
	newState := &stateWhitespaceID

	if IsWhitespace(b) {
		newState = &stateWhitespaceID

	} else if IsDigit(b) {
		newState = &stateIntegerID

	} else if IsLetter(b) {
		newState = &stateIdentifierID

	} else if isStartOfOperator(b) {
		newState = &stateOperatorID

	} else {
		l.PushBack(createError("Invalid token '"+string(b)+"' (code: "+strconv.Itoa(int(b))+")", lineNum))
	}

	return newState, newState != &stateOperatorID
}

// When switching states, we interpret the input we've just gotten.
// eg. if we switched from stateInteger -> stateWhitespace, then we just read an integer
// eg. if we switched from stateIdentifier -> stateWhitespace, it could be either an identifier or a keyword
// eg. if we switched from stateWhitespace -> stateInteger, we only had white spaces. Ignore those
func interpretInput(slice string, length int, state *stateFunction, sc *Scanner, l *list.List, line int) {
	token := Token{tokenID: -1, tokenVal: nil, lineNum: line}

	switch state {
	case &stateIdentifierID:

		// Check if identifier could be a keyword or not
		if IsKeyword(slice) {
			token.tokenID = MapKeyword(slice)
		} else {
			token.tokenID = TokenIdentifier
			token.tokenVal = slice
		}

	case &stateIntegerID:
		token.tokenID = TokenInteger
		i, _ := strconv.Atoi(slice)
		token.tokenVal = i

	case &stateOperatorID:

		replace := false
		// Found an operator
		// Note that it could be multiple operators clustered together. So we have to pull those apart
		start, end := 0, 1
		for ; end <= length; end++ {
			outcome := isPossibleOperator(slice[start:end])

			// if outcome == -1, then we probably finished our last valid operator and start with a new one
			if outcome == -1 {
				s := slice[start : end-1]
				id := MapOperator(s)

				if id < 0 {
					l.PushBack(createError("\""+s+"\" is not a valid operator", line))
				} else {
					if replace {
						token = Token{tokenID: id, tokenVal: nil}
					} else {
						token.tokenID = id
					}
					sc.scannedTokens = append(sc.scannedTokens, token)
					replace = true
				}

				start = end - 1
			}
		}

		// there could be an operator at the end
		if start != length {
			s := slice[start:length]
			id := MapOperator(s)

			if id < 0 {
				l.PushBack(createError("\""+s+"\" is not a valid operator", line))
			} else {
				if replace {
					token = Token{tokenID: id, tokenVal: nil}
				} else {
					token.tokenID = id
				}

				sc.scannedTokens = append(sc.scannedTokens, token)
			}
		}

		return
	}

	if token.tokenID >= 0 {
		sc.scannedTokens = append(sc.scannedTokens, token)
	}
}

// ScanFile scans file
func scanFile(sc *Scanner) *list.List {
	dat, err := ioutil.ReadFile(sc.filename)
	errList := list.New()

	if err != nil {
		errList.PushBack(err)
		return errList
	}

	// convert file content to string.
	str := string(dat)

	if len(str) == 0 {
		errList.PushBack(createError("Empty file", 1))
		return errList
	}

	// initialize state IDs with corresponding function
	stateWhitespaceID = stateWhitespace
	stateIntegerID = stateInteger
	stateIdentifierID = stateIdentifier
	stateOperatorID = stateOperator

	// initialize keyword and operator array
	keyword = GetKeywords()
	operator = GetOperators()

	// initialize operatorStart with first character of each operator (eg. "&&" -> '&')
	for _, val := range operator {
		r := []rune(val)[0]
		if !arrContainsRune(r, operatorStart) {
			operatorStart = append(operatorStart, r)
		}
	}

	// lineNum starts at 1
	lineNum = 1

	// for substrings we always slice it from start to the current index
	start := 0

	// currentState set to white space at the beginning
	var currentState = &stateWhitespaceID
	// newState always set for each character we read
	var newState *stateFunction
	// bool to check, if we changed from our last state
	var changed bool

	// slices to take apart our tokens
	var slice string

	// we use index and value after our loop, so we declare those outside
	var index int
	var value rune

	// for correct line numbers I need to subtract 1 from lineNum if last token was a linebreak
	var diff int

	for index, value = range []rune(str) {
		newState, changed = (*currentState)(value, errList)

		// state didn't change, read next character
		if !changed {
			continue
		}

		diff = 0
		if value == '\n' {
			diff = 1
		}

		// state changed. Interpret our slice
		slice = str[start:index]
		interpretInput(slice, index-start, currentState, sc, errList, lineNum-diff)

		// start set to current index, the start of our next token
		start = index

		// currentState changed to the new one
		currentState = newState
	}

	// loop doesn't get very last operator, so we have to include it here
	interpretInput(str[start:index+1], index-start+1, currentState, sc, errList, lineNum)

	if errList.Len() != 0 {
		return errList
	}
	return nil
}
