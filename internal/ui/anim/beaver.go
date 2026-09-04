// Package anim provides small terminal animation helpers for BVR-CLI.
// In addition to the cursor trail, it hosts the idle beaver mascot that
// lives on the homescreen and reuses the boot splash frames.
package anim

import (
	"charm.land/lipgloss/v2"

	"github.com/richavery/bvr-cli/internal/ui/boot"
)

// mascotStyle renders the beaver mascot in BVR neon green (#39f66b) with bold
// weight, matching the boot splash treatment so the homescreen mascot looks
// identical to the intro.
var mascotStyle = lipgloss.NewStyle().Foreground(boot.MascotColor).Bold(true)

// BeaverFrame returns the BVR-CLI beaver mascot for the given tick index,
// cycling through the boot sequence (center, left, center, right, center).
//
// It reuses the boot splash mascot so the homescreen beaver is identical to
// the intro, and is driven by the existing banner ticker (pass m.bannerFrame
// from the UI model) so it idles on the homescreen without its own tick loop.
func BeaverFrame(tick int) string {
	seq := boot.BootSequence
	if len(seq) == 0 {
		return ""
	}
	if frame, ok := boot.BeaverFrames[seq[tick%len(seq)]]; ok {
		return mascotStyle.Render(frame)
	}
	return ""
}
