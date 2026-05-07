package entity

type Coordinate struct {
	X float64
	Y float64
}

func (c *Coordinate) IsEmpty() bool {
	return c.X == 0 && c.Y == 0
}

func NewCoordinate(x, y float64) Coordinate {
	return Coordinate{
		X: x,
		Y: y,
	}
}
