package decorator

import (
	"bytes"
	"context"

	"github.com/jodi-ivan/numbered-notation-xml/adapter"
	"github.com/jodi-ivan/numbered-notation-xml/svc/repository"
	"github.com/jodi-ivan/numbered-notation-xml/utils/params"
)

type HymnVarianDirectReport struct {
	Repo repository.Repository
	Next *adapter.RenderString

	TotalVariant int
	TotalVerse   int
}

func (hvd *HymnVarianDirectReport) RenderHymn(ctx context.Context, buf *bytes.Buffer, number int, variant ...string) (string, error) {
	param, _ := params.GetParamFromContext(ctx)
	param.DirectReport = &params.DirectReport{
		TotalVerse: make(chan int, 1),
	}

	rctx := params.NewParamContext(ctx, param)
	variants, err := hvd.Repo.GetHymnVariant(rctx, number)
	if err != nil {
		return "", err
	}

	if len(variants) > 0 {
		hvd.TotalVariant = len(variants)
		if len(variant) == 0 {
			variant = []string{"a"}
		}
	}

	go func() {
		hvd.TotalVerse = <-param.DirectReport.TotalVerse
	}()
	return hvd.Next.RenderHymn(rctx, buf, number, variant...)
}
