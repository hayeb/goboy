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
	Debug   bool
	Speed   int
}

type Gameboy struct {
	cartridgeInfo  *cartridgeInfo
	instructionMap *map[uint8]*instruction
	cbInstruction  *map[uint8]*cbInstruction
	mem            *memory
	graphics       *graphics
	reg            *register
	options        *Options
	cartridge      []uint8

	timer *timer
	input *input

	interruptMaster           bool
	interruptEnableScheduled  bool
	interruptDisableScheduled bool
	halted                    bool
	stopped                   bool
	bootromSwapped            bool
}

func Initialize(cart []uint8, renderer *sdl.Surface, options *Options) *Gameboy {
	instructionMap := createInstructionMap()
	cbInstrucionMap := createCBInstructionMap()
	mem := memInit(cart)
	graphics := createGraphics(mem.videoRam[:], mem.ioPorts[:], mem.spriteAttribMemory[:], renderer, options.Speed, options.Scaling)
	registers := new(register)
	return &Gameboy{
		cartridgeInfo:   createCartridgeInfo(cart),
		instructionMap:  instructionMap,
		cbInstruction:   cbInstrucionMap,
		mem:             mem,
		graphics:        graphics,
		reg:             registers,
		options:         options,
		cartridge:       cart,
		timer:           new(timer),
		interruptMaster: true,
		input:           new(input),
	}
}

func (gb *Gameboy) Run() {
	for true {
		if gb.halted || gb.stopped {
			gb.updateTimer(4)

			if gb.halted {
				gb.graphics.updateGraphics(4)
				if gb.handleInterrupts() {
					gb.halted = false
				}

			}

			if gb.stopped {
				if gb.handleInput() {
					gb.stopped = false
				}
			}

			continue
		}

		oldPC := gb.reg.PC
		instrLength, name := gb.executeInstruction()

		if gb.interruptEnableScheduled {
			gb.interruptEnableScheduled = false
			gb.interruptMaster = true
		} else if gb.interruptDisableScheduled {
			gb.interruptDisableScheduled = false
			gb.interruptMaster = false
		}

		if name == "DI" {
			gb.interruptDisableScheduled = true
		} else if name == "EI" {
			gb.interruptEnableScheduled = true
		} else if name == "RETI" {
			gb.interruptMaster = true
		} else if name == "HALT" {
			gb.halted = true
		} else if name == "stop" {
			gb.stopped = true
		}

		gb.updateTimer(instrLength)
		gb.graphics.updateGraphics(instrLength)
		if gb.handleInterrupts() {
			gb.halted = false
		}
		gb.handleInput()

		// Swap out the boot rom
		if oldPC == 0xfe && !gb.bootromSwapped {
			gb.mem.swapBootRom(gb.cartridge)
			gb.bootromSwapped = true
		}
	}
}

func (gb *Gameboy) updateTimer(cycles int) {
	gb.timer.t += cycles

	if gb.timer.t >= 16 {
		gb.timer.m++
		gb.timer.t -= 16
		gb.timer.d++
		if gb.timer.d == 16 {
			gb.timer.d = 0
			val := gb.mem.ioPorts[0x04]
			gb.mem.ioPorts[0x04] = val + 1
		}
	}

	tac := gb.mem.ioPorts[0x07]

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

		if gb.timer.m >= t {
			gb.timer.m = 0
			gb.mem.ioPorts[0x05]++

			if gb.mem.ioPorts[0x05] == 0 {
				gb.mem.ioPorts[0x05] = gb.mem.ioPorts[0x06]
				gb.mem.ioPorts[0x0F] |= 0x4
			}
		}
	}
}

func (gb *Gameboy) handleInput() bool {
	event := sdl.PollEvent()
	switch t := event.(type) {
	case *sdl.KeyboardEvent:
		switch t.Keysym.Sym {
		case sdl.K_a:
			if t.Type == sdl.KEYDOWN {
				gb.input.A = true
			} else {
				gb.input.A = false
			}
		case sdl.K_b:
			if t.Type == sdl.KEYDOWN {
				gb.input.B = true
			} else {
				gb.input.B = false
			}
		case sdl.K_LEFT:
			if t.Type == sdl.KEYDOWN {
				gb.input.LEFT = true
			} else {
				gb.input.LEFT = false
			}
		case sdl.K_RIGHT:
			if t.Type == sdl.KEYDOWN {
				gb.input.RIGHT = true
			} else {
				gb.input.RIGHT = false
			}
		case sdl.K_UP:
			if t.Type == sdl.KEYDOWN {
				gb.input.UP = true
			} else {
				gb.input.UP = false
			}
		case sdl.K_DOWN:
			if t.Type == sdl.KEYDOWN {
				gb.input.DOWN = true
			} else {
				gb.input.DOWN = false
			}
		case sdl.K_RETURN:
			if t.Type == sdl.KEYDOWN {
				gb.input.ENTER = true
			} else {
				gb.input.ENTER = false
			}
		case sdl.K_SPACE:
			if t.Type == sdl.KEYDOWN {
				gb.input.SPACE = true
			} else {
				gb.input.SPACE = false
			}
		}
	}

	return gb.updateJoyReg()
}

