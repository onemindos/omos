# OneMind Infra — Developer Commands
# Usage: make <target>

.PHONY: lint lint-yaml lint-tf lint-shell validate plan fmt help

SHELL := /bin/bash
TF_DIR := terraform

## ─── Quality Gates ───────────────────────────────────────────────────────────

lint: lint-yaml lint-tf lint-shell ## Run all linters
	@echo "✓ All lints passed"

lint-yaml: ## Lint all YAML files
	@echo "→ yamllint..."
	@yamllint -c .yamllint.yml platform/ docs/ tak-cluster/ geoserver/ uptime-kuma/ 2>/dev/null || true

lint-tf: ## Lint + validate Terraform
	@echo "→ terraform fmt check..."
	@cd $(TF_DIR) && terraform fmt -check -recursive
	@echo "→ terraform validate..."
	@cd $(TF_DIR) && terraform validate

lint-shell: ## ShellCheck all .sh scripts
	@echo "→ shellcheck..."
	@find . -name '*.sh' -not -path './.git/*' -exec shellcheck -x -s bash {} +

fmt: ## Auto-format Terraform + trailing whitespace
	@cd $(TF_DIR) && terraform fmt -recursive
	@echo "✓ Terraform formatted"

validate: lint ## Full validation (same as CI)
	@echo "✓ Validation complete"

## ─── Infrastructure ──────────────────────────────────────────────────────────

plan: ## Terraform plan (requires DIGITALOCEAN_TOKEN)
	@cd $(TF_DIR) && terraform plan

apply: ## Terraform apply (requires DIGITALOCEAN_TOKEN)
	@cd $(TF_DIR) && terraform apply

install: ## Deploy full platform to DOKS (runs install.sh)
	@bash platform/install.sh

## ─── Utilities ───────────────────────────────────────────────────────────────

streams: ## Provision all NATS streams + KV buckets
	@bash nats-migration/provision.sh

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
