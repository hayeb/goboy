package gameboy

import (
	"github.com/davecgh/go-spew/spew"
)
type flags struct {
	Z bool
	N bool
	H bool
	C bool
}

type register struct {
	A uint8
	B uint8
	C uint8
	D uint8
	E uint8
	F uint8
	H uint8
	L uint8

	SP uint16
	PC uint16

	Flag flags
}

func initRegisters() register {
	return register{0, 0, 0, 0, 0, 0, 0, 0, 0, 0,flags{false, false, false, false} }
}

func Run(bootrom []byte, cartridge []byte) {
	registers := initRegisters()
	memory := memInit()
	spew.Dump(registers)
	spew.Dump(memory)

}
