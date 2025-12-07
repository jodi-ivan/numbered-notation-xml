package decorator

import (
	"context"

	"github.com/jodi-ivan/numbered-notation-xml/svc/repository"
	"github.com/jodi-ivan/numbered-notation-xml/svc/usecase"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

type HTTPMiddlewareAdapter func(usecase.Usecase) usecase.Usecase

type HymnVariantDecorator struct {
	repo repository.Repository
	next usecase.Usecase
}

func (hvd *HymnVariantDecorator) RenderHymn(ctx context.Context, canv canvas.Canvas, hymnNum int, variant ...string) error {
	if len(variant) == 0 {
		variants, err := hvd.repo.GetHymnVariant(ctx, hymnNum)
		if err != nil {
			return err
		}

		if len(variants) > 0 {
			return repository.ErrHymnHasMoreThanOneVariant
		}
	}
	return hvd.next.RenderHymn(ctx, canv, hymnNum, variant...)
}

func WithVariantRedirect(repo repository.Repository) HTTPMiddlewareAdapter {
	return func(next usecase.Usecase) usecase.Usecase {
		return &HymnVariantDecorator{repo: repo, next: next}
	}
}
