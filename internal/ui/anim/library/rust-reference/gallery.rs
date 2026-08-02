//! Animation gallery metadata — CRUSH ecosystem port reference.

pub const GALLERY_NAMES: &[&str] = &[
    "Harmonica Spring",
    "Gradient Splash",
    "Doom Fire",
    "Blinking Eyes",
    "Char Cycling",
];

pub fn gallery_tab_line(active: usize) -> String {
    GALLERY_NAMES
        .iter()
        .enumerate()
        .map(|(i, n)| if i == active { format!("[{}]", n) } else { n.to_string() })
        .collect::<Vec<_>>()
        .join("  ")
}
