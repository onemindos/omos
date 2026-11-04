---
name: ocr-and-documents
description: "Extract text from PDFs and scanned documents (pymupdf for text PDFs, marker-pdf for OCR). Use when read_file shows a coverage warning or the PDF is a scan."
version: 2.4.0
author: Hermes Agent
license: MIT
platforms: [linux, macos, windows]
metadata:
  hermes:
    tags: [PDF, Documents, Research, Arxiv, Text-Extraction, OCR]
    related_skills: [pdf, docx, powerpoint, nano-pdf]
---

# PDF and document text extraction

Source: the upstream Hermes skill, current form. The standalone `ocr-and-documents`
skill was merged into the bundled `pdf` skill in Hermes v0.22.5
(`/opt/hermes/skills/productivity/pdf/references/ocr-extraction.md`); the
procedure and the two scripts below are that upstream text and code. It is kept
here as its own skill because the owner asked for it by name, with a section on
what this pod can and cannot run. The Legacy-pod copy (v2.3.0, June 2026) was the
older standalone version; this one replaces it.

For DOCX: see the `docx` skill (create/edit) or use `python-docx` for structured reads.
For PPTX: see the `powerpoint` skill.
For PDF manipulation (merge, split, forms, watermarks, creation): see the `pdf` skill.
This skill covers **text extraction from PDFs and scanned documents**.

> **Coming from a `read_file` EXTRACTION COVERAGE WARNING?** `read_file` converts
> local PDFs but reads the text layer only; the warning footer lists the pages that
> yielded no text (scanned images). For a handful of pages, render and use vision
> (`pdftoppm -jpeg -r 150 -f N -l N file.pdf $TMPDIR/page`, then `vision_analyze`
> each image). For bulk OCR of many pages, use marker-pdf (Step 2).

## What this pod has (checked 2026-10-30, `onemind-hermes`)

The image ships none of the extraction tools. Verified missing: `pymupdf`,
`pymupdf4llm`, `marker-pdf`, `pypdf`, `pdfplumber`, `pypdfium2`, `reportlab`,
`python-docx`, `python-pptx`, `openpyxl`, `pytesseract`, and the binaries
`pdftoppm`/`pdftotext` (poppler), `tesseract`, `ocrmypdf`, `qpdf`, `gs`,
`pandoc`, `libreoffice`, `pip`, `uv`. Present: Python 3.14, Pillow, `curl`, `node`,
and a `pip` bundle reachable with `python3 -m ensurepip`.

What follows from that:

- **Do the cheap paths first**: `web_extract` for any document with a URL (Step 1)
  needs no local tool, and `read_file` already handles text-layer PDFs.
- **Scanned PDF with no URL**: no local OCR is available. Do not pretend. Tell
  the owner which pages are image-only and offer: (a) vision on rendered pages,
  which also needs `pdftoppm` or `pypdfium2` and is therefore blocked here, (b) the
  owner running OCR on the Mac and sending text, or (c) adding the tools to the
  image (the real fix, below).
- **Installing at runtime is a change outside Git.** The venv is root-owned and the
  pod runs as root, so `python3 -m ensurepip` then
  `python3 -m pip install --target /opt/data/pylibs pymupdf` would work (egress 443
  is open), and it persists on the retained volume. Ask the owner first, install
  only `pymupdf` (about 25 MB), never `marker-pdf` (3 to 5 GB of PyTorch and
  models; the memory-saturated pods and a 25 GiB volume are not the place), and
  set `PYTHONPATH=/opt/data/pylibs` for the call.
- **The durable fix** is a PR that adds the packages to the Hermes image or a
  sidecar job: `pymupdf`, `pymupdf4llm`, `pypdf`, `pdfplumber`, `poppler-utils`,
  and `tesseract-ocr` if local OCR is wanted. That changes the pinned image
  digest, so it is the owner's call.

## Step 1: Remote URL available?

If the document has a URL, **always try `web_extract` first**:

```
web_extract(urls=["https://arxiv.org/pdf/2402.03300"])
web_extract(urls=["https://example.com/report.pdf"])
```

