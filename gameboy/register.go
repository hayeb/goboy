package gameboy

import "fmt"

type ByteRegister uint8
type HalfWordRegister uint16
type flags struct {
	Z bool
	N bool
	H bool
	C bool
}

type register struct {
	A    ByteRegister
	B    ByteRegister
	C    ByteRegister
	D    ByteRegister
	E    ByteRegister
	F    ByteRegister
	H    ByteRegister
	L    ByteRegister

	SP   HalfWordRegister
	PC   HalfWordRegister

	Flag flags
}

type duoRegister int
const (
	REG_AF duoRegister= iota
	REG_BC
	REG_DE
	REG_HL
)

func (reg register) duoRegs(duo duoRegister) (ByteRegister, ByteRegister) {
	switch (duo) {
	case REG_AF: return reg.A, reg.F
	case REG_BC: return reg.B, reg.C
	case REG_DE: return reg.D, reg.E
	case REG_HL: return reg.H, reg.L
	default: panic(fmt.Sprintf("Unknown duo register %d", duo))
	}
}

func (reg register) readDuo(duo duoRegister) uint16 {
	l, r := reg.duoRegs(duo)
	return uint16(l) << 8 | uint16(r)
}

func (reg *register) decrDuo(duo duoRegister) {
	l, r := reg.duoRegs(duo)
	reg.writeDuo(REG_AF , reg.readDuo(duo) - 1)
	switch (duo) {
	case REG_AF:
		reg.writeDuo(REG_AF, reg.readDuo(duo) - 1)
	case REG_BC:
		reg.writeDuo(REG_BC, reg.readDuo(duo) - 1)
	case REG_DE:
		reg.writeDuo(REG_DE, reg.readDuo(duo) - 1)
	case REG_HL:
		reg.writeDuo(REG_HL, reg.readDuo(duo) - 1)
	default: panic(fmt.Sprintf("Attempt to decrease duo %d", duo ))
	}
}

func (reg *register) incrDuo(duo duoRegister) {
	switch (duo) {
	case REG_AF:
		reg.writeDuo(REG_AF, reg.readDuo(duo) + 1)
	case REG_BC:
		reg.writeDuo(REG_BC, reg.readDuo(duo) + 1)
	case REG_DE:
		reg.writeDuo(REG_DE, reg.readDuo(duo) + 1)
	case REG_HL:
		reg.writeDuo(REG_HL, reg.readDuo(duo) + 1)
	default: panic(fmt.Sprintf("Attempt to increase duo %d", duo ))
	}
}

func (reg *register) writeDuo(duo duoRegister, val uint16 ) {
	left := uint8(val >> 8)
	right := uint8(val)
	switch (duo) {
	case REG_AF:
		reg.A = left
		reg.F = right
	case REG_BC:
		reg.B = left
		reg.C = right
	case REG_DE:
		reg.D = left
		reg.E = right
	case REG_HL:
		reg.H = left
		reg.L = right
	}
}
