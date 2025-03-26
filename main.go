package main

import (
	"image/color"
	"math/rand"
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/basicfont"
)

const (
	screenWidth  = 300
	screenHeight = 600
	blockSize    = 30
	gridWidth    = 10
	gridHeight   = 20
)

var (
	colors = []color.RGBA{
		{0, 0, 0, 255},     // Empty (black)
		{0, 255, 255, 255}, // I (cyan)
		{0, 0, 255, 255},   // J (blue)
		{255, 165, 0, 255}, // L (orange)
		{255, 255, 0, 255}, // O (yellow)
		{0, 255, 0, 255},   // S (green)
		{128, 0, 128, 255}, // T (purple)
		{255, 0, 0, 255},   // Z (red)
	}
)

type Game struct {
	grid          [gridWidth][gridHeight]int
	currentPiece  *Piece
	nextPiece     *Piece
	gameOver      bool
	score         int
	level         int
	linesCleared  int
	dropSpeed     int
	lastDropTime  time.Time
	dropInterval  time.Duration
	moveDownDelay time.Duration
	lastMoveDown  time.Time
}

type Piece struct {
	shape    [][]int
	x, y     int
	colorIdx int
}

var pieceShapes = [][][]int{
	// I
	{
		{0, 0, 0, 0},
		{1, 1, 1, 1},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
	},
	// J
	{
		{2, 0, 0},
		{2, 2, 2},
		{0, 0, 0},
	},
	// L
	{
		{0, 0, 3},
		{3, 3, 3},
		{0, 0, 0},
	},
	// O
	{
		{4, 4},
		{4, 4},
	},
	// S
	{
		{0, 5, 5},
		{5, 5, 0},
		{0, 0, 0},
	},
	// T
	{
		{0, 6, 0},
		{6, 6, 6},
		{0, 0, 0},
	},
	// Z
	{
		{7, 7, 0},
		{0, 7, 7},
		{0, 0, 0},
	},
}

func NewPiece(shapeIdx int) *Piece {
	shape := pieceShapes[shapeIdx]
	return &Piece{
		shape:    shape,
		x:        gridWidth/2 - len(shape[0])/2,
		y:        0,
		colorIdx: shapeIdx + 1,
	}
}

func (g *Game) Init() {
	rand.Seed(time.Now().UnixNano())
	g.currentPiece = NewPiece(rand.Intn(len(pieceShapes)))
	g.nextPiece = NewPiece(rand.Intn(len(pieceShapes)))
	g.gameOver = false
	g.score = 0
	g.level = 1
	g.linesCleared = 0
	g.dropSpeed = 1000
	g.dropInterval = time.Duration(g.dropSpeed) * time.Millisecond
	g.lastDropTime = time.Now()
	g.moveDownDelay = 100 * time.Millisecond
	g.lastMoveDown = time.Now()
}

func (g *Game) validMove(p *Piece, x, y int, rotatedShape [][]int) bool {
	if rotatedShape == nil {
		rotatedShape = p.shape
	}

	for py, row := range rotatedShape {
		for px, val := range row {
			if val != 0 {
				newX, newY := p.x+x+px, p.y+y+py
				if newX < 0 || newX >= gridWidth || newY >= gridHeight {
					return false
				}
				if newY >= 0 && g.grid[newX][newY] != 0 {
					return false
				}
			}
		}
	}
	return true
}

func (g *Game) rotatePiece() {
	if g.currentPiece == nil {
		return
	}

	rows := len(g.currentPiece.shape)
	cols := len(g.currentPiece.shape[0])
	rotated := make([][]int, cols)
	for i := range rotated {
		rotated[i] = make([]int, rows)
	}

	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			rotated[x][rows-1-y] = g.currentPiece.shape[y][x]
		}
	}

	if g.validMove(g.currentPiece, 0, 0, rotated) {
		g.currentPiece.shape = rotated
	}
}

func (g *Game) mergePiece() {
	for y, row := range g.currentPiece.shape {
		for x, val := range row {
			if val != 0 {
				gridX, gridY := g.currentPiece.x+x, g.currentPiece.y+y
				if gridY >= 0 {
					g.grid[gridX][gridY] = g.currentPiece.colorIdx
				}
			}
		}
	}
}

func (g *Game) clearLines() {
	linesCleared := 0
	for y := gridHeight - 1; y >= 0; y-- {
		lineFull := true
		for x := 0; x < gridWidth; x++ {
			if g.grid[x][y] == 0 {
				lineFull = false
				break
			}
		}

		if lineFull {
			linesCleared++
			// Move all lines above down
			for ny := y; ny > 0; ny-- {
				for x := 0; x < gridWidth; x++ {
					g.grid[x][ny] = g.grid[x][ny-1]
				}
			}
			// Clear top line
			for x := 0; x < gridWidth; x++ {
				g.grid[x][0] = 0
			}
			y++ // Check the same line again
		}
	}

	if linesCleared > 0 {
		g.linesCleared += linesCleared
		switch linesCleared {
		case 1:
			g.score += 100 * g.level
		case 2:
			g.score += 300 * g.level
		case 3:
			g.score += 500 * g.level
		case 4:
			g.score += 800 * g.level
		}

		// Level up every 10 lines
		g.level = g.linesCleared/10 + 1
		g.dropSpeed = 1000 - (g.level-1)*100
		if g.dropSpeed < 100 {
			g.dropSpeed = 100
		}
		g.dropInterval = time.Duration(g.dropSpeed) * time.Millisecond
	}
}

