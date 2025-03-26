package main

import (
	"image/color"
	"tetris/game"
	"tetris/ui"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type App struct {
	game *game.Game
}

func (a *App) Update() error {
	return a.game.Update()
}

func (a *App) Draw(screen *ebiten.Image) {
	// Draw background
	for y := 0; y < game.ScreenHeight; y++ {
		alpha := uint8(30 + y/4)
		lineClr := color.RGBA{15, 15, 20, alpha}
		vector.StrokeLine(screen, 0, float32(y), float32(game.ScreenWidth), float32(y), 1, lineClr, false)
	}

	// Draw game area background
	vector.DrawFilledRect(screen, float32(game.GridOffsetX-5), float32(game.GridOffsetY-5),
		float32(game.GridWidth*game.BlockSize+10), float32(game.GridHeight*game.BlockSize+10), color.RGBA{25, 25, 35, 200}, false)

	// Draw panel background
	vector.DrawFilledRect(screen, float32(game.GridOffsetX+game.GridWidth*game.BlockSize+10), 0,
		float32(ui.PanelWidth), float32(game.ScreenHeight), ui.PanelColor, false)

	// Draw game elements
	a.game.Draw(screen)

	// Draw UI
	ui.DrawPanel(screen, a.game.Score, a.game.Level, a.game.LinesCleared, a.game.NextPiece)

	// Draw game over if needed
	if a.game.GameOver {
		ui.DrawGameOver(screen, a.game.Score)
	}
}

func (a *App) Layout(outsideWidth, outsideHeight int) (int, int) {
	return game.ScreenWidth, game.ScreenHeight
}

func main() {
	// Initialize game constants
	game.ScreenWidth = 700
	game.ScreenHeight = 600
	game.BlockSize = 25
	game.GridOffsetX = 20
	game.GridOffsetY = 30

	app := &App{
		game: &game.Game{},
	}
	app.game.Init()

	ebiten.SetWindowSize(game.ScreenWidth*2, game.ScreenHeight*2)
	ebiten.SetWindowTitle("Tetris in Go")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	if err := ebiten.RunGame(app); err != nil {
		panic(err)
	}
}
