# Reviewed workflows

Open [Create](/create) to prepare one of four review packages:

| Workflow | Produces | Review before promotion |
| --- | --- | --- |
| GitOps service | Hardened service proposal manifests | Image, probes, resources, network and storage |
| App onboarding | Proposed catalog descriptor and checklist | Canonical URL, account context, owner and category |
| Deployment change | Pinned image change request and rollback plan | Actual GitOps owner, migrations, capacity and acceptance criteria |
| Recovery preparation | Recovery objective, evidence reference and rehearsal gates | Compatible backups, keys, isolated recovery and operator approval |

## Access and credentials

Only administrators may execute mutation workflows under the existing backend
permission policy. Machine MCP credentials cannot invoke Scaffolder. Read access
or a visible template is not execution permission.

Use a short-lived fine-grained GitHub token limited to `onemindos/omos`, with
Contents and Pull requests write, in the password-style Secret field. Never put
it in purpose, URLs, evidence, or rollback fields. It is used for this task's
publish action; it is not a shared cluster deployment credential. Follow the
existing task-retention policy and revoke a temporary token after use.

## Review packages versus deployment

Each workflow creates a **draft** PR under `proposals/`, outside the live catalog
and Argo watched deployment paths. Repeating the same workflow and proposal name
updates the same open proposal branch; use a new name after a proposal has been
closed or merged. Read the generated files before approving them.

Merging a proposal stores the review package; it does not provision an account,
register an app, deploy an image, restore a database, or change permissions.
Promotion is a separate reviewed change against the actual owning source files.
The proposal's reviewer field records accountability but does not assign or
notify a GitHub user.

Do not assume an image rollback can undo a database migration. Do not infer that
an available backup can be restored until an isolated rehearsal produces evidence.
Preserve production volumes and encryption/sealing keys. Portal backup evidence
is not identity-system recovery evidence.

## Before release

Review branch changes, build the portal container, publish it through the existing
release workflow, and pin its digest in GitOps. Confirm the Create page lists all
four templates and use a deliberately scoped token for a first real proposal.
