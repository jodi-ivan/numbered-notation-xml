package usecase

import (
	"context"
	"fmt"

	"github.com/jodi-ivan/numbered-notation-xml/internal/renderer"
	"github.com/jodi-ivan/numbered-notation-xml/svc/repository"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
	"github.com/jodi-ivan/numbered-notation-xml/utils/config"
)

type Usecase interface {
	RenderHymn(ctx context.Context, canv canvas.Canvas, hymnNum int, variant ...string) error
}

type interactor struct {
	config   config.Config
	repo     repository.Repository
	renderer renderer.Renderer
}

func New(config config.Config, repo repository.Repository, renderer renderer.Renderer) Usecase {
	return &interactor{
		config:   config,
		repo:     repo,
		renderer: renderer,
	}
}

func (i *interactor) RenderHymn(ctx context.Context, canv canvas.Canvas, hymnNum int, variant ...string) error {
	filepath := fmt.Sprintf("%s%s-%03d.musicxml", i.config.MusicXML.Path, i.config.MusicXML.FilePrefix, hymnNum)
	if len(variant) > 0 {
		filepath = fmt.Sprintf("%s%s-%03d%s.musicxml", i.config.MusicXML.Path, i.config.MusicXML.FilePrefix, hymnNum, variant[0])
	}
	music, err := i.repo.GetMusicXML(ctx, filepath)
	if err != nil {
		flow := canv.Delegator().OnError(err)
		if flow != canvas.DelegatorErrorFlowControlIgnore {
			return err
		}
	}
	metaData, err := i.repo.GetHymnMetaData(ctx, hymnNum, variant...)
	if err != nil {
		flow := canv.Delegator().OnError(err)
		if flow != canvas.DelegatorErrorFlowControlIgnore {
			return err
		}
	}

	canv.Delegator().OnBeforeStartWrite()

	i.renderer.Render(ctx, music, canv, metaData)

	return nil
}
