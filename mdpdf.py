import os
import sys

import markdown
from weasyprint import HTML

src, dst = sys.argv[1], sys.argv[2]

with open(src, encoding="utf-8") as f:
    text = f.read()

body = markdown.markdown(text, extensions=["extra", "sane_lists"])

css = """
body { font-family: Georgia, serif; max-width: 780px; margin: 40px auto; line-height: 1.6; color: #222; }
h1, h2, h3 { font-family: Helvetica, Arial, sans-serif; }
img { max-width: 100%; display: block; margin: 24px auto; border-radius: 6px; }
code { background: #f2f2f2; padding: 2px 5px; border-radius: 3px; }
pre { background: #f2f2f2; padding: 12px; border-radius: 6px; overflow-x: auto; }
blockquote { border-left: 3px solid #ccc; margin: 0; padding-left: 16px; color: #555; }
"""

html = f"<html><head><meta charset='utf-8'><style>{css}</style></head><body>{body}</body></html>"

HTML(string=html, base_url=os.path.dirname(os.path.abspath(src))).write_pdf(dst)
