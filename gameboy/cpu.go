package gameboy

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/banthar/Go-SDL/sdl"
)

var _ = spew.Config

type input struct {
	A     bool
	B     bool
	LEFT  bool
	RIGHT bool
	UP    bool
	DOWN  bool
	ENTER bool
	SPACE bool
}

type timer struct {
	t int
	m int
	d int
}

type Options struct {
	Scaling int
	Debug bool
	Speed int
}

type Gameboy struct {
	cartridgeInfo  *cartridgeInfo
	instructionMap *map[uint8]*instruction
	cbInstruction  *map[uint8]*cbInstruction
	mem            *memory
	graphics       *graphics
	reg            *register
}

func Initialize(cart []uint8, renderer *sdl.Surface, options *Options) *Gameboy{
	cartridgeInfo := createCartridgeInfo(cart)
	instructionMap := createInstructionMap()
	cbInstrucionMap := createCBInstructionMap()
	mem := memInit(cart)
	graphics := createGraphics(mem, renderer, cartridgeInfo)
	registers := new(register)
	return &Gameboy{
		cartridgeInfo: cartridgeInfo,
		instructionMap:instructionMap,
		cbInstruction:cbInstrucionMap,
		mem:mem,
		graphics:graphics,
		reg:registers,
	}
}

func (gb *Gameboy) Run(cart []uint8, renderer *sdl.Surface) {
	input := input{}

	//if ci.ramSize != ram_none || ci.romSize != rom_kbit_256 {
	//	panic("Cartridge not supported")
	//}
	interruptMaster := true
	interruptEnableScheduled := false
	interruptDisableScheduled := false
	timer := new(timer)
	halted := false
	stopped := false
	bootromSwapped := false
	for true {
		if halted || stopped {
			updateTimer(gb.mem, timer, 4)

			if halted {
				gb.graphics.updateGraphics(4)
				if handleInterrupts(gb.mem, gb.reg, interruptMaster) {
					halted = false
				}

			}

			if stopped {
				if handleInput(&input, gb.mem) {
					stopped = false
				}
			}

			continue
		}

		oldPC := gb.reg.PC
		instrLength, name := executeInstruction(gb.mem, gb.reg, gb.instructionMap, gb.cbInstruction)

		if oldPC == 0x20b {
			fmt.Println("BREAK!")
		}

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
		} else if name == "RETI" {
			interruptMaster = true
		} else if name == "HALT" {
			halted = true
		} else if name == "stop" {
			stopped = true
		}

		updateTimer(gb.mem, timer, instrLength)
		gb.graphics.updateGraphics(instrLength)
		if handleInterrupts(gb.mem, gb.reg, interruptMaster) {
			halted = false
		}
		handleInput(&input, gb.mem)

		// Swap out the boot rom
		if oldPC == 0xfe && !bootromSwapped {
			gb.mem.swapBootRom(cart)
			bootromSwapped = true
		}
	}
}

func updateTimer(mem *memory, timer *timer, cycles int) {
	timer.t += cycles

	if timer.t >= 16 {
		timer.m++
		timer.t -= 16
		timer.d++
		if timer.d == 16 {
			timer.d = 0
			val := mem.ioPorts[0x04]
			mem.ioPorts[0x04] = val + 1
		}
	}

	tac := mem.ioPorts[0x07]

	// If the timer is turned on,
	t := 0
	if testBit(tac, 2) {
		val := tac & 0x3
		switch val {
		case 0:
			t = 64
		case 1:
			t = 1
		case 2:
			t = 4
		case 3:
			t = 16
		}

		if timer.m >= t {
			timer.m = 0
			mem.ioPorts[0x05]++

			if mem.ioPorts[0x05] == 0 {
				mem.ioPorts[0x05] = mem.ioPorts[0x06]
				mem.ioPorts[0x0F] |= 0x4
			}
		}
	}
}

func handleInput(input *input, mem *memory) bool {
	event := sdl.PollEvent()
	switch t := event.(type) {
	case *sdl.KeyboardEvent:
		switch t.Keysym.Sym {
		case sdl.K_a:
			if t.Type == sdl.KEYDOWN {
				input.A = true
			} else {
				input.A = false
			}
		case sdl.K_b:
			if t.Type == sdl.KEYDOWN {
				input.B = true
			} else {
				input.B = false
			}
		case sdl.K_LEFT:
			if t.Type == sdl.KEYDOWN {
				input.LEFT = true
			} else {
				input.LEFT = false
			}
		case sdl.K_RIGHT:
			if t.Type == sdl.KEYDOWN {
				input.RIGHT = true
			} else {
				input.RIGHT = false
			}
		case sdl.K_UP:
			if t.Type == sdl.KEYDOWN {
				input.UP = true
			} else {
				input.UP = false
			}
		case sdl.K_DOWN:
			if t.Type == sdl.KEYDOWN {
				input.DOWN = true
			} else {
				input.DOWN = false
			}
		case sdl.K_RETURN:
			if t.Type == sdl.KEYDOWN {
				input.ENTER = true
			} else {
				input.ENTER = false
			}
		case sdl.K_SPACE:
			if t.Type == sdl.KEYDOWN {
				input.SPACE = true
			} else {
				input.SPACE = false
			}
		}
	}

	return updateJoyReg(input, mem)
}

