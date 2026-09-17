# Chainguard: Software Supply Chain Security

**Type:** System Architecture / DevOps  
**Category:** Security, Containerization  
**Created:** 2026-01-20  
**Status:** Reference  
**Website:** [chainguard.dev](https://chainguard.dev)

---

## Overview

**Chainguard** is a software supply chain security company that builds "secure-by-default" open source artifacts. They provide a suite of minimal, hardened container images, libraries, and tools designed to eliminate vulnerabilities (CVEs) before they reach your production environment.

**Core Promise:**

> "Stop wasting time patching CVEs. Start with zero-CVE images."

---

## The Problem: "CVE Fatigue"

Most developers build apps on top of standard base images (like `node:latest`, `python:3.9`, or `ubuntu`).

- **Bloat:** These images contain hundreds of unnecessary packages (bash, package managers, shells).
- **Vulnerabilities:** A standard image often comes with 50-500+ known vulnerabilities (CVEs) out of the box.
- **Toil:** DevOps teams waste thousands of hours scanning, triaging, and patching non-critical bugs.

---

## The Solution: Chainguard Images ("Wolfi")

Chainguard builds images based on **Wolfi**, their proprietary (but open-source based) Linux distribution designed specifically for containers.

### Key Features

1.  **Zero-CVE Promise:**
    Images are rebuilt daily from upstream sources. If a vulnerability is found, they patch it immediately.
    - _Result:_ Most Chainguard images have 0 CVEs.

2.  **Distroless / Minimalist:**
    Images contain _only_ what is needed to run the application.
    - No shell (`/bin/sh` or `/bin/bash`)
    - No package manager (`apk`, `apt`) in the final runtime
    - Significantly reduced attack surface.

3.  **Hardenened Security:**
    - All images are signed (Sigstore).
    - SBOMs (Software Bill of Materials) generated at build time.
    - FIPS-verified versions available for compliance (Gov/FedRAMP).

4.  **Daily Rebuilds:**
    They don't patch old binaries; they rebuild the entire universe every single day.

---

## Technical Architecture

### 1. Wolfi (The OS)

Wolfi is the "undistro" — a Linux distribution designed for container environments, not hardware. It uses the `apk` package format (like Alpine) but is glibc-compatible (unlike Alpine), making it compatible with most standard software.

### 2. Melange (The Builder)

A declarative build tool that creates APK packages from source. It replaces complex Makefiles or Dockerfiles with clean YAML configurations.

### 3. Apko (The Assembler)

A tool that takes APK packages and composes them into a commercially valid OCI container image _without_ needing `docker build`. It generates the image deterministically and outputs an SBOM.

---

## Comparison

### Standard vs. Chainguard

| Feature          | Standard Image (`python:3.11`) | Chainguard Image (`cgr.dev/chainguard/python`) |
| :--------------- | :----------------------------- | :--------------------------------------------- |
| **Base OS**      | Debian/Ubuntu                  | Wolfi (Distroless)                             |
| **Size**         | ~300MB+                        | ~50MB                                          |
| **CVE Count**    | Often 100+                     | Usually 0                                      |
| **Shell Access** | Yes (`/bin/bash`)              | No (Harder for attackers)                      |
| **Root User**    | Often default                  | Non-root by default                            |
| **Updates**      | Irregular                      | Daily                                          |

---

## How to Use It

### 1. Public Images (Developer Tier)

Chainguard offers a free tier of images for development (`:latest-dev`).

- **Warning:** The `:latest` tag in the free tier pulls from the _developer_ registry, which might have CVEs if not updated frequently.
- **Production:** Companies pay to get access to the pinned, signed, zero-CVE production registry.

**Example `Dockerfile` Replacement:**

_Old:_

```dockerfile
FROM python:3.11
COPY . /app
RUN pip install -r requirements.txt
CMD ["python", "app.py"]
```

_New (Chainguard):_

```dockerfile
FROM cgr.dev/chainguard/python:latest
COPY . /app
RUN pip install -r requirements.txt
CMD ["python", "app.py"]
```

_(Note: Since there is no shell/package manager, complex installs might require a multi-stage build)._

### 2. Multi-Stage Build Strategy

Since the runtime image has no tools, you use a "build" image to compile, then copy artifacts to the "run" image.

```dockerfile
# Builder Stage
FROM cgr.dev/chainguard/python:latest-dev as builder
WORKDIR /app
COPY requirements.txt .
RUN pip install --user -r requirements.txt

# Runtime Stage
FROM cgr.dev/chainguard/python:latest
WORKDIR /app
COPY --from=builder /home/nonroot/.local /home/nonroot/.local
COPY app.py .
ENV PATH="/home/nonroot/.local/bin:$PATH"
CMD ["python", "app.py"]
```

---

## Relevance to OneMind-Codex

**Chainguard aligns with the "Cyberpunk/Solarpunk" philosophy:**

- **Cyberpunk (High Tech/Defense):** Hardened, minimal, secure-by-default architecture. It assumes the network is hostile and minimizes surface area.
- **Solarpunk (Efficiency/Regeneration):** Daily rebuilds ensure the ecosystem is always fresh ("regenerated") rather than patching old, decaying systems.

**Recommendation:**
For **Sovereign Operations (UI)** services:

- Use Chainguard images for all production containers.
- Eliminates "noise" in security scanning.
- Ensures extreme stability and authorized-only code execution.

---

**Tags:** `#security` `#containers` `#devops` `#supply-chain` `#wolfi`

**Related Docs:**

- `system-container-strategy.md` (Coming soon)
- `process-vulnerability-management.md` (Coming soon)


---

## Related

- [[00-24 SO (Sovereign Operations)/04 Security (Cyber & Digital Shield) (The Shield)/04 Security|04 Security]]
