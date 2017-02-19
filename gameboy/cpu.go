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
	ci, mem, reg, instrMap, cbInstrMap := initializeSystem(cart, bootrom)
	if ci.ramSize != ram_none || ci.romSize != rom_kbit_256 {
		panic("Cartridge not supported")
	}

	for {
		instructionCode := mem.read8(reg.PC.val())
		instr, ok := (*instrMap)[instructionCode]

		if !ok {
			panic(fmt.Sprintf("Unrecognized instruction %x at address %#04x", instructionCode, reg.PC.val()))
		}

		fmt.Printf("%#04x \t %s\n", reg.PC.val(), instr.name)

		if instr.name != "CB" {
			instr.executor(mem, reg)
			reg.PC = halfWordRegister(reg.PC.val() + uint16(instr.bytes))
		} else {
			// Look up the CB instruction
			cbCode := mem.read8(reg.PC.val() + 1)
			cb, ok := (*cbInstrMap)[cbCode]
			if !ok {
				panic(fmt.Sprintf("Unrecognized cb instruction %x at address %#04x", cbCode, reg.PC.val()+1))
			}
			cb.executor(mem, reg)
			reg.PC = halfWordRegister(reg.PC.val() + uint16(cb.bytes))
		}

		spew.Dump(reg)
	}
}

func pushStack8(mem *memory, regs *register, val uint8) {
	mem.write8(regs.SP.val(), val)
	regs.decSP(1)
}

func pushStack16(mem *memory, reg *register, val uint16) {
	left := mostSig16(val)
	right := leastSig16(val)
	pushStack8(mem, reg, left)
	pushStack8(mem, reg, right)
}

func incRegister8(reg *byteRegister) {
	*reg = byteRegister(reg.val() + 1)
}

func decrRegister8(reg *byteRegister) {
	*reg = byteRegister(reg.val() - 1)
}

// Read a byte from memory from address SP + offset and returns the value
func readArgByte(mem *memory, reg *register, offset int) uint8 {
	return mem.read8(reg.PC.val() + uint16(offset))
}

// Read a halfword from memory from address SP + offset and returns the value
func readArgHalfword(mem *memory, reg *register, offset int) uint16 {
	return mem.read16(reg.PC.val() + uint16(offset))
}

func initializeSystem(cart []uint8, bootrom []uint8) (*cartridgeInfo, *memory, *register, *map[uint8]instruction, *map[uint8]cbInstruction) {
	cartridgeInfo := createCartridgeInfo(cart)
	instructionMap := createInstructionMap()
	cbInstrucionMap := createCBInstructionMap()
	mem := memInit(bootrom, cart)
	registers := new(register)
	return cartridgeInfo, mem, registers, instructionMap, cbInstrucionMap
}
