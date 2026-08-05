package common

import (
	"github.com/richavery/donk-cli/internal/ui/diffview"
	"github.com/richavery/donk-cli/internal/ui/styles"
)

// DiffFormatter returns a diff formatter with the given styles that can be
// used to format diff outputs.
func DiffFormatter(s *styles.Styles) *diffview.DiffView {
	formatDiff := diffview.New()
	diff := formatDiff.ChromaStyle(ChromaStyle(s, nil)).Style(s.Diff).TabWidth(4)
	return diff
}
