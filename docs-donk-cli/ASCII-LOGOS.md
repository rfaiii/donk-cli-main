# DONK ASCII logos — pinned catalog (MUST FOLLOW)

**Source of truth for art:** `ref/ascii/DONK-ASCII-TXT.txt`  
**Runtime crate:** `crates/donk-assets` (`AsciiLogo`, `LogoSlot`)  
**Do not invent new logo variants in chrome without updating this file + the enum.**

Companion: [ANIMATIONS.md](ANIMATIONS.md) · [LAYOUT-FLOORPLAN.md](LAYOUT-FLOORPLAN.md) · [CRUSH-CHROME.md](CRUSH-CHROME.md) · [ASCII-FONTS.md](ASCII-FONTS.md)

Logo still comps (neon Crush look): `ref/cli-design/donk-ascii1.png` … `donk-ascii4.png` — map to existing enum variants before adding new ones.

**FIGlet / TAAG provenance:** [ASCII-FONTS.md](ASCII-FONTS.md) · fonts in `ref/ascii/fonts/` · regenerate with `scripts/gen-ascii-samples.ps1`.

---

## Locked product wiring (Layout A)

| UI location | Space | `AsciiLogo` | Notes |
|-------------|-------|-------------|-------|
| Brand header wordmark | **S1** | `ClearVariant` → `TopCentered` → `Generic` | Width fallback via `chrome::header_logo_for_width` |
| Brand header mascot | **S1** | `TinyDonkey` | Always paired with wordmark |
| Keyboard footer | **S10** | `Keyboard` | `donk-cli` keycaps |
| Boot corner hero | pre-S | `ClearVariant` (`boot_corner_hero`) | Swaps via `boot_corner_alts` |
| Boot corner alts | pre-S | `ClearVariant`, `LargeClean`, `IntroBlock`, `TopCentered` | Wide-term rotation |
| Boot logo reel | pre-S | `AsciiLogo::boot_reel()` | Order locked below |
| Loading / technical | S3/S4 | `Technical` | Status / load wordmark |
| Gallery · logos | `/animations` | `AsciiLogo::all()` cycle | Tab through every mark |
| Secret easter egg | — | `Binary` | Reel + gallery only |

---

## Full enum (bible IDs)

| Enum | Bible # | Slot | Ship role |
|------|---------|------|-----------|
| `BoldDonkey` | 1A | MainMark | Hero mascot (reel) |
| `RegularDonkey` | 1B | MainMark | Backup if bold too heavy |
| `TinyDonkey` | 1C | CompactMark | **S1 mascot** + boot |
| `Signature` | 2 | Signature | Sign-out / personal |
| `SmallWordmark` | 3 | ChromeAlt | Horizontal-tight |
| `LargeWordmark` | 4 | ChromeAlt | Standard readable WM |
| `ThinWordmark` | 6 | ChromeAlt | Vertical-scarce |
| `IntroBlock` | 7 | Intro | Splash / corner alt |
| `IntroCompact` | 8.2 | Intro | Scrolling shrink of #7 |
| `EasyRead` | 12 | Intro | Multicolor-friendly |
| `TopCentered` | 14 | ChromeAlt | **S1 mid-width** fallback |
| `Technical` | 17 | Loading | Load / technical status |
| `Outline` | 20 | ChromeAlt | Light / outline |
| `Keyboard` | 36 | Keyboard | **S10 footer** |
| `Generic` | 39 | ChromeAlt | **S1 narrow** fallback |
| `Binary` | 40 | Secret | Easter egg |
| `ClearVariant` | 41 | Intro | **S1 preferred** + boot hero |
| `LargeClean` | 42 | Intro | Large clean block |
| `VerySimple` | 43 | ChromeAlt | Minimal mark |

Embedded slice files (when not inline): `crates/donk-assets/assets/logos/*.txt`

---

## Boot reel order (locked)

`AsciiLogo::boot_reel()` — do not reorder casually:

1. `TinyDonkey`
2. `Technical`
3. `BoldDonkey`
4. `EasyRead`
5. `IntroCompact`
6. `Generic`
7. `Outline`
8. `Signature`
9. `LargeWordmark`
10. `Keyboard`
11. `VerySimple`
12. `Binary`

Timing (in `donk-anim::boot`): hold **720ms**, fade **220ms**, corner swap **2600ms**.

---

## Theme rule for logos

- Colorize through `DonkTheme` (`donk_purple` / `donk_pink` / `donk_green`) — never hardcode neon green as default.
- `art_lines()` skips catalog labels like `NO COLOR:`.
- Prefer `art_width()` / `art_height()` before painting into a `Rect`.

---

## Mac / Linux parity checklist

When touching logos on another host:

- [ ] Same `AsciiLogo` enum + `boot_reel()` order
- [ ] UTF-8 terminal (Courier / Nerd Font optional; glyphs must render)
- [ ] Width fallbacks still work in narrow panes (`ClearVariant` → `TopCentered` → `Generic`)
- [ ] `/animations` logos pane still cycles `AsciiLogo::all()`
- [ ] No platform-specific logo forks — one catalog for all OS

---

## Change policy

1. Edit art in `ref/ascii/DONK-ASCII-TXT.txt` (or logo `.txt` slices)
2. Mirror into `crates/donk-assets/assets/`
3. Update this doc + enum/`boot_reel` if adding/removing variants
4. Smoke: `donk` boot splash + `/animations` logos

---

## Rendered FIGlet logos (v2.6.0)

Generated from the 26 vendored `.flf` fonts in `ref/ascii/fonts/`:

| Directory | Files | Content |
|-----------|-------|---------|
| `ref/ascii/logos_rendered/` | 26 | `DONK-CLI` in all 26 fonts |
| `ref/ascii/logos_rendered/short/` | 26 | `DONK` in all 26 fonts |
| `ref/ascii/logos_rendered/underscores/` | 26 | `DONK_CLI` in all 26 fonts |
| `ref/ascii/animated/` | 300 | 8 animation styles (wipe, scanline, glitch, rain, typewriter, dissolve, bounce, curtain) |
| `ref/ascii/animated/donkey_walk/` | 5 | Donkey mascot walking cycle |
| `ref/ascii/animated/tiny_donkey_walk/` | 4 | Tiny donkey walking |
| `ref/ascii/animated/donkey_gradient/` | 5 | Gradient-colored donkey (ANSI truecolor) |
| `ref/ascii/gradient/` | 126 | 7 logos × 6 color presets × 3 directions |
| `ref/ascii/gradient/animated/` | 240 | 10-frame gradient cycling loops |

### Best fonts for DONK brand
- **ansi_shadow** — Block shadow, best for boot splash
- **slant** — Italic, best for chrome header
- **doom** — Bold geometric, best for standby
- **standard** — Classic FIGlet, best default wordmark
- **cyberlarge** — Cyberpunk, best for tech/loading

### Color presets (gradient logos)
- `cyan_purple` — Cyan to purple
- `pink_purple` — Crush colors (#F967DC to #6B50FF)
- `green_cyan` — Green to cyan
- `fire` — Orange to yellow
- `ice` — Blue to ice white
- `donk_brand` — Donk green to blue
