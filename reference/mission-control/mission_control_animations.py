from __future__ import annotations

from dataclasses import dataclass
from html.parser import HTMLParser
from math import cos, hypot, pi, sin
from urllib.error import URLError
from urllib.parse import urlparse
from urllib.request import Request, urlopen
from typing import Callable

from rich import box
from rich.align import Align
from rich.panel import Panel
from rich.text import Text


FRAME_RAMP = " .,:;irsXA253hMHGS#9B&@"


@dataclass(frozen=True)
class AnimationScene:
    name: str
    title: str
    description: str
    frame_count: int
    renderer: Callable[[int, int, int], list[str]]
    source_type: str = "local"
    url: str | None = None


@dataclass(frozen=True)
class RemoteAnimationSpec:
    url: str
    name: str | None = None
    title: str | None = None
    description: str | None = None


def _clamp(value: float, minimum: float = 0.0, maximum: float = 1.0) -> float:
    return max(minimum, min(maximum, value))


def _ramp_char(value: float) -> str:
    value = _clamp(value)
    index = int(value * (len(FRAME_RAMP) - 1))
    return FRAME_RAMP[index]


def _blank_canvas(width: int, height: int, fill: str = " ") -> list[list[str]]:
    return [[fill for _ in range(width)] for _ in range(height)]


def _to_lines(canvas: list[list[str]]) -> list[str]:
    return ["".join(row).rstrip() for row in canvas]


def _stamp(canvas: list[list[str]], x: int, y: int, char: str) -> None:
    if not canvas or not canvas[0]:
        return
    if 0 <= y < len(canvas) and 0 <= x < len(canvas[0]):
        canvas[y][x] = char


def _draw_text(canvas: list[list[str]], x: int, y: int, text: str) -> None:
    for offset, char in enumerate(text):
        _stamp(canvas, x + offset, y, char)


def _frame_dimensions(width: int, height: int) -> tuple[int, int]:
    return max(24, width), max(10, height)


def _short_domain(url: str | None) -> str:
    if not url:
        return "remote"
    parsed = urlparse(url)
    return parsed.netloc or url


def _slugify(text: str) -> str:
    slug = []
    for char in text.lower():
        if char.isalnum():
            slug.append(char)
        elif slug and slug[-1] != "-":
            slug.append("-")
    value = "".join(slug).strip("-")
    return value or "clip"


def _normalize_lines(text: str) -> list[str]:
    return [line.rstrip("\n\r") for line in text.replace("\r\n", "\n").replace("\r", "\n").split("\n")]


class _RemoteClipParser(HTMLParser):
    def __init__(self) -> None:
        super().__init__()
        self.title = ""
        self.description = ""
        self.pre_blocks: list[str] = []
        self.og_image = ""
        self._in_title = False
        self._in_pre = False
        self._pre_buffer: list[str] = []

    def handle_starttag(self, tag: str, attrs) -> None:
        attr_map = {key: value for key, value in attrs}
        if tag == "title":
            self._in_title = True
        elif tag == "pre":
            self._in_pre = True
            self._pre_buffer = []
        elif tag == "meta":
            key = attr_map.get("name") or attr_map.get("property") or ""
            content = attr_map.get("content") or ""
            if key == "description" and content and not self.description:
                self.description = content
            elif key == "og:description" and content:
                self.description = content
            elif key == "og:title" and content and not self.title:
                self.title = content
            elif key == "og:image" and content:
                self.og_image = content

    def handle_endtag(self, tag: str) -> None:
        if tag == "title":
            self._in_title = False
        elif tag == "pre":
            if self._in_pre:
                block = "".join(self._pre_buffer).rstrip()
                if block.strip():
                    self.pre_blocks.append(block)
            self._in_pre = False
            self._pre_buffer = []

    def handle_data(self, data: str) -> None:
        if self._in_title:
            self.title += data
        if self._in_pre:
            self._pre_buffer.append(data)


