#!/usr/bin/env python3
# name=Image info
# desc=Report the size and header of each selected image file
# path=Image
import os
import sys

# Minimal, dependency-free header sniffing for the common formats — enough to demonstrate a
# per-file image tool without pulling in Pillow. A real GIMP/plugin launcher would set
# terminal=true and hand these paths to the external tool instead.
SIGNS = {
    b"\x89PNG\r\n\x1a\n": "PNG",
    b"\xff\xd8\xff": "JPEG",
    b"GIF87a": "GIF",
    b"GIF89a": "GIF",
    b"BM": "BMP",
}


def kind(path):
    try:
        with open(path, "rb") as f:
            head = f.read(16)
    except OSError as e:
        return f"unreadable ({e})"
    for sig, name in SIGNS.items():
        if head.startswith(sig):
            return name
    return "unknown"


def main(paths):
    if not paths:
        print("no files selected")
        return 1
    for p in paths:
        try:
            size = os.path.getsize(p)
        except OSError:
            size = -1
        print(f"{kind(p):<10} {size:>10}  {p}")
    return 0


if __name__ == "__main__":
    sys.exit(main(sys.argv[1:]))
