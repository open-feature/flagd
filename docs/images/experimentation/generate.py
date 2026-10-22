"""Generates the figures for docs/guides/experimentation.md.

Run from the repo root: python3 docs/images/experimentation/generate.py
"""
import sys
from pathlib import Path

OUT = Path(sys.argv[1]) if len(sys.argv) > 1 else Path(__file__).parent
OUT.mkdir(parents=True, exist_ok=True)

W = 680          # svg width
X0, X1 = 140, 660  # bar extent
BAR_H = 30
FONT = "font-family='-apple-system,Segoe UI,Roboto,Helvetica,Arial,sans-serif'"
MUTED = "#8a8a99"   # readable on light and dark backgrounds
COLORS = {
    "excluded": "#868e96", "out": "#868e96",
    "test": "#5d5dff", "control": "#d9480f", "in": "#2f9e44",
}


def x(f):
    return X0 + f * (X1 - X0)


def bar(y, label, arms):
    """arms: list of (name, weight). Returns svg for a labelled bar."""
    total = sum(w for _, w in arms)
    s = [f"<text x='{X0 - 12}' y='{y + BAR_H / 2 + 5}' text-anchor='end' font-size='13' fill='{MUTED}'>{label}</text>"]
    acc = 0
    for name, w in arms:
        if w == 0:
            continue
        a, b = acc / total, (acc + w) / total
        s.append(f"<rect x='{x(a):.1f}' y='{y}' width='{x(b) - x(a):.1f}' height='{BAR_H}' fill='{COLORS[name]}'/>")
        if x(b) - x(a) > 40:
            s.append(f"<text x='{(x(a) + x(b)) / 2:.1f}' y='{y + BAR_H / 2 + 5}' text-anchor='middle' font-size='13' font-weight='600' fill='#fff'>{name}</text>")
        acc += w
    # thin separators between slices
    acc = 0
    for _, w in arms[:-1]:
        acc += w
        s.append(f"<line x1='{x(acc / total):.1f}' y1='{y}' x2='{x(acc / total):.1f}' y2='{y + BAR_H}' stroke='#fff' stroke-width='2'/>")
    return "\n".join(s)


def user(f, bar_ys, label=None):
    """A user at hash position f: a small marker above each bar."""
    s = []
    for y in bar_ys:
        cx = x(f)
        s.append(f"<path d='M{cx - 6:.1f},{y - 10} L{cx + 6:.1f},{y - 10} L{cx:.1f},{y - 2} Z' fill='{MUTED}'/>")
    if label:
        s.append(f"<text x='{x(f):.1f}' y='{bar_ys[0] - 15}' text-anchor='middle' font-size='12' fill='{MUTED}'>{label}</text>")
    return "\n".join(s)


def band(a, b, y_top, y_bot, text=None, text_y=None):
    """Highlight a hash range across bars."""
    s = [f"<rect x='{x(a):.1f}' y='{y_top}' width='{x(b) - x(a):.1f}' height='{y_bot - y_top}' fill='{MUTED}' fill-opacity='0.18' stroke='{MUTED}' stroke-dasharray='4 3'/>"]
    if text:
        s.append(f"<text x='{(x(a) + x(b)) / 2:.1f}' y='{text_y}' text-anchor='middle' font-size='12' fill='{MUTED}'>{text}</text>")
    return "\n".join(s)


def axis(y, ticks):
    s = [f"<line x1='{X0}' y1='{y}' x2='{X1}' y2='{y}' stroke='{MUTED}'/>"]
    for t in ticks:
        s.append(f"<line x1='{x(t):.1f}' y1='{y}' x2='{x(t):.1f}' y2='{y + 4}' stroke='{MUTED}'/>")
        s.append(f"<text x='{x(t):.1f}' y='{y + 17}' text-anchor='middle' font-size='11' fill='{MUTED}'>{t:g}</text>")
    return "\n".join(s)


def svg(name, h, body, title):
    (OUT / name).write_text(
        f"<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 {W} {h}' width='{W}' height='{h}' {FONT} role='img'>\n"
        f"<title>{title}</title>\n{body}\n</svg>\n")


# 1. How a user gets an arm
svg("bucketing.svg", 110, "\n".join([
    bar(40, "weights 80/10/10", [("excluded", 80), ("test", 10), ("control", 10)]),
    user(0.84, [40], "a user"),
    axis(80, [0, 0.8, 0.9, 1]),
]), "Each user hashes to a fixed point; arms are slices in list order")

# 2. Flat ramp: test users become control
svg("flat-ramp-reshuffle.svg", 170, "\n".join([
    band(0.8, 0.9, 30, 140, "these users switch arms", 22),
    bar(40, "before 80/10/10", [("excluded", 80), ("test", 10), ("control", 10)]),
    bar(100, "after 20/40/40", [("excluded", 20), ("test", 40), ("control", 40)]),
    axis(150, [0, 0.2, 0.6, 0.8, 0.9, 1]),
]), "Ramping a flat fractional moves test users into control")

# 3. Gate + split: ramp the gate, the split never moves
svg("nested-ramp.svg", 240, "\n".join([
    f"<text x='{X0}' y='18' font-size='12' font-weight='600' fill='{MUTED}'>gate — seed with -enrollment-gate-</text>",
    bar(40, "before 20/80", [("in", 20), ("out", 80)]),
    bar(90, "after 40/60", [("in", 40), ("out", 60)]),
    user(0.12, [40, 90]),
    f"<text x='{X0}' y='150' font-size='12' font-weight='600' fill='{MUTED}'>split — seed with -arms-split-, never edited</text>",
    bar(172, "1/1", [("test", 1), ("control", 1)]),
    user(0.71, [172]),
    f"<text x='{X0}' y='225' font-size='12' fill='{MUTED}'>One user, two independent positions. Ramping the gate never changes the arm.</text>",
]), "Gate and split use different seeds; ramping the gate leaves the split untouched")

# 4. Same seed on both stages: control starves
svg("nested-same-seed.svg", 190, "\n".join([
    band(0, 0.2, 18, 150, None),
    f"<text x='{x(0.1):.1f}' y='168' text-anchor='middle' font-size='12' fill='{MUTED}'>enrolled users</text>",
    bar(30, "gate 20/80", [("in", 20), ("out", 80)]),
    bar(100, "split 1/1", [("test", 1), ("control", 1)]),
    f"<text x='{x(0.75):.1f}' y='92' text-anchor='middle' font-size='12' fill='{MUTED}'>same seed = same position on both lines</text>",
    f"<text x='{x(0.75):.1f}' y='168' text-anchor='middle' font-size='12' fill='{MUTED}'>control never receives anyone</text>",
]), "With the same seed on gate and split, all enrolled users land in test")
