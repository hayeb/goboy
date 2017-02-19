package gameboy

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
}

func createInstructionMap() *map[uint8]instruction {
	return &map[uint8]instruction{
		0xc:  newInstrucion(inc_c, 1, 4),
		0xe:  newInstrucion(ld_c, 2, 8),
		0x11: newInstrucion(ld_de_d16, 3, 12),
		0x1a: newInstrucion(ld_a_DE, 1, 8),
		0x20: newConditionalInstruction(jr, 2, 12, 8),
		0x21: newInstrucion(ld_hl, 3, 12),
		0x31: newInstrucion(ld_sp, 3, 12),
		0x32: newInstrucion(ldd_hl_a, 1, 8),
		0x3e: newInstrucion(ld_a, 2, 8),
		0x77: newInstrucion(ld_HL_a, 1, 8),
		0xAF: newInstrucion(xor_a, 1, 4),
		0xCB: newInstrucion(cb, 1, 4),
		0xCD: newInstrucion(call_nn, 3, 12),
		0xE0: newInstrucion(ldh_a8_A, 2, 12),
		0xE2: newInstrucion(ld_C_a, 1, 8),
	}
}

func newConditionalInstruction(name string, length int, actionDuration int, noopDuration int) instruction {
	return instruction{
		name:            name,
		bytes:           length,
		duration_action: actionDuration,
		duration_noop:   noopDuration,
	}
}

func newInstrucion(name string, length int, duration int) instruction {
	return newConditionalInstruction(name, length, duration, 0)
}

const (
	cb        = "CB"
	call_nn   = "CALL_nn"
	inc_c     = "INC_C"
	jr        = "JR"
	ld_a      = "LD_A"
	ld_a_DE   = "LD_A_(DE)"
	ld_C_a    = "LD_(C)_A"
	ld_c      = "LD_C"
	ld_de_d16 = "LD_DE_d16"
	ld_sp     = "LD_SP"
	ld_hl     = "LD_HL"
	ld_HL_a   = "LD_(HL)_a"
	ldh_a8_A  = "LDH_a8_A"
	ldd_hl_a  = "LDD_(HL)_A"
	xor_a     = "XOR_A"
)