def _remote_source_configs(config: dict[str, object]) -> list[RemoteAnimationSpec]:
    animation_cfg = config.get("animations", {})
    if isinstance(animation_cfg, dict):
        raw_sources = animation_cfg.get("sources", [])
    else:
        raw_sources = []

    if not raw_sources:
        raw_sources = [
            {"url": "https://www.asciiart.eu/animations/ascii-fireworks", "name": "fireworks"},
            {"url": "https://www.asciiart.eu/animations/ascii-horizon", "name": "horizon"},
            {"url": "https://www.asciiart.eu/animations/ascii-radio-waves", "name": "radio"},
            {"url": "https://www.asciiart.eu/animations/ascii-rain-drops", "name": "rain"},
            {"url": "https://www.asciiart.eu/animations/ascii-starfield", "name": "starfield"},
            {"url": "https://www.asciiart.eu/animations/ascii-rotating-galaxy", "name": "galaxy"},
            {"url": "https://www.asciiart.eu/animations/ascii-tunnel", "name": "tunnel"},
            {"url": "https://www.asciiart.eu/animations/ascii-plasma", "name": "plasma"},
            {"url": "https://ascii.co.uk/animated-art/artist-animated-ascii-art-by-hey-duggee.html", "name": "hey-duggee"},
            {"url": "https://ascii-motion.app/community/project/5097d8dd-9b38-4786-b7dd-dcb18754584d", "name": "ascii-motion-squirrel"},
        ]

    specs: list[RemoteAnimationSpec] = []
    for entry in raw_sources:
        if isinstance(entry, str):
            specs.append(RemoteAnimationSpec(url=entry))
            continue
        if not isinstance(entry, dict):
            continue
        url = str(entry.get("url", "")).strip()
        if not url:
            continue
        specs.append(
            RemoteAnimationSpec(
                url=url,
                name=str(entry.get("name", "")).strip() or None,
                title=str(entry.get("title", "")).strip() or None,
                description=str(entry.get("description", "")).strip() or None,
            )
        )
    return specs


def _fetch_remote_html(url: str, timeout: float = 4.0) -> str:
    request = Request(url, headers={"User-Agent": "Mozilla/5.0 MissionControl/1.0"})
    with urlopen(request, timeout=timeout) as response:
        return response.read().decode("utf-8", errors="ignore")


def _build_remote_renderer(
    *,
    title: str,
    description: str,
    url: str,
    lines: list[str],
    source_type: str,
) -> Callable[[int, int, int], list[str]]:
    cleaned = [line.rstrip() for line in lines if line.rstrip()]
    if not cleaned:
        cleaned = ["(empty clip)"]

    def render(frame: int, width: int, height: int) -> list[str]:
        canvas = _blank_canvas(width, height)
        content_width = max(8, width - 4)
        content_height = max(3, height - 5)
        header = f"{source_type.upper()} // {title}"
        subhead = description or url
        footer = f"{_short_domain(url)}  frame {frame + 1}"

        _draw_text(canvas, 2, 0, header[:content_width])
        _draw_text(canvas, 2, 1, subhead[:content_width])

        max_line = max(len(line) for line in cleaned)
        if max_line > content_width:
            scroll = frame % max(1, max_line - content_width + 1)
        else:
            scroll = 0
        if len(cleaned) > content_height:
            line_offset = frame % max(1, len(cleaned) - content_height + 1)
        else:
            line_offset = 0

        for row in range(content_height):
            source_index = line_offset + row
            if source_index >= len(cleaned):
                break
            line = cleaned[source_index]
            if max_line > content_width:
                segment = line[scroll : scroll + content_width]
            else:
                segment = line
            segment = segment.ljust(content_width)[:content_width]
            _draw_text(canvas, 2, 3 + row, segment)

        spark_row = height - 2
        if spark_row > 2:
            spark_col = 2 + (frame % max(1, content_width))
            _stamp(canvas, min(width - 2, spark_col), spark_row, "•")
        _draw_text(canvas, 2, height - 1, footer[:content_width])
        return _to_lines(canvas)

    frame_count = 18 if source_type == "remote-meta" else max(16, min(42, len(cleaned) * 2))
    return render


