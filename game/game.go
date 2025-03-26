package game

import (
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var (
	BlockSize    = 30
	ScreenWidth  = 700
	ScreenHeight = 600
	GridOffsetX  = 20
	GridOffsetY  = 20
)

const (
	GridWidth  = 10
	GridHeight = 20
)

type Game struct {
	Grid          [GridWidth][GridHeight]int
	CurrentPiece  *Piece
	NextPiece     *Piece
	GameOver      bool
	Score         int
	Level         int
	LinesCleared  int
	DropSpeed     int
	LastDropTime  time.Time
	DropInterval  time.Duration
	MoveDownDelay time.Duration
	LastMoveDown  time.Time
	FlashTimer    float64
	FlashLines    []int
	ShakeOffset   float64
	ShakeTimer    float64
	RotationState int
}

func (g *Game) Init() {
	rand.Seed(time.Now().UnixNano())
	g.CurrentPiece = NewPiece(rand.Intn(7))
	g.NextPiece = NewPiece(rand.Intn(7))
	g.GameOver = false
	g.Score = 0
	g.Level = 1
	g.LinesCleared = 0
	g.DropSpeed = 1000
	g.DropInterval = time.Duration(g.DropSpeed) * time.Millisecond
	g.LastDropTime = time.Now()
	g.MoveDownDelay = 100 * time.Millisecond
	g.LastMoveDown = time.Now()

	for x := 0; x < GridWidth; x++ {
		for y := 0; y < GridHeight; y++ {
			g.Grid[x][y] = 0
		}
	}
}

func (g *Game) validMove(p *Piece, x, y int, rotatedShape [][]int) bool {
	if rotatedShape == nil {
		rotatedShape = p.Shape
	}

	for py, row := range rotatedShape {
		for px, val := range row {
			if val != 0 {
				newX, newY := p.X+x+px, p.Y+y+py
				if newX < 0 || newX >= GridWidth || newY >= GridHeight {
					return false
				}
				if newY >= 0 && g.Grid[newX][newY] != 0 {
					return false
				}
			}
		}
	}
	return true
}

func (g *Game) rotatePiece() {
	if g.CurrentPiece == nil || g.GameOver {
		return
	}

	oldX := g.CurrentPiece.X
	oldY := g.CurrentPiece.Y
	oldShape := g.CurrentPiece.Shape
	oldRotation := g.RotationState

	newRotation := (g.RotationState + 1) % 4
	if g.CurrentPiece.ColorIdx-1 == 3 {
		newRotation = 0
	} else if g.CurrentPiece.ColorIdx-1 == 0 ||
		g.CurrentPiece.ColorIdx-1 == 4 ||
		g.CurrentPiece.ColorIdx-1 == 6 {
		newRotation = (g.RotationState + 1) % 2
	}

	var newShape [][]int
	switch g.CurrentPiece.ColorIdx - 1 {
	case 0:
		newShape = pieceShapes[newRotation%2]
	case 1:
		newShape = pieceShapes[2+newRotation%4]
	case 2:
		newShape = pieceShapes[6+newRotation%4]
	case 3:
		newShape = pieceShapes[10]
	case 4:
		newShape = pieceShapes[11+newRotation%2]
	case 5:
		newShape = pieceShapes[13+newRotation%4]
	case 6:
		newShape = pieceShapes[17+newRotation%2]
	}

	g.CurrentPiece.Shape = newShape
	g.RotationState = newRotation
	g.CurrentPiece.X += (len(oldShape[0]) - len(newShape[0])) / 2

	if !g.validMove(g.CurrentPiece, 0, 0, nil) {
		if g.validMove(g.CurrentPiece, -1, 0, nil) {
			g.CurrentPiece.X--
		} else if g.validMove(g.CurrentPiece, 1, 0, nil) {
			g.CurrentPiece.X++
		} else if g.validMove(g.CurrentPiece, 0, -1, nil) {
			g.CurrentPiece.Y--
		} else {
			g.CurrentPiece.X = oldX
			g.CurrentPiece.Y = oldY
			g.CurrentPiece.Shape = oldShape
			g.RotationState = oldRotation
		}
	}
}

func (g *Game) mergePiece() {
	for y, row := range g.CurrentPiece.Shape {
		for x, val := range row {
			if val != 0 {
				gridX, gridY := g.CurrentPiece.X+x, g.CurrentPiece.Y+y
				if gridY >= 0 {
					g.Grid[gridX][gridY] = g.CurrentPiece.ColorIdx
				}
			}
		}
	}
}

func (g *Game) clearLines() {
	g.FlashLines = nil
	linesCleared := 0

	for y := GridHeight - 1; y >= 0; y-- {
		lineFull := true
		for x := 0; x < GridWidth; x++ {
			if g.Grid[x][y] == 0 {
				lineFull = false
				break
			}
		}

		if lineFull {
			g.FlashLines = append(g.FlashLines, y)
			linesCleared++
		}
	}

	if linesCleared > 0 {
		g.FlashTimer = 0.5
		g.ShakeTimer = 0.3
	}
}

func (g *Game) removeLines() {
	linesCleared := 0
	for y := GridHeight - 1; y >= 0; y-- {
		lineFull := true
		for x := 0; x < GridWidth; x++ {
			if g.Grid[x][y] == 0 {
				lineFull = false
				break
			}
		}

		if lineFull {
			linesCleared++
			for ny := y; ny > 0; ny-- {
				for x := 0; x < GridWidth; x++ {
					g.Grid[x][ny] = g.Grid[x][ny-1]
				}
			}
			for x := 0; x < GridWidth; x++ {
				g.Grid[x][0] = 0
			}
			y++
		}
	}

	if linesCleared > 0 {
		g.LinesCleared += linesCleared
		switch linesCleared {
		case 1:
			g.Score += 100 * g.Level
		case 2:
			g.Score += 300 * g.Level
		case 3:
			g.Score += 500 * g.Level
		case 4:
			g.Score += 800 * g.Level
		}

		g.Level = g.LinesCleared/10 + 1
		g.DropSpeed = max(100, 1000-(g.Level-1)*100)
		g.DropInterval = time.Duration(g.DropSpeed) * time.Millisecond
	}
}

func (g *Game) newPiece() {
	g.CurrentPiece = g.NextPiece
	g.NextPiece = NewPiece(rand.Intn(7))
	g.RotationState = 0

	if !g.validMove(g.CurrentPiece, 0, 0, nil) {
		g.GameOver = true
		g.ShakeTimer = 1.0
	}
}

func (g *Game) Update() error {
	if g.GameOver {
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.Init()
		}
		return nil
	}

	dt := 1.0 / 60.0
	if g.FlashTimer > 0 {
		g.FlashTimer -= dt
		if g.FlashTimer <= 0 {
			g.removeLines()
		}
	}

	if g.ShakeTimer > 0 {
		g.ShakeTimer -= dt
		g.ShakeOffset = math.Sin(g.ShakeTimer*50) * 5 * g.ShakeTimer
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		if g.validMove(g.CurrentPiece, -1, 0, nil) {
			g.CurrentPiece.X--
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		if g.validMove(g.CurrentPiece, 1, 0, nil) {
			g.CurrentPiece.X++
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		g.LastMoveDown = time.Now()
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) && time.Since(g.LastMoveDown) > g.MoveDownDelay {
		if g.validMove(g.CurrentPiece, 0, 1, nil) {
			g.CurrentPiece.Y++
		}
		g.LastMoveDown = time.Now()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		g.rotatePiece()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		for g.validMove(g.CurrentPiece, 0, 1, nil) {
			g.CurrentPiece.Y++
		}
		g.mergePiece()
		g.clearLines()
		g.newPiece()
		g.LastDropTime = time.Now()
	}

	if time.Since(g.LastDropTime) > g.DropInterval {
		if g.validMove(g.CurrentPiece, 0, 1, nil) {
			g.CurrentPiece.Y++
		} else {
			g.mergePiece()
			g.clearLines()
			g.newPiece()
		}
		g.LastDropTime = time.Now()
	}

	return nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
