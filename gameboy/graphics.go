package gameboy

import (
	"github.com/veandco/go-sdl2/sdl"
	"fmt"
)

const (
	WHITE = 1
	LIGHT_GRAY = 2
	DARK_GRAY = 3
	BLACK = 4
)

type graphics struct {
	memory        *memory
	renderer      *sdl.Renderer
	cartridgeInfo *cartridgeInfo

	scanlineCounter int
	lcd_status 	int
}

func createGraphics(mem *memory, rend *sdl.Renderer, ci *cartridgeInfo) *graphics {
	return &graphics{
		memory:        mem,
		renderer:      rend,
		cartridgeInfo: ci,
		scanlineCounter: 456,
	}
}

func (graphics *graphics) updateGraphics(instructionLength int) {
	graphics.setLCDStatus()

	if testBit(graphics.memory.read8(0xff40), 7) {
		graphics.scanlineCounter -= instructionLength
	}

	if graphics.memory.read8(0xff44) > 0x99 {
		graphics.memory.io_ports[0xff44-0xFF00] = 0
	}

	if graphics.scanlineCounter <= 0 {
		graphics.drawCurrentLine()
	}
}

func (graphics *graphics) setLCDStatus() {
	status := graphics.memory.read8(0xFF41)

	if !graphics.isLCDEnabled() {
		graphics.scanlineCounter = 456
		graphics.memory.io_ports[0xff44 - 0xFF00] = 0
		graphics.lcd_status &= 252
		status = setBit(status, 0)
		graphics.memory.write8(0xff41, status)
		return
	}
	currentLine := graphics.memory.read8(0xff44)
	currentMode := status & 0x3

	mode := uint8(0)
	reqInt := false

	if currentLine >= 0x90 {
		mode = 1
		status = setBit(status, 0)
		status = resetBit(status, 1)
		reqInt = testBit(status, 4)
	} else {
		mode2bounds := 456 - 80
		mode3bounds := mode2bounds - 172

		if (graphics.scanlineCounter >= mode2bounds) {
			mode = 2
			status = setBit(status, 1)
			status = resetBit(status, 0)
			reqInt = testBit(status, 5)
		} else if (graphics.scanlineCounter >= mode3bounds) {
			mode = 3
			status = setBit(status, 1)
			status = setBit(status, 0)
		} else {
			mode = 0
			status = resetBit(status, 1)
			status = resetBit(status, 0)
			reqInt = testBit(status, 3)
		}
	}

	if reqInt && (mode != currentMode) {
		graphics.memory.requestInterupt(1)
	}

	if currentLine == graphics.memory.read8(0xff45) {
		status = setBit(status, 2)
		if testBit(status, 6) {
			graphics.memory.requestInterupt(1)
		}
	} else {
		status = resetBit(status, 2)
	}
	graphics.memory.write8(0xff41, status)
}

func (graphics *graphics) isLCDEnabled() bool {
	return testBit(graphics.memory.read8(0xff40), 7)
}

func (graphics *graphics) drawCurrentLine() {
	control := graphics.memory.read8(0xff40)

	if !testBit(control, 7) {
		return
	}

	graphics.memory.io_ports[0xff44-0xFF00] = graphics.memory.io_ports[0xff44-0xFF00] + 1
	graphics.scanlineCounter = 456

	scanLine := graphics.memory.read8(0xff44)
	if scanLine == 0x90 {
		graphics.memory.requestInterupt(0)
	}

	if scanLine > 0x99 {
		graphics.memory.write8(0xff44, 0)
	}

	if scanLine < 0x90 {
		graphics.drawScanLine()
	}
}

func (graphics *graphics) drawScanLine() {
	var lcdControl uint8 = graphics.memory.read8(0xff40)
	if testBit(lcdControl, 7) {
		graphics.renderTiles(lcdControl)
		graphics.renderSprites(lcdControl)
	}
}

