// Package boot renders the animated DONK-CLI beaver intro sequence before the
// Bubble Tea main loop takes over the terminal.
package boot

import (
	"fmt"
	"time"

	"charm.land/lipgloss/v2"
)

// beaverFrames are the static-width ASCII mascot frames. Keeping the width
// constant across frames prevents the terminal grid from collapsing or
// jittering while the animation plays.
var beaverFrames = map[string]string{
	"center": `  (\.---./)  
 /   0_0   \ 
|   -(v)-   |
 \   [=]   / 
   '-------'  `,
	"left": `  (\.---./)  
 / 0_0     \ 
| -(v)-     |
 \ [=]     / 
   '-------'  `,
	"right": `  (\.---./)  
 /     0_0 \ 
|     -(v)- |
 \     [=] / 
   '-------'  `,
}

// bootSequence is the frame order for the intro animation.
var bootSequence = []string{"center", "left", "center", "right", "center"}

// mascotColor is the DONK neon green (#39f66b) used for the beaver so it pops
// vibrantly against the dark terminal background.
var mascotColor = lipgloss.Color("#39f66b")

// loadingColor is a secondary accent used for the loading line beneath the
// mascot.
var loadingColor = lipgloss.Color("#b972ff")

var mascotStyle = lipgloss.NewStyle().Foreground(mascotColor).Bold(true)
var loadingStyle = lipgloss.NewStyle().Foreground(loadingColor)

// Run plays the beaver boot animation. It clears the terminal first, renders
// each mascot frame in neon green with a secondary-accent loading line beneath,
// then returns so the caller can start the Bubble Tea main event loop.
//
// Run is a no-op guard: callers should only invoke it from an interactive
// terminal (see IsInteractive) and should clear any prior program output first.
func Run() {
	// Clear the screen and home the cursor so the mascot starts top-left.
	fmt.Print("\033[H\033[2J")
	for _, frame := range bootSequence {
		fmt.Print("\033[0;0H") // cursor to top-left
		fmt.Println(mascotStyle.Render(beaverFrames[frame]))
		fmt.Println()
		fmt.Println(loadingStyle.Render("Initializing donk-cli components..."))
		time.Sleep(400 * time.Millisecond)
	}
}
