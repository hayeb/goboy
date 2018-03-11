package gameboy

type cbInstruction struct {
	name  string
	bytes int

	actionDuration int
	noopDuration   int

	executor cbInstructionExecutor
}

type cbInstructionExecutor func(mem *memory, reg *register, cbInstr *cbInstruction) int

func createCBInstructionMap() *map[uint8]*cbInstruction {
	return &map[uint8]*cbInstruction{
		0x11: newCBInstruction("RL C", 2, 8, rl_c),
		0x1a: newCBInstruction("RR D", 2, 8, rrD),
		0x1b: newCBInstruction("RR E", 2, 8, rrE),
		0x19: newCBInstruction("RR C", 2, 8, rrC),

		0x26: newCBInstruction("SLA (HL)", 2, 16, slaHL),
		0x27: newCBInstruction("SLA A", 2, 8, slaA),

		0x30: newCBInstruction("SWAP B", 2, 8, swapB),
		0x33: newCBInstruction("SWAP E", 2, 8, swapE),
		0x37: newCBInstruction("SWAP A", 2, 8, swapA),
		0x38: newCBInstruction("SRL B", 2, 8, srlB),
		0x3a: newCBInstruction("SRL D", 2, 8, srlD),
		0x3b: newCBInstruction("SRL E", 2, 8, srlE),
		0x3f: newCBInstruction("SRL A", 2, 8, srlA),

		0x40: newCBInstruction("BIT 0,B", 2, 8, bit_0_b),
		0x41: newCBInstruction("BIT 0,C", 2, 8, bit_0_c),
		0x47: newCBInstruction("BIT 0,A", 2, 8, bit_0_a),
		0x48: newCBInstruction("BIT 1,B", 2, 8, bit_1_b),
		0x4e: newCBInstruction("BIT 1,(HL)", 2, 8, bit_1_hl),
		0x4f: newCBInstruction("BIT 1,A", 2, 8, bit_1_a),

		0x50: newCBInstruction("BIT 2,B", 2, 8, bit_2_b),
		0x57: newCBInstruction("BIT 2,A", 2, 8, bit_2_a),
		0x58: newCBInstruction("BIT 3,B", 2, 8, bit_3_b),
		0x5f: newCBInstruction("BIT 3,A", 2, 8, bit_3_a),

		0x60: newCBInstruction("BIT 4,B", 2, 8, bit_4_b),
		0x61: newCBInstruction("BIT 4,C", 2, 8, bit_4_c),
		0x68: newCBInstruction("BIT 5,B", 2, 8, bit_5_b),
		0x69: newCBInstruction("BIT 5,C", 2, 8, bit_5_c),
		0x6e: newCBInstruction("BIt 5,D", 2, 8, bit_5_d),
		0x6f: newCBInstruction("BIT 5,A", 2, 8, bit_5_a),

		0x70: newCBInstruction("BIT 6,B", 2, 8, bit_6_b),
		0x71: newCBInstruction("BIT 6,C", 2, 8, bit_6_c),
		0x74: newCBInstruction("BIT 6,H", 2, 8, bit_6_h),
		0x76: newCBInstruction("BIT 6,(HL)", 2, 8, bit_6_hl),
		0x77: newCBInstruction("BIT 6,A", 2, 8, bit_6_a),
		0x78: newCBInstruction("BIT 7,B", 2, 8, bit_7_b),
		0x79: newCBInstruction("BIT 7,C", 2, 8, bit_7_c),
		0x7a: newCBInstruction("BIT 7,D", 2, 8, bit_7_d),
		0x7b: newCBInstruction("BIT 7,E", 2, 8, bit_7_e),
		0x7c: newCBInstruction("BIT 7,H", 2, 8, bit_7_h),
		0x7d: newCBInstruction("BIT 7,L", 2, 8, bit_7_l),
		0x7e: newCBInstruction("BIT 7,(HL)", 2, 16, bit_7_hl),
		0x7f: newCBInstruction("BIT 7,A", 2, 8, bit_7_a),

		0x86: newCBInstruction("RES 0,(HL)", 2, 16, res0hl),
		0x87: newCBInstruction("RES 0,A", 2, 8, res0A),
		0x8e: newCBInstruction("RES 1,(HL)", 2, 16, res1hl),

		0x96: newCBInstruction("RES 2,(HL)", 2, 16, res2hl),
		0x9e: newCBInstruction("RES 3,(HL)", 2, 8, res3hl),

		0xae: newCBInstruction("RES 5,(HL)", 2, 16, res5hl),

		0xb6: newCBInstruction("RES 6,(HL)", 2, 16, res6hl),
		0xbe: newCBInstruction("RES 7,(HL)", 2, 8, res7hl),

		0xc6: newCBInstruction("SET 0,(HL)", 2, 16, set0hl),
		0xc7: newCBInstruction("SET 0,A", 2, 8, set0A),
		0xce: newCBInstruction("SET 1,(HL)", 2, 16, set1hl),
		0xcf: newCBInstruction("SET 1,A", 2, 8, set1A),

		0xd6: newCBInstruction("SET 2,(HL)", 2, 16, set2hl),
		0xd7: newCBInstruction("SET 2,A", 2, 8, set2A),
		0xde: newCBInstruction("SET 3,(HL)", 2, 16, set3hl),

		0xe7: newCBInstruction("SET 4,A", 2, 8, set4A),
		0xee: newCBInstruction("SET 5,(HL)", 2, 16, set5hl),

		0xf7: newCBInstruction("SET 6,A", 2, 8, set6A),
		0xfe: newCBInstruction("SET 7,(HL)", 2, 16, set_7_hl),
	}
}

