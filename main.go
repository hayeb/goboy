package main

import (
	"github.com/hayeb/goboy/gameboy"
	"github.com/banthar/Go-SDL/sdl"
	"io/ioutil"
	"flag"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {

	rom := flag.String("rom", "", "Rom to be loaded")
	scale := flag.Int("scale", 4, "Scaling factor to be used. Default is 4, resulting in 4*160 x 4*144 resolution")
	debug := flag.Bool("debug", false, "Whether to emit debug information during execution")
	speed := flag.Int("speed", 1, "Speed factor")

	flag.Parse()

	cartridge, error2 := ioutil.ReadFile(*rom)
	check(error2)

	sdl.Init(sdl.INIT_EVERYTHING)

	window := sdl.SetVideoMode(4*160, 4*144, 32, sdl.HWACCEL)
	sdl.JoystickEventState(sdl.DISABLE)

	gb := gameboy.Initialize(cartridge, window, &gameboy.Options{Scaling: *scale, Debug: *debug, Speed: *speed})
	gb.Run(cartridge, window)
}
