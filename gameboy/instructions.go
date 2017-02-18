package gameboy

type instruction struct {
	name  string
	bytes int
}

func createInstructionMap() *map[uint8]instruction {
	return &map[uint8]instruction{
		0x20: {jr, 2},
		0x21: {ld_hl, 3},
		0x31: {ld_sp, 3},
		0x32: {ldd_hl_a, 1},
		0xAF: {xor_a, 1},
		0xCB: {cb, 1},
	}
}

const (
	jr       = "JR"
	cb       = "CB"
	ld_sp    = "LD SP"
	ld_hl    = "LD_HL"
	ldd_hl_a = "LDD_(HL)_A"
	xor_a    = "XOR A"
)
