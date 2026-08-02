#!/usr/bin/env python3
from __future__ import annotations

import argparse
import copy
import json
import os
import platform
import socket
import subprocess
import threading
import time
import urllib.error
import urllib.parse
import urllib.request
from datetime import datetime, timedelta
from pathlib import Path
from typing import Any

import psutil
from rich import box
from rich.console import Console, Group
from rich.layout import Layout
from rich.live import Live
from rich.panel import Panel
from rich.table import Table
from rich.text import Text

try:
    import yaml
except Exception:  # pragma: no cover - fallback keeps the app usable without PyYAML
    yaml = None

try:
    from terminaltexteffects.effects.effect_decrypt import Decrypt
except Exception:  # pragma: no cover - optional flourish only
    Decrypt = None


BASE_DIR = Path(__file__).resolve().parent
DEFAULT_CONFIG_PATH = BASE_DIR / "dashboard.yaml"
DEFAULT_MISSIONS_PATH = BASE_DIR / "active_missions.json"

console = Console()

STATE_LOCK = threading.Lock()
MISSION_COUNTS: dict[str, int] = {}
MISSION_TIMESTAMPS: dict[str, float] = {}
MISSION_RATES: dict[str, float] = {}
PROCESS_CACHE: dict[str, Any] = {"updated": 0.0, "items": []}
WEATHER_CACHE: dict[str, Any] = {"updated": 0.0, "data": None}
BORDER_ANIMS: dict[str, Any] = {}
RUNTIME_LOG_PATH = BASE_DIR / "mission_control.log"

DEFAULT_DASHBOARD: dict[str, Any] = {
    "title": "MISSION CONTROL // ECHO",
    "workspace_root": str(BASE_DIR.parent),
    "refresh_hz": 4.0,
    "weather": {
        "enabled": True,
        "city": "Detroit",
        "latitude": 42.3314,
        "longitude": -83.0458,
        "units": "fahrenheit",
    },
    "missions": [
        {
            "id": "legacy-core",
            "name": "Legacy Mission Core",
            "target_dir": str(BASE_DIR.parent / "mission_control_app"),
            "file_ext": ".py",
            "total_delta": 3,
            "color": "cyan",
            "note": "Older terminal dashboard wired into the new pipeline",
            "exclude_dirs": ["venv", ".venv", "__pycache__", ".git"],
        },
        {
            "id": "archive-v1",
            "name": "V1 Archive",
            "target_dir": str(BASE_DIR.parent / "V1"),
            "file_ext": ".py",
            "total_delta": 2,
            "color": "magenta",
            "note": "Original mission control reference build",
            "exclude_dirs": ["venv", ".venv", "__pycache__", ".git"],
        },
        {
            "id": "new-control-room",
            "name": "New Control Room",
            "target_dir": str(BASE_DIR),
            "file_ext": ".py",
            "total_delta": 5,
            "color": "green",
            "note": "The rebuilt dashboard itself",
            "exclude_dirs": ["venv", ".venv", "__pycache__", ".git"],
        },
    ],
    "services": [
        {
            "name": "Mission Control",
            "kind": "process",
            "match": "mission_control.py",
            "note": "Terminal dashboard runtime",
        },
        {
            "name": "VS Code",
            "kind": "process",
            "match": "Code",
            "note": "Editor session",
        },
        {
            "name": "Python Runner",
            "kind": "process",
            "match": "python",
            "note": "Any active python pipeline workers",
        },
        {
            "name": "Local Port 3000",
            "kind": "port",
            "host": "127.0.0.1",
            "port": 3000,
            "note": "Preview or local web app",
        },
    ],
    "endpoints": [
        {
            "name": "Local Health",
            "url": "http://127.0.0.1:3000/health",
            "note": "Healthcheck for the local preview app",
        },
        {
            "name": "Mission API",
            "url": "http://127.0.0.1:8000/health",
            "note": "Placeholder for your local pipeline API",
        },
    ],
    "shortcuts": [
        {
            "label": "Open Workspace",
            "command": f"code {BASE_DIR.parent}",
            "note": "Launch this folder in VS Code",
        },
        {
            "label": "Run Mission Control",
            "command": f"cd {BASE_DIR} && ./start.sh --embedded",
            "note": "Start the terminal dashboard in-place",
        },
        {
            "label": "Edit Pipeline",
            "command": f"code {BASE_DIR / 'dashboard.yaml'}",
            "note": "Tweak widgets and mission targets",
        },
        {
            "label": "Tail Missions",
            "command": f"cat {BASE_DIR / 'active_missions.json'}",
            "note": "Inspect the live pipeline seed",
        },
    ],
}

WEATHER_CODE_MAP = {
    0: "Clear",
    1: "Mostly clear",
    2: "Partly cloudy",
    3: "Overcast",
    45: "Fog",
    48: "Rime fog",
    51: "Light drizzle",
    53: "Drizzle",
    55: "Heavy drizzle",
    61: "Light rain",
    63: "Rain",
    65: "Heavy rain",
    71: "Light snow",
    73: "Snow",
    75: "Heavy snow",
    80: "Rain showers",
    81: "Heavy showers",
    82: "Violent showers",
    95: "Thunderstorm",
    96: "Thunderstorm with hail",
    99: "Severe thunderstorm with hail",
}


