package main

import (
	"github.com/hayeb/goboy/gameboy"
	"github.com/veandco/go-sdl2/sdl"
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

	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		4 * 160, 4 * 144, sdl.WINDOW_SHOWN)

	mode := sdl.DisplayMode{}
	window.GetDisplayMode(&mode)
	mode.RefreshRate = 60
	window.SetDisplayMode(&mode)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED | sdl.RENDERER_PRESENTVSYNC)
	gameboy.Run(cartridge, bootrom, renderer)
}