func newCBInstruction(name string, length int, duration int, fp func(mem *memory, reg *register, cbInstr *cbInstruction) int) *cbInstruction {
	return newCBConditionalInstruction(name, length, duration, 0, fp)
}

func newCBConditionalInstruction(name string, length int, actionDuration int, noopDuration int, fp func(mem *memory, reg *register, cbInstr *cbInstruction) int) *cbInstruction {
	return &cbInstruction{
		name:           name,
		bytes:          length,
		actionDuration: actionDuration,
		noopDuration:   noopDuration,
		executor:       fp,
	}
}

func bit_7_h(_ *memory, reg *register, cbInstr *cbInstruction) int {
	reg.bit(7, reg.H)

	reg.incPC(cbInstr.bytes)
	return cbInstr.actionDuration
}

func bit_7_l(_ *memory, reg *register, cbInstr *cbInstruction) int {
	reg.bit(7, reg.L)

	reg.incPC(cbInstr.bytes)
	return cbInstr.actionDuration
}

func bit_0_a(_ *memory, reg *register, cbInstr *cbInstruction) int {
	reg.bit(0, reg.A)

	reg.incPC(cbInstr.bytes)
	return cbInstr.actionDuration
}

func bit_0_b(_ *memory, reg *register, cbInstr *cbInstruction) int {
	reg.bit(0, reg.B)

	reg.incPC(cbInstr.bytes)
	return cbInstr.actionDuration
}

func bit_0_c(_ *memory, reg *register, cbInstr *cbInstruction) int {
	reg.bit(0, reg.C)

	reg.incPC(cbInstr.bytes)
	return cbInstr.actionDuration
}

func bit_1_b(_ *memory, reg *register, cbInstr *cbInstruction) int {
	reg.bit(1, reg.B)

	reg.incPC(cbInstr.bytes)
	return cbInstr.actionDuration
}

func bit_1_hl(mem *memory, reg *register, cbInstr *cbInstruction) int {
	reg.bit(1, mem.read8(reg.readDuo(REG_HL)))

	reg.incPC(cbInstr.bytes)
	return cbInstr.actionDuration
}

func bit_1_a(_ *memory, reg *register, cbInstr *cbInstruction) int {
	reg.bit(1, reg.A)

	reg.incPC(cbInstr.bytes)
	return cbInstr.actionDuration
}

func bit_2_a(_ *memory, reg *register, cbInstr *cbInstruction) int {
	reg.bit(2, reg.A)

	reg.incPC(cbInstr.bytes)
	return cbInstr.actionDuration
}

func bit_2_b(_ *memory, reg *register, cbInstr *cbInstruction) int {
	reg.bit(2, reg.B)

	reg.incPC(cbInstr.bytes)
	return cbInstr.actionDuration
}

func bit_3_a(_ *memory, reg *register, cbInstr *cbInstruction) int {
	reg.bit(3, reg.A)

	reg.incPC(cbInstr.bytes)
	return cbInstr.actionDuration
}

func bit_3_b(_ *memory, reg *register, cbInstr *cbInstruction) int {
	reg.bit(3, reg.B)

	reg.incPC(cbInstr.bytes)
	return cbInstr.actionDuration
}

