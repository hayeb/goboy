package gameboy

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/veandco/go-sdl2/sdl"
	"time"
)

var _ = spew.Config

const (
	screen_update_cycles         = 69905
	fps                  float64 = 60.0
)

func Run(cart *[]uint8, bootrom *[]uint8, renderer *sdl.Renderer) {
	ci, mem, reg, instrMap, cbInstrMap := initializeSystem(cart, bootrom)
	if ci.ramSize != ram_none || ci.romSize != rom_kbit_256 {
		panic("Cartridge not supported")
	}
	running := true
	for running {
		i := 0
		start := time.Now()
		for i < screen_update_cycles {
			instrLength := executeInstruction(mem, reg, instrMap, cbInstrMap)

			i += instrLength
			// TODO: Update timers
			// TODO: Update Graphics
			// TODO: Handle interrupts
		}
		updateScreen(renderer)

		elapsed := time.Since(start)
		var ticks float64 = float64(elapsed.Nanoseconds()) * float64(1E-6)
		if ticks < 1000.0/fps {
			wait := (1000.0 / fps) - ticks
			sdl.Delay(uint32(wait))
		}

	}
}

func updateScreen(renderer *sdl.Renderer) {
	renderer.Present()
}

// Executes the next instruction at the PC. Returns the length (in cycles) of the instructione
func executeInstruction(mem *memory, reg *register, instrMap *map[uint8](*instruction), cbInstrMap *map[uint8](*cbInstruction)) int {
	instructionCode := mem.read8(reg.PC.val())
	instr, ok := (*instrMap)[instructionCode]

	if !ok {
		panic(fmt.Sprintf("Unrecognized instruction %#02x at address %#04x", instructionCode, reg.PC.val()))
	}

	if instr.name != "CB" {
		fmt.Printf("%#04x\t%s\n", reg.PC.val(), instr.name)
		cycles := instr.executor(mem, reg, instr)
		reg.PC = halfWordRegister(reg.PC.val() + uint16(instr.bytes))
		return cycles
	} else {
		cbCode := mem.read8(reg.PC.val() + 1)
		cb, ok := (*cbInstrMap)[cbCode]
		if !ok {
			panic(fmt.Sprintf("Unrecognized cb instruction %x at address %#04x", cbCode, reg.PC.val()+1))
		}
		fmt.Printf("%#04x\t%s %s\n", reg.PC.val(), instr.name, cb.name)
		cycles := cb.executor(mem, reg, cb)
		reg.PC = halfWordRegister(reg.PC.val() + uint16(cb.bytes))
		return cycles + 4
	}
}

func pushStack8(mem *memory, regs *register, val uint8) {
	mem.write8(regs.SP.val(), val)
	regs.decSP(1)
}

func pushStack16(mem *memory, reg *register, val uint16) {
	pushStack8(mem, reg, mostSig16(val))
	pushStack8(mem, reg, leastSig16(val))
}

func popStack8(mem *memory, reg *register) uint8 {
	reg.incSP(1)
	return mem.read8(reg.SP.val())
}

func popStack16(mem *memory, reg *register) uint16 {
	least := popStack8(mem, reg)
	most := popStack8(mem, reg)
	val := uint16(most)<<8 | uint16(least)
	return val
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

func initializeSystem(cart *[]uint8, bootrom *[]uint8) (*cartridgeInfo, *memory, *register, *map[uint8]*instruction, *map[uint8]*cbInstruction) {
	cartridgeInfo := createCartridgeInfo(cart)
	instructionMap := createInstructionMap()
	cbInstrucionMap := createCBInstructionMap()
	mem := memInit(bootrom, cart, int(cartridgeInfo.CartType))
	registers := new(register)
	return cartridgeInfo, mem, registers, instructionMap, cbInstrucionMap
}
