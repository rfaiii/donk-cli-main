//! Charm Bubbles–style spinner frames (port of presets, not a Go runtime).

use ratatui::prelude::*;
use std::time::{Duration, Instant};

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum SpinnerKind {
    Line,
    Dot,
    MiniDot,
    Jump,
    Pulse,
    Points,
    Moon,
}

impl SpinnerKind {
    pub fn all() -> &'static [SpinnerKind] {
        &[
            Self::Line,
            Self::Dot,
            Self::MiniDot,
            Self::Jump,
            Self::Pulse,
            Self::Points,
            Self::Moon,
        ]
    }

    pub fn next(self) -> Self {
        let all = Self::all();
        let i = all.iter().position(|k| *k == self).unwrap_or(0);
        all[(i + 1) % all.len()]
    }

    pub fn label(self) -> &'static str {
        match self {
            Self::Line => "line",
            Self::Dot => "dot",
            Self::MiniDot => "minidot",
            Self::Jump => "jump",
            Self::Pulse => "pulse",
            Self::Points => "points",
            Self::Moon => "moon",
        }
    }

    fn frames(self) -> &'static [&'static str] {
        match self {
            Self::Line => &["|", "/", "-", "\\"],
            Self::Dot => &["⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"],
            Self::MiniDot => &["⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"],
            Self::Jump => &["⢄", "⢂", "⢁", "⡁", "⡈", "⡐", "⡠"],
            Self::Pulse => &["▏", "▎", "▍", "▌", "▋", "▊", "▉", "▊", "▋", "▌", "▍", "▎"],
            Self::Points => &["∙∙∙", "●∙∙", "∙●∙", "∙∙●"],
            Self::Moon => &["🌑", "🌒", "🌓", "🌔", "🌕", "🌖", "🌗", "🌘"],
        }
    }

    fn interval(self) -> Duration {
        match self {
            Self::Line => Duration::from_millis(100),
            Self::Dot | Self::MiniDot => Duration::from_millis(80),
            Self::Jump => Duration::from_millis(100),
            Self::Pulse => Duration::from_millis(120),
            Self::Points => Duration::from_millis(160),
            Self::Moon => Duration::from_millis(200),
        }
    }
}

#[derive(Debug, Clone)]
pub struct Spinner {
    kind: SpinnerKind,
    index: usize,
    last: Instant,
}

impl Default for Spinner {
    fn default() -> Self {
        Self::new()
    }
}

impl Spinner {
    pub fn new() -> Self {
        Self::with_kind(SpinnerKind::Dot)
    }

    pub fn with_kind(kind: SpinnerKind) -> Self {
        Self {
            kind,
            index: 0,
            last: Instant::now(),
        }
    }

    pub fn kind(&self) -> SpinnerKind {
        self.kind
    }

    pub fn set_kind(&mut self, kind: SpinnerKind) {
        self.kind = kind;
        self.index = 0;
        self.last = Instant::now();
    }

    pub fn cycle_kind(&mut self) {
        self.set_kind(self.kind.next());
    }

    pub fn update(&mut self) {
        if self.last.elapsed() >= self.kind.interval() {
            let frames = self.kind.frames();
            self.index = (self.index + 1) % frames.len();
            self.last = Instant::now();
        }
    }

    pub fn frame(&self) -> &'static str {
        let frames = self.kind.frames();
        frames[self.index % frames.len()]
    }

    pub fn view_lines(&self) -> Vec<Line<'static>> {
        let preview: String = self
            .kind
            .frames()
            .iter()
            .map(|f| format!("{f} "))
            .collect();
        vec![
            Line::from(Span::styled(
                format!("spinner · {}", self.kind.label()),
                Style::default().add_modifier(Modifier::BOLD),
            )),
            Line::from(""),
            Line::from(Span::styled(
                format!("  {}  ", self.frame()),
                Style::default()
                    .fg(Color::Rgb(0xF9, 0x67, 0xDC))
                    .add_modifier(Modifier::BOLD),
            )),
            Line::from(""),
            Line::from(Span::styled(
                preview,
                Style::default().fg(Color::DarkGray),
            )),
            Line::from(""),
            Line::from("space cycle preset · tab next"),
        ]
    }
}
