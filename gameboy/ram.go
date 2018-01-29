package gameboy

import (
	"fmt"
	"os"
)

type memory struct {
	bank0                   [16 * 1024]uint8    // 0x0000 (16 kB)
	switchableRomBank       [16 * 1024]uint8    // 0x4000 (16 kB)
	videoRam                [8 * 1024]uint8     // 0x8000 (8 kB)
	switchableRamBank       [4 * 8 * 1024]uint8 // 0xA000 (8 kB)
	internalRam8kb          [8 * 1024]uint8     // 0xC000 (8 kB)
	echoInternalRam         [8 * 1024]uint8     // 0xE000 (8 kB)
	spriteAttribMemory      [7680]uint8         // 0xFE00 (7680 B)
	empty1                  [96]uint8           // 0xFEA0 (96 B)
	ioPorts                 [76]uint8           // 0xFF00 (67 B)
	empty2                  [52]uint8           // 0xFF4C (52 B)
	internalRam             [127]uint8          // 0xFF80 (127 B)
	interruptEnableRegister uint8               // 0xFFFF (1 B)
	memorySettings          memorySettings
}

type memorySettings struct {
	mbc1 bool
	mbc2 bool

	romBanking bool
	ramEnabled bool

	currentROMBank int
	currentRAMBank int

	dma_pending_time int
}

const (
	bank0                   = 0
	switchableRomBank       = 1
	videoRam                = 2
	switchableRamBank       = 3
	internalRam8kb          = 4
	echoInternalRam8kb      = 5
	spriteAttribMemory      = 6
	empty1                  = 7
	ioPorts                 = 8
	empty2                  = 9
	internalRam             = 10
	interruptEnableRegister = 11
)

var bootrom = []uint8{
  0x31, 0xfe, 0xff, 0xaf, 0x21, 0xff, 0x9f, 0x32, 0xcb, 0x7c, 0x20, 0xfb,
  0x21, 0x26, 0xff, 0x0e, 0x11, 0x3e, 0x80, 0x32, 0xe2, 0x0c, 0x3e, 0xf3,
  0xe2, 0x32, 0x3e, 0x77, 0x77, 0x3e, 0xfc, 0xe0, 0x47, 0x11, 0x04, 0x01,
  0x21, 0x10, 0x80, 0x1a, 0xcd, 0x95, 0x00, 0xcd, 0x96, 0x00, 0x13, 0x7b,
  0xfe, 0x34, 0x20, 0xf3, 0x11, 0xd8, 0x00, 0x06, 0x08, 0x1a, 0x13, 0x22,
  0x23, 0x05, 0x20, 0xf9, 0x3e, 0x19, 0xea, 0x10, 0x99, 0x21, 0x2f, 0x99,
  0x0e, 0x0c, 0x3d, 0x28, 0x08, 0x32, 0x0d, 0x20, 0xf9, 0x2e, 0x0f, 0x18,
  0xf3, 0x67, 0x3e, 0x64, 0x57, 0xe0, 0x42, 0x3e, 0x91, 0xe0, 0x40, 0x04,
  0x1e, 0x02, 0x0e, 0x0c, 0xf0, 0x44, 0xfe, 0x90, 0x20, 0xfa, 0x0d, 0x20,
  0xf7, 0x1d, 0x20, 0xf2, 0x0e, 0x13, 0x24, 0x7c, 0x1e, 0x83, 0xfe, 0x62,
  0x28, 0x06, 0x1e, 0xc1, 0xfe, 0x64, 0x20, 0x06, 0x7b, 0xe2, 0x0c, 0x3e,
  0x87, 0xe2, 0xf0, 0x42, 0x90, 0xe0, 0x42, 0x15, 0x20, 0xd2, 0x05, 0x20,
  0x4f, 0x16, 0x20, 0x18, 0xcb, 0x4f, 0x06, 0x04, 0xc5, 0xcb, 0x11, 0x17,
  0xc1, 0xcb, 0x11, 0x17, 0x05, 0x20, 0xf5, 0x22, 0x23, 0x22, 0x23, 0xc9,
  0xce, 0xed, 0x66, 0x66, 0xcc, 0x0d, 0x00, 0x0b, 0x03, 0x73, 0x00, 0x83,
  0x00, 0x0c, 0x00, 0x0d, 0x00, 0x08, 0x11, 0x1f, 0x88, 0x89, 0x00, 0x0e,
  0xdc, 0xcc, 0x6e, 0xe6, 0xdd, 0xdd, 0xd9, 0x99, 0xbb, 0xbb, 0x67, 0x63,
  0x6e, 0x0e, 0xec, 0xcc, 0xdd, 0xdc, 0x99, 0x9f, 0xbb, 0xb9, 0x33, 0x3e,
  0x3c, 0x42, 0xb9, 0xa5, 0xb9, 0xa5, 0x42, 0x3c, 0x21, 0x04, 0x01, 0x11,
  0xa8, 0x00, 0x1a, 0x13, 0xbe, 0x20, 0xfe, 0x23, 0x7d, 0xfe, 0x34, 0x20,
  0xf5, 0x06, 0x19, 0x78, 0x86, 0x23, 0x05, 0x20, 0xfb, 0x86, 0x20, 0xfe,
  0x3e, 0x01, 0xe0, 0x50}

