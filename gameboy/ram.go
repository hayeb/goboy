package gameboy

import "fmt"

type memory struct {
	bank_0                    [16 * 1024]uint8 // 0x0000 (16 kB)
	switchable_rom_bank       [16 * 1024]uint8 // 0x4000 (16 kB)
	video_ram                 [8 * 1024]uint8  // 0x8000 (8 kB)
	switchable_ram_bank       [8 * 1024]uint8  // 0xA000 (8 kB)
	internal_ram_8kb          [8 * 1024]uint8  // 0xC000 (8 kB)
	echo_internal_ram         [8 * 1024]uint8  // 0xE000 (8 kB)
	sprite_attrib_memory      [7680]uint8      // 0xFE00 (7680 B)
	empty1                    [96]uint8        // 0xFEA0 (96 B)
	io_ports                  [67]uint8        // 0xFF00 (67 B)
	empty2                    [52]uint8        // 0xFF4C (52 B)
	internal_ram              [127]uint8       // 0xFF80 (127 B)
	interrupt_enable_register uint8            // 0xFFFF (1 B)
}

const (
	bank_0                    = 0
	switchable_rom_bank       = 1
	video_ram                 = 2
	switchable_ram_bank       = 3
	internal_ram_8kb          = 4
	echo_internal_ram         = 5
	sprite_attrib_memory      = 6
	empty1                    = 7
	io_ports                  = 8
	empty2                    = 9
	internal_ram              = 10
	interrupt_enable_register = 11
)

func memInit(bootrom []uint8) *memory {
	fb := [16 * 1024]uint8{}
	for index, item := range bootrom {
		fb[index] = item
	}
	return &memory{
		fb,
		[16 * 1024]uint8{},
		[8 * 1024]uint8{},
		[8 * 1024]uint8{},
		[8 * 1024]uint8{},
		[8 * 1024]uint8{},
		[7680]uint8{},
		[96]uint8{},
		[67]uint8{},
		[52]uint8{},
		[127]uint8{},
		0,
	}

}

func map_addr(addr uint16) int {
	if addr < 0x4000 {
		return bank_0
	} else if addr >= 0x4000 && addr < 0x8000 {
		return switchable_rom_bank
	} else if addr >= 8000 && addr < 0xA000 {
		return video_ram
	} else if addr >= 0xA000 && addr < 0xC000 {
		return switchable_ram_bank
	} else if addr >= 0xC000 && addr < 0xE000 {
		return internal_ram_8kb
	} else if addr >= 0xE000 && addr < 0xFE00 {
		return echo_internal_ram
	} else if addr >= 0xFE00 && addr < 0xFEA0 {
		return sprite_attrib_memory
	} else if addr >= 0xFEA0 && addr < 0xFF00 {
		return empty1
	} else if addr >= 0xFF00 && addr < 0xFF4C {
		return io_ports
	} else if addr >= 0xFF4C && addr < 0xFF80 {
		return empty2
	} else if addr >= 0xFF80 && addr < 0xFFFF {
		return internal_ram
	} else if addr == 0xFFFF {
		return interrupt_enable_register
	} else {
		panic(fmt.Sprintf("Unknown memory address %x\n", addr))
	}
}

func (mem memory) LoadROM(cartridge []uint8) {

}

// TODO: Fix
func (memory memory) read8(addres uint16) uint8 {
	switch map_addr(addres) {
	case 0:
		return memory.bank_0[addres]
	case 1:
		return memory.switchable_rom_bank[addres-0x4000]
	case 2:
		return memory.bank_0[int(addres)]
	case 3:
		return memory.bank_0[int(addres)]
	case 4:
		return memory.bank_0[int(addres)]
	case 5:
		return memory.bank_0[int(addres)]
	case 6:
		return memory.bank_0[int(addres)]
	case 7:
		return memory.bank_0[int(addres)]
	case 8:
		return memory.bank_0[int(addres)]
	case 9:
		return memory.bank_0[int(addres)]
	case 10:
		return memory.bank_0[int(addres)]
	default:
		panic(fmt.Sprintf("Read requested outside memory: %x", addres))
	}
}

func (memory memory) read16(addres uint16) uint16 {
	switch map_addr(addres) {
	case 0:
		return uint16(memory.bank_0[addres]) | (uint16(memory.bank_0[addres+1]) << 8)
	case 1:
		return uint16(memory.switchable_rom_bank[addres-0x4000])
	case 2:
		panic(fmt.Sprintf("Read requested unimplemented memory: %x", addres))
	case 3:
		panic(fmt.Sprintf("Read requested unimplemented memory: %x", addres))
	case 4:
		panic(fmt.Sprintf("Read requested unimplemented memory: %x", addres))
	case 5:
		panic(fmt.Sprintf("Read requested unimplemented memory: %x", addres))
	case 6:
		panic(fmt.Sprintf("Read requested unimplemented memory: %x", addres))
	case 7:
		panic(fmt.Sprintf("Read requested unimplemented memory: %x", addres))
	case 8:
		panic(fmt.Sprintf("Read requested unimplemented memory: %x", addres))
	case 9:
		panic(fmt.Sprintf("Read requested unimplemented memory: %x", addres))
	case 10:
		panic(fmt.Sprintf("Read requested unimplemented memory: %x", addres))
	default:
		panic(fmt.Sprintf("Read requested unimplemented memory: %x", addres))
	}
}

func (memory memory) write8(Saddress uint16, val uint8) {

}

func (memory memory) write16(address uint16, val uint16) {

}