def deep_merge(base: dict[str, Any], override: dict[str, Any]) -> dict[str, Any]:
    merged = copy.deepcopy(base)
    for key, value in override.items():
        if (
            key in merged
            and isinstance(merged[key], dict)
            and isinstance(value, dict)
        ):
            merged[key] = deep_merge(merged[key], value)
        else:
            merged[key] = value
    return merged


def load_structured_file(path: Path) -> dict[str, Any]:
    if not path.exists():
        return {}

    text = path.read_text(encoding="utf-8")
    if path.suffix.lower() in {".yml", ".yaml"}:
        if yaml is None:
            return {}
        loaded = yaml.safe_load(text) or {}
        if isinstance(loaded, dict):
            return loaded
        return {}

    loaded = json.loads(text or "{}")
    if isinstance(loaded, dict):
        return loaded
    return {}


def load_dashboard_config(path: Path) -> dict[str, Any]:
    config = copy.deepcopy(DEFAULT_DASHBOARD)
    try:
        config = deep_merge(config, load_structured_file(path))
    except Exception as exc:
        console.print(f"[yellow]Could not read config {path}: {exc}[/yellow]")
    return config


def save_json(path: Path, payload: Any) -> None:
    path.write_text(json.dumps(payload, indent=2, ensure_ascii=True) + "\n", encoding="utf-8")


def load_missions(missions_path: Path, config: dict[str, Any]) -> list[dict[str, Any]]:
    if missions_path.exists():
        try:
            payload = json.loads(missions_path.read_text(encoding="utf-8") or "[]")
            if isinstance(payload, list) and payload:
                return payload
        except Exception as exc:
            console.print(f"[yellow]Could not read missions {missions_path}: {exc}[/yellow]")

    seed_missions: list[dict[str, Any]] = []
    for mission in config.get("missions", []):
        if not isinstance(mission, dict):
            continue
        current = count_target(
            mission.get("target_dir", ""),
            mission.get("file_ext", ""),
            mission.get("exclude_dirs", []),
        )
        total_delta = max(int(mission.get("total_delta", 1) or 1), 1)
        seed_missions.append(
            {
                **mission,
                "total": max(current + total_delta, current + 1),
            }
        )

    if seed_missions:
        save_json(missions_path, seed_missions)
    return seed_missions


def count_target(target: str, file_ext: str = "", exclude_dirs: list[str] | tuple[str, ...] | None = None) -> int:
    target_path = Path(target)
    if not target_path.exists():
        return 0

    exclude = {name for name in (exclude_dirs or [])}
    if target_path.is_file():
        if file_ext and target_path.suffix != file_ext:
            return 0
        try:
            with target_path.open("r", encoding="utf-8", errors="ignore") as handle:
                return sum(1 for _ in handle)
        except Exception:
            return 1

    total = 0
    for root, dirs, files in os.walk(target_path):
        dirs[:] = [d for d in dirs if d not in exclude and not d.startswith(".")]
        for filename in files:
            if filename.startswith(".") and not file_ext:
                continue
            if file_ext and not filename.endswith(file_ext):
                continue
            total += 1
    return total


def refresh_mission_state(mission: dict[str, Any]) -> int:
    mission_id = mission["id"]
    now = time.time()
    poll_interval = float(mission.get("poll_interval", 5.0))

    def worker() -> None:
        current = count_target(
            mission.get("target_dir", ""),
            mission.get("file_ext", ""),
            mission.get("exclude_dirs", []),
        )
        with STATE_LOCK:
            old_count = MISSION_COUNTS.get(mission_id, current)
            old_time = MISSION_TIMESTAMPS.get(mission_id, now)
            if mission_id in MISSION_COUNTS and now > old_time:
                delta = current - old_count
                elapsed = max(now - old_time, 0.001)
                if delta >= 0:
                    rate = delta / elapsed
                    prior = MISSION_RATES.get(mission_id, 0.0)
                    MISSION_RATES[mission_id] = prior * 0.65 + rate * 0.35
            MISSION_COUNTS[mission_id] = current
            MISSION_TIMESTAMPS[mission_id] = time.time()

    if mission_id not in MISSION_COUNTS:
        worker()
    elif now - MISSION_TIMESTAMPS.get(mission_id, 0.0) > poll_interval:
        threading.Thread(target=worker, daemon=True).start()
        with STATE_LOCK:
            MISSION_TIMESTAMPS[mission_id] = now

    return MISSION_COUNTS.get(mission_id, 0)


def fmt_bytes(num: float) -> str:
    units = ["B", "KB", "MB", "GB", "TB", "PB"]
    value = float(num)
    for unit in units:
        if value < 1024.0:
            return f"{value:.1f} {unit}"
        value /= 1024.0
    return f"{value:.1f} EB"


def fmt_duration(seconds: float) -> str:
    delta = timedelta(seconds=int(seconds))
    days = delta.days
    hours, remainder = divmod(delta.seconds, 3600)
    minutes, secs = divmod(remainder, 60)
    parts = []
    if days:
        parts.append(f"{days}d")
    if hours or parts:
        parts.append(f"{hours}h")
    if minutes or parts:
        parts.append(f"{minutes}m")
    parts.append(f"{secs}s")
    return " ".join(parts)


