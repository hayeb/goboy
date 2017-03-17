package gameboy

type cbInstruction struct {
	name  string
	bytes int

	actionDuration int
	noopDuration   int

	executor cbInstructionExecutor
}

type cbInstructionExecutor func(mem *memory, reg *register, cbInstr *cbInstruction) int

func createCBInstructionMap() *map[uint8]cbInstruction {
	return &map[uint8]cbInstruction{
		0x7c: newCBInstruction("BIT 7,H", 2, 8, bit_7_h),
		0x11: newCBInstruction("RL C", 2, 8, rl_c),
	}
}

func newCBInstruction(name string, length int, duration int, fp func(mem *memory, reg *register, cbInstr *cbInstruction) int) cbInstruction {
	return newCBConditionalInstruction(name, length, duration, 0, fp)
}

func newCBConditionalInstruction(name string, length int, actionDuration int, noopDuration int, fp func(mem *memory, reg *register, cbInstr *cbInstruction) int) cbInstruction {
	return cbInstruction{
		name:           name,
		bytes:          length,
		actionDuration: actionDuration,
		noopDuration:   noopDuration,
		executor:       fp,
	}
}

func bit_7_h(_ *memory, reg *register, cbInstr *cbInstruction) int {
	reg.bit(7, reg.H.val())
	return cbInstr.actionDuration
}

func rl_c(mem *memory, reg *register, cbInstr *cbInstruction) int {
	newVal, carry := rLeftCarry(reg.C.val(), reg.Flag.C)
	reg.C = byteRegister(newVal)
	reg.Flag.Z = reg.C == 0
	reg.Flag.N = false
	reg.Flag.H = false
	reg.Flag.C = carry
	return cbInstr.actionDuration
}
