package game

type Piece struct {
	Shape    [][]int
	X, Y     int
	ColorIdx int
}

var pieceShapes = [][][]int{
	// I (2 rotation states)
	{
		{0, 0, 0, 0},
		{1, 1, 1, 1},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
	},
	{
		{0, 0, 1, 0},
		{0, 0, 1, 0},
		{0, 0, 1, 0},
		{0, 0, 1, 0},
	},
	// J (4 rotation states)
	{
		{2, 0, 0},
		{2, 2, 2},
		{0, 0, 0},
	},
	{
		{0, 2, 2},
		{0, 2, 0},
		{0, 2, 0},
	},
	{
		{0, 0, 0},
		{2, 2, 2},
		{0, 0, 2},
	},
	{
		{0, 2, 0},
		{0, 2, 0},
		{2, 2, 0},
	},
	// L (4 rotation states)
	{
		{0, 0, 3},
		{3, 3, 3},
		{0, 0, 0},
	},
	{
		{0, 3, 0},
		{0, 3, 0},
		{0, 3, 3},
	},
	{
		{0, 0, 0},
		{3, 3, 3},
		{3, 0, 0},
	},
	{
		{3, 3, 0},
		{0, 3, 0},
		{0, 3, 0},
	},
	// O (1 rotation state)
	{
		{4, 4},
		{4, 4},
	},
	// S (2 rotation states)
	{
		{0, 5, 5},
		{5, 5, 0},
		{0, 0, 0},
	},
	{
		{0, 5, 0},
		{0, 5, 5},
		{0, 0, 5},
	},
	// T (4 rotation states)
	{
		{0, 6, 0},
		{6, 6, 6},
		{0, 0, 0},
	},
	{
		{0, 6, 0},
		{0, 6, 6},
		{0, 6, 0},
	},
	{
		{0, 0, 0},
		{6, 6, 6},
		{0, 6, 0},
	},
	{
		{0, 6, 0},
		{6, 6, 0},
		{0, 6, 0},
	},
	// Z (2 rotation states)
	{
		{7, 7, 0},
		{0, 7, 7},
		{0, 0, 0},
	},
	{
		{0, 0, 7},
		{0, 7, 7},
		{0, 7, 0},
	},
}

func NewPiece(shapeIdx int) *Piece {
	baseShapeIdx := 0
	switch shapeIdx {
	case 0:
		baseShapeIdx = 0
	case 1:
		baseShapeIdx = 2
	case 2:
		baseShapeIdx = 6
	case 3:
		baseShapeIdx = 10
	case 4:
		baseShapeIdx = 11
	case 5:
		baseShapeIdx = 13
	case 6:
		baseShapeIdx = 17
	}

	shape := pieceShapes[baseShapeIdx]
	return &Piece{
		Shape:    shape,
		X:        GridWidth/2 - len(shape[0])/2,
		Y:        0,
		ColorIdx: shapeIdx + 1,
	}
}