def build_progress_bar(label: str, count: int, total: int, color: str) -> str:
    width = 34
    total = max(total, 1)
    ratio = min(max(count / total, 0.0), 1.0)
    filled = int(width * ratio)
    bar = "█" * filled + "░" * (width - filled)
    state = "✅" if count >= total else "⏳"
    return f"[{color}]{state} {label:<24}[/{color}] |{bar}| [white]{count}/{total}[/white] ({ratio * 100:5.1f}%)"


def build_marquee(key: str, text: str, width: int, style: str) -> Text:
    if Decrypt is None:
        return Text(text[:width].ljust(width), style=style, no_wrap=True, overflow="crop")

    if key not in BORDER_ANIMS:
        try:
            BORDER_ANIMS[key] = iter(Decrypt((text + " ") * ((width // len(text)) + 2)))
        except Exception:
            return Text(text[:width].ljust(width), style=style, no_wrap=True, overflow="crop")

    try:
        frame = next(BORDER_ANIMS[key])
    except StopIteration:
        try:
            BORDER_ANIMS[key] = iter(Decrypt((text + " ") * ((width // len(text)) + 2)))
            frame = next(BORDER_ANIMS[key])
        except Exception:
            return Text(text[:width].ljust(width), style=style, no_wrap=True, overflow="crop")

    plain = Text.from_ansi(frame).plain
    return Text(plain[:width].ljust(width), style=style, no_wrap=True, overflow="crop")


def make_table(rows: list[tuple[str, str]], title: str | None = None, box_style: box.Box = box.SIMPLE_HEAVY) -> Table:
    table = Table(box=box_style, show_header=False, expand=True, pad_edge=False)
    table.add_column("Key", style="bold white", no_wrap=True)
    table.add_column("Value", style="white", ratio=1)
    for left, right in rows:
        table.add_row(left, right)
    if title:
        table.caption = title
    return table


def build_header_panel(config: dict[str, Any]) -> Panel:
    now = datetime.now().astimezone().strftime("%H:%M:%S")
    workspace_root = config.get("workspace_root", str(BASE_DIR.parent))
    header_left = Text.from_markup(
        f"[bold white]{config.get('title', 'MISSION CONTROL')}[/bold white] [dim]//[/dim] [bright_cyan]Terminal + VS Code[/bright_cyan]"
    )
    header_right = Text.from_markup(
        f"[white]{workspace_root}[/white] [dim]|[/dim] [cyan]{platform.system()} {platform.machine()}[/cyan] [dim]|[/dim] [bold white]{now}[/bold white]"
    )
    return Panel(Group(header_left, header_right), border_style="bright_cyan", box=box.DOUBLE)


def get_cpu_panel() -> Panel:
    cpu_percents = psutil.cpu_percent(interval=None, percpu=True)
    table = Table(box=box.SIMPLE_HEAVY, show_header=False, expand=True)
    table.add_column("Core L", style="cyan")
    table.add_column("Core R", style="green")

    half = (len(cpu_percents) + 1) // 2
    for index in range(half):
        left = f"C{index:02d}: {cpu_percents[index]:>4.1f}%"
        right_index = index + half
        right = f"C{right_index:02d}: {cpu_percents[right_index]:>4.1f}%" if right_index < len(cpu_percents) else ""
        table.add_row(left, right)

    load = "n/a"
    if hasattr(os, "getloadavg"):
        try:
            avg = os.getloadavg()
            count = max(psutil.cpu_count() or 1, 1)
            load = f"{avg[0] / count * 100:.1f}% / {avg[1] / count * 100:.1f}% / {avg[2] / count * 100:.1f}%"
        except OSError:
            load = "n/a"

    footer = Text.from_markup(f"[bold cyan]Load Avg:[/bold cyan] [white]{load}[/white]")
    body = Group(build_marquee("cpu", "SYS.CORE.YIELD // AURA.SYNC.ACTIVE // ", 80, "cyan"), table, footer)
    return Panel(body, title="[b]CPU[/b]", border_style="cyan", box=box.ROUNDED)


def get_memory_panel() -> Panel:
    mem = psutil.virtual_memory()
    swap = psutil.swap_memory()
    disk = psutil.disk_usage(str(BASE_DIR.parent))
    uptime = fmt_duration(time.time() - psutil.boot_time())

    rows = [
        ("RAM Total", f"{fmt_bytes(mem.total)}"),
        ("RAM Used", f"[red]{fmt_bytes(mem.used)} ({mem.percent:.1f}%)[/red]"),
        ("RAM Free", f"[green]{fmt_bytes(mem.available)}[/green]"),
        ("Swap Used", f"{fmt_bytes(swap.used)} / {fmt_bytes(swap.total)}"),
        ("Disk Used", f"[magenta]{fmt_bytes(disk.used)} ({disk.percent:.1f}%)[/magenta]"),
        ("Disk Free", f"{fmt_bytes(disk.free)}"),
        ("Uptime", f"{uptime}"),
    ]
    table = make_table(rows)
    return Panel(table, title="[b]Memory & Disk[/b]", border_style="magenta", box=box.ROUNDED)


def scan_processes(limit: int = 10) -> list[dict[str, Any]]:
    now = time.time()
    if PROCESS_CACHE["items"] and now - PROCESS_CACHE["updated"] < 2.0:
        return PROCESS_CACHE["items"]

    items: list[dict[str, Any]] = []
    for proc in psutil.process_iter(["name", "cpu_percent", "memory_percent", "cmdline"]):
        try:
            info = proc.info
            items.append(
                {
                    "name": info.get("name") or "unknown",
                    "cpu": info.get("cpu_percent") or 0.0,
                    "mem": info.get("memory_percent") or 0.0,
                    "cmdline": " ".join(info.get("cmdline") or []),
                }
            )
        except (psutil.NoSuchProcess, psutil.AccessDenied):
            continue

    items.sort(key=lambda item: item["cpu"], reverse=True)
    PROCESS_CACHE["items"] = items[:limit]
    PROCESS_CACHE["updated"] = now
    return PROCESS_CACHE["items"]


def get_processes_panel() -> Panel:
    items = scan_processes()
    table = Table(box=box.SIMPLE_HEAVY, show_header=True, expand=True)
    table.add_column("Process", style="white")
    table.add_column("CPU%", style="red", justify="right")
    table.add_column("MEM%", style="magenta", justify="right")

    for item in items:
        table.add_row(item["name"][:24], f"{item['cpu']:.1f}", f"{item['mem']:.1f}")

    net = psutil.net_io_counters()
    subtitle = f"Net▼ {fmt_bytes(net.bytes_recv)} | Net▲ {fmt_bytes(net.bytes_sent)}"
    return Panel(table, title="[b]Top Processes[/b]", border_style="yellow", subtitle=subtitle, box=box.ROUNDED)


def fetch_weather(config: dict[str, Any]) -> dict[str, Any]:
    weather_cfg = config.get("weather", {})
    if not weather_cfg.get("enabled", False):
        return {"enabled": False}

    latitude = weather_cfg.get("latitude")
    longitude = weather_cfg.get("longitude")
    if latitude is None or longitude is None:
        return {"enabled": False, "error": "missing coordinates"}

    cache = WEATHER_CACHE.get("data")
    if cache and time.time() - WEATHER_CACHE["updated"] < 900:
        if cache.get("latitude") == latitude and cache.get("longitude") == longitude:
            return cache

    params = urllib.parse.urlencode(
        {
            "latitude": latitude,
            "longitude": longitude,
            "current": "temperature_2m,wind_speed_10m,weather_code",
            "temperature_unit": "fahrenheit" if weather_cfg.get("units", "fahrenheit") == "fahrenheit" else "celsius",
            "wind_speed_unit": "mph" if weather_cfg.get("units", "fahrenheit") == "fahrenheit" else "kmh",
            "timezone": "auto",
        }
    )
    url = f"https://api.open-meteo.com/v1/forecast?{params}"

    try:
        with urllib.request.urlopen(url, timeout=4) as response:
            data = json.loads(response.read().decode("utf-8"))
        current = data.get("current", {})
        weather = {
            "enabled": True,
            "latitude": latitude,
            "longitude": longitude,
            "city": weather_cfg.get("city", "Local"),
            "temperature": current.get("temperature_2m"),
            "units": "F" if weather_cfg.get("units", "fahrenheit") == "fahrenheit" else "C",
            "wind": current.get("wind_speed_10m"),
            "code": current.get("weather_code"),
            "description": WEATHER_CODE_MAP.get(current.get("weather_code"), f"Code {current.get('weather_code')}"),
            "updated": data.get("current_time") or datetime.now().astimezone().isoformat(),
        }
        WEATHER_CACHE["data"] = weather
        WEATHER_CACHE["updated"] = time.time()
        return weather
    except (urllib.error.URLError, TimeoutError, json.JSONDecodeError, ValueError, OSError) as exc:
        return {"enabled": True, "error": str(exc), "city": weather_cfg.get("city", "Local")}


def get_clock_panel(config: dict[str, Any]) -> Panel:
    now = datetime.now().astimezone()
    boot = datetime.fromtimestamp(psutil.boot_time()).astimezone()
    rows = [
        ("Local Time", now.strftime("%I:%M:%S %p")),
        ("Date", now.strftime("%A, %B %d, %Y")),
        ("Timezone", now.tzname() or "local"),
        ("Boot", boot.strftime("%Y-%m-%d %H:%M")),
        ("Mode", "Terminal / VS Code"),
    ]
    table = make_table(rows)
    body = Group(build_marquee("clock", "CLOCK.SYNC // ECHO ONLINE // ", 80, "bright_cyan"), table)
    return Panel(body, title="[b]Clock[/b]", border_style="bright_cyan", box=box.ROUNDED)


def get_weather_panel(config: dict[str, Any]) -> Panel:
    weather = fetch_weather(config)
    weather_cfg = config.get("weather", {})
    if not weather.get("enabled"):
        rows = [
            ("Weather", "Disabled in dashboard.yaml"),
            ("Tip", "Set latitude/longitude to enable live weather"),
        ]
        table = make_table(rows)
        return Panel(table, title="[b]Weather[/b]", border_style="blue", box=box.ROUNDED)

    if weather.get("error"):
        rows = [
            ("Weather", f"[red]Offline[/red]"),
            ("City", str(weather_cfg.get("city", "Local"))),
            ("Error", weather["error"]),
        ]
        table = make_table(rows)
        return Panel(table, title="[b]Weather[/b]", border_style="blue", box=box.ROUNDED)

    rows = [
        ("City", str(weather.get("city", "Local"))),
        ("Temp", f"{weather.get('temperature')}°{weather.get('units', '')}"),
        ("Conditions", str(weather.get("description", "Unknown"))),
        ("Wind", f"{weather.get('wind')}"),
        ("Updated", str(weather.get("updated", ""))[:19].replace("T", " ")),
    ]
    table = make_table(rows)
    return Panel(table, title="[b]Weather[/b]", border_style="blue", box=box.ROUNDED)


def get_workspace_panel(config: dict[str, Any], missions: list[dict[str, Any]]) -> Panel:
    rows: list[tuple[str, str]] = [("Workspace", config.get("workspace_root", str(BASE_DIR.parent)))]
    for mission in missions:
        count = refresh_mission_state(mission)
        total = max(int(mission.get("total", 1)), 1)
        label = mission.get("name", mission["id"])
        rows.append((label, f"{count} / {total} {mission.get('file_ext', '')} files"))
    table = make_table(rows)
    return Panel(table, title="[b]Workspace[/b]", border_style="green", box=box.ROUNDED)


def get_services_panel(config: dict[str, Any]) -> Panel:
    rows: list[tuple[str, str]] = []
    for service in config.get("services", []):
        if not isinstance(service, dict):
            continue
        name = service.get("name", "Service")
        kind = service.get("kind", "process")
        if kind == "process":
            match = service.get("match", "")
            online = process_matches(match)
            detail = f"process:{match}"
        elif kind == "port":
            host = service.get("host", "127.0.0.1")
            port = int(service.get("port", 0) or 0)
            online = port_is_open(host, port)
            detail = f"{host}:{port}"
        elif kind == "path":
            path = Path(str(service.get("path", "")))
            online = path.exists()
            detail = str(path)
        else:
            online = False
            detail = kind

        status = "[bold green]ONLINE[/bold green]" if online else "[bold red]OFFLINE[/bold red]"
        note = service.get("note", "")
        rows.append((name, f"{status}  [dim]{detail}[/dim] {f' - {note}' if note else ''}"))

    table = make_table(rows or [("Services", "No service checks configured")])
    return Panel(table, title="[b]Services[/b]", border_style="yellow", box=box.ROUNDED)


def get_shortcuts_panel(config: dict[str, Any]) -> Panel:
    rows: list[tuple[str, str]] = []
    for shortcut in config.get("shortcuts", []):
        if not isinstance(shortcut, dict):
            continue
        label = shortcut.get("label", "Shortcut")
        command = shortcut.get("command", "")
        note = shortcut.get("note", "")
        rows.append((label, f"[white]{command}[/white]{f'  [dim]{note}[/dim]' if note else ''}"))

    table = make_table(rows or [("Shortcuts", "No shortcuts configured")])
    return Panel(table, title="[b]Shortcuts[/b]", border_style="bright_magenta", box=box.ROUNDED)


def _service_alert_counts(config: dict[str, Any]) -> tuple[int, int]:
    service_rows = config.get("services", [])
    endpoint_rows = config.get("endpoints", [])
    offline_services = 0
    offline_endpoints = 0

    for service in service_rows:
        if not isinstance(service, dict):
            continue
        kind = service.get("kind", "process")
        if kind == "process":
            online = process_matches(str(service.get("match", "")))
        elif kind == "port":
            online = port_is_open(str(service.get("host", "127.0.0.1")), int(service.get("port", 0) or 0))
        elif kind == "path":
            online = Path(str(service.get("path", ""))).exists()
        else:
            online = False
        if not online:
            offline_services += 1

    for endpoint in endpoint_rows:
        if not isinstance(endpoint, dict):
            continue
        url = str(endpoint.get("url", ""))
        if not url:
            continue
        if check_endpoint(url).get("state") != "ONLINE":
            offline_endpoints += 1

    return offline_services, offline_endpoints


def get_alerts_panel(config: dict[str, Any], missions: list[dict[str, Any]]) -> Panel:
    offline_services, offline_endpoints = _service_alert_counts(config)
    incomplete_missions = sum(
        1
        for mission in missions
        if refresh_mission_state(mission) < max(int(mission.get("total", 1)), 1)
    )
    weather = fetch_weather(config)
    alerts = []
    if offline_services:
        alerts.append(f"{offline_services} service{'s' if offline_services != 1 else ''} offline")
    if offline_endpoints:
        alerts.append(f"{offline_endpoints} endpoint{'s' if offline_endpoints != 1 else ''} offline")
    if incomplete_missions:
        alerts.append(f"{incomplete_missions} mission{'s' if incomplete_missions != 1 else ''} in flight")
    if weather.get("error"):
        alerts.append("weather feed unavailable")

    rows = [
        ("Signal", f"{len(alerts)} active alert{'s' if len(alerts) != 1 else ''}" if alerts else "[bold green]CLEAR[/bold green]"),
        ("Alerts", ", ".join(alerts) if alerts else "No active alerts"),
    ]
    table = make_table(rows)
    return Panel(table, title="[b]Alerts[/b]", border_style="red", box=box.ROUNDED)


def get_health_panel(config: dict[str, Any]) -> Panel:
    cpu = psutil.cpu_percent(interval=0.0)
    memory = psutil.virtual_memory()
    disk = psutil.disk_usage(BASE_DIR.parent)
    uptime_seconds = max(0, int(time.time() - psutil.boot_time()))
    rows = [
        ("CPU", f"{cpu:.1f}%"),
        ("RAM", f"{memory.percent:.1f}%"),
        ("Disk", f"{disk.percent:.1f}%"),
        ("Uptime", fmt_duration(uptime_seconds)),
        ("Mode", "Terminal / VS Code"),
    ]
    table = make_table(rows)
    return Panel(table, title="[b]Health[/b]", border_style="bright_green", box=box.ROUNDED)


def get_status_panel(config: dict[str, Any], missions: list[dict[str, Any]], *, active_tab: str = "snapshot", animation_source: str = "local") -> Panel:
    services_online = 0
    for service in config.get("services", []):
        if not isinstance(service, dict):
            continue
        kind = service.get("kind", "process")
        if kind == "process":
            online = process_matches(str(service.get("match", "")))
        elif kind == "port":
            online = port_is_open(str(service.get("host", "127.0.0.1")), int(service.get("port", 0) or 0))
        elif kind == "path":
            online = Path(str(service.get("path", ""))).exists()
        else:
            online = False
        if online:
            services_online += 1

    mission_count = len(missions)
    completed = sum(1 for mission in missions if refresh_mission_state(mission) >= max(int(mission.get("total", 1)), 1))
    rows = [
        ("Tab", active_tab.upper()),
        ("Workspace", config.get("workspace_root", str(BASE_DIR.parent))),
        ("Services", f"{services_online}/{len(config.get('services', []))} online"),
        ("Missions", f"{completed}/{mission_count} complete"),
        ("Animation", animation_source),
    ]
    table = make_table(rows)
    return Panel(table, title="[b]Status[/b]", border_style="cyan", box=box.ROUNDED)


def _shortcut_match(shortcut: dict[str, Any]) -> str:
    explicit = str(shortcut.get("match", "")).strip()
    if explicit:
        return explicit

    command = str(shortcut.get("command", ""))
    if command.startswith("code "):
        return "Code"
    if "mission_control_tui.py" in command:
        return "mission_control_tui.py"
    if command.startswith("python") or "python " in command:
        return "python"
    return command.split()[0] if command.split() else ""


def _shortcut_mode(shortcut: dict[str, Any]) -> str:
    mode = str(shortcut.get("mode", "")).strip().lower()
    if mode in {"one-shot", "persistent"}:
        return mode
    command = str(shortcut.get("command", ""))
    if command.startswith("cat ") or command.startswith("tail "):
        return "one-shot"
    return "persistent"


def _launcher_state_label(shortcut: dict[str, Any], launch_state: dict[str, Any] | None = None) -> str:
    launch_state = launch_state or {}
    mode = _shortcut_mode(shortcut)
    match = _shortcut_match(shortcut)
    last_launched = float(launch_state.get("last_launched", 0.0) or 0.0)
    launched_recently = last_launched > 0 and (time.time() - last_launched) < 30.0
    live = process_matches(match) if match else False

    if mode == "one-shot":
        if launched_recently:
            return "[bold green]DONE[/bold green]"
        return "[dim]IDLE[/dim]"

    if live:
        return "[bold green]LIVE[/bold green]"
    if launched_recently:
        return "[bold yellow]WAITING[/bold yellow]"
    return "[dim]OFFLINE[/dim]"


def get_launcher_monitor_panel(config: dict[str, Any], launch_state: dict[int, dict[str, Any]] | None = None) -> Panel:
    rows: list[tuple[str, str]] = []
    for index, shortcut in enumerate(config.get("shortcuts", []), start=1):
        if not isinstance(shortcut, dict):
            continue
        state = _launcher_state_label(shortcut, (launch_state or {}).get(index))
        note = shortcut.get("note", "")
        rows.append(
            (
                f"{index}. {shortcut.get('label', 'Launcher')}",
                f"{state}  [dim]{_shortcut_match(shortcut) or 'manual'}[/dim]{f'  [dim]- {note}[/dim]' if note else ''}",
            )
        )
    table = make_table(rows or [("Launchers", "No shortcuts configured")])
    return Panel(table, title="[b]Launch Monitor[/b]", border_style="magenta", box=box.ROUNDED)


def get_pipeline_panel(missions: list[dict[str, Any]]) -> Panel:
    lines: list[str] = ["[bold cyan]ACTIVE MISSIONS[/bold cyan]"]
    all_complete = True

    for mission in missions:
        count = refresh_mission_state(mission)
        total = max(int(mission.get("total", 1)), 1)
        if count < total:
            all_complete = False
        lines.append(build_progress_bar(mission.get("name", mission["id"]), count, total, mission.get("color", "blue")))
        note = mission.get("note")
        if note:
            lines.append(f"[dim]  {note}[/dim]")

    if not missions:
        lines.append("[italic]No missions loaded in active_missions.json[/italic]")
        all_complete = False

    lines.append("")
    if all_complete and missions:
        lines.append("[bold green]MISSION COMPLETE - ALL PIPELINE STAGES ARE READY[/bold green]")
        lines.append("[bold magenta]\\(^o^)/   ECHO APPROVES THE RUN   \\(^o^)/[/bold magenta]")
    else:
        lines.append("[bold yellow]Pipeline is live. The dashboard will keep watching active_missions.json.[/bold yellow]")

    body = Group(build_marquee("pipeline", "BEAVERY.PIPELINE.SYNC // FLIGHT.READY // ", 80, "green"), Text.from_markup("\n".join(lines)))
    return Panel(body, title="[b]Pipeline[/b]", border_style="green", padding=(1, 2), box=box.ROUNDED)


def get_telemetry_panel(missions: list[dict[str, Any]]) -> Panel:
    total_rendered = sum(MISSION_COUNTS.values())
    total_rate = sum(MISSION_RATES.values())
    complete = sum(1 for mission in missions if MISSION_COUNTS.get(mission["id"], 0) >= max(int(mission.get("total", 1)), 1))
    total = len(missions) if missions else 0

    rows = [
        ("Rendered", f"{total_rendered:,}"),
        ("Velocity", f"{total_rate:.2f} items/sec"),
        ("Missions", f"{complete}/{total} complete"),
    ]

    if missions and complete == total:
        rows.append(("Status", "[bold green]ALL SYSTEMS GO[/bold green]"))
    else:
        rows.append(("Status", "[bold yellow]Monitoring[/bold yellow]"))

    table = make_table(rows)
    body = Group(build_marquee("telemetry", "TELEMETRY.UPLINK.ESTABLISHED // LINK.OK // ", 80, "cyan"), table)
    return Panel(body, title="[b]Telemetry[/b]", border_style="cyan", box=box.ROUNDED)


def process_matches(pattern: str) -> bool:
    if not pattern:
        return False

    needle = pattern.lower()
    for proc in psutil.process_iter(["name", "cmdline"]):
        try:
            name = (proc.info.get("name") or "").lower()
            cmdline = " ".join(proc.info.get("cmdline") or []).lower()
            if needle in name or needle in cmdline:
                return True
        except (psutil.NoSuchProcess, psutil.AccessDenied):
            continue
    return False


def port_is_open(host: str, port: int, timeout: float = 0.35) -> bool:
    if port <= 0:
        return False
    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as sock:
        sock.settimeout(timeout)
        return sock.connect_ex((host, port)) == 0


def append_runtime_log(message: str) -> None:
    stamp = datetime.now().astimezone().strftime("%Y-%m-%d %H:%M:%S")
    line = f"[{stamp}] {message}\n"
    try:
        with RUNTIME_LOG_PATH.open("a", encoding="utf-8") as handle:
            handle.write(line)
    except Exception:
        pass


def read_recent_log_lines(limit: int = 18) -> list[str]:
    if not RUNTIME_LOG_PATH.exists():
        return []
    try:
        lines = RUNTIME_LOG_PATH.read_text(encoding="utf-8", errors="ignore").splitlines()
        return lines[-limit:]
    except Exception:
        return []


def find_git_root(start: Path) -> Path | None:
    current = start.resolve()
    if current.is_file():
        current = current.parent
    for candidate in [current, *current.parents]:
        if (candidate / ".git").exists():
            return candidate
    return None


def get_git_status(start: Path) -> dict[str, str]:
    repo = find_git_root(start)
    if repo is None:
        return {"repo": "none", "branch": "n/a", "status": "No git repository found"}

    def run_git(*args: str) -> str:
        try:
            output = subprocess.check_output(["git", "-C", str(repo), *args], text=True, stderr=subprocess.DEVNULL)
            return output.strip()
        except Exception:
            return ""

    branch = run_git("branch", "--show-current") or "detached"
    status = run_git("status", "--short")
    commit = run_git("log", "-1", "--pretty=%s")
    dirty = len([line for line in status.splitlines() if line.strip()])
    return {
        "repo": str(repo),
        "branch": branch,
        "status": f"{dirty} dirty files" if dirty else "Clean working tree",
        "commit": commit or "No commit found",
    }


def get_git_status_panel(config: dict[str, Any]) -> Panel:
    git = get_git_status(Path(config.get("workspace_root", str(BASE_DIR.parent))))
    rows = [
        ("Repo", git.get("repo", "none")),
        ("Branch", git.get("branch", "n/a")),
        ("State", git.get("status", "n/a")),
    ]
    if git.get("commit"):
        rows.append(("Last Commit", git["commit"]))
    table = make_table(rows)
    return Panel(table, title="[b]Git Status[/b]", border_style="bright_blue", box=box.ROUNDED)


def check_endpoint(url: str, timeout: float = 3.0) -> dict[str, str]:
    request = urllib.request.Request(url, method="GET")
    try:
        with urllib.request.urlopen(request, timeout=timeout) as response:
            status = getattr(response, "status", 200)
            return {"status": str(status), "state": "ONLINE"}
    except Exception as exc:
        return {"status": "offline", "state": "OFFLINE", "error": str(exc)}


def get_endpoints_panel(config: dict[str, Any]) -> Panel:
    rows: list[tuple[str, str]] = []
    for endpoint in config.get("endpoints", []):
        if not isinstance(endpoint, dict):
            continue
        url = str(endpoint.get("url", ""))
        result = check_endpoint(url) if url else {"state": "OFFLINE", "status": "n/a"}
        state = "[bold green]ONLINE[/bold green]" if result.get("state") == "ONLINE" else "[bold red]OFFLINE[/bold red]"
        note = endpoint.get("note", "")
        rows.append(
            (
                endpoint.get("name", "Endpoint"),
                f"{state}  [dim]{url}[/dim]  [dim]{result.get('status', '')}[/dim]{f' - {note}' if note else ''}",
            )
        )
    table = make_table(rows or [("Endpoints", "No endpoint checks configured")])
    return Panel(table, title="[b]Endpoints[/b]", border_style="cyan", box=box.ROUNDED)


def get_logs_panel(limit: int = 18) -> Panel:
    lines = read_recent_log_lines(limit=limit)
    if not lines:
        table = make_table([("Logs", "No runtime log yet"), ("Tip", "Launch a tool to start writing events")])
        return Panel(table, title="[b]Logs[/b]", border_style="white", box=box.ROUNDED)

    text = Text("\n".join(lines), style="white")
    return Panel(text, title="[b]Logs[/b]", border_style="white", box=box.ROUNDED)


def build_layout(config: dict[str, Any], missions: list[dict[str, Any]]) -> Layout:
    layout = Layout(name="root")
    layout.split_column(
        Layout(name="header", size=4),
        Layout(name="strip", size=5),
        Layout(name="row1", ratio=1),
        Layout(name="row2", ratio=1),
        Layout(name="row3", ratio=1),
        Layout(name="launchers", size=8),
        Layout(name="footer", size=9),
    )

    layout["row1"].split_row(
        Layout(name="clock"),
        Layout(name="weather"),
        Layout(name="workspace"),
    )

    layout["row2"].split_row(
        Layout(name="cpu"),
        Layout(name="memory"),
        Layout(name="processes"),
    )

    layout["row3"].split_row(
        Layout(name="pipeline"),
        Layout(name="services"),
        Layout(name="shortcuts"),
    )

    layout["header"].update(build_header_panel(config))
    layout["strip"].split_row(
        Layout(name="alerts"),
        Layout(name="health"),
        Layout(name="status"),
    )

    layout["clock"].update(get_clock_panel(config))
    layout["weather"].update(get_weather_panel(config))
    layout["workspace"].update(get_workspace_panel(config, missions))
    layout["cpu"].update(get_cpu_panel())
    layout["memory"].update(get_memory_panel())
    layout["processes"].update(get_processes_panel())
    layout["pipeline"].update(get_pipeline_panel(missions))
    layout["services"].update(get_services_panel(config))
    layout["shortcuts"].update(get_shortcuts_panel(config))
    layout["alerts"].update(get_alerts_panel(config, missions))
    layout["health"].update(get_health_panel(config))
    layout["status"].update(get_status_panel(config, missions))
    layout["launchers"].update(get_launcher_monitor_panel(config))
    layout["footer"].update(get_telemetry_panel(missions))
    return layout


def render_snapshot(config: dict[str, Any], missions: list[dict[str, Any]]) -> Group:
    return Group(
        build_header_panel(config),
        get_alerts_panel(config, missions),
        get_health_panel(config),
        get_status_panel(config, missions),
        get_clock_panel(config),
        get_weather_panel(config),
        get_workspace_panel(config, missions),
        get_git_status_panel(config),
        get_endpoints_panel(config),
        get_logs_panel(),
        get_cpu_panel(),
        get_memory_panel(),
        get_processes_panel(),
        get_pipeline_panel(missions),
        get_services_panel(config),
        get_shortcuts_panel(config),
        get_launcher_monitor_panel(config),
        get_telemetry_panel(missions),
    )


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Terminal-first Mission Control dashboard")
    parser.add_argument("--config", default=str(DEFAULT_CONFIG_PATH), help="Path to dashboard.yaml or dashboard.json")
    parser.add_argument("--missions", default=str(DEFAULT_MISSIONS_PATH), help="Path to active_missions.json")
    parser.add_argument("--refresh", type=float, default=None, help="Refresh rate in frames per second")
    parser.add_argument("--embedded", action="store_true", help="Do not use the alternate screen buffer")
    parser.add_argument("--once", action="store_true", help="Render one frame and exit")
    return parser.parse_args()


def main() -> int:
    args = parse_args()
    config_path = Path(args.config).expanduser().resolve()
    missions_path = Path(args.missions).expanduser().resolve()
    config = load_dashboard_config(config_path)
    missions = load_missions(missions_path, config)

    refresh_hz = float(args.refresh or config.get("refresh_hz", 4.0) or 4.0)
    refresh_hz = max(0.5, refresh_hz)
    refresh_sleep = 1.0 / refresh_hz

    psutil.cpu_percent(interval=None)

    if args.once:
        console.print(render_snapshot(config, missions))
        return 0

    try:
        with Live(
            build_layout(config, missions),
            console=console,
            refresh_per_second=refresh_hz,
            screen=not args.embedded,
        ) as live:
            while True:
                time.sleep(refresh_sleep)
                live.update(build_layout(config, missions))
    except KeyboardInterrupt:
        return 0


if __name__ == "__main__":
    raise SystemExit(main())
