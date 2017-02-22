package scanner

import "fmt"
import "container/list"

// Scanner contains all scanned tokens and the current index it's on
type Scanner struct {
	filename      string
	scannedTokens []Token
	length        int
	index         int
}

// NewScanner reads file and returns scanner
func NewScanner(path string) (Scanner, *list.List) {
	s := Scanner{
		filename:      path,
		scannedTokens: []Token{},
		length:        0,
		index:         0,
	}

	err := scanFile(&s)
	s.length = len(s.scannedTokens)

	return s, err
}

// Filename returns filename of scanner
func (s Scanner) Filename() string {
	return s.filename
}

// PrintTokens prints all tokens from scanner
func (s Scanner) PrintTokens() {
	for x, t := range s.scannedTokens {
		fmt.Printf("Token #%3d: Type: %2d (= %-12s)", x, t.tokenID, GetStringFromTokenID(t.tokenID))

		if t.tokenVal != nil {
			fmt.Printf(", Value: %v", t.tokenVal)
		}

		fmt.Println()
	}
}

// NextToken increments token index and returns new token. Returns false if end has been already reached
func (s *Scanner) NextToken() (*Token, bool) {
	if s.index >= s.length {
		return nil, false
	}

	t := &s.scannedTokens[s.index]
	s.index++
	return t, true
}

// HasNextToken returns false, when end of token array is reached
func (s Scanner) HasNextToken() bool {
	return s.index < s.length-1
}
