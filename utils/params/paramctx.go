package params

import (
	"context"
	"time"

	"github.com/jodi-ivan/numbered-notation-xml/utils/config"
)

const PARAM_CTX_KEY = "param"

type Param struct {
	Gregorian       bool
	Verse           int
	SingleVerseMode bool

	Diagnostic   *DiagParam
	Render       *RenderParam
	DirectReport *DirectReport
}

type paramCtx struct {
	ctx   context.Context
	Param *Param
}

func (pc *paramCtx) Deadline() (deadline time.Time, ok bool) {
	return pc.ctx.Deadline()
}
func (pc *paramCtx) Done() <-chan struct{} {
	return pc.ctx.Done()
}
func (pc *paramCtx) Err() error {
	return pc.ctx.Err()
}
func (pc *paramCtx) Value(key any) any {
	strKey, ok := key.(string)
	if ok && strKey == PARAM_CTX_KEY {
		return pc.Param
	}

	return pc.ctx.Value(key)
}

func NewParamContext(ctx context.Context, param *Param) context.Context {
	ctx = context.WithValue(ctx, PARAM_CTX_KEY, param)
	return &paramCtx{
		ctx:   ctx,
		Param: param,
	}
}
func GetParamFromContext(ctx context.Context) (*Param, bool) {
	pCtx, ok := ctx.(*paramCtx)
	if ok && pCtx != nil {
		return pCtx.Param, true
	}

	param, ok := ctx.Value(PARAM_CTX_KEY).(*Param)
	if !ok {
		return &Param{}, false
	}

	return param, true

}
func GetGregorianToggle(ctx context.Context) bool {
	p, ok := GetParamFromContext(ctx)
	/*
		|  ok  |   greogrian  |  want  |
		+------+--------------+--------+
		|   F  |   ? (any)    |   T    |
		|   T  |      T       |   T    |
		|   T  |      F       |   F    |

	*/
	return !ok || p.Gregorian
}

func GetDeviceInformation(ctx context.Context) *config.DeviceDimensions {
	p, ok := GetParamFromContext(ctx)
	if !ok || p.Render == nil {
		return nil
	}
	dims, found := config.DeviceTypeMap[p.Render.DeviceType]
	if !found {
		return nil
	}
	return &dims
}

// GetDeviceDimensions retrieves the device dimensions based on context and provided parameters.
func GetDeviceDimensions(ctx context.Context) (int, bool) {
	p, ok := GetParamFromContext(ctx)
	if !ok || p.Render == nil {
		return -1, false
	}

	dims, found := config.DeviceTypeMap[p.Render.DeviceType]
	if !found {
		return -1, false
	}

	if p.Render.Orientation == OrientationPortrait {
		return dims.PortraitW, true
	}
	return dims.LandscapeW, true
}

// SetMaxNoteWith sets the maximum note width for the RenderParam.
func SetMaxNoteWith(ctx context.Context, value int) {
	p, ok := GetParamFromContext(ctx)
	if !ok || p.Render == nil {
		return
	}

	if value > p.Render.MaxNoteWidth {
		p.Render.MaxNoteWidth = value
	}
}

// GetMaxPositionX retrieves the maximum note position X from the context.
func GetMaxPositionX(ctx context.Context) int {
	p, ok := GetParamFromContext(ctx)
	if !ok || p.Render == nil {
		return 0
	}
	return p.Render.MaxPositionX
}

// SetMaxPositionX sets the maximum note position X for the RenderParam in the context.
func SetMaxPositionX(ctx context.Context, value int) {
	p, ok := GetParamFromContext(ctx)
	if !ok || p.Render == nil {
		return
	}
	p.Render.MaxPositionX = value
}

func IgnoreNewLine(ctx context.Context) bool {
	p, ok := GetParamFromContext(ctx)
	if !ok || p.Render == nil {
		return false
	}

	return true
}
