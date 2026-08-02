# DONK animations — hunt catalog & port rules

## Binding rule

**Ship only in Rust** (`donk-anim` + ratatui).  
Hunt algorithms anywhere (Go / Rust crates / Python / GLSL) → rewrite as `(state, dt, theme) → Lines`.

Floor plan for **where** anims appear: [LAYOUT-FLOORPLAN.md](LAYOUT-FLOORPLAN.md) (S3 busy, boot reel, `/animations`).  
Crush chrome screen patterns: [CRUSH-CHROME.md](CRUSH-CHROME.md).  
Engineering guide (PDF): [terminal_animation_guide.pdf](terminal_animation_guide.pdf).

---

## Terminal animation guide (PDF) — principles we follow

Source: `docs/terminal_animation_guide.pdf` (also mirrored from Downloads).

| Principle | DONK application |
|-----------|------------------|
| Frame buffer / tick loop (50–100ms) | `donk-anim` update + ratatui redraw |
| Diff/flush, no scroll spam | Overwrite canvas regions only |
| Go Bubble Tea / Lip Gloss = best Charm DNA | Hunt Go → port Rust |
| Rust Ratatui = ship runtime | `donk-anim` + `donk-tui` |
| Braille / high-res ASCII | Future particles / spring trails |
| TrueColor gradients | Splash + Crush `Cycling` ramps |
| Easing (not linear) | Harmonica / spring for W6/W7 |

Do **not** add Python Textual/Rich as a ship dependency — inspiration only.

---

## Language map

| Role | Language | Why |
|------|----------|-----|
| **Ship** | Rust (`donk-anim`) | App runtime |
| **Best Crush / Charm DNA** | Go | Crush + Bubble Tea |
| **Shader-like FX eval** | Rust crates (`tachyonfx`, …) | Only if ANSI-friendly |
| **Host cursor trails** | GLSL | Ghostty — not ratatui |
| **Inspiration only** | Python / C / demos | Steal math, don’t ship interpreters |

---

## Shipped in DONK (`donk-anim`) — pinned

| Effect | Module | Product slot | Status |
|--------|--------|--------------|--------|
| Boot logo reel | `boot.rs` | Pre-S splash (`BootSplash`) | **Pinned** — see [ASCII-LOGOS.md](ASCII-LOGOS.md) |
| Splash gradient | `splash.rs` | Boot underlay + gallery | **Pinned** |
| Doom fire | `doomfire.rs` | Gallery / standby | **Pinned** |
| Crush cycling | `cycling.rs` | **S3 busy** (`Responding`) + gallery | **Pinned** — scramble→label→ellipsis |
| Spinner | `spinner.rs` | **S4 busy** + gallery | **Pinned** |
| Spring | `spring.rs` | Gallery / W6 (harmonica coeffs) | **Pinned** |
| Progress | `progress.rs` | Gallery / long tools | **Pinned** |
| Eyes | `eyes.rs` | Gallery | **Pinned** |
| Matrix rain | `matrix.rs` | Gallery / standby | **Pinned** |
| Space | `space.rs` | Gallery standby | **Pinned** |
| Vanish | `vanish.rs` | Gallery / logo dissolve | **Pinned** |
| Logos | `donk-assets` | Gallery · logos | **Pinned** catalog |
| Scenes (Mission Control) | `mission.rs` | Gallery / standby | **Pinned** — 8 procedural scenes |

`/animations` Tab order (locked):  
`boot → splash → doomfire → cycling → spinner → spring → progress → eyes → matrix → space → vanish → logos → scenes`

### Mission Control Pro scenes (`mission.rs`, shipped v2.6.0)

Ported from Mission Control Pro. 8 canvas-based procedural scenes:

| Scene | Algorithm | Use |
|-------|-----------|-----|
| Fireworks | Expanding bursts + particle rings | Celebration |
| Horizon | Sinusoidal sun + wave horizon + stars | Standby |
| Radio | Concentric expanding rings | Gallery |
| Rain | Phase-offset falling streaks + splash | Standby |
| Starfield | Golden angle warp-drive distribution | Gallery / boot |
| Galaxy | 4-arm spiral with radial glow | Standby |
| Tunnel | Chebyshev distance concentric bands | Gallery |
| Plasma | 4-term sine wave interference | Gallery / standby |

Brightness ramp: ` .,:;irsXA253hMHGS#9B&@` (24 chars).

### Live chrome wiring (do not drift)

| Trigger | Animation | Code |
|---------|-----------|------|
| Cold start | `BootSplash` (reel + corner hero + gradient) | `app.rs` boot |
| AI busy | `Cycling` in **S3** + `Spinner` in **S4** | `busy_anim` / `busy_spinner` |
| `/animations` | Full gallery cycle | `tools.rs` `AnimKind` |
| Tool success (future) | confetti / spring poke | hunt queue |

Logo art rules: **[ASCII-LOGOS.md](ASCII-LOGOS.md)** (MUST FOLLOW).

---

## In-tree references (already vendored under `resources/` / `ref/`)

**Full Crush/Charm → DONK map:** [RESOURCES.md](../RESOURCES.md) (fancy UI queue lives there).

| Source | Path | Port next |
|--------|------|-----------|
| Crush anim | `ref/crush-anim/` · `ref/anim/crush/` | Fidelity vs `.mov` |
| Crush video | `ref/design/cli/crush-cli-animations.mov` | Timing / ellipsis feel |
| Bubble Tea eyes/fire/splash/vanish/space | `resources/bubbletea/examples/` | P1 chrome examples in RESOURCES |
| Bubbles spinners/progress | Charm upstream + BT examples | S4 ribbon / C7 (**spinner/progress shipped**) |
| Harmonica | `resources/harmonica/` | W6/W7 motion |
| Huh spinners | `resources/huh/spinner/` | Setup busy |
| Lip Gloss | `resources/lipgloss/` | Gradient ramps · pills · lists |
| Ghostty shaders | `ref/shaders/` | Host only |
| VHS | `resources/vhs/` | Record demos, not runtime |

