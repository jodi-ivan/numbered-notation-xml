package utils

import (
	"math"

	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
)

func QuadBezierLength(p0, p1, p2 entity.Coordinate, steps int) float64 {
	var length float64
	prev := p0

	for i := 1; i <= steps; i++ {
		t := float64(i) / float64(steps)

		x := math.Pow(1-t, 2)*p0.X +
			2*(1-t)*t*p1.X +
			math.Pow(t, 2)*p2.X

		y := math.Pow(1-t, 2)*p0.Y +
			2*(1-t)*t*p1.Y +
			math.Pow(t, 2)*p2.Y

		dx := x - prev.X
		dy := y - prev.Y

		length += math.Hypot(dx, dy)
		prev = entity.Coordinate{X: x, Y: y}
	}

	return length
}
