package vminstruction

import (
	"strconv"
	"strings"
)

type Instruction struct {
	Name string
	Args []int64
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

func convStrings(arr []string) []int64 {
	res := make([]int64, len(arr))

	for i, s := range arr {
		conv, err := strconv.ParseInt(s, 10, 64)

		if err != nil {
			panic(err)
		}

		res[i] = conv
	}

	return res
}
