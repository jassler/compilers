package vm

import "github.com/simplang/vminstruction"
import "strings"
import "fmt"
import "os"

type VirtualMachine struct {
	Instructions   []*vminstruction.Instruction
	Funcs          []func(...int64)
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
		Funcs:          make([]func(...int64), len(instr)),
		ProgramCounter: 0,
		CallStack:      []call{},
		Values:         make([]int64, 1000000),
		ValPointer:     0,
	}

	type tuple struct {
		f    func(...int64)
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

func (vm *VirtualMachine) write(offset int64, val int64) {
	vm.Values[vm.ValPointer+offset] = val
}

// Move DST SRC
func (vm *VirtualMachine) Move(args ...int64) {
	vm.write(args[0], vm.read(args[1]))
}

// Set DST NUMBER
func (vm *VirtualMachine) Set(args ...int64) {
	vm.write(args[0], args[1])
}

// Add DST SRC1 SRC2
func (vm *VirtualMachine) Add(args ...int64) {
	vm.write(args[0], vm.read(args[1])+vm.read(args[2]))
}

// Multiply DST SRC1 SRC2
func (vm *VirtualMachine) Multiply(args ...int64) {
	vm.write(args[0], vm.read(args[1])+vm.read(args[2]))
}

// Negate DST SRC
func (vm *VirtualMachine) Negate(args ...int64) {
	vm.write(args[0], -vm.read(args[1]))
}

// Not DST SRC
func (vm *VirtualMachine) Not(args ...int64) {
	if vm.read(args[1]) != 0 {
		vm.write(args[0], 0)
	} else {
		vm.write(args[0], 1)
	}
}

// Jump INS
func (vm *VirtualMachine) Jump(args ...int64) {
	vm.ProgramCounter = args[0]
}

// JumpIfZero SRC INS
func (vm *VirtualMachine) JumpIfZero(args ...int64) {
	if vm.read(args[0]) == 0 {
		vm.ProgramCounter = args[1]
	}
}

// Call INS NUMBER DST
func (vm *VirtualMachine) Call(args ...int64) {
	vm.CallStack = append(vm.CallStack, call{pc: vm.ProgramCounter, vp: vm.ValPointer, dst: args[2]})
	vm.ValPointer += args[1]
	vm.ProgramCounter = args[0]
}

// Return SRC
func (vm *VirtualMachine) Return(args ...int64) {
	res := vm.read(args[0])

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
func (vm *VirtualMachine) LessThan(args ...int64) {
	if args[1] < args[2] {
		vm.write(args[0], 1)
	} else {
		vm.write(args[0], 0)
	}
}

// Equals DST SRC1 SRC2
func (vm *VirtualMachine) Equals(args ...int64) {
	if args[1] == args[2] {
		vm.write(args[0], 1)
	} else {
		vm.write(args[0], 0)
	}
}
