package gameboy

import "fmt"

type CartridgeTypeCode int
const (
	ROM_ONLY CartridgeTypeCode = iota
	ROM_MBC1
	ROM_MBC1_RAM
	ROM_MBC1_RAM_BAT
	ROM_MBC2
	ROM_MBC2_BATT
	ROM_RAM
	ROM_RAM_BATTERY
	ROM_MMM01
	ROM_MMM01_SRAM
	ROM_MMM01_SRAM_BATT
	ROM_MBC3_RAM
	ROM_MBC3_RAM_BATT
	ROM_MBC5
	ROM_MBC5_RAM
	ROM_MBC5_RAM_BATT
	ROM_MBC5_RUMBLE
	ROM_MBC5_RUMBLE_SRAM
	ROM_MBC5_RUMBLE_SRAM_BATT
	POCKET_CAMERA
	BANDAI_TAMA5
	HUDSON_HUC3
)

type GameBoyType int
const (
	GAMEBOY GameBoyType = iota
	SUPER_GAMEBOY
)

func gameboyType(typeCode uint8) GameBoyType {
	switch (typeCode) {
	case 0x00: return GAMEBOY
	case 0x03: return SUPER_GAMEBOY
	default: panic(fmt.Sprintf("Unknown system typecode %d", typeCode))
	}
}

func gameboyTypeString(gameBoyType GameBoyType) string {
	switch (gameBoyType) {
	case GAMEBOY: return "Gameboy"
	case SUPER_GAMEBOY: return "Super Gameboy"
	default: panic(fmt.Sprintf("No string for unknown system typecode %d", gameBoyType))
	}
}

func typeCode(typeCode uint8) CartridgeTypeCode {
	switch typeCode {
	case 0x0: return ROM_ONLY
	case 0x1: return ROM_MBC1
	case 0x2: return ROM_MBC1_RAM
	case 0x3: return ROM_MBC1_RAM_BAT
	case 0x5: return ROM_MBC2
	case 0x6: return ROM_MBC2_BATT
	case 0x8: return ROM_RAM
	case 0x9: return ROM_RAM_BATTERY
	case 0xB: return ROM_MMM01
	case 0xC: return ROM_MMM01_SRAM
	case 0xD: return ROM_MMM01_SRAM_BATT
	case 0x12: return ROM_MBC3_RAM
	case 0x13: return ROM_MBC3_RAM_BATT
	case 0x19: return ROM_MBC5
	case 0x1A: return ROM_MBC5_RAM
	case 0x1B: return ROM_MBC5_RAM_BATT
	case 0x1C: return ROM_MBC5_RUMBLE
	case 0x1D: return ROM_MBC5_RUMBLE_SRAM
	case 0x1E: return ROM_MBC5_RUMBLE_SRAM_BATT
	case 0x1F: return POCKET_CAMERA
	case 0xFD: return BANDAI_TAMA5
	case 0xFE: return HUDSON_HUC3
	default:
		panic(fmt.Sprintf("Unknown type code %d", typeCode))
	}
}

func typeCodeString(typecode CartridgeTypeCode) string {
	switch (typecode) {
	case ROM_ONLY: return "ROM Only"
	case ROM_MBC1: return "ROM+MBC1"
	case ROM_MBC1_RAM: return "ROM+MBC1+RAM"
	case ROM_MBC1_RAM_BAT: return "ROM+MBC1+RAM+BATT"
	case ROM_MBC2: return "ROM+MBC2"
	case ROM_MBC2_BATT: return "ROM+MBC2+BATT"
	case ROM_RAM: return "ROM+RAM"
	case ROM_RAM_BATTERY: return "ROM+RAM+BATT"
	case ROM_MMM01: return "ROM+MMM01"
	case ROM_MMM01_SRAM: return "ROM+MMM01+SRAM"
	case ROM_MMM01_SRAM_BATT: return "ROM+MMMM01+SRAM+BATT"
	case ROM_MBC3_RAM: return "ROM+MBC3+RAM"
	case ROM_MBC3_RAM_BATT: return "ROM+MNC3+RAM+BATT"
	case ROM_MBC5: return "ROM~MBC5"
	case ROM_MBC5_RAM: return "ROM+MBC5+RAM"
	case ROM_MBC5_RAM_BATT: return "ROM+MBC5+RAM+BATT"
	case ROM_MBC5_RUMBLE: return "ROM+MBC5+RUMBLE"
	case ROM_MBC5_RUMBLE_SRAM: return "ROM+MBC5+RUMBLE+SRAM"
	case ROM_MBC5_RUMBLE_SRAM_BATT: return "ROM+MBC5+RUMBLE+SRAM+BATT"
	case POCKET_CAMERA: return "Pocket Camera"
	case BANDAI_TAMA5: return "Bandai TAMA5"
	case HUDSON_HUC3: return "Hudson HuC-3"
	default:
		panic(fmt.Sprintf("No string for unknown cartridge type %d", typecode ))
	}
}

func romCodeString(romcode int) string {
	switch (romcode) {
	case 0: return "256 Kbit"
	case 1: return "512 Kbit"
	case 2: return "1 Mbit"
	case 3: return "2 Mbit"
	case 4: return "4 Mit"
	case 5: return "8 Mbit"
	case 6: return "16 Mbit"
	default: panic(fmt.Sprintf("No rom size for code %d", romcode))
	}
}

func ramCodeString(ramcode int) string {
	switch (ramcode) {
	case 0: return "None"
	case 1: return "16 Kbit"
	case 2: return "64 Kbit"
	case 3: return "256 Kbit"
	case 4: return "1 Mbit"
	default: panic(fmt.Sprintf("No ram size for code %d", ramcode))

	}
}

func cartridgeTitle(cartridge []byte) string {
	title := ""
	for i := 0x134; i < 0x0142; i++ {
		title += string(cartridge[i])
	}
	return title
}

func localization(code uint8) string {
	if code == 0x1 {
		return "Non-Japanese"
	} else if code == 0x0{
		return "Japanese"
	} else {
		panic(fmt.Sprintf("No localization for code %d", code))
	}
}

type CartridgeInfo struct {
	Name         string
	CartType     string
	System       string
	ROMsize      string
	RAMsize      string
	Localization string
}

func GetCartridgeInfo(cartridge []byte) CartridgeInfo{
	return CartridgeInfo{
		Name: cartridgeTitle(cartridge),
		CartType: typeCodeString(typeCode(cartridge[0x147])),
		System: gameboyTypeString(gameboyType(cartridge[0x146])),
		ROMsize: romCodeString(int(cartridge[0x148])),
		RAMsize: ramCodeString(int(cartridge[0x149])),
		Localization: localization(cartridge[0x14A]),
	}
}
