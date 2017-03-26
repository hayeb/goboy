package gameboy

import (
	"fmt"
	"os"
)

type memory struct {
	bank_0                    [16 * 1024]uint8    // 0x0000 (16 kB)
	switchable_rom_bank       [16 * 1024]uint8    // 0x4000 (16 kB)
	video_ram                 [8 * 1024]uint8     // 0x8000 (8 kB)
	switchable_ram_bank       [4 * 8 * 1024]uint8 // 0xA000 (8 kB)
	internal_ram_8kb          [8 * 1024]uint8     // 0xC000 (8 kB)
	echo_internal_ram         [8 * 1024]uint8     // 0xE000 (8 kB)
	sprite_attrib_memory      [7680]uint8         // 0xFE00 (7680 B)
	empty1                    [96]uint8           // 0xFEA0 (96 B)
	io_ports                  [76]uint8           // 0xFF00 (67 B)
	empty2                    [52]uint8           // 0xFF4C (52 B)
	internal_ram              [127]uint8          // 0xFF80 (127 B)
	interrupt_enable_register uint8               // 0xFFFF (1 B)
	memory_settings           memorySettings
}

type memorySettings struct {
	mbc1 bool
	mbc2 bool

	romBanking bool
	ramEnabled bool

	currentROMBank int
	currentRAMBank int
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

// Specific (Special) memory addresses
const (
	lcd_control_address      = 0xff40
	lcd_status_address       = 0xff41
	current_scanline_address = 0xff44
	target_scanline_address  = 0xff45
)

func memInit(bootrom *[]uint8, cartridge *[]uint8, rom_type_code int) *memory {
	b0 := [16 * 1024]uint8{}
	sw := [16 * 1024]uint8{}
	for index, item := range *bootrom {
		b0[index] = item
	}
	for i := 0xFF; i < len(b0); i++ {
		b0[i] = (*cartridge)[i]
	}
	for j := 0; j < len(sw); j++ {
		sw[j] = (*cartridge)[j+len(sw)]
	}

	return &memory{
		bank_0:                    b0,
		switchable_rom_bank:       sw,
		video_ram:                 [8 * 1024]uint8{},
		switchable_ram_bank:       [4 * 8 * 1024]uint8{},
		internal_ram_8kb:          [8 * 1024]uint8{},
		echo_internal_ram:         [8 * 1024]uint8{},
		sprite_attrib_memory:      [7680]uint8{},
		empty1:                    [96]uint8{},
		io_ports:                  [76]uint8{},
		empty2:                    [52]uint8{},
		internal_ram:              [127]uint8{},
		interrupt_enable_register: 0,
		memory_settings: memorySettings{
			romBanking: true,
		},
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

func (memory *memory) read8(address uint16) uint8 {
	switch map_addr(address) {
	case bank_0:
		return memory.bank_0[address]
	case switchable_rom_bank:
		return memory.switchable_rom_bank[address-0x4000+uint16(memory.memory_settings.currentROMBank)*0x4000]
	case video_ram:
		if address > 0x9900 {
			fmt.Println(os.Stderr, "Possible interesting vedio read..S")
		}
		//fmt.Printf("Reading from video ram: %#04x\n", address)
		return memory.video_ram[address - 0x8000]
	case switchable_ram_bank:
		return memory.switchable_ram_bank[address-0xA000+uint16(memory.memory_settings.currentRAMBank*0x2000)]
	case io_ports:
		return memory.io_ports[address-0xff00]
	case internal_ram:
		return memory.internal_ram[address-0xff80]
	case interrupt_enable_register:
		return memory.interrupt_enable_register
	default:
		panic(fmt.Sprintf("Read byte requested outside implemented memory: %x", address))
	}
}

func (memory *memory) read16(address uint16) uint16 {
	switch map_addr(address) {
	case bank_0:
		return uint16(memory.bank_0[address]) | (uint16(memory.bank_0[address+1]) << 8)
	case switchable_rom_bank:
		taddress := address - 0x4000 + uint16(memory.memory_settings.currentROMBank)*0x4000
		return uint16(memory.bank_0[taddress]) | (uint16(memory.bank_0[taddress+1]) << 8)
	case switchable_ram_bank:
		taddr := address - 0xA000 + uint16(memory.memory_settings.currentRAMBank*0x2000)
		return uint16(memory.switchable_ram_bank[taddr]) | uint16(memory.switchable_ram_bank[taddr+1])
	case io_ports:
		return uint16(memory.io_ports[address-0xff00]) | (uint16(memory.io_ports[address-0xff00+1]) << 8)
	default:
		panic(fmt.Sprintf("Read halfword requested unimplemented memory: %x", address))
	}
}

func (mem *memory) write8(address uint16, val uint8) {
	if mem.handleSpecificAddress(address, val) {
		return
	}
	switch map_addr(address) {
	case bank_0:
	case switchable_rom_bank:
		mem.doBankingAction(address, val)
	case video_ram:
		mem.video_ram[address-0x8000] = val
	case switchable_ram_bank:
		mem.switchable_ram_bank[address-0xa000] = val
	case internal_ram_8kb:
		mem.internal_ram_8kb[address-0xc000] = val
	case echo_internal_ram_8kb:
		mem.internal_ram_8kb[address-0x2000-0xc000] = val
	case sprite_attrib_memory:
		panic(fmt.Sprint("Write to sprite attribute memory not implemented"))
		mem.sprite_attrib_memory[address-0xfe00] = val
	case empty1:
		mem.empty1[address-0xfea0] = val
	case io_ports:
		mem.io_ports[address-0xff00] = val
	case empty2:
		mem.empty2[address-0xff4c] = val
	case internal_ram:
		mem.internal_ram[address-0xFF80] = val
	default:
		panic(fmt.Sprintf("Write byte not yet implemented for address: %#04x on %d", address, map_addr(address)))
	}
}

func (mem *memory) handleSpecificAddress(address uint16, val uint8) bool {
	if address == 0xff44 {
		// Reset the scanline to 0
		fmt.Println("Resetting scanline register 0xff44 to 0")
		mem.io_ports[0xff44-0xFF00] = 0
		return true
	} else if address == 0xff46 {
		panic("DMAtransfer")
	}
	return false
}

func (mem *memory) write16(address uint16, val uint16) {
	switch map_addr(address) {
	default:
		panic(fmt.Sprintf("Write halfword not yet implemented for address %#04x", address))
	}
}

func (mem *memory) doBankingAction(address uint16, val uint8) {
	settings := mem.memory_settings
	if address < 0x2000 && settings.mbc1 || settings.mbc2 {
		mem.doRAMBankEnable(address, val)
	} else if address >= 0x2000 && address < 0x4000 && settings.mbc1 || settings.mbc2 {
		mem.doChangeLoROMBank(address, val)
	} else if address > 0x4000 && address < 0x6000 {
		if settings.mbc1 {
			if settings.romBanking {
				mem.doChangeHiRomBank(val)
			} else {
				mem.doRAMBankChange(val)
			}
		}
	} else if address >= 0x6000 && address < 0x8000 && settings.mbc1 {
		mem.doChangeROMRAMMode(val)
	}
}

func (mem *memory) doRAMBankEnable(address uint16, val uint8) {
	fmt.Println(os.Stderr, "Performing RAM bank enable switch")
	settings := &mem.memory_settings
	if settings.mbc2 && address&(1<<4) == 1 {
		return
	}

	var test uint8 = val & 0xf
	if test == 0xA {
		settings.ramEnabled = true
	} else if test == 0x0 {
		settings.ramEnabled = false
	}
}

func (mem *memory) doChangeLoROMBank(address uint16, val uint8) {
	fmt.Println(os.Stderr, "Performing LoROM bank switching")
	if mem.memory_settings.mbc2 {
		mem.memory_settings.currentROMBank = int(val & 0xF)
		return
	}
	var lower5 uint8 = val & 31
	mem.memory_settings.currentROMBank &= 224
	mem.memory_settings.currentROMBank |= int(lower5)
}

func (mem *memory) doChangeHiRomBank(val uint8) {
	fmt.Println(os.Stderr, "Performing HiROM bank switching")
	mem.memory_settings.currentROMBank &= 31

	val &= 224
	mem.memory_settings.currentROMBank |= int(val)
}

func (mem *memory) doRAMBankChange(val uint8) {
	fmt.Println(os.Stderr, "Performing RAM bank change")
	mem.memory_settings.currentRAMBank = int(val & 0x3)
}

func (mem *memory) doChangeROMRAMMode(val uint8) {
	fmt.Println(os.Stderr, "Performing ROM/RAM mode change")
	newData := val & 0x1
	mem.memory_settings.romBanking = newData == 0
	if mem.memory_settings.romBanking {
		mem.memory_settings.currentRAMBank = 0
	}
}

func (mem *memory) requestInterupt(interupt_type int) {
	//fmt.Println("Interrupt requested")

	mem.write8(0xff0f, setBit(mem.read8(0xff0f), uint(interupt_type)))
}

func (mem *memory) swapBootRom(cartridge *[]uint8) {
	fmt.Println("Swapping out bootrom")
	for i := 0; i < 0xff; i += 1 {
		mem.bank_0[i] = (*cartridge)[i]
	}
}
