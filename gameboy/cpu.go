package gameboy

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/veandco/go-sdl2/sdl"
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

	fmt.Println("Open window")
	sdl.Init(sdl.INIT_EVERYTHING)

	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		800, 600, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	surface, err := window.GetSurface()
	rect := sdl.Rect{0, 0, 200, 200}
	surface.FillRect(&rect,0xffff0000)
	window.UpdateSurface()

	for {
		instructionCode := mem.read8(reg.PC.val())
		instr, ok := (*instrMap)[instructionCode]

		fmt.Printf("Mem at 0x0104: %#02x", mem.read8(0x0104))

		if !ok {
			spew.Dump(mem)
			spew.Dump(reg)
			panic(fmt.Sprintf("Unrecognized instruction %#02x at address %#04x", instructionCode, reg.PC.val()))
		}

		if instr.name != "CB" {
			fmt.Printf("%#04x\t%s\n", reg.PC.val(), instr.name)
			instr.executor(mem, reg)
			reg.PC = halfWordRegister(reg.PC.val() + uint16(instr.bytes))
		} else {
			cbCode := mem.read8(reg.PC.val() + 1)
			cb, ok := (*cbInstrMap)[cbCode]
			if !ok {
				spew.Dump(mem)
				spew.Dump(reg)
				panic(fmt.Sprintf("Unrecognized cb instruction %x at address %#04x", cbCode, reg.PC.val() + 1))
			}
			fmt.Printf("%#04x\t%s %s\n", reg.PC.val(), instr.name, cb.name)
			cb.executor(mem, reg)
			reg.PC = halfWordRegister(reg.PC.val() + uint16(cb.bytes))
		}
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
	val := uint16(most) << 8 | uint16(least)
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

func initializeSystem(cart []uint8, bootrom []uint8) (*cartridgeInfo, *memory, *register, *map[uint8]instruction, *map[uint8]cbInstruction) {
	cartridgeInfo := createCartridgeInfo(cart)
	instructionMap := createInstructionMap()
	cbInstrucionMap := createCBInstructionMap()
	mem := memInit(bootrom, cart)
	registers := new(register)
	return cartridgeInfo, mem, registers, instructionMap, cbInstrucionMap
}
