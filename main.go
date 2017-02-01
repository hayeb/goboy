package main

import (
	"io/ioutil"
	"github.com/hayeb/goboy/gameboy"
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

	gameboy.Run(bootrom, cartridge);
}
