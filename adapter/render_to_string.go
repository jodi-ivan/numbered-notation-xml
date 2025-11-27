package adapter

import (
	"bytes"
	"context"

	svg "github.com/ajstarks/svgo"
	"github.com/jodi-ivan/numbered-notation-xml/svc/usecase"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

type RenderStringCanvasDelegator struct{}

func (rscd *RenderStringCanvasDelegator) OnBeforeStartWrite() {}
func (rscd *RenderStringCanvasDelegator) OnError(err error) canvas.DelegatorErrorFlowControl {
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

func (rs *RenderString) RenderHymn(ctx context.Context, buf *bytes.Buffer, number int) (string, error) {
	canv := canvas.NewCanvasWithDelegator(svg.New(buf), &RenderStringCanvasDelegator{})
	err := rs.usecase.RenderHymn(ctx, canv, number)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
