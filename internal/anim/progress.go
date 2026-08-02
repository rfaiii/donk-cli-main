// Package anim provides a spring-eased progress bar ported from
// the Rust donk-anim crate.
package anim

import (
	"math"
	"strconv"
	"strings"
)

// Progress is a spring-physics progress bar.
type Progress struct {
	Spring *Spring
	Width  int
	Label  string
}

// NewProgress creates a progress bar with default parameters.
func NewProgress() *Progress {
	return &Progress{
		Spring: NewSpring(),
		Width:  40,
		Label:  "progress",
	}
}

// SetWidth sets the minimum bar width.
func (p *Progress) SetWidth(width int) {
	if width > p.Width {
		p.Width = width
	}
	p.Spring.setWidth(p.Width)
}

// SetPercent sets the spring target as a 0-1 fraction.
func (p *Progress) SetPercent(pct float64) {
	if pct < 0 {
		pct = 0
	}
	if pct > 1 {
		pct = 1
	}
	p.Spring.Target = pct
}

// IncrPercent increments the spring target by delta.
func (p *Progress) IncrPercent(delta float64) {
	p.SetPercent(p.Spring.Target + delta)
}

// Poke toggles progress and injects spring velocity.
func (p *Progress) Poke() {
	p.Spring.Poke()
	if p.Spring.Target >= 1-math.Nextafter(1, 2) {
		p.Spring.Target = 0
		p.Spring.Pos = 0
		p.Spring.Vel = 0
	}
}

// Update advances the spring and frame counter.
func (p *Progress) Update() {
	p.Spring.Update()
}

// Percent returns the current spring position clamped to 0-1.
func (p *Progress) Percent() float64 {
	if p.Spring.Pos < 0 {
		return 0
	}
	if p.Spring.Pos > 1 {
		return 1
	}
	return p.Spring.Pos
}

// Render returns a text progress bar with label.
func (p *Progress) Render() string {
	w := p.Width
	if w > 48 {
		w = 48
	}
	fill := int(math.Round(p.Percent() * float64(w)))
	if fill > w {
		fill = w
	}
	bar := strings.Repeat("█", fill) + strings.Repeat("░", w-fill)
	pct := int(math.Round(p.Percent() * 100))
	return p.Label + "\n\n[" + bar + "] " + strconv.Itoa(pct) + "%\n\ntarget=" + strconv.FormatFloat(p.Spring.Target*100, 'f', 0, 64) + "% · harmonica spring · space +25%"
}