This converts PDF to markdown with no local dependencies. Use local extraction
only when the file is local, `web_extract` fails, or you need batch processing.

## Step 2: Choose a local extractor

| Feature | pymupdf (~25MB) | marker-pdf (~3-5GB) |
|---------|-----------------|---------------------|
| Text-based PDF | yes | yes |
| Scanned PDF (OCR) | no | yes (90+ languages) |
| Tables | basic | high accuracy |
| Equations / LaTeX | no | yes |
| Code blocks | no | yes |
| Forms | no | yes |
| Headers/footers removal | no | yes |
| Reading order detection | no | yes |
| Images extraction | embedded | with context |
| Images to text (OCR) | no | yes |
| EPUB | yes | yes |
| Markdown output | via pymupdf4llm | native, higher quality |
| Speed | instant | ~1-14s/page CPU, ~0.2s/page GPU |

**Decision**: use pymupdf unless you need OCR, equations, forms or complex layout.

If marker is needed and the disk is short, say so:
"This document needs OCR/advanced extraction (marker-pdf), which needs about 5GB
for PyTorch and models. Options: free space, give me a URL so I can use
`web_extract`, or I can try pymupdf, which works for text PDFs but not scans or
equations."

## pymupdf (lightweight)

```bash
python3 scripts/extract_pymupdf.py document.pdf              # plain text
python3 scripts/extract_pymupdf.py document.pdf --markdown   # markdown (pymupdf4llm)
python3 scripts/extract_pymupdf.py document.pdf --tables     # tables (needs pandas)
python3 scripts/extract_pymupdf.py document.pdf --images out/  # embedded images
python3 scripts/extract_pymupdf.py document.pdf --metadata   # title, author, pages
python3 scripts/extract_pymupdf.py document.pdf --pages 0-4  # 0-indexed range
```

Inline:

```bash
python3 -c "
import pymupdf
doc = pymupdf.open('document.pdf')
for page in doc:
    print(page.get_text())
"
```

## marker-pdf (high-quality OCR)

```bash
python3 scripts/extract_marker.py --check                    # disk check first
python3 scripts/extract_marker.py scanned.pdf                # markdown, OCR included
python3 scripts/extract_marker.py document.pdf --json        # JSON with metadata
python3 scripts/extract_marker.py document.pdf --output_dir out/   # save images
```

Do not use `--use_llm` here: it sends page content to a hosted model with its own
key and provider, which this instance does not have and which must not be an
Anthropic provider. If it is wanted, route it through LiteLLM as an owner decision.
marker downloads about 2.5 GB of models to `~/.cache/huggingface/` on first use.

## arXiv papers

```
web_extract(urls=["https://arxiv.org/abs/2402.03300"])   # abstract only (fast)
web_extract(urls=["https://arxiv.org/pdf/2402.03300"])   # full paper
web_search(query="arxiv GRPO reinforcement learning 2026")
```

The bundled `arxiv` skill is usually the better entry for papers.

## Split, merge and search (pymupdf)

```python
import pymupdf
doc = pymupdf.open("report.pdf")

# Split: pages 1-5 to a new PDF
new = pymupdf.open()
for i in range(5):
    new.insert_pdf(doc, from_page=i, to_page=i)
new.save("pages_1-5.pdf")

# Merge
result = pymupdf.open()
for path in ["a.pdf", "b.pdf", "c.pdf"]:
    result.insert_pdf(pymupdf.open(path))
result.save("merged.pdf")

# Search every page
for i, page in enumerate(doc):
    if page.search_for("revenue"):
        print(f"Page {i+1}")
```

## Notes

- `web_extract` first for URLs; pymupdf is the safe default; marker only for scans,
  equations and complex layouts.
- Both scripts accept `--help`.
- Documents may be private (family, estate, financial). Extracted text stays in the
  conversation or a path the owner named. Do not paste it into a skill, memory, a
  public site or a third-party service.
- Writing output next to the source or creating files on the pod volume is fine;
  deleting the source after extraction is not (needs "execute" or "override").
