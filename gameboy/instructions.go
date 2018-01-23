package gameboy

import "fmt"

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
		0x04: newInstruction("INC B", 1, 4, incB),
		0x05: newInstruction("DEC B", 1, 4, decB),
		0x06: newInstruction("LD B, d8", 2, 8, ldBD8),
		0x0b: newInstruction("DEC BC", 1, 8, decBC),
		0x0c: newInstruction("INC C", 1, 4, incC),
		0x0d: newInstruction("DEC C", 1, 4, decC),
		0x0e: newInstruction("LD C", 2, 8, ldC),
		0x11: newInstruction("LD DE,d16", 3, 12, ldDeD16),
		0x13: newInstruction("INC DE", 1, 8, incDe),
		0x15: newInstruction("DEC D", 1, 4, decD),
		0x16: newInstruction("LD D,d8", 2, 8, ldDD8),
		0x17: newInstruction("RLA", 1, 4, rla),
		0x18: newInstruction("JR r8", 2, 12, jrR8),
		0x1a: newInstruction("LD A,(DE)", 1, 8, ldADE),
		0x1d: newInstruction("DEC E", 1, 4, decE),
		0x1e: newInstruction("LD E,d8", 2, 8, ldED8),
		0x20: newConditionalInstruction("JR NZ,r8", 2, 12, 8, jrNz),
		0x21: newInstruction("LD HL", 3, 12, ldHl),
		0x22: newInstruction("LD (HL+),A", 1, 8, ldHLPA),
		0x23: newInstruction("INC HL", 1, 8, incHl),
		0x24: newInstruction("INC H", 1, 4, incH),
		0x28: newConditionalInstruction("JR Z,r8", 2, 12, 8, jrZR8),
		0x2a: newInstruction("LD A,(HL+)", 1, 8, ldAHLP),
		0x2e: newInstruction("LD L,d8", 2, 8, ldLD8),
		0x31: newInstruction("LD SP", 3, 12, ldSp),
		0x32: newInstruction("LDD (HL-),A", 1, 8, lddHLA),
		0x36: newInstruction("LD (HL),n", 1, 12, lddHLn),
		0x3d: newInstruction("DEC A", 1, 4, decA),
		0x3e: newInstruction("LD A", 2, 8, ldA),
		0x4f: newInstruction("LD C,A", 1, 4, ldcA),
		0x57: newInstruction("LD D,A", 1, 4, ldDA),
		0x67: newInstruction("LD H,A", 1, 4, ldHA),
		0x77: newInstruction("LD (HL),A", 1, 8, ldHLA),
		0x78: newInstruction("LD A,B", 1, 4, ldAB),
		0x7b: newInstruction("LD A,E", 1, 4, ldAE),
		0x7c: newInstruction("LD A,H", 1, 4, ldAH),
		0x7d: newInstruction("LD A,L", 1, 4, ldAL),
		0x86: newInstruction("ADD A,(HL)", 1, 8, addAHL),
		0x90: newInstruction("SUB B", 1, 4, subB),
		0xa7: newInstruction("AND A", 1, 4, andA),
		0xaf: newInstruction("XOR A", 1, 4, xorA),
		0xb1: newInstruction("OR C", 1, 4, orC),
		0xbe: newInstruction("CP (HL)", 1, 8, cpHL),
		0xc0: newConditionalInstruction("RET NZ", 1, 20, 8, retNz),
		0xc1: newInstruction("POP BC", 1, 12, popBc),
		0xc3: newInstruction("JP a16", 0, 16, jpnn), // JP nn no length: interference with updating PC
		0xc5: newInstruction("PUSH BC", 1, 16, pushBc),
		0xc9: newInstruction("RET", 0, 16, ret), // RET no length: interference with updating PC
		0xcb: newInstruction("CB", 1, 4, nil),
		0xcc: newConditionalInstruction("CALL Z,a16", 3, 24, 12, callZa16),
		0xcd: newInstruction("CALL a16", 0, 12, callNn), // CALL has no length, as it interferes with updating PC
		0xd5: newInstruction("PUSH DE", 1, 16, pushDe),
		0xe0: newInstruction("LDH a8,A", 2, 12, ldhA8A),
		0xe2: newInstruction("LD (C),A", 1, 8, ldCA),
		0xe5: newInstruction("PUSH HL", 1, 8, pushHl),
		0xea: newInstruction("LD (a16),A", 3, 16, ldA16A),
		0xf0: newInstruction("LDH A,(a8)", 2, 12, ldAA8),
		0xf3: newInstruction("DI", 1, 4, di),
		0xf5: newInstruction("PUSH AF", 1, 16, pushAf),
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