Bubble Tea P1 for chrome (not anim): `list-fancy`, `help`, `tabs`, `glamour`, `package-manager`, `autocomplete` — see RESOURCES §3.

---

## External hunt — high value (2025–2026)

### Priority ports (algorithms → `donk-anim`)

| Project | Lang | Why it matters | Candidate DONK use |
|---------|------|----------------|--------------------|
| [tachyonfx](https://github.com/ratatui/tachyonfx) | Rust | **50+** composable cell FX (dissolve, sweep, explode, coalesce) + FTL editor | W6/W7 transitions · S1 dissolve |
| [termflix](https://github.com/paulrobello/termflix) | Rust | ~60 procedural FX (plasma, aurora, rain, starfield, lightning, …) | Gallery + standby underlays |
| [matrix-rain](https://crates.io/crates/matrix-rain) | Rust | Drop-in ratatui Matrix widget (themes/charsets) | Deepen `matrix` or standby |
| [tui-rain](https://crates.io/crates/tui-rain) | Rust | Rain / snow / “emoji soup” | Gallery weather FX |
| [tui-shimmer](https://github.com/vinhnx/tui-shimmer) | Rust | Shimmer text | S3 loading / skeleton |
| [tui-skeleton](https://crates.io/crates/tui-skeleton) | Rust | Pulse/sweep placeholders | Busy chat placeholders |
| [tui-big-text](https://crates.io/crates/tui-big-text) | Rust | Huge FIGlet-ish cells | S1 brand alternate |
| [confetty_rs](https://github.com/Handfish/confetty_rs) | Rust | Particles / fireworks / stars | Celebrate tool success |
| [theattyr](https://github.com/orhun/theattyr) | Rust | VT100 art player | Boot / gallery sequences |
| [sigye](https://github.com/am2rican5/sigye) | Rust | Clock + animated backgrounds | Standby mode |
| [sysc-Go](https://github.com/Nomadcxx/sysc-Go) | Go | Fire/matrix/rain/fireworks + **text effects** (pour, blackhole, fire-text) | Logo reveal flair |
| [erik-adelbert/flame](https://github.com/erik-adelbert/flame) | Go + Bubble Tea | High-perf doom fire | Deepen `doomfire` |
| [erik-adelbert/firework](https://github.com/erik-adelbert/firework) | Go + Bubble Tea | Particle fireworks | Celebrate / vanish outro |
| [charmbracelet/bubbles](https://github.com/charmbracelet/bubbles) | Go | Spinner frames + harmonica progress | S4 / C7 |
| [ratatui-cheese](https://crates.io/crates/ratatui-cheese) | Rust | Bubbletea-inspired spinner/help | Chrome polish |
| [awesome-ratatui](https://github.com/ratatui/awesome-ratatui) | list | Index of TUI/anim crates | Ongoing hunt |

### Text / logo flair (sysc-Go style — great for DONK ASCII)

Steal these *behaviors* for brand marks:

- **Pour** — glyphs fall into logo positions  
- **Print** — typewriter wordmark  
- **Fire-text** — logo consumed by fire  
- **Matrix-art** — streams through glyph mask  
- **Blackhole / ring** — dramatic boot/outro  

Pair with `AsciiLogo::boot_reel()` / vanish.

---

## Local machines (“~300 repos”)

This workspace’s `Projects/` only has a handful of folders (`donk-cli`, `DONK-NODE`, …) — the larger animation clone set is **not** on this path right now. When you point us at the drive/folder (or mount it), drop notes/symlinks here:

| Local path | Notes |
|------------|-------|
| _(add as you find them)_ | Prefer symlink or copy samples into `ref/anim/inbox/` |

Until then, **in-tree** Charm/Bubble Tea samples under `resources/` plus the external hunt table above are the catalog.

---

## Next port queue (ordered)

1. ~~Bubbles spinner set~~ — **shipped 2.5.3** (S4 + gallery)  
2. ~~Harmonica progress~~ — **shipped 2.5.3**  
3. ~~Bubble Tea `space`~~ — **shipped 2.5.3**  
4. **tachyonfx eval** — dissolve/sweep/explode for W6/W7 (prefer algo port or thin dep)  
5. **termflix** cherry-picks: plasma, rain, starfield, aurora (one at a time)  
6. **sysc-Go text effects** — pour/print for logo reel  
7. **confetty_rs** / particles — tool-card success pop  
8. **matrix-rain** fidelity pass on our `matrix` module  

---

## Mac / Linux handoff

Pinned on Windows hybrid (`phase-0-hybrid`, v2.5.3+):

- Logos + slots → [ASCII-LOGOS.md](ASCII-LOGOS.md)
- Layout/theme → [LAYOUT-FLOORPLAN.md](LAYOUT-FLOORPLAN.md)
- This hunt list stays shared; **ship only Rust ports**

On Mac/Linux next:

1. Pull `phase-0-hybrid` / latest PR — verify boot + `/animations` + S3 busy scramble  
2. Host pack: Alacritty (or Ghostty later) — shaders stay host-only  
3. Prefer ports from queue #1–3 first (spinners / progress / space)  
4. Do **not** fork logo enums per OS 

---

## Design-phase rule (unchanged)

1. Hunt (prefer Go Charm / Rust ratatui FX)  
2. Prototype pure frames  
3. Ship only in `donk-anim`  
4. Wire to S3 / boot / `/animations` per [LAYOUT-FLOORPLAN.md](LAYOUT-FLOORPLAN.md)
