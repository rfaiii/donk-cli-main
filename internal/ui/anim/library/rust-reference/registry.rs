//! Animation registry for DONK CLI
//!
//! This module provides a unified registry for all animation effects
//! across the application, including:
//! - Terminal UI animations (splash, loading, transitions)
//! - Background/standby animations (doomfire, matrix, etc.)
//! - Cursor/caret animations (shader references)
//! - Audio-reactive animations (if audio context available)

use crate::boot::BootSplash;
use crate::cycling::Cycling;
use crate::doomfire::DoomFire;
use crate::eyes::Eyes;
use crate::matrix::MatrixRain;
use crate::progress::Progress;
use crate::showcase::Showcase;
use crate::space::SpaceField;
use crate::splash::Splash;
use crate::spinner::Spinner;
use crate::spring::Spring;
use crate::vanish::Vanish;

#[derive(Debug, Clone, Copy, PartialEq, Eq, Hash)]
pub enum AnimationId {
    Boot,
    Splash,
    Cycling,
    Spinner,
    DoomFire,
    Spring,
    Progress,
    Eyes,
    Matrix,
    Space,
    Vanish,
    Showcase,
}

#[derive(Debug)]
pub enum Animation {
    Boot(BootSplash),
    Splash(Splash),
    Cycling(Cycling),
    Spinner(Spinner),
    DoomFire(DoomFire),
    Spring(Spring),
    Progress(Progress),
    Eyes(Eyes),
    Matrix(MatrixRain),
    Space(SpaceField),
    Vanish(Vanish),
    Showcase(Showcase),
}

impl Animation {
    pub fn update(&mut self) {
        match self {
            Animation::Boot(b) => {
                b.update();
            }
            Animation::Splash(s) => {
                s.update();
            }
            Animation::Cycling(c) => c.update(std::time::Duration::ZERO),
            Animation::Spinner(s) => s.update(),
            Animation::DoomFire(d) => {
                d.update();
            }
            Animation::Spring(s) => {
                s.update();
            }
            Animation::Progress(p) => p.update(),
            Animation::Eyes(e) => {
                e.update();
            }
            Animation::Matrix(m) => {
                m.update();
            }
            Animation::Space(s) => s.update(),
            Animation::Vanish(v) => {
                v.update();
            }
            Animation::Showcase(s) => {
                s.update();
            }
        }
    }
}
