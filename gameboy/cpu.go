package gameboy

import (
	"fmt"
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
	defer func () {
		if r := recover(); r != nil {
			fmt.Printf("Encountered an error: %s", r)
		}
	}()
	cartridgeInfo := GetCartridgeInfo(cartridge)
	fmt.Println(CartridgeInfoString(cartridgeInfo))

	if cartridgeInfo.RAMSize != RAM_NONE || cartridgeInfo.ROMSize != ROM_KBIT_256 {
		panic("Cartridge not supported")
	}

	registers := initRegisters()



	// TODO: First run bootrom, then load cartridge!

	// Load the rom in memory locations 0x000 - 0x7FFF

	memory := memInit(cartridgeInfo.CartType, cartridge)

	// Run the boot rom

}
