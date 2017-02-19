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
		0x11: newCBInstruction("RL_C", 2, 8, rl_c),
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

func rl_c(mem *memory, reg *register) {
	// TODO: Implement shifting THROUGH carry flag.
	val := reg.C.val()
	reg.C = byteRegister(rotateLeftThroughCarry(val, 1))
	reg.Flag.Z = reg.C == 0
	reg.Flag.N = false
	reg.Flag.H = false
	reg.Flag.C = val>>7 == 1
}

func rotateLeftThroughCarry(val uint8, n uint) uint8 {
	return val<<n | val>>(8-n)
}
