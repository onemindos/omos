# Developer Infrastructure Audit — Checklist

Use this when the owner asks "is it organized?" or "make it 100% dev-ready" for any OneMind repo.

---

## Step 1: Identify Repo Type

| Type | Signal | Tooling Profile |
|------|--------|-----------------|
| **Frontend (React/Vite)** | `package.json`, `src/`, `vite.config.ts` | ESLint 9, Prettier, Vitest, Husky+lint-staged, CI |
| **IaC / GitOps** | `terraform/`, `platform/`, K8s YAML, Helm | yamllint, shellcheck, terraform fmt, pre-commit, Makefile, CI |
| **Python service** | `pyproject.toml`, `requirements.txt` | ruff, pytest, pre-commit, Makefile, CI |
| **Go service** | `go.mod` | golangci-lint, go test, Makefile, CI |

---

## Frontend (React/Vite) Repos

### What to Check

| Category | Expected | If Missing |
|----------|----------|------------|
| **ESLint** | `eslint.config.js` (flat config, v9+) | Create with TS + React + hooks rules |
| **Prettier** | `.prettierrc` + `.prettierignore` | Add (semi:false, singleQuote, trailing commas) |
| **Tests** | `vitest.config.ts` + `src/test/setup.ts` | Add Vitest + Testing Library + jsdom |
| **Pre-commit** | `.husky/pre-commit` + lint-staged in package.json | Add husky + lint-staged |
| **CI quality gate** | Lint + typecheck + test BEFORE build/deploy | Add quality job to workflow |
| **CI branch trigger** | Triggers on actual default branch | Fix (check: is it `main` or `feat/*`?) |
| **README.md** | Quick start, scripts table, arch overview, env vars | Write one |
| **Error boundary** | App-level crash recovery | Add `ErrorBoundary` in main.tsx |
| **Shared types** | `src/shared/types/index.ts` | Create for cross-module types |
| **TypeScript strict** | `strict: true` in tsconfig | Enable |
| **Path aliases** | `@/*` → `./src/*` | Add to tsconfig + vite config |
| **Scripts** | `validate`, `format`, `lint:fix`, `test:watch`, `test:coverage` | Add to package.json |

### Standard package.json Scripts

```json
{
  "dev": "vite",
  "build": "tsc -b && vite build",
  "preview": "vite preview",
  "typecheck": "tsc --noEmit",
  "lint": "eslint src/",
  "lint:fix": "eslint src/ --fix",
  "format": "prettier --write 'src/**/*.{ts,tsx,css}'",
  "format:check": "prettier --check 'src/**/*.{ts,tsx,css}'",
  "test": "vitest run",
  "test:watch": "vitest",
  "test:coverage": "vitest run --coverage",
  "validate": "npm run typecheck && npm run lint && npm run test",
  "prepare": "husky"
}
```

### Standard devDependencies to Add

```
@eslint/js, typescript-eslint, eslint-plugin-react, eslint-plugin-react-hooks
prettier, husky, lint-staged
vitest, @vitest/coverage-v8, @testing-library/react, @testing-library/jest-dom, jsdom
```

### CI Workflow Pattern (GitHub Actions)

```yaml
jobs:
  quality:          # ← runs first, fast fail
    steps: [checkout, setup-node, npm ci, typecheck, lint, format:check, test]
  build-push:       # ← only on push (not PR), needs quality
    steps: [checkout, ghcr-login, docker-build-push]
  deploy:           # ← only on default branch push, needs build-push
    steps: [doctl, kubectl rollout restart]
```

---

## IaC / GitOps Repos

### What to Check

| Category | Expected | If Missing |
|----------|----------|------------|
| **EditorConfig** | `.editorconfig` (2-space, LF, UTF-8) | Create |
| **yamllint** | `.yamllint.yml` (200 char line, 2-space, no doc-start) | Create |
| **pre-commit** | `.pre-commit-config.yaml` (yaml check, yamllint, shellcheck, tf fmt/validate, key detection) | Create |
| **Makefile** | `make lint`, `make fmt`, `make plan`, `make validate`, `make install` | Create |
| **CI quality gate** | yamllint + shellcheck + terraform validate + secret scan | Add workflow |
| **README.md** | Architecture, layout, workflow, adding-a-service pattern, namespaces | Write one |
| **CONTRIBUTING.md** | Code standards, branch strategy, review checklist, NATS grammar | Write one |
| **Secret scanning** | CI checks for private keys, hardcoded passwords | Add to workflow |

### Standard Makefile Targets

```makefile
lint: lint-yaml lint-tf lint-shell    # All linters
lint-yaml:                            # yamllint platform/ docs/
lint-tf:                              # terraform fmt -check + validate
lint-shell:                           # shellcheck all .sh files
fmt:                                  # terraform fmt -recursive
validate: lint                        # Full gate (same as CI)
plan:                                 # terraform plan
apply:                                # terraform apply
install:                              # Deploy platform (install.sh)
help:                                 # Show targets
```

### pre-commit Config Hooks

```yaml
repos:
  - repo: pre-commit-hooks          # trailing-whitespace, end-of-file, check-yaml, detect-private-key
  - repo: yamllint                   # yamllint with .yamllint.yml config
  - repo: shellcheck-py              # shellcheck -x -s bash
  - repo: pre-commit-terraform       # terraform_fmt + terraform_validate
```

### CI Workflow Pattern (GitHub Actions)

```yaml
jobs:
  lint:
    steps:
      - checkout
      - Install yamllint (pip)
      - Lint YAML
      - ShellCheck (ludeeus/action-shellcheck)
      - Setup Terraform
      - terraform fmt -check
      - terraform init -backend=false && terraform validate
      - Secret scan (grep for private keys + hardcoded passwords)
```

### yamllint Config

```yaml
extends: default
rules:
  line-length: { max: 200, allow-non-breakable-inline-mappings: true }
  truthy: { check-keys: false }
  comments: { min-spaces-from-content: 1 }
  document-start: disable
  indentation: { spaces: 2, indent-sequences: consistent }
  empty-lines: { max: 2 }
```

---

## Pitfalls (All Repo Types)

1. **ESLint 9 uses flat config** — do NOT use `.eslintrc.*` or `--ext` flag (v8 syntax)
2. **lint-staged config** goes in `package.json` under `"lint-staged"` key (no separate file needed)
3. **Husky v9** only needs `.husky/pre-commit` containing `npx lint-staged` (no `husky install`)
4. **Commit identity per repo** — set the author identity the owner chose before the first commit
5. **Never embed a token in a remote URL** — push through the credential helper or deploy key the owner provisioned
6. **CI default branch mismatch** — many repos have `feat/agents` or similar as default, NOT `main`
7. **IaC repos don't need package.json** — use Makefile + pre-commit instead of npm scripts
8. **yamllint exclude patterns** — NATS stream YAML and Helm values often fail strict lint; exclude them
9. **terraform validate needs init** — use `terraform init -backend=false` in CI (no state access needed)
10. **Secret scanning** — always add to CI for IaC repos (private keys, hardcoded passwords in YAML)
