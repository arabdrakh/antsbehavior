import sys
import os
import json
import time
import re
import math

RESET  = "\033[0m"
BOLD   = "\033[1m"
YELLOW = "\033[33m"
GREEN  = "\033[32m"
CYAN   = "\033[36m"
DIM    = "\033[2m"
ORANGE = "\033[38;5;208m"

CLEAR_SCREEN = "\033[2J\033[H"


def parse_text_input(lines: list[str]) -> dict:
    """Parse lem-in stdout (not JSON) into the same farm dict shape."""
    farm = {
        "num_ants": 0, "start": None, "end": None,
        "rooms": [], "links": [], "turns": [],
    }
    next_is_start = next_is_end = in_moves = False
    ant_line_re = re.compile(r'^L\d+-\S+')

    for raw in lines:
        line = raw.strip()
        if not line:
            continue
        if ant_line_re.match(line):
            in_moves = True
            farm["turns"].append(line.split())
            continue
        if in_moves:
            continue
        if line == "##start":
            next_is_start = True; continue
        if line == "##end":
            next_is_end = True; continue
        if line.startswith("#"):
            continue
        if farm["num_ants"] == 0 and re.match(r'^\d+$', line):
            farm["num_ants"] = int(line); continue
        if "-" in line and " " not in line:
            a, b = line.split("-", 1)
            if a and b:
                farm["links"].append({"from": a, "to": b})
            continue
        parts = line.split()
        if len(parts) == 3:
            try:
                name, x, y = parts[0], int(parts[1]), int(parts[2])
                farm["rooms"].append({"name": name, "x": x, "y": y})
                if next_is_start:
                    farm["start"] = name; next_is_start = False
                if next_is_end:
                    farm["end"] = name; next_is_end = False
            except ValueError:
                pass
    return farm


def iso_project(x: int, y: int, z: float, scale: float = 2.0) -> tuple[int, int]:
    """Classic isometric projection: 2D grid coords -> screen (col, row)."""
    col = (x - y) * scale
    row = (x + y) * (scale * 0.5) - z
    return col, row


