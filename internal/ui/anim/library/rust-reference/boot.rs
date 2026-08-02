//! Gemini/Crush-style boot reel — cycle ASCII brand marks over splash art.
//!
//! Logos come from `donk_assets::AsciiLogo::boot_reel()`. Transitions use a
//! light dissolve (vanish DNA); chrome uses Crush-style cycling text.

use donk_assets::AsciiLogo;
use ratatui::prelude::*;
use std::time::{Duration, Instant};

use crate::cycling::Cycling;
use crate::splash::Splash;

const HOLD_MS: u64 = 720;
const FADE_MS: u64 = 220;
const CORNER_SWAP_MS: u64 = 2600;

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum BootCorner {
    TopLeft,
    TopRight,
}

/// Animated boot sequence: corner hero + gradient art + cycling logo reel.
#[derive(Debug, Clone)]
pub struct BootSplash {
    splash: Splash,
    cycling: Cycling,
    reel: Vec<AsciiLogo>,
    index: usize,
    hold_started: Instant,
    boot_started: Instant,
    frame: u64,
    width: usize,
    height: usize,
}

impl Default for BootSplash {
    fn default() -> Self {
        Self::new()
    }
}

impl BootSplash {
    pub fn new() -> Self {
        Self {
            splash: Splash::new(),
            cycling: Cycling::with_label(" Starting"),
            reel: AsciiLogo::boot_reel().to_vec(),
            index: 0,
            hold_started: Instant::now(),
            boot_started: Instant::now(),
            frame: 0,
            width: 80,
            height: 24,
        }
    }

    pub fn set_size(&mut self, width: usize, height: usize) {
        self.width = width.max(20);
        self.height = height.max(8);
        // Gradient fills most of the frame; logo floats on top.
        self.splash.set_size(self.width, self.height.max(1));
    }

    pub fn update(&mut self) {
        self.frame = self.frame.wrapping_add(1);
        self.splash.update();
        self.cycling
            .update(Duration::from_millis(16));

        let elapsed = self.hold_started.elapsed();
        if elapsed >= Duration::from_millis(HOLD_MS) && !self.reel.is_empty() {
            self.index = (self.index + 1) % self.reel.len();
            self.hold_started = Instant::now();
        }
    }

    pub fn elapsed(&self) -> Duration {
        self.boot_started.elapsed()
    }

    pub fn current_logo(&self) -> AsciiLogo {
        self.reel
            .get(self.index)
            .copied()
            .unwrap_or(AsciiLogo::TinyDonkey)
    }

    /// Giant DONK-CLI mark for the top corner (catalog clear / large / intro).
    pub fn corner_hero(&self) -> AsciiLogo {
        let alts = AsciiLogo::boot_corner_alts();
        let i = (self.boot_started.elapsed().as_millis() as u64 / CORNER_SWAP_MS) as usize
            % alts.len().max(1);
        let candidate = alts.get(i).copied().unwrap_or_else(AsciiLogo::boot_corner_hero);
        // Fall back to the top-L/R mark if the giant won't fit.
        if candidate.art_width() + 2 > self.width {
            AsciiLogo::TopCentered
        } else {
            candidate
        }
    }

    pub fn corner_side(&self) -> BootCorner {
        // Alternate top-left / top-right during the boot (~catalog #14 usage).
        if (self.boot_started.elapsed().as_millis() as u64 / CORNER_SWAP_MS) % 2 == 0 {
            BootCorner::TopLeft
        } else {
            BootCorner::TopRight
        }
    }

    pub fn corner_rect(&self, area: Rect) -> Rect {
        let hero = self.corner_hero();
        let w = (hero.art_width() as u16 + 2).min(area.width);
        let h = (hero.art_height() as u16 + 1).min(area.height.saturating_sub(2).max(1));
        let x = match self.corner_side() {
            BootCorner::TopLeft => area.x,
            BootCorner::TopRight => area.x + area.width.saturating_sub(w),
        };
        Rect {
            x,
            y: area.y,
            width: w,
            height: h,
        }
    }

