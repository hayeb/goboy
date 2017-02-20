package gameboy

import "fmt"

/*
CPU Instruction structure. Has two durations for some instructions: action and noop.
When action noop is 0, the instruction always takes the action duration
*/
type instruction struct {
	name string
	// The length of the instruction in bytes, including the opcode
	bytes int

	// Duration of the instruction when an action is taken.
	duration_action int

	// Duration of the instruction when no action is taken.
	duration_noop int

	// Pointer to the function to be executed
	executor instructionExecutor
}

type instructionExecutor func(mem *memory, reg *register)

func createInstructionMap() *map[uint8]instruction {
	return &map[uint8]instruction{
		0xc:  newInstruction("INC_C", 1, 4, inc_c),
		0xe:  newInstruction("LD_C", 2, 8, ld_c),
		0x11: newInstruction("LD_DE_d16", 3, 12, ld_de_d16),
		0x17: newInstruction("RLA", 1, 4, rla),
		0x1a: newInstruction("LD_A_(DE)", 1, 8, ld_a_DE),
		0x20: newConditionalInstruction("JR", 2, 12, 8, jr),
		0x21: newInstruction("LD_HL", 3, 12, ld_hl),
		0x31: newInstruction("LD_SP", 3, 12, ld_sp),
		0x32: newInstruction("LDD_(HL)_A", 1, 8, ldd_HL_a),
		0x3e: newInstruction("LD_A", 2, 8, ld_a),
		0x77: newInstruction("LD_(HL)_a", 1, 8, ld_HL_a),
		0xAF: newInstruction("XOR_A", 1, 4, xor_a),
		0xC1: newInstruction("POP_BC", 1, 12, pop_bc),
		0xC5: newInstruction("PUSH_BC", 1, 16, push_bc),
		0xCB: newInstruction("CB", 1, 4, nil),
		0xCD: newInstruction("CALL_nn", 3, 12, call_nn),
		0xE0: newInstruction("LDH_a8_A", 2, 12, ldh_a8_A),
		0xE2: newInstruction("LD_(C)_A", 1, 8, ld_C_a),
	}
}

func newConditionalInstruction(name string, length int, actionDuration int, noopDuration int, fp func(mem *memory, reg *register)) instruction {
	return instruction{
		name:            name,
		bytes:           length,
		duration_action: actionDuration,
		duration_noop:   noopDuration,
		executor:        fp,
	}
}

func newInstruction(name string, length int, duration int, fp func(mem *memory, reg *register)) instruction {
	return newConditionalInstruction(name, length, duration, 0, fp)
}

// Returns a uint8 with the 8 least signigicant bits of i
func leastSig16(i uint16) uint8 {
	return uint8(i & ((1 << 8) - 1))
}

// Returns a uint8 with the 8 most signigicant bits of i
func mostSig16(i uint16) uint8 {
	return uint8(i >> 8)
}

func call_nn(mem *memory, reg *register) {
	left := mostSig16(reg.PC.val() + uint16(3))
	right := leastSig16((reg.PC.val() + uint16(3)))
	pushStack8(mem, reg, left)
	pushStack8(mem, reg, right)

	reg.PC = halfWordRegister(readArgHalfword(mem, reg, 1))
}

func jr(mem *memory, reg *register) {
	n := int8(readArgByte(mem, reg, 1))
	if !reg.Flag.Z {
		reg.PC = halfWordRegister(int(reg.PC.val()) + int(n))
	}
}

func inc_c(mem *memory, reg *register) {
	val := reg.C
	incRegister8(&reg.C)
	reg.Flag.Z = reg.C == 0
	reg.Flag.N = false
	// TODO: hacky half-carry. Do differently?
	reg.Flag.H = val < 8 && reg.C >= 16
}

func ld_a(mem *memory, reg *register) {
	reg.A = byteRegister(readArgByte(mem, reg, 1))
}

func ld_a_DE(mem *memory, reg *register) {
	reg.A = byteRegister(mem.read8(reg.readDuo(reg_de)))

}

func ld_c(mem *memory, reg *register) {
	reg.C = byteRegister(readArgByte(mem, reg, 1))
}

func ld_C_a(mem *memory, reg *register) {
	mem.write8(0xFF00+uint16(reg.C.val()), reg.A.val())
}

func ld_de_d16(mem *memory, reg *register) {
	reg.writeDuo(reg_de, readArgHalfword(mem, reg, 1))
}

func ld_sp(mem *memory, reg *register) {
	arg := readArgHalfword(mem, reg, 1)
	fmt.Printf("SP: %#04x\n", arg)
	reg.SP = halfWordRegister(arg)
}

func ld_hl(mem *memory, reg *register) {
	reg.writeDuo(reg_hl, readArgHalfword(mem, reg, 1))
}

func ld_HL_a(mem *memory, reg *register) {
	mem.write8(mem.read16(reg.readDuo(reg_hl)), readArgByte(mem, reg, 1))
}

func ldd_HL_a(mem *memory, reg *register) {
	mem.write8(reg.readDuo(reg_hl), reg.A.val())
	reg.decrDuo(reg_hl)
}

func ldh_a8_A(mem *memory, reg *register) {
	reg.A = byteRegister(mem.read8(uint16(readArgByte(mem, reg, 1)) + 0xff00))
}

func xor_a(mem *memory, reg *register) {
	reg.A = reg.A ^ reg.A
	if reg.A == 0 {
		reg.Flag.Z = true
	}
}

func push_bc(mem *memory, reg *register) {
	pushStack16(mem, reg, reg.readDuo(reg_bc))
}

func rla(mem *memory, reg *register) {
	r, c := rLeft(reg.A.val())
	reg.A = byteRegister(r)
	reg.Flag.Z = reg.A == 0
	reg.Flag.N = false
	reg.Flag.H = false
	reg.Flag.C = c
}

func pop_bc(mem *memory, reg *register) {
	reg.writeDuo(reg_bc, popStack16(mem, reg))
}
