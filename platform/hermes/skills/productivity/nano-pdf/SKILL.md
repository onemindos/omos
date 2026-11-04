---
name: nano-pdf
description: "Edit text, typos and titles on one page of a PDF with the nano-pdf CLI and a natural-language instruction. Needs the CLI and a model key that this pod does not have yet."
version: 1.1.0
author: community
license: MIT
platforms: [linux, macos, windows]
metadata:
  hermes:
    tags: [PDF, Documents, Editing, NLP, Productivity]
    homepage: https://pypi.org/project/nano-pdf/
    related_skills: [pdf, ocr-and-documents]
---

# nano-pdf

Source: the upstream Hermes skill, current form. Like `ocr-and-documents`, the
standalone skill was merged into the bundled `pdf` skill in Hermes v0.22.5
(`/opt/hermes/skills/productivity/pdf/references/nano-pdf-editing.md`); the
usage text below is that upstream text. It is kept as its own skill because the
owner asked for it by name, with a section on this pod's gaps.

Edit PDFs with natural-language instructions: point it at a page and describe the
change. For structural PDF work (merge, split, forms, watermarks, creation) use the
`pdf` skill; for text extraction from scans use `ocr-and-documents`.

## Before you use it: what this pod lacks (checked 2026-10-30, `onemind-hermes`)

- **The CLI is not installed.** `command -v nano-pdf` is empty, and so are `pip`,
  `uv` and the PDF libraries it builds on. The upstream line `uv pip install
  nano-pdf` ("already available in Hermes") does not hold on this image.
- **It needs a model key.** Upstream says the tool uses an LLM under the hood and
  needs an API key (`nano-pdf --help` shows the config). No such key is provisioned
  for this instance, and the provider must not be an Anthropic one. The tool
  sends the page to that provider, so it is a data-leaving action: do not use it
  on family, estate, financial or client documents without the owner's say-so.
  If the owner wants it, the key goes in a SealedSecret and the provider choice is
  the owner's; confirm what it calls with `--help` rather than from this text.
- **Installing is a change outside Git** (see "What this pod has" in
  `ocr-and-documents`): `python3 -m ensurepip`, then `python3 -m pip install
  --target /opt/data/pylibs nano-pdf`, and `PYTHONPATH=/opt/data/pylibs` plus the
  `bin` it creates. Ask first. The durable fix is the image or a sidecar.

If you cannot do the edit, say which of these is missing. Offer the fallbacks:
the owner edits in a PDF editor, or you produce corrected text for them to apply,
or for a document you generated yourself, fix the source and regenerate (the `pdf`
skill creates PDFs from a spec).

## Usage

```bash
nano-pdf edit <file.pdf> <page_number> "<instruction>"
```

## Examples

```bash
# Change a title on page 1
nano-pdf edit deck.pdf 1 "Change the title to 'Q3 Results' and fix the typo in the subtitle"

# Update a date on a specific page
nano-pdf edit report.pdf 3 "Update the date from January to February 2026"

# Fix content
nano-pdf edit contract.pdf 2 "Change the client name from 'Acme Corp' to 'Acme Industries'"
```

## Notes

- Page numbers may be 0-based or 1-based depending on version. If the edit hits
  the wrong page, retry with one more or one less.
- **Verify the output**: re-extract the page text and compare it to the original
  (`ocr-and-documents`), and have the owner look at it before it is sent anywhere.
  Never present an edited PDF as the original.
- Edit a copy. Write the result to a new file name; do not overwrite the source,
  and do not delete it afterwards (deletion needs "execute" or "override").
- Contracts, legal and financial documents: changing them is the owner's decision,
  not a convenience. Confirm the exact change before running it, and note it in
  your report.
- Works well for text changes; complex layout changes need another approach.
