package params

import (
	"context"
	"time"
)

const PARAM_CTX_KEY = "param"

type Param struct {
	DisableGregorian bool
	Verse            int
	SingleVerseMode  bool
}

type paramCtx struct {
	ctx   context.Context
	Param Param
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

func NewParamContext(ctx context.Context, param Param) context.Context {
	ctx = context.WithValue(ctx, PARAM_CTX_KEY, param)
	return &paramCtx{
		ctx:   ctx,
		Param: param,
	}
}
func GetParamFromContext(ctx context.Context) (Param, bool) {
	pCtx, ok := ctx.(*paramCtx)
	if ok && pCtx != nil {
		return pCtx.Param, true
	}

	param, ok := ctx.Value(PARAM_CTX_KEY).(*Param)
	if !ok {
		return Param{}, false
	}

	return *param, true

}
