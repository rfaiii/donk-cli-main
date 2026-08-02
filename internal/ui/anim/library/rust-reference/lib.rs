//! DONK Animations — Crush/Charm ports + custom effects (Rust / ratatui).
//!
//! Hunt DNA in Go (`ref/crush-anim/`, Bubble Tea examples); ship here.

pub mod boot;
pub mod cycling;
pub mod doomfire;
pub mod eyes;
pub mod matrix;
pub mod mission;
pub mod progress;
pub mod registry;
pub mod showcase;
pub mod space;
pub mod splash;
pub mod spinner;
pub mod spring;
pub mod vanish;

pub use boot::{BootCorner, BootSplash};
pub use cycling::Cycling;
pub use doomfire::DoomFire;
pub use eyes::Eyes;
pub use matrix::MatrixRain;
pub use mission::{SceneKind, ScenePlayer};
pub use progress::Progress;
pub use showcase::Showcase;
pub use space::SpaceField;
pub use splash::Splash;
pub use spinner::{Spinner, SpinnerKind};
pub use spring::Spring;
pub use vanish::Vanish;
