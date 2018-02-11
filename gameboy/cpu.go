package gameboy

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
)

var _ = spew.Config

func (gb *Gameboy) Step() {
	if gb.halted {
		gb.updateTimer(4)
		gb.graphics.updateGraphics(4)
		if gb.handleInterrupts() {
			gb.halted = false
		}

		return
	}

	oldPC := gb.reg.PC

	instrLength, name := gb.executeInstruction()

	if gb.interruptEnableScheduled {
		if gb.options.Debug {
			fmt.Printf("%#04x: Enable interrupts\n", oldPC)
		}
		gb.interruptEnableScheduled = false
		gb.interruptMaster = true
	} else if gb.interruptDisableScheduled {
		gb.interruptDisableScheduled = false
		gb.interruptMaster = false
		if gb.options.Debug {
			fmt.Printf("%#04x: Disable interrupts\n", oldPC)
		}
	}

	// TODO: Move to respective instructions
	if name == "DI" {
		gb.interruptDisableScheduled = true
		if gb.options.Debug {
			fmt.Println("Scheduling disable interrupts")
		}
	} else if name == "EI" {
		gb.interruptEnableScheduled = true
		if gb.options.Debug {
			fmt.Println("Scheduling enabling interrupts")
		}
	} else if name == "RETI" {
		gb.interruptMaster = true
		if gb.options.Debug {
			fmt.Println("Return enabling interrupts")
		}
	} else if name == "HALT" {
		gb.halted = true
	}

	gb.updateTimer(instrLength)
	gb.graphics.updateGraphics(instrLength)
	if gb.handleInterrupts() {
		gb.halted = false
	}

	// Swap out the boot rom
	if oldPC == 0xFE && !gb.bootromSwapped {
		gb.mem.swapBootRom(gb.cartridge)
		gb.bootromSwapped = true
	}
}

func (gb *Gameboy) PC() uint16 {
	return gb.reg.PC
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

func (gb *Gameboy) HandleInput(input *Input) bool {
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

	gb.mem.ioPorts[0x00] = joyPadReg

	return joyPadReg < 0xf
}

func (gb *Gameboy) handleInterrupts() bool {
	if !gb.interruptMaster {
		return false
	}
	req := gb.mem.read8(0xFF0F)
	enabled := gb.mem.read8(0xFFFF)
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

	// TODO: Implement serial transfer
	switch i {
	case 0:
		if gb.options.Debug {
			fmt.Println("Servicing VBLANK interrupt")
		}
		gb.reg.PC = 0x40
	case 1:
		if gb.options.Debug {
			fmt.Println("Servicing LCD STAT interrupt")
		}
		gb.reg.PC = 0x48
	case 2:
		if gb.options.Debug {
			fmt.Println("Servicing timer overflow interrupt")
		}
		gb.reg.PC = 0x50
	case 4:
		if gb.options.Debug {
			fmt.Println("Servicing JOYPAD interrupt")
		}
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
		panic(fmt.Sprintf("Unrecognized instruction %#02x at address %#04x", instructionCode, gb.reg.PC))
	}

	if instr.name != "CB" {
		if gb.options.Debug && gb.bootromSwapped {
			if instr.bytes == 1 {
				fmt.Printf("%#04x %-12s\n", gb.reg.PC, instr.name)
			} else if instr.bytes == 2 {
				fmt.Printf("%#04x %-12s%-#02x\n", gb.reg.PC, instr.name, uint16(readArgByte(gb.mem, gb.reg)))
			} else if instr.bytes == 3 {
				fmt.Printf("%#04x %-12s%-#04x\n", gb.reg.PC, instr.name, uint16(readArgHalfword(gb.mem, gb.reg)))
			}
		}

		cycles := instr.executor(gb.mem, gb.reg, instr)

		return cycles, instr.name
	} else {
		cbCode := gb.mem.read8(gb.reg.PC + 1)
		cb, ok := (*gb.cbInstruction)[cbCode]
		if !ok {
			panic(fmt.Sprintf("Unrecognized cb instruction %x at address %#04x", cbCode, gb.reg.PC))
		}
		if gb.options.Debug && gb.bootromSwapped {
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
func readArgByte(mem *memory, reg *register) uint8 {
	return mem.read8(reg.PC + uint16(1))
}

// Read a halfword from memory from address SP + offset and returns the value
func readArgHalfword(mem *memory, reg *register) uint16 {
	return mem.read16(reg.PC + uint16(1))
}

func (gb *Gameboy) Regs() *register {
	return gb.reg
}