func bit_4_b(_ *memory, reg *register, cbInstr *cbInstruction) int {
	reg.bit(4, reg.B)

	reg.incPC(cbInstr.bytes)
	return cbInstr.actionDuration
}

func bit_4_c(_ *memory, reg *register, cbInstr *cbInstruction) int {
	reg.bit(4, reg.C)

	reg.incPC(cbInstr.bytes)
	return cbInstr.actionDuration
}

func bit_5_a(_ *memory, reg *register, cbInstr *cbInstruction) int {
	reg.bit(5, reg.A)

	reg.incPC(cbInstr.bytes)
	return cbInstr.actionDuration
}

func bit_6_b(_ *memory, reg *register, cbInstr *cbInstruction) int {
	reg.bit(6, reg.B)

	reg.incPC(cbInstr.bytes)
	return cbInstr.actionDuration
}

func bit_6_c(_ *memory, reg *register, cbInstr *cbInstruction) int {
	reg.bit(6, reg.B)

	reg.incPC(cbInstr.bytes)
	return cbInstr.actionDuration
}

func bit_6_h(_ *memory, reg *register, cbInstr *cbInstruction) int {
	reg.bit(6, reg.H)

	reg.incPC(cbInstr.bytes)
	return cbInstr.actionDuration
}

func bit_6_hl(mem *memory, reg *register, cbInstr *cbInstruction) int {
	reg.bit(6, mem.read8(reg.readDuo(REG_HL)))

	reg.incPC(cbInstr.bytes)
	return cbInstr.actionDuration
}

func bit_5_b(_ *memory, reg *register, cbInstr *cbInstruction) int {
	reg.bit(5, reg.B)

	reg.incPC(cbInstr.bytes)
	return cbInstr.actionDuration
}

func bit_5_c(_ *memory, reg *register, cbInstr *cbInstruction) int {
	reg.bit(5, reg.C)

	reg.incPC(cbInstr.bytes)
	return cbInstr.actionDuration
}

func bit_5_d(_ *memory, reg *register, cbInstr *cbInstruction) int {
	reg.bit(5, reg.D)

	reg.incPC(cbInstr.bytes)
	return cbInstr.actionDuration
}

func bit_6_a(_ *memory, reg *register, cbInstr *cbInstruction) int {
	reg.bit(6, reg.A)

	reg.incPC(cbInstr.bytes)
	return cbInstr.actionDuration
}

func bit_7_b(_ *memory, reg *register, cbInstr *cbInstruction) int {
	reg.bit(7, reg.B)

	reg.incPC(cbInstr.bytes)
	return cbInstr.actionDuration
}

func bit_7_c(_ *memory, reg *register, cbInstr *cbInstruction) int {
	reg.bit(7, reg.C)

	reg.incPC(cbInstr.bytes)
	return cbInstr.actionDuration
}

func bit_7_d(_ *memory, reg *register, cbInstr *cbInstruction) int {
	reg.bit(7, reg.D)

	reg.incPC(cbInstr.bytes)
	return cbInstr.actionDuration
}

func bit_7_e(_ *memory, reg *register, cbInstr *cbInstruction) int {
	reg.bit(7, reg.E)

	reg.incPC(cbInstr.bytes)
	return cbInstr.actionDuration
}

func bit_7_a(_ *memory, reg *register, cbInstr *cbInstruction) int {
	reg.bit(7, reg.A)

	reg.incPC(cbInstr.bytes)
	return cbInstr.actionDuration
}

func bit_7_hl(mem *memory, reg *register, cbInstr *cbInstruction) int {
	reg.bit(7, mem.read8(reg.readDuo(REG_HL)))

	reg.incPC(cbInstr.bytes)
	return cbInstr.actionDuration
}

func rl_c(_ *memory, reg *register, cbInstr *cbInstruction) int {
	isCarrySet := reg.isC()
	isMSBSet := testBit(reg.C, 7)

	reg.setZ( false)
	reg.setN( false)
	reg.setH( false)

	newVal := reg.C << 1

	reg.setC( isMSBSet)

	if isCarrySet {
		newVal = setBit(newVal, 0)
	}

	if newVal == 0 {
		reg.setZ( true)
	}
	reg.C = newVal
	reg.incPC(cbInstr.bytes)
	return cbInstr.actionDuration
}

func res0A(_ *memory, reg *register, cbInstr *cbInstruction) int {
	reg.A = resetBit(reg.A, 0)
	reg.incPC(cbInstr.bytes)
	return cbInstr.actionDuration
}

