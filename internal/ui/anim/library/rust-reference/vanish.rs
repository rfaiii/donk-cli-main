//! Vanish / dissolve — bubbletea vanish DNA (glyph dissolve over time).

use ratatui::prelude::*;
use donk_assets::{AsciiLogo, LogoSlot};

#[derive(Debug, Clone)]
pub struct Vanish {
    step: u32,
    logo: AsciiLogo,
    base: Vec<String>,
}

impl Default for Vanish {
    fn default() -> Self {
        Self::new()
    }
}

impl Vanish {
    pub fn new() -> Self {
        let logo = AsciiLogo::for_slot(LogoSlot::Intro);
        let base = logo
            .as_str()
            .lines()
            .map(|l| l.to_string())
            .filter(|l| !l.trim().is_empty())
            .collect();
        Self {
            step: 0,
            logo,
            base,
        }
    }

    pub fn update(&mut self) {
        self.step = self.step.wrapping_add(1);
        if self.step % 90 == 0 {
            self.logo = self.logo.next();
            // Prefer wordmarks for vanish readability
            if matches!(
                self.logo,
                AsciiLogo::Binary | AsciiLogo::BoldDonkey | AsciiLogo::RegularDonkey
            ) {
                self.logo = AsciiLogo::IntroCompact;
            }
            self.base = self
                .logo
                .as_str()
                .lines()
                .map(|l| l.to_string())
                .filter(|l| !l.trim().is_empty())
                .collect();
        }
    }

    pub fn view_lines(&self) -> Vec<Line<'static>> {
        let phase = ((self.step % 60) as f32) / 60.0;
        let dissolve = if phase < 0.5 {
            phase * 2.0
        } else {
            (1.0 - phase) * 2.0
        };
        let mut lines = vec![
            Line::from(Span::styled(
                format!("vanish · {}", self.logo.label()),
                Style::default().add_modifier(Modifier::BOLD),
            )),
            Line::from(""),
        ];
        for (row_i, row) in self.base.iter().enumerate() {
            let mut spans = Vec::new();
            for (col_i, ch) in row.chars().enumerate() {
                let noise = ((row_i * 31 + col_i * 17 + self.step as usize) % 100) as f32 / 100.0;
                let show = noise > dissolve * 0.85;
                let out = if show {
                    ch
                } else if ch == ' ' {
                    ' '
                } else {
                    match (row_i + col_i + self.step as usize) % 4 {
                        0 => '·',
                        1 => '░',
                        2 => '▒',
                        _ => ' ',
                    }
                };
                spans.push(Span::styled(
                    out.to_string(),
                    Style::default().fg(Color::Rgb(0xF9, 0x67, 0xDC)),
                ));
            }
            lines.push(Line::from(spans));
        }
        lines
    }
}
