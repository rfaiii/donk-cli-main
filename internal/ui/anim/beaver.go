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
// Pass errored=true to render the x-ray (X_X) Beta variant when the agent has
// errored; otherwise the normal (0_0) Alpha variant is shown. Both reuse the
// boot splash mascot so the homescreen beaver matches the intro, and both are
// driven by the existing banner ticker (m.bannerFrame from the UI model) so
// they idle without their own tick loop.
func BeaverFrame(tick int, errored bool) string {
	frames := boot.BeaverFramesDenseAlpha
	if errored {
		frames = boot.BeaverFramesDenseBeta
	}
	seq := boot.BootSequence
	if len(seq) == 0 {
		return ""
	}
	if frame, ok := frames[seq[tick%len(seq)]]; ok {
		return mascotStyle.Render(frame)
	}
	return ""
}
