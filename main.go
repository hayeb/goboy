package main

import (
	"github.com/hayeb/goboy/gameboy"
	"github.com/banthar/Go-SDL/sdl"
	"io/ioutil"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	bootrom, error1 := ioutil.ReadFile("resources/DMG_ROM.bin")
	check(error1)

	cartridge, error2 := ioutil.ReadFile("resources/tetris.gb")
	check(error2)

	sdl.Init(sdl.INIT_EVERYTHING)

	window := sdl.SetVideoMode(4 * 160, 4 * 144, 32, 0)
	sdl.JoystickEventState(sdl.DISABLE)

	gameboy.Run(cartridge, bootrom, window)
}
