#!/usr/bin/env python3
# name=Word count
# desc=Count lines, words, and bytes in each selected file
# path=Text
import sys


def main(paths):
    if not paths:
        print("no files selected")
        return 1
    total_lines = total_words = total_bytes = 0
    for p in paths:
        try:
            with open(p, "rb") as f:
                data = f.read()
        except OSError as e:
            print(f"skip {p}: {e}")
            continue
        lines = data.count(b"\n")
        words = len(data.split())
        n = len(data)
        total_lines += lines
        total_words += words
        total_bytes += n
        print(f"{lines:>7} {words:>7} {n:>9}  {p}")
    print(f"{total_lines:>7} {total_words:>7} {total_bytes:>9}  total")
    return 0


if __name__ == "__main__":
    sys.exit(main(sys.argv[1:]))
