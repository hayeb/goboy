package gameboy

type Instruction struct {
	name  string
	bytes int
}

func CreateInstructionMap() *map[uint8]Instruction {
	return &map[uint8]Instruction{
		0x21: {LD_HL, 3},
		0x31: {LD_SP, 3},
		0x32: {LDD_HL_A, 1},
		0xAF: {XOR_A, 1},
		0xCB: {CB, 1},
	}
}

const (
	CB       = "CB"
	LD_SP    = "LD SP"
	LD_HL    = "LD_HL"
	LDD_HL_A = "LDD_(HL)_A"
	XOR_A    = "XOR A"
)
