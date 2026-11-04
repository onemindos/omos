# Karakeep: key only, no gateway wiring yet

`@karakeep/mcp` (ghcr.io/karakeep-app/karakeep-mcp:0.34.1) speaks **stdio only**. The gateway image has no node and must not run stdio, and
`mcp-proxy` (ghcr.io/sparfenyuk/mcp-proxy v0.14.0) is a Python image without node, so no existing image turns it into streamable HTTP. A
combined image has to be built (node + the pinned `@karakeep/mcp` + `mcp-proxy`), pushed to GHCR and pinned by digest first.

Until then the keys session can seal `gateway-target-karakeep` (key `api_key`), but there is no target fragment, so nothing uses it.
Prepared for that day (`tools/list` of 0.34.1, read from the pinned image 2026-11-04): `classification.yaml` here, to be copied to
`policy/classification/karakeep.yaml` together with a `targets/karakeep.yaml` fragment (`groups: [work]`).
