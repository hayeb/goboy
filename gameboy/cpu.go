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
		case call_nn:
			// Push adress of next instruction onto the stack
			// TODO: See if correct?
			left := mostSig16(regs.PC.val()+uint16(3))
			right := leastSig16((regs.PC.val()+uint16(3)))
			pushStack(mem, regs, left)
			pushStack(mem, regs, right)

			regs.PC = halfWordRegister(readArgHalfword(mem, regs, 1))
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
		case ld_a_DE:
			regs.A = byteRegister(mem.read8(regs.readDuo(reg_de)))
		case ld_c:
			regs.C = byteRegister(readArgByte(mem, regs, 1))
		case ld_C_a:
			mem.write8(0xFF00+uint16(regs.C.val()), regs.A.val())
		case ld_de_d16:
			regs.writeDuo(reg_de, readArgHalfword(mem, regs, 1))
		case ld_sp:
			arg := readArgHalfword(mem, regs, 1)
			fmt.Printf("SP: %#04x\n", arg)
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

		regs.PC = halfWordRegister(regs.PC.val() + uint16(instr.bytes))
		spew.Dump(mem.internal_ram)
	}
}

func pushStack(mem *memory, regs *register, val uint8) {
	mem.write8(regs.SP.val(), val)
	regs.decSP(1)
}

func incRegister8(reg *byteRegister) {
	*reg = byteRegister(reg.val() + 1)
}

func decrRegister8(reg *byteRegister) {
	*reg = byteRegister(reg.val() - 1)
}

// reads a byte from memory from address SP + offset and returns the value
func readArgByte(mem *memory, reg *register, offset int) uint8 {
	return mem.read8(reg.PC.val() + uint16(offset))
}

// Reads a halfword from memory from address SP + offset and returns the value
func readArgHalfword(mem *memory, reg *register, offset int) uint16 {
	return mem.read16(reg.PC.val() + uint16(offset))
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
	mem := memInit(bootrom, cart)
	registers := new(register)
	return cartridgeInfo, mem, registers, instructionMap
}
// Returns a uint8 with the 8 least signigicant bits of i
func leastSig16(i uint16) uint8 {
	return uint8(i & ((1 << 8) - 1))
}

// Returns a uint8 with the 8 most signigicant bits of i
func mostSig16(i uint16) uint8 {
	return uint8(i>>8)
}