func incRegister(br *byteRegister, reg *register) {
	val := br.val()
	*br = byteRegister(val + 1)
	reg.Flag.Z = *br == 0
	reg.Flag.N = false
	reg.Flag.H = *br > 0xf
}

func decRegister(br *byteRegister, reg *register) {
	val := br.val()
	*br = byteRegister(val - uint8(1))
	reg.Flag.Z = *br == 0
	reg.Flag.N = true
	reg.Flag.H = *br < 0xf
}

func subRegister(reg *register, val uint8) {
	before := reg.A.val()
	reg.A = byteRegister(reg.A.val() - val)
	reg.Flag.Z = reg.A == 0
	reg.Flag.N = true
	reg.Flag.H = before < val
}

func callNn(mem *memory, reg *register, instr *instruction) int {
	pushStack16(mem, reg, reg.PC.val()+uint16(3))
	reg.PC = halfWordRegister(readArgHalfword(mem, reg, 1))
	// Does not affect flags
	return instr.durationAction
}

func jrNz(mem *memory, reg *register, instr *instruction) int {
	n := readArgByte(mem, reg, 1)
	if !reg.Flag.Z {
		reg.PC = halfWordRegister(int(reg.PC.val()) + int(int8(n)))
		return instr.durationAction
	}
	// Does not affect flags
	return instr.durationNoop
}

func incB(mem *memory, reg *register, instr *instruction) int {
	incRegister(&reg.B, reg)
	return instr.durationAction
}

func incC(mem *memory, reg *register, instr *instruction) int {
	incRegister(&reg.C, reg)
	return instr.durationAction
}

func ldA(mem *memory, reg *register, instr *instruction) int {
	reg.A = byteRegister(readArgByte(mem, reg, 1))
	return instr.durationAction
}

func ldADE(mem *memory, reg *register, instr *instruction) int {
	val := mem.read8(reg.readDuo(reg_de))
	reg.A = byteRegister(val)
	// Does not affect flags
	return instr.durationAction
}

func ldC(mem *memory, reg *register, instr *instruction) int {
	reg.C = byteRegister(readArgByte(mem, reg, 1))
	// Does not affect flags
	return instr.durationAction
}

func ldCA(mem *memory, reg *register, instr *instruction) int {
	mem.write8(0xFF00+uint16(reg.C.val()), reg.A.val())
	// Does not affect flags
	return instr.durationAction
}

func ldDeD16(mem *memory, reg *register, instr *instruction) int {
	reg.writeDuo(reg_de, readArgHalfword(mem, reg, 1))
	// Does not affect flags
	return instr.durationAction
}

func ldSp(mem *memory, reg *register, instr *instruction) int {
	arg := readArgHalfword(mem, reg, 1)
	reg.SP = halfWordRegister(arg)
	// Does not affect flags
	return instr.durationAction
}

func ldHl(mem *memory, reg *register, instr *instruction) int {
	reg.writeDuo(reg_hl, readArgHalfword(mem, reg, 1))
	// Does not affect flags
	return instr.durationAction
}

func ldHLA(mem *memory, reg *register, instr *instruction) int {
	mem.write8(reg.readDuo(reg_hl), reg.A.val())
	// Does not affect flags
	return instr.durationAction
}

func lddHLA(mem *memory, reg *register, instr *instruction) int {
	mem.write8(reg.readDuo(reg_hl), reg.A.val())
	reg.decrDuo(reg_hl)
	// Does not affect flags
	return instr.durationAction
}

func ldhA8A(mem *memory, reg *register, instr *instruction) int {
	arg := readArgByte(mem, reg, 1)
	mem.write8(0xff00+uint16(arg), reg.A.val())
	// Does not affect flags
	return instr.durationAction
}

func xorA(mem *memory, reg *register, instr *instruction) int {
	reg.A = reg.A ^ reg.A
	if reg.A == 0 {
		reg.Flag.Z = true
	}
	reg.Flag.N = false
	reg.Flag.H = false
	reg.Flag.C = false
	return instr.durationAction
}

func pushBc(mem *memory, reg *register, instr *instruction) int {
	pushStack16(mem, reg, reg.readDuo(reg_bc))
	// Does not affect flags
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
	return instr.durationAction
}

