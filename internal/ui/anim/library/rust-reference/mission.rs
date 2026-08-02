//! Procedural ASCII animation scenes — ported from Mission Control Pro.
//!
//! Eight canvas-based scenes: fireworks, horizon, radio, rain, starfield,
//! galaxy, tunnel, plasma. Each renderer takes `(frame, width, height)` and
//! returns `Vec<String>` of canvas rows.

use std::f64::consts::PI;

/// 24-character brightness ramp from dark to bright.
const RAMP: &str = " .,:;irsXA253hMHGS#9B&@";

fn ramp_char(v: f64) -> char {
    let v = v.clamp(0.0, 1.0);
    let idx = (v * (RAMP.len() - 1) as f64).round() as usize;
    RAMP.as_bytes()[idx] as char
}

fn blank_canvas(width: usize, height: usize) -> Vec<Vec<char>> {
    vec![vec![' '; width]; height]
}

fn to_lines(canvas: &[Vec<char>]) -> Vec<String> {
    canvas.iter().map(|row| row.iter().collect()).collect()
}

fn stamp(canvas: &mut [Vec<char>], x: i32, y: i32, ch: char) {
    if y >= 0 && (y as usize) < canvas.len() {
        let row = &mut canvas[y as usize];
        if x >= 0 && (x as usize) < row.len() {
            row[x as usize] = ch;
        }
    }
}

fn draw_text(canvas: &mut [Vec<char>], x: i32, y: i32, text: &str) {
    for (i, ch) in text.chars().enumerate() {
        stamp(canvas, x + i as i32, y, ch);
    }
}

/// Scene identifiers matching Mission Control Pro.
#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum SceneKind {
    Fireworks,
    Horizon,
    Radio,
    Rain,
    Starfield,
    Galaxy,
    Tunnel,
    Plasma,
}

impl SceneKind {
    pub fn all() -> &'static [SceneKind] {
        &[
            SceneKind::Fireworks,
            SceneKind::Horizon,
            SceneKind::Radio,
            SceneKind::Rain,
            SceneKind::Starfield,
            SceneKind::Galaxy,
            SceneKind::Tunnel,
            SceneKind::Plasma,
        ]
    }

    pub fn next(self) -> Self {
        let all = Self::all();
        let i = all.iter().position(|&s| s == self).unwrap_or(0);
        all[(i + 1) % all.len()]
    }

    pub fn label(self) -> &'static str {
        match self {
            SceneKind::Fireworks => "fireworks",
            SceneKind::Horizon => "horizon",
            SceneKind::Radio => "radio",
            SceneKind::Rain => "rain",
            SceneKind::Starfield => "starfield",
            SceneKind::Galaxy => "galaxy",
            SceneKind::Tunnel => "tunnel",
            SceneKind::Plasma => "plasma",
        }
    }

    pub fn title(self) -> &'static str {
        match self {
            SceneKind::Fireworks => "Fireworks",
            SceneKind::Horizon => "Horizon",
            SceneKind::Radio => "Radio Waves",
            SceneKind::Rain => "Rain",
            SceneKind::Starfield => "Warp Starfield",
            SceneKind::Galaxy => "Galaxy",
            SceneKind::Tunnel => "Tunnel",
            SceneKind::Plasma => "Plasma",
        }
    }

    pub fn description(self) -> &'static str {
        match self {
            SceneKind::Fireworks => "Expanding circular bursts with particle rings",
            SceneKind::Horizon => "Animated sun over wave horizon",
            SceneKind::Radio => "Concentric expanding rings from center",
            SceneKind::Rain => "Falling streaks with splash effects",
            SceneKind::Starfield => "Warp-drive starfield (golden angle)",
            SceneKind::Galaxy => "4-arm spiral with radial glow",
            SceneKind::Tunnel => "Depth-based concentric bands",
            SceneKind::Plasma => "4-term sine wave interference",
        }
    }
}

/// Animation scene renderer.
#[derive(Debug, Clone)]
pub struct ScenePlayer {
    pub kind: SceneKind,
    pub frame: u64,
    width: usize,
    height: usize,
}

impl Default for ScenePlayer {
    fn default() -> Self {
        Self::new()
    }
}

impl ScenePlayer {
    pub fn new() -> Self {
        Self {
            kind: SceneKind::Fireworks,
            frame: 0,
            width: 40,
            height: 12,
        }
    }

    pub fn set_size(&mut self, width: usize, height: usize) {
        self.width = width.max(20);
        self.height = height.max(8);
    }

    pub fn set_kind(&mut self, kind: SceneKind) {
        self.kind = kind;
        self.frame = 0;
    }