def _remote_animation_from_html(spec: RemoteAnimationSpec, html: str) -> AnimationScene | None:
    parser = _RemoteClipParser()
    parser.feed(html)
    title = (spec.title or parser.title or spec.url).strip()
    description = (spec.description or parser.description or "").strip()
    name = spec.name or _slugify(title)

    if parser.pre_blocks:
        lines = _normalize_lines("\n\n".join(parser.pre_blocks))
        renderer = _build_remote_renderer(
            title=title,
            description=description,
            url=spec.url,
            lines=lines,
            source_type="remote-pre",
        )
        return AnimationScene(
            name=name,
            title=title,
            description=description or "live fetched ASCII clip",
            frame_count=max(18, min(48, len(lines) * 2 + 12)),
            renderer=renderer,
            source_type="remote",
            url=spec.url,
        )

    lines = [title]
    if description:
        lines.append(description)
    if parser.og_image:
        lines.append(f"preview: {parser.og_image}")
    lines.append(spec.url)
    renderer = _build_remote_renderer(
        title=title,
        description=description or "live fetched remote clip",
        url=spec.url,
        lines=lines,
        source_type="remote-meta",
    )
    return AnimationScene(
        name=name,
        title=title,
        description=description or "live fetched remote clip",
        frame_count=18,
        renderer=renderer,
        source_type="remote",
        url=spec.url,
    )


