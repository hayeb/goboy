package gameboy

import (
	"fmt"
	"github.com/banthar/Go-SDL/sdl"
)

type Options struct {
	Scaling int
	Debug   bool
	Speed   int
}

type Input struct {
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

	interruptMaster           bool
	interruptEnableScheduled  bool
	interruptDisableScheduled bool
	halted                    bool
	bootromSwapped            bool
}

func Initialize(cart []uint8, renderer *sdl.Surface, options *Options) Gameboy {
	cartInfo := createCartridgeInfo(cart)
	instructionMap := createInstructionMap()
	cbInstrucionMap := createCBInstructionMap()
	mem := memInit(cart, cartInfo)
	graphics := createGraphics(mem.videoRam[:], mem.ioPorts[:], mem.spriteAttribMemory[:], renderer, options.Speed, options.Scaling)
	registers := new(register)


	gameboy := Gameboy{
		cartridgeInfo:   cartInfo,
		instructionMap:  instructionMap,
		cbInstruction:   cbInstrucionMap,
		mem:             mem,
		graphics:        graphics,
		reg:             registers,
		options:         options,
		cartridge:       cart,
		timer:           new(timer),
		interruptMaster: true,
	}

	fmt.Printf("GoBoy initialized:\n%s", cartridgeInfoString(*cartInfo))
	return gameboy
}
