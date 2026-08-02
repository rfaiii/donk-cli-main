package anim

import (
	"math"
	"testing"
	"time"
)

func TestNewSpinnerDefaultsToDot(t *testing.T) {
	s := NewSpinner(SpinnerKind(""))
	if s == nil {
		t.Fatal("expected spinner")
	}
	if s.Kind() != SpinnerKindDot {
		t.Fatalf("expected dot, got %s", s.Kind())
	}
	if s.Frame() == "" {
		t.Fatal("expected non-empty frame")
	}
}

func TestSpinnerSetKindUpdatesFrame(t *testing.T) {
	s := NewSpinner(SpinnerKindDot)
	before := s.Frame()
	s.SetKind(SpinnerKindLine)
	if s.Kind() != SpinnerKindLine {
		t.Fatalf("expected line, got %s", s.Kind())
	}
	if s.Frame() == before {
		t.Fatal("expected frame to change after set kind")
	}
}

func TestSpinnerCycleKind(t *testing.T) {
	s := NewSpinner(SpinnerKindLine)
	s.CycleKind()
	if s.Kind() == SpinnerKindLine {
		t.Fatal("expected kind to change")
	}
}

func TestSpinnerUpdateAdvancesFrame(t *testing.T) {
	s := NewSpinner(SpinnerKindLine)
	before := s.Frame()
	s.Update(time.Now().Add(time.Second))
	if s.Frame() == before {
		t.Fatal("expected frame to advance after long interval")
	}
}

func TestSpinnerUpdateSkipsEarlyTick(t *testing.T) {
	s := NewSpinner(SpinnerKindLine)
	before := s.Frame()
	s.Update(time.Now())
	if s.Frame() != before {
		t.Fatal("expected frame to hold within interval")
	}
}

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) < 1e-9
}

func TestNewSpringDefaults(t *testing.T) {
	s := NewSpring()
	if !almostEqual(s.Target, 1) {
		t.Fatalf("expected target 1, got %v", s.Target)
	}
	if !almostEqual(s.Pos, 0) {
		t.Fatalf("expected pos 0, got %v", s.Pos)
	}
	if !almostEqual(s.AngularFrequency, 6) {
		t.Fatalf("expected omega 6, got %v", s.AngularFrequency)
	}
	if !almostEqual(s.DampingRatio, 0.45) {
		t.Fatalf("expected zeta 0.45, got %v", s.DampingRatio)
	}
}

func TestSpringUpdateMovesTowardTarget(t *testing.T) {
	s := NewSpring()
	for i := 0; i < 600; i++ {
		s.Update()
	}
	if s.Pos <= 0 {
		t.Fatalf("expected positive pos after updates, got %v", s.Pos)
	}
	if s.Pos > 1.01 {
		t.Fatalf("expected pos <= 1.01, got %v", s.Pos)
	}
}

func TestSpringPoke(t *testing.T) {
	s := NewSpring()
	s.Poke()
	if s.Target != 0 {
		t.Fatalf("expected target 0 after poke, got %v", s.Target)
	}
	s.Poke()
	if s.Target != 1 {
		t.Fatalf("expected target 1 after second poke, got %v", s.Target)
	}
}

func TestNewProgressDefaults(t *testing.T) {
	p := NewProgress()
	if p == nil || p.Spring == nil {
		t.Fatal("expected progress with spring")
	}
	if p.Width != 40 || p.Label != "progress" {
		t.Fatalf("unexpected defaults: width=%d label=%s", p.Width, p.Label)
	}
}

func TestProgressRenderContainsBar(t *testing.T) {
	p := NewProgress()
	p.SetPercent(0.5)
	view := p.Render()
	if view == "" {
		t.Fatal("expected non-empty render")
	}
}