func res0hl(mem *memory, reg *register, cbInstr *cbInstruction) int {
	mem.write8(reg.readDuo(REG_HL), resetBit(mem.read8(reg.readDuo(REG_HL)), 0))
	reg.incPC(cbInstr.bytes)
	return cbInstr.actionDuration
}

func res1hl(mem *memory, reg *register, cbInstr *cbInstruction) int {
	mem.write8(reg.readDuo(REG_HL), resetBit(mem.read8(reg.readDuo(REG_HL)), 1))
	reg.incPC(cbInstr.bytes)
	return cbInstr.actionDuration
}

func res2hl(mem *memory, reg *register, cbInstr *cbInstruction) int {
	mem.write8(reg.readDuo(REG_HL), resetBit(mem.read8(reg.readDuo(REG_HL)), 2))
	reg.incPC(cbInstr.bytes)
	return cbInstr.actionDuration
}

func res3hl(mem *memory, reg *register, cbInstr *cbInstruction) int {
	mem.write8(reg.readDuo(REG_HL), resetBit(mem.read8(reg.readDuo(REG_HL)), 3))
	reg.incPC(cbInstr.bytes)
	return cbInstr.actionDuration
}

func res5hl(mem *memory, reg *register, cbInstr *cbInstruction) int {
	mem.write8(reg.readDuo(REG_HL), resetBit(mem.read8(reg.readDuo(REG_HL)), 5))
	reg.incPC(cbInstr.bytes)
	return cbInstr.actionDuration
}

func res6hl(mem *memory, reg *register, cbInstr *cbInstruction) int {
	mem.write8(reg.readDuo(REG_HL), resetBit(mem.read8(reg.readDuo(REG_HL)), 6))
	reg.incPC(cbInstr.bytes)
	return cbInstr.actionDuration
}

func res7hl(mem *memory, reg *register, cbInstr *cbInstruction) int {
	mem.write8(reg.readDuo(REG_HL), resetBit(mem.read8(reg.readDuo(REG_HL)), 7))
	reg.incPC(cbInstr.bytes)
	return cbInstr.actionDuration
}

func swapA(_ *memory, reg *register, cbInstr *cbInstruction) int {
	val := reg.A

	reg.A = val<<4 | val>>4
	reg.setZ( reg.A == 0)
	reg.setN( false)
	reg.setH( false)
	reg.setC( false)
	reg.incPC(cbInstr.bytes)
	return cbInstr.actionDuration
}

func swapB(_ *memory, reg *register, cbInstr *cbInstruction) int {
	val := reg.B

	reg.B = val<<4 | val>>4
	reg.setZ( reg.B == 0)
	reg.setN( false)
	reg.setH( false)
	reg.setC( false)
	reg.incPC(cbInstr.bytes)
	return cbInstr.actionDuration
}

func swapE(_ *memory, reg *register, cbInstr *cbInstruction) int {
	val := reg.E

	reg.E = val<<4 | val>>4
	reg.setZ( reg.E == 0)
	reg.setN( false)
	reg.setH( false)
	reg.setC( false)
	reg.incPC(cbInstr.bytes)
	return cbInstr.actionDuration
}

func slaHL(mem *memory, reg *register, cbInstr *cbInstruction) int {
	val := mem.read8(reg.readDuo(REG_HL))
	reg.setC( val>>7 == 1)
	reg.A = val << 1

	reg.setZ( reg.A == 0)
	reg.setN( false)
	reg.setH( false)
	reg.incPC(cbInstr.bytes)
	return cbInstr.actionDuration
}

func slaA(_ *memory, reg *register, cbInstr *cbInstruction) int {
	reg.setC( reg.A>>7 == 1)
	reg.A = reg.A << 1

	reg.setZ( reg.A == 0)
	reg.setN( false)
	reg.setH( false)
	reg.incPC(cbInstr.bytes)
	return cbInstr.actionDuration
}

func set1hl(mem *memory, reg *register, cbInstr *cbInstruction) int {
	mem.write8(reg.readDuo(REG_HL), setBit(mem.read8(reg.readDuo(REG_HL)), 1))
	reg.incPC(cbInstr.bytes)
	return cbInstr.actionDuration
}

func set0A(_ *memory, reg *register, cbInstr *cbInstruction) int {
	reg.A = setBit(reg.A, 0)
	reg.incPC(cbInstr.bytes)
	return cbInstr.actionDuration
}

