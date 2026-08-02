//! Char-cycling gradient animation — Crush `internal/ui/anim` DNA in Rust.
//!
//! Layout: `[scramble][gap][label][ellipsis]` · label resolves without mid-scramble.

use ratatui::prelude::*;
use ratatui::style::Color;
use std::time::Duration;

const ELLIPSIS_FRAMES: &[&str] = &[".", "..", "...", ""];
const ELLIPSIS_SPEED_STEPS: u64 = 8;
const MAX_BIRTH_STEPS: u64 = 20;

/// Char-cycling gradient animation (CRUSH-style "Generating..." indicator)
#[derive(Debug, Clone)]
pub struct Cycling {
    step: u64,
    cycling_chars: Vec<CyclingChar>,
    label_chars: Vec<CyclingChar>,
    birth_steps: Vec<u64>,
    ramp: Vec<Color>,
    cycle_colors: bool,
    ellipsis_on: bool,
    label: &'static str,
}

#[derive(Debug, Clone)]
struct CyclingChar {
    final_value: char,
    current_value: char,
}

impl Default for Cycling {
    fn default() -> Self {
        Self::new()
    }
}

impl Cycling {
    pub fn new() -> Self {
        Self::with_label(" Generating")
    }

    pub fn with_label(label: &'static str) -> Self {
        let ramp = lerp_ramp(
            Color::Rgb(0xFF, 0x00, 0x00),
            Color::Rgb(0x00, 0x00, 0xFF),
            12,
        );
        let n = 10usize;
        let birth_steps: Vec<u64> = (0..n)
            .map(|i| ((i as u64) * MAX_BIRTH_STEPS) / n.max(1) as u64)
            .collect();

        let cycling_chars = (0..n)
            .map(|_| CyclingChar {
                final_value: '\0',
                current_value: '.',
            })
            .collect();

        let label_chars = label
            .chars()
            .map(|ch| CyclingChar {
                final_value: ch,
                current_value: '.',
            })
            .collect();

        Self {
            step: 0,
            cycling_chars,
            label_chars,
            birth_steps,
            ramp,
            cycle_colors: true,
            ellipsis_on: false,
            label,
        }
    }

    pub fn update(&mut self, _elapsed: Duration) {
        self.step = self.step.wrapping_add(1);

        if self.cycle_colors && self.step % 3 == 0 {
            self.ramp.rotate_left(1);
        }

        for (i, c) in self.cycling_chars.iter_mut().enumerate() {
            let birth = self.birth_steps.get(i).copied().unwrap_or(0);
            if self.step < birth {
                c.current_value = '.';
            } else {
                c.current_value = random_char();
            }
        }

        // Label: dots until birth step, then final glyph (no mid-scramble).
        let label_birth_base = MAX_BIRTH_STEPS / 2;
        let mut label_done = true;
        for (i, c) in self.label_chars.iter_mut().enumerate() {
            let birth = label_birth_base + (i as u64 % 8);
            if self.step < birth {
                c.current_value = '.';
                label_done = false;
            } else {
                c.current_value = c.final_value;
            }
        }

        if label_done && !self.ellipsis_on {
            self.ellipsis_on = true;
        }
    }

    pub fn view_line(&self) -> Line<'static> {
        let mut spans = Vec::new();
        // Crush order: scramble → gap → label → ellipsis
        for (i, c) in self.cycling_chars.iter().enumerate() {
            let color = self
                .ramp
                .get(i % self.ramp.len())
                .cloned()
                .unwrap_or(Color::White);
            spans.push(Span::styled(
                c.current_value.to_string(),
                Style::default().fg(color),
            ));
        }
        spans.push(Span::raw(" "));
        for c in &self.label_chars {
            spans.push(Span::styled(
                c.current_value.to_string(),
                Style::default().fg(Color::Rgb(0xF9, 0x67, 0xDC)),
            ));
        }
        if self.ellipsis_on {
            spans.push(Span::styled(
                self.ellipsis_frame().to_string(),
                Style::default().fg(Color::Rgb(0x6B, 0x50, 0xFF)),
            ));
        }
        Line::from(spans)
    }

    pub fn view_lines(&self) -> Vec<Line<'static>> {
        vec![
            Line::from(Span::styled(
                format!("crush anim · {}", self.label.trim()),
                Style::default().add_modifier(Modifier::BOLD),
            )),
            Line::from(""),
            self.view_line(),
            Line::from(""),
            Line::from(Span::styled(
                "scramble · label · ellipsis (Go → Rust)",
                Style::default().fg(Color::DarkGray),
            )),
        ]
    }

    pub fn view(&self) -> String {
        let mut buf = String::with_capacity(64);
        for (i, c) in self.cycling_chars.iter().enumerate() {
            if let Some(color) = self.ramp.get(i) {
                buf.push_str(&format!(
                    "\x1b[38;2;{};{};{}m{}\x1b[0m",
                    color_r(color),
                    color_g(color),
                    color_b(color),
                    c.current_value
                ));
            } else {
                buf.push(c.current_value);
            }
        }
        buf.push(' ');
        for c in &self.label_chars {
            buf.push(c.current_value);
        }
        if self.ellipsis_on {
            buf.push_str(self.ellipsis_frame());
        }
        buf
    }

    fn ellipsis_frame(&self) -> &'static str {
        let i = ((self.step / ELLIPSIS_SPEED_STEPS) as usize) % ELLIPSIS_FRAMES.len();
        ELLIPSIS_FRAMES[i]
    }
}

fn lerp_ramp(a: Color, b: Color, n: usize) -> Vec<Color> {
    let (ar, ag, ab) = rgb_parts(a);
    let (br, bg, bb) = rgb_parts(b);
    let mut out = Vec::with_capacity(n);
    for i in 0..n {
        let t = if n <= 1 {
            0.0
        } else {
            i as f32 / (n - 1) as f32
        };
        let (mr, mg, mb) = (0xF9_u8, 0x67_u8, 0xDC_u8);
        let (r1, g1, b1) = if t < 0.5 {
            let u = t * 2.0;
            (lerp_u8(ar, mr, u), lerp_u8(ag, mg, u), lerp_u8(ab, mb, u))
        } else {
            let u = (t - 0.5) * 2.0;
            (lerp_u8(mr, br, u), lerp_u8(mg, bg, u), lerp_u8(mb, bb, u))
        };
        out.push(Color::Rgb(r1, g1, b1));
    }
    out
}

fn lerp_u8(a: u8, b: u8, t: f32) -> u8 {
    (a as f32 + (b as f32 - a as f32) * t).round() as u8
}

fn rgb_parts(c: Color) -> (u8, u8, u8) {
    match c {
        Color::Rgb(r, g, b) => (r, g, b),
        _ => (255, 255, 255),
    }
}

fn random_char() -> char {
    const CHARS: &[char] = &[
        '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f', 'A', 'B',
        'C', 'D', 'E', 'F', '~', '!', '@', '#', '$', '£', '€', '%', '^', '&', '*', '(', ')', '+',
        '=', '_',
    ];
    let idx = rand::random_range(0..CHARS.len());
    CHARS[idx]
}

fn color_r(c: &Color) -> u8 {
    match c {
        Color::Rgb(r, _, _) => *r,
        _ => 255,
    }
}
fn color_g(c: &Color) -> u8 {
    match c {
        Color::Rgb(_, g, _) => *g,
        _ => 255,
    }
}
fn color_b(c: &Color) -> u8 {
    match c {
        Color::Rgb(_, _, b) => *b,
        _ => 255,
    }
}
