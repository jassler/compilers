package vm

import "github.com/simplang/vminstruction"
import "strings"
import "fmt"
import "os"

type VirtualMachine struct {
	Instructions   []*vminstruction.Instruction
	Funcs          []func(...*vminstruction.Arg)
	ProgramCounter int64
	CallStack      []call
	Values         []int64
	ValPointer     int64
}

type call struct {
	pc  int64 // index of call instruction (program counter)
	vp  int64 // original index of the value pointer
	dst int64 // destination where to write the result
}

func New(instr []*vminstruction.Instruction) *VirtualMachine {
	vm := &VirtualMachine{
		Instructions:   instr,
		Funcs:          make([]func(...*vminstruction.Arg), len(instr)),
		ProgramCounter: 0,
		CallStack:      []call{},
		Values:         make([]int64, 1000000),
		ValPointer:     0,
	}

	type tuple struct {
		f    func(...*vminstruction.Arg)
		args int
	}

	instrMap := map[string]tuple{
		"move":       tuple{args: 2, f: vm.Move},
		"set":        tuple{args: 2, f: vm.Set},
		"add":        tuple{args: 3, f: vm.Add},
		"multiply":   tuple{args: 3, f: vm.Multiply},
		"negate":     tuple{args: 2, f: vm.Negate},
		"not":        tuple{args: 2, f: vm.Not},
		"jump":       tuple{args: 1, f: vm.Jump},
		"jumpifzero": tuple{args: 2, f: vm.JumpIfZero},
		"call":       tuple{args: 3, f: vm.Call},
		"return":     tuple{args: 1, f: vm.Return},
		"lessthan":   tuple{args: 3, f: vm.LessThan},
		"equals":     tuple{args: 3, f: vm.Equals},
	}

	valid := true
	for ind, instr := range vm.Instructions {
		val, ok := instrMap[strings.ToLower(instr.Name)]

		if !ok {
			valid = false
			fmt.Printf("vm Error: %s is not a valid instruction (index %d)\n", instr.Name, ind)
			continue
		}

		if len(instr.Args) != val.args {
			valid = false
			fmt.Printf("vm Error: %s does not have correct amount of arguments. expected=%d, got=%d, args=[%v] (index %d)", instr.Name, val.args, len(instr.Args), instr.Args, ind)
			continue
		}

		vm.Funcs[ind] = val.f
	}

	if !valid {
		os.Exit(2)
	}

	return vm
}

func (vm *VirtualMachine) read(offset int64) int64 {
	return vm.Values[vm.ValPointer+offset]
}

func (vm *VirtualMachine) getVal(arg *vminstruction.Arg) int64 {
	if arg.IsAbsolute {
		return arg.Value
	}

	return vm.read(arg.Value)
}

func (vm *VirtualMachine) write(offset int64, val int64) {
	vm.Values[vm.ValPointer+offset] = val
}

// Move DST SRC
func (vm *VirtualMachine) Move(args ...*vminstruction.Arg) {
	vm.write(args[0].Value, vm.getVal(args[1]))
}

// Set DST NUMBER
func (vm *VirtualMachine) Set(args ...*vminstruction.Arg) {
	vm.write(args[0].Value, vm.getVal(args[1]))
}

// Add DST SRC1 SRC2
func (vm *VirtualMachine) Add(args ...*vminstruction.Arg) {
	vm.write(args[0].Value, vm.getVal(args[1])+vm.getVal(args[2]))
}

// Multiply DST SRC1 SRC2
func (vm *VirtualMachine) Multiply(args ...*vminstruction.Arg) {
	vm.write(args[0].Value, vm.getVal(args[1])*vm.getVal(args[2]))
}

// Negate DST SRC
func (vm *VirtualMachine) Negate(args ...*vminstruction.Arg) {
	vm.write(args[0].Value, -vm.getVal(args[1]))
}

// Not DST SRC
func (vm *VirtualMachine) Not(args ...*vminstruction.Arg) {
	if vm.getVal(args[1]) != 0 {
		vm.write(args[0].Value, 0)
	} else {
		vm.write(args[0].Value, 1)
	}
}

// Jump INS
func (vm *VirtualMachine) Jump(args ...*vminstruction.Arg) {
	vm.ProgramCounter = vm.getVal(args[0])
}

// JumpIfZero SRC INS
func (vm *VirtualMachine) JumpIfZero(args ...*vminstruction.Arg) {
	if vm.getVal(args[0]) == 0 {
		vm.ProgramCounter = vm.getVal(args[1])
	}
}

// Call INS NUMBER DST
func (vm *VirtualMachine) Call(args ...*vminstruction.Arg) {
	vm.CallStack = append(vm.CallStack, call{pc: vm.ProgramCounter, vp: vm.ValPointer, dst: vm.getVal(args[2])})
	vm.ValPointer += vm.getVal(args[1])
	vm.ProgramCounter = vm.getVal(args[0])
}

// Return SRC
func (vm *VirtualMachine) Return(args ...*vminstruction.Arg) {
	res := vm.getVal(args[0])

	if len(vm.CallStack) == 0 {
		fmt.Println(res)
		os.Exit(0)
	}

	// last index
	li := len(vm.CallStack) - 1
	lastcall := vm.CallStack[li]
	vm.CallStack = vm.CallStack[:li]

	vm.ProgramCounter = lastcall.pc
	vm.ValPointer = lastcall.vp
	vm.write(lastcall.dst, res)
}

// LessThan DST SRC1 SRC2
func (vm *VirtualMachine) LessThan(args ...*vminstruction.Arg) {
	if vm.getVal(args[1]) < vm.getVal(args[2]) {
		vm.write(args[0].Value, 1)
	} else {
		vm.write(args[0].Value, 0)
	}
}

// Equals DST SRC1 SRC2
func (vm *VirtualMachine) Equals(args ...*vminstruction.Arg) {
	if vm.getVal(args[1]) == vm.getVal(args[2]) {
		vm.write(args[0].Value, 1)
	} else {
		vm.write(args[0].Value, 0)
	}
}
