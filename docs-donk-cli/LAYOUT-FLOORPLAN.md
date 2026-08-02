# DONK floor plan — layout, theme, behavior (MUST FOLLOW)

**This document is the binding floor plan.**  
If code and this file disagree, **change the code** (or open a doc PR that explicitly revises a rule).  
Do not invent ad-hoc regions, colors, or fold order outside these rules.

Companion indexes:

- Tile IDs & shrink map: [LAYOUT-TILES.md](LAYOUT-TILES.md)
- ASCII logos (S1/S10/boot): [ASCII-LOGOS.md](ASCII-LOGOS.md)
- Anim hunt / ports: [ANIMATIONS.md](ANIMATIONS.md)
- Ref library map: [`ref/README.md`](../ref/README.md)
- Interactive diagram: Cursor canvas `donk-layout-tiles.canvas.tsx`

---

## 1. Product DNA (locked)

| Layer | Choice | Rule |
|-------|--------|------|
| Default chrome | **Layout A** (Gemini single-column) | New-user default; chat owns `Min(0)` |
| Pro chrome | **Layout B** (DESIGN-BOXES) | Opt-in via `/layout pro` (C8) when built |
| Ship palette | **`donk-theme` pink/purple/cyan/green** | Wireframe neon greens are **partition cues only** |
| Render stack | ratatui + `donk-tui`/`chrome.rs` | ANSI/Unicode only — **no background bitmaps** |
| Tools | Slash overlays (W8) | No always-on 28-col Crush sidebar in Layout A |
| Animations | Port → `donk-anim` (Rust) | Hunt Go/Python/GLSL; **never** ship extra runtimes |

---

## 2. Theme floor plan

### 2.1 Required tokens (`DonkTheme`)

| Token | Role in chrome |
|-------|----------------|
| `background` / `foreground` | Base canvas / body text |
| `donk_purple` | Brand wordmarks, assistant label, S2/S3 accents |
| `donk_pink` | Focus border, Crush scramble, tool-card strings, keys art |
| `donk_green` | Online / success / tool-card borders / cwd |
| `donk_blue` | User `>` prompt, mode chip, file-focus |
| `border` / `border_focused` | B1 / B6 |
| `user_fg` / `assistant_fg` / `system_fg` / `error_fg` | Chat roles |
| `header_*` / `footer_*` | Optional strip fills (may equal background) |

### 2.2 Theme rules (MUST)

1. **One active theme** at a time (`ThemeRegistry`).
2. **Error badges** use `error_fg` (or high-contrast fill + light text) — never pink for errors.
3. **Busy / Crush scramble** uses pink→purple ramp (B5) — not neon matrix green unless a future `neon` theme pack is selected.
4. **Tool cards**: green border + pink string accents (Gemini mockup) mapped through theme tokens.
5. **Neon / classic comps** under `ref/design/classic` and `ref/mood/neon` are **R&D only** until promoted to a named theme TOML.

### 2.3 Theme packs (roadmap)

| Pack | Source | Status |
|------|--------|--------|
| `donk-dark` | `crates/donk-theme/resources/themes/` | **Ship default** |
| Vault CLI themes | `ref/themes/cli-themes/` | Import via `/themes` |
| DONK-E master | `ref/themes/master/` | Tokens / docs — not auto-loaded |
| Neon green | `ref/mood/neon/` + classic comps | Future optional theme |

---

## 3. Space floor plan (S1–S10)

Every visible chrome region has exactly one ID. Do not invent S11+ without updating this file.

| ID | Name | Owns | Default height | Empty behavior (B4) |
|----|------|------|----------------|---------------------|
| **S1** | Brand / ASCII | Wordmark + donkey cluster | Tall at boot; may fold after first user msg | Never 0 at cold start |
| **S2** | Settings & Commands | Tips / slash chips | Tips until first user msg; else 0–1 | 0h when tips hidden |
| **S3** | AI Model + Connection | Model · net · Crush Responding | 1 idle / 2 busy | Never 0 |
| **S4** | Interactive Keys / Loading | Hotkey ribbon / LOADING>>> | 0 idle chat; show on tool/busy | 0h when idle |
| **S5** | File Focus | `$USER CURRENT FILE…` | 0 until path hot | 0h when no focus |
| **S6** | Prompt Label | TYPE YOUR MESSAGE | Prefer merge into input title | 0h if merged |
| **S7** | User Chatbox | Messages + input (`Min(0)`) | Claims leftover | Never 0 |
| **S8** | CLI Viewport | Live PTY | Hidden in A; overlay W8 | 0w/0h until tool/pro |
| **S9** | Node Connection | Peers / sync | Hidden in A | 0w until nodes/pro |
| **S10** | Status / Keys | cwd · sandbox · mode · keyboard ASCII | Always ≥1 (status) | Keys may fold |

### 3.1 Layout A stack (top → bottom) — DEFAULT

```
S1 Brand
S2 Tips/Commands   (B4 → 0 after first user message)
S3 Meta (+ Crush when busy)
S4 Keys            (B4 → 0 when idle)
S5 File focus      (B4 → 0 when no path)
S6 Label           (usually merged into input)
S7 Chat body Min(0) + input
S10 Status + keyboard footer
```

Tools (`/terminal`, `/files`, …) **replace S7 body only** via **W8 Overlay**. S1/S3/S10 stay.

### 3.2 Layout B stack (pro) — FUTURE

```
S1a | S1b     (W2 split)
S8 CLI        (W3 Expand mid)
S2 / S3 / S4  (thin strips)
S5a | S5b     (W2 split)
S7 | S9       (W2 chat | node)
S10 status
```

