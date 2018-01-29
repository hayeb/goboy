package gameboy

import (
	"bytes"
	"fmt"
)

type cartridgeInfo struct {
	Name         string
	CartType     cartridgeTypeCode
	System       gameBoyType
	romSize      romSizeCode
	ramSize      ramSizeCode
	Localization string
}

type cartridgeTypeCode int

const (
	rom_only cartridgeTypeCode = iota
	rom_mbc1
	rom_mbc1_ram
	rom_mbc1_ram_bat
	rom_mbc2
	rom_mbc2_batt
	rom_ram
	rom_ram_battery
	rom_mmm01
	rom_mmm01_sram
	rom_mmm01_sram_batt
	rom_mbc3_ram
	rom_mbc3_ram_batt
	rom_mbc5
	rom_mbc5_ram
	rom_mbc5_ram_batt
	rom_mbc5_rumble
	rom_mbc5_rumble_sram
	rom_mbc5_rumble_sram_batt
	pocket_camera
	bandai_tama5
	hudson_huc3
)

type ramSizeCode int

const (
	ram_none ramSizeCode = iota
	ram_kbit_16
	ram_kbit_64
	ram_kbit_256
	ram_mbit_1
)

type romSizeCode int

const (
	rom_kbit_256 romSizeCode = iota
	rom_kbit_512
	rom_mbit_1
	rom_mbit_2
	rom_mbit_4
	rom_mbit_8
	rom_mbit_16
)

type gameBoyType int

const (
	type_gameboy gameBoyType = iota
	type_super_gameboy
)

func (cartInfo *cartridgeInfo) gameboyTypeString() string {
	switch cartInfo.System {
	case type_gameboy:
		return "Gameboy"
	case type_super_gameboy:
		return "Super Gameboy"
	default:
		return ""
	}
}

func (cartInfo *cartridgeInfo) cartridgeTypeCodeString() string {
	switch cartInfo.CartType {
	case rom_only:
		return "ROM Only"
	case rom_mbc1:
		return "ROM+MBC1"
	case rom_mbc1_ram:
		return "ROM+MBC1+RAM"
	case rom_mbc1_ram_bat:
		return "ROM+MBC1+RAM+BATT"
	case rom_mbc2:
		return "ROM+MBC2"
	case rom_mbc2_batt:
		return "ROM+MBC2+BATT"
	case rom_ram:
		return "ROM+RAM"
	case rom_ram_battery:
		return "ROM+RAM+BATT"
	case rom_mmm01:
		return "ROM+MMM01"
	case rom_mmm01_sram:
		return "ROM+MMM01+SRAM"
	case rom_mmm01_sram_batt:
		return "ROM+MMMM01+SRAM+BATT"
	case rom_mbc3_ram:
		return "ROM+MBC3+RAM"
	case rom_mbc3_ram_batt:
		return "ROM+MNC3+RAM+BATT"
	case rom_mbc5:
		return "ROM~MBC5"
	case rom_mbc5_ram:
		return "ROM+MBC5+RAM"
	case rom_mbc5_ram_batt:
		return "ROM+MBC5+RAM+BATT"
	case rom_mbc5_rumble:
		return "ROM+MBC5+RUMBLE"
	case rom_mbc5_rumble_sram:
		return "ROM+MBC5+RUMBLE+SRAM"
	case rom_mbc5_rumble_sram_batt:
		return "ROM+MBC5+RUMBLE+SRAM+BATT"
	case pocket_camera:
		return "Pocket Camera"
	case bandai_tama5:
		return "Bandai TAMA5"
	case hudson_huc3:
		return "Hudson HuC-3"
	default:
		return ""
	}
}

func (cartInfo *cartridgeInfo) romSizeCodeString() string {
	switch cartInfo.romSize {
	case rom_kbit_256:
		return "256 Kbit"
	case rom_kbit_512:
		return "512 Kbit"
	case rom_mbit_1:
		return "1 Mbit"
	case rom_mbit_2:
		return "2 Mbit"
	case rom_mbit_4:
		return "4 Mit"
	case rom_mbit_8:
		return "8 Mbit"
	case rom_mbit_16:
		return "16 Mbit"
	default:
		return ""
	}
}

