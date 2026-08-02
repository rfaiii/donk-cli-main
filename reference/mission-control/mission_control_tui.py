#!/usr/bin/env python3
from __future__ import annotations

import argparse
import platform
import subprocess
import threading
import time
from pathlib import Path

from rich import box
from rich.console import Console
from rich.panel import Panel
from textual.app import App, ComposeResult
from textual.containers import Horizontal, Vertical, VerticalScroll, Grid
from textual.widgets import Button, Footer, Header, RichLog, Static, TabbedContent, TabPane

import mission_control as core
import mission_control_animations as motion


BASE_DIR = Path(__file__).resolve().parent
console = Console()


class DashboardCard(Static):
    """Focusable card button that renders a live Rich panel."""

    def __init__(self, kind: str, *, id: str) -> None:
        super().__init__(id=id)
        self.kind = kind

    def render(self):
        app = self.app
        if self.kind == "clock":
            return core.get_clock_panel(app.dashboard_config)
        if self.kind == "weather":
            return core.get_weather_panel(app.dashboard_config)
        if self.kind == "workspace":
            return core.get_workspace_panel(app.dashboard_config, app.missions)
        if self.kind == "cpu":
            return core.get_cpu_panel()
        if self.kind == "memory":
            return core.get_memory_panel()
        if self.kind == "processes":
            return core.get_processes_panel()
        if self.kind == "pipeline":
            return core.get_pipeline_panel(app.missions)
        if self.kind == "telemetry":
            return core.get_telemetry_panel(app.missions)
        if self.kind == "services":
            return core.get_services_panel(app.dashboard_config)
        if self.kind == "endpoints":
            return core.get_endpoints_panel(app.dashboard_config)
        if self.kind == "git":
            return core.get_git_status_panel(app.dashboard_config)
        if self.kind == "logs":
            return core.get_logs_panel()
        if self.kind == "shortcuts":
            return core.get_shortcuts_panel(app.dashboard_config)
        if self.kind == "alerts":
            return core.get_alerts_panel(app.dashboard_config, app.missions)
        if self.kind == "health":
            return core.get_health_panel(app.dashboard_config)
        if self.kind == "status":
            current = app.animation_sources[app.animation_source_index].name if app.animation_sources else "local"
            return core.get_status_panel(
                app.dashboard_config,
                app.missions,
                active_tab=app.active_tab,
                animation_source=current,
            )
        if self.kind == "launch_monitor":
            return core.get_launcher_monitor_panel(app.dashboard_config, app.launch_states)
        return core.get_logs_panel()


class AnimationCard(Static):
    """Dedicated looping ASCII animation stage."""

    def render(self):
        app = self.app
        scene = app.animation_sources[app.animation_source_index]
        width = max(40, self.size.width - 4) if self.size.width else 56
        height = max(12, self.size.height - 4) if self.size.height else 16
        return motion.render_animation_source(
            scene,
            app.animation_frame_index,
            paused=app.animation_paused,
            width=width,
            height=height,
        )


class AnimationFavoriteCard(Static):
    """Compact selectable animation favorite."""

    def __init__(self, slot_index: int, *, id: str) -> None:
        super().__init__(id=id)
        self.slot_index = slot_index

    def render(self):
        app = self.app
        if not app.animation_sources or self.slot_index <= 0:
            return Panel(
                core.make_table([("State", "No clips available")]),
                title=f"[b]{self.slot_index}. Favorite[/b]",
                border_style="dim",
                box=box.ROUNDED,
            )

        source_index = self.slot_index - 1
        if source_index >= len(app.animation_sources):
            rows = [
                ("Slot", str(self.slot_index)),
                ("State", "No clip"),
            ]
            return Panel(core.make_table(rows), title="[b]Favorite[/b]", border_style="dim", box=box.ROUNDED)

        scene = app.animation_sources[source_index]
        active = source_index == app.animation_source_index
        rows = [
            ("Slot", str(self.slot_index)),
            ("Clip", scene.title),
            ("Source", scene.source_type.upper()),
            ("Hint", "press this number"),
        ]
        border = "bright_magenta" if active else "cyan"
        return Panel(
            core.make_table(rows),
            title=f"[b]{self.slot_index}. Favorite[/b]",
            border_style=border,
            box=box.ROUNDED,
        )


