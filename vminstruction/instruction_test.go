package vminstruction

import (
	"reflect"
	"testing"
)

func TestRead(t *testing.T) {
	input := `1    Add 0, -1, 0
2    Set 1, 10
3    LessThan 0, 0, 1
12    Equals 1,0,5
13    Multiply 0, -1,	0
14    Return -142`

	expected := []*Instruction{
		&Instruction{Name: "Add", Args: []int64{0, -1, 0}},
		&Instruction{Name: "Set", Args: []int64{1, 10}},
		&Instruction{Name: "LessThan", Args: []int64{0, 0, 1}},
		&Instruction{Name: "Equals", Args: []int64{1, 0, 5}},
		&Instruction{Name: "Multiply", Args: []int64{0, -1, 0}},
		&Instruction{Name: "Return", Args: []int64{-142}},
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
