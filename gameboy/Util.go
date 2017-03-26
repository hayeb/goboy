package gameboy

func rLeftCarry(val uint8, c bool) (uint8, bool) {
	return val<<1 | uint8(btoi(c)), uint8tob(val >> 7)
}

func rRightCarry(val uint8, c bool) (uint8, bool) {
	return val>>1 | uint8(btoi(c)<<7), uint8tob(val >> 7)
}

func rLeft(val uint8) (uint8, bool) {
	return val<<1 | val>>7, uint8tob(val >> 7)
}

func rRight(val uint8) (uint8, bool) {
	return val>>1 | val<<7, uint8tob(val << 7)
}

func btoi(b bool) int {
	if b {
		return 1
	} else {
		return 0
	}
}

func uint8tob(i uint8) bool {
	return i != 0
}

func setBit(bit uint8, n uint) uint8 {
	return bit | (1 << n)
}

func resetBit(bit uint8, n uint) uint8 {
	return bit &^ (1 << n)
}

func testBit(bit uint8, n uint) bool {
	return bit & (1 << n) != 0
}

func getBitN(but uint8, n uint) uint8 {
	val := but & (1 << n)
	if val == 0 {
		return 0
	} else {
		return 1
	}
}
