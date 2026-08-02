//! Blinking eyes — port of bubbletea examples/eyes (ESP32 blink math).

use ratatui::prelude::*;
use std::time::{Duration, Instant};

const EYE_W: i32 = 11;
const EYE_H: i32 = 7;
const BLINK_FRAMES: i32 = 16;

#[derive(Debug, Clone)]
pub struct Eyes {
    width: usize,
    height: usize,
    last_blink: Instant,
    open_for: Duration,
    blinking: bool,
    blink_state: i32,
    frame: u64,
}

impl Default for Eyes {
    fn default() -> Self {
        Self::new()
    }
}

impl Eyes {
    pub fn new() -> Self {
        Self {
            width: 48,
            height: 14,
            last_blink: Instant::now(),
            open_for: Duration::from_millis(1200 + (rand::random::<u64>() % 2200)),
            blinking: false,
            blink_state: 0,
            frame: 0,
        }
    }

    pub fn set_size(&mut self, width: usize, height: usize) {
        self.width = width.max(24);
        self.height = height.max(8);
    }

    pub fn update(&mut self) {
        self.frame = self.frame.wrapping_add(1);
        let now = Instant::now();
        if !self.blinking && now.duration_since(self.last_blink) >= self.open_for {
            self.blinking = true;
            self.blink_state = 0;
        }
        if self.blinking {
            self.blink_state += 1;
            if self.blink_state >= BLINK_FRAMES {
                self.blinking = false;
                self.last_blink = now;
                self.open_for = Duration::from_millis(800 + (rand::random::<u64>() % 2800));
                if rand::random::<u8>() % 10 == 0 {
                    self.open_for = Duration::from_millis(280);
                }
            }
        }
    }

    fn eye_height(&self) -> i32 {
        if !self.blinking {
            return EYE_H;
        }
        let half = BLINK_FRAMES / 2;
        let progress = if self.blink_state < half {
            let t = self.blink_state as f64 / half as f64;
            1.0 - t * t
        } else {
            let t = (self.blink_state - half) as f64 / half as f64;
            t * (2.0 - t)
        };
        ((EYE_H as f64) * progress).max(1.0) as i32
    }

    pub fn view_lines(&self) -> Vec<Line<'static>> {
        let w = self.width;
        let h = self.height.min(18);
        let mut canvas = vec![vec![' '; w]; h];
        let cy = (h as i32) / 2;
        let left_x = (w as i32) / 2 - 10;
        let right_x = (w as i32) / 2 + 10;
        let ry = self.eye_height();
        draw_ellipse(&mut canvas, left_x, cy, EYE_W, ry);
        draw_ellipse(&mut canvas, right_x, cy, EYE_W, ry);

        let mut lines = vec![
            Line::from(Span::styled(
                "eyes · bubbletea port".to_string(),
                Style::default().add_modifier(Modifier::BOLD),
            )),
            Line::from(""),
        ];
        for row in canvas {
            let s: String = row.into_iter().collect();
            lines.push(Line::from(Span::styled(
                s,
                Style::default().fg(Color::Rgb(0xF0, 0xF0, 0xF0)),
            )));
        }
        lines
    }
}

fn draw_ellipse(canvas: &mut [Vec<char>], x0: i32, y0: i32, rx: i32, ry: i32) {
    let h = canvas.len() as i32;
    let w = canvas.first().map(|r| r.len() as i32).unwrap_or(0);
    for y in -ry..=ry {
        let denom = (ry as f64).max(1.0);
        let width = ((rx as f64) * (1.0 - (y as f64 / denom).powi(2)).sqrt()) as i32;
        for x in -width..=width {
            let cx = x0 + x;
            let cy = y0 + y;
            if cx >= 0 && cy >= 0 && cx < w && cy < h {
                canvas[cy as usize][cx as usize] = '●';
            }
        }
    }
}
