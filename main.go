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
	debug := flag.Bool("debug", false, "Whether to start the debugger")
	speed := flag.Int("speed", 1, "Speed factor, should be >= 1. Default is 1")

	flag.Parse()

	if rom == nil || *rom == "" {
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

	if *debug {
		gameboy.RunDebugger(&gb)
	} else {
		input := gameboy.Input{}
		for true {
			gb.Step()
			updateInput(&input)
			gb.HandleInput(&input)
		}
	}
}

func updateInput(input *gameboy.Input) {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.KeyboardEvent:
			switch t.Keysym.Sym {
			case sdl.K_a:
				if t.Type == sdl.KEYDOWN {
					input.A = true
				} else {
					input.A = false
				}
			case sdl.K_b:
				if t.Type == sdl.KEYDOWN {
					input.B = true
				} else {
					input.B = false
				}
			case sdl.K_LEFT:
				if t.Type == sdl.KEYDOWN {
					input.LEFT = true
				} else {
					input.LEFT = false
				}
			case sdl.K_RIGHT:
				if t.Type == sdl.KEYDOWN {
					input.RIGHT = true
				} else {
					input.RIGHT = false
				}
			case sdl.K_UP:
				if t.Type == sdl.KEYDOWN {
					input.UP = true
				} else {
					input.UP = false
				}
			case sdl.K_DOWN:
				if t.Type == sdl.KEYDOWN {
					input.DOWN = true
				} else {
					input.DOWN = false
				}
			case sdl.K_RETURN:
				if t.Type == sdl.KEYDOWN {
					input.ENTER = true
				} else {
					input.ENTER = false
				}
			case sdl.K_SPACE:
				if t.Type == sdl.KEYDOWN {
					input.SPACE = true
				} else {
					input.SPACE = false
				}
			}
		case *sdl.QuitEvent:
			os.Exit(0)
		}
	}
}
