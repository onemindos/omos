# Test fixtures

`policy_vectors.json` holds command envelopes built and signed exactly as the policy
service (omos PR #554, `services/policy`, commit `66900872`) builds them in
`Service.publishCommandID`: `policy.Signer.Sign`, `policy.ParamsHash` and
`policy.Canonicalize`. `tests/test_command.py` and `tests/test_jcs.py` verify them with the
Python port, so the bridge and the policy service cannot drift apart silently.

They were produced by `genvec.go.txt` (kept as text so it is not built in this repo),
run from a copy of the policy module:

```
git archive origin/feat/policy-service services/policy services/schemas platform/security/openfga | tar -x -C /tmp/p
mkdir -p /tmp/p/services/policy/cmd/genvec && cp genvec.go.txt /tmp/p/services/policy/cmd/genvec/main.go
cd /tmp/p/services/policy && go run ./cmd/genvec > policy_vectors.json
```

The signing key is derived from a public label (`sha256("onemind-h3-test-vector")` as the
Ed25519 seed). It is a test vector, not a secret, and no service trusts it.
