package entity

type Coordinate struct {
	X float64
	Y float64
}

func NewCoordinate(x, y float64) Coordinate {
	return Coordinate{
		X: x,
		Y: y,
	}
}