def _fireworks_frame(frame: int, width: int, height: int) -> list[str]:
    canvas = _blank_canvas(width, height)
    bursts = [
        (width // 4, height // 3, 0),
        (width // 2, height // 4, 6),
        (3 * width // 4, height // 3, 12),
    ]
    max_radius = max(4, min(width, height) // 2)
    for index, (cx, cy, offset) in enumerate(bursts):
        radius = (frame + offset) % max_radius
        thickness = 0.7 + (index % 2) * 0.25
        for y in range(height):
            for x in range(width):
                distance = hypot(x - cx, (y - cy) * 1.25)
                if abs(distance - radius) <= thickness:
                    _stamp(canvas, x, y, "*" if distance > 1.5 else "@")
                elif radius > 1 and abs(distance - radius * 0.6) <= 0.35:
                    _stamp(canvas, x, y, ".")
        if radius < 3:
            _draw_text(canvas, max(0, cx - 1), cy, "@@")
    _draw_text(canvas, 2, height - 2, "fireworks loop")
    return _to_lines(canvas)


def _horizon_frame(frame: int, width: int, height: int) -> list[str]:
    canvas = _blank_canvas(width, height)
    horizon_y = int(height * 0.65)
    sun_x = width // 2 + int(sin(frame * 0.11) * max(2, width * 0.12))
    sun_y = max(2, horizon_y - 5 + int(cos(frame * 0.09) * 2))
    for y in range(height):
        for x in range(width):
            wave = sin(x * 0.26 + frame * 0.18) * 1.3 + sin(x * 0.1 + frame * 0.07) * 0.8
            line_y = horizon_y + int(wave)
            if y == line_y:
                _stamp(canvas, x, y, "~")
            elif y > line_y and (y - line_y) < 4:
                _stamp(canvas, x, y, "=" if (x + y + frame) % 2 == 0 else ":")
            elif y < line_y and (x + y + frame) % 11 == 0:
                _stamp(canvas, x, y, ".")
    for y in range(-3, 4):
        for x in range(-5, 6):
            if x * x + y * y <= 16:
                _stamp(canvas, sun_x + x, sun_y + y, "O" if abs(x) + abs(y) < 5 else "o")
    _draw_text(canvas, 2, height - 2, "horizon loop")
    return _to_lines(canvas)


def _radio_frame(frame: int, width: int, height: int) -> list[str]:
    canvas = _blank_canvas(width, height)
    cx, cy = width // 2, height // 2
    max_radius = min(width, height) // 2 - 1
    radii = [(frame * 0.85 + offset) % (max_radius + 3) for offset in (0, 4, 8, 12)]
    for y in range(height):
        for x in range(width):
            distance = hypot(x - cx, (y - cy) * 1.45)
            if distance < 1.2:
                _stamp(canvas, x, y, "@")
            else:
                for radius in radii:
                    if abs(distance - radius) <= 0.55:
                        _stamp(canvas, x, y, "(" if x < cx else ")")
                        break
                else:
                    if abs(distance - max_radius) < 0.35 and (x + y + frame) % 7 == 0:
                        _stamp(canvas, x, y, ".")
    _draw_text(canvas, 2, height - 2, "radio waves loop")
    return _to_lines(canvas)


def _rain_frame(frame: int, width: int, height: int) -> list[str]:
    canvas = _blank_canvas(width, height)
    for x in range(width):
        phase = (frame * 2 + x * 5) % (height + 4)
        drop_y = phase - 2
        if 0 <= drop_y < height:
            _stamp(canvas, x, drop_y, "|")
        if 0 <= drop_y - 1 < height and (x + frame) % 4 == 0:
            _stamp(canvas, x, drop_y - 1, ".")
        if drop_y >= height - 1:
            for dx, char in ((-1, "\\"), (0, "_"), (1, "/")):
                _stamp(canvas, x + dx, height - 1, char)
    for x in range(0, width, 7):
        _stamp(canvas, x, height - 1, "~")
    _draw_text(canvas, 2, height - 2, "rain drops loop")
    return _to_lines(canvas)


def _starfield_frame(frame: int, width: int, height: int) -> list[str]:
    canvas = _blank_canvas(width, height)
    cx, cy = width / 2.0, height / 2.0
    max_radius = hypot(cx, cy)
    for star in range(1, 90):
        angle = (star * 0.61803398875) % 1.0 * (2 * pi)
        depth = 0.15 + ((star * 17) % 100) / 100.0 * 0.85
        speed = 0.35 + depth * 1.65
        radius = (frame * speed + star * 0.9) % max_radius
        x = int(cx + cos(angle) * radius)
        y = int(cy + sin(angle) * radius * 0.58)
        char = "." if depth < 0.35 else "*" if depth < 0.7 else "@"
        _stamp(canvas, x, y, char)
    _draw_text(canvas, 2, height - 2, "starfield loop")
    return _to_lines(canvas)


def _galaxy_frame(frame: int, width: int, height: int) -> list[str]:
    canvas = _blank_canvas(width, height)
    cx, cy = width / 2.0, height / 2.0
    max_radius = hypot(cx, cy)
    arms = 4
    for y in range(height):
        for x in range(width):
            dx = (x - cx) / max(1.0, cx)
            dy = (y - cy) / max(1.0, cy)
            radius = hypot(dx, dy)
            angle = (frame * 0.09) + (dy * 2.4) + (dx * 1.4)
            spiral = sin(angle * arms + radius * 4.6)
            glow = _clamp(1.0 - radius / 1.05)
            brightness = glow * 0.45 + (spiral + 1.0) * 0.27
            if radius <= 1.05 and brightness > 0.18:
                _stamp(canvas, x, y, _ramp_char(brightness))
    for ring in range(3):
        ring_radius = (ring + 1) * max_radius / 5.5
        for y in range(height):
            for x in range(width):
                distance = hypot(x - cx, (y - cy) * 1.18)
                if abs(distance - ring_radius) < 0.4 and (x + y + frame + ring) % 3 == 0:
                    _stamp(canvas, x, y, ".")
    _draw_text(canvas, 2, height - 2, "galaxy loop")
    return _to_lines(canvas)


def _tunnel_frame(frame: int, width: int, height: int) -> list[str]:
    canvas = _blank_canvas(width, height)
    cx, cy = width // 2, height // 2
    for y in range(height):
        for x in range(width):
            dx = abs((x - cx) / max(1, cx))
            dy = abs((y - cy) / max(1, cy))
            depth = max(dx, dy)
            bands = depth * 8 + frame * 0.18
            phase = bands - int(bands)
            if phase < 0.06:
                _stamp(canvas, x, y, "#")
            elif phase < 0.12:
                _stamp(canvas, x, y, "+")
            elif depth < 0.98 and (x + y + frame) % 11 == 0:
                _stamp(canvas, x, y, ".")
    _draw_text(canvas, 2, height - 2, "tunnel loop")
    return _to_lines(canvas)


def _plasma_frame(frame: int, width: int, height: int) -> list[str]:
    canvas = _blank_canvas(width, height)
    cx, cy = width / 2.0, height / 2.0
    for y in range(height):
        for x in range(width):
            value = (
                sin(x * 0.22 + frame * 0.18)
                + sin(y * 0.27 - frame * 0.11)
                + sin((x + y) * 0.11 + frame * 0.07)
                + sin(hypot(x - cx, y - cy) * 0.18 - frame * 0.15)
            ) / 4.0
            _stamp(canvas, x, y, _ramp_char((value + 1.0) / 2.0))
    _draw_text(canvas, 2, height - 2, "plasma loop")
    return _to_lines(canvas)


ANIMATION_SCENES: list[AnimationScene] = [
    AnimationScene("fireworks", "Fireworks", "expanding bursts", 18, _fireworks_frame),
    AnimationScene("horizon", "Horizon", "sunrise sweep", 20, _horizon_frame),
    AnimationScene("radio", "Radio Waves", "signal rings", 18, _radio_frame),
    AnimationScene("rain", "Rain Drops", "falling streaks", 16, _rain_frame),
    AnimationScene("starfield", "Starfield", "warp drift", 24, _starfield_frame),
    AnimationScene("galaxy", "Rotating Galaxy", "spiral dust", 24, _galaxy_frame),
    AnimationScene("tunnel", "Tunnel", "moving depth", 20, _tunnel_frame),
    AnimationScene("plasma", "Plasma", "wave flow", 22, _plasma_frame),
]


def animation_scene_names() -> list[str]:
    return [scene.name for scene in ANIMATION_SCENES]


def get_animation_scene(name: str) -> AnimationScene:
    for scene in ANIMATION_SCENES:
        if scene.name == name:
            return scene
    return ANIMATION_SCENES[0]


def render_animation_source(
    scene: AnimationScene,
    frame_index: int,
    *,
    paused: bool = False,
    width: int = 56,
    height: int = 16,
) -> Panel:
    width, height = _frame_dimensions(width, height)
    frame = frame_index % max(1, scene.frame_count)
    lines = scene.renderer(frame, width, height)
    body = Text("\n".join(lines), style="bold cyan")
    source_label = "REMOTE" if scene.source_type != "local" else "LOCAL"
    status = f"{source_label}  {scene.title}  [{scene.name}]"
    hint = "space pause | n next | b back"
    suffix = "PAUSED" if paused else f"frame {frame + 1}/{scene.frame_count}"
    caption = f"{status}  {suffix}  |  {hint}"
    return Panel(
        Align.left(body),
        title=f"[b]Animations[/b]  {caption}",
        border_style="magenta",
        box=box.ROUNDED,
        padding=(0, 1),
    )


def render_animation_panel(
    scene_name: str,
    frame_index: int,
    *,
    paused: bool = False,
    width: int = 56,
    height: int = 16,
) -> Panel:
    return render_animation_source(get_animation_scene(scene_name), frame_index, paused=paused, width=width, height=height)


def build_animation_catalog(config: dict[str, object]) -> list[AnimationScene]:
    remote_catalog: list[AnimationScene] = []
    for spec in _remote_source_configs(config):
        try:
            html = _fetch_remote_html(spec.url)
            source = _remote_animation_from_html(spec, html)
        except (URLError, TimeoutError, OSError, ValueError):
            source = None
        if source is not None:
            remote_catalog.append(source)

    if remote_catalog:
        return remote_catalog

    return list(ANIMATION_SCENES)


def build_local_animation_catalog() -> list[AnimationScene]:
    return list(ANIMATION_SCENES)
