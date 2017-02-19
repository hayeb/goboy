package gameboy

import "fmt"

type byteRegister uint8
type halfWordRegister uint16
type flags struct {
	Z bool
	N bool
	H bool
	C bool
}

func (bytereg byteRegister) val() uint8 {
	return uint8(bytereg)
}

func (hwreg halfWordRegister) val() uint16 {
	return uint16(hwreg)
}

type register struct {
	A byteRegister
	B byteRegister
	C byteRegister
	D byteRegister
	E byteRegister
	F byteRegister
	H byteRegister
	L byteRegister

	SP halfWordRegister
	PC halfWordRegister

	Flag flags
}

type duoRegister int

const (
	reg_af duoRegister = iota
	reg_bc
	reg_de
	reg_hl
)

func (reg *register) duoRegs(duo duoRegister) (byteRegister, byteRegister) {
	switch duo {
	case reg_af:
		return reg.A, reg.F
	case reg_bc:
		return reg.B, reg.C
	case reg_de:
		return reg.D, reg.E
	case reg_hl:
		return reg.H, reg.L
	default:
		panic(fmt.Sprintf("Unknown duo register %d", duo))
	}
}

func duoRegisterValue(left byteRegister, right byteRegister) uint16 {
	return uint16(left.val())<<8 | uint16(right.val())
}

func (reg *register) readDuo(duo duoRegister) uint16 {
	l, r := reg.duoRegs(duo)
	return uint16(l.val())<<8 | uint16(r.val())
}

func (reg *register) decrDuo(duo duoRegister) {
	l, r := reg.duoRegs(duo)
	reg.writeDuo(duo, duoRegisterValue(l, r)-1)
}

func (reg *register) incrDuo(duo duoRegister) {
	l, r := reg.duoRegs(duo)
	reg.writeDuo(duo, duoRegisterValue(l, r)+1)
}

func (reg *register) writeDuo(duo duoRegister, val uint16) {
	left := uint8(val >> 8)
	right := uint8(val)
	switch duo {
	case reg_af:
		reg.A = byteRegister(left)
		reg.F = byteRegister(right)
	case reg_bc:
		reg.B = byteRegister(left)
		reg.C = byteRegister(right)
	case reg_de:
		reg.D = byteRegister(left)
		reg.E = byteRegister(right)
	case reg_hl:
		reg.H = byteRegister(left)
		reg.L = byteRegister(right)
	}
}

func (reg *register) bit(bit uint8, val uint8) {
	cond := val & bit
	if cond != 0 {
		reg.Flag.Z = false
	} else {
		reg.Flag.Z = true
	}
	reg.Flag.N = false
	reg.Flag.H = true
}

func (reg *register) incSP(n int) {
	reg.SP = halfWordRegister(reg.SP.val() + uint16(n))
}

func (reg *register) decSP(n int) {
	reg.SP = halfWordRegister(reg.SP.val() - uint16(n))
}
