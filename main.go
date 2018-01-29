package main

import (
	"github.com/hayeb/goboy/gameboy"
	"github.com/banthar/Go-SDL/sdl"
	"io/ioutil"
	"flag"
	"fmt"
	"os"
)

func check(e error) {
	if e != nil {
		fmt.Println(e.Error())
		os.Exit(1)
	}
}

func main() {

	rom := flag.String("rom", "", "Rom to be loaded")
	scale := flag.Int("scale", 4, "Scaling factor to be used. Default is 4, resulting in 4*160 x 4*144 resolution")
	debug := flag.Bool("debug", false, "Whether to emit debug information during execution, default is true")
	speed := flag.Int("speed", 1, "Speed factor, should be >= 1. Default is 1")

	flag.Parse()

	if rom == nil || *rom == ""{
		fmt.Println("Please specify a rom using -rom")
		os.Exit(1)
	}

	if *scale <= 0 {
		fmt.Println("Invalid scale")
		os.Exit(1)
	}

	if *speed < 1 {
		fmt.Println("Invalid speed")
		os.Exit(1)
	}


	cartridge, error2 := ioutil.ReadFile(*rom)
	check(error2)

	sdl.Init(sdl.INIT_EVERYTHING)

	window := sdl.SetVideoMode(*scale*160, *scale*144, 32, sdl.HWACCEL)
	defer window.Free()
	sdl.JoystickEventState(sdl.DISABLE)

	gb := gameboy.Initialize(cartridge, window, &gameboy.Options{Scaling: *scale, Debug: *debug, Speed: *speed})
	gb.Run(cartridge, window)
}
