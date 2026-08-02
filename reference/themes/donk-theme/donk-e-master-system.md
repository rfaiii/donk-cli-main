# Donk-E Master Design & Color System

> **Permanent Visual Identity & Theme Architecture**
> *Theme Base:* Blood Moon Alacritty Configuration (`donk-e-theme.toml`) [cite: 1]
> *Aesthetic Archetype:* High-Contrast Cyber-Gothic Dark Mode / Terminal Operator

---

## 1. Core Philosophy & Architecture

The **Donk-E** color system is built around uncompromising dark-mode ergonomics, high-contrast terminal syntax readability, and heavy metal visual intensity. It utilizes an ultra-deep obsidian background (`#060606`) paired with an electric neon-green foreground (`#67ff67`), accented by deep purples, acid yellows, and sharp cyan highlights.

### Global Design Rules
* **Background Supremacy:** Always use `#060606` as the absolute canvas base. Never use pure white or generic grey for primary containers.
* **Text Luminance:** Primary running text or active code must leverage the high-visibility neon green (`#67ff67`) or neutral bright white (`#e6e6e6`).
* **Visual Anchoring:** Use high-contrast purple selections (`#671070`) and glowing magenta cursors (`#9f6eff`) to maintain tactile terminal familiarity across all media formats (documents, web apps, IDE themes).

---

## 2. The Master Palette Reference

### Primary & Window Properties
| Element | Hex / Value | Usage Context |
| :--- | :--- | :--- |
| **Window Background** | `#060606` [cite: 1] | Global app background, canvas base, terminal background |
| **Foreground / Primary Text** | `#67ff67` [cite: 1] | Main text, active prompts, primary brand headers |
| **Window Opacity** | `0.9` [cite: 1] | Glassmorphism, desktop compositor layers |

### Accent, Cursor & Selection
| Element | Hex / Value | Usage Context |
| :--- | :--- | :--- |
| **Primary Cursor** | `#9f6eff` [cite: 1] | Text cursor, focus rings, interactive highlights |
| **VI Mode Cursor** | `#1b1d1e` [cite: 1] | Inactive / block cursor state |
| **Selection Background** | `#671070` [cite: 1] | Text highlight background, active menu states |

### Normal Palette (Base Spectrum)
| Color Name | Hex Code | Terminal / UI Mapping |
| :--- | :--- | :--- |
| **Black** | `#1b1b1b` [cite: 1] | Sub-cards, borders, dark containers |
| **Red** | `#f88a8a` [cite: 1] | Warnings, errors, destructive actions, soft highlights |
| **Green** | `#6af85f` [cite: 1] | Success states, strings, positive indicators |
| **Yellow** | `#fcff6f` [cite: 1] | Warnings, variables, highlighted constants |
| **Blue** | `#00acff` [cite: 1] | Links, informational badges, tags |
| **Magenta** | `#9f6eff` [cite: 1] | Functions, special keywords, cursor highlights |
| **Cyan** | `#67f5ff` [cite: 1] | Operators, types, boolean flags |
| **White** | `#e6e6e6` [cite: 1] | Secondary text, high-legibility body copy |

### Bright Palette (High-Intensity Accents)
| Color Name | Hex Code | Terminal / UI Mapping |
| :--- | :--- | :--- |
| **Bright Black** | `#1d1d1d` [cite: 1] | Borders, dividers, subtle depth shading |
| **Bright Red** | `#dd1010` [cite: 1] | Critical alerts, hard errors, metal blood accents |
| **Bright Green** | `#2fd422` [cite: 1] | High-visibility success, active status dots |
| **Bright Yellow**| `#d5d918` [cite: 1] | Focused metrics, glowing highlights |
| **Bright Blue** | `#0000ff` [cite: 1] | Deep structural highlights, hyperlinks |
| **Bright Magenta**| `#5c19df` [cite: 1] | Deep shadow purple, deep component backgrounds |
| **Bright Cyan** | `#67f5ff` [cite: 1] | Neon highlights, active state glows |
| **Bright White** | `#dcdcdc` [cite: 1] | Pure headings, emphasized white text |

---

## 3. Typography & Environment

* **Font Family:** `JetBrainsMono Nerd Font` [cite: 1]
* **Style:** Bold [cite: 1]
* **Base Size:** `14.0pt` [cite: 1] (Scale up proportionally for headers: H1: 28pt, H2: 20pt, H3: 16pt)
* **Import Base:** `~/.config/alacritty/themes/themes/blood_moon.toml` [cite: 1]
