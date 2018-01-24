package gameboy

import (
	"fmt"
	"os"
)

type memory struct {
	bank0              [16 * 1024]uint8    // 0x0000 (16 kB)
	switchableRomBank  [16 * 1024]uint8    // 0x4000 (16 kB)
	videoRam           [8 * 1024]uint8     // 0x8000 (8 kB)
	switchableRamBank  [4 * 8 * 1024]uint8 // 0xA000 (8 kB)
	internalRam8kb     [8 * 1024]uint8     // 0xC000 (8 kB)
	echoInternalRam    [8 * 1024]uint8     // 0xE000 (8 kB)
	spriteAttribMemory [7680]uint8         // 0xFE00 (7680 B)
	empty1             [96]uint8           // 0xFEA0 (96 B)
	ioPorts            [76]uint8           // 0xFF00 (67 B)
	empty2             [52]uint8           // 0xFF4C (52 B)
	internalRam         [127]uint8         // 0xFF80 (127 B)
	interruptEnableRegister uint8          // 0xFFFF (1 B)
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

// Specific (Special) memory addresses
const (
	lcdControlAddress      = 0xff40
	lcdStatusAddress       = 0xff41
	currentScanlineAddress = 0xff44
	targetScanlineAddress  = 0xff45
)

func memInit(bootrom []uint8, cartridge []uint8) *memory {
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

	return &memory{
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
			romBanking: true,
			dma_pending_time: 0,
		},
	}
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
		return memory.videoRam[address - 0x8000]
	case switchableRamBank:
		return memory.switchableRamBank[address-0xA000+uint16(memory.memorySettings.currentRAMBank*0x2000)]
	case internalRam8kb:
		return memory.internalRam8kb[address - 0xC000]
	case echoInternalRam8kb:
		return memory.echoInternalRam[address - 0xE000]
	case ioPorts:
		return memory.ioPorts[address-0xff00]
	case internalRam:
		return memory.internalRam[address-0xff80]
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
		return uint16(memory.switchableRamBank[taddr]) | uint16(memory.switchableRamBank[taddr+1])
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
	if address == 0xff44 {
		// Reset the scanline to 0
		fmt.Println("Resetting scanline register 0xff44 to 0")
		return true
	} else if address == 0xff46 {
		fmt.Printf("DMA transfer from %#04x", val)
		for i := 0; i < 0xA0; i++ {
			memory.spriteAttribMemory[i] = memory.read8(uint16(val) + uint16(i))
		}
	}
	return false
}

func (memory *memory) write16(address uint16, val uint16) {
	switch mapAddr(address) {
	case bank0:
		memory.bank0[address] = uint8(val & 0xff)
		memory.bank0[address+1] = uint8(val  >> 8)
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

func (memory *memory) requestInterupt(interuptType int) {
	//fmt.Println("Interrupt requested")

	memory.write8(0xff0f, setBit(memory.read8(0xff0f), uint(interuptType)))
}

func (memory *memory) swapBootRom(cartridge []uint8) {
	fmt.Println("Swapping out bootrom")
	for i := 0; i < 0x100; i += 1 {
		memory.bank0[i] = cartridge[i]
	}
}
