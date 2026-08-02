//! Rotating gradient splash screen (ported from CRUSH/bubbletea examples).

use ratatui::prelude::*;
use ratatui::text::{Line, Span};
use std::time::Instant;

/// Rotating half-block gradient splash.
#[derive(Debug, Clone)]
pub struct Splash {
    width: usize,
    height: usize,
    rate: f64,
    started: Instant,
}

impl Default for Splash {
    fn default() -> Self {
        Self::new()
    }
}

impl Splash {
    pub fn new() -> Self {
        Self {
            width: 0,
            height: 0,
            rate: 90.0,
            started: Instant::now(),
        }
    }

    pub fn set_size(&mut self, width: usize, height: usize) {
        self.width = width;
        // Half-block rows encode 2 pixels of color height each.
        self.height = height.max(1);
    }

    pub fn update(&mut self) -> bool {
        true
    }

    /// ANSI string view (useful for raw terminal dumps / tests).
    pub fn view(&self) -> String {
        if self.width == 0 || self.height == 0 {
            return "Resize terminal to see gradient splash…".to_string();
        }

        let mut out = String::with_capacity(self.width * self.height * 24);
        for line in self.color_rows() {
            for (c1, c2) in line {
                out.push_str(&format!(
                    "\x1b[38;2;{};{};{};48;2;{};{};{}m▀\x1b[0m",
                    c1.0, c1.1, c1.2, c2.0, c2.1, c2.2
                ));
            }
            out.push('\n');
        }
        out
    }

    /// Ratatui lines for embedding in a Frame.
    pub fn view_lines(&self) -> Vec<Line<'static>> {
        if self.width == 0 || self.height == 0 {
            return vec![Line::from("Resize terminal to see gradient splash…")];
        }

        self.color_rows()
            .into_iter()
            .map(|row| {
                let spans: Vec<Span> = row
                    .into_iter()
                    .map(|(c1, c2)| {
                        Span::styled(
                            "▀",
                            Style::default()
                                .fg(Color::Rgb(c1.0, c1.1, c1.2))
                                .bg(Color::Rgb(c2.0, c2.1, c2.2)),
                        )
                    })
                    .collect();
                Line::from(spans)
            })
            .collect()
    }

    fn color_rows(&self) -> Vec<Vec<((u8, u8, u8), (u8, u8, u8))>> {
        let t = self.started.elapsed().as_secs_f64() * self.rate;
        let angle = -t * std::f64::consts::PI / 180.0;
        let (sin_a, cos_a) = angle.sin_cos();
        let cx = self.width as f64 / 2.0;
        let cy = self.height as f64;
        let colors = SPLASH_COLORS;

        let mut rows = Vec::with_capacity(self.height);
        for line_y in 0..self.height {
            let py = line_y as f64 * 2.0 - cy;
            let mut px = 0.0 - cx;
            let x1 = (cx + (px * cos_a - py * sin_a)) / self.width as f64;
            let x2 = (cx + (px * cos_a - (py + 1.0) * sin_a)) / self.width as f64;
            px = self.width as f64 - cx;
            let end_x1 = (cx + (px * cos_a - py * sin_a)) / self.width as f64;
            let delta_x = (end_x1 - x1) / self.width as f64;

            let mut row = Vec::with_capacity(self.width);
            if delta_x.abs() < 0.0001 {
                let c1 = splash_color_at(x1, &colors);
                let c2 = splash_color_at(x2, &colors);
                for _ in 0..self.width {
                    row.push((c1, c2));
                }
            } else {
                for x in 0..self.width {
                    let pos1 = x1 + x as f64 * delta_x;
                    let pos2 = x2 + x as f64 * delta_x;
                    row.push((splash_color_at(pos1, &colors), splash_color_at(pos2, &colors)));
                }
            }
            rows.push(row);
        }
        rows
    }
}

const SPLASH_COLORS: [(u8, u8, u8); 12] = [
    (0x88, 0x11, 0x77),
    (0xaa, 0x33, 0x55),
    (0xcc, 0x66, 0x66),
    (0xee, 0x99, 0x44),
    (0xee, 0xdd, 0x00),
    (0x99, 0xdd, 0x55),
    (0x44, 0xdd, 0x88),
    (0x22, 0xcc, 0xbb),
    (0x00, 0xbb, 0xcc),
    (0x00, 0x99, 0xcc),
    (0x33, 0x66, 0xbb),
    (0x66, 0x33, 0x99),
];

fn splash_color_at(pos: f64, colors: &[(u8, u8, u8)]) -> (u8, u8, u8) {
    let pos = pos.clamp(0.0, 1.0);
    let idx = pos * (colors.len() - 1) as f64;
    let i1 = idx.floor() as usize % colors.len();
    let i2 = (i1 + 1) % colors.len();
    let t = idx - idx.floor();
    lerp_color(colors[i1], colors[i2], t)
}

fn lerp_color(c1: (u8, u8, u8), c2: (u8, u8, u8), t: f64) -> (u8, u8, u8) {
    (
        (c1.0 as f64 * (1.0 - t) + c2.0 as f64 * t) as u8,
        (c1.1 as f64 * (1.0 - t) + c2.1 as f64 * t) as u8,
        (c1.2 as f64 * (1.0 - t) + c2.2 as f64 * t) as u8,
    )
}