func updateJoyReg(input *input, mem *memory) bool {
	joyPadReg := mem.ioPorts[0x00]
	if !testBit(joyPadReg, 4) && !testBit(joyPadReg, 5) {
		// Do nothing when input is not polled
		return false
	} else if testBit(joyPadReg, 4) && testBit(joyPadReg, 5) {
		// Reset the register
		mem.write8(0xFF00, 0xf)
		return false
	}

	if testBit(joyPadReg, 4) {
		if input.ENTER {
			joyPadReg = resetBit(joyPadReg, 3)
		} else {
			joyPadReg = setBit(joyPadReg, 3)
		}
		if input.SPACE {
			joyPadReg = resetBit(joyPadReg, 2)
		} else {
			joyPadReg = setBit(joyPadReg, 2)
		}
		if input.B {
			joyPadReg = resetBit(joyPadReg, 1)
		} else {
			joyPadReg = setBit(joyPadReg, 1)
		}
		if input.A {
			joyPadReg = resetBit(joyPadReg, 0)
		} else {
			joyPadReg = setBit(joyPadReg, 0)
		}
	} else if testBit(joyPadReg, 5) {
		if input.DOWN {
			joyPadReg = resetBit(joyPadReg, 3)
		} else {
			joyPadReg = setBit(joyPadReg, 3)
		}
		if input.UP {
			joyPadReg = resetBit(joyPadReg, 2)
		} else {
			joyPadReg = setBit(joyPadReg, 2)
		}
		if input.LEFT {
			joyPadReg = resetBit(joyPadReg, 1)
		} else {
			joyPadReg = setBit(joyPadReg, 1)
		}
		if input.RIGHT {
			joyPadReg = resetBit(joyPadReg, 0)
		} else {
			joyPadReg = setBit(joyPadReg, 0)
		}
	}

	joyPadReg = resetBit(joyPadReg, 4)
	joyPadReg = resetBit(joyPadReg, 5)

	mem.ioPorts[0x00] = joyPadReg

	return joyPadReg < 0xf
}

func handleInterrupts(mem *memory, reg *register, master bool) bool {
	if !master {
		return false
	}
	req := mem.read8(0xff0f)
	enabled := mem.read8(0xffff)
	handled := false
	if req > 0 {
		for i := 0; i < 5; i += 1 {
			if testBit(req, uint(i)) && testBit(enabled, uint(i)) {
				serviceInterrupt(mem, reg, i, req)
				handled = true
			}
		}
	}
	return handled
}

func serviceInterrupt(mem *memory, reg *register, i int, requested uint8) {
	mem.write8(0xff0f, resetBit(requested, uint(i)))
	pushStack16(mem, reg, reg.PC)

	switch i {
	case 0:
		reg.PC = 0x40
	case 1:
		fmt.Println("Servicing LCD interrupt")
		reg.PC = 0x48
	case 2:
		reg.PC = 0x50
	case 4:
		fmt.Println("Servicing JOYPAD interrupt")
		reg.PC = 0x60
	default:
		panic(fmt.Sprintf("Servicing unknown interupt %d", i))
	}
}

// Executes the next instruction at the PC. Returns the length (in cycles) of the instruction
func executeInstruction(mem *memory, reg *register, instrMap *map[uint8]*instruction, cbInstrMap *map[uint8]*cbInstruction) (int, string) {
	instructionCode := mem.read8(reg.PC)
	instr, ok := (*instrMap)[instructionCode]

	if !ok {
		//spew.Dump(mem.videoRam)
		panic(fmt.Sprintf("Unrecognized instruction %#02x at address %#04x", instructionCode, reg.PC))
	}

	if instr.name != "CB" {
		fmt.Printf("%#04x\t%s\n", reg.PC, instr.name)
		cycles := instr.executor(mem, reg, instr)

		return cycles, instr.name
	} else {
		cbCode := mem.read8(reg.PC + 1)
		cb, ok := (*cbInstrMap)[cbCode]
		if !ok {
			panic(fmt.Sprintf("Unrecognized cb instruction %x at address %#04x", cbCode, reg.PC))
		}
		fmt.Printf("%#04x\t%s %s\n", reg.PC, instr.name, cb.name)
		cycles := cb.executor(mem, reg, cb)
		return cycles + 4, cb.name
	}
}

func pushStack8(mem *memory, regs *register, val uint8) {
	mem.write8(regs.SP, val)
	regs.decSP(1)
}

func pushStack16(mem *memory, reg *register, val uint16) {
	pushStack8(mem, reg, leastSig16(val))
	pushStack8(mem, reg, mostSig16(val))

}

func popStack8(mem *memory, reg *register) uint8 {
	reg.incSP(1)
	return mem.read8(reg.SP)
}

func popStack16(mem *memory, reg *register) uint16 {
	most := popStack8(mem, reg)
	least := popStack8(mem, reg)
	val := uint16(most)<<8 | uint16(least)
	return val
}

// Read a byte from memory from address SP + offset and returns the value
func readArgByte(mem *memory, reg *register, offset int) uint8 {
	return mem.read8(reg.PC + uint16(offset))
}

// Read a halfword from memory from address SP + offset and returns the value
func readArgHalfword(mem *memory, reg *register, offset int) uint16 {
	return mem.read16(reg.PC + uint16(offset))
}