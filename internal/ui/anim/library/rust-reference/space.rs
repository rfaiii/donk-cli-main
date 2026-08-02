//! Half-block space scroller — Bubble Tea `examples/space` port.

use ratatui::prelude::*;
use rand::Rng;

#[derive(Debug, Clone)]
pub struct SpaceField {
    width: usize,
    height: usize,
    /// Double-height grayscale field (0–255).
    field: Vec<Vec<u8>>,
    frame: u64,
}

impl Default for SpaceField {
    fn default() -> Self {
        Self::new()
    }
}

impl SpaceField {
    pub fn new() -> Self {
        let mut s = Self {
            width: 40,
            height: 12,
            field: Vec::new(),
            frame: 0,
        };
        s.rebuild();
        s
    }

    pub fn set_size(&mut self, width: usize, height: usize) {
        let w = width.max(8);
        let h = height.max(4);
        if w != self.width || h != self.height {
            self.width = w;
            self.height = h;
            self.rebuild();
        }
    }

    fn rebuild(&mut self) {
        let rows = self.height * 2;
        let mut rng = rand::rng();
        self.field = Vec::with_capacity(rows);
        for y in 0..rows {
            let mut row = Vec::with_capacity(self.width);
            let falloff = 1.0 - (y as f64 / rows as f64);
            for _x in 0..self.width {
                let base = falloff * falloff;
                let jitter: f64 = rng.random_range(-0.1..0.1);
                let v = (base + jitter).clamp(0.0, 1.0);
                row.push((v * 255.0) as u8);
            }
            self.field.push(row);
        }
    }

    pub fn update(&mut self) {
        self.frame = self.frame.wrapping_add(1);
    }

    pub fn view_lines(&self) -> Vec<Line<'static>> {
        let mut lines = Vec::with_capacity(self.height + 2);
        lines.push(Line::from(Span::styled(
            format!("space · frame {}", self.frame),
            Style::default().add_modifier(Modifier::BOLD),
        )));
        for y in 0..self.height {
            let mut spans = Vec::with_capacity(self.width);
            for x in 0..self.width {
                let xi = (x + self.frame as usize) % self.width.max(1);
                let fg = self.field.get(y * 2).and_then(|r| r.get(xi)).copied().unwrap_or(0);
                let bg = self
                    .field
                    .get(y * 2 + 1)
                    .and_then(|r| r.get(xi))
                    .copied()
                    .unwrap_or(0);
                spans.push(Span::styled(
                    "▀",
                    Style::default()
                        .fg(Color::Rgb(fg, fg, fg))
                        .bg(Color::Rgb(bg, bg, bg)),
                ));
            }
            lines.push(Line::from(spans));
        }
        lines
    }
}
