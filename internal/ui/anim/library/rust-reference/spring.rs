//! Harmonica-style spring physics (Ryan Juckett / Charm harmonica DNA).

use ratatui::prelude::*;

#[derive(Debug, Clone)]
pub struct Spring {
    pub pos: f64,
    pub vel: f64,
    pub target: f64,
    pub angular_frequency: f64,
    pub damping_ratio: f64,
    pub width: usize,
    pub frame: u64,
    pos_pos: f64,
    pos_vel: f64,
    vel_pos: f64,
    vel_vel: f64,
}

impl Default for Spring {
    fn default() -> Self {
        Self::new()
    }
}

impl Spring {
    pub fn new() -> Self {
        let mut s = Self {
            pos: 0.0,
            vel: 0.0,
            target: 1.0,
            angular_frequency: 6.0,
            damping_ratio: 0.45,
            width: 48,
            frame: 0,
            pos_pos: 1.0,
            pos_vel: 0.0,
            vel_pos: 0.0,
            vel_vel: 1.0,
        };
        s.recompute_coeffs(1.0 / 60.0);
        s
    }

    pub fn set_width(&mut self, width: usize) {
        self.width = width.max(16);
    }

    fn recompute_coeffs(&mut self, dt: f64) {
        let eps = f64::EPSILON * 8.0;
        let omega = self.angular_frequency.max(0.0);
        let zeta = self.damping_ratio.max(0.0);
        if omega < eps {
            self.pos_pos = 1.0;
            self.pos_vel = 0.0;
            self.vel_pos = 0.0;
            self.vel_vel = 1.0;
            return;
        }
        if zeta > 1.0 + eps {
            let za = -omega * zeta;
            let zb = omega * (zeta * zeta - 1.0).sqrt();
            let z1 = za - zb;
            let z2 = za + zb;
            let e1 = (z1 * dt).exp();
            let e2 = (z2 * dt).exp();
            let inv = 1.0 / (2.0 * zb);
            let e1_h = e1 * inv;
            let e2_h = e2 * inv;
            self.pos_pos = e1_h * z2 - z2 * e2_h + e2;
            self.pos_vel = -e1_h + e2_h;
            self.vel_pos = (z1 * e1_h - z2 * e2_h + e2) * z2;
            self.vel_vel = -z1 * e1_h + z2 * e2_h;        } else if zeta < 1.0 - eps {
            let omega_zeta = omega * zeta;
            let alpha = omega * (1.0 - zeta * zeta).sqrt();
            let exp_term = (-omega_zeta * dt).exp();
            let cos_term = (alpha * dt).cos();
            let sin_term = (alpha * dt).sin();
            let inv_alpha = 1.0 / alpha;
            let exp_sin = exp_term * sin_term;
            let exp_cos = exp_term * cos_term;
            let exp_oz_sin = exp_term * omega_zeta * sin_term * inv_alpha;
            self.pos_pos = exp_cos + exp_oz_sin;
            self.pos_vel = exp_sin * inv_alpha;
            self.vel_pos = -exp_sin * alpha - omega_zeta * exp_oz_sin;
            self.vel_vel = exp_cos - exp_oz_sin;
        } else {
            let exp_term = (-omega * dt).exp();
            let time_exp = dt * exp_term;
            let time_exp_freq = time_exp * omega;
            self.pos_pos = time_exp_freq + exp_term;
            self.pos_vel = time_exp;
            self.vel_pos = -omega * time_exp_freq;
            self.vel_vel = -time_exp_freq + exp_term;
        }
    }

    pub fn poke(&mut self) {
        self.target = if self.target > 0.5 { 0.0 } else { 1.0 };
        self.vel += if self.target > 0.5 { 2.0 } else { -2.0 };
    }

    pub fn update(&mut self) {
        self.recompute_coeffs(1.0 / 60.0);
        let old_pos = self.pos - self.target;
        let old_vel = self.vel;
        self.pos = old_pos * self.pos_pos + old_vel * self.pos_vel + self.target;
        self.vel = old_pos * self.vel_pos + old_vel * self.vel_vel;
        self.frame = self.frame.wrapping_add(1);
    }

    pub fn view_lines(&self) -> Vec<Line<'static>> {
        let w = self.width;
        let x = ((self.pos.clamp(0.0, 1.0)) * (w.saturating_sub(1) as f64)).round() as usize;
        let mut track = vec!['·'; w];
        if x < track.len() {
            track[x] = '●';
        }
        let track: String = track.into_iter().collect();
        let bar_fill = ((self.pos.clamp(0.0, 1.0)) * 20.0).round() as usize;
        let bar = format!("[{}{}]", "█".repeat(bar_fill), "░".repeat(20 - bar_fill));

        vec![
            Line::from(""),
            Line::from(Span::styled(
                "harmonica spring".to_string(),
                Style::default().add_modifier(Modifier::BOLD),
            )),
            Line::from(""),
            Line::from(track),
            Line::from(format!("  {bar}  pos={:.2} vel={:.2}", self.pos, self.vel)),
            Line::from(""),
            Line::from(format!("frame {} · space poke · tab next", self.frame)),
        ]
    }

    pub fn view(&self) -> String {
        self.view_lines()
            .into_iter()
            .map(|l| {
                l.spans
                    .iter()
                    .map(|s| s.content.as_ref())
                    .collect::<String>()
            })
            .collect::<Vec<_>>()
            .join("\n")
    }
}