class MissionControlApp(App[None]):
    """Interactive Mission Control dashboard."""

    CSS = """
    Screen {
        background: #06101c;
        color: #e8eef9;
    }

    Header {
        dock: top;
    }

    Footer {
        dock: bottom;
    }

    TabbedContent {
        height: 1fr;
    }

    TabPane {
        padding: 1;
    }

    VerticalScroll {
        height: 1fr;
    }

    Horizontal {
        height: 1fr;
    }

    Vertical {
        width: 1fr;
        height: 1fr;
    }

    DashboardCard {
        width: 1fr;
        margin: 0 0 1 0;
        min-height: 8;
    }

    AnimationCard {
        width: 1fr;
        height: 1fr;
        min-height: 24;
    }

    AnimationFavoriteCard {
        width: 1fr;
        min-height: 7;
    }

    #home-stack {
        height: 1fr;
    }

    #animation-favorites {
        grid-size: 3 3;
        grid-gutter: 1;
        height: 11;
        margin-top: 1;
    }

    TabPane#wall {
        padding: 0;
    }

    #wall-stage {
        height: 1fr;
        padding: 1;
    }

    #launcher-grid {
        grid-size: 2 2;
        grid-gutter: 1;
        height: auto;
        margin-top: 1;
    }

    #activity-log {
        height: 1fr;
        border: round $accent;
        padding: 1;
    }
    """

    BINDINGS = [
        ("q", "quit", "Quit"),
        ("r", "refresh_dashboard", "Refresh"),
        ("h", "show_tab('home')", "Home"),
        ("o", "show_tab('observability')", "Observability"),
        ("p", "show_tab('pipeline')", "Pipeline"),
        ("m", "show_tab('wall')", "Wall"),
        ("l", "show_tab('logs')", "Logs"),
        ("a", "show_tab('launchers')", "Launchers"),
        ("space", "toggle_animation_pause", "Pause"),
        ("n", "next_animation_scene", "Next scene"),
        ("b", "previous_animation_scene", "Prev scene"),
    ]

    def __init__(self) -> None:
        super().__init__()
        config_path = BASE_DIR / "dashboard.yaml"
        missions_path = BASE_DIR / "active_missions.json"
        self.dashboard_config = core.load_dashboard_config(config_path)
        self.missions_path = missions_path
        self.missions = core.load_missions(missions_path, self.dashboard_config)
        self.shortcuts = list(self.dashboard_config.get("shortcuts", []))
        self.active_tab = "home"
        self.launch_states: dict[int, dict[str, float | str]] = {}
        self.animation_sources = motion.build_local_animation_catalog()
        self.animation_source_index = 0
        self.animation_frame_index = 0
        self.animation_paused = False
        self.animation_scene_ticks = 0
        self._animation_scene_hold = 64
        self._animation_catalog_refreshing = False
        self._last_log_lines = 0

    def compose(self) -> ComposeResult:
        yield Header()
        with TabbedContent(initial="home", id="tabs"):
            with TabPane("HOME", id="home"):
                with VerticalScroll(id="home-stack"):
                    with Horizontal():
                        yield DashboardCard("alerts", id="card-alerts-home")
                        yield DashboardCard("health", id="card-health-home")
                        yield DashboardCard("status", id="card-status-home")
                    yield DashboardCard("launch_monitor", id="card-launch-monitor-home")
                    with Horizontal():
                        with Vertical():
                            yield DashboardCard("clock", id="card-clock-home")
                            yield DashboardCard("weather", id="card-weather-home")
                            yield DashboardCard("workspace", id="card-workspace-home")
                        with Vertical():
                            yield DashboardCard("pipeline", id="card-pipeline-home")
                            yield DashboardCard("telemetry", id="card-telemetry-home")
                            yield DashboardCard("shortcuts", id="card-shortcuts-home")
                        with Vertical():
                            yield DashboardCard("git", id="card-git-home")
                            yield DashboardCard("services", id="card-services-home")
                            yield DashboardCard("endpoints", id="card-endpoints-home")
                    with Horizontal():
                        with Vertical():
                            yield DashboardCard("cpu", id="card-cpu-home")
                            yield DashboardCard("memory", id="card-memory-home")
                        with Vertical():
                            yield DashboardCard("processes", id="card-processes-home")
                            yield DashboardCard("logs", id="card-logs-home")
                        with Vertical():
                            yield AnimationCard(id="card-animation-home")

            with TabPane("Pipeline", id="pipeline"):
                with Horizontal():
                    with Vertical():
                        yield DashboardCard("pipeline", id="card-pipeline")
                        yield DashboardCard("telemetry", id="card-telemetry")
                    with Vertical():
                        yield DashboardCard("services", id="card-services")
                        yield DashboardCard("endpoints", id="card-endpoints")

            with TabPane("Observability", id="observability"):
                with Horizontal():
                    with Vertical():
                        yield DashboardCard("cpu", id="card-cpu")
                        yield DashboardCard("memory", id="card-memory")
                    with Vertical():
                        yield DashboardCard("processes", id="card-processes")
                        yield DashboardCard("git", id="card-git-observe")

            with TabPane("Logs", id="logs"):
                with Horizontal():
                    yield RichLog(id="activity-log", markup=False, highlight=True, wrap=True, auto_scroll=True)
                    with Vertical():
                        yield DashboardCard("services", id="card-services-logs")
                        yield DashboardCard("endpoints", id="card-endpoints-logs")
                        yield DashboardCard("workspace", id="card-workspace-logs")

            with TabPane("WALL", id="wall"):
                with Vertical(id="wall-stage"):
                    yield AnimationCard(id="card-animation-wall")
                    with Grid(id="animation-favorites"):
                        for index in range(1, min(9, len(self.animation_sources)) + 1):
                            yield AnimationFavoriteCard(index, id=f"animation-favorite-{index}")

            with TabPane("Launchers", id="launchers"):
                with VerticalScroll():
                    yield DashboardCard("launch_monitor", id="card-launch-monitor-launchers")
                    with Grid(id="launcher-grid"):
                        for index, shortcut in enumerate(self.shortcuts, start=1):
                            label = shortcut.get("label", f"Launcher {index}")
                            command = shortcut.get("command", "")
                            note = shortcut.get("note", "")
                            button = Button(f"{index}. {label}", id=f"launcher-{index}", variant="success")
                            button.tooltip = f"{command}\n{note}".strip()
                            yield button

        yield Footer()

    def on_mount(self) -> None:
        self.title = self.dashboard_config.get("title", "MISSION CONTROL")
        self.sub_title = f"{self.dashboard_config.get('workspace_root', BASE_DIR.parent)} | {platform.system()} {platform.machine()}"
        core.append_runtime_log("Mission Control Textual booted")
        self.set_interval(1.0, self.refresh_dashboard)
        self.set_interval(0.11, self.advance_animation)
        self.set_interval(180.0, self.refresh_animation_catalog_async)
        self.refresh_animation_catalog_async()
        self.refresh_dashboard()

    def action_show_tab(self, tab: str) -> None:
        self.active_tab = tab
        tabs = self.query_one("#tabs", TabbedContent)
        tabs.active = tab

    def on_tabbed_content_tab_activated(self, event: TabbedContent.TabActivated) -> None:
        self.active_tab = event.tabbed_content.active or self.active_tab

    def action_refresh_dashboard(self) -> None:
        self.refresh_dashboard()

    def action_toggle_animation_pause(self) -> None:
        self.animation_paused = not self.animation_paused
        label = "paused" if self.animation_paused else "resumed"
        core.append_runtime_log(f"Animation loop {label}")
        self.refresh_dashboard()

    def action_next_animation_scene(self) -> None:
        if not self.animation_sources:
            return
        self.animation_source_index = (self.animation_source_index + 1) % len(self.animation_sources)
        self.animation_frame_index = 0
        self.animation_scene_ticks = 0
        core.append_runtime_log(f"Animation scene -> {self.animation_sources[self.animation_source_index].name}")
        self.refresh_dashboard()

    def action_previous_animation_scene(self) -> None:
        if not self.animation_sources:
            return
        self.animation_source_index = (self.animation_source_index - 1) % len(self.animation_sources)
        self.animation_frame_index = 0
        self.animation_scene_ticks = 0
        core.append_runtime_log(f"Animation scene -> {self.animation_sources[self.animation_source_index].name}")
        self.refresh_dashboard()

    def action_pick_animation_favorite(self, index: int) -> None:
        if index <= 0 or index > len(self.animation_sources):
            return
        self.animation_source_index = index - 1
        self.animation_frame_index = 0
        self.animation_scene_ticks = 0
        core.append_runtime_log(f"Animation favorite -> {self.animation_sources[self.animation_source_index].name}")
        self.refresh_dashboard()

    def action_launch_shortcut(self, index: int) -> None:
        if index <= 0 or index > len(self.shortcuts):
            return
        shortcut = self.shortcuts[index - 1]
        self._launch_command(index, shortcut)

    def on_key(self, event) -> None:
        key = getattr(event, "key", "")
        if len(key) == 1 and key.isdigit():
            index = int(key)
            if self.active_tab == "wall":
                self.action_pick_animation_favorite(index)
            else:
                self.action_launch_shortcut(index)
            if hasattr(event, "stop"):
                event.stop()

    def on_button_pressed(self, event: Button.Pressed) -> None:
        button_id = event.button.id or ""
        if button_id.startswith("launcher-"):
            try:
                index = int(button_id.split("-", 1)[1])
            except ValueError:
                return
            if 1 <= index <= len(self.shortcuts):
                self._launch_command(index, self.shortcuts[index - 1])

    def _launch_command(self, index: int, shortcut: dict[str, str]) -> None:
        label = shortcut.get("label", f"Launcher {index}")
        command = shortcut.get("command", "")
        if not command:
            return

        self.launch_states[index] = {
            "last_launched": time.time(),
            "last_command": command,
        }
        core.append_runtime_log(f"Launching {label}: {command}")
        log = self.query_one("#activity-log", RichLog)
        log.write(f"Launching {label}: {command}")

        if platform.system() == "Windows":
            subprocess.Popen(["cmd", "/c", command], cwd=str(BASE_DIR.parent))
        else:
            subprocess.Popen(["/bin/zsh", "-lc", command], cwd=str(BASE_DIR.parent))

    def refresh_dashboard(self) -> None:
        self.missions = core.load_missions(self.missions_path, self.dashboard_config) or self.missions

        for card in self.query(DashboardCard):
            card.refresh()

        for card in self.query(AnimationCard):
            card.refresh()
        for card in self.query(AnimationFavoriteCard):
            card.refresh()

        log_widget = self.query_one("#activity-log", RichLog)
        if log_widget:
            lines = core.read_recent_log_lines()
            if len(lines) < self._last_log_lines:
                log_widget.clear()
                self._last_log_lines = 0
            for line in lines[self._last_log_lines :]:
                log_widget.write(line)
            self._last_log_lines = len(lines)

            log_widget.scroll_end(animate=False)

    def advance_animation(self) -> None:
        if self.animation_paused or not self.animation_sources:
            return

        scene = self.animation_sources[self.animation_source_index]
        self.animation_frame_index = (self.animation_frame_index + 1) % max(1, scene.frame_count)
        self.animation_scene_ticks += 1
        if self.animation_scene_ticks >= self._animation_scene_hold:
            self.animation_source_index = (self.animation_source_index + 1) % len(self.animation_sources)
            self.animation_frame_index = 0
            self.animation_scene_ticks = 0
            core.append_runtime_log(f"Animation scene -> {self.animation_sources[self.animation_source_index].name}")
        for card in self.query(AnimationCard):
            card.refresh()
        for card in self.query(AnimationFavoriteCard):
            card.refresh()

    def refresh_animation_catalog_async(self) -> None:
        if self._animation_catalog_refreshing:
            return

        self._animation_catalog_refreshing = True

        def worker() -> None:
            try:
                catalog = motion.build_animation_catalog(self.dashboard_config)
            except Exception:
                catalog = motion.build_local_animation_catalog()
            self.call_from_thread(self._apply_animation_catalog, catalog)

        threading.Thread(target=worker, daemon=True).start()

    def _apply_animation_catalog(self, catalog: list[motion.AnimationScene]) -> None:
        self._animation_catalog_refreshing = False
        if not catalog:
            catalog = motion.build_local_animation_catalog()

        current_name = self.animation_sources[self.animation_source_index].name if self.animation_sources else ""
        next_index = 0
        for index, source in enumerate(catalog):
            if source.name == current_name:
                next_index = index
                break

        self.animation_sources = catalog
        self.animation_source_index = next_index
        self.animation_frame_index = 0
        self.animation_scene_ticks = 0
        core.append_runtime_log(f"Animation catalog -> {len(catalog)} sources")
        for card in self.query(AnimationCard):
            card.refresh()


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Interactive Mission Control control room")
    parser.add_argument("--embedded", action="store_true", help="Run inline in the current terminal pane")
    parser.add_argument("--once", action="store_true", help="Render a single snapshot and exit")
    return parser.parse_args()


def main() -> int:
    args = parse_args()
    config = core.load_dashboard_config(BASE_DIR / "dashboard.yaml")
    missions = core.load_missions(BASE_DIR / "active_missions.json", config)
    if args.once:
        snapshot = core.render_snapshot(config, missions)
        console.print(snapshot)
        return 0
    app = MissionControlApp()
    app.run(inline=args.embedded)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
