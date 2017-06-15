package vminstruction

import (
	"reflect"
	"testing"
)

func TestRead(t *testing.T) {
	input := `1    Add $0, -1, 0
2    Set $1, $10
3    LessThan $0, 0, 1
12    Equals $1,0,5
13    Multiply $0, -1,	$0
14    Return -142`

	expected := []*Instruction{
		createInstruction("Add", []int64{0, -1, 0}, []bool{false, true, true}),
		createInstruction("Set", []int64{1, 10}, []bool{false, false}),
		createInstruction("LessThan", []int64{0, 0, 1}, []bool{false, true, true}),
		createInstruction("Equals", []int64{1, 0, 5}, []bool{false, true, true}),
		createInstruction("Multiply", []int64{0, -1, 0}, []bool{false, true, false}),
		createInstruction("Return", []int64{-142}, []bool{true}),
	}

	got := ReadInstructions(input)

	for i, instr := range got {

		if instr.Name != expected[i].Name {
			t.Fatalf("tests[%d] - Name wrong. expected=%q, got=%q", i, expected[i].Name, instr.Name)
		}

		if !reflect.DeepEqual(instr.Args, expected[i].Args) {
			t.Fatalf("tests[%d] - Arguments wrong. expected=%v, got=%v", i, expected[i].Args, instr.Args)
		}
	}
}

func createInstruction(name string, vals []int64, isAbsolutes []bool) *Instruction {
	ret := make([]*Arg, len(vals))

	for i, v := range vals {
		ret[i] = &Arg{Value: v, IsAbsolute: isAbsolutes[i]}
	}

	return &Instruction{Name: name, Args: ret}
}
