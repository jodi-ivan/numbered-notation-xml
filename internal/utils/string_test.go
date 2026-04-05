package utils_test

import (
	"testing"

	"github.com/jodi-ivan/numbered-notation-xml/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestReplaceItalicToSpan(t *testing.T) {

	got := utils.ReplaceItalicToSpan("<i>unittest</i>")
	assert.Equal(t, "<tspan font-style=\"italic\">unittest</tspan>", got, "TestReplaceItalicToSpan()")

}

func TestReplaceItalicToSpanWithClean(t *testing.T) {

	got, got2 := utils.ReplaceItalicToSpanWithClean("<i>unittest</i>")
	assert.Equal(t, "<tspan font-style=\"italic\">unittest</tspan>", got, "TestReplaceItalicToSpan() --> tspan")
	assert.Equal(t, "unittest", got2, "TestReplaceItalicToSpan() --> clean")

}

func TestCleanSpan(t *testing.T) {
	got := utils.CleanSpan("<tspan font-style=\"italic\">unittest</tspan>")
	assert.Equal(t, "unittest", got, "TestCleanSpan()")

}
