package gameboy

import (
	"github.com/veandco/go-sdl2/sdl"
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
			graphics.drawSprites()
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

	// Where is the screen relative to the background in memory?
	var scY = uint16(graphics.memory.ioPorts[SCROLL_Y])
	var scX = uint16(graphics.memory.ioPorts[SCROLL_X])

	// We adjust the tilemap address. There are 32*32 tiles total in the background,
	// and every tile has 8 lines.
	tileMapAddress = tileMapAddress + (scY+uint16(graphics.line))/8*32

	// The x and y values of the point in the tile
	var y = (uint8(scY) + graphics.line) % 8
	var x = uint8(scX) % 8

	// The offset in the current line of tiles according to the x-scroll,
	// there are 8 pixels width in a tile
	var offsetInLine = scX / 8

	tileNumber := graphics.memory.videoRam[tileMapAddress+offsetInLine]
	for i := 0; i < 160; i++ {
		var dataAddr = uint16(tileNumber)*16 + uint16(y*2)

		lowerByte := graphics.memory.videoRam[dataAddr]
		higherByte := graphics.memory.videoRam[dataAddr+1]

		// The tile are layn out in memory as you would expect, so
		// the 7th bit in a line is the left-most pixel
		lowerBit := lowerByte >> (7 - x) & 0x1
		higherBit := higherByte >> (7 - x) & 0x1

		colorByte := higherBit<<1 | lowerBit

		colour := graphics.getColor(colorByte, 0xff47)

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
			offsetInLine = offsetInLine + 1
			tileNumber = graphics.memory.videoRam[tileMapAddress+offsetInLine]
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
	// TODO: Properly handle palettes/colors
	if n > 0 {
		return 0
	}
	return 1
}

func (graphics *graphics) drawSprites() {

}
