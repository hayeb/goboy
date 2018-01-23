package gameboy

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/veandco/go-sdl2/sdl"
)

var _ = spew.Config

func Run(cart []uint8, bootrom []uint8, renderer *sdl.Renderer) {
	ci, mem, reg, instrMap, cbInstrMap, graphics := initializeSystem(cart, bootrom, renderer)

	if ci.ramSize != ram_none || ci.romSize != rom_kbit_256 {
		panic("Cartridge not supported")
	}
	interruptMaster := false
	interruptEnableScheduled := false
	interruptDisableScheduled := false
	for true {
		oldPC := reg.PC.val()

		if oldPC > 0x100 {
			fmt.Println("oeps")
		}

		instrLength, name := executeInstruction(mem, reg, instrMap, cbInstrMap)

		if interruptEnableScheduled {
			interruptEnableScheduled = false
			interruptMaster = true
		} else if interruptDisableScheduled {
			interruptDisableScheduled = false
			interruptMaster = false
		}

		if name == "DI" {
			interruptDisableScheduled = true
		} else if name == "EI" {
			interruptEnableScheduled = true
		}

		graphics.updateGraphics(instrLength)
		handleInterupts(mem, reg, interruptMaster)
		// TODO: Update timers

		// Swap out the boot rom
		if oldPC == 0xfe {
			mem.swapBootRom(cart)
		}
	}
}

func handleInterupts(mem *memory, reg *register, master bool) {
	if !master {
		return
	}
	req := mem.read8(0xff0f)
	enabled := mem.read8(0xffff)
	if req > 0 {
		for i := 0; i < 5; i += 1 {
			if testBit(req, uint(i)) && testBit(enabled, uint(i)) {
				serviceInterupt(mem, reg, i, req)
			}
		}
	}
}

func serviceInterupt(mem *memory, reg *register, i int, requested uint8) {
	mem.write8(0xff0f, resetBit(requested, uint(i)))
	pushStack16(mem, reg, reg.PC.val())

	switch i {
	case 0:
		fmt.Println("Servicing V-BLANK interrupt")
		reg.PC = halfWordRegister(0x40)
	case 1:
		fmt.Println("Servicing LCD interrupt")
		reg.PC = halfWordRegister(0x48)
	case 2:
		fmt.Println("Servicing TIMER interrupt")
		reg.PC = halfWordRegister(0x50)
	case 4:
		fmt.Println("Servicing JOYPAD interrupt")
		reg.PC = halfWordRegister(0x60)
	default:
		panic(fmt.Sprintf("Servicing unknown interupt %d", i))
	}
}

// Executes the next instruction at the PC. Returns the length (in cycles) of the instruction
func executeInstruction(mem *memory, reg *register, instrMap *map[uint8]*instruction, cbInstrMap *map[uint8]*cbInstruction) (int, string) {
	instructionCode := mem.read8(reg.PC.val())
	instr, ok := (*instrMap)[instructionCode]

	if !ok {
		//spew.Dump(mem.videoRam)
		panic(fmt.Sprintf("Unrecognized instruction %#02x at address %#04x", instructionCode, reg.PC.val()))
	}

	if instr.name != "CB" {
		//fmt.Printf("%#04x\t%s\n", reg.PC.val(), instr.name)

		cycles := instr.executor(mem, reg, instr)
		return cycles, instr.name
	} else {
		cbCode := mem.read8(reg.PC.val() + 1)
		cb, ok := (*cbInstrMap)[cbCode]
		if !ok {
			panic(fmt.Sprintf("Unrecognized cb instruction %x at address %#04x", cbCode, reg.PC.val()+1))
		}
		//fmt.Printf("%#04x\t%s %s\n", reg.PC.val(), instr.name, cb.name)
		cycles := cb.executor(mem, reg, cb)
		reg.PC = halfWordRegister(reg.PC.val() + 1)
		return cycles + 4, cb.name
	}
}

func pushStack8(mem *memory, regs *register, val uint8) {
	mem.write8(regs.SP.val(), val)
	regs.decSP(1)
}

func pushStack16(mem *memory, reg *register, val uint16) {
	pushStack8(mem, reg, leastSig16(val))
	pushStack8(mem, reg, mostSig16(val))

}

func popStack8(mem *memory, reg *register) uint8 {
	reg.incSP(1)
	return mem.read8(reg.SP.val())
}

func popStack16(mem *memory, reg *register) uint16 {
	most := popStack8(mem, reg)
	least := popStack8(mem, reg)
	val := uint16(most)<<8 | uint16(least)
	return val
}

// Read a byte from memory from address SP + offset and returns the value
func readArgByte(mem *memory, reg *register, offset int) uint8 {
	return mem.read8(reg.PC.val() + uint16(offset))
}

// Read a halfword from memory from address SP + offset and returns the value
func readArgHalfword(mem *memory, reg *register, offset int) uint16 {
	return mem.read16(reg.PC.val() + uint16(offset))
}

func initializeSystem(cart []uint8, bootrom []uint8, ren *sdl.Renderer) (*cartridgeInfo, *memory, *register, *map[uint8]*instruction, *map[uint8]*cbInstruction, *graphics) {
	cartridgeInfo := createCartridgeInfo(cart)
	instructionMap := createInstructionMap()
	cbInstrucionMap := createCBInstructionMap()
	mem := memInit(bootrom, cart)
	graphics := createGraphics(mem, ren, cartridgeInfo)
	registers := new(register)
	return cartridgeInfo, mem, registers, instructionMap, cbInstrucionMap, graphics
}