func memInit(cartridge []uint8) *memory {
	b0 := [16 * 1024]uint8{}
	sw := [16 * 1024]uint8{}
	for index, item := range bootrom {
		b0[index] = item
	}
	for i := 0x100; i < len(b0); i++ {
		b0[i] = cartridge[i]
	}
	for j := 0; j < len(sw); j++ {
		sw[j] = cartridge[j+len(sw)]
	}

	mem := &memory{
		bank0:                   b0,
		switchableRomBank:       sw,
		videoRam:                [8 * 1024]uint8{},
		switchableRamBank:       [4 * 8 * 1024]uint8{},
		internalRam8kb:          [8 * 1024]uint8{},
		echoInternalRam:         [8 * 1024]uint8{},
		spriteAttribMemory:      [7680]uint8{},
		empty1:                  [96]uint8{},
		ioPorts:                 [76]uint8{},
		empty2:                  [52]uint8{},
		internalRam:             [127]uint8{},
		interruptEnableRegister: 0,
		memorySettings: memorySettings{
			romBanking:       true,
			dma_pending_time: 0,
		},
	}
	return mem
}

func mapAddr(addr uint16) int {
	if addr < 0x4000 {
		return bank0
	} else if addr >= 0x4000 && addr < 0x8000 {
		return switchableRomBank
	} else if addr >= 8000 && addr < 0xA000 {
		return videoRam
	} else if addr >= 0xA000 && addr < 0xC000 {
		return switchableRamBank
	} else if addr >= 0xC000 && addr < 0xE000 {
		return internalRam8kb
	} else if addr >= 0xE000 && addr < 0xFE00 {
		return echoInternalRam8kb
	} else if addr >= 0xFE00 && addr < 0xFEA0 {
		return spriteAttribMemory
	} else if addr >= 0xFEA0 && addr < 0xFF00 {
		return empty1
	} else if addr >= 0xFF00 && addr < 0xFF4C {
		return ioPorts
	} else if addr >= 0xFF4C && addr < 0xFF80 {
		return empty2
	} else if addr >= 0xFF80 && addr < 0xFFFF {
		return internalRam
	} else if addr == 0xFFFF {
		return interruptEnableRegister
	} else {
		panic(fmt.Sprintf("Unknown memory address %x\n", addr))
	}
}

func (memory *memory) read8(address uint16) uint8 {
	switch mapAddr(address) {
	case bank0:
		return memory.bank0[address]
	case switchableRomBank:
		return memory.switchableRomBank[address-0x4000+uint16(memory.memorySettings.currentROMBank)*0x4000]
	case videoRam:
		return memory.videoRam[address-0x8000]
	case switchableRamBank:
		return memory.switchableRamBank[address-0xA000+uint16(memory.memorySettings.currentRAMBank*0x2000)]
	case internalRam8kb:
		return memory.internalRam8kb[address-0xC000]
	case echoInternalRam8kb:
		return memory.echoInternalRam[address-0xE000]
	case ioPorts:
		return memory.ioPorts[address-0xFF00]
	case empty2:
		return memory.empty2[address - 0xFF4C]
	case internalRam:
		return memory.internalRam[address-0xFF80]
	case interruptEnableRegister:
		return memory.interruptEnableRegister
	default:
		panic(fmt.Sprintf("Read byte requested outside implemented memory: %x", address))
	}
}

func (memory *memory) read16(address uint16) uint16 {
	switch mapAddr(address) {
	case bank0:
		return uint16(memory.bank0[address]) | (uint16(memory.bank0[address+1]) << 8)
	case switchableRomBank:
		taddress := address - 0x4000 + uint16(memory.memorySettings.currentROMBank)*0x4000
		return uint16(memory.switchableRomBank[taddress]) | (uint16(memory.switchableRomBank[taddress+1]) << 8)
	case switchableRamBank:
		taddr := address - 0xA000 + uint16(memory.memorySettings.currentRAMBank*0x2000)
		return uint16(memory.switchableRamBank[taddr]) | uint16(memory.switchableRamBank[taddr+1])<<8
	case internalRam8kb:
		return uint16(memory.internalRam8kb[address-0xC000]) | uint16(memory.internalRam8kb[address-0xC000+1]) << 8
	case ioPorts:
		return uint16(memory.ioPorts[address-0xff00]) | (uint16(memory.ioPorts[address-0xff00+1]) << 8)
	default:
		panic(fmt.Sprintf("Read halfword requested unimplemented memory: %x", address))
	}
}

