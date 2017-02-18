package gameboy

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
)

var _ = spew.Config

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
			panic(fmt.Sprintf("Unrecognized instruction %x at address %#04x", instructionCode, regs.PC.val()))
		}

		fmt.Printf("%#04x \t %s\n", regs.PC.val(), instr.name)
		switch instr.name {
		case cb:
			nb := readArgByte(mem, regs, 1)
			cbInstruction(mem, regs, nb)
		case jr:
			n := int8(readArgByte(mem, regs, 1))
			if !regs.Flag.Z {
				regs.PC = halfWordRegister(int(regs.PC.val()) + int(n))
			}
		case inc_c:
			incRegister8(&regs.C)
		case ld_a:
			regs.A = byteRegister(readArgByte(mem, regs, 1))
		case ld_c:
			regs.C = byteRegister(readArgByte(mem, regs, 1))
		case ld_C_a:
			mem.write8(0xFF00+uint16(regs.C.val()), regs.A.val())
		case ld_sp:
			arg := readArgByte(mem, regs, 1)
			regs.SP = halfWordRegister(arg)
		case ld_hl:
			regs.writeDuo(reg_hl, readArgHalfword(mem, regs, 1))
		case ld_HL_a:
			mem.write8(mem.read16(regs.readDuo(reg_hl)), readArgByte(mem, regs, 1))
		case ldd_hl_a:
			mem.write8(regs.readDuo(reg_hl), regs.A.val())
			regs.decrDuo(reg_hl)
		case ldh_a8_A:
			// TODO: A is hierna leeg? Uitzoeken!
			regs.A = byteRegister(mem.read8(uint16(readArgByte(mem, regs, 1)) + 0xff00))
		case xor_a:
			regs.A = regs.A ^ regs.A
			if regs.A == 0 {
				regs.Flag.Z = true
			}
		default:
			panic(fmt.Sprintf("Instuction not implemented: %s", instr.name))
		}

		spew.Dump(regs)
		//fmt.Printf("REG AF: %#08x\n", regs.readDuo(reg_af))
		//fmt.Printf("REG BC: %#08x\n", regs.readDuo(reg_bc))
		//fmt.Printf("REG DE: %#08x\n", regs.readDuo(reg_de))
		//fmt.Printf("REG HL: %#08x\n", regs.readDuo(reg_hl))

		regs.PC = halfWordRegister(regs.PC.val() + uint16(instr.bytes))
	}
}

func incRegister8(reg *byteRegister) {
	*reg = byteRegister(reg.val() + 1)
}

func decrRegister8(reg *byteRegister) {
	*reg = byteRegister(reg.val() - 1)
}

func readArgByte(mem *memory, reg *register, arg int) uint8 {
	return mem.read8(reg.PC.val() + uint16(arg))
}

func readArgHalfword(mem *memory, reg *register, arg int) uint16 {
	return mem.read16(reg.PC.val() + uint16(arg))
}

func cbInstruction(mem *memory, regs *register, cbCode uint8) {
	switch cbCode {
	case 0x7c:
		// BIT 7, H (Check bit 7 in register H, if 0, set Z if 0
		// 0b01000000 >> 7
		regs.bit(1<<7, regs.H.val())
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