func popBc(mem *memory, reg *register, instr *instruction) int {
	val := popStack16(mem, reg)
	reg.writeDuo(reg_bc, val)
	// Does not affect flags
	return instr.durationAction
}

func decB(mem *memory, reg *register, instr *instruction) int {
	decRegister(&reg.B, reg)
	return instr.durationAction
}

func ldcA(mem *memory, reg *register, instr *instruction) int {
	reg.C = reg.A
	return instr.durationAction
}

func ldBD8(mem *memory, reg *register, instr *instruction) int {
	arg := readArgByte(mem, reg, 1)
	reg.B = byteRegister(arg)
	return instr.durationAction
}

func ldHLPA(mem *memory, reg *register, instr *instruction) int {
	mem.write8(reg.readDuo(reg_hl), reg.A.val())
	reg.incrDuo(reg_hl)
	return instr.durationAction
}

func incHl(mem *memory, reg *register, instr *instruction) int {
	reg.incrDuo(reg_hl)
	return instr.durationAction
}

func incDe(mem *memory, reg *register, instr *instruction) int {
	reg.incrDuo(reg_de)
	return instr.durationAction
}

func ret(mem *memory, reg *register, instr *instruction) int {
	addr := popStack16(mem, reg)
	reg.PC = halfWordRegister(addr)
	return instr.durationAction
}

func ldAE(mem *memory, reg *register, instr *instruction) int {
	reg.A = reg.E
	return instr.durationAction
}

func cpD8(mem *memory, reg *register, instr *instruction) int {
	var arg uint8 = readArgByte(mem, reg, 1)
	reg.Flag.Z = reg.A.val() == arg
	reg.Flag.N = true
	reg.Flag.H = true // TODO: Check borrow
	reg.Flag.C = reg.A.val() < arg
	return instr.durationAction
}

func ldA16A(mem *memory, reg *register, instr *instruction) int {
	mem.write8(readArgHalfword(mem, reg, 1), reg.A.val())
	return instr.durationAction
}

func decA(mem *memory, reg *register, instr *instruction) int {
	decRegister(&reg.A, reg)
	return instr.durationAction
}

func jrZR8(mem *memory, reg *register, instr *instruction) int {
	if reg.Flag.Z {
		reg.PC = halfWordRegister(reg.PC.val() + uint16(readArgByte(mem, reg, 1)))
		return instr.durationAction
	}
	return instr.durationNoop
}

func decC(mem *memory, reg *register, instr *instruction) int {
	decRegister(&reg.C, reg)
	return instr.durationAction
}

func ldLD8(mem *memory, reg *register, instr *instruction) int {
	reg.L = byteRegister(readArgByte(mem, reg, 1))
	return instr.durationAction
}

func jrR8(mem *memory, reg *register, instr *instruction) int {
	oldPc := reg.PC.val()
	arg := int8(readArgByte(mem, reg, 1))
	reg.PC = halfWordRegister(uint16(int(oldPc) + int(arg)))
	return instr.durationAction
}

func ldHA(mem *memory, reg *register, instr *instruction) int {
	reg.H = reg.A
	return instr.durationAction
}

func ldDA(mem *memory, reg *register, instr *instruction) int {
	reg.D = reg.A
	return instr.durationAction
}

func ldED8(mem *memory, reg *register, instr *instruction) int {
	d8 := readArgByte(mem, reg, 1)
	reg.E = byteRegister(d8)
	return instr.durationAction
}

func ldAA8(mem *memory, reg *register, instr *instruction) int {
	address := 0xFF00 + uint16(readArgByte(mem, reg, 1))
	reg.A = byteRegister(mem.read8(address))
	return instr.durationAction
}

func decE(mem *memory, reg *register, instr *instruction) int {
	decRegister(&reg.E, reg)
	return instr.durationAction
}

func incH(mem *memory, reg *register, instr *instruction) int {
	incRegister(&reg.H, reg)
	return instr.durationAction
}

func ldAH(mem *memory, reg *register, instr *instruction) int {
	reg.A = reg.H
	return instr.durationAction
}

func subB(mem *memory, reg *register, instr *instruction) int {
	subRegister(reg, reg.B.val())
	return instr.durationAction
}

func decD(mem *memory, reg *register, instr *instruction) int {
	decRegister(&reg.D, reg)
	return instr.durationAction
}