func set1A(_ *memory, reg *register, cbInstr *cbInstruction) int {
	reg.A = setBit(reg.A, 1)
	reg.incPC(cbInstr.bytes)
	return cbInstr.actionDuration
}

func set2A(_ *memory, reg *register, cbInstr *cbInstruction) int {
	reg.A = setBit(reg.A, 2)
	reg.incPC(cbInstr.bytes)
	return cbInstr.actionDuration
}

func set4A(_ *memory, reg *register, cbInstr *cbInstruction) int {
	reg.A = setBit(reg.A, 4)
	reg.incPC(cbInstr.bytes)
	return cbInstr.actionDuration
}

func set6A(_ *memory, reg *register, cbInstr *cbInstruction) int {
	reg.A = setBit(reg.A, 6)
	reg.incPC(cbInstr.bytes)
	return cbInstr.actionDuration
}

func set_7_hl(mem *memory, reg *register, cbInstr *cbInstruction) int {
	address := reg.readDuo(REG_HL)
	val := mem.read8(address)
	mem.write8(address, setBit(val, 7))
	reg.incPC(cbInstr.bytes)
	return cbInstr.actionDuration
}

func set0hl(mem *memory, reg *register, cbInstr *cbInstruction) int {
	address := reg.readDuo(REG_HL)
	val := mem.read8(address)
	mem.write8(address, setBit(val, 0))
	reg.incPC(cbInstr.bytes)
	return cbInstr.actionDuration
}

func set2hl(mem *memory, reg *register, cbInstr *cbInstruction) int {
	address := reg.readDuo(REG_HL)
	val := mem.read8(address)
	mem.write8(address, setBit(val, 2))
	reg.incPC(cbInstr.bytes)
	return cbInstr.actionDuration
}

func set3hl(mem *memory, reg *register, cbInstr *cbInstruction) int {
	address := reg.readDuo(REG_HL)
	val := mem.read8(address)
	mem.write8(address, setBit(val, 3))
	reg.incPC(cbInstr.bytes)
	return cbInstr.actionDuration
}

func set5hl(mem *memory, reg *register, cbInstr *cbInstruction) int {
	address := reg.readDuo(REG_HL)
	val := mem.read8(address)
	mem.write8(address, setBit(val, 5))
	reg.incPC(cbInstr.bytes)
	return cbInstr.actionDuration
}

func srlA(_ *memory, reg *register, cbInstr *cbInstruction) int {
	reg.setC( reg.A&0x1 == 1)
	reg.A >>= 1

	reg.setZ( reg.A == 0)
	reg.setN( false)
	reg.setH( false)
	reg.incPC(cbInstr.bytes)
	return cbInstr.actionDuration
}

func srlB(_ *memory, reg *register, cbInstr *cbInstruction) int {
	reg.setC( reg.B&0x1 == 1)
	reg.B >>= 1

	reg.setZ( reg.B == 0)
	reg.setN( false)
	reg.setH( false)
	reg.incPC(cbInstr.bytes)
	return cbInstr.actionDuration
}

func srlD(_ *memory, reg *register, cbInstr *cbInstruction) int {
	reg.setC( reg.D&0x1 == 1)
	reg.D >>= 1

	reg.setZ( reg.D == 0)
	reg.setN( false)
	reg.setH( false)
	reg.incPC(cbInstr.bytes)
	return cbInstr.actionDuration
}


func srlE(_ *memory, reg *register, cbInstr *cbInstruction) int {
	reg.setC( reg.E&0x1 == 1)
	reg.E >>= 1

	reg.setZ( reg.E == 0)
	reg.setN( false)
	reg.setH( false)
	reg.incPC(cbInstr.bytes)
	return cbInstr.actionDuration
}

func rrC(_ *memory, reg *register, _ *cbInstruction) int {
	return rotateRight(&reg.C, reg)
}

func rrD(_ *memory, reg *register, _ *cbInstruction) int {
	return rotateRight(&reg.D, reg)
}

func rrE(_ *memory, reg *register, _ *cbInstruction) int {
	return rotateRight(&reg.E, reg)
}

func rotateRight(r *uint8, reg *register) int {
	val := *r
	carry := reg.isC()

	reg.C = val >> 1
	if carry {
		setBit(reg.C, 7)
	}

	reg.setZ( reg.C == 0)
	reg.setN( false)
	reg.setH( false)
	reg.setC( reg.C&0x1 == 1)
	reg.incPC(2)
	return 8
}
