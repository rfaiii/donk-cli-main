//! Matrix-style rain (standby / crush DNA).

use ratatui::prelude::*;
use rand::Rng;

#[derive(Debug, Clone)]
pub struct MatrixRain {
    width: usize,
    height: usize,
    cols: Vec<Col>,
    frame: u64,
}

#[derive(Debug, Clone)]
struct Col {
    y: f32,
    speed: f32,
    trail: usize,
}

impl Default for MatrixRain {
    fn default() -> Self {
        Self::new()
    }
}

impl MatrixRain {
    pub fn new() -> Self {
        Self {
            width: 40,
            height: 16,
            cols: Vec::new(),
            frame: 0,
        }
    }

    pub fn set_size(&mut self, width: usize, height: usize) {
        let w = width.max(12);
        let h = height.max(6);
        if w != self.width || h != self.height || self.cols.len() != w {
            self.width = w;
            self.height = h;
            let mut rng = rand::rng();
            self.cols = (0..w)
                .map(|_| Col {
                    y: rng.random_range(-(h as f32)..0.0),
                    speed: rng.random_range(0.25..1.1),
                    trail: rng.random_range(3..10),
                })
                .collect();
        }
    }

    pub fn update(&mut self) {
        self.frame = self.frame.wrapping_add(1);
        let h = self.height as f32;
        let mut rng = rand::rng();
        for c in &mut self.cols {
            c.y += c.speed;
            if c.y - c.trail as f32 > h {
                c.y = rng.random_range(-8.0..0.0);
                c.speed = rng.random_range(0.25..1.1);
                c.trail = rng.random_range(3..10);
            }
        }
    }

    pub fn view_lines(&self) -> Vec<Line<'static>> {
        let mut grid = vec![vec![(' ', Color::Reset); self.width]; self.height];
        const GLYPHS: &[char] = &[
            '0', '1', 'ｱ', 'ｲ', 'ｳ', 'ｴ', 'ｵ', 'ｶ', 'ｷ', 'ｸ', 'ｹ', 'ｺ', 'ｻ', 'ｼ', 'ｽ',
        ];
        for (x, col) in self.cols.iter().enumerate() {
            let head = col.y as i32;
            for t in 0..col.trail {
                let y = head - t as i32;
                if y < 0 || y >= self.height as i32 {
                    continue;
                }
                let ch = GLYPHS[(x * 7 + t * 3 + self.frame as usize) % GLYPHS.len()];
                let color = if t == 0 {
                    Color::Rgb(0xC8, 0xFF, 0xC8)
                } else if t < 3 {
                    Color::Rgb(0x40, 0xE0, 0x70)
                } else {
                    Color::Rgb(0x1A, 0x7A, 0x3A)
                };
                grid[y as usize][x] = (ch, color);
            }
        }
        let mut lines = vec![Line::from(Span::styled(
            "matrix rain".to_string(),
            Style::default().fg(Color::Green).add_modifier(Modifier::BOLD),
        ))];
        for row in grid {
            let spans: Vec<Span> = row
                .into_iter()
                .map(|(ch, color)| Span::styled(ch.to_string(), Style::default().fg(color)))
                .collect();
            lines.push(Line::from(spans));
        }
        lines
    }
}
