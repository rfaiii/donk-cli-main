# DONK ASCII fonts — TAAG / FIGlet map (MUST FOLLOW for authoring)

**Bible:** [`ref/ascii/DONK-ASCII-TXT.txt`](../ref/ascii/DONK-ASCII-TXT.txt)  
**Tool:** [patorjk TAAG](https://patorjk.com/software/taag/#p=display&f=Graffiti&t=DONK-CLI&x=none&v=4&h=4&w=80&we=false)  
**Vendored `.flf`:** [`ref/ascii/fonts/`](../ref/ascii/fonts/)  
**Runtime:** still **pre-rendered** strings in `donk-assets` — do **not** add a figlet dependency to the `donk` binary.

Companion: [ASCII-LOGOS.md](ASCII-LOGOS.md)

---

## Bible → font provenance

| Bible # | Sample cue | TAAG / FIGlet font | Local `.flf` | Ship role |
|---------|------------|--------------------|--------------|-----------|
| 1A–1C | Hand donkey | **Custom** | — | Mascot |
| 2 | Box-drawing signature | **Custom** Unicode | — | Signature |
| 3–4 | Half-block / `██` | Custom compositor | — | Compact chrome |
| **5** | `.o8` / `oooo` | **Roman** | `roman.flf` | XL marketing |
| **6** | stick `,-.` | Letters / JS Stick (verify) | — | Vertical-scarce |
| 7–8.2 | Dense ▓▒░ | Custom ANSI-density | — | Intro |
| **17** | technical `┼` | Custom / Cyber* adjacent | `cyberlarge.flf` family | Loading |
| **36** | `\|\|d \|\|\|o \|\|` | **Keyboard** | `keyboard.flf` | S10 footer |
| **37** | `▗▄▖` half-blocks | **ANSI Regular** | `ansi_regular.flf` | Standard layout |
| **38** | `▓█████▄ ▒█████` | **Doom** / Bloody | `doom.flf`, `bloody.flf` | Hype |
| **39** | Compact `▄████` | Custom | — | Generic |
| **40** | Binary | Custom | — | Secret |
| **41** | `e88~\888` / `.d88888` | **Soft** + **ANSI Shadow** | `soft.flf`, `ansi_shadow.flf` | ClearVariant / boot |
| **42** | `█████` / `░░███` | **ANSI Shadow** | `ansi_shadow.flf` | Large clean / Crush header |
| **44** | `░██` shadow | **ANSI Shadow** / Shadow | `ansi_shadow.flf`, `shadow.flf` | Shadow variant |

---

## Must-have pack (vendored)

Priority fonts under `ref/ascii/fonts/`:

| File | Why |
|------|-----|
| `ansi_shadow.flf` | Crush / #42 / S1 hero energy |
| `ansi_regular.flf` | #37-class half-block marks |
| `doom.flf` / `bloody.flf` | #38 bloody |
| `roman.flf` | #5 XL |
| `keyboard.flf` | #36 S10 |
| `soft.flf` | #41 clear variant |
| `graffiti.flf` | TAAG default flair |
| `big.flf` `standard.flf` `slant.flf` `shadow.flf` `block.flf` | Staples |
| `colossal.flf` `larry3d.flf` `banner3-D.flf` `doh.flf` `3-d.flf` | Showy boot |
| `ogre.flf` `crawford.flf` | Heavy / alt |
| `cyberlarge.flf` `cybermedium.flf` `cybersmall.flf` | Tech / loading family |
| `poison.flf` `broadway.flf` `starwars.flf` | Optional gallery spice |

Sources: [xero/figlet-fonts](https://github.com/xero/figlet-fonts) (TAAG-aligned names) + [cmatsuoka/figlet-fonts](https://github.com/cmatsuoka/figlet-fonts).

---

## Regenerate samples

```powershell
.\scripts\gen-ascii-samples.ps1
# requires figlet on PATH (Windows: scoop/choco, or WSL)
```

Writes `ref/ascii/samples/<font>-DONK-CLI.txt`. Diff against bible slices before promoting into `donk-assets`.

### TAAG deep-links

- [ANSI Shadow · DONK-CLI](https://patorjk.com/software/taag/#p=display&f=ANSI%20Shadow&t=DONK-CLI)
- [Doom · DONK-CLI](https://patorjk.com/software/taag/#p=display&f=Doom&t=DONK-CLI)
- [Roman · DONK](https://patorjk.com/software/taag/#p=display&f=Roman&t=DONK)
- [Keyboard · donk-cli](https://patorjk.com/software/taag/#p=display&f=Keyboard&t=donk-cli)
- [Graffiti · DONK](https://patorjk.com/software/taag/#p=display&f=Graffiti&t=DONK)
- [Soft · DONK](https://patorjk.com/software/taag/#p=display&f=Soft&t=DONK)

---

## Change policy

1. Prefer editing art in the bible / TAAG → bake into `donk-assets`
2. New FIGlet variants: add `.flf` here + row in this doc + sample
3. Never fork fonts per OS — one catalog for Windows / Mac / Linux
