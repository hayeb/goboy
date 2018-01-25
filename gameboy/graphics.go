package gameboy

import (
	"github.com/banthar/Go-SDL/sdl"
	"fmt"
)

const (
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
	renderer      *sdl.Surface
	cartridgeInfo *cartridgeInfo

	screen [144][160]drawColor

	mode      int
	modeclock int
	line      uint8
}

type drawColor struct {
	r uint8
	g uint8
	b uint8
}

func createGraphics(mem *memory, rend *sdl.Surface, ci *cartridgeInfo) *graphics {
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
				graphics.memory.write8(0xff0f, setBit(graphics.memory.read8(0xff0f), 0))
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
	lcdc := graphics.memory.read8(0xFF40)
	if !testBit(lcdc, 7) {
		return
	}

	if testBit(lcdc, 0) {
		graphics.drawBackground()
	}

	if testBit(lcdc, 1) {
		graphics.drawSprites()
	}
}

func (graphics *graphics) drawSprites() {
	var mem = graphics.memory

	bytesInSprite := uint8(16)
	if testBit(mem.read8(0xFF00), 2) {
		fmt.Println("8*16 sprites")
		bytesInSprite = 32
	}
	for i := 0; i < 40; i++ {
		address := i * 4

		// Y top left on screen
		y := mem.spriteAttribMemory[address] - 16
		// X top left on screen
		x := mem.spriteAttribMemory[address+1] - 8

		// Should we draw a line in the current sprite?
		if y > graphics.line || graphics.line >= y+8 {
			continue
		}

		// Tile number
		tileNumber := mem.spriteAttribMemory[address+2]

		// Several options for this sprite
		spriteFlags := mem.spriteAttribMemory[address+3]

		rowAddress := tileNumber * bytesInSprite
		if testBit(spriteFlags, 6) {
			rowAddress += (bytesInSprite/2 - 1) - (graphics.line - y)
		} else {
			rowAddress += graphics.line - y
		}

		for j := uint8(0); j < 8; j++ {
			if x+j < 0 || x+j > 160 {
				continue
			}

			lowerByte := mem.videoRam[rowAddress]
			higherByte := mem.videoRam[rowAddress+1]

			var rowShift = 7 - j
			if testBit(spriteFlags, 5) {
				rowShift = j
			}

			lowerBit := (lowerByte >> rowShift) & 0x1
			higherBit := (higherByte >> rowShift) & 0x1

			colorByte := higherBit<<1 | lowerBit

			paletteAddress := uint16(0xFF48)
			if testBit(spriteFlags, 4) {
				paletteAddress++
			}

			priority := testBit(spriteFlags, 7)
			color := graphics.getColor(colorByte, paletteAddress)
			backgroundColor := graphics.getColor(0, 0xFF47)
			hidden := backgroundColor.r == color.r && backgroundColor.g == color.g && backgroundColor.b == color.b

			if !(priority && hidden) {
				continue
			}
			graphics.screen[graphics.line][x+j] = color
		}
	}
}

func (graphics *graphics) drawBackground() {

	// TODO: Draw window
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

		// The tile are lain out in memory as you would expect, so
		// the 7th bit in a line is the left-most pixel
		lowerBit := lowerByte >> (7 - x) & 0x1
		higherBit := higherByte >> (7 - x) & 0x1

		colorByte := higherBit<<1 | lowerBit

		graphics.screen[graphics.line][i] = graphics.getColor(colorByte, 0xff47)

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
			rect := sdl.Rect{X: int16(i * 4), Y: int16(j * 4), W: 4, H: 4}
			graphics.renderer.FillRect(&rect, sdl.MapRGBA(graphics.renderer.Format, graphics.screen[j][i].r, graphics.screen[j][i].g, graphics.screen[j][i].b, 0xff))
		}
	}
	graphics.renderer.Flip()
}

func (graphics *graphics) isLCDEnabled() bool {
	return testBit(graphics.memory.read8(0xff40), 7)
}

func (graphics *graphics) getColor(colorByte uint8, paletteAddress uint16) drawColor {
	palette := graphics.memory.read8(paletteAddress)

	colorNo := (palette >> (colorByte * 2)) & 0x3

	switch colorNo {
	case 0:
		return drawColor{255, 255, 255}
	case 1:
		return drawColor{0xcc, 0xcc, 0xcc}
	case 2:
		return drawColor{0x77, 0x77, 0x77}
	case 3:
		return drawColor{0, 0, 0}
	default:
		panic("Unknown color!")
	}
}
