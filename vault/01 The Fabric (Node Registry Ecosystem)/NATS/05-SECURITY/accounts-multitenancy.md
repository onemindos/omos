---
title: "Accounts & Multi-Tenancy"
section: "05-SECURITY"
type: note
tags: [nats, security, accounts, multi-tenancy, exports, imports, system-account]
source: https://docs.nats.io/running-a-nats-service/configuration/securing_nats/accounts
created: 2026-06-14
---

# Accounts & Multi-Tenancy

**Accounts** are NATS's hard isolation boundary: each account is its own **subject space**. Subjects in account A are completely invisible to account B — true multi-tenancy in a single server. Controlled sharing happens only through explicit **exports/imports**.

See also: [[NATS/05-SECURITY/authorization|Authorization]] · [[NATS/05-SECURITY/jwt-nkeys-decentralized|JWT & NKEYS]] · [[NATS/02-JETSTREAM/streams|Streams]]

---

## The Isolation Model

```
┌─────────── Account A ───────────┐   ┌─────────── Account B ───────────┐
│ subjects: orders.>, _INBOX.>     │   │ subjects: orders.>, _INBOX.>     │
│ users: a1, a2                    │   │ users: b1                        │
│ JetStream: its own streams       │   │ JetStream: its own streams       │
└──────────────────────────────────┘   └──────────────────────────────────┘
        A.orders.new  ≠  B.orders.new   (same subject string, different universe)
```

- Two accounts can both use the subject `orders.new` with **zero** crosstalk.
- Each account has its **own JetStream** assets (streams/KV/OS are per-account).
- A user belongs to exactly **one** account; its subjects are scoped to that account.

This is far stronger than per-user permissions — it's a *namespace wall*, not a filter.

---

## Defining Accounts (config)

```text
accounts {
  SYS: { users: [{ user: admin, password: admin }] }
  CORE: {
    jetstream: enabled
    users: [{ user: core, password: "..." }]
  }
  FABRIC: {
    jetstream: enabled
    users: [{ user: fabric, password: "..." }]
  }
}
system_account: SYS
```

---

## Exports & Imports — controlled sharing

To let data cross the account wall, the **owning** account **exports** a subject (stream) or service, and the **consuming** account **imports** it. Two kinds:

| Kind | Direction | Meaning |
|------|-----------|---------|
| **stream** | push | Account B exports a *stream* of messages; account A imports it to receive them |
| **service** | request-reply | Account B exports a *service* subject; account A imports it to make requests |

```text
accounts: {
  A: {
    users: [{user: a, password: a}]
    imports: [{ stream: { account: B, subject: "foo" } }]
  },
  B: {
    users: [{user: b, password: b}]
    exports: [{ stream: "foo" }]
  },
}
```

Imports can **remap** the subject locally (`to:`), so the consumer sees it under its own naming. Services use `response_type` (e.g. `stream` for streamed responses). Cross-account **JetStream** sharing (e.g. a leaf node sourcing a hub stream) is done by exporting the relevant `$JS.<domain>.API.*` and `$JS.ACK.*`/`$JS.FC.*` service subjects — see [[NATS/02-JETSTREAM/mirrors-and-sources|Mirrors & Sources]].

---

## The System Account (`$SYS`)

The **system account** is special — it receives server **events, advisories, and monitoring** on `$SYS.>`:
- `$SYS.SERVER.<id>.STATSZ` — server stats events
- `$SYS.ACCOUNT.<acct>.CONNS` — connection events
- Server discovery, account connz, JetStream cluster events

Admin tooling (`nats server report`, `nats server ls`, `nats server info`) talks to `$SYS`. Only the system account's users can access it. Always define a dedicated `SYS` account and set `system_account`.

> **OneMind pitfall (verified, critical):** the SYS user's credentials are **separate** from regular bus creds. `nats server report jetstream` fails with normal creds — it needs the **system-account** user. Keep a distinct `--user/--password` (or context) for `$SYS` operations. Also: changing `system_account` requires a **full server restart**, *not* a SIGHUP reload.

---

## `no_auth_user` and the "visibility ≠ persistence" trap

`no_auth_user` maps anonymous clients to a configured user — and therefore that user's **account**. The OneMind gotcha: an anonymous tool lands in the **GLOBAL** account and can *see* core subjects, but JetStream streams created under the **CORE** account are in a different subject space — so a message can be "seen" on the wire yet not persisted to the stream you expected, because publisher and stream are in different accounts. **Lesson:** always know which account a client is in; subject visibility does not imply it's hitting the right account's JetStream.

---

## OneMind Account Map

OneMind runs a multi-account bus (per `~/onemindos-core/docs/NATS-DOCTRINE.md`): **CORE, FABRIC, GLOBAL, INTERFACE, SAAS, VOICE** + **SYS**. Each domain agent and subsystem is scoped to its account; cross-account flows (e.g. SAAS results into CORE, VOICE events to INTERFACE) go through explicit exports/imports. This is the production reason to master accounts: the whole OneMind isolation/locality doctrine is built on them.

## Sources
- Accounts — https://docs.nats.io/running-a-nats-service/configuration/securing_nats/accounts
- Exports/imports — https://docs.nats.io/running-a-nats-service/configuration/securing_nats/accounts#exporting-and-importing
- System events — https://docs.nats.io/running-a-nats-service/configuration/sys_accounts
- Cross-account JetStream — https://docs.nats.io/running-a-nats-service/configuration/leafnodes/jetstream_leafnodes