    pub fn next_scene(&mut self) {
        self.kind = self.kind.next();
        self.frame = 0;
    }

    pub fn prev_scene(&mut self) {
        let all = SceneKind::all();
        let i = all.iter().position(|&s| s == self.kind).unwrap_or(0);
        let prev = if i == 0 { all.len() - 1 } else { i - 1 };
        self.kind = all[prev];
        self.frame = 0;
    }

    pub fn update(&mut self) {
        self.frame = self.frame.wrapping_add(1);
    }

    pub fn view_lines(&self) -> Vec<String> {
        let f = self.frame as i32;
        let w = self.width as i32;
        let h = self.height as i32;
        match self.kind {
            SceneKind::Fireworks => fireworks(f, w, h),
            SceneKind::Horizon => horizon(f, w, h),
            SceneKind::Radio => radio(f, w, h),
            SceneKind::Rain => rain(f, w, h),
            SceneKind::Starfield => starfield(f, w, h),
            SceneKind::Galaxy => galaxy(f, w, h),
            SceneKind::Tunnel => tunnel(f, w, h),
            SceneKind::Plasma => plasma(f, w, h),
        }
    }
}

// ─── Renderers ───────────────────────────────────────────────────────────────

fn fireworks(frame: i32, w: i32, h: i32) -> Vec<String> {
    let mut canvas = blank_canvas(w as usize, h as usize);
    let bursts = [
        (w / 4, h / 3, frame % 60),
        (w / 2, h / 2, (frame + 20) % 60),
        (3 * w / 4, h / 3, (frame + 40) % 60),
    ];
    for (cx, cy, phase) in bursts {
        let radius = phase as f64 * 0.8;
        let brightness = 1.0 - (phase as f64 / 60.0);
        for theta in (0..360).step_by(15) {
            let rad = theta as f64 * PI / 180.0;
            let px = cx as f64 + radius * rad.cos();
            let py = cy as f64 + radius * rad.sin() * 0.5;
            let ch = ramp_char(brightness);
            stamp(&mut canvas, px as i32, py as i32, ch);
        }
        for theta in (0..360).step_by(45) {
            let rad = theta as f64 * PI / 180.0;
            let r2 = radius * 1.3;
            let px = cx as f64 + r2 * rad.cos();
            let py = cy as f64 + r2 * rad.sin() * 0.5;
            let ch = ramp_char(brightness * 0.5);
            stamp(&mut canvas, px as i32, py as i32, ch);
        }
    }
    to_lines(&canvas)
}

fn horizon(frame: i32, w: i32, h: i32) -> Vec<String> {
    let mut canvas = blank_canvas(w as usize, h as usize);
    let t = frame as f64 * 0.03;

    for x in 0..w {
        for y in 0..h {
            let star = ((x as f64 * 0.7 + y as f64 * 1.3 + t * 0.5).sin() * 0.5 + 0.5)
                * ((x as f64 * 0.3).sin() * 0.5 + 0.5);
            if star > 0.93 {
                stamp(&mut canvas, x, y, '.');
            }
        }
    }

    let sun_x = (w as f64 / 2.0 + (t * 3.0).sin() * w as f64 * 0.3) as i32;
    let sun_y = (h as f64 / 3.0 + (t * 2.0).cos() * h as f64 * 0.15) as i32;
    for dy in -2..=2 {
        for dx in -3..=3 {
            let dist = ((dx * dx + dy * dy * 4) as f64).sqrt();
            if dist < 3.0 {
                stamp(&mut canvas, sun_x + dx, sun_y + dy, ramp_char(1.0 - dist / 3.0));
            }
        }
    }

    let horizon_y = (h as f64 * 0.65) as i32;
    for x in 0..w {
        let wave = ((x as f64 * 0.3 + t * 2.0).sin() * 2.0
            + (x as f64 * 0.1 + t).sin() * 1.0) as i32;
        let y = horizon_y + wave;
        stamp(&mut canvas, x, y, '~');
        if y + 1 < h {
            stamp(&mut canvas, x, y + 1, '≈');
        }
    }

    to_lines(&canvas)
}

fn radio(frame: i32, w: i32, h: i32) -> Vec<String> {
    let mut canvas = blank_canvas(w as usize, h as usize);
    let cx = w as f64 / 2.0;
    let cy = h as f64 / 2.0;
    let t = frame as f64 * 0.05;

    for ring in 0..8 {
        let radius = (t + ring as f64 * 1.5) % 12.0;
        let brightness = 1.0 - radius / 12.0;
        for theta in (0..360).step_by(5) {
            let rad = theta as f64 * PI / 180.0;
            let px = cx + radius * rad.cos();
            let py = cy + radius * rad.sin() * 0.5;
            let ch = ramp_char(brightness);
            stamp(&mut canvas, px as i32, py as i32, ch);
        }
    }

    stamp(&mut canvas, cx as i32, cy as i32, '+');
    to_lines(&canvas)
}

fn rain(frame: i32, w: i32, h: i32) -> Vec<String> {
    let mut canvas = blank_canvas(w as usize, h as usize);
    let t = frame as f64;

    for x in 0..w {
        let phase = (x as f64 * 0.7 + t * 0.5) % h as f64;
        let col_height = (phase) as i32;
        for y in 0..h {
            let drop_y = (y + col_height) % h;
            if (y + col_height) / h == col_height / h {
                let ch = if y == h - 1 || drop_y == h - 1 {
                    's'
                } else {
                    '|'
                };
                stamp(&mut canvas, x, drop_y, ch);
            }
        }
    }

    to_lines(&canvas)
}

fn starfield(frame: i32, w: i32, h: i32) -> Vec<String> {
    let mut canvas = blank_canvas(w as usize, h as usize);
    let t = frame as f64 * 0.02;
    let golden = 0.61803398875_f64;
    let num_stars = 80;

    for i in 0..num_stars {
        let seed = i as f64 * golden;
        let angle = seed * 2.0 * PI;
        let depth = ((t + i as f64 * 0.1) % 1.0);
        let speed = 0.5 + depth * 2.0;
        let r = depth * (w as f64).min(h as f64 * 2.0) * 0.5;

        let px = w as f64 / 2.0 + r * angle.cos();
        let py = h as f64 / 2.0 + r * angle.sin() * 0.5;

        let ch = if speed > 1.5 {
            ramp_char(depth)
        } else {
            '.'
        };
        stamp(&mut canvas, px as i32, py as i32, ch);
    }

    to_lines(&canvas)
}

fn galaxy(frame: i32, w: i32, h: i32) -> Vec<String> {
    let mut canvas = blank_canvas(w as usize, h as usize);
    let cx = w as f64 / 2.0;
    let cy = h as f64 / 2.0;
    let t = frame as f64 * 0.02;

    for arm in 0..4 {
        let arm_offset = arm as f64 * PI / 2.0;
        for i in 0..60 {
            let r = i as f64 * 0.15;
            let theta = arm_offset + r * 0.3 + t;
            let px = cx + r * theta.cos();
            let py = cy + r * theta.sin() * 0.5;
            let brightness = 1.0 - (i as f64 / 60.0);
            stamp(&mut canvas, px as i32, py as i32, ramp_char(brightness));
        }
    }

    for r in [3.0_f64, 6.0, 9.0] {
        for theta in (0..360).step_by(20) {
            let rad = theta as f64 * PI / 180.0 + t * 0.5;
            let px = cx + r * rad.cos();
            let py = cy + r * rad.sin() * 0.5;
            stamp(&mut canvas, px as i32, py as i32, '.');
        }
    }

    stamp(&mut canvas, cx as i32, cy as i32, '@');
    to_lines(&canvas)
}

fn tunnel(frame: i32, w: i32, h: i32) -> Vec<String> {
    let mut canvas = blank_canvas(w as usize, h as usize);
    let cx = w as f64 / 2.0;
    let cy = h as f64 / 2.0;
    let t = frame as f64 * 0.05;

    for y in 0..h {
        for x in 0..w {
            let dx = (x as f64 - cx).abs();
            let dy = (y as f64 - cy).abs() * 2.0;
            let dist = dx.max(dy);
            let band = (dist - t * 3.0).rem_euclid(6.0);
            let ch = if band < 1.0 {
                '#'
            } else if band < 2.0 {
                '+'
            } else if band < 3.0 {
                ':'
            } else {
                ' '
            };
            stamp(&mut canvas, x, y, ch);
        }
    }

    to_lines(&canvas)
}

fn plasma(frame: i32, w: i32, h: i32) -> Vec<String> {
    let mut canvas = blank_canvas(w as usize, h as usize);
    let t = frame as f64 * 0.08;

    for y in 0..h {
        for x in 0..w {
            let xf = x as f64;
            let yf = y as f64;
            let v = ((xf * 0.1 + t).sin()
                + (yf * 0.15 + t * 1.3).sin()
                + ((xf + yf) * 0.08 + t * 0.7).sin()
                + (((w as f64 - xf) * 0.1 - yf * 0.05) * 0.05 + t * 1.5).sin())
                / 4.0;
            let brightness = (v * 0.5 + 0.5).clamp(0.0, 1.0);
            stamp(&mut canvas, x, y, ramp_char(brightness));
        }
    }

    to_lines(&canvas)
}
