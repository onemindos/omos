---
id: ui-02-20260111053925
title: "N8N As The Tool Layer (Mcp Integration)"
type: concept
pillar: UI
domain: 00-24
folder: "00-24 UI (Unified Intelligence)"
status: active
created: 2026-01-11
modified: 2026-01-11
owner: zeus
tags: [domain/ui, folder/02]
todoist_project: "00-24 UI (Unified Intelligence)"
ai_summary: "N8N As The Tool Layer (Mcp Integration)"
---

n8n Workflow Registry
─────────────────────

## Core Workflows (All Agents)
- todoist_sync           # Sync tasks between vault and Todoist
- vault_read             # Read file from Obsidian vault
- vault_write            # Write/update file in vault
- vault_search           # Search vault content
- notify_zeus            # Send notification to you
- log_action             # Log agent action

## HP Workflows
- health_metrics_log     # Log health data
- workout_tracker        # Track workout completion
- nutrition_log          # Log meals

## LE Workflows
- homestead_daily_check  # Generate daily homestead checklist
- family_calendar_sync   # Sync family calendar
- egg_production_log     # Log egg counts

## GE Workflows
- revenue_tracker        # Log revenue/expenses
- community_metrics      # Pull Circle.so analytics
- client_followup        # Generate client followup tasks

## IT Workflows
- infra_health_check     # Check Docker containers, services
- backup_trigger         # Trigger backup workflows
- deploy_update          # Deploy updates to systems

## Capture & Processing
- capture_intake         # Receive capture from any source
- process_capture        # AI categorization + routing
- approval_request       # Send to Zeus for approval
- finalize_to_vault      # Convert to markdown, add to vault