package ui

import (
	"image/color"
	"strconv"
	"tetris/game"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/basicfont"
)

const (
	PanelWidth   = 140 // Reduced from 150
	PanelPadding = 10  // Added padding constant
)

var (
	PanelColor     = color.RGBA{25, 25, 35, 255}
	HighlightColor = color.RGBA{60, 60, 80, 255}
)

func DrawPanel(screen *ebiten.Image, score, level, lines int, nextPiece *game.Piece) {
	panelX := game.GridOffsetX + game.GridWidth*game.BlockSize + PanelPadding

	// Game logo - made more compact
	logo := "TETRIS"
	for i, c := range logo {
		clr := game.Colors[(i%6)+1]
		text.Draw(screen, string(c), basicfont.Face7x13, panelX+i*12, 30, clr) // Reduced spacing
	}

	// Next piece preview - more compact
	nextX, nextY := panelX, 50 // Reduced vertical space
	text.Draw(screen, "NEXT:", basicfont.Face7x13, nextX, nextY-5, color.White)

	previewWidth := 80  // Reduced from 100
	previewHeight := 60 // Reduced from 80
	vector.DrawFilledRect(screen, float32(nextX), float32(nextY),
		float32(previewWidth), float32(previewHeight), HighlightColor, false)
	vector.StrokeRect(screen, float32(nextX), float32(nextY),
		float32(previewWidth), float32(previewHeight), 1, color.White, false)

	if nextPiece != nil {
		pieceWidth := len(nextPiece.Shape[0]) * (game.BlockSize - 5) // Smaller preview
		pieceHeight := len(nextPiece.Shape) * (game.BlockSize - 5)
		startX := nextX + (previewWidth-pieceWidth)/2
		startY := nextY + (previewHeight-pieceHeight)/2

		for y, row := range nextPiece.Shape {
			for x, val := range row {
				if val != 0 {
					px, py := startX+x*(game.BlockSize-5), startY+y*(game.BlockSize-5)
					game.DrawBlock(screen, px, py, game.Colors[nextPiece.ColorIdx], false)
				}
			}
		}
	}

	// Stats display - more compact layout
	statsStartY := nextY + previewHeight + 15 // Reduced spacing
	drawStat(screen, panelX, statsStartY, "SCORE", strconv.Itoa(score))
	drawStat(screen, panelX, statsStartY+25, "LEVEL", strconv.Itoa(level)) // Reduced spacing
	drawStat(screen, panelX, statsStartY+50, "LINES", strconv.Itoa(lines)) // Reduced spacing

	// Controls - more compact
	controls := []string{
		"CONTROLS",
		"←→ : Move",
		"↑ : Rotate",
		"↓ : Soft Drop",
		"SPACE : Drop",
	}
	controlsStartY := statsStartY + 85 // Adjusted position
	for i, ctrl := range controls {
		text.Draw(screen, ctrl, basicfont.Face7x13, panelX, controlsStartY+i*15, color.White) // Reduced line height
	}
}

func drawStat(screen *ebiten.Image, x, y int, label, value string) {
	text.Draw(screen, label+":", basicfont.Face7x13, x, y, color.White)
	text.Draw(screen, value, basicfont.Face7x13, x+PanelWidth-30, y, color.White) // Right-aligned
}

func DrawGameOver(screen *ebiten.Image, score int) {
	overlay := color.RGBA{0, 0, 0, 200}
	gameWidth := game.GridWidth * game.BlockSize
	gameHeight := game.GridHeight * game.BlockSize

	vector.DrawFilledRect(screen,
		float32(game.GridOffsetX), float32(game.GridOffsetY),
		float32(gameWidth), float32(gameHeight), overlay, false)

	centerX := game.GridOffsetX + gameWidth/2
	centerY := game.GridOffsetY + gameHeight/2

	text.Draw(screen, "GAME OVER", basicfont.Face7x13, centerX-35, centerY-20, color.White)
	text.Draw(screen, "Score: "+strconv.Itoa(score), basicfont.Face7x13, centerX-30, centerY, color.White)
	text.Draw(screen, "SPACE to restart", basicfont.Face7x13, centerX-45, centerY+20, color.White)
}