func (g *Game) newPiece() {
	g.currentPiece = g.nextPiece
	g.nextPiece = NewPiece(rand.Intn(len(pieceShapes)))

	// Check if game over
	if !g.validMove(g.currentPiece, 0, 0, nil) {
		g.gameOver = true
	}
}

func (g *Game) Update() error {
	if g.gameOver {
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.Init()
		}
		return nil
	}

	// Handle input
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		if g.validMove(g.currentPiece, -1, 0, nil) {
			g.currentPiece.x--
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		if g.validMove(g.currentPiece, 1, 0, nil) {
			g.currentPiece.x++
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		g.lastMoveDown = time.Now()
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) && time.Since(g.lastMoveDown) > g.moveDownDelay {
		if g.validMove(g.currentPiece, 0, 1, nil) {
			g.currentPiece.y++
		}
		g.lastMoveDown = time.Now()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		g.rotatePiece()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		// Hard drop
		for g.validMove(g.currentPiece, 0, 1, nil) {
			g.currentPiece.y++
		}
		g.mergePiece()
		g.clearLines()
		g.newPiece()
		g.lastDropTime = time.Now()
	}

	// Automatic dropping
	if time.Since(g.lastDropTime) > g.dropInterval {
		if g.validMove(g.currentPiece, 0, 1, nil) {
			g.currentPiece.y++
		} else {
			g.mergePiece()
			g.clearLines()
			g.newPiece()
		}
		g.lastDropTime = time.Now()
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw grid background
	for x := 0; x < gridWidth; x++ {
		for y := 0; y < gridHeight; y++ {
			clr := colors[g.grid[x][y]]
			drawBlock(screen, x, y, clr)
		}
	}

	// Draw current piece
	if g.currentPiece != nil {
		for y, row := range g.currentPiece.shape {
			for x, val := range row {
				if val != 0 {
					px, py := g.currentPiece.x+x, g.currentPiece.y+y
					if py >= 0 {
						drawBlock(screen, px, py, colors[g.currentPiece.colorIdx])
					}
				}
			}
		}
	}

	// Draw next piece preview
	nextX, nextY := gridWidth*blockSize+20, 60
	text.Draw(screen, "Next:", basicfont.Face7x13, nextX, nextY-10, color.White)
	if g.nextPiece != nil {
		for y, row := range g.nextPiece.shape {
			for x, val := range row {
				if val != 0 {
					px, py := nextX+x*blockSize, nextY+y*blockSize
					drawBlockAt(screen, px, py, colors[g.nextPiece.colorIdx])
				}
			}
		}
	}

	// Draw score and level
	scoreText := "Score: " + strconv.Itoa(g.score)
	levelText := "Level: " + strconv.Itoa(g.level)
	linesText := "Lines: " + strconv.Itoa(g.linesCleared)
	text.Draw(screen, scoreText, basicfont.Face7x13, gridWidth*blockSize+20, 180, color.White)
	text.Draw(screen, levelText, basicfont.Face7x13, gridWidth*blockSize+20, 200, color.White)
	text.Draw(screen, linesText, basicfont.Face7x13, gridWidth*blockSize+20, 220, color.White)

	// Draw controls help
	controls := []string{
		"Controls:",
		"Left/Right: Move",
		"Up: Rotate",
		"Down: Soft Drop",
		"Space: Hard Drop",
	}
	for i, ctrl := range controls {
		text.Draw(screen, ctrl, basicfont.Face7x13, gridWidth*blockSize+20, 260+i*20, color.White)
	}

	// Draw game over screen
	if g.gameOver {
		vector.DrawFilledRect(screen, 0, 0, float32(screenWidth), float32(screenHeight), color.RGBA{0, 0, 0, 200}, false)
		text.Draw(screen, "GAME OVER", basicfont.Face7x13, screenWidth/2-40, screenHeight/2-20, color.White)
		text.Draw(screen, "Press SPACE to restart", basicfont.Face7x13, screenWidth/2-80, screenHeight/2+10, color.White)
	}
}

func drawBlock(screen *ebiten.Image, x, y int, clr color.Color) {
	vector.DrawFilledRect(screen, float32(x*blockSize), float32(y*blockSize), blockSize, blockSize, clr, false)
	vector.StrokeRect(screen, float32(x*blockSize), float32(y*blockSize), blockSize, blockSize, 1, color.RGBA{50, 50, 50, 255}, false)
}

func drawBlockAt(screen *ebiten.Image, x, y int, clr color.Color) {
	vector.DrawFilledRect(screen, float32(x), float32(y), blockSize, blockSize, clr, false)
	vector.StrokeRect(screen, float32(x), float32(y), blockSize, blockSize, 1, color.RGBA{50, 50, 50, 255}, false)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	game := &Game{}
	game.Init()

	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Tetris in Go")
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
