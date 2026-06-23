package adapter

import (
	"bytes"
	"context"

	"github.com/jodi-ivan/numbered-notation-xml/svc/usecase"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

type CanvasDelegator struct{}

func (rscd *CanvasDelegator) OnBeforeStartWrite() {
	// no pre operation needed for string operation
}
func (rscd *CanvasDelegator) OnError(err error) canvas.DelegatorErrorFlowControl {
	return canvas.DelegatorErrorFlowControlStop
}

type RenderString struct {
	usecase usecase.Usecase
}

func NewRenderString(u usecase.Usecase) *RenderString {

	return &RenderString{
		usecase: u,
	}
}

func (rs *RenderString) RenderHymn(ctx context.Context, buf *bytes.Buffer, number int, variant ...string) (string, error) {
	canv := canvas.NewBufferedCanvas(buf, &CanvasDelegator{})
	err := rs.usecase.RenderHymn(ctx, canv, number, variant...)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