    /// Giant pinned corner wordmark (no dissolve — always solid).
    pub fn corner_logo_lines(&self, theme_pink: Color, theme_purple: Color) -> Vec<Line<'static>> {
        let logo = self.corner_hero();
        let pulse = ((self.frame as f32 * 0.08).sin() * 0.5 + 0.5).clamp(0.0, 1.0);
        let accent = lerp_color(theme_pink, theme_purple, pulse * 0.45);
        let mut lines = Vec::new();
        for (row_i, row) in logo.art_lines().enumerate() {
            let mut spans = Vec::new();
            for (col_i, ch) in row.chars().enumerate() {
                let sweep = ((col_i as f32 * 0.035 + self.frame as f32 * 0.04).sin() * 0.5 + 0.5)
                    .clamp(0.0, 1.0);
                let fg = if ch == ' ' {
                    Color::Reset
                } else {
                    lerp_color(accent, theme_purple, sweep * 0.65)
                };
                let style = if ch == ' ' {
                    Style::default()
                } else {
                    Style::default().fg(fg).add_modifier(Modifier::BOLD)
                };
                spans.push(Span::styled(ch.to_string(), style));
            }
            let _ = row_i;
            lines.push(Line::from(spans));
        }
        lines
    }

    /// Center reel + Crush starting line (below the corner hero).
    pub fn view_overlay_lines(&self, theme_pink: Color, theme_purple: Color) -> Vec<Line<'static>> {
        let logo = self.current_logo();
        let fade = self.dissolve_amount();
        let pulse = ((self.frame as f32 * 0.12).sin() * 0.5 + 0.5).clamp(0.0, 1.0);
        let accent = lerp_color(theme_pink, theme_purple, pulse);

        let mut art: Vec<Line<'static>> = Vec::new();
        for (row_i, row) in logo.art_lines().enumerate() {
            let mut spans = Vec::new();
            for (col_i, ch) in row.chars().enumerate() {
                let noise =
                    ((row_i * 31 + col_i * 17 + self.frame as usize * 3) % 100) as f32 / 100.0;
                let show = noise > fade * 0.9;
                let out = if ch == ' ' {
                    ' '
                } else if show {
                    ch
                } else {
                    match (row_i + col_i + self.frame as usize) % 5 {
                        0 => '·',
                        1 => '░',
                        2 => '▒',
                        3 => '▓',
                        _ => ' ',
                    }
                };
                let sweep = ((col_i as f32 * 0.04 + self.frame as f32 * 0.05).sin() * 0.5 + 0.5)
                    .clamp(0.0, 1.0);
                let fg = lerp_color(accent, theme_purple, sweep * 0.55);
                spans.push(Span::styled(out.to_string(), Style::default().fg(fg)));
            }
            art.push(Line::from(spans));
        }

        let mut lines = Vec::new();
        lines.push(Line::from(vec![
            Span::styled(
                " DONK ",
                Style::default()
                    .fg(theme_pink)
                    .add_modifier(Modifier::BOLD),
            ),
            Span::styled(
                format!("· {} ", logo.label()),
                Style::default().fg(Color::DarkGray),
            ),
        ]));
        lines.push(Line::from(""));
        lines.extend(art);
        lines.push(Line::from(""));
        lines.push(self.cycling.view_line());
        lines.push(Line::from(Span::styled(
            format!(
                "{}/{}  ·  press any key",
                self.index + 1,
                self.reel.len().max(1)
            ),
            Style::default().fg(Color::DarkGray),
        )));
        lines
    }

    pub fn backdrop_lines(&self) -> Vec<Line<'static>> {
        self.splash.view_lines()
    }

    fn dissolve_amount(&self) -> f32 {
        let ms = self.hold_started.elapsed().as_millis() as u64;
        if ms < FADE_MS {
            // Entering
            1.0 - (ms as f32 / FADE_MS as f32)
        } else if ms + FADE_MS > HOLD_MS {
            // Leaving
            let leave = HOLD_MS.saturating_sub(ms);
            1.0 - (leave as f32 / FADE_MS as f32)
        } else {
            0.0
        }
        .clamp(0.0, 1.0)
    }
}

fn lerp_color(a: Color, b: Color, t: f32) -> Color {
    let (ar, ag, ab) = match a {
        Color::Rgb(r, g, b) => (r, g, b),
        _ => (0xF9, 0x67, 0xDC),
    };
    let (br, bg, bb) = match b {
        Color::Rgb(r, g, b) => (r, g, b),
        _ => (0x6B, 0x50, 0xFF),
    };
    let t = t.clamp(0.0, 1.0);
    Color::Rgb(
        (ar as f32 + (br as f32 - ar as f32) * t) as u8,
        (ag as f32 + (bg as f32 - ag as f32) * t) as u8,
        (ab as f32 + (bb as f32 - ab as f32) * t) as u8,
    )
}
