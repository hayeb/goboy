package gameboy

type cbInstruction struct {
	name  string
	bytes int

	actionDuration int
	noopDuration   int

	executor cbInstructionExecutor
}

type cbInstructionExecutor func(mem *memory, reg *register, cbInstr *cbInstruction) int

func createCBInstructionMap() *map[uint8]*cbInstruction {
	return &map[uint8]*cbInstruction{
		0x11: newCBInstruction("RL C", 2, 8, rl_c),
		0x27: newCBInstruction("SLA A", 2, 8, slaA),
		0x37: newCBInstruction("SWAP A", 2, 8, swapA),
		0x50: newCBInstruction("BIT 2,B", 2, 8, bit_2_b),
		0x58: newCBInstruction("BIT 3,B", 2, 8, bit_3_b),
		0x60: newCBInstruction("BIT 4,B", 2, 8, bit_4_b),
		0x68: newCBInstruction("BIT 5,B", 2, 8, bit_5_b),
		0x7c: newCBInstruction("BIT 7,H", 2, 8, bit_7_h),
		0x87: newCBInstruction("RES 0, a", 2, 8, res0A),

	}
}

func newCBInstruction(name string, length int, duration int, fp func(mem *memory, reg *register, cbInstr *cbInstruction) int) *cbInstruction {
	return newCBConditionalInstruction(name, length, duration, 0, fp)
}

func newCBConditionalInstruction(name string, length int, actionDuration int, noopDuration int, fp func(mem *memory, reg *register, cbInstr *cbInstruction) int) *cbInstruction {
	return &cbInstruction{
		name:           name,
		bytes:          length,
		actionDuration: actionDuration,
		noopDuration:   noopDuration,
		executor:       fp,
	}
}

func bit_7_h(_ *memory, reg *register, cbInstr *cbInstruction) int {
	reg.bit(7, reg.H)

	reg.incPC(cbInstr.bytes)
	return cbInstr.actionDuration
}

func bit_2_b(_ *memory, reg *register, cbInstr *cbInstruction) int {
	reg.bit(2, reg.B)

	reg.incPC(cbInstr.bytes)
	return cbInstr.actionDuration
}

func bit_3_b(_ *memory, reg *register, cbInstr *cbInstruction) int {
	reg.bit(3, reg.B)

	reg.incPC(cbInstr.bytes)
	return cbInstr.actionDuration
}

func bit_4_b(_ *memory, reg *register, cbInstr *cbInstruction) int {
	reg.bit(4, reg.B)

	reg.incPC(cbInstr.bytes)
	return cbInstr.actionDuration
}

func bit_5_b(_ *memory, reg *register, cbInstr *cbInstruction) int {
	reg.bit(5, reg.B)

	reg.incPC(cbInstr.bytes)
	return cbInstr.actionDuration
}

func rl_c(mem *memory, reg *register, cbInstr *cbInstruction) int {
	isCarrySet := reg.Flag.C
	isMSBSet := testBit(reg.C, 7)

	reg.Flag.Z = false
	reg.Flag.N = false
	reg.Flag.H = false
	reg.Flag.C = false

	newVal := reg.C << 1

	reg.Flag.C = isMSBSet

	if isCarrySet {
		newVal = setBit(newVal, 0)
	}

	if newVal == 0 {
		reg.Flag.Z = true
	}
	reg.C = newVal
	reg.incPC(cbInstr.bytes)
	return cbInstr.actionDuration
}

func res0A(mem *memory, reg *register, cbInstr *cbInstruction) int {
	reg.A = resetBit(reg.A, 0)
	reg.incPC(cbInstr.bytes)
	return cbInstr.actionDuration
}

func swapA(mem *memory, reg *register, cbInstr *cbInstruction) int {
	val := reg.A

	reg.A = val << 4 | val >> 4
	reg.Flag.Z = reg.A == 0
	reg.Flag.N = false
	reg.Flag.H = false
	reg.Flag.C = false
	reg.incPC(cbInstr.bytes)
	return cbInstr.actionDuration
}

func slaA(mem *memory, reg *register, cbInstr *cbInstruction) int {
	reg.Flag.C = reg.A & 0x80 == 1
	reg.A = reg.A << 1

	reg.Flag.Z = reg.A == 0
	reg.Flag.N = false
	reg.Flag.H = false
	reg.incPC(cbInstr.bytes)
	return cbInstr.actionDuration
}
