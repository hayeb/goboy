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
		0x37: newCBInstruction("SWAP A", 2, 8, swapA),
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
	reg.bit(7, reg.H.val())

	reg.incPC(cbInstr.bytes)
	return cbInstr.actionDuration
}

func rl_c(mem *memory, reg *register, cbInstr *cbInstruction) int {
	isCarrySet := reg.Flag.C
	isMSBSet := testBit(reg.C.val(), 7)

	reg.Flag.Z = false
	reg.Flag.N = false
	reg.Flag.H = false
	reg.Flag.C = false

	newVal := reg.C.val() << 1

	reg.Flag.C = isMSBSet

	if isCarrySet {
		newVal = setBit(newVal, 0)
	}

	if newVal == 0 {
		reg.Flag.Z = true
	}
	reg.C = byteRegister(newVal)
	reg.incPC(cbInstr.bytes)
	return cbInstr.actionDuration
}

func res0A(mem *memory, reg *register, cbInstr *cbInstruction) int {
	reg.A = byteRegister(resetBit(reg.A.val(), 0))
	reg.incPC(cbInstr.bytes)
	return cbInstr.actionDuration
}

func swapA(mem *memory, reg *register, cbInstr *cbInstruction) int {
	val := reg.A.val()

	reg.A = byteRegister(val << 4 | val >> 4)
	reg.Flag.Z = reg.A == 0
	reg.Flag.N = false
	reg.Flag.H = false
	reg.Flag.C = false
	reg.incPC(cbInstr.bytes)
	return cbInstr.actionDuration
}
