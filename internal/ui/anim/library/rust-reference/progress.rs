//! Spring-eased progress bar (Bubble Tea progress-animated DNA).

use ratatui::prelude::*;

use crate::spring::Spring;

#[derive(Debug, Clone)]
pub struct Progress {
    spring: Spring,
    width: usize,
    label: String,
}

impl Default for Progress {
    fn default() -> Self {
        Self::new()
    }
}

impl Progress {
    pub fn new() -> Self {
        let mut spring = Spring::new();
        spring.target = 0.0;
        spring.pos = 0.0;
        Self {
            spring,
            width: 40,
            label: "progress".into(),
        }
    }

    pub fn set_width(&mut self, width: usize) {
        self.width = width.max(12);
        self.spring.set_width(self.width);
    }

    pub fn set_percent(&mut self, pct: f64) {
        self.spring.target = pct.clamp(0.0, 1.0);
    }

    pub fn incr_percent(&mut self, delta: f64) {
        self.set_percent(self.spring.target + delta);
    }

    pub fn poke(&mut self) {
        self.incr_percent(0.25);
        if self.spring.target >= 1.0 - f64::EPSILON {
            self.spring.target = 0.0;
            self.spring.pos = 0.0;
            self.spring.vel = 0.0;
        }
    }

    pub fn update(&mut self) {
        self.spring.update();
    }

    pub fn percent(&self) -> f64 {
        self.spring.pos.clamp(0.0, 1.0)
    }

    pub fn view_lines(&self) -> Vec<Line<'static>> {
        let w = self.width.min(48);
        let fill = ((self.percent() * w as f64).round() as usize).min(w);
        let bar = format!("{}{}", "█".repeat(fill), "░".repeat(w - fill));
        let pct = (self.percent() * 100.0).round() as i32;
        vec![
            Line::from(Span::styled(
                self.label.clone(),
                Style::default().add_modifier(Modifier::BOLD),
            )),
            Line::from(""),
            Line::from(Span::styled(
                format!("[{bar}] {pct}%"),
                Style::default().fg(Color::Rgb(0x6B, 0x50, 0xFF)),
            )),
            Line::from(""),
            Line::from(format!(
                "target={:.0}% · harmonica spring · space +25%",
                self.spring.target * 100.0
            )),
        ]
    }
}
