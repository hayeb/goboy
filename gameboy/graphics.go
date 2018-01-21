package gameboy

import (
	"github.com/veandco/go-sdl2/sdl"
	"fmt"
)

const (
	WHITE      = 1
	LIGHT_GRAY = 2
	DARK_GRAY  = 3
	BLACK      = 4

	ADDRESS_IO_PORTS  = 0xFF00
	ADDRESS_VIDEO_RAM = 0x8000

	LCDC      uint16 = 0xFF40 - ADDRESS_IO_PORTS
	LCDC_STAT uint16 = 0xFF41 - ADDRESS_IO_PORTS
	SCROLL_Y  uint16 = 0xFF42 - ADDRESS_IO_PORTS
	SCROLL_X  uint16 = 0xFF43 - ADDRESS_IO_PORTS
	SCANLINE  uint16 = 0xFF44 - ADDRESS_IO_PORTS
)

type graphics struct {
	memory        *memory
	renderer      *sdl.Renderer
	cartridgeInfo *cartridgeInfo

	screen [144][160]drawColor

	mode      int
	modeclock int
	line      uint8
}

type drawColor struct {
	r int
	g int
	b int
}

func createGraphics(mem *memory, rend *sdl.Renderer, ci *cartridgeInfo) *graphics {
	return &graphics{
		memory:        mem,
		renderer:      rend,
		cartridgeInfo: ci,
	}
}

func (graphics *graphics) updateGraphics(instructionLength int) {
	graphics.modeclock += instructionLength

	graphics.memory.ioPorts[SCANLINE] = graphics.line

	// TODO: Write the LCD status to memory
	switch graphics.mode {
	case 0:
		// HBLANK mode:
		if graphics.modeclock >= 204 {
			graphics.modeclock = 0
			graphics.line += 1

			// If this is the last line, enter VBLANK
			if graphics.line == 143 {
				graphics.mode = 1
				graphics.showData()
			} else {
				graphics.mode = 2
			}
		}
	case 1:
		// VBLANK mode:
		if graphics.modeclock >= 456 {
			graphics.modeclock = 0
			graphics.line++

			if graphics.line > 153 {
				graphics.mode = 2
				graphics.line = 0
			}
		}
	case 2:
		if graphics.modeclock >= 80 {
			graphics.modeclock = 0
			graphics.mode = 3
		}
	case 3:
		if graphics.modeclock >= 172 {
			graphics.modeclock = 0
			graphics.mode = 0

			graphics.drawCurrentLine()
		}
	}
}

/**
Draws current line to the screen buffer
 */
func (graphics *graphics) drawCurrentLine() {
	var tileMapAddress uint16 = 0x1800

	if testBit(graphics.memory.ioPorts[LCDC], 3) {
		tileMapAddress = 0x1C00
	}

	// Adjust for the current line and the current place of the screen in the background
	var scY = uint16(graphics.memory.ioPorts[SCROLL_Y])
	var scX = uint16(graphics.memory.ioPorts[SCROLL_X])
	tileMapAddress += ((scY + uint16(graphics.line)) >> 3) << 5

	// The offset in the current line of tiles according to the x-scroll
	var lineOffset = scX >> 3

	// The x and y values of the point in the background
	var y = uint8(graphics.line + uint8(scY))
	var x = uint8(scX)

	tileNumber := int16(graphics.memory.videoRam[tileMapAddress+lineOffset])

	for i := 0; i < 160; i++ {
		var dataAddr = uint16(tileNumber * 16)

		lowerByte := graphics.memory.videoRam[dataAddr+uint16((y%8)*2)]
		higherByte := graphics.memory.videoRam[dataAddr+uint16((y%8)*2)+1]

		var colourNum = ((higherByte >> x & 0x1) << 1) | ((lowerByte >> x) & 0x1)

		colour := graphics.getColor(colourNum, 0xff47)

		r, g, b := 0, 0, 0
		switch colour {
		case WHITE:
			r, g, b = 255, 255, 255
		case LIGHT_GRAY:
			r, g, b = 0xcc, 0xcc, 0xcc
		case DARK_GRAY:
			r, g, b = 0x77, 0x77, 0x77
		}

		graphics.screen[graphics.line][i] = drawColor{
			r: r,
			g: g,
			b: b,
		}

		x++
		if x == 8 {
			x = 0
			lineOffset = (lineOffset + 1) & 31
			var tileAddress = tileMapAddress + lineOffset
			tileNumber = int16(graphics.memory.videoRam[tileAddress])
		}
	}
}

