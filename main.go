package main

import (
	"fmt"
	"github.com/hayeb/goboy/gameboy"
	"github.com/veandco/go-sdl2/sdl"
	"io/ioutil"
	"os"
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
		800, 600, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)

	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "Error: %s", r)
		}
	}()

	gameboy.Run(&cartridge, &bootrom, renderer)

}
