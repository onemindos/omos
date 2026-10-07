# Inventory refresh and retirement review

The six-hour discovery workflow proposes changes through one reusable pull request. It does not deploy, mutate provider resources or directly write to the catalog API. Enable `ONEMIND_INVENTORY_REFRESH_ENABLED=true` only after configuring a dedicated self-hosted runner with read-only discovery access. GitHub Actions must permit workflow-created pull requests.

The runner requires Node 22, kubectl context `do-nyc3-omos-cloud`, Tailscale status access and GitHub repository metadata access. The default workflow token sees repositories permitted to that token, not necessarily every private organization repository. For complete private organization discovery, configure narrowly scoped GitHub metadata access on a dedicated runner and use that credential for discovery only; keep proposal writes restricted to omos. Do not install an organization-wide write token. Keep the workflow disabled until visibility is verified.

`refresh-inventory.cjs` requires successful Kubernetes, mesh and repository discovery plus curated Fabric reconciliation. On any source failure, all generated files are restored. Entity names remain unchanged. Missing resources retain their last observation, relationship targets and documentation. Three complete observations mark a retirement-review candidate; removal always requires human review. A missing result never proves resource deletion.

Review observed-at and last-checked-at annotations independently of live health. Fresh snapshots still depend on review/merge before production catalog updates. A failed scheduled run is visible in GitHub Actions; configure repository notifications for failures. Repository listings at the discovery limit fail closed instead of publishing truncated results.

Official references: [catalog integrations](https://backstage.io/docs/features/software-catalog/external-integrations/), [scheduled workflows](https://docs.github.com/en/actions/reference/workflows-and-actions/events-that-trigger-workflows#schedule).
