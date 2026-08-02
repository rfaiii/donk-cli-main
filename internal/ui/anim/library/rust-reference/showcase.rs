//! Showcase catalog stub (CRUSH ecosystem reference)

#[derive(Debug, Clone)]
pub struct Showcase;

impl Showcase {
    pub fn new() -> Self {
        Self
    }

    pub fn update(&mut self) {}

    pub fn view(&self) -> &'static str {
        "CRUSH Showcase catalog"
    }
}
