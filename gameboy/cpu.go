package gameboy

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
)

func initializeSystem(cartridge []uint8, bootrom []uint8) (_ *CartridgeInfo, _ *Memory, _ *register, _ *map[uint8]instruction) {
	cartridgeInfo := CreateCartridgeInfo(cartridge)
	instructionMap := CreateInstructionMap()
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
			nb := mem.Read8(regs.PC + 1)
			// TODO: Correct interface
			cbInstruction(nb )
		default: panic(fmt.Sprintf("Instuction not implemented: %s", instr.name))
		}
		fmt.Printf("Instr: %s\n", instr.name)
		regs.PC += uint16(instr.bytes)

	}

}
