package gameboy

type instruction struct {
	name string
	// The length of the instruction in bytes, including the opcode
	bytes int
}

func createInstructionMap() *map[uint8]instruction {
	return &map[uint8]instruction{
		0xc:  {inc_c, 1},
		0xe:  {ld_c, 2},
		0x20: {jr, 2},
		0x21: {ld_hl, 3},
		0x31: {ld_sp, 3},
		0x32: {ldd_hl_a, 1},
		0x3e: {ld_a, 2},
		0x77: {ld_HL_a, 1},
		0xAF: {xor_a, 1},
		0xCB: {cb, 1},
		0xE0: {ldh_a8_A, 2},
		0xE2: {ld_C_a, 1},
	}
}

const (
	cb       = "CB"
	inc_c    = "INC_C"
	jr       = "JR"
	ld_a     = "LD_A"
	ld_C_a   = "LD_(C)_A"
	ld_c     = "LD_C"
	ld_sp    = "LD_SP"
	ld_hl    = "LD_HL"
	ld_HL_a  = "LD_(HL)_a"
	ldh_a8_A = "LDH_a8_A"
	ldd_hl_a = "LDD_(HL)_A"
	xor_a    = "XOR_A"
)
