package gameboy

import "fmt"

type register struct {
	A uint8
	B uint8
	C uint8
	D uint8
	E uint8
	F uint8
	H uint8
	L uint8

	SP uint16
	PC uint16
}

type duoRegister int

const (
	REG_AF duoRegister = iota
	REG_BC
	REG_DE
	REG_HL
)

func (reg *register) setZ(val bool) {
	if val {
		reg.F = setBit(reg.F, 7)
	} else {
		reg.F = resetBit(reg.F, 7)
	}
}

func (reg *register) isZ() bool {
	return (reg.F>>7)&0x1 == 1
}

func (reg *register) setN(val bool) {
	if val {
		reg.F = setBit(reg.F, 6)
	} else {
		reg.F = resetBit(reg.F, 6)
	}
}

func (reg *register) isN() bool {
	return (reg.F>>6)&0x1 == 1
}

func (reg *register) setH(val bool) {
	if val {
		reg.F = setBit(reg.F, 5)
	} else {
		reg.F = resetBit(reg.F, 5)
	}
}

func (reg *register) isH() bool {
	return (reg.F>>5)&0x1 == 1
}

func (reg *register) setC(val bool) {
	if val {
		reg.F = setBit(reg.F, 4)
	} else {
		reg.F = resetBit(reg.F, 4)
	}
}

func (reg *register) isC() bool {
	return (reg.F>>4)&0x1 == 1
}

func (reg *register) duoRegs(duo duoRegister) (uint8, uint8) {
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

func duoRegisterValue(left uint8, right uint8) uint16 {
	return uint16(left)<<8 | uint16(right)
}

func (reg *register) readDuo(duo duoRegister) uint16 {
	l, r := reg.duoRegs(duo)
	return uint16(l)<<8 | uint16(r)
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
	right := uint8(val & 0xff)
	switch duo {
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

func (reg *register) bit(bit int, val uint8) {
	cond := val & (1 << uint(bit))
	reg.setZ(cond == 0)
	reg.setN(false)
	reg.setH(true)
}

func (reg *register) incSP(n int) {
	reg.SP = reg.SP + uint16(n)
}

func (reg *register) decSP(n int) {
	reg.SP = reg.SP - uint16(n)
}

func (reg *register) incPC(n int) {
	reg.PC = reg.PC + uint16(n)
}

func (reg *register) decPC(n int) {
	reg.PC = reg.PC - uint16(n)
}
