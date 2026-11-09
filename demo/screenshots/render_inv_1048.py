#!/usr/bin/env python3
"""Render demo/screenshots/inv-1048-scale.png with the stdlib only.

The CLI demo reads this file by path. The words that must stay readable
are PLAN, SCALE, and $249.00 (Scale invoice inv_1048).
"""

import struct
import zlib
from pathlib import Path

# 5x7 glyphs, row-major, '#' is ink. Public-domain style bitmap, drawn for this demo.
GLYPHS = {
    " ": [".....", ".....", ".....", ".....", ".....", ".....", "....."],
    "$": ["..#..", ".####", "#.#..", ".###.", "..#.#", "####.", "..#.."],
    ".": [".....", ".....", ".....", ".....", ".....", "..#..", "..#.."],
    "_": [".....", ".....", ".....", ".....", ".....", ".....", "#####"],
    "0": [".###.", "#...#", "#..##", "#.#.#", "##..#", "#...#", ".###."],
    "1": ["..#..", ".##..", "..#..", "..#..", "..#..", "..#..", "#####"],
    "2": [".###.", "#...#", "....#", "..##.", ".#...", "#....", "#####"],
    "4": ["...#.", "..##.", ".#.#.", "#..#.", "#####", "...#.", "...#."],
    "8": [".###.", "#...#", "#...#", ".###.", "#...#", "#...#", ".###."],
    "9": [".###.", "#...#", "#...#", ".####", "....#", "#...#", ".###."],
    "A": [".###.", "#...#", "#...#", "#####", "#...#", "#...#", "#...#"],
    "B": ["####.", "#...#", "#...#", "####.", "#...#", "#...#", "####."],
    "C": [".###.", "#...#", "#....", "#....", "#....", "#...#", ".###."],
    "D": ["####.", "#...#", "#...#", "#...#", "#...#", "#...#", "####."],
    "E": ["#####", "#....", "#....", "####.", "#....", "#....", "#####"],
    "G": [".###.", "#...#", "#....", "#.###", "#...#", "#...#", ".###."],
    "H": ["#...#", "#...#", "#...#", "#####", "#...#", "#...#", "#...#"],
    "I": ["#####", "..#..", "..#..", "..#..", "..#..", "..#..", "#####"],
    "L": ["#....", "#....", "#....", "#....", "#....", "#....", "#####"],
    "M": ["#...#", "##.##", "#.#.#", "#...#", "#...#", "#...#", "#...#"],
    "N": ["#...#", "##..#", "#.#.#", "#..##", "#...#", "#...#", "#...#"],
    "O": [".###.", "#...#", "#...#", "#...#", "#...#", "#...#", ".###."],
    "P": ["####.", "#...#", "#...#", "####.", "#....", "#....", "#...."],
    "R": ["####.", "#...#", "#...#", "####.", "#.#..", "#..#.", "#...#"],
    "S": [".####", "#....", "#....", ".###.", "....#", "....#", "####."],
    "T": ["#####", "..#..", "..#..", "..#..", "..#..", "..#..", "..#.."],
    "U": ["#...#", "#...#", "#...#", "#...#", "#...#", "#...#", ".###."],
    "V": ["#...#", "#...#", "#...#", "#...#", "#...#", ".#.#.", "..#.."],
    "W": ["#...#", "#...#", "#...#", "#.#.#", "#.#.#", "##.##", "#...#"],
    "Y": ["#...#", "#...#", ".#.#.", "..#..", "..#..", "..#..", "..#.."],
}


def write_png(path: Path, width: int, height: int, rgb: bytearray) -> None:
    raw = bytearray()
    stride = width * 3
    for y in range(height):
        raw.append(0)
        raw.extend(rgb[y * stride : (y + 1) * stride])

    def chunk(tag: bytes, data: bytes) -> bytes:
        return struct.pack(">I", len(data)) + tag + data + struct.pack(">I", zlib.crc32(tag + data) & 0xFFFFFFFF)

    ihdr = struct.pack(">IIBBBBB", width, height, 8, 2, 0, 0, 0)
    png = b"\x89PNG\r\n\x1a\n" + chunk(b"IHDR", ihdr) + chunk(b"IDAT", zlib.compress(bytes(raw), 9)) + chunk(b"IEND", b"")
    path.write_bytes(png)


def main() -> None:
    width, height = 960, 540
    bg = (18, 20, 26)
    card = (28, 36, 48)
    bar = (36, 48, 68)
    ink = (244, 247, 250)
    muted = (168, 180, 196)
    accent = (125, 255, 176)
    rgb = bytearray(bg * (width * height))

    def pix(x: int, y: int, color: tuple[int, int, int]) -> None:
        if 0 <= x < width and 0 <= y < height:
            i = (y * width + x) * 3
            rgb[i : i + 3] = bytes(color)

    def rect(x0: int, y0: int, x1: int, y1: int, color: tuple[int, int, int]) -> None:
        for y in range(y0, y1):
            for x in range(x0, x1):
                pix(x, y, color)

    def text(s: str, x: int, y: int, scale: int, color: tuple[int, int, int]) -> None:
        cursor = x
        for ch in s:
            glyph = GLYPHS[ch]
            for row, bits in enumerate(glyph):
                for col, bit in enumerate(bits):
                    if bit == "#":
                        rect(
                            cursor + col * scale,
                            y + row * scale,
                            cursor + (col + 1) * scale - 1,
                            y + (row + 1) * scale - 1,
                            color,
                        )
            cursor += 6 * scale

    rect(40, 36, 920, 504, card)
    rect(40, 36, 920, 92, bar)
    rect(64, 56, 84, 72, (255, 95, 86))
    rect(92, 56, 112, 72, (255, 189, 46))
    rect(120, 56, 140, 72, (39, 201, 63))
    text("LEDGERLY", 168, 50, 4, ink)
    text("INVOICE INV_1048", 72, 120, 4, muted)
    text("BRIGHTWELL LABS", 72, 168, 4, ink)
    text("PLAN", 72, 240, 3, muted)
    text("SCALE", 72, 276, 8, ink)
    text("AMOUNT", 72, 360, 3, muted)
    text("$249.00", 72, 396, 8, accent)
    text("DISPUTED", 700, 128, 3, (255, 180, 120))

    out = Path(__file__).with_name("inv-1048-scale.png")
    write_png(out, width, height, rgb)
    print(out)


if __name__ == "__main__":
    main()