func ldDD8(mem *memory, reg *register, instr *instruction) int {
	reg.D = byteRegister(readArgByte(mem, reg, 1))
	return instr.durationAction
}

func cpHL(mem *memory, reg *register, instr *instruction) int {
	value := mem.read8(reg.readDuo(reg_hl))

	reg.Flag.Z = reg.A.val() == value
	reg.Flag.N = true
	reg.Flag.H = true // TODO: Borrow?
	reg.Flag.C = reg.A.val() < value

	return instr.durationAction
}

func ldAL(mem *memory, reg *register, instr *instruction) int {
	reg.A = byteRegister(reg.L.val())
	return instr.durationAction
}

func ldAB(mem *memory, reg *register, instr *instruction) int {
	reg.A = byteRegister(reg.B.val())
	return instr.durationAction
}

func addAHL(mem *memory, reg *register, instr *instruction) int {
	a := reg.A.val()
	value := mem.read8(reg.readDuo(reg_hl))
	reg.A = byteRegister(reg.A.val() + value)

	reg.Flag.Z = reg.A.val() == 0
	reg.Flag.N = false
	reg.Flag.H = a <= 15 && value > 15
	reg.Flag.C = true // TODO: Check carry

	return instr.durationAction
}

func noop(mem *memory, reg *register, instr *instruction) int {
	return instr.durationAction
}

func jpnn(mem *memory, reg *register, instr *instruction) int {
	address := readArgHalfword(mem, reg, 1)
	fmt.Printf("Jumping to %#04x\n", address)
	reg.PC = halfWordRegister(address)
	return instr.durationAction
}

func di(mem *memory, reg *register, instr *instruction) int {
	return instr.durationAction
}

func ei(mem *memory, reg *register, instr *instruction) int {
	return instr.durationAction
}

func lddHLn(mem *memory, reg *register, instr *instruction) int {
	address := reg.readDuo(reg_hl)
	value := readArgByte(mem, reg, 1)

	mem.write8(address, value)
	return instr.durationAction
}

func ldAHLP(mem *memory, reg *register, instr *instruction) int {
	value := mem.read8(reg.readDuo(reg_hl))
	reg.A = byteRegister(value)
	reg.incrDuo(reg_hl)

	return instr.durationAction
}

func ldBCnn(mem *memory, reg *register, instr *instruction) int {
	arg := readArgHalfword(mem, reg, 1)
	reg.writeDuo(reg_bc, arg)
	return instr.durationAction
}

func decBC(mem *memory, reg *register, instr *instruction) int {
	reg.decrDuo(reg_bc)
	return instr.durationAction
}

func orC(mem *memory, reg *register, instr *instruction) int {
	reg.A = reg.C | reg.A

	reg.Flag.Z = reg.A == 0
	reg.Flag.N = false
	reg.Flag.H = false
	reg.Flag.C = false
	return instr.durationAction
}

func pushAf(mem *memory, reg *register, instr *instruction) int {
	pushStack16(mem, reg, reg.readDuo(reg_af))
	// Does not affect flags
	return instr.durationAction
}

func pushDe(mem *memory, reg *register, instr *instruction) int {
	pushStack16(mem, reg, reg.readDuo(reg_de))
	// Does not affect flags
	return instr.durationAction
}

func pushHl(mem *memory, reg *register, instr *instruction) int {
	pushStack16(mem, reg, reg.readDuo(reg_hl))
	// Does not affect flags
	return instr.durationAction
}

func andA(mem *memory, reg *register, instr *instruction) int {
	reg.A = byteRegister(reg.A.val() & reg.A.val())

	reg.Flag.Z = reg.A == 0
	reg.Flag.N = false
	reg.Flag.H = true
	reg.Flag.C = false
	return instr.durationAction
}

func retNz(mem *memory, reg *register, instr *instruction) int {
	if !reg.Flag.Z {
		addr := popStack16(mem, reg)
		reg.PC = halfWordRegister(addr)
		return instr.durationAction
	}
	return instr.durationNoop
}

func callZa16(mem *memory, reg *register, instr *instruction) int {
	if reg.Flag.Z {
		pushStack16(mem, reg, reg.PC.val()+uint16(3))
		addr := readArgHalfword(mem, reg, 1)
		reg.PC = halfWordRegister(addr)
		return instr.durationAction
	}
	return instr.durationNoop
}