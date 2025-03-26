package game

import (
	"image/color"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var Colors = []color.RGBA{
	{20, 20, 20, 255},   // Empty
	{0, 240, 240, 255},  // I
	{0, 100, 255, 255},  // J
	{255, 165, 0, 255},  // L
	{240, 240, 0, 255},  // O
	{50, 205, 50, 255},  // S
	{138, 43, 226, 255}, // T
	{220, 20, 60, 255},  // Z
}

func DrawBlock(screen *ebiten.Image, x, y int, clr color.Color, isCurrent bool) {
	blockClr := clr.(color.RGBA)

	for i := 0; i < BlockSize; i++ {
		alpha := uint8(200 + 55*i/BlockSize)
		lineClr := color.RGBA{blockClr.R, blockClr.G, blockClr.B, alpha}
		vector.StrokeLine(screen, float32(x), float32(y+i), float32(x+BlockSize), float32(y+i), 1, lineClr, false)
	}

	highlight := color.RGBA{
		uint8(min(255, int(blockClr.R)+70)),
		uint8(min(255, int(blockClr.G)+70)),
		uint8(min(255, int(blockClr.B)+70)),
		200,
	}
	vector.StrokeRect(screen, float32(x+1), float32(y+1), float32(BlockSize-2), float32(BlockSize-2), 1, highlight, false)

	vector.StrokeLine(screen, float32(x+1), float32(y+1), float32(x+4), float32(y+1), 1, highlight, false)
	vector.StrokeLine(screen, float32(x+1), float32(y+1), float32(x+1), float32(y+4), 1, highlight, false)
	vector.StrokeLine(screen, float32(x+BlockSize-2), float32(y+1), float32(x+BlockSize-5), float32(y+1), 1, highlight, false)
	vector.StrokeLine(screen, float32(x+BlockSize-2), float32(y+1), float32(x+BlockSize-2), float32(y+4), 1, highlight, false)

	shadow := color.RGBA{
		uint8(max(0, int(blockClr.R)-40)),
		uint8(max(0, int(blockClr.G)-40)),
		uint8(max(0, int(blockClr.B)-40)),
		150,
	}
	vector.StrokeRect(screen, float32(x+2), float32(y+2), float32(BlockSize-4), float32(BlockSize-4), 1, shadow, false)

	if isCurrent {
		pulse := uint8(100 + 155*(1+math.Sin(float64(time.Now().UnixNano())/500000000)/2))
		vector.StrokeRect(screen, float32(x), float32(y), float32(BlockSize), float32(BlockSize), 1, color.RGBA{255, 255, 255, pulse}, false)
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	shakeX := float32(g.ShakeOffset)

	// Draw grid
	for x := 0; x < GridWidth; x++ {
		for y := 0; y < GridHeight; y++ {
			px := float32(GridOffsetX+x*BlockSize) + shakeX
			py := float32(GridOffsetY + y*BlockSize)
			vector.DrawFilledRect(screen, px, py, float32(BlockSize), float32(BlockSize), color.RGBA{40, 40, 50, 255}, false)
			vector.StrokeRect(screen, px, py, float32(BlockSize), float32(BlockSize), 1, color.RGBA{60, 60, 70, 255}, false)
		}
	}

	// Draw placed blocks
	for x := 0; x < GridWidth; x++ {
		for y := 0; y < GridHeight; y++ {
			if g.Grid[x][y] != 0 {
				shouldFlash := false
				for _, line := range g.FlashLines {
					if y == line && g.FlashTimer > 0 {
						shouldFlash = true
						break
					}
				}

				if !shouldFlash {
					px := float32(GridOffsetX+x*BlockSize) + shakeX
					py := float32(GridOffsetY + y*BlockSize)
					DrawBlock(screen, int(px), int(py), Colors[g.Grid[x][y]], false)
				}
			}
		}
	}

	// Draw current piece
	if g.CurrentPiece != nil {
		for y, row := range g.CurrentPiece.Shape {
			for x, val := range row {
				if val != 0 {
					px := float32(GridOffsetX+(g.CurrentPiece.X+x)*BlockSize) + shakeX
					py := float32(GridOffsetY + (g.CurrentPiece.Y+y)*BlockSize)
					if g.CurrentPiece.Y+y >= 0 {
						DrawBlock(screen, int(px), int(py), Colors[g.CurrentPiece.ColorIdx], true)
					}
				}
			}
		}
	}

	// Draw flashing lines
	if g.FlashTimer > 0 {
		flashAlpha := uint8(150 + 105*math.Sin(g.FlashTimer*20))
		for _, y := range g.FlashLines {
			py := float32(GridOffsetY + y*BlockSize)
			vector.DrawFilledRect(screen, float32(GridOffsetX)+shakeX, py,
				float32(GridWidth*BlockSize), float32(BlockSize), color.RGBA{255, 255, 255, flashAlpha}, false)
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
