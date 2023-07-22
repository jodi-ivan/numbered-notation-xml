package canvas

import (
	"io"

	svg "github.com/ajstarks/svgo"
)

type Canvas interface {
	Start(w int, h int, ns ...string)
	End()

	Def()
	DefEnd()

	Group(s ...string)
	Gend()

	Circle(x int, y int, r int, s ...string)
	Line(x1 int, y1 int, x2 int, y2 int, s ...string)
	Path(d string, s ...string)
	Qbez(sx int, sy int, cx int, cy int, ex int, ey int, s ...string)
	Qbezier(sx int, sy int, cx int, cy int, ex int, ey int, tx int, ty int, s ...string)
	Text(x int, y int, t string, s ...string)
	Writer() io.Writer
}

type _canvas struct {
	s *svg.SVG
}

func (c *_canvas) Start(w int, h int, ns ...string) {
	c.s.Start(w, h, ns...)
}
func (c *_canvas) End() {
	c.s.End()
}
func (c *_canvas) Def() {
	c.s.Def()
}
func (c *_canvas) DefEnd() {
	c.s.DefEnd()
}
func (c *_canvas) Group(s ...string) {
	c.s.Group(s...)
}
func (c *_canvas) Gend() {
	c.s.Gend()
}
func (c *_canvas) Circle(x int, y int, r int, s ...string) {
	c.s.Circle(x, y, r, s...)
}
func (c *_canvas) Line(x1 int, y1 int, x2 int, y2 int, s ...string) {
	c.s.Line(x1, y1, x2, y2, s...)
}
func (c *_canvas) Path(d string, s ...string) {
	c.s.Path(d, s...)
}
func (c *_canvas) Qbez(sx int, sy int, cx int, cy int, ex int, ey int, s ...string) {
	c.s.Qbez(sx, sy, cx, cy, ex, ey, s...)
}
func (c *_canvas) Qbezier(sx int, sy int, cx int, cy int, ex int, ey int, tx int, ty int, s ...string) {
	c.s.Qbezier(sx, sy, cx, cy, ex, ey, tx, ty, s...)
}
func (c *_canvas) Text(x int, y int, t string, s ...string) {
	c.s.Text(x, y, t, s...)
}

func (c *_canvas) Writer() io.Writer {
	return c.s.Writer
}

func NewCanvas(s *svg.SVG) Canvas {

	return &_canvas{
		s: s,
	}
}
