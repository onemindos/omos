# OneMind OS — Developer Commands
# Usage: make <target>

.PHONY: lint lint-yaml lint-tf lint-shell validate init plan apply fmt install streams help

SHELL := /bin/bash
TF_DIR := infra
YAML_DIRS := platform docs

## ─── Quality Gates ───────────────────────────────────────────────────────────

lint: lint-yaml lint-tf lint-shell ## Run all linters
	@echo "✓ All lints passed"

lint-yaml: ## Lint all YAML files
	@echo "→ yamllint..."
	@yamllint -c .yamllint.yml $(YAML_DIRS)

lint-tf: ## Lint + validate Terraform
	@echo "→ terraform fmt check..."
	@terraform -chdir=$(TF_DIR) fmt -check -recursive
	@terraform -chdir=$(TF_DIR) init -backend=false -input=false
	@echo "→ terraform validate..."
	@terraform -chdir=$(TF_DIR) validate

lint-shell: ## ShellCheck all .sh scripts
	@echo "→ shellcheck..."
	@find . -path './.git' -prune -o -path './vault' -prune -o -type f -name '*.sh' -exec shellcheck -x -s bash {} +

fmt: ## Format Terraform files
	@terraform -chdir=$(TF_DIR) fmt -recursive
	@echo "✓ Terraform formatted"

validate: lint ## Run local lint and validation checks
	@echo "✓ Validation complete"

## ─── Infrastructure ──────────────────────────────────────────────────────────

init: ## Initialize Terraform providers in infra/
	@terraform -chdir=$(TF_DIR) init -input=false

plan: ## Terraform plan (requires DIGITALOCEAN_TOKEN)
	@terraform -chdir=$(TF_DIR) init -input=false
	@terraform -chdir=$(TF_DIR) plan

apply: ## Terraform apply (requires DIGITALOCEAN_TOKEN)
	@terraform -chdir=$(TF_DIR) init -input=false
	@terraform -chdir=$(TF_DIR) apply

install: ## Deploy full platform to DOKS (runs install.sh)
	@bash platform/install.sh

## ─── Utilities ───────────────────────────────────────────────────────────────

streams: ## Show the GitOps source for NATS streams and KV buckets
	@echo "NATS streams and KV buckets are declared under platform/nats/streams/ and reconciled by Argo CD."

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
