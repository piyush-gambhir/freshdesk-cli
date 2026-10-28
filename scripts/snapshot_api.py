#!/usr/bin/env python3
"""Extract endpoint facts from Freshdesk's public reference using only stdlib."""
import argparse
import hashlib
import json
import re
from html.parser import HTMLParser
from pathlib import Path
from urllib.request import urlopen

SOURCE = "https://developers.freshdesk.com/api/"


class Node:
    def __init__(self, tag="", attrs=(), parent=None):
        self.tag, self.attrs, self.parent = tag, dict(attrs), parent
        self.children = []

    def text(self):
        return " ".join(" ".join(c.text() if isinstance(c, Node) else c for c in self.children).split())

    def find(self, tag=None, **attrs):
        for child in self.children:
            if isinstance(child, Node):
                if (tag is None or child.tag == tag) and all(v in child.attrs.get(k, "").split() for k, v in attrs.items()):
                    yield child
                yield from child.find(tag, **attrs)


class Parser(HTMLParser):
    def __init__(self):
        super().__init__()
        self.root = self.current = Node()

    def handle_starttag(self, tag, attrs):
        n = Node(tag, attrs, self.current)
        self.current.children.append(n)
        if tag not in {"area", "base", "br", "col", "embed", "hr", "img", "input", "link", "meta", "param", "source", "track", "wbr"}:
            self.current = n

    def handle_endtag(self, tag):
        n = self.current
        while n.parent:
            if n.tag == tag:
                self.current = n.parent
                return
            n = n.parent

    def handle_data(self, data):
        self.current.children.append(data)


def extract(html):
    parser = Parser()
    parser.feed(html)
    endpoints = []
    for section in parser.root.find("div", **{"class": "api-content-main"}):
        anchor = section.parent.attrs.get("id", "")
        headings = list(section.find("h2"))
        if not headings:
            continue
        for api in section.find("div", **{"class": "api-url"}):
            label = next(api.find("i"), None)
            route = next(api.find("h6"), None)
            if label is None or route is None:
                continue
            method = label.text().upper()
            if method == "UPDATE":
                method = "PUT"
            if method == "CLONE":
                method = "GET"
            if method not in {"GET", "POST", "PUT", "DELETE", "PATCH"}:
                continue
            fields = []
            for table in section.find("table"):
                for tr in table.find("tr"):
                    cells = [c for c in tr.children if isinstance(c, Node) and c.tag == "td"]
                    if len(cells) >= 2:
                        name = cells[0].text()
                        name = re.split(r"\s|\*|†", name)[0]
                        if re.fullmatch(r"[a-zA-Z][a-zA-Z0-9_]*", name):
                            fields.append({"name": name, "type": cells[1].text(), "required": "mandatory" in cells[0].text().lower()})
            endpoints.append({"anchor": anchor, "title": headings[0].text(), "method": method, "path": route.text(), "fields": fields})
    return endpoints


def main():
    p = argparse.ArgumentParser(description=__doc__)
    p.add_argument("--html", type=Path, help="Read a previously downloaded reference")
    p.add_argument("--output", type=Path, default=Path(".scratch/endpoints.json"))
    args = p.parse_args()
    raw = args.html.read_bytes() if args.html else urlopen(SOURCE, timeout=30).read()
    endpoints = extract(raw.decode())
    snapshot = {"source": SOURCE, "source_sha256": hashlib.sha256(raw).hexdigest(), "operations": endpoints}
    args.output.parent.mkdir(parents=True, exist_ok=True)
    args.output.write_text(json.dumps(snapshot, indent=2) + "\n")
    print(f"Extracted {len(endpoints)} operations into {args.output}")


if __name__ == "__main__":
    main()
