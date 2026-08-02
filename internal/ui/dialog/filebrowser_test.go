package dialog

import (
	"testing"

	"github.com/charmbracelet/x/ansi"
	"github.com/stretchr/testify/require"
)

func TestFixedLinesKeepsViewportStable(t *testing.T) {
	lines := fixedLines([]string{
		"/Users/richavery/Projects/donk-cli-go/this-is-a-very-long-path-that-must-not-wrap",
		"second line",
		"third line should be clipped",
	}, 20, 2)

	require.Len(t, lines, 2)
	for _, line := range lines {
		require.Equal(t, 20, len([]rune(line)))
	}
	require.Contains(t, lines[0], "…")
	require.Equal(t, "second line         ", lines[1])
}

func TestFixedLinesHandlesEmptyViewport(t *testing.T) {
	require.Empty(t, fixedLines([]string{"content"}, 0, 3))
	require.Empty(t, fixedLines([]string{"content"}, 10, 0))
}

func TestFixedLinesClipsUnbrokenTextAndNormalizesCarriageReturns(t *testing.T) {
	lines := fixedLines([]string{"0123456789abcdef\rOVERWRITE"}, 8, 1)

	require.Len(t, lines, 1)
	require.Equal(t, 8, ansi.StringWidth(lines[0]))
	require.NotContains(t, lines[0], "\r")
	require.Equal(t, "0123456…", lines[0])
}

func TestFixedLinesDoesNotPreserveRtfControlSequencesAsLayoutCommands(t *testing.T) {
	lines := fixedLines([]string{"{\\rtf1\\ansi\\ansicpg1252\\cocoartf2868"}, 16, 1)

	require.Len(t, lines, 1)
	require.Equal(t, 16, ansi.StringWidth(lines[0]))
	require.NotContains(t, lines[0], "\n")
	require.NotContains(t, lines[0], "\r")
}

func TestFileBrowserKeepsSelectedEntryVisible(t *testing.T) {
	f := &FileBrowser{
		entries:  make([]fileBrowserEntry, 12),
		selected: 0,
	}

	f.selected = 8
	f.ensureSelectedVisible(4)
	require.Equal(t, 5, f.scroll)

	f.selected = 2
	f.ensureSelectedVisible(4)
	require.Equal(t, 2, f.scroll)

	f.selected = 11
	f.ensureSelectedVisible(4)
	require.Equal(t, 8, f.scroll)
}

func TestScrollbarLinesShowsMovingBlock(t *testing.T) {
	top := scrollbarLines(5, 20, 0)
	bottom := scrollbarLines(5, 20, 15)

	require.Len(t, top, 5)
	require.Len(t, bottom, 5)
	require.Equal(t, "█", top[0])
	require.Equal(t, "█", bottom[4])
	require.NotEqual(t, top, bottom)
}

func TestCloseLabelUsesUppercaseButton(t *testing.T) {
	require.Equal(t, "[X]", closeLabel(40))
	require.Equal(t, "X", closeLabel(8))
}

func TestFinderPaneWidthsLeavesRightGutter(t *testing.T) {
	left, right, list := finderPaneWidths(80)

	require.Equal(t, left-1, list)
	require.Equal(t, 77, left+1+right)
}

func TestFinderPaneWidthsHandlesNarrowContent(t *testing.T) {
	left, right, list := finderPaneWidths(2)

	require.GreaterOrEqual(t, left, 1)
	require.GreaterOrEqual(t, right, 1)
	require.GreaterOrEqual(t, list, 1)
}
