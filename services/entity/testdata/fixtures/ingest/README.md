# omos-ingest parser outputs (copied from PR #553)

These files are the golden parser outputs of omos-ingest (`services/bridge/ingest/testdata/golden/`
on branch `feat/ingest-framework`, commit `61d8bc298861a653606a57ee584c4acf75886dd7`), copied
unchanged. Each row is one `Record` (`subject`, `msgId`, `ts`, `data`, or `error`).

`internal/adapt/envelope_test.go` wraps each row in the bus envelope exactly as omos-ingest's
`NewEnvelope` does (`{v:1, id:StableUUID(msgId), ts, type:<subject>, source:"ingest.<id>",
subject, data, msgId}`) and checks the entity adapter's output against
`testdata/golden/intel-envelope.json`. Rows with `error` are dead-lettered by omos-ingest and never
published, so the test skips them; the malformed rows with `data` would fail omos-ingest's schema
validation, and the adapter must still reject or skip them on its own.

When #553 changes its parsers, copy the new files here and run
`go test ./internal/adapt -run Envelope -update`.