func (memory *memory) write8(address uint16, val uint8) {
	if memory.handleSpecificAddress(address, val) {
		return
	}
	switch mapAddr(address) {
	case bank0:
	case switchableRomBank:
		memory.doBankingAction(address, val)
	case videoRam:
		memory.videoRam[address-0x8000] = val
	case switchableRamBank:
		memory.switchableRamBank[address-0xa000] = val
	case internalRam8kb:
		memory.internalRam8kb[address-0xc000] = val
	case echoInternalRam8kb:
		memory.internalRam8kb[address-0x2000-0xc000] = val
		// TODO: Not accessible while display is updating
	case spriteAttribMemory:
		memory.spriteAttribMemory[address-0xfe00] = val
	case empty1:
		memory.empty1[address-0xfea0] = val
	case ioPorts:
		memory.ioPorts[address-0xff00] = val
	case empty2:
		memory.empty2[address-0xff4c] = val
	case internalRam:
		memory.internalRam[address-0xFF80] = val
	case interruptEnableRegister:
		memory.interruptEnableRegister = val
	default:
		panic(fmt.Sprintf("Write byte not yet implemented for address: %#04x on %d", address, mapAddr(address)))
	}
}

func (memory *memory) handleSpecificAddress(address uint16, val uint8) bool {
	switch address {
	case 0xFF40:
		memory.ioPorts[0x40] = val
		return true
	case 0xFF44:
		// Reset the scanline to 0
		memory.ioPorts[0x44] = 0
		return true
	case 0xFF46:
		for i := 0; i < 0xA0; i++ {
			memory.spriteAttribMemory[i] = memory.read8(uint16(val)*0x100 + uint16(i))
		}
		return true
	case 0xFF04:
		fmt.Println("Reset divider to 0")
		memory.ioPorts[0x04] = 0
		return true
	case 0xFF05:
		fmt.Printf("Writing to counter: %#04x\n", val)
		memory.ioPorts[0x05] = val
		return true
	case 0xFF06:
		fmt.Printf("Writing to modulo: %#04x\n", val)
		memory.ioPorts[0x06] = val
		return true
	case 0xFF07:
		fmt.Printf("Writing to timer control: %b\n", val)
		memory.ioPorts[0x07] = val
		return true
	}
	return false
}

func (memory *memory) write16(address uint16, val uint16) {
	switch mapAddr(address) {
	case bank0:
		memory.bank0[address] = uint8(val)
		memory.bank0[address+1] = uint8(val >> 8)
	case internalRam8kb:
		memory.internalRam8kb[address - 0xC000] = uint8(val)
		memory.internalRam8kb[address - 0xC000 + 1] = uint8(val >> 8)
	default:
		panic(fmt.Sprintf("Write halfword not yet implemented for address %#04x", address))
	}
}

func (memory *memory) doBankingAction(address uint16, val uint8) {
	settings := memory.memorySettings
	if address < 0x2000 && settings.mbc1 || settings.mbc2 {
		memory.doRAMBankEnable(address, val)
	} else if address >= 0x2000 && address < 0x4000 && settings.mbc1 || settings.mbc2 {
		memory.doChangeLoROMBank(address, val)
	} else if address > 0x4000 && address < 0x6000 {
		if settings.mbc1 {
			if settings.romBanking {
				memory.doChangeHiRomBank(val)
			} else {
				memory.doRAMBankChange(val)
			}
		}
	} else if address >= 0x6000 && address < 0x8000 && settings.mbc1 {
		memory.doChangeROMRAMMode(val)
	}
}

func (memory *memory) doRAMBankEnable(address uint16, val uint8) {
	fmt.Println(os.Stderr, "Performing RAM bank enable switch")
	settings := &memory.memorySettings
	if settings.mbc2 && address&(1<<4) == 1 {
		return
	}

	var test = val & 0xf
	if test == 0xA {
		settings.ramEnabled = true
	} else if test == 0x0 {
		settings.ramEnabled = false
	}
}

func (memory *memory) doChangeLoROMBank(address uint16, val uint8) {
	fmt.Println(os.Stderr, "Performing LoROM bank switching")
	if memory.memorySettings.mbc2 {
		memory.memorySettings.currentROMBank = int(val & 0xF)
		return
	}
	var lower5 = val & 31
	memory.memorySettings.currentROMBank &= 224
	memory.memorySettings.currentROMBank |= int(lower5)
}

func (memory *memory) doChangeHiRomBank(val uint8) {
	fmt.Println(os.Stderr, "Performing HiROM bank switching")
	memory.memorySettings.currentROMBank &= 31

	val &= 224
	memory.memorySettings.currentROMBank |= int(val)
}

func (memory *memory) doRAMBankChange(val uint8) {
	fmt.Println(os.Stderr, "Performing RAM bank change")
	memory.memorySettings.currentRAMBank = int(val & 0x3)
}

func (memory *memory) doChangeROMRAMMode(val uint8) {
	fmt.Println(os.Stderr, "Performing ROM/RAM mode change")
	newData := val & 0x1
	memory.memorySettings.romBanking = newData == 0
	if memory.memorySettings.romBanking {
		memory.memorySettings.currentRAMBank = 0
	}
}

func (memory *memory) requestInterupt(interruptType int) {
	fmt.Printf("Interrupt %d requested", interruptType)

	memory.write8(0xff0f, setBit(memory.read8(0xff0f), uint(interruptType)))
}

func (memory *memory) swapBootRom(cartridge []uint8) {
	for i := 0; i < 0x100; i += 1 {
		memory.bank0[i] = cartridge[i]
	}
}
