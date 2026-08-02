# Crush classic chrome — visual target (from layout screens)

**Primary screen:** `ref/cli-design/DONK-CRUSH-LAYOUT.png`  
(also mirrored: `ref/design/cli/DONK-CRUSH-LAYOUT.png`)

Companion logo stills: `ref/cli-design/donk-ascii1.png` … `donk-ascii4.png`  
Anim engineering PDF: [terminal_animation_guide.pdf](terminal_animation_guide.pdf)

This is a **layout + chrome pattern** target. Default ship theme remains `donk-theme` pink/purple/cyan ([LAYOUT-FLOORPLAN.md](LAYOUT-FLOORPLAN.md)). Neon green here maps to a future **`crush-neon`** / classic theme pack — **do not** replace the default palette until that pack ships.

---

## Stack (maps to spaces)

```
S1  Giant DONK-CLI block wordmark (neon green on dark texture)
S3  Status slash bar:  donk-cli  ////////////////  ~ • 0% • ctrl+d open
S7  Message cards (divider-separated interaction blocks)
```

| Region | Pattern | DONK space |
|--------|---------|------------|
| Giant wordmark | Block FIGlet / `ClearVariant`–class mark | **S1** |
| Status slash filler | `name` + `/` run + cwd · pct · hint | **S3** (or S10 if folded) |
| User turn | `\|` prefix + prompt text | **S7** user line |
| Error badge | Solid **blue pill** `ERROR` + title | **S7** system/error card |
| Detail lines | Indented green mono under title | **S7** card body |
| Dividers | Full-width hairlines between turns | **S7** separators |

---

## Interaction card recipe (S7)

```
| <user text>
[ERROR] <Short Title>
  <detail line 1>
  <detail line 2>
────────────────────────
```

Rules:

1. User lines start with `|` (Crush-style), not only `>` — Gemini Layout A may keep `>` until a `/chrome crush` mode lands.
2. Status labels (`ERROR`, later `OK` / `WARN`) are **filled pills**, never plain colored text alone.
3. Error pill fill = high-contrast blue; title + details follow active theme greens (or error_fg for severity).
4. One card per turn; hairline divider between cards.
5. No empty chrome shells (B4).

---

## Logo stills (ascii1–4)

| File | Read as |
|------|---------|
| `donk-ascii1.png` | Larger regular / easy-read block DONK-CLI |
| `donk-ascii2.png` | Alt weight / spacing |
| `donk-ascii3.png` | Alt composition |
| `donk-ascii4.png` | Tight neon block mark (matches Crush header energy) |

Prefer matching these to existing `AsciiLogo` variants (`ClearVariant`, `LargeClean`, `EasyRead`, `Generic`) before inventing new enum cases — see [ASCII-LOGOS.md](ASCII-LOGOS.md).

---

## Theme pack note

| Pack | When |
|------|------|
| `donk-dark` (default) | Pink/purple/cyan chrome today |
| `crush-neon` (future) | Green-on-black + blue ERROR pills from this screen |

Until `crush-neon` exists, use this file for **structure** (slash bar, `|` turns, pills, dividers), not for forcing neon green into the default theme.
