package gameboy

import (
	"fmt"
)

func Run(ci *CartridgeInfo, mem *Memory, regs *Register, instrMap *map[uint8]Instruction) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Encountered an error: %s", r)
		}
	}()

	fmt.Println(CartridgeInfoString(*ci))

	if ci.RAMSize != RAM_NONE || ci.ROMSize != ROM_KBIT_256 {
		panic("Cartridge not supported")
	}

	for {
		instructionCode := mem.Read8(regs.PC.val())
		instr := (*instrMap)[instructionCode]

		if instr == (Instruction{}) {
			panic(fmt.Sprintf("Unrecognized instruction %x", instructionCode))
		}

		switch instr.name {
		case LD_SP:
			arg := mem.Read16(regs.PC.val() + 1)
			regs.SP = HalfWordRegister(arg)
		case LD_HL:
			regs.writeDuo(REG_HL, mem.Read16(regs.PC.val()+1))
		case LDD_HL_A:
			mem.Write8(regs.readDuo(REG_HL), regs.A.val())
			regs.decrDuo(REG_HL)
		case XOR_A:
			regs.A = regs.A ^ regs.A
			if regs.A == 0 {
				regs.Flag.Z = true
			}
		case CB:
			nb := mem.Read8(regs.PC.val() + 1)
			cbInstruction(nb)
		default:
			panic(fmt.Sprintf("Instuction not implemented: %s", instr.name))
		}
		fmt.Printf("Instr: %s\n", instr.name)
		regs.PC = HalfWordRegister(regs.PC.val() + uint16(instr.bytes))

	}

}