---

## 4. Behavior codes (authoritative)

### 4.1 Window (W) — how tiles change size/presence

| Code | Name | Contract |
|------|------|----------|
| **W1** Fold | Height → 0 or 1 line; sibling `Min` grows |
| **W2** Halve | Axis ÷ 2; sibling claims freed half |
| **W3** Expand | Claims `Min(0)`; neighbors W1/W2 |
| **W4** Stretch | Fills remainder after constraints resolve |
| **W5** Move | Reorder vertical index |
| **W6** Fly-in | Appear from edge 120–220ms |
| **W7** Fly-out | Dismiss; Stretch reclaim |
| **W8** Overlay | Full-bleed over S7 (Esc restores) |
| **W9** Persist | Save layout mode in config |

**Shrink rule:** If a box is smaller than in the previous layout, it used W1/W2/W5/W7 so another could W3/W4/W6.

### 4.2 Box (B) — how a single tile draws

| Code | Name | Contract |
|------|------|----------|
| **B1** Border | Theme `border`; focused → `border_focused` |
| **B2** Title | Optional border title |
| **B3** Pad | 1 col horizontal; 0–1 row vertical |
| **B4** Empty | No data → **height 0** (no empty chrome shells) |
| **B5** Busy | S3 owns Crush scramble; S4 may pulse |
| **B6** Focus | Exactly one interactive focus (input or tool) |
| **B7** Split | S1/S5 may bisect without new IDs |
| **B8** Priority | Short terminal: fold **S2 → S6 → S4 → S10-keys** before S7/S8 |

### 4.3 Command list (C) — how users drive tiles

| Code | Audience | Contract |
|------|----------|----------|
| **C1** New | Tips on; full labels; tools = W8; Esc exits |
| **C2** Pro | Dense chips; S7\|S9; persistent S8; fold tips |
| **C3** Dock | S2 slash list; Enter runs; Tab cycles |
| **C4** Tool | Open tool → W8/W3; S4 keys; S7 folds |
| **C5** File | Path hot → S5 Expand; optional tool-card in S7 |
| **C6** Nodes | S9 Fly-in; S7 Halve |
| **C7** Busy AI | S3 Expand + scramble; input lock; S4 loading |
| **C8** Mode | `/layout new\|pro\|auto` + W9 |

---

## 5. Chat content rules (inside S7)

| Role | Prefix | Style |
|------|--------|-------|
| User | `>` (Layout A) or `\|` (Crush chrome target) | User text foreground |
| Assistant | `✦` + `DONK` | Purple label + assistant_fg / markdown |
| System | none | system_fg |
| Tool | Bordered card | ✓ title · green border · pink strings |
| Error | Blue **pill** `ERROR` + title | See [CRUSH-CHROME.md](CRUSH-CHROME.md) |

Tool cards are **inline tiles inside S7**, not new S-ids.  
Crush layout screen (`DONK-CRUSH-LAYOUT.png`) defines divider-separated turn cards + slash status bar — structure we converge on; neon green stays optional theme pack.

---

## 6. Enforcement checklist (for agents & PRs)

Before merging chrome/anim changes:

- [ ] New UI belongs to an existing **S#** (or doc updated first)
- [ ] Fold order respects **B8**
- [ ] Empty optional strips use **B4** (no hollow boxes)
- [ ] Colors come from **DonkTheme** tokens (section 2)
- [ ] Busy path uses **C7** + Crush cycling in S3
- [ ] Tools use **W8** in Layout A (do not resurrect sidebar)
- [ ] New effects land in **`donk-anim`** and are listed in [ANIMATIONS.md](ANIMATIONS.md)
- [ ] Design assets referenced from `ref/design/` or `ref/ascii/` (see [ref/README.md](../ref/README.md))

---

## 7. Code map (Layout A today)

| Space | Module |
|-------|--------|
| S1 | `chrome::render_brand_header` |
| S2 | `chrome::render_tips` (partial — tips only) |
| S3 | `chrome::render_meta_strip` + `Cycling` |
| S4 | partial via status help — **gap** |
| S5 | `chrome::render_file_focus` |
| S6 | textarea title/placeholder |
| S7 | `chrome::render_messages` + `textarea` |
| S8/S9 | slash overlays / flows |
| S10 | `render_status_bar` + `render_keyboard_footer` |
| Boot | `donk_anim::BootSplash` (not S-stack) |

Gaps vs this floor plan are tracked in [LAYOUT-TILES.md](LAYOUT-TILES.md) roadmap — implement toward this file, not around it.

---

## 8. Design sources

| Asset | Path |
|-------|------|
| DESIGN-04 | `ref/cli-design/DONK-CLI-DESIGN-04.png` |
| BOXES pro | `ref/cli-design/DONK-DESIGN-BOXES.png` |
| **Crush layout screen** | `ref/cli-design/DONK-CRUSH-LAYOUT.png` |
| Crush chrome patterns | [CRUSH-CHROME.md](CRUSH-CHROME.md) |
| Logo stills | `ref/cli-design/donk-ascii1.png` … `donk-ascii4.png` |
| Gemini/Crush mockups | `ref/cli-design/DONK-CLI-DESIGN-0*.png` |
| Crush anim video | `ref/cli-design/crush-cli-animations.mov` |
| Anim engineering PDF | [terminal_animation_guide.pdf](terminal_animation_guide.pdf) |
| ASCII bible | `ref/ascii/DONK-ASCII-TXT.txt` |

Canonical paths also under `ref/design/cli/` when mirrored. Legacy alias: `ref/cli-design/`.
