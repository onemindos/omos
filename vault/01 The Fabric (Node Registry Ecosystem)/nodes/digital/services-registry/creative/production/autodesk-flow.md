---
title: "Autodesk Flow Production Tracking"
type: reference
domain: SO
status: active
created: 2026-05-27
updated: 2026-05-27
tags: [platform, creative, media-production, pipeline, saas, autodesk, shotgrid]
---

# Autodesk Flow Production Tracking

## What It Is
Flow Production Tracking (formerly ShotGrid) is Autodesk's production tracking and pipeline management platform for film, television, games, and VFX. It is a centralized system for managing workflows, tracking assets, shots, and tasks across the full production pipeline.

Official developer docs: https://developers.shotgridsoftware.com
REST API: https://developers.shotgridsoftware.com/rest-api/

## Core Capabilities

| Feature | Description |
|---------|-------------|
| Entity Management | Full CRUD on Projects, Assets, Shots, Tasks, Sequences, Cuts, Users |
| Production Tracking | Task assignment, status tracking (Hold, In Progress, Lost, etc.), milestones |
| Project Organization | Hierarchical structure — sequences → shots → assets → tasks |
| Search & Filtering | Complex queries with logical operators (and/or), pagination, field selection |
| Collaboration | User assignments, project collaborators, shared labels, created-by/updated-by tracking |
| Visual Assets | Thumbnail management, filmstrip images, billboard attachments |
| Integration | REST API (v1 and v1.1), OpenAPI v3 spec, webhooks, event notifications |
| Licensing | Seat-based — assigned, free, and total seat types |

## Entity Types
Projects, Assets, Shots, Tasks, Sequences, Cuts, CutItems, HumanUser, Phases, Tags, TaskTemplates, custom entities.

## API
- RESTful, JSON responses
- Base endpoint: `https://{site}.shotgridsoftware.com/api/v1.1/`
- Authentication via access tokens
- Batch operations, sorting, filter presets, pagination with next/prev links
- OpenAPI v3 spec available in YAML and JSON
- MCP server implementation available (`shotgrid-mcp-server`)

## Licensing
Seat-based. Pricing not publicly listed in documentation — available on Autodesk's product pages at autodesk.com.

## OneMind OS Angle
Relevant if the Real Time Story Creation project scales to a multi-collaborator pipeline needing structured version control for renders, AI-generated sequences, and asset reviews. The REST API and webhook support means it could integrate with NATS event streams for automated status updates.

## Related Concepts
- [[touch-designer]] — real-time generative layer that produces assets tracked in Flow
- [[real-time-story-creation]] — production context where Flow becomes relevant

## Notes
<!-- Add: Autodesk account status, trial info, pricing tier for project scale -->