func (graphics *graphics) renderTiles(lcdControl uint8) {
	if !testBit(lcdControl, 0) {
		return
	}

	var tileData uint16 = 0
	var backgroundMemory uint16 = 0
	var unsig bool = true

	//var sY uint8 = graphics.memory.read8(0xff42)
	var sY uint8 = 0
	//fmt.Printf("sY: %#02x\n", sY)
	var sX uint8 = graphics.memory.read8(0xff43)
	fmt.Println("ScrolyX: ", sX)
	fmt.Println("ScrolyY: ", sY)
	var wY uint8 = graphics.memory.read8(0xff4a)
	var wX uint8 = graphics.memory.read8(0xff4b) - 7
	usingWindow := false

	if testBit(lcdControl, 5) && wY <= graphics.memory.read8(0xff44) {
		usingWindow = true
	}

	if testBit(lcdControl, 4) {
		tileData = 0x8000
	} else {
		tileData = 0x8800
		unsig = false
	}

	if !usingWindow {
		if testBit(lcdControl, 3) {
			backgroundMemory = 0x9c00
		} else {
			backgroundMemory = 0x9800
		}
	} else {
		if testBit(lcdControl, 6) {
			backgroundMemory = 0x9c00
		} else {
			backgroundMemory = 0x9800
		}
	}

	var yPos uint8 = 0

	if !usingWindow {
		yPos = sY + graphics.memory.read8(0xff44)
	} else {
		yPos = graphics.memory.read8(0xff44) - wY
	}

	var tileRow uint16 = uint16(yPos/8 * 32)

	for pixel := uint8(0); pixel < 160; pixel += 1 {
		var xPos uint8 = pixel + sX

		if usingWindow && pixel >= wX {
			xPos = pixel - wX
		}

		var tileCol uint16 = uint16(xPos) / 8
		var tileNum int16 = 0

		if unsig {
			b := graphics.memory.read8(backgroundMemory + tileRow + tileCol)
			//fmt.Printf("B: %#02x (%#04x + %#04x + %#04x)\n", b, backgroundMemory, tileRow, tileCol)
			tileNum = int16(uint16(b))
		} else {
			tileNum = int16(int8(graphics.memory.read8(backgroundMemory + tileRow + tileCol)))
		}

		var tileLocation uint16 = tileData

		if unsig {
			tileLocation += uint16(int(tileNum) * 16)
		} else {
			tileLocation += uint16(int((tileNum + 128) * 16))
		}

		var line uint8 = 2 * (yPos % 8)
		//fmt.Printf("Reading tile from %#04x, %#04x\n", tileLocation + uint16(line), tileLocation + uint16(line)+1)
		d1 := graphics.memory.read8(tileLocation + uint16(line))
		d2 := graphics.memory.read8(tileLocation + uint16(line) + 1)

		var colourBit int = int((int(xPos) % 8) - 7) * -1
		var colourNum uint8 = getBitN(d2, uint(colourBit)) << 1 | getBitN(d1, uint(colourBit))
		var colour int = graphics.getColor(colourNum, 0xff47)

		r, g, b := uint8(0), uint8(0), uint8(0)

		switch colour {
		case WHITE:
			r, g, b = 255, 255, 255
		case LIGHT_GRAY:
			r, g, b = 0xcc, 0xcc, 0xcc
		case DARK_GRAY:
			r, g, b = 0x77, 0x77, 0x77
		}
		f := graphics.memory.read8(0xff44)

		graphics.renderer.SetDrawColor(r, g, b, 255)
		graphics.renderer.DrawPoint(int(pixel), int(f))
	}
}

func (graphics *graphics) getColor(n uint8, a uint16) int {
	p := graphics.memory.read8(a)
	hi := uint8(0)
	lo := uint8(0)

	switch n {
	case 0: hi, lo = 1, 0
	case 1: hi, lo = 3, 2
	case 2: hi, lo = 5, 4
	case 3: hi, lo = 7, 6
	}

	colour := getBitN(p, uint(hi)) << 1 | getBitN(p, uint(lo))
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
		yPos := graphics.memory.read8(0xfe00 + uint16(index)) - 16
		xPos := graphics.memory.read8(0xfe00 + uint16(index) + 1) - 8
		tileLocation := graphics.memory.read8(0xfe00 + uint16(index) + 2)
		attributes := graphics.memory.read8(0xfe00 + uint16(index) + 3)

		yFlip := testBit(attributes, 6)
		xFlip := testBit(attributes, 5)

		scanLine := graphics.memory.read8(0xff44)

		ysize := 8

		if use8x16 {
			ysize = 16
		}

		if scanLine >= yPos && scanLine < yPos + uint8(ysize) {
			var line int = int(scanLine - yPos)
			if yFlip {
				line = -1 * (line - ysize)
			}
			line = line * 2
			d1 := graphics.memory.read8(0x8000 + uint16(int((tileLocation * 16)) + line))
			d2 := graphics.memory.read8(0x8000 + uint16(int((tileLocation * 16)) + line + 1))

			for tilePixel := 7; tilePixel >=0; tilePixel-- {
				colourBit := tilePixel
				if xFlip {
					colourBit = -1 * (colourBit - 7)
				}
				var colourNum uint8 = ((d2 >> uint8(colourBit) & 0x1) << 1) | ((d1 >> uint8(colourBit)) & 0x1)
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
				case WHITE: r, g, b = 255, 255, 255
				case LIGHT_GRAY: r, g, b = 0xcc, 0xcc, 0xcc
				case DARK_GRAY: r, g, b = 0x77, 0x77, 0x77
				}

				xPix := 0 - tilePixel + 7
				pixel := int(xPos) + xPix

				graphics.renderer.SetDrawColor(r, g,b, 255)
				graphics.renderer.DrawPoint(int(scanLine), pixel)
			}
		}

	}
}
