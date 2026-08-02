// Package anim re-exports animation types from the internal/anim library.
// Existing code may continue to use this package path during migration.
package anim

import (
	"github.com/rfaiii/donk-cli-main/internal/anim"
)

// Spring is an alias for anim.Spring.
type Spring = anim.Spring

// NewSpring is an alias for anim.NewSpring.
var NewSpring = anim.NewSpring

// Progress is an alias for anim.Progress.
type Progress = anim.Progress

// NewProgress is an alias for anim.NewProgress.
var NewProgress = anim.NewProgress

// Spinner is an alias for anim.Spinner.
type Spinner = anim.Spinner

// SpinnerKind is an alias for anim.SpinnerKind.
type SpinnerKind = anim.SpinnerKind

// SpinnerKind* constants are aliases for anim.SpinnerKind* constants.
const (
	SpinnerKindLine     = anim.SpinnerKindLine
	SpinnerKindDot      = anim.SpinnerKindDot
	SpinnerKindMiniDot  = anim.SpinnerKindMiniDot
	SpinnerKindJump     = anim.SpinnerKindJump
	SpinnerKindPulse    = anim.SpinnerKindPulse
	SpinnerKindPoints   = anim.SpinnerKindPoints
	SpinnerKindMoon     = anim.SpinnerKindMoon
)

// NewSpinner is an alias for anim.NewSpinner.
var NewSpinner = anim.NewSpinner

// GalleryName is an alias for anim.GalleryName.
type GalleryName = anim.GalleryName

// GalleryNames is an alias for anim.GalleryNames.
var GalleryNames = anim.GalleryNames

// GalleryTabLine is an alias for anim.GalleryTabLine.
var GalleryTabLine = anim.GalleryTabLine
