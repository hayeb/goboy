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

func (bytereg ByteRegister) val() uint8 {
	return uint8(bytereg)
}

func (hwreg HalfWordRegister) val() uint16 {
	return uint16(hwreg)
}

type Register struct {
	A ByteRegister
	B ByteRegister
	C ByteRegister
	D ByteRegister
	E ByteRegister
	F ByteRegister
	H ByteRegister
	L ByteRegister

	SP HalfWordRegister
	PC HalfWordRegister

	Flag flags
}

type duoRegister int

const (
	REG_AF duoRegister = iota
	REG_BC
	REG_DE
	REG_HL
)

func (reg Register) duoRegs(duo duoRegister) (ByteRegister, ByteRegister) {
	switch duo {
	case REG_AF:
		return reg.A, reg.F
	case REG_BC:
		return reg.B, reg.C
	case REG_DE:
		return reg.D, reg.E
	case REG_HL:
		return reg.H, reg.L
	default:
		panic(fmt.Sprintf("Unknown duo register %d", duo))
	}
}

func duoRegisterValue(left ByteRegister, right ByteRegister) uint16 {
	return uint16(left.val())<<8 | uint16(right.val())
}

func (reg Register) readDuo(duo duoRegister) uint16 {
	l, r := reg.duoRegs(duo)
	return uint16(l)<<8 | uint16(r)
}

func (reg *Register) decrDuo(duo duoRegister) {
	l, r := reg.duoRegs(duo)
	reg.writeDuo(duo, duoRegisterValue(l, r)-1)
}

func (reg *Register) incrDuo(duo duoRegister) {
	l, r := reg.duoRegs(duo)
	reg.writeDuo(duo, duoRegisterValue(l, r)+1)
}

func (reg *Register) writeDuo(duo duoRegister, val uint16) {
	left := uint8(val >> 8)
	right := uint8(val)
	switch duo {
	case REG_AF:
		reg.A = ByteRegister(left)
		reg.F = ByteRegister(right)
	case REG_BC:
		reg.B = ByteRegister(left)
		reg.C = ByteRegister(right)
	case REG_DE:
		reg.D = ByteRegister(left)
		reg.E = ByteRegister(right)
	case REG_HL:
		reg.H = ByteRegister(left)
		reg.L = ByteRegister(right)
	}
}