def build_frame(farm: dict, ant_rooms: dict, turn_num: int, total: int) -> list[str]:
    W, H = 84, 30
    grid = [[" "] * W for _ in range(H)]

    def put(row, col, ch, color=""):
        r, c = int(row), int(col)
        if 0 <= r < H and 0 <= c < W:
            grid[r][c] = (color + ch + RESET) if color else ch

    def put_str(row, col, s, color=""):
        for i, ch in enumerate(s):
            put(row, col + i, ch, color)

    rooms = farm["rooms"]
    if not rooms:
        return ["No rooms to display."]

    xs = [r["x"] for r in rooms]
    ys = [r["y"] for r in rooms]
    minx, maxx = min(xs), max(xs)
    miny, maxy = min(ys), max(ys)
    rx = max(maxx - minx, 1)
    ry = max(maxy - miny, 1)

    # Normalize room coords into a small grid (0..16) before iso projection
    positions = {}
    for r in rooms:
        nx = (r["x"] - minx) / rx * 16
        ny = (r["y"] - miny) / ry * 16
        positions[r["name"]] = (nx, ny)

    origin_col = W // 2
    origin_row = 6

    screen_pos = {}
    for name, (nx, ny) in positions.items():
        z = 3 if name in (farm["start"], farm["end"]) else 1
        col, row = iso_project(nx, ny, z)
        screen_pos[name] = (origin_row + row, origin_col + col)

    # Draw tunnels (simple stepped line between iso points)
    for link in farm["links"]:
        a, b = link["from"], link["to"]
        if a not in screen_pos or b not in screen_pos:
            continue
        r1, c1 = screen_pos[a]
        r2, c2 = screen_pos[b]
        steps = int(max(abs(r2 - r1), abs(c2 - c1), 1))
        for s in range(1, steps):
            t = s / steps
            rr = r1 + (r2 - r1) * t
            cc = c1 + (c2 - c1) * t
            put(rr, cc, "·", DIM)

    # Draw rooms as little "pillars" (3D feel: a vertical tick + label)
    for name, (row, col) in screen_pos.items():
        is_start = name == farm["start"]
        is_end = name == farm["end"]
        color = CYAN if is_start else (GREEN if is_end else "")
        # base
        put(row, col, "▲" if (is_start or is_end) else "·", color or DIM)
        label = name if (is_start or is_end or len(farm["rooms"]) <= 14) else ""
        if label:
            put_str(row + 1, col - len(label) // 2, label, color or DIM)

    # Draw ants, grouped per room
    by_room = {}
    for ant_id, room in ant_rooms.items():
        by_room.setdefault(room, []).append(ant_id)

    for room, ids in by_room.items():
        if room not in screen_pos:
            continue
        row, col = screen_pos[room]
        label = ",".join(str(i) for i in sorted(ids)[:4])
        if len(ids) > 4:
            label += "…"
        put_str(row - 1, col - len(label) // 2, label, ORANGE + BOLD)

    frame = []
    progress = int((turn_num / max(total, 1)) * 36)
    bar = "█" * progress + "░" * (36 - progress)
    frame.append(f"{BOLD}{CYAN}lem-in 3D (isometric){RESET}  turn {BOLD}{turn_num}{RESET}/{total}  [{YELLOW}{bar}{RESET}]")
    frame.append(DIM + "─" * W + RESET)
    for row in grid:
        frame.append("".join(row))
    frame.append(DIM + "─" * W + RESET)
    frame.append(f"  {CYAN}▲{RESET} start [{farm['start']}]   {GREEN}▲{RESET} end [{farm['end']}]   {ORANGE}numbers{RESET} = ants   {DIM}·{RESET} tunnel/room")
    return frame


def animate(farm: dict, delay: float) -> None:
    ant_rooms = {i + 1: farm["start"] for i in range(farm["num_ants"])}
    total = len(farm["turns"])

    def show(turn_num):
        frame = build_frame(farm, ant_rooms, turn_num, total)
        sys.stdout.write(CLEAR_SCREEN)
        sys.stdout.write("\n".join(frame) + "\n")
        sys.stdout.flush()

    show(0)
    time.sleep(delay)

    for idx, moves in enumerate(farm["turns"], start=1):
        for token in moves:
            m = re.match(r'^L(\d+)-(.+)$', token)
            if m:
                ant_rooms[int(m.group(1))] = m.group(2)
        show(idx)
        time.sleep(delay)

    print(f"\n{GREEN}{BOLD}All ants reached [{farm['end']}]!{RESET}")


def main():
    delay = 0.4
    if "--fast" in sys.argv:
        delay = 0.1
    elif "--slow" in sys.argv:
        delay = 0.8

    json_arg = None
    for arg in sys.argv[1:]:
        if arg.endswith(".json"):
            json_arg = arg
            break

    if json_arg:
        if not os.path.exists(json_arg):
            print(f"ERROR: file not found: {json_arg}", file=sys.stderr)
            sys.exit(1)
        try:
            with open(json_arg) as f:
                farm = json.load(f)
        except json.JSONDecodeError as e:
            print(f"ERROR: invalid JSON in {json_arg}: {e}", file=sys.stderr)
            sys.exit(1)
    else:
        if sys.stdin.isatty():
            print(
                "ERROR: no input provided. Either pipe lem-in output into this "
                "script, or pass a colony.json path as an argument.",
                file=sys.stderr,
            )
            sys.exit(1)
        lines = sys.stdin.read().splitlines()
        if not lines:
            print("ERROR: received empty input on stdin", file=sys.stderr)
            sys.exit(1)
        farm = parse_text_input(lines)

    if not farm.get("rooms"):
        print("ERROR: could not find any rooms in the input", file=sys.stderr)
        sys.exit(1)
    if farm.get("num_ants", 0) == 0:
        print("ERROR: could not find a valid ant count in the input", file=sys.stderr)
        sys.exit(1)
    if not farm.get("start"):
        print("ERROR: could not find ##start room in the input", file=sys.stderr)
        sys.exit(1)
    if not farm.get("end"):
        print("ERROR: could not find ##end room in the input", file=sys.stderr)
        sys.exit(1)

    try:
        animate(farm, delay)
    except KeyboardInterrupt:
        sys.stdout.write(RESET + "\n")


if __name__ == "__main__":
    main()