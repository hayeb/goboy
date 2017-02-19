package gameboy

type cbInstruction struct {
	name  string
	bytes int

	actionDuration int
	noopDuration   int

	executor instructionExecutor
}

func createCBInstructionMap() *map[uint8]cbInstruction {
	return &map[uint8]cbInstruction{
		0x7c: newCBInstruction("BIT_7_H", 2, 8, bit_7_h),
	}
}

func newCBInstruction(name string, length int, duration int, fp func(mem *memory, reg *register)) cbInstruction {
	return newCBConditionalInstruction(name, length, duration, 0, fp)
}

func newCBConditionalInstruction(name string, length int, actionDuration int, noopDuration int, fp func(mem *memory, reg *register)) cbInstruction {
	return cbInstruction{
		name:           name,
		bytes:          length,
		actionDuration: actionDuration,
		noopDuration:   noopDuration,
		executor:       fp,
	}
}

type cbInstructionExecutor func(mem *memory, reg *register)

func bit_7_h(_ *memory, reg *register) {
	reg.bit(1<<7, reg.H.val())
}
