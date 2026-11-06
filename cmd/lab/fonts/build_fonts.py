# build_fonts.py  (run: uv run --with fonttools --with brotli build_fonts.py)
import base64, io
from fontTools import subset
from fontTools.ttLib import TTFont

def make_subset(path, unicodes):
    opts = subset.Options()
    opts.flavor = "woff2"
    opts.layout_features = ["*"]
    opts.hinting = False
    font = subset.load_font(path, opts)
    s = subset.Subsetter(opts)
    s.populate(unicodes=unicodes)
    s.subset(font)
    buf = io.BytesIO()
    subset.save_font(font, buf, opts)
    return base64.b64encode(buf.getvalue()).decode()

def cps(*parts):
    """cps('0-9', 'A-Z', ' ,.:') -> sorted codepoints. 'x-y' = range, otherwise literal chars."""
    out = set()
    for p in parts:
        if len(p) == 3 and p[1] == "-":
            out.update(range(ord(p[0]), ord(p[2]) + 1))
        else:
            out.update(map(ord, p))
    return sorted(out)

TEXT = cps("A-Z", "a-z", "0-9", " -!.,:'=")

FONTS = {
    # family name         file                          codepoints
    "Noto Music":    ("files/var/www/fonts/noto-music.ttf",       [0x1D100, 0x1D101, 0x1D102, 0x1D103, 0x1D106, 0x1D107, 0x1D110]),
    "mozart11":  ("files/var/www/fonts/mozart11.ttf",             [0xF026, 0xF02B, 0xF02D, 0xF02E, *range(0xF031, 0xF037),
                                                    0xF038, 0xF039, 0xF060, 0xF062, 0xF063, 0xF064,
                                                    0xF069, 0xF06D, 0xF074, 0xF075, 0xF076, 0xF0E2]),
    "Old Standard TT": ("files/var/www/fonts/old-standard-tt.ttf", cps("0-9", ".'")),
    "Caladea":     ("files/var/www/fonts/caladea.ttf",          TEXT),
    "Figtree":    ("files/var/www/fonts/figtree.ttf",         TEXT),
}

faces = []
for family, (path, codepoints) in FONTS.items():
    b64 = make_subset(path, codepoints)
    print(f"{family}: {len(b64) * 3 // 4} bytes")
    faces.append(f"@font-face {{ font-family:'{family}'; "
                 f"src:url(data:font/woff2;base64,{b64}) format('woff2'); }}")

css = "<style>\n" + "\n".join(faces) + "\n</style>"
open("fonts.css.txt", "w").write(css)