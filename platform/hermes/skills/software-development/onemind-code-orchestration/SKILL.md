---
name: onemind-code-orchestration
description: "Use when orchestrating coding tasks, worker agents or PR workflows across OneMind repos."
version: 2.0.0
author: Legacy
metadata:
  hermes:
    tags: [coding, orchestration, pr, code-review, worktree]
    related_skills: [github-pr-workflow, github-code-review, claude-code, codex]
---

# OneMind Code Orchestration

This is Legacy's **SO / Forge** mode (see `onemind-domains`). Legacy is the single
orchestrator; workers are short-lived coding sub-agents, not separate personas.

## Workflow

1. Read the assigned work (GitHub Issues/Projects or the owner's request).
2. Dispatch each worker into its own isolated git worktree with one explicit
   scope: exact PR head and base, allowed writes, test commands, evidence path,
   and a prohibition on merging unless the owner designated that worker as merger.
3. Workers implement, fix and test. Review the PRs: code quality, tests, UI.
4. Merge only under the owner's current authorization. With explicit merge
   authority, mark a PR ready only after its final-source CI and review pass;
   otherwise label it for owner review.
5. One designated lead executes final merges and verifies external state.
   Parallelize review, fixes and testing, not competing merges. Report exact
   merged SHAs and aggregate genuine blockers.

## Repos

- `onemindos/omos`: the canonical consolidated repository (GitOps, `platform/`,
  `docs/`, `services/`). `platform/hermes/` holds this agent's definition.
- The older `omos-legacy` and `omos-infra` repositories are archived or stale;
  do not task workers against them. Verify any other repository name and
  responsibility from live remotes and the active request.

## Rules

- Use worktrees; never push directly to `main`; every PR has a clear description
  linking the issue.
- Use the worker and reviewer capabilities currently configured; do not assume
  historical provider model names are still available.
- Recheck live PR state before assigning preparation. If another harness already
  merged a candidate, verify read-only instead of redoing the work.
- Respect existing harness ownership: if coordination records show a new owner
  instruction, run collision checks before tasking a second editor on a branch.
- Worker reports are self-reports: independently verify PR state, the exact CI
  source, the merge SHA and main ancestry before claiming completion.
- While background workers run, finish independent checks, then yield; do not
  poll artifacts or transcripts to wait.
- Never ship broken code: tests pass first.
- Never squash merges. Never run `terraform apply`. Cluster changes are PRs to
  `platform/`, applied by Argo CD, not `kubectl apply` or `helm install`.
- Deleting a branch, worktree, volume or any resource needs the owner to write
  "execute" or "override" for that specific deletion.
