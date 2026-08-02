# Donk-E Multi-Language Color Implementation Guide

> **Developer Specification for Syntax Highlighting & UI Theming**
> *Applies to:* Web (CSS/Tailwind), Python, Rust, JavaScript/TypeScript, Shell/Bash, and Markdown/Documentation.

---

## 1. Semantic Mapping Matrix

To maintain visual consistency across different programming languages and document types, map your syntax and UI elements to the Donk-E palette as follows:

| Semantic Token | Donk-E Palette Role | Hex Code | Visual Intent |
| :--- | :--- | :--- | :--- |
| **Keywords / Control Flow** | Magenta (Normal) | `#9f6eff` [cite: 1] | Stands out clearly against dark background |
| **Strings / Text Literals** | Green (Normal) | `#6af85f` [cite: 1] | High readability, natural contrast |
| **Functions / Methods** | Bright Magenta | `#5c19df` [cite: 1] / `#9f6eff` [cite: 1] | Draws the eye to executable blocks |
| **Variables / Properties** | Yellow (Normal) | `#fcff6f` [cite: 1] | Warm pop for identifiers |
| **Types / Classes** | Cyan (Normal) | `#67f5ff` [cite: 1] | Sharp, technical distinction |
| **Numbers / Constants** | Bright Yellow | `#d5d918` [cite: 1] | Distinct numeric highlighting |
| **Comments / Docstrings** | Black (Bright) | `#1d1d1d` [cite: 1] | Recessed, non-distracting background tone |
| **Errors / Exceptions** | Red (Bright) | `#dd1010` [cite: 1] | Immediate visual warning |

---

## 2. Language-Specific Implementation Blueprints

### A. CSS / Web UI Variables (`:root`)
Inject these custom properties into your web projects, UI templates, or documentation styles:

```css
:root {
  /* Primary Canvas */
  --donk-bg: #060606;
  --donk-fg: #67ff67;
  --donk-surface: #1b1b1b;
  --donk-surface-bright: #1d1d1d;

  /* Accents & Interaction */
  --donk-cursor: #9f6eff;
  --donk-selection-bg: #671070;
  --donk-selection-fg: #e6e6e6;

  /* Syntax & Terminal Spectrum */
  --donk-red: #f88a8a;
  --donk-bright-red: #dd1010;
  --donk-green: #6af85f;
  --donk-bright-green: #2fd422;
  --donk-yellow: #fcff6f;
  --donk-bright-yellow: #d5d918;
  --donk-blue: #00acff;
  --donk-bright-blue: #0000ff;
  --donk-magenta: #9f6eff;
  --donk-bright-magenta: #5c19df;
  --donk-cyan: #67f5ff;
  --donk-white: #e6e6e6;
  --donk-bright-white: #dcdcdc;
  
  --donk-font: 'JetBrainsMono Nerd Font', monospace;
}
```

### B. Python Configuration & Logging Palette
Use this dictionary for terminal outputs, Rich/Click formatting, or GUI application themes in Python:

```python
DONK_PALETTE = {
    "background": "#060606",
    "foreground": "#67ff67",
    "cursor": "#9f6eff",
    "selection": "#671070",
    "syntax": {
        "keyword": "#9f6eff",
        "string": "#6af85f",
        "variable": "#fcff6f",
        "type": "#67f5ff",
        "comment": "#1d1d1d",
        "error": "#dd1010",
        "number": "#d5d918",
    }
}
```

### C. Rust Theme Constants
For custom terminal apps, TUI interfaces, or game engines written in Rust:

```rust
pub struct DonkETheme {
    pub background: &'static str,
    pub foreground: &'static str,
    pub accent_magenta: &'static str,
    pub accent_green: &'static str,
    pub accent_cyan: &'static str,
}

pub const DONK_THEME: DonkETheme = DonkETheme {
    background: "#060606",
    foreground: "#67ff67",
    accent_magenta: "#9f6eff",
    accent_green: "#6af85f",
    accent_cyan: "#67f5ff",
};
```

### D. JavaScript / TypeScript Theme Object
For web apps, VS Code extension themes, or Node.js CLI tools:

```typescript
export const DonkETheme = {
  name: "donk-e-permanent",
  type: "dark",
  colors: {
    "editor.background": "#060606",
    "editor.foreground": "#67ff67",
    "editorCursor.foreground": "#9f6eff",
    "editor.selectionBackground": "#671070",
    "terminal.ansiBlack": "#1b1b1b",
    "terminal.ansiRed": "#f88a8a",
    "terminal.ansiGreen": "#6af85f",
    "terminal.ansiYellow": "#fcff6f",
    "terminal.ansiBlue": "#00acff",
    "terminal.ansiMagenta": "#9f6eff",
    "terminal.ansiCyan": "#67f5ff",
    "terminal.ansiWhite": "#e6e6e6",
  }
};
```

---

## 3. Application Guidelines Across Media

1. **Documents & Reports:** Use `#060606` or ultra-dark slate backgrounds with `#67ff67` headers and `#e6e6e6` body text. Highlight key data callouts with `#671070` containers and `#67f5ff` borders.
2. **Web Pages / Dashboards:** Implement CSS custom properties universally. Ensure contrast ratios on text blocks (use `#e6e6e6` for dense paragraphs instead of neon green to prevent eye strain during long-form reading).
3. **Code Editors & Terminals:** Link directly back to your Alacritty config [cite: 1] or generate corresponding syntax themes using the mapping table above.
