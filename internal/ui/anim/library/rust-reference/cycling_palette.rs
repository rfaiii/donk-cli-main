//! Cycling animation palette — CRUSH gradient metadata.

pub const CYCLING_RAMP: &[&str] = &[
    "#F967DC", "#F967DC", "#D967F9", "#A967F9", "#6B50FF",
    "#6B50FF", "#A967F9", "#D967F9", "#F967DC",
];

pub fn cycling_gradient(len: usize) -> Vec<&'static str> {
    let mut out = Vec::new();
    for i in 0..len {
        let idx = i % CYCLING_RAMP.len();
        out.push(CYCLING_RAMP[idx]);
    }
    out
}
