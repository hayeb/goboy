package gameboy

import "fmt"

func cbInstruction(cbCode uint8) {
	switch cbCode {
	case 0x7c:
		fmt.Println("Execute 0x7c instruction")
		panic("CB instruction 0x7c not yet implemented")
	default:
		panic(fmt.Sprintf("Unknown cb instruction %x", cbCode))
	}
}
