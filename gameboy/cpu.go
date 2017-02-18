package gameboy

import (
	"fmt"
)

func Run(cart []uint8, bootrom []uint8) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("There was an error: %s", r)
		}
	}()
	ci, mem, regs, instrMap := initializeSystem(cart, bootrom)
	if ci.ramSize != ram_none || ci.romSize != rom_kbit_256 {
		panic("Cartridge not supported")
	}

	for {
		instructionCode := mem.read8(regs.PC.val())
		instr := (*instrMap)[instructionCode]

		if instr == (instruction{}) {
			panic(fmt.Sprintf("Unrecognized instruction %x", instructionCode))
		}

		fmt.Printf("%#08x \t %s\n", regs.PC.val(), instr.name)
		switch instr.name {
		case jr:
			n := mem.read8(regs.PC.val() + 1)
			if !regs.Flag.Z {
				regs.PC = regs.PC + halfWordRegister(n)
			}
		case ld_sp:
			arg := mem.read16(regs.PC.val() + 1)
			regs.SP = halfWordRegister(arg)
		case ld_hl:
			regs.writeDuo(reg_hl, mem.read16(regs.PC.val()+1))
		case ldd_hl_a:
			mem.write8(regs.readDuo(reg_hl), regs.A.val())
			regs.decrDuo(reg_hl)
		case xor_a:
			regs.A = regs.A ^ regs.A
			if regs.A == 0 {
				regs.Flag.Z = true
			}
		case cb:
			nb := mem.read8(regs.PC.val() + 1)
			cbInstruction(mem, regs, nb)
		default:
			panic(fmt.Sprintf("Instuction not implemented: %s", instr.name))
		}

		regs.PC = halfWordRegister(regs.PC.val() + uint16(instr.bytes))

	}
}

func cbInstruction(mem *memory, regs *register, cbCode uint8) {
	switch cbCode {
	case 0x7c:
		// BIT 7, H (Check bit 7 in register H, if 0, set Z if 0
		// 0b01000000 >> 7
		t := regs.H.val() >> 7 & 0x1
		if t == 0x0 {
			regs.Flag.Z = true
		}
		regs.Flag.N = false
		regs.Flag.H = true
		regs.PC++
	default:
		panic(fmt.Sprintf("Unknown cb instruction %x", cbCode))
	}
}

func initializeSystem(cart []uint8, bootrom []uint8) (*cartridgeInfo, *memory, *register, *map[uint8]instruction) {
	cartridgeInfo := createCartridgeInfo(cart)
	instructionMap := createInstructionMap()
	mem := memInit(bootrom)
	registers := new(register)
	return cartridgeInfo, mem, registers, instructionMap
}
