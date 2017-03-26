package gameboy

/*
CPU Instruction structure. Has two durations for some instructions: action and noop.
When action noop is 0, the instruction always takes the action duration.
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

type instructionExecutor func(mem *memory, reg *register, instr *instruction) int

func createInstructionMap() *map[uint8]*instruction {
	return &map[uint8]*instruction{
		0x04: newInstruction("INC B", 1, 4, inc_b),
		0x05: newInstruction("DEC B", 1, 4, dec_b),
		0x06: newInstruction("LD B, d8", 2, 8, ld_b_d8),
		0x0c: newInstruction("INC C", 1, 4, inc_c),
		0x0d: newInstruction("DEC C", 1, 4, dec_c),
		0x0e: newInstruction("LD C", 2, 8, ld_c),
		0x11: newInstruction("LD DE,d16", 3, 12, ld_de_d16),
		0x13: newInstruction("INC DE", 1, 8, inc_de),
		0x15: newInstruction("DEC D", 1, 4, dec_d),
		0x16: newInstruction("LD D,d8", 2, 8, ld_d_d8),
		0x17: newInstruction("RLA", 1, 4, rla),
		0x18: newInstruction("JR r8", 2, 12, jr_r8),
		0x1a: newInstruction("LD A,(DE)", 1, 8, ld_a_DE),
		// HIER WAS IK
		0x1d: newInstruction("DEC E", 1, 4, dec_e),
		0x1e: newInstruction("LD E,d8", 2, 8, ld_e_d8),
		0x20: newConditionalInstruction("JR NZ,r8", 2, 12, 8, jr_nz),
		0x21: newInstruction("LD HL", 3, 12, ld_hl),
		0x22: newInstruction("LD (HL+),A", 1, 8, ld_HLP_a),
		0x23: newInstruction("INC HL", 1, 8, inc_hl),
		0x24: newInstruction("INC H", 1, 4, inc_h),
		0x28: newConditionalInstruction("JR Z,r8", 2, 12, 8, jr_z_r8),
		0x2e: newInstruction("LD L,d8", 2, 8, ld_l_d8),
		0x31: newInstruction("LD SP", 3, 12, ld_sp),
		0x32: newInstruction("LDD (HL-),A", 1, 8, ldd_HL_a),
		0x3d: newInstruction("DEC A", 1, 4, dec_a),
		0x3e: newInstruction("LD A", 2, 8, ld_a),
		0x4f: newInstruction("LD C,A", 1, 4, ld_c_a),
		0x57: newInstruction("LD D,A", 1, 4, ld_d_a),
		0x67: newInstruction("LD H,A", 1, 4, ld_h_a),
		0x77: newInstruction("LD (HL),A", 1, 8, ld_HL_a),
		0x7b: newInstruction("LD A,E", 1, 4, ld_a_e),
		0x7c: newInstruction("LD A,H", 1, 4, ld_a_h),
		0x90: newInstruction("SUB B", 1, 4, sub_b),
		0xaf: newInstruction("XOR A", 1, 4, xor_a),
		0xc1: newInstruction("POP BC", 1, 12, pop_bc),
		0xc5: newInstruction("PUSH BC", 1, 16, push_bc),
		0xc9: newInstruction("RET", 0, 16, ret), // RET no length: interference with updating PC
		0xcb: newInstruction("CB", 1, 4, nil),
		0xcd: newInstruction("CALL a16", 0, 12, call_nn), // CALL has no length, as it interferes with updating PC
		0xe0: newInstruction("LDH a8,A", 2, 12, ldh_a8_A),
		0xe2: newInstruction("LD (C),A", 1, 8, ld_C_a),
		0xea: newInstruction("LD (a16),A", 3, 16, ld_A16_A),
		0xf0: newInstruction("LDH A,(a8)", 2, 12, ld_a_a8),
		0xfe: newInstruction("CP d8", 2, 8, cp_d8),
	}
}

func newConditionalInstruction(name string, length int, actionDuration int, noopDuration int, fp func(mem *memory, reg *register, instr *instruction) int) *instruction {
	return &instruction{
		name:            name,
		bytes:           length,
		duration_action: actionDuration,
		duration_noop:   noopDuration,
		executor:        fp,
	}
}

func newInstruction(name string, length int, duration int, fp func(mem *memory, reg *register, instr *instruction) int) *instruction {
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

func inc_register(br *byteRegister, reg *register) {
	val := br.val()
	*br = byteRegister(val + 1)
	reg.Flag.Z = *br== 0
	reg.Flag.N = false
	reg.Flag.H = (val&0xf)+(1&0xf)&0x10 == 0x10
}

func dec_register(br *byteRegister, reg *register) {
	val := br.val()
	*br = byteRegister(val - uint8(1))
	reg.Flag.Z = val - 1 == 0
	reg.Flag.N = true
	reg.Flag.H = (val&0xf0)-(1&0xf0)&0x8 == 0x8
}

func sub_register(reg *register, val uint8) {
	before := reg.A.val()
	reg.A = byteRegister(reg.A.val() - val)
	reg.Flag.Z = reg.A == 0
	reg.Flag.N = true
	reg.Flag.H = before < val
}

func call_nn(mem *memory, reg *register, instr *instruction) int {
	pushStack16(mem, reg, reg.PC.val()+uint16(3))
	reg.PC = halfWordRegister(readArgHalfword(mem, reg, 1))
	// Does not affect flags
	return instr.duration_action
}

func jr_nz(mem *memory, reg *register, instr *instruction) int {
	n := readArgByte(mem, reg, 1)
	if !reg.Flag.Z {
		reg.PC = halfWordRegister(int(reg.PC.val()) + int(int8(n)))
		return instr.duration_action
	}
	// Does not affect flags
	return instr.duration_noop
}

func inc_b(mem *memory, reg *register, instr *instruction) int {
	inc_register(&reg.B, reg)
	return instr.duration_action
}

func inc_c(mem *memory, reg *register, instr *instruction) int {
	inc_register(&reg.C, reg)
	return instr.duration_action
}

func ld_a(mem *memory, reg *register, instr *instruction) int {
	reg.A = byteRegister(readArgByte(mem, reg, 1))
	return instr.duration_action
}

func ld_a_DE(mem *memory, reg *register, instr *instruction) int {
	val := mem.read8(reg.readDuo(reg_de))
	reg.A = byteRegister(val)
	// Does not affect flags
	return instr.duration_action
}


func ld_c(mem *memory, reg *register, instr *instruction) int {
	reg.C = byteRegister(readArgByte(mem, reg, 1))
	// Does not affect flags
	return instr.duration_action
}

func ld_C_a(mem *memory, reg *register, instr *instruction) int {
	mem.write8(0xFF00+uint16(reg.C.val()), reg.A.val())
	// Does not affect flags
	return instr.duration_action
}

func ld_de_d16(mem *memory, reg *register, instr *instruction) int {
	reg.writeDuo(reg_de, readArgHalfword(mem, reg, 1))
	// Does not affect flags
	return instr.duration_action
}

func ld_sp(mem *memory, reg *register, instr *instruction) int {
	arg := readArgHalfword(mem, reg, 1)
	reg.SP = halfWordRegister(arg)
	// Does not affect flags
	return instr.duration_action
}

func ld_hl(mem *memory, reg *register, instr *instruction) int {
	reg.writeDuo(reg_hl, readArgHalfword(mem, reg, 1))
	// Does not affect flags
	return instr.duration_action
}

func ld_HL_a(mem *memory, reg *register, instr *instruction) int {
	mem.write8(reg.readDuo(reg_hl), reg.A.val())
	// Does not affect flags
	return instr.duration_action
}

func ldd_HL_a(mem *memory, reg *register, instr *instruction) int {
	mem.write8(reg.readDuo(reg_hl), reg.A.val())
	reg.decrDuo(reg_hl)
	// Does not affect flags
	return instr.duration_action
}

func ldh_a8_A(mem *memory, reg *register, instr *instruction) int {
	mem.write8(0xff00+uint16(readArgByte(mem, reg, 1)), reg.A.val())
	// Does not affect flags
	return instr.duration_action
}

func xor_a(mem *memory, reg *register, instr *instruction) int {
	reg.A = reg.A ^ reg.A
	if reg.A == 0 {
		reg.Flag.Z = true
	}
	reg.Flag.N = false
	reg.Flag.H = false
	reg.Flag.C = false
	return instr.duration_action
}

func push_bc(mem *memory, reg *register, instr *instruction) int {
	pushStack16(mem, reg, reg.readDuo(reg_bc))
	// Does not affect flags
	return instr.duration_action
}

func rla(mem *memory, reg *register, instr *instruction) int {
	isCarrySet := reg.Flag.C
	isMSBSet := testBit(reg.A.val(), 7)

	reg.Flag.Z = false
	reg.Flag.N = false
	reg.Flag.H = false
	reg.Flag.C = false

	reg.A = byteRegister(reg.A.val() << 1)

	if isMSBSet {
		reg.Flag.C = true
	}
	if isCarrySet {
		reg.A = byteRegister(setBit(reg.A.val(), 0))
	}

	if reg.A.val() == 0 {
		reg.Flag.Z = true
	}
	return instr.duration_action
}

func pop_bc(mem *memory, reg *register, instr *instruction) int {
	val := popStack16(mem, reg)
	reg.writeDuo(reg_bc, val)
	// Does not affect flags
	return instr.duration_action
}

func dec_b(mem *memory, reg *register, instr *instruction) int {
	dec_register(&reg.B, reg)
	return instr.duration_action
}

func ld_c_a(mem *memory, reg *register, instr *instruction) int {
	reg.C = reg.A
	return instr.duration_action
}

func ld_b_d8(mem *memory, reg *register, instr *instruction) int {
	arg := readArgByte(mem, reg, 1)
	reg.B = byteRegister(arg)
	return instr.duration_action
}

func ld_HLP_a(mem *memory, reg *register, instr *instruction) int {
	mem.write8(reg.readDuo(reg_hl), reg.A.val())
	reg.incrDuo(reg_hl)
	return instr.duration_action
}

func inc_hl(mem *memory, reg *register, instr *instruction) int {
	reg.incrDuo(reg_hl)
	return instr.duration_action
}

func inc_de(mem *memory, reg *register, instr *instruction) int {
	reg.incrDuo(reg_de)
	return instr.duration_action
}

func ret(mem *memory, reg *register, instr *instruction) int {
	addr := popStack16(mem, reg)
	reg.PC = halfWordRegister(addr)
	return instr.duration_action
}

func ld_a_e(mem *memory, reg *register, instr *instruction) int {
	reg.A = reg.E
	return instr.duration_action
}

func cp_d8(mem *memory, reg *register, instr *instruction) int {
	var arg uint8 = readArgByte(mem, reg, 1)
	reg.Flag.Z = reg.A.val() == arg
	reg.Flag.N = true
	reg.Flag.H = false // TODO: Check borrow
	reg.Flag.C = reg.A.val() < arg
	return instr.duration_action
}

func ld_A16_A(mem *memory, reg *register, instr *instruction) int {
	mem.write8(readArgHalfword(mem, reg, 1), reg.A.val())
	return instr.duration_action
}

func dec_a(mem *memory, reg *register, instr *instruction) int {
	dec_register(&reg.A, reg)
	return instr.duration_action
}

func jr_z_r8(mem *memory, reg *register, instr *instruction) int {
	if reg.Flag.Z {
		reg.PC = halfWordRegister(reg.PC.val() + uint16(readArgByte(mem, reg, 1)))
		return instr.duration_action
	}
	return instr.duration_noop
}

func dec_c(mem *memory, reg *register, instr *instruction) int {
	dec_register(&reg.C, reg)
	return instr.duration_action
}

func ld_l_d8(mem *memory, reg *register, instr *instruction) int {
	reg.L = byteRegister(readArgByte(mem, reg, 1))
	return instr.duration_action
}

func jr_r8(mem *memory, reg *register, instr *instruction) int {
	old_pc := reg.PC.val()
	arg := int8(readArgByte(mem, reg, 1))
	reg.PC = halfWordRegister(uint16(int(old_pc) + int(arg)))
	return instr.duration_action
}

func ld_h_a(mem *memory, reg *register, instr *instruction) int {
	reg.H = reg.A
	return instr.duration_action
}

func ld_d_a(mem *memory, reg *register, instr *instruction) int {
	reg.D = reg.A
	return instr.duration_action
}

func ld_e_d8(mem *memory, reg *register, instr *instruction) int {
	d8 := readArgByte(mem, reg, 1)
	reg.E = byteRegister(d8)
	return instr.duration_action
}

func ld_a_a8(mem *memory, reg *register, instr *instruction) int {
	address := 0xFF00 + uint16(readArgByte(mem, reg, 1))
	reg.A = byteRegister(mem.read8(address))
	return instr.duration_action
}

func dec_e(mem *memory, reg *register, instr *instruction) int {
	dec_register(&reg.E, reg)
	return instr.duration_action
}

func inc_h(mem *memory, reg *register, instr *instruction) int {
	inc_register(&reg.H, reg)
	return instr.duration_action
}

func ld_a_h(mem *memory, reg *register, instr *instruction) int {
	reg.A = reg.H
	return instr.duration_action
}

func sub_b(mem *memory, reg *register, instr *instruction) int {
	sub_register(reg, reg.B.val())
	return instr.duration_action
}

func dec_d(mem *memory, reg *register, instr *instruction) int {
	dec_register(&reg.D, reg)
	return instr.duration_action
}

func ld_d_d8(mem *memory, reg *register, instr *instruction) int {
	reg.D = byteRegister(readArgByte(mem, reg, 1))
	return instr.duration_action
}
