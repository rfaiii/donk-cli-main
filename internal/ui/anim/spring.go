// Package anim provides harmonica-style spring physics ported from
// the Rust donk-anim crate.
package anim

import "math"

// Spring is a damped harmonic oscillator used by animated progress bars
// and spring-physics effects.
type Spring struct {
	Pos               float64
	Vel               float64
	Target            float64
	AngularFrequency  float64
	DampingRatio      float64
	Width             int
	Frame             uint64

	posPos, posVel float64
	velPos, velVel float64
}

// NewSpring creates a spring with default Charm-style parameters.
func NewSpring() *Spring {
	s := &Spring{
		Pos:              0,
		Vel:              0,
		Target:           1,
		AngularFrequency: 6,
		DampingRatio:     0.45,
		Width:            48,
		Frame:            0,
	}
	s.recomputeCoeffs(1.0 / 60.0)
	return s
}

func (s *Spring) setWidth(width int) {
	if width > s.Width {
		s.Width = width
	}
}

func (s *Spring) recomputeCoeffs(dt float64) {
	eps := math.Nextafter(1, 2) - 1
	eps *= 8
	omega := s.AngularFrequency
	if omega < 0 {
		omega = 0
	}
	zeta := s.DampingRatio
	if zeta < 0 {
		zeta = 0
	}

	if omega < eps {
		s.posPos, s.posVel = 1, 0
		s.velPos, s.velVel = 0, 1
		return
	}

	if zeta > 1+eps {
		za := -omega * zeta
		zb := omega * math.Sqrt(zeta*zeta-1)
		z1, z2 := za-zb, za+zb
		e1, e2 := math.Exp(z1*dt), math.Exp(z2*dt)
		inv := 1 / (2 * zb)
		e1h, e2h := e1*inv, e2*inv
		s.posPos = e1h*z2 - z2*e2h + e2
		s.posVel = -e1h + e2h
		s.velPos = (z1*e1h - z2*e2h + e2) * z2
		s.velVel = -z1*e1h + z2*e2h
	} else if zeta < 1-eps {
		omegaZeta := omega * zeta
		alpha := omega * math.Sqrt(1-zeta*zeta)
		exp := math.Exp(-omegaZeta * dt)
		cos := math.Cos(alpha * dt)
		sin := math.Sin(alpha * dt)
		invAlpha := 1 / alpha
		expSin := exp * sin
		expCos := exp * cos
		expOzSin := exp * omegaZeta * sin * invAlpha
		s.posPos = expCos + expOzSin
		s.posVel = expSin * invAlpha
		s.velPos = -expSin*alpha - omegaZeta*expOzSin
		s.velVel = expCos - expOzSin
	} else {
		exp := math.Exp(-omega * dt)
		te := dt * exp
		tef := te * omega
		s.posPos = tef + exp
		s.posVel = te
		s.velPos = -omega * tef
		s.velVel = -tef + exp
	}
}

// Poke toggles the spring target and injects velocity.
func (s *Spring) Poke() {
	if s.Target > 0.5 {
		s.Target = 0
		s.Vel -= 2
	} else {
		s.Target = 1
		s.Vel += 2
	}
}

// Update steps the spring simulation by one tick.
func (s *Spring) Update() {
	s.recomputeCoeffs(1.0 / 60.0)
	oldPos := s.Pos - s.Target
	oldVel := s.Vel
	s.Pos = oldPos*s.posPos + oldVel*s.posVel + s.Target
	s.Vel = oldPos*s.velPos + oldVel*s.velVel
	s.Frame++
}