func (gb *Gameboy) updateJoyReg() bool {
	joyPadReg := gb.mem.ioPorts[0x00]
	if !testBit(joyPadReg, 4) && !testBit(joyPadReg, 5) {
		// Do nothing when input is not polled
		return false
	} else if testBit(joyPadReg, 4) && testBit(joyPadReg, 5) {
		// Reset the register
		gb.mem.write8(0xFF00, 0xf)
		return false
	}

	if testBit(joyPadReg, 4) {
		if gb.input.ENTER {
			joyPadReg = resetBit(joyPadReg, 3)
		} else {
			joyPadReg = setBit(joyPadReg, 3)
		}
		if gb.input.SPACE {
			joyPadReg = resetBit(joyPadReg, 2)
		} else {
			joyPadReg = setBit(joyPadReg, 2)
		}
		if gb.input.B {
			joyPadReg = resetBit(joyPadReg, 1)
		} else {
			joyPadReg = setBit(joyPadReg, 1)
		}
		if gb.input.A {
			joyPadReg = resetBit(joyPadReg, 0)
		} else {
			joyPadReg = setBit(joyPadReg, 0)
		}
	} else if testBit(joyPadReg, 5) {
		if gb.input.DOWN {
			joyPadReg = resetBit(joyPadReg, 3)
		} else {
			joyPadReg = setBit(joyPadReg, 3)
		}
		if gb.input.UP {
			joyPadReg = resetBit(joyPadReg, 2)
		} else {
			joyPadReg = setBit(joyPadReg, 2)
		}
		if gb.input.LEFT {
			joyPadReg = resetBit(joyPadReg, 1)
		} else {
			joyPadReg = setBit(joyPadReg, 1)
		}
		if gb.input.RIGHT {
			joyPadReg = resetBit(joyPadReg, 0)
		} else {
			joyPadReg = setBit(joyPadReg, 0)
		}
	}

	joyPadReg = resetBit(joyPadReg, 4)
	joyPadReg = resetBit(joyPadReg, 5)

	gb.mem.ioPorts[0x00] = joyPadReg

	return joyPadReg < 0xf
}

func (gb *Gameboy) handleInterrupts() bool {
	if !gb.interruptMaster {
		return false
	}
	req := gb.mem.read8(0xff0f)
	enabled := gb.mem.read8(0xffff)
	handled := false
	if req > 0 {
		for i := 0; i < 5; i += 1 {
			if testBit(req, uint(i)) && testBit(enabled, uint(i)) {
				gb.serviceInterrupt(i, req)
				handled = true
			}
		}
	}
	return handled
}

func (gb *Gameboy) serviceInterrupt(i int, requested uint8) {
	gb.mem.write8(0xff0f, resetBit(requested, uint(i)))
	pushStack16(gb.mem, gb.reg, gb.reg.PC)

	switch i {
	case 0:
		gb.reg.PC = 0x40
	case 1:
		fmt.Println("Servicing LCD interrupt")
		gb.reg.PC = 0x48
	case 2:
		gb.reg.PC = 0x50
	case 4:
		fmt.Println("Servicing JOYPAD interrupt")
		gb.reg.PC = 0x60
	default:
		panic(fmt.Sprintf("Servicing unknown interupt %d", i))
	}
}

// Executes the next instruction at the PC. Returns the length (in cycles) of the instruction
func (gb *Gameboy) executeInstruction() (int, string) {
	instructionCode := gb.mem.read8(gb.reg.PC)
	instr, ok := (*gb.instructionMap)[instructionCode]

	if !ok {
		//spew.Dump(mem.videoRam)
		panic(fmt.Sprintf("Unrecognized instruction %#02x at address %#04x", instructionCode, gb.reg.PC))
	}

	if instr.name != "CB" {
		if gb.options.Debug {
			fmt.Printf("%#04x\t%s\n", gb.reg.PC, instr.name)
		}

		cycles := instr.executor(gb.mem, gb.reg, instr)

		return cycles, instr.name
	} else {
		cbCode := gb.mem.read8(gb.reg.PC + 1)
		cb, ok := (*gb.cbInstruction)[cbCode]
		if !ok {
			panic(fmt.Sprintf("Unrecognized cb instruction %x at address %#04x", cbCode, gb.reg.PC))
		}
		if gb.options.Debug {
			fmt.Printf("%#04x\t%s %s\n", gb.reg.PC, instr.name, cb.name)
		}

		cycles := cb.executor(gb.mem, gb.reg, cb)
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
