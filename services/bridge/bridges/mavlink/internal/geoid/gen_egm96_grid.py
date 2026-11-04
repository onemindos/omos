#!/usr/bin/env python3
"""Regenerate egm96-30min.bin.gz from the GeographicLib EGM96 15-minute grid.

Source: https://downloads.sourceforge.net/project/geographiclib/geoids-distrib/egm96-15.tar.bz2
(sha256 8b1ebad1ebae0a045502d0edb9cc51553da1d3914f01e07470c11b3bed75048e, read 2026-10-29).
EGM96 is an NGA product (https://earth-info.nga.mil/); GeographicLib's PGM packaging is MIT.

The PGM stores big-endian uint16 values v, with height N = -108 + 0.003 * v metres,
1440 x 721 points, origin 90N 0E, 0.25 degree spacing (header of egm96-15.pgm).
This script keeps every second row and column (0.5 degree spacing, 361 x 720 points),
writes the same uint16 encoding gzip-compressed, and prints the bilinear
interpolation error of the 0.5 degree grid against the 15-minute grid.

Usage: python3 -I gen_egm96_grid.py geoids/egm96-15.pgm egm96-30min.bin.gz
"""
import gzip
import math
import struct
import sys

W, H, OFF, SC = 1440, 721, -108.0, 0.003


def main():
    pgm, out = sys.argv[1], sys.argv[2]
    raw = open(pgm, "rb").read()
    start = raw.index(b"65535\n") + 6
    vals = struct.unpack(">%dH" % (W * H), raw[start:start + W * H * 2])

    def bil(get, step, lat, lon):
        y = (90.0 - lat) / step
        x = (lon % 360.0) / step
        r0 = min(int(math.floor(y)), int(180 / step) - 1)
        c0 = int(math.floor(x))
        fy, fx = y - r0, x - c0
        a, b = get(r0, c0), get(r0, c0 + 1)
        c, d = get(r0 + 1, c0), get(r0 + 1, c0 + 1)
        return (1 - fy) * ((1 - fx) * a + fx * b) + fy * ((1 - fx) * c + fx * d)

    w, h = 720, 361
    sub = [vals[(r * 2) * W + (c * 2)] for r in range(h) for c in range(w)]

    def g15(r, c):
        return OFF + SC * vals[r * W + (c % W)]

    def g30(r, c):
        return OFF + SC * sub[r * w + (c % w)]

    worst, sq, n = 0.0, 0.0, 0
    for r in range(720):
        for c in range(1440):
            lat, lon = 90 - r * 0.25 - 0.125, c * 0.25 + 0.125
            e = abs(bil(g30, 0.5, lat, lon) - bil(g15, 0.25, lat, lon))
            worst, sq, n = max(worst, e), sq + e * e, n + 1
    open(out, "wb").write(gzip.compress(struct.pack(">%dH" % len(sub), *sub), 9, mtime=0))
    print("0.5 deg grid %dx%d: max error %.2f m, rms %.3f m vs 15-minute grid" % (h, w, worst, math.sqrt(sq / n)))


if __name__ == "__main__":
    main()
