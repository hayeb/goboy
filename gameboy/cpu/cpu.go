package cpu

import (
	"fmt"
	"github.com/hayeb/goboy/gameboy/cartridge"
	"github.com/hayeb/goboy/gameboy/memory"
)

func Run(cart []uint8, bootrom []uint8  ) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("There was an error: %s", r)
		}
	}();
	ci, mem, regs, instrMap := initializeSystem(cart, bootrom)
	if ci.RAMSize != cartridge.RAM_NONE || ci.ROMSize != cartridge.ROM_KBIT_256 {
		panic("Cartridge not supported")
	}

	for {
		instructionCode := mem.Read8(regs.PC.val())
		instr := (*instrMap)[instructionCode]

		if instr == (Instruction{}) {
			panic(fmt.Sprintf("Unrecognized instruction %x", instructionCode))
		}

		fmt.Printf("Instr: %s\n", instr.name)
		switch instr.name {
		case JR:
			n := mem.Read8(regs.PC.val() + 1)
			if !regs.Flag.Z {
				regs.PC = regs.PC + HalfWordRegister(n)
			}
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
			cbInstruction(mem, regs, nb)
		default:
			panic(fmt.Sprintf("Instuction not implemented: %s", instr.name))
		}

		regs.PC = HalfWordRegister(regs.PC.val() + uint16(instr.bytes))

	}
}

func cbInstruction(mem *memory.Memory, regs *Register, cbCode uint8) {
	switch cbCode {
	case 0x7c:
		// BIT 7, H (Check bit 7 in register H, if 0, set Z if 0
		// 0b01000000 >> 7
		if regs.H.val() >> 7 == 0 {
			regs.Flag.Z = true
		} else {
			regs.Flag.Z = false
		}
		regs.Flag.N = false
		regs.Flag.H = true
		regs.PC++
	default:
		panic(fmt.Sprintf("Unknown cb instruction %x", cbCode))
	}
}

func initializeSystem(cart []uint8, bootrom []uint8) (*cartridge.CartridgeInfo, *memory.Memory, *Register, *map[uint8]Instruction) {
	cartridgeInfo := cartridge.CreateCartridgeInfo(cart)
	instructionMap := CreateInstructionMap()
	mem := memory.MemInit(bootrom)
	registers := new(Register)
	return cartridgeInfo, mem, registers, instructionMap
}