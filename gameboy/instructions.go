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
	durationAction int

	// Duration of the instruction when no action is taken.
	durationNoop int

	// Pointer to the function to be executed
	executor instructionExecutor
}

type instructionExecutor func(mem *memory, reg *register, instr *instruction) int

func createInstructionMap() *map[uint8]*instruction {
	return &map[uint8]*instruction{
		0x00: newInstruction("NOOP", 1, 4, noop),
		0x01: newInstruction("LD BC,d16", 3, 12, ldBCnn),
		0x03: newInstruction("INC BC", 1, 8, incBC),
		0x04: newInstruction("INC B", 1, 4, incB),
		0x05: newInstruction("DEC B", 1, 4, decB),
		0x06: newInstruction("LD B, d8", 2, 8, ldBD8),
		0x07: newInstruction("RLCA", 1, 4, rlca),
		0x09: newInstruction("ADD HL,BC", 1, 8, addHlBc),
		0x0a: newInstruction("LD A,(BC)", 1, 8, ldABC),
		0x0b: newInstruction("DEC BC", 1, 8, decBC),
		0x0c: newInstruction("INC C", 1, 4, incC),
		0x0d: newInstruction("DEC C", 1, 4, decC),
		0x0e: newInstruction("LD C", 2, 8, ldC),
		0x11: newInstruction("LD DE,d16", 3, 12, ldDeD16),
		0x12: newInstruction("LD (DE),A", 1, 8, ldDeA),
		0x13: newInstruction("INC DE", 1, 8, incDE),
		0x15: newInstruction("DEC D", 1, 4, decD),
		0x16: newInstruction("LD D,d8", 2, 8, ldDD8),
		0x17: newInstruction("RLA", 1, 4, rla),
		0x18: newInstruction("JR r8", 2, 12, jrR8),
		0x19: newInstruction("ADD HL,BE", 1, 8, addHlDe),
		0x1a: newInstruction("LD A,(DE)", 1, 8, ldADE),
		0x1c: newInstruction("INC E", 1, 8, incE),
		0x1d: newInstruction("DEC E", 1, 4, decE),
		0x1e: newInstruction("LD E,d8", 2, 8, ldED8),
		0x20: newConditionalInstruction("JR NZ,r8", 2, 12, 8, jrNz),
		0x21: newInstruction("LD HL,d16", 3, 12, ldHl),
		0x22: newInstruction("LD (HL+),A", 1, 8, ldHLPA),
		0x23: newInstruction("INC HL", 1, 8, incHl),
		0x24: newInstruction("INC H", 1, 4, incH),
		0x28: newConditionalInstruction("JR Z,r8", 2, 12, 8, jrZR8),
		0x2a: newInstruction("LD A,(HL+)", 1, 8, ldAHLP),
		0x2c: newInstruction("INC L", 1, 8, incL),
		0x2d: newInstruction("DEC L", 1, 4, decL),
		0x2e: newInstruction("LD L,d8", 2, 8, ldLD8),
		0x2f: newInstruction("CPL", 1, 4, cpl),
		0x31: newInstruction("LD SP", 3, 12, ldSp),
		0x32: newInstruction("LDD (HL-),A", 1, 8, lddHLA),
		0x34: newInstruction("INC (HL)", 1, 12, incHL),
		0x35: newInstruction("DEC (HL)", 1, 12, decHl),
		0x36: newInstruction("LD (HL),n", 2, 12, lddHLn),
		0x3c: newInstruction("INC A", 1, 4, incA),
		0x3d: newInstruction("DEC A", 1, 4, decA),
		0x3e: newInstruction("LD A", 2, 8, ldA),
		0x40: newInstruction("LD B,B", 1, 4, ldBB),
		0x46: newInstruction("LD B,(HL)", 1, 8, ldBHL),
		0x47: newInstruction("LD B,A", 1, 4, ldBA),
		0x4f: newInstruction("LD C,A", 1, 4, ldcA),
		0x4e: newInstruction("LD C,(HL)", 1, 8, ldCHL),
		0x50: newInstruction("LD D,B", 1, 4, ldDB),
		0x56: newInstruction("LD D,(HL)", 1, 8, ldDHL),
		0x57: newInstruction("LD D,A", 1, 4, ldDA),
		0x5e: newInstruction("LD E,(HL)", 1, 8, ldEHL),
		0x5f: newInstruction("LD E,A", 1, 4, ldEA),
		0x60: newInstruction("LD H,B", 1, 4, ldHB),
		0x67: newInstruction("LD H,A", 1, 4, ldHA),
		0x69: newInstruction("LD E,C", 1, 4, ldEC),
		0x6f: newInstruction("LD L,A", 1, 4, ldLA),
		0x77: newInstruction("LD (HL),A", 1, 8, ldHLA),
		0x78: newInstruction("LD A,B", 1, 4, ldAB),
		0x79: newInstruction("LD A, C", 1, 4, ldAC),
		0x7b: newInstruction("LD A,E", 1, 4, ldAE),
		0x7c: newInstruction("LD A,H", 1, 4, ldAH),
		0x7d: newInstruction("LD A,L", 1, 4, ldAL),
		0x7e: newInstruction("LD A,(HL)", 1, 8, ldAHL),
		0x7f: newInstruction("LD A,A", 1, 4, ldAA),
		0x80: newInstruction("ADD A,B", 1, 4, addAB),
		0x85: newInstruction("ADD A,L", 1, 4, addAL),
		0x86: newInstruction("ADD A,(HL)", 1, 8, addAHL),
		0x87: newInstruction("ADD A,A", 1, 8, addAA),
		0x89: newInstruction("ADC A,C", 1, 4, adcAC),
		0x90: newInstruction("SUB B", 1, 4, subB),
		0xa1: newInstruction("AND C", 1, 4, andC),
		0xa7: newInstruction("AND A", 1, 4, andA),
		0xa9: newInstruction("XOR C", 1, 4, xorC),
		0xaf: newInstruction("XOR A", 1, 4, xorA),
		0xb0: newInstruction("OR B", 1, 4, orB),
		0xb1: newInstruction("OR C", 1, 4, orC),
		0xbe: newInstruction("CP (HL)", 1, 8, cpHL),
		0xc0: newConditionalInstruction("RET NZ", 1, 20, 8, retNz),
		0xc1: newInstruction("POP BC", 1, 12, popBc),
		0xc3: newInstruction("JP a16", 3, 16, jpnn),
		0xc8: newConditionalInstruction("RET Z", 1, 20, 8, retZ),
		0xc5: newInstruction("PUSH BC", 1, 16, pushBc),
		0xc7: newInstruction("RST 00", 1, 16, rst00),
		0xc9: newInstruction("RET", 1, 16, ret),
		0xca: newConditionalInstruction("JP Z,a16", 3, 16, 12, jpZa16),
		0xcb: newInstruction("CB", 1, 4, nil),
		0xcc: newConditionalInstruction("CALL Z,a16", 3, 24, 12, callZa16),
		0xcd: newInstruction("CALL a16", 3, 12, callNn),
		0xd1: newInstruction("POP DE", 1, 12, popDe),
		0xd5: newInstruction("PUSH DE", 1, 16, pushDe),
		0xd6: newInstruction("SUB d8", 2, 8, subD8),
		0xd9: newInstruction("RETI", 1, 16, reti),
		0xe0: newInstruction("LDH a8,A", 2, 12, ldhA8A),
		0xe1: newInstruction("POP HL", 1, 12, popHl),
		0xe2: newInstruction("LD (C),A", 1, 8, ldCA),
		0xe5: newInstruction("PUSH HL", 1, 8, pushHl),
		0xe6: newInstruction("AND d8", 2, 8, andd8),
		0xe9: newInstruction("JP (HL)", 1, 8, jphl),
		0xea: newInstruction("LD (a16),A", 3, 16, ldA16A),
		0xee: newInstruction("XOR d8", 2, 8, xord8),
		0xef: newInstruction("RST 28H", 1, 16, rst28),
		0xf0: newInstruction("LDH A,(a8)", 2, 12, ldAA8),
		0xf1: newInstruction("POP AF", 1, 12, popAf),
		0xf3: newInstruction("DI", 1, 4, di),
		0xf5: newInstruction("PUSH AF", 1, 16, pushAf),
		0xfa: newInstruction("LD A,(a16)", 3, 16, ldAa16),
		0xfb: newInstruction("EI", 1, 4, ei),
		0xfe: newInstruction("CP d8", 2, 8, cpD8),
	}
}

