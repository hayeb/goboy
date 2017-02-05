package main

import (
	"github.com/hayeb/goboy/gameboy"
	"io/ioutil"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func initializeSystem(cartridge []uint8, bootrom []uint8) (_ *gameboy.CartridgeInfo, _ *gameboy.Memory, _ *gameboy.Register, _ *map[uint8]gameboy.Instruction) {
	cartridgeInfo := gameboy.CreateCartridgeInfo(cartridge)
	instructionMap := gameboy.CreateInstructionMap()
	memory := gameboy.MemInit(bootrom)
	registers := new(gameboy.Register)
	return cartridgeInfo, memory, registers, instructionMap
}

func main() {
	bootrom, error1 := ioutil.ReadFile("resources/DMG_ROM.bin")
	check(error1)

	cartridge, error2 := ioutil.ReadFile("resources/tetris.gb")
	check(error2)

	gameboy.Run(initializeSystem(cartridge, bootrom))
}
