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
	io_ports                  [76]uint8        // 0xFF00 (67 B)
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
	echo_internal_ram_8kb     = 5
	sprite_attrib_memory      = 6
	empty1                    = 7
	io_ports                  = 8
	empty2                    = 9
	internal_ram              = 10
	interrupt_enable_register = 11
)

func memInit(bootrom []uint8, cartridge []uint8) *memory {
	b0 := [16 * 1024]uint8{}
	sw := [16 * 1024]uint8{}
	for index, item := range bootrom {
		b0[index] = item
	}
	// TODO: Implement switching of the bootrom page when memory address 0xFE is executed
	for i := 0xFF; i < len(b0); i++ {
		b0[i] = cartridge[i-0xFF]
	}
	for j := 0; j < len(sw); j++ {
		sw[j] = cartridge[j+len(sw)]
	}
	return &memory{
		b0,
		sw,
		[8 * 1024]uint8{},
		[8 * 1024]uint8{},
		[8 * 1024]uint8{},
		[8 * 1024]uint8{},
		[7680]uint8{},
		[96]uint8{},
		[76]uint8{},
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
		return echo_internal_ram_8kb
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

func (memory memory) read8(address uint16) uint8 {
	switch map_addr(address) {
	case bank_0:
		return memory.bank_0[address]
	case switchable_rom_bank:
		return memory.switchable_rom_bank[address-0x4000]
	case io_ports:
		return memory.io_ports[address-0xff00]
	default:
		panic(fmt.Sprintf("Read byte requested outside implemented memory: %x", address))
	}
}

func (memory *memory) read16(address uint16) uint16 {
	switch map_addr(address) {
	case bank_0:
		return uint16(memory.bank_0[address]) | (uint16(memory.bank_0[address+1]) << 8)
	case switchable_rom_bank:
		return uint16(memory.bank_0[address-0x4000]) | (uint16(memory.bank_0[address-0x4000+1]) << 8)
	case io_ports:
		return uint16(memory.io_ports[address-0xff00]) | (uint16(memory.io_ports[address-0xff00+1]) << 8)
	default:
		panic(fmt.Sprintf("Read halfword requested unimplemented memory: %x", address))
	}
}

func (mem *memory) write8(address uint16, val uint8) {
	switch map_addr(address) {
	case bank_0:
		mem.bank_0[address] = val
	case switchable_rom_bank:
		mem.switchable_rom_bank[address-0x4000] = val
	case video_ram:
		mem.video_ram[address-0x8000] = val
	case io_ports:
		mem.io_ports[address-0xff00] = val
	case echo_internal_ram_8kb:
		mem.internal_ram_8kb[address-0x2000-0xc000] = val
	case internal_ram:
		mem.internal_ram[address-0xFF80] = val
	default:
		panic(fmt.Sprintf("Write byte not yet implemented for address: %x", address))
	}
}

func (mem *memory) write16(address uint16, val uint16) {
	switch map_addr(address) {
	default:
		panic(fmt.Sprintf("Write halfword not yet implemented for address %#04x", address))
	}
}
