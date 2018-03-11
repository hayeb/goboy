package gameboy

import (
	"bufio"
	"os"

	"github.com/alecthomas/participle"
	"fmt"
)

type Command struct {
	StepCommand       bool        `@("s"|"")`
	QuitCommand       bool        `| @"q"`
	BreakpointCommand *Breakpoint `| "b" @@`
	PrintRegsCommand  bool        `| @"p"`
	StatCommand       *Stat       `| "t" @@`
	RunCommand        bool        `| @"r"`
	HelpCommand       bool        `| @"h"`
}

type Stat struct {
	InterruptStatCommand bool `@"i"`
	MemoryStatCommand    bool `| @"m"`
	StackStatCommand     bool `| @"s"`
}

type Breakpoint struct {
	SetBreakpoint   *int `@Int`
	ListBreakpoints bool `| @"l"`
}

type Debugger struct {
	breakpoints []uint16
	gb          *Gameboy
	input       *Input
}

func RunDebugger(gb *Gameboy, updateInputFunction func(input *Input)) {
	fmt.Println("Welcome to the GoBoy debugger.")
	reader := bufio.NewReader(os.Stdin)
	stopped := false

	parser, err := participle.Build(&Command{}, nil)
	if err != nil {
		panic(err)
	}

	debugger := Debugger{
		gb:          gb,
		breakpoints: make([]uint16, 0),
		input:       new(Input),
	}

	for !stopped {
		fmt.Print(">: ")
		command := &Command{}
		commandLine, e := reader.ReadString('\n')

		if e != nil {
			fmt.Println(e)
			return
		}

		pErr := parser.ParseString(commandLine, command)
		if pErr != nil {
			fmt.Println(pErr)
			continue
		}
		stopped = debugger.handleCommand(*command, updateInputFunction)
	}
}

func (debugger *Debugger) handleCommand(command Command, updateInputFunction func(input *Input)) bool {
	if command.StepCommand {
		debugger.gb.Step()
		updateInputFunction(debugger.input)
		debugger.gb.HandleInput(debugger.input)
	} else if command.QuitCommand {
		return true
	} else if command.BreakpointCommand != nil {
		breakPointCommand := command.BreakpointCommand
		if breakPointCommand.SetBreakpoint != nil {
			bp := uint16(*breakPointCommand.SetBreakpoint)
			debugger.breakpoints = append(debugger.breakpoints, bp)
			fmt.Printf("Set breakpoint at %#04x\n", bp)

		} else if breakPointCommand.ListBreakpoints {
			for _, bp := range debugger.breakpoints {
				fmt.Printf("%#04x\n", bp)
			}
		}
	} else if command.RunCommand {
		hit, bp := false, uint16(0)
		for !hit {
			debugger.gb.Step()
			updateInputFunction(debugger.input)
			debugger.gb.HandleInput(debugger.input)
			hit, bp = debugger.breakpointHit(debugger.gb.PC())
		}
		fmt.Printf("Hit breakpoint at %#04x\n", bp)
	} else if command.PrintRegsCommand {
		regs := debugger.gb.Regs()
		fmt.Printf("A: %#04x\t", regs.A)
		fmt.Printf("B: %#04x\n", regs.B)
		fmt.Printf("C: %#04x\t", regs.C)
		fmt.Printf("D: %#04x\n", regs.D)
		fmt.Printf("E: %#04x\t", regs.E)
		fmt.Printf("F: %#04x\n", regs.F)
		fmt.Printf("H: %#04x\t", regs.H)
		fmt.Printf("L: %#04x\n", regs.L)

		fmt.Printf("AF: %#04x\t", regs.readDuo(REG_AF))
		fmt.Printf("BC: %#04x\n", regs.readDuo(REG_BC))
		fmt.Printf("DE: %#04x\t", regs.readDuo(REG_DE))
		fmt.Printf("HL: %#04x\n", regs.readDuo(REG_HL))

		fmt.Printf("PC: %#04x\tSP: %#04x\n", regs.PC, regs.SP)
	} else if command.StatCommand != nil {
		if command.StatCommand.InterruptStatCommand {
			ifReg := debugger.gb.mem.ioPorts[0x0F]
			ieReg := debugger.gb.mem.interruptEnableRegister
			fmt.Printf("Interrupt master: %t\n", debugger.gb.interruptMaster)
			fmt.Println("[Name] [Enabled] [Requested]")
			fmt.Printf("VBLANK: %t %t\n", testBit(ieReg, 0), testBit(ifReg, 0))
			fmt.Printf("LCDC: %t %t\n", testBit(ieReg, 1), testBit(ifReg, 1))
			fmt.Printf("Timer overflow: %t %t\n", testBit(ieReg, 2), testBit(ifReg, 2))
			fmt.Printf("Serial transfer complete: %t %t\n", testBit(ieReg, 3), testBit(ifReg, 3))
		} else if command.StatCommand.MemoryStatCommand {
			mem := debugger.gb.mem

			fmt.Printf("Memor settings:\n")
			fmt.Printf("MBC1: %t\n", mem.memorySettings.mbc1)
			fmt.Printf("MBC2: %t\n", mem.memorySettings.mbc2)
			fmt.Printf("MBC3: %t\n", mem.memorySettings.mbc3)
			if mem.memorySettings.bankingMode == romBankingMode {
				fmt.Printf("Banking mode: ROM\n")
			} else if mem.memorySettings.bankingMode == ramBankingMode {
				fmt.Printf("Banking mode: RAM\n")
			}
			fmt.Printf("Current ROM bank: %d\n", mem.memorySettings.currentROMBank)
			fmt.Printf("Current RAM bank: %d\n", mem.memorySettings.currentRAMBank)
		} else if command.StatCommand.StackStatCommand {
			mem := debugger.gb.mem
			reg := debugger.gb.reg
			fmt.Println("Stack dump:")
			for i := 0; i < 10; i++ {
				fmt.Printf("%#04x %#02x\n",reg.SP + uint16(i), mem.read8(reg.SP + uint16(i)))
			}

		}
	} else if command.HelpCommand {
		// TODO: Implement help
	} else {
		fmt.Println("Command not yet implemented!")
		return true
	}
	return false
}

func (debugger *Debugger) breakpointHit(pc uint16) (bool, uint16) {
	for _, bp := range debugger.breakpoints {
		if pc == bp {
			return true, bp
		}
	}
	return false, 0
}