func newConditionalInstruction(name string, length int, actionDuration int, noopDuration int, fp func(mem *memory, reg *register, instr *instruction) int) *instruction {
	return &instruction{
		name:           name,
		bytes:          length,
		durationAction: actionDuration,
		durationNoop:   noopDuration,
		executor:       fp,
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

func incRegister(br *uint8, reg *register) {
	*br += uint8(1)
	reg.Flag.Z = *br == 0
	reg.Flag.N = false
	reg.Flag.H = *br > 0xf
}

func decRegister(br *uint8, reg *register) {
	*br -= uint8(1)
	reg.Flag.Z = *br == 0
	reg.Flag.N = true
	reg.Flag.H = *br < 0xf
}

func subRegister(reg *register, val uint8) {
	before := reg.A
	reg.A = reg.A - val
	reg.Flag.Z = reg.A == 0
	reg.Flag.N = true
	reg.Flag.H = before < val
}

func callNn(mem *memory, reg *register, instr *instruction) int {
	pushStack16(mem, reg, reg.PC+uint16(instr.bytes))
	reg.PC = readArgHalfword(mem, reg, 1)
	return instr.durationAction
}

func jrNz(mem *memory, reg *register, instr *instruction) int {
	var n = int8(readArgByte(mem, reg, 1))
	reg.incPC(instr.bytes)
	if !reg.Flag.Z {
		reg.PC = uint16(int(reg.PC) + int(n))
		return instr.durationAction
	}
	return instr.durationNoop
}

func jpZa16(mem *memory, reg *register, instr *instruction) int {
	if reg.Flag.Z {
		reg.PC = readArgHalfword(mem, reg, 1)
		return instr.durationAction
	}
	reg.incPC(instr.bytes)
	return instr.durationNoop
}

func incB(mem *memory, reg *register, instr *instruction) int {
	incRegister(&reg.B, reg)
	reg.incPC(instr.bytes)
	return instr.durationAction
}

func incC(mem *memory, reg *register, instr *instruction) int {
	incRegister(&reg.C, reg)
	reg.incPC(instr.bytes)
	return instr.durationAction
}

func ldA(mem *memory, reg *register, instr *instruction) int {
	reg.A = readArgByte(mem, reg, 1)
	reg.incPC(instr.bytes)
	return instr.durationAction
}

func ldADE(mem *memory, reg *register, instr *instruction) int {
	reg.A = mem.read8(reg.readDuo(REG_DE))
	// Does not affect flags
	reg.incPC(instr.bytes)
	return instr.durationAction
}

func ldC(mem *memory, reg *register, instr *instruction) int {
	reg.C = readArgByte(mem, reg, 1)
	// Does not affect flags
	reg.incPC(instr.bytes)
	return instr.durationAction
}

func ldCA(mem *memory, reg *register, instr *instruction) int {
	mem.write8(0xFF00+uint16(reg.C), reg.A)
	// Does not affect flags
	reg.incPC(instr.bytes)
	return instr.durationAction
}

func ldDeD16(mem *memory, reg *register, instr *instruction) int {
	reg.writeDuo(REG_DE, readArgHalfword(mem, reg, 1))
	// Does not affect flags
	reg.incPC(instr.bytes)
	return instr.durationAction
}

func ldSp(mem *memory, reg *register, instr *instruction) int {
	reg.SP = readArgHalfword(mem, reg, 1)
	// Does not affect flags
	reg.incPC(instr.bytes)
	return instr.durationAction
}

func ldHl(mem *memory, reg *register, instr *instruction) int {
	reg.writeDuo(REG_HL, readArgHalfword(mem, reg, 1))
	// Does not affect flags
	reg.incPC(instr.bytes)
	return instr.durationAction
}

func ldHLA(mem *memory, reg *register, instr *instruction) int {
	mem.write8(reg.readDuo(REG_HL), reg.A)
	// Does not affect flags
	reg.incPC(instr.bytes)
	return instr.durationAction
}

func lddHLA(mem *memory, reg *register, instr *instruction) int {
	mem.write8(reg.readDuo(REG_HL), reg.A)
	reg.decrDuo(REG_HL)
	// Does not affect flags
	reg.incPC(instr.bytes)
	return instr.durationAction
}

func ldhA8A(mem *memory, reg *register, instr *instruction) int {
	arg := readArgByte(mem, reg, 1)
	mem.write8(0xFF00+uint16(arg), reg.A)
	// Does not affect flags
	reg.incPC(instr.bytes)
	return instr.durationAction
}

func xorA(mem *memory, reg *register, instr *instruction) int {
	reg.A ^= reg.A

	reg.Flag.Z = reg.A == 0
	reg.Flag.N = false
	reg.Flag.H = false
	reg.Flag.C = false

	reg.incPC(instr.bytes)
	return instr.durationAction
}

func xord8(mem *memory, reg *register, instr *instruction) int {
	reg.A ^= readArgByte(mem, reg, 1)

	reg.Flag.Z = reg.A == 0
	reg.Flag.N = false
	reg.Flag.H = false
	reg.Flag.C = false

	reg.incPC(instr.bytes)
	return instr.durationAction
}

func pushBc(mem *memory, reg *register, instr *instruction) int {
	pushStack16(mem, reg, reg.readDuo(REG_BC))
	// Does not affect flags
	reg.incPC(instr.bytes)
	return instr.durationAction
}

func rla(mem *memory, reg *register, instr *instruction) int {
	carry := reg.Flag.C
	reg.Flag.C = (reg.A&0x80)>>7 == 1

	reg.A = reg.A << 1
	if carry {
		reg.A = reg.A | 0x1
	}
	reg.Flag.Z = reg.A == 0
	reg.incPC(instr.bytes)
	return instr.durationAction
}

func rlca(mem *memory, reg *register, instr *instruction) int {
	reg.Flag.C = testBit(reg.A, 7)
	reg.A <<= 1

	reg.Flag.Z = reg.A == 0
	reg.Flag.N = false
	reg.Flag.H = false
	reg.incPC(instr.bytes)
	return instr.durationAction
}

func popBc(mem *memory, reg *register, instr *instruction) int {
	val := popStack16(mem, reg)
	reg.writeDuo(REG_BC, val)
	// Does not affect flags
	reg.incPC(instr.bytes)
	return instr.durationAction
}

func decB(mem *memory, reg *register, instr *instruction) int {
	decRegister(&reg.B, reg)
	reg.incPC(instr.bytes)
	return instr.durationAction
}

func ldcA(mem *memory, reg *register, instr *instruction) int {
	reg.C = reg.A
	reg.incPC(instr.bytes)
	return instr.durationAction
}

func ldBD8(mem *memory, reg *register, instr *instruction) int {
	reg.B = readArgByte(mem, reg, 1)
	reg.incPC(instr.bytes)
	return instr.durationAction
}

func ldHLPA(mem *memory, reg *register, instr *instruction) int {
	mem.write8(reg.readDuo(REG_HL), reg.A)
	reg.incrDuo(REG_HL)
	reg.incPC(instr.bytes)
	return instr.durationAction
}

func incHl(mem *memory, reg *register, instr *instruction) int {
	reg.incrDuo(REG_HL)
	reg.incPC(instr.bytes)
	return instr.durationAction
}

func decHl(mem *memory, reg *register, instr *instruction) int {
	addr := reg.readDuo(REG_HL)
	val := mem.read8(addr)
	mem.write8(addr, val+1)

	reg.Flag.Z = reg.readDuo(REG_HL) == 0
	reg.Flag.N = true
	reg.Flag.H = val == 0xf

	reg.incPC(instr.bytes)
	return instr.durationAction
}

func incDE(mem *memory, reg *register, instr *instruction) int {
	reg.incrDuo(REG_DE)
	reg.incPC(instr.bytes)
	return instr.durationAction
}

func incBC(mem *memory, reg *register, instr *instruction) int {
	reg.incrDuo(REG_BC)
	reg.incPC(instr.bytes)
	return instr.durationAction
}

func ret(mem *memory, reg *register, instr *instruction) int {
	reg.PC = popStack16(mem, reg)
	return instr.durationAction
}

func ldAE(mem *memory, reg *register, instr *instruction) int {
	reg.A = reg.E
	reg.incPC(instr.bytes)
	return instr.durationAction
}

func cpD8(mem *memory, reg *register, instr *instruction) int {
	var arg = readArgByte(mem, reg, 1)
	reg.Flag.Z = reg.A == arg
	reg.Flag.N = true
	reg.Flag.H = true
	reg.Flag.C = reg.A < arg
	reg.incPC(instr.bytes)
	return instr.durationAction
}

func ldA16A(mem *memory, reg *register, instr *instruction) int {
	mem.write8(readArgHalfword(mem, reg, 1), reg.A)
	reg.incPC(instr.bytes)
	return instr.durationAction
}

func decA(mem *memory, reg *register, instr *instruction) int {
	decRegister(&reg.A, reg)
	reg.incPC(instr.bytes)
	return instr.durationAction
}

func decL(mem *memory, reg *register, instr *instruction) int {
	decRegister(&reg.L, reg)
	reg.incPC(instr.bytes)
	return instr.durationAction
}

func jrZR8(mem *memory, reg *register, instr *instruction) int {
	arg := int8(readArgByte(mem, reg, 1))
	reg.incPC(instr.bytes)
	if reg.Flag.Z {
		reg.PC = uint16(int(reg.PC) + int(arg))
		return instr.durationAction
	}
	return instr.durationNoop
}

func decC(mem *memory, reg *register, instr *instruction) int {
	decRegister(&reg.C, reg)
	reg.incPC(instr.bytes)
	return instr.durationAction
}

func ldLD8(mem *memory, reg *register, instr *instruction) int {
	reg.L = readArgByte(mem, reg, 1)
	reg.incPC(instr.bytes)
	return instr.durationAction
}

func jrR8(mem *memory, reg *register, instr *instruction) int {
	oldPc := reg.PC
	arg := int8(readArgByte(mem, reg, 1))
	reg.PC = uint16(int(oldPc) + int(arg))

	reg.incPC(instr.bytes)
	return instr.durationAction
}

func ldHA(mem *memory, reg *register, instr *instruction) int {
	reg.H = reg.A
	reg.incPC(instr.bytes)
	return instr.durationAction
}

func ldDA(mem *memory, reg *register, instr *instruction) int {
	reg.D = reg.A
	reg.incPC(instr.bytes)
	return instr.durationAction
}

func ldDB(mem *memory, reg *register, instr *instruction) int {
	reg.D = reg.B
	reg.incPC(instr.bytes)
	return instr.durationAction
}

func ldEA(mem *memory, reg *register, instr *instruction) int {
	reg.E = reg.A
	reg.incPC(instr.bytes)
	return instr.durationAction
}

func ldED8(mem *memory, reg *register, instr *instruction) int {
	reg.E = readArgByte(mem, reg, 1)
	reg.incPC(instr.bytes)
	return instr.durationAction
}

func ldAA8(mem *memory, reg *register, instr *instruction) int {
	reg.A = mem.read8(0xFF00 + uint16(readArgByte(mem, reg, 1)))
	reg.incPC(instr.bytes)
	return instr.durationAction
}

func decE(mem *memory, reg *register, instr *instruction) int {
	decRegister(&reg.E, reg)
	reg.incPC(instr.bytes)
	return instr.durationAction
}

func incH(mem *memory, reg *register, instr *instruction) int {
	incRegister(&reg.H, reg)
	reg.incPC(instr.bytes)
	return instr.durationAction
}

func ldAH(mem *memory, reg *register, instr *instruction) int {
	reg.A = reg.H
	reg.incPC(instr.bytes)
	return instr.durationAction
}

func subB(mem *memory, reg *register, instr *instruction) int {
	subRegister(reg, reg.B)
	reg.incPC(instr.bytes)
	return instr.durationAction
}

func subD8(mem *memory, reg *register, instr *instruction) int {
	subRegister(reg, readArgByte(mem, reg, 1))
	reg.incPC(instr.bytes)
	return instr.durationAction
}

func decD(mem *memory, reg *register, instr *instruction) int {
	decRegister(&reg.D, reg)
	reg.incPC(instr.bytes)
	return instr.durationAction
}

func ldDD8(mem *memory, reg *register, instr *instruction) int {
	reg.D = readArgByte(mem, reg, 1)
	reg.incPC(instr.bytes)
	return instr.durationAction
}

func cpHL(mem *memory, reg *register, instr *instruction) int {
	value := mem.read8(reg.readDuo(REG_HL))

	reg.Flag.Z = reg.A == value
	reg.Flag.N = true
	reg.Flag.H = true
	reg.Flag.C = reg.A < value

	reg.incPC(instr.bytes)
	return instr.durationAction
}

func ldAL(mem *memory, reg *register, instr *instruction) int {
	reg.A = reg.L
	reg.incPC(instr.bytes)
	return instr.durationAction
}

func ldAA(mem *memory, reg *register, instr *instruction) int {
	reg.incPC(instr.bytes)
	return instr.durationAction
}

func ldAB(mem *memory, reg *register, instr *instruction) int {
	reg.A = reg.B
	reg.incPC(instr.bytes)
	return instr.durationAction
}

func ldAC(mem *memory, reg *register, instr *instruction) int {
	reg.A = reg.C
	reg.incPC(instr.bytes)
	return instr.durationAction
}

func addAHL(mem *memory, reg *register, instr *instruction) int {
	a := reg.A
	value := mem.read8(reg.readDuo(REG_HL))

	reg.Flag.Z = reg.A + value == 0
	reg.Flag.N = false
	reg.Flag.H = a&0xf+value&0xf > 0xf
	reg.Flag.C = reg.A+value > 0xff

	reg.A = reg.A + value

	reg.incPC(instr.bytes)
	return instr.durationAction
}

func addAA(mem *memory, reg *register, instr *instruction) int {
	var val = uint16(reg.A) + uint16(reg.A)

	reg.Flag.Z = reg.A == 0
	reg.Flag.N = false
	reg.Flag.H = reg.A&0xf+reg.A&0xf > 0xf
	reg.Flag.C = val > 0xff

	reg.A = uint8(val)

	reg.incPC(instr.bytes)
	return instr.durationAction
}

func addAB(mem *memory, reg *register, instr *instruction) int {
	var val = uint16(reg.A) + uint16(reg.B)

	reg.Flag.Z = reg.A == 0
	reg.Flag.N = false
	reg.Flag.H = reg.A&0xf+reg.B&0xf > 0xf
	reg.Flag.C = val > 0xff

	reg.A = uint8(val)

	reg.incPC(instr.bytes)
	return instr.durationAction
}

func noop(mem *memory, reg *register, instr *instruction) int {
	reg.incPC(instr.bytes)
	return instr.durationAction
}

func jpnn(mem *memory, reg *register, instr *instruction) int {
	address := readArgHalfword(mem, reg, 1)
	reg.PC = address
	return instr.durationAction
}

func di(mem *memory, reg *register, instr *instruction) int {
	reg.incPC(instr.bytes)
	return instr.durationAction
}

func ei(mem *memory, reg *register, instr *instruction) int {
	reg.incPC(instr.bytes)
	return instr.durationAction
}

func lddHLn(mem *memory, reg *register, instr *instruction) int {
	address := reg.readDuo(REG_HL)
	value := readArgByte(mem, reg, 1)

	mem.write8(address, value)
	reg.incPC(instr.bytes)
	return instr.durationAction
}

func ldAHLP(mem *memory, reg *register, instr *instruction) int {
	reg.A = mem.read8(reg.readDuo(REG_HL))
	reg.incrDuo(REG_HL)

	reg.incPC(instr.bytes)
	return instr.durationAction
}

func ldBCnn(mem *memory, reg *register, instr *instruction) int {
	arg := readArgHalfword(mem, reg, 1)
	reg.writeDuo(REG_BC, arg)

	reg.incPC(instr.bytes)
	return instr.durationAction
}

func decBC(mem *memory, reg *register, instr *instruction) int {
	reg.decrDuo(REG_BC)
	reg.incPC(instr.bytes)
	return instr.durationAction
}

func orC(mem *memory, reg *register, instr *instruction) int {
	reg.A = reg.C | reg.A

	reg.Flag.Z = reg.A == 0
	reg.Flag.N = false
	reg.Flag.H = false
	reg.Flag.C = false
	reg.incPC(instr.bytes)
	return instr.durationAction
}

func pushAf(mem *memory, reg *register, instr *instruction) int {
	pushStack16(mem, reg, reg.readDuo(REG_AF))
	// Does not affect flags
	reg.incPC(instr.bytes)
	return instr.durationAction
}

func pushDe(mem *memory, reg *register, instr *instruction) int {
	pushStack16(mem, reg, reg.readDuo(REG_DE))
	// Does not affect flags
	reg.incPC(instr.bytes)
	return instr.durationAction
}

func pushHl(mem *memory, reg *register, instr *instruction) int {
	pushStack16(mem, reg, reg.readDuo(REG_HL))
	// Does not affect flags
	reg.incPC(instr.bytes)
	return instr.durationAction
}

func andA(mem *memory, reg *register, instr *instruction) int {
	reg.A = reg.A & reg.A

	reg.Flag.Z = reg.A == 0
	reg.Flag.N = false
	reg.Flag.H = true
	reg.Flag.C = false

	reg.incPC(instr.bytes)
	return instr.durationAction
}

func andC(mem *memory, reg *register, instr *instruction) int {
	reg.C &= reg.C

	reg.Flag.Z = reg.C == 0
	reg.Flag.N = false
	reg.Flag.H = true
	reg.Flag.C = false

	reg.incPC(instr.bytes)
	return instr.durationAction
}

func retNz(mem *memory, reg *register, instr *instruction) int {
	if !reg.Flag.Z {
		addr := popStack16(mem, reg)
		reg.PC = addr
		return instr.durationAction
	}
	reg.incPC(instr.bytes)
	return instr.durationNoop
}

func callZa16(mem *memory, reg *register, instr *instruction) int {
	if reg.Flag.Z {
		pushStack16(mem, reg, reg.PC+uint16(3))
		addr := readArgHalfword(mem, reg, 1)
		reg.PC = addr
		return instr.durationAction
	}
	reg.incPC(instr.bytes)
	return instr.durationNoop
}

func ldAa16(mem *memory, reg *register, instr *instruction) int {
	reg.A = mem.read8(readArgHalfword(mem, reg, 1))

	reg.incPC(instr.bytes)
	return instr.durationAction
}

func retZ(mem *memory, reg *register, instr *instruction) int {
	if reg.Flag.Z {
		addr := popStack16(mem, reg)
		reg.PC = addr
		return instr.durationAction
	}
	reg.incPC(instr.bytes)
	return instr.durationNoop
}

func incHL(mem *memory, reg *register, instr *instruction) int {
	val := mem.read8(reg.readDuo(REG_HL)) + 1
	mem.write8(reg.readDuo(REG_HL), val)

	reg.Flag.Z = val == 0
	reg.Flag.N = false
	reg.Flag.H = val == 0xf

	reg.incPC(instr.bytes)
	return instr.durationAction
}

func incA(mem *memory, reg *register, instr *instruction) int {
	incRegister(&reg.A, reg)
	reg.incPC(instr.bytes)
	return instr.durationAction
}

func incE(mem *memory, reg *register, instr *instruction) int {
	incRegister(&reg.C, reg)
	reg.incPC(instr.bytes)
	return instr.durationAction
}

func popAf(mem *memory, reg *register, instr *instruction) int {
	reg.writeDuo(REG_AF, popStack16(mem, reg))
	reg.incPC(instr.bytes)
	return instr.durationAction
}

func popDe(mem *memory, reg *register, instr *instruction) int {
	reg.writeDuo(REG_DE, popStack16(mem, reg))
	// Does not affect flags
	reg.incPC(instr.bytes)
	return instr.durationAction
}

func popHl(mem *memory, reg *register, instr *instruction) int {
	reg.writeDuo(REG_HL, popStack16(mem, reg))
	// Does not affect flags
	reg.incPC(instr.bytes)
	return instr.durationAction
}

func reti(mem *memory, reg *register, instr *instruction) int {
	reg.PC = popStack16(mem, reg)
	// Does not affect flags
	return instr.durationAction
}

func cpl(mem *memory, reg *register, instr *instruction) int {
	reg.A = ^reg.A
	reg.Flag.N = true
	reg.Flag.H = true
	reg.incPC(instr.bytes)
	return instr.durationAction
}

func andd8(mem *memory, reg *register, instr *instruction) int {
	arg := readArgByte(mem, reg, 1)
	reg.A = arg & reg.A
	reg.Flag.Z = reg.A == 0
	reg.Flag.N = false
	reg.Flag.H = true
	reg.Flag.C = false
	reg.incPC(instr.bytes)
	return instr.durationAction
}

func ldBA(mem *memory, reg *register, instr *instruction) int {
	reg.B = reg.A
	reg.incPC(instr.bytes)
	return instr.durationAction
}

func ldBB(mem *memory, reg *register, instr *instruction) int {
	reg.incPC(instr.bytes)
	return instr.durationAction
}

func orB(mem *memory, reg *register, instr *instruction) int {
	reg.A = reg.B | reg.A

	reg.Flag.Z = reg.A == 0
	reg.Flag.N = false
	reg.Flag.H = false
	reg.Flag.C = false

	reg.incPC(instr.bytes)
	return instr.durationAction
}

func xorC(mem *memory, reg *register, instr *instruction) int {
	reg.A ^= reg.C
	reg.Flag.Z = reg.C == 0
	reg.Flag.N = false
	reg.Flag.H = false
	reg.Flag.C = false

	reg.incPC(instr.bytes)
	return instr.durationAction
}

func rst28(mem *memory, reg *register, instr *instruction) int {
	pushStack16(mem, reg, reg.PC+1)
	reg.PC = 0x28
	return instr.durationAction
}

func rst00(mem *memory, reg *register, instr *instruction) int {
	pushStack16(mem, reg, reg.PC+1)
	reg.PC = 0x00
	return instr.durationAction
}

func addHlDe(mem *memory, reg *register, instr *instruction) int {
	val := reg.readDuo(REG_HL)
	toAdd := reg.readDuo(REG_DE)

	reg.Flag.N = false
	reg.Flag.H = val&0xf+toAdd&0xf > 0xf
	reg.Flag.C = val+toAdd > 0xff

	reg.writeDuo(REG_HL, val+toAdd)

	reg.incPC(instr.bytes)
	return instr.durationAction
}

func ldEHL(mem *memory, reg *register, instr *instruction) int {
	reg.E = mem.read8(reg.readDuo(REG_HL))

	reg.incPC(instr.bytes)
	return instr.durationAction
}

func ldDHL(mem *memory, reg *register, instr *instruction) int {
	reg.D = mem.read8(reg.readDuo(REG_HL))

	reg.incPC(instr.bytes)
	return instr.durationAction
}

func jphl(mem *memory, reg *register, instr *instruction) int {
	reg.PC = reg.readDuo(REG_HL)
	return instr.durationAction
}

func ldDeA(mem *memory, reg *register, instr *instruction) int {
	mem.write8(reg.readDuo(REG_DE), reg.A)
	reg.incPC(instr.bytes)
	return instr.durationAction
}

func ldAHL(mem *memory, reg *register, instr *instruction) int {
	reg.A = mem.read8(reg.readDuo(REG_HL))
	reg.incPC(instr.bytes)
	return instr.durationAction
}

func incL(mem *memory, reg *register, instr *instruction) int {
	reg.L = reg.L + 1

	reg.Flag.Z = reg.L == 0
	reg.Flag.N = false
	reg.Flag.H = reg.L == 0x10

	reg.incPC(instr.bytes)
	return instr.durationAction
}

func addHlBc(mem *memory, reg *register, instr *instruction) int {
	val := reg.readDuo(REG_HL)
	reg.writeDuo(REG_HL, val+reg.readDuo(REG_BC))
	reg.Flag.N = false
	reg.Flag.H = val <= 0xf && reg.readDuo(REG_HL) > 0xf
	reg.Flag.C = val <= 0xff && reg.readDuo(REG_HL) > 0xff

	reg.incPC(instr.bytes)
	return instr.durationAction
}

func ldCHL(mem *memory, reg *register, instr *instruction) int {
	reg.C = mem.read8(reg.readDuo(REG_HL))

	reg.incPC(instr.bytes)
	return instr.durationAction
}

func ldBHL(mem *memory, reg *register, instr *instruction) int {
	reg.B = mem.read8(reg.readDuo(REG_HL))

	reg.incPC(instr.bytes)
	return instr.durationAction
}

func ldEC(mem *memory, reg *register, instr *instruction) int {
	reg.E = reg.C

	reg.incPC(instr.bytes)
	return instr.durationAction
}

func ldHB(mem *memory, reg *register, instr *instruction) int {
	reg.H = reg.B

	reg.incPC(instr.bytes)
	return instr.durationAction
}

func ldLA(mem *memory, reg *register, instr *instruction) int {
	reg.L = reg.A

	reg.incPC(instr.bytes)
	return instr.durationAction
}

func ldABC(mem *memory, reg *register, instr *instruction) int {
	reg.A = mem.read8(reg.readDuo(REG_BC))

	reg.incPC(instr.bytes)
	return instr.durationAction
}

func addAL(mem *memory, reg *register, instr *instruction) int {
	val := reg.A
	reg.A += reg.L

	reg.Flag.Z = reg.A == 0
	reg.Flag.N = false
	reg.Flag.H = val <= 0xf && reg.A > 0xf
	reg.Flag.C = val >= 0x80 && reg.A < 0x80

	reg.incPC(instr.bytes)
	return instr.durationAction
}

func adcAC(mem *memory, reg *register, instr *instruction) int {
	var i = uint16(reg.A) + uint16(reg.C)

	if reg.Flag.C {
		i++
	}

	hc := reg.A & 0xf + reg.C & 0xf
	if reg.Flag.C {
		hc++
	}

	reg.Flag.Z = i == 0
	reg.Flag.N = false
	reg.Flag.H = hc > 0xf
	reg.Flag.C = i > 0xff

	reg.A = uint8(i)

	reg.incPC(instr.bytes)
	return instr.durationAction
}
