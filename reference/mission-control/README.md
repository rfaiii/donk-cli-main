# Mission Control Pro

Terminal-first mission control for the pipeline, with an editor-friendly layout for VS Code and a live dashboard for the terminal.

## What it shows

- Clock and host status
- Weather widget inspired by Dashy-style dashboard cards
- Workspace counts for the legacy pipeline folders
- CPU, memory, disk, and top process monitors
- Git status and local API/service endpoint widgets
- Pipeline progress using the older file-count mission model
- Runtime logs and quick command shortcuts for the workspace
- A big HOME hub with grouped summary boxes and a mini animation window
- Interactive Textual tabs with keyboard navigation and launch buttons
- A dedicated looping ASCII animation stage with scene switching and pause control
- Live-fetched remote animation clips when source URLs are reachable, with local fallback loops when they are not

## Run it

```bash
cd /Users/ravery/Documents/MISSION_CONTROL/mission_control_pro
./start.sh --embedded
```

Use `--embedded` when you want the dashboard to stay inside the current terminal pane, which is usually the nicer mode for VS Code.

Add `--legacy` if you want the older static Rich dashboard instead of the new Textual control room.
Add `--once` if you want a one-shot snapshot for a quick terminal check.

Keyboard shortcuts in the Textual view:

- `h` opens the home hub
- `m` opens the animations tab
- `space` pauses or resumes the animation loop
- `n` moves to the next animation scene
- `b` moves to the previous animation scene
- `o` opens the observability tab
- `1` through `9` launch the configured shortcuts, except on the WALL tab where they pick animation favorites
- the Launchers tab shows the live shortcut monitor so you can see which tools are running or waiting

## Edit the dashboard

- [`dashboard.yaml`](./dashboard.yaml) controls the widget layout, services, shortcuts, and mission seeds.
- [`active_missions.json`](./active_missions.json) is the live mission list the dashboard reads on startup.
