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

// BeaverFrame returns the homescreen beaver mascot, facing toward the
// cursor/prompt instead of cycling through the boot sequence on its own.
//
// It looks left when the cursor/prompt sits on the left half of the terminal
// and right when it sits on the right half, so the mascot tracks the user
// (a slow, cursor-following idle) rather than spazzing out on the banner
// ticker. resting renders the neutral center-facing "rest" pose the mascot
// briefly strikes between direction changes. Pass errored=true to render the
// x-ray Beta variant (X_X eyes) when the agent has errored; otherwise the
// normal Alpha variant is shown.
//
// cursorX is the terminal cell column of the cursor/mouse and width is the
// terminal width in cells. When cursorX hasn't been observed yet (<=0) the
// mascot defaults to facing right, toward the typing area.
func BeaverFrame(cursorX, width int, errored, resting bool) string {
	frames := boot.BeaverFramesDenseAlpha
	if errored {
		frames = boot.BeaverFramesDenseBeta
	}

	if resting {
		return mascotStyle.Render(frames["center"])
	}

	// Face the cursor/prompt: left half of the screen -> look left,
	// right half -> look right. Default (no cursor yet) -> right.
	facing := "right"
	if cursorX > 0 && cursorX < width/2 {
		facing = "left"
	} else if cursorX >= width/2 {
		facing = "right"
	}
	frame, ok := frames[facing]
	if !ok {
		frame = frames["center"]
	}
	return mascotStyle.Render(frame)
}
