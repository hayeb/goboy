package gameboy

import "testing"

func dummyMemory() *memory {
	bootrom := [256]uint8{}
	cartridge := [32 * 1024]uint8{}
	return memInit(bootrom[:], cartridge[:])
}

func dummyRegs() *register {
	return new(register)
}

func TestIncB(t *testing.T) {
	mem := dummyMemory()
	regs := dummyRegs()
	regs.B = byteRegister(0x0)

	result := dummyRegs()
	result.B = byteRegister(0x1)
	result.Flag.Z = false
	result.Flag.N = false

	testInstruction(t, mem, regs, incB, result, mem, "INC B")

	regs.B = byteRegister(0xf)
	result.B = byteRegister(0x10)
	result.Flag.H = true

	testInstruction(t, mem, regs, incB, result, mem, "INC B")
}

func TestIncC(t *testing.T) {
	mem := dummyMemory()
	regs := dummyRegs()
	regs.C = byteRegister(0x0)

	result := dummyRegs()
	result.C = byteRegister(0x1)
	result.Flag.Z = false
	result.Flag.N = false

	testInstruction(t, mem, regs, incC, result, mem, "INC C")

	regs.C = byteRegister(0xf)
	result.C = byteRegister(0x10)
	result.Flag.H = true

	testInstruction(t, mem, regs, incC, result, mem, "INC C")
}

func TestDecB(t *testing.T) {
	mem := dummyMemory()
	regs := dummyRegs()
	regs.B = byteRegister(0x1)

	result := dummyRegs()
	result.B = byteRegister(0x0)
	result.Flag.Z = true
	result.Flag.N = true
	result.Flag.H = true

	testInstruction(t, mem, regs, decB, result, mem, "DEC B")

	regs.B = byteRegister(0x10)
	result.B = byteRegister(0xf)
	result.Flag.Z = false
	result.Flag.H = false

	testInstruction(t, mem, regs, decB, result, mem, "DEC B")
}

func TestDecC(t *testing.T) {
	mem := dummyMemory()
	regs := dummyRegs()
	regs.C = byteRegister(0x1)

	result := dummyRegs()
	result.C = byteRegister(0x0)
	result.Flag.Z = true
	result.Flag.N = true
	result.Flag.H = true

	testInstruction(t, mem, regs, decC, result, mem,"DEC C")

	regs.C = byteRegister(0x10)
	result.C = byteRegister(0xf)
	result.Flag.Z = false
	result.Flag.H = false

	testInstruction(t, mem, regs, decC, result, mem,"DEC C")
}

func TestLDSP(t *testing.T) {
	regs := dummyRegs()
	mem := dummyMemory()
	mem.write16(0x1, 0xfefe)

	resultReg := dummyRegs()
	resultReg.SP = halfWordRegister(0xfefe)

	testInstruction(t, mem, regs, ldSPd16, resultReg, mem, "LD SP")
}

func TestXOR(t *testing.T) {
	regs := dummyRegs()
	mem := dummyMemory()

	regs.A = byteRegister(0xfe)

	resultRegs := dummyRegs()
	resultRegs.A = 0
	resultRegs.Flag.Z = true
	resultRegs.Flag.N = false
	resultRegs.Flag.H = false
	resultRegs.Flag.C = false

	testInstruction(t, mem, regs, xorA, resultRegs, mem, "XOR A")
}

func TestLDHL(t *testing.T) {
	regs := dummyRegs()
	mem := dummyMemory()
	mem.write16(0x1, 0xfefe)

	resultReg := dummyRegs()
	resultReg.writeDuo(REG_HL, 0xfefe)

	testInstruction(t, mem, regs, ldHL, resultReg, mem, "LD HL")
}

func TestLDDHLA(t *testing.T) {
	regs := dummyRegs()
	mem := dummyMemory()

	regs.writeDuo(REG_HL, 0x1000)
	regs.A = 0xfe

	resultMem := dummyMemory()
	resultMem.write8(0x1000, 0xfe)

	resultReg := dummyRegs()
	resultReg.A = byteRegister(0xfe)
	resultReg.writeDuo(REG_HL, 0xfff)

	testInstruction(t, mem, regs, lddHLA, resultReg, mem, "LD (HL-) A")
}

func TestBit7H(t *testing.T) {
	regs := dummyRegs()
	mem := dummyMemory()

	regs.H = 0xf0

	resultReg := dummyRegs()
	resultReg.H = 0xf0
	resultReg.Flag.Z = false
	resultReg.Flag.N = false
	resultReg.Flag.H = true

	testCbInstruction(t, mem, regs, bit_7_h, resultReg, mem, "BIT 7,H")
}

func TestJRNZ(t *testing.T) {
	regs := dummyRegs()
	mem := dummyMemory()

	mem.write16(0x1, 0xffee)

	resultRegs := dummyRegs()
	resultRegs.PC = 0xffee
	testInstruction(t, mem, regs, jrNZr8, resultRegs, mem, "JR NZ")

	mem.write16(0x1, 0xffee)
	regs.Flag.Z = true
	testInstruction(t, mem, regs, jrNZr8, regs, mem, "JR NZ")
}

func testCbInstruction(t *testing.T, mem *memory, regs *register, executor cbInstructionExecutor, resultReg *register, resultMem *memory, name string) {
	executor(mem, regs, new(cbInstruction))

	if *regs != *resultReg {
		t.Errorf("Instruction %s failed, registers does not match:\nExpected:\n%+v\n\nGot:\n%+v", name, resultReg, regs)
	}

	if *mem != *resultMem {
		t.Errorf("Instruction %s failed, memory does not match:\nExpected:\n%+v\n\nGot:\n%+v", name, resultMem, mem)
	}
}


func testInstruction(t *testing.T, mem *memory, regs *register, executor instructionExecutor, resultReg *register, resultMem *memory, name string) {
	executor(mem, regs, new(instruction))

	if *regs != *resultReg {
		t.Errorf("Instruction %s failed, registers does not match:\nExpected:\n%+v\n\nGot:\n%+v", name, resultReg, regs)
	}

	if *mem != *resultMem {
		t.Errorf("Instruction %s failed, memory does not match:\nExpected:\n%+v\n\nGot:\n%+v", name, resultMem, mem)
	}
}