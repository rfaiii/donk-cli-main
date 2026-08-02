/// Classic doom fire effect for standby/background mode
#[derive(Debug, Clone)]
pub struct DoomFire {
    width: usize,
    height: usize,
    fire_pixels: Vec<u8>,
    last_tick: std::time::Instant,
    tick_interval: std::time::Duration,
}

impl Default for DoomFire {
    fn default() -> Self {
        Self::new()
    }
}

impl DoomFire {
    pub fn new() -> Self {
        let mut fire = Self {
            width: 80,
            height: 20,
            fire_pixels: vec![0; 80 * 20],
            last_tick: std::time::Instant::now(),
            tick_interval: std::time::Duration::from_millis(50),
        };
        fire.seed_bottom();
        fire
    }

    pub fn set_size(&mut self, width: usize, height: usize) {
        self.width = width.max(1);
        self.height = height.max(1);
        self.fire_pixels = vec![0; self.width * self.height];
        self.seed_bottom();
    }

    fn seed_bottom(&mut self) {
        let y = self.height.saturating_sub(1);
        for x in 0..self.width {
            self.fire_pixels[y * self.width + x] = 36;
        }
    }

    pub fn update(&mut self) -> bool {
        if self.last_tick.elapsed() < self.tick_interval {
            return false;
        }
        self.last_tick = std::time::Instant::now();
        self.seed_bottom();

        for x in 0..self.width {
            for y in 1..self.height {
                let src = y * self.width + x;
                let pixel = self.fire_pixels[src];
                if pixel == 0 {
                    self.fire_pixels[src - self.width] = 0;
                } else {
                    let rand = (rand::random::<u32>() % 3) as usize;
                    let dst = src
                        .saturating_sub(self.width)
                        .saturating_sub(rand)
                        .saturating_add(1);
                    if dst < self.fire_pixels.len() {
                        self.fire_pixels[dst] = pixel.saturating_sub((rand & 1) as u8);
                    }
                }
            }
        }
        true
    }

    pub fn view(&self) -> String {
        let mut out = String::with_capacity(self.width * self.height * 8);
        for y in 0..self.height {
            for x in 0..self.width {
                let intensity = self.fire_pixels[y * self.width + x];
                let (r, g, b) = fire_color(intensity);
                out.push_str(&format!("\x1b[48;2;{};{};{}m ", r, g, b));
            }
            out.push_str("\x1b[0m");
            out.push('\n');
        }
        out
    }

    /// Ratatui half-block lines for embedding in a Frame.
    pub fn view_lines(&self) -> Vec<ratatui::text::Line<'static>> {
        use ratatui::prelude::*;
        let mut lines = Vec::with_capacity(self.height);
        for y in 0..self.height {
            let mut spans = Vec::with_capacity(self.width);
            for x in 0..self.width {
                let intensity = self.fire_pixels[y * self.width + x];
                let (r, g, b) = fire_color(intensity);
                spans.push(Span::styled(
                    "█",
                    Style::default().fg(Color::Rgb(r, g, b)),
                ));
            }
            lines.push(Line::from(spans));
        }
        lines
    }
}

fn fire_color(intensity: u8) -> (u8, u8, u8) {
    let i = intensity as u32 * 3;
    match i {
        0..=84 => (0, 0, 0),
        85..=170 => (i as u8, 0, 0),
        171..=255 => (255, (i - 170) as u8, 0),
        _ => (255, 255, (i - 255) as u8),
    }
}