func (graphics *graphics) showData() {
	for j := 0; j < len(graphics.screen); j++ {
		for i := 0; i < len(graphics.screen[0]); i++ {
			graphics.renderer.SetDrawColor(uint8(graphics.screen[j][i].r), uint8(graphics.screen[j][i].g), uint8(graphics.screen[j][i].b), 255)
			graphics.renderer.DrawPoint(i, j)
		}
	}
	graphics.renderer.Present()
}

func (graphics *graphics) isLCDEnabled() bool {
	return testBit(graphics.memory.read8(0xff40), 7)
}

func (graphics *graphics) getColor(n uint8, a uint16) int {
	p := graphics.memory.read8(a)
	hi := uint8(0)
	lo := uint8(0)

	switch n {
	case 0:
		hi, lo = 1, 0
	case 1:
		hi, lo = 3, 2
	case 2:
		hi, lo = 5, 4
	case 3:
		hi, lo = 7, 6
	}

	colour := getBitN(p, uint(hi))<<1 | getBitN(p, uint(lo))
	switch colour {
	case 0:
		return WHITE
	case 1:
		return LIGHT_GRAY
	case 2:
		return DARK_GRAY
	case 3:
		return BLACK
	default:
		panic(fmt.Sprintf("Could not determine colour for value %#02x", colour))
	}
}

func (graphics *graphics) renderSprites(lcdControl uint8) {
	if !testBit(lcdControl, 1) {
		return
	}

	use8x16 := false

	if testBit(lcdControl, 2) {
		use8x16 = true
	}

	for sprite := 0; sprite < 40; sprite += 1 {
		index := sprite * 4
		yPos := graphics.memory.read8(0xfe00+uint16(index)) - 16
		xPos := graphics.memory.read8(0xfe00+uint16(index)+1) - 8
		tileLocation := graphics.memory.read8(0xfe00 + uint16(index) + 2)
		attributes := graphics.memory.read8(0xfe00 + uint16(index) + 3)

		yFlip := testBit(attributes, 6)
		xFlip := testBit(attributes, 5)

		scanLine := graphics.memory.read8(0xff44)

		ysize := 8

		if use8x16 {
			ysize = 16
		}

		if scanLine >= yPos && scanLine < yPos+uint8(ysize) {
			var line = int(scanLine - yPos)
			if yFlip {
				line = -1 * (line - ysize)
			}
			line = line * 2
			d1 := graphics.memory.read8(0x8000 + uint16(int((tileLocation * 16))+line))
			d2 := graphics.memory.read8(0x8000 + uint16(int((tileLocation * 16))+line+1))

			for tilePixel := 7; tilePixel >= 0; tilePixel-- {
				colourBit := tilePixel
				if xFlip {
					colourBit = -1 * (colourBit - 7)
				}
				var colourNum = ((d2 >> uint8(colourBit) & 0x1) << 1) | ((d1 >> uint8(colourBit)) & 0x1)
				address := 0xff48
				if testBit(attributes, 4) {
					address = 0xff49
				}

				colour := graphics.getColor(colourNum, uint16(address))
				if colour == WHITE {
					continue
				}
				r, g, b := uint8(0), uint8(0), uint8(0)
				switch colour {
				case WHITE:
					r, g, b = 255, 255, 255
				case LIGHT_GRAY:
					r, g, b = 0xcc, 0xcc, 0xcc
				case DARK_GRAY:
					r, g, b = 0x77, 0x77, 0x77
				}

				xPix := 0 - tilePixel + 7
				pixel := int(xPos) + xPix

				graphics.renderer.SetDrawColor(r, g, b, 255)
				graphics.renderer.DrawPoint(int(scanLine), pixel)
			}
		}
	}
}
