# Aircraft and map website references

This batch registers external references only, not operational integrations.

| Exact requested link | Canonical Component | Home category | System |
| --- | --- | --- | --- |
| https://www.adsbexchange.com | `adsb-exchange` (existing) | Intelligence | `intelligence` |
| https://www.openstreetmap.org | `openstreetmap` | Geospatial | `geospatial` |
| https://opensky-network.org | `opensky-network` | Intelligence | `intelligence` |

The ADS-B Exchange `www` URL is a supporting link. Its original
`https://adsbexchange.com` reference, launcher, description and identity are
preserved. The original mapping regression still asserts the original URL and
launcher; its expected link list gains only the requested supporting URL.

Current-main inspection also found the historical Resources
`fabric-openstreetmap-91ae68a8` and `fabric-opensky-20188a97` in `cluster.yaml`.
These are archived documented-capability records, not Home Components or verified
services. They remain unchanged. New Home titles are disambiguated from those
records. No new Systems or workload relationships are added.

## Primary-source evidence

All three supplied homepages were retrieved successfully with `web_extract`:

- ADS-B Exchange describes an independent aircraft flight-data receiver network.
  Its existing conservative description remains unchanged; provider coverage
  and operational claims have not been adopted as OneMind facts.
- OpenStreetMap describes a world map created by people and available under an
  open license. The card describes collaborative geographic data and maps.
- OpenSky Network says “Open Air Traffic Data for Research” and lists ADS-B,
  Mode-S, ADS-C, FLARM and VHF data. It is an aircraft-tracking research reference,
  not an OpenStreetMap/geocoding service.

The exact submitted URLs are retained regardless of website redirects. Page
retrieval is source evidence only, not an API/data query, authentication check,
feed approval, coverage-completeness or flight-safety assessment. No installation,
activation, payment, cluster write or OneMind deployment is part of this change.
`external-provider-unverified` remains explicit. The existing System's
`pillar-review: approved` inherited by Intelligence references is classification
metadata, not approval of an API, integration or provider.

## Regression scope and local execution

`test-aircraft-map-references.cjs` covers all three exact links, canonical
uniqueness, prior ADS-B URL/launcher preservation, safe HTTPS destinations,
explicit Home metadata, durable placement, archived-record separation,
regeneration byte stability and source/mounted mirror parity. It evaluates the
actual source expressions used for Home eligibility/category selection; it is
**not** a mounted React-hook or rendered-UI test. It is imported additively by
`test-topology.cjs`.

Local validation uses Node v26.12.1 and scratch-only dependencies at
`/opt/data/profiles/builder/cache/scratch/aircraft-map-validation/node_modules`:
`yaml@2.9.1`, `prettier@2.8.8`, `@backstage/filter-predicates@0.1.5`. The genuine
`@backstage/cli@0.36.6` npm tarball supplies its `config/prettier.json` and package
metadata without installing the CLI or executing package scripts. These versions
come from the portal lockfile. No dependencies or lockfile changes are committed.

From the repository root, with `NODE_PATH` set to that scratch dependency path:

```sh
node --test services/onemind-portal/scripts/test-aircraft-map-references.cjs
node --test \
  services/onemind-portal/scripts/test-aircraft-map-references.cjs \
  services/onemind-portal/scripts/test-mapping-references.cjs \
  services/onemind-portal/scripts/test-tak-ecosystem.cjs \
  services/onemind-portal/scripts/test-hermes-ecosystem.cjs \
  services/onemind-portal/scripts/test-intelligence-references.cjs \
  services/onemind-portal/scripts/test-audit-catalog.cjs \
  services/onemind-portal/scripts/test-catalog-structure.cjs
node --test services/onemind-portal/scripts/test-topology.cjs
git diff --check
make validate
```

Results on this batch: new suite 9/9; focused invocation 102/102; topology
invocation 119/119, no skipped tests. Counts describe each invocation and must not
be summed as unique tests: suites overlap/import other suites.

Unchanged base `f76114dab9b6e01c6c5e6b9a4b8c83a3a8cec279` in isolated
`/opt/data/profiles/builder/cache/scratch/aircraft-map-base` passes the same
focused selection excluding the absent new suite (93/93), and topology (110/110),
with identical dependencies. Both base and branch `make validate` stop at
`yamllint: command not found` (Makefile target exits with error 127, make exits 2).
That is a local tooling blocker, not a passing full repository validation.
Initial missing scratch dependencies were resolved before these passing Node
results; no generated-config stubs or synthetic provider responses were used.

GitHub CI and PR state are reported separately at publication. Source tests do
not establish catalog ingestion, deployment or live Home visibility.
