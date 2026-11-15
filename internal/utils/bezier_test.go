package utils

import (
	"testing"

	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestQuadBezierLength(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		p0    entity.Coordinate
		p1    entity.Coordinate
		p2    entity.Coordinate
		steps int
		want  float64
	}{
		{
			name:  "Straight line segment",
			p0:    entity.Coordinate{X: 0, Y: 0},
			p1:    entity.Coordinate{X: 5, Y: 0},
			p2:    entity.Coordinate{X: 10, Y: 0},
			steps: 10,
			want:  10.0,
		},
		{
			name:  "Simple parabolic arc (2 steps)",
			p0:    entity.Coordinate{X: 0, Y: 0},
			p1:    entity.Coordinate{X: 5, Y: 10},
			p2:    entity.Coordinate{X: 10, Y: 0},
			steps: 2,
			want:  14.142135623730951, // sqrt(50) + sqrt(50)
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := QuadBezierLength(tt.p0, tt.p1, tt.p2, tt.steps)
			assert.InDelta(t, tt.want, got, 1e-9)
		})
	}
}