func (cartInfo *cartridgeInfo) ramSizeCodeString() string {
	switch cartInfo.ramSize {
	case ram_none:
		return "None"
	case ram_kbit_16:
		return "16 Kbit"
	case ram_kbit_64:
		return "64 Kbit"
	case ram_kbit_256:
		return "256 Kbit"
	case ram_mbit_1:
		return "1 Mbit"
	default:
		return ""
	}
}

func gameboyType(typeCode uint8) gameBoyType {
	switch typeCode {
	case 0x00:
		return type_gameboy
	case 0x03:
		return type_super_gameboy
	default:
		panic(fmt.Sprintf("Unknown system typecode %d", typeCode))
	}
}

func typeCode(typeCode uint8) cartridgeTypeCode {
	switch typeCode {
	case 0x0:
		return rom_only
	case 0x1:
		return rom_mbc1
	case 0x2:
		return rom_mbc1_ram
	case 0x3:
		return rom_mbc1_ram_bat
	case 0x5:
		return rom_mbc2
	case 0x6:
		return rom_mbc2_batt
	case 0x8:
		return rom_ram
	case 0x9:
		return rom_ram_battery
	case 0xB:
		return rom_mmm01
	case 0xC:
		return rom_mmm01_sram
	case 0xD:
		return rom_mmm01_sram_batt
	case 0x12:
		return rom_mbc3_ram
	case 0x13:
		return rom_mbc3_ram_batt
	case 0x19:
		return rom_mbc5
	case 0x1A:
		return rom_mbc5_ram
	case 0x1B:
		return rom_mbc5_ram_batt
	case 0x1C:
		return rom_mbc5_rumble
	case 0x1D:
		return rom_mbc5_rumble_sram
	case 0x1E:
		return rom_mbc5_rumble_sram_batt
	case 0x1F:
		return pocket_camera
	case 0xFD:
		return bandai_tama5
	case 0xFE:
		return hudson_huc3
	default:
		panic(fmt.Sprintf("Unknown type code %d", typeCode))
	}
}

func uint8ToromSizeCode(romcode uint8) romSizeCode {
	switch romcode {
	case 0:
		return rom_kbit_256
	case 1:
		return rom_kbit_512
	case 2:
		return rom_mbit_1
	case 3:
		return rom_mbit_2
	case 4:
		return rom_mbit_4
	case 5:
		return rom_mbit_8
	case 6:
		return rom_mbit_16
	default:
		panic(fmt.Sprintf("Unknown ROM size code %d", romcode))
	}
}

func uint8ToramSizeCode(ramSizeCode uint8) ramSizeCode {
	switch ramSizeCode {
	case 0:
		return ram_none
	case 1:
		return ram_kbit_16
	case 2:
		return ram_kbit_64
	case 3:
		return ram_kbit_256
	case 4:
		return ram_mbit_1
	default:
		panic(fmt.Sprintf("Unknown RAM size code %d", ramSizeCode))
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
	} else if code == 0x0 {
		return "Japanese"
	} else {
		panic(fmt.Sprintf("No localization for code %d", code))
	}
}

func createCartridgeInfo(cartridge []byte) *cartridgeInfo {
	return &cartridgeInfo{
		Name:         cartridgeTitle(cartridge),
		CartType:     typeCode(cartridge[0x147]),
		System:       gameboyType(cartridge[0x146]),
		romSize:      uint8ToromSizeCode(cartridge[0x148]),
		ramSize:      uint8ToramSizeCode(cartridge[0x149]),
		Localization: localization(cartridge[0x14A]),
	}
}

func cartridgeInfoString(cartridgeInfo cartridgeInfo) string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("Cartridge name: %s\n", cartridgeInfo.Name))
	buffer.WriteString(fmt.Sprintf("Cartridge type: %s\n", cartridgeInfo.cartridgeTypeCodeString()))
	buffer.WriteString(fmt.Sprintf("System type: %s\n", cartridgeInfo.gameboyTypeString()))
	buffer.WriteString(fmt.Sprintf("ROM size: %s\n", cartridgeInfo.romSizeCodeString()))
	buffer.WriteString(fmt.Sprintf("RAM size: %s\n", cartridgeInfo.ramSizeCodeString()))
	buffer.WriteString(fmt.Sprintf("Localization: %s\n", cartridgeInfo.Localization))

	return buffer.String()
}
