package vminstruction

import (
	"fmt"
	"strconv"
	"strings"
)

type Instruction struct {
	Name string
	Args []*Arg
}

type Arg struct {
	Value      int64
	IsAbsolute bool
}

func (arg *Arg) String() string {
	return fmt.Sprintf("(%d, %v)", arg.Value, arg.IsAbsolute)
}

func ReadInstructions(input string) []*Instruction {
	lines := strings.Split(strings.Replace(input, ",", " ", -1), "\n")

	res := make([]*Instruction, len(lines))

	// example line
	// 7    Set $2, -1
	for index, line := range lines {
		args := strings.Fields(line)[1:]
		res[index] = &Instruction{Name: args[0], Args: convStrings(args[1:])}
	}

	return res
}

func convStrings(arr []string) []*Arg {
	res := make([]*Arg, len(arr))

	for i, s := range arr {
		arg := &Arg{}
		var err error

		if s[0] == '$' {
			arg.IsAbsolute = false
			arg.Value, err = strconv.ParseInt(s[1:], 10, 64)
		} else {
			arg.IsAbsolute = true
			arg.Value, err = strconv.ParseInt(s, 10, 64)
		}

		if err != nil {
			panic(err)
		}

		res[i] = arg
	}

	return res
}
