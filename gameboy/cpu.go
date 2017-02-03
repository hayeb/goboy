package gameboy

import (
	"fmt"
	"io"
	"github.com/davecgh/go-spew/spew"
)

const (
	CB = "CB"
	LD_SP = "LD SP"
	LD_HL = "LD_HL"
	LDD_HL_A = "LDD_(HL)_A"
	XOR_A = "XOR A"
)

type duoRegister int
const (
	REG_AF duoRegister= iota
	REG_BC
	REG_DE
	REG_HL
)

func createInstructionMap() *map[uint8]instruction {
	return &map[uint8]instruction{
		0x21: {LD_HL, 3},
		0x31: {LD_SP, 3},
		0x32: {LDD_HL_A, 1},
		0xAF: {XOR_A, 1},
		0xCB: {CB, 1},
	}
}

type instruction struct {
	name  string
	bytes int
}

func (instr *instruction) Write(w io.Writer) {
	str := fmt.Sprintf("%s %d\n", instr.name, instr.bytes)
	w.Write([]byte(str))
}

type flags struct {
	Z bool
	N bool
	H bool
	C bool
}

type register struct {
	A    uint8
	B    uint8
	C    uint8
	D    uint8
	E    uint8
	F    uint8
	H    uint8
	L    uint8

	SP   uint16
	PC   uint16

	Flag flags
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

func (reg register) readDuo(duo duoRegister) uint16 {
	switch (duo) {
	case REG_AF:
		return uint16(reg.A) << 8 | uint16(reg.F)
	case REG_BC:
		return uint16(reg.B) << 8 | uint16(reg.C)
	case REG_DE:
		return uint16(reg.D) << 8 | uint16(reg.E)
	case REG_HL:
		return uint16(reg.H) << 8 | uint16(reg.L)
	default: panic(fmt.Sprintf("attempt to write to register duo %d", duo ))
	}
}

func (reg *register) decrDuo(duo duoRegister) {
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

func initializeSystem(cartridge []uint8, bootrom []uint8) (_ *CartridgeInfo, _ *Memory, _ *register, _ *map[uint8]instruction) {
	cartridgeInfo := CreateCartridgeInfo(cartridge)
	instructionMap := createInstructionMap()
	memory := memInit(bootrom)
	registers := new(register)
	return cartridgeInfo, memory, registers, instructionMap
}

func cbInstruction(cbCode uint8) {
	switch cbCode {
	case 0x7c:
		fmt.Println("Execute 0x7c instruction")
		panic("CB instruction 0x7c not yet implemented")
	default:
		panic(fmt.Sprintf("Unknown cb instruction %x", cbCode))
	}
}

func Run(bootrom []byte, cartridge []byte) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Encountered an error: %s", r)
		}
	}()

	ci, mem, regs, instrMap := initializeSystem(cartridge, bootrom)
	fmt.Println(CartridgeInfoString(*ci))

	if ci.RAMSize != RAM_NONE || ci.ROMSize != ROM_KBIT_256 {
		panic("Cartridge not supported")
	}

	for {
		spew.Dump(regs)
		instructionCode := mem.Read8(regs.PC)
		instr := (*instrMap)[instructionCode]

		if instr == (instruction{}) {
			panic(fmt.Sprintf("Unrecognized instruction %x", instructionCode));
		}

		switch instr.name {
		case LD_SP:
			arg := mem.Read16(regs.PC + 1)
			regs.SP = arg
		case LD_HL:
			regs.writeDuo(REG_HL, mem.Read16(regs.PC + 1))
		case LDD_HL_A:
			mem.Write8(regs.readDuo(REG_HL), regs.A)
			regs.decrDuo(REG_HL)
		case XOR_A:
			result := regs.A ^ regs.A
			regs.A = result
			if result == 0 {
				regs.Flag.Z = true
			}
		case CB:
			// Read the next byte, handle cb instruction
			nb := mem.Read8(regs.PC + 1)
			fmt.Printf("Execute DB instruction %x\n" , nb)
			// TODO: Correct interface
			cbInstruction(nb )
		default: panic(fmt.Sprintf("Instuction not implemented: %s", instr.name))
		}
		fmt.Printf("Instr: %s\n", instr.name)
		regs.PC += uint16(instr.bytes)

	}

}
