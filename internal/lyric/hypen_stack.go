package lyric

import (
	"context"

	"github.com/golang-collections/collections/stack"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
)

type HypenStack struct {
	container *stack.Stack
}

func NewHypenStack() *HypenStack {
	return &HypenStack{
		container: stack.New(),
	}
}

func (hs *HypenStack) Process(ctx context.Context, syllabic musicxml.LyricSyllabic) {
	switch syllabic {
	case musicxml.LyricSyllabicTypeEnd:
		hs.container.Pop()
	case musicxml.LyricSyllabicTypeBegin:
		hs.container.Push(syllabic)
	case musicxml.LyricSyllabicTypeSingle:
		lastRaw := hs.container.Peek()

		last, _ := lastRaw.(musicxml.LyricSyllabic)
		if last == musicxml.LyricSyllabicTypeBegin || last == musicxml.LyricSyllabicTypeMiddle {
			hs.container.Pop()
		}

	}
}

func (hs *HypenStack) IsEmpty() bool {
	return hs.container.Len() == 0
}
