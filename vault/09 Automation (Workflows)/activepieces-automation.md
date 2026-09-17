# Activepieces: Open Source Automation

**Type:** System Architecture / Automation  
**Category:** No-Code/Low-Code, Workflow Engine  
**Created:** 2026-01-20  
**Status:** Reference  
**Website:** [activepieces.com](https://www.activepieces.com/)  
**License:** MIT (Open Source)

---

## Overview

**Activepieces** is an open-source, no-code business automation tool. It acts as a self-hosted alternative to Zapier or Make.com. It allows you to build complex workflows by connecting various apps (Google Sheets, OpenAI, Slack, etc.) using a visual drag-and-drop interface.

**Core Philosophy:**

> "AI-first automation for every team."  
> It emphasizes not just connecting APIs but integrating AI agents directly into the flow.

---

## Why Activepieces for OneMind?

For the **Sovereign Operations (UI)** pillar, we prioritize systems that are:

1.  **Self-Hostable:** We own the infrastructure and data.
2.  **Extensible:** We can write custom code (TypeScript/Node.js) if a connector is missing.
3.  **Cost-Effective:** No "task run" limits like Zapier when self-hosted.
4.  **AI-Native:** Built to work seamlessly with LLMs (ChatGPT, Claude) for intelligent processing.

---

## Architecture & Concepts

### 1. Flows

A **Flow** is a sequence of steps triggered by an event.

- **Trigger:** What starts the flow (e.g., "New Email", "Webhook Received", "Schedule: Every Morning").
- **Actions:** What happens next (e.g., "Send to Slack", "Update Database", "Ask ChatGPT").

### 2. Pieces (Connectors)

**Pieces** are the building blocks. There are 100+ pre-built pieces (OpenAI, Google Drive, Telegram, etc.).

- _Key Feature:_ You can install pieces dynamically via npm or build your own custom pieces.

### 3. The Engine

Activepieces runs on a standard stack:

- **Frontend:** Angular (Visual Builder)
- **Backend:** Node.js / NestJS
- **Database:** PostgreSQL (stores flow definitions and run history)
- **Queue:** Redis (manages job execution)

---

## Self-Hosting Guide (Docker Compose)

For our infrastructure (`08 Infrastructure`), we use Docker Compose to deploy Activepieces alongside our other services.

**Recommended `docker-compose.yml` Config:**

```yaml
version: "3"

services:
  activepieces:
    image: activepieces/activepieces:latest
    container_name: activepieces
    restart: always
    ports:
      - "8080:80"
    environment:
      # Database Configuration
      - AP_POSTGRES_HOST=postgres
      - AP_POSTGRES_PORT=5432
      - AP_POSTGRES_DATABASE=activepieces
      - AP_POSTGRES_USERNAME=postgres
      - AP_POSTGRES_PASSWORD=your_password
      # Redis Configuration
      - AP_REDIS_HOST=redis
      - AP_REDIS_PORT=6379
      # Security (Generate random strings for these)
      - AP_APP_KEY=random_secret_key_here
      - AP_JWT_SECRET=another_random_secret
      # Execution Mode
      - AP_EXECUTION_MODE=UNSANDBOXED # Use SANDBOXED for multi-tenant security
    depends_on:
      - postgres
      - redis

  postgres:
    image: postgres:14-alpine
    container_name: activepieces_db
    restart: always
    environment:
      - POSTGRES_DB=activepieces
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=your_password
    volumes:
      - activepieces_db_data:/var/lib/postgresql/data

  redis:
    image: redis:7-alpine
    container_name: activepieces_redis
    restart: always
    volumes:
      - activepieces_redis_data:/data

volumes:
  activepieces_db_data:
  activepieces_redis_data:
```

---

## Integration Strategy

### 1. The "Glue" of OneMind

Activepieces serves as the central nervous system that connects our disparate self-hosted tools.

- **Input:** Webhooks from NocoDB, Obsidian (via community plugins), or external APIs.
- **Processing:** Activepieces cleans data, runs Python scripts, or queries an LLM.
- **Output:** Updates dashboards, sends Telegram notifications, or writes back to Vault.

### 2. Use Case Examples

**A. Holistic Performance (Health Sync)**

- **Trigger:** Webhook from Oura/Apple Health output.
- **Action:** Activepieces formats the JSON data.
- **Action:** Write row to NocoDB "Health Metrics".
- **Action:** If `Sleep Score < 70`, send Telegram msg: "Low energy expected today. Reschedule deep work?"

**B. Legacy Evolution (Family Archives)**

- **Trigger:** New photo uploaded to "Family Drop" folder (Google Drive/Nextcloud).
- **Action:** AI Vision (OpenAI) analyzes photo for faces/context.
- **Action:** Renames file with date and description (e.g., `2026-01-20_birthday_party.jpg`).
- **Action:** Moves to "Processed" folder.

**C. Generational Endeavor (Lead Ops)**

- **Trigger:** New lead in Typeform.
- **Action:** Use "Company Search" piece to enrich data.
- **Action:** Draft personalized email using Claude.
- **Action:** Create task in Todoist `81 Sales Strategy`.

---

## Advanced: "Code Pieces"

One of Activepieces' strongest features for developers is the **Code Piece**.
You can drop in a block of TypeScript/Node.js to do _anything_ not covered by a standard connector.

```javascript
export const code = async (inputs) => {
  const data = inputs.data;
  // Custom logic here
  const score = data.value * 1.5;
  return { processed_score: score };
};
```

---

## Comparison to Alternatives

| Feature            | Activepieces        | Zapier         | n8n                    |
| :----------------- | :------------------ | :------------- | :--------------------- |
| **Self-Hostable**  | ✅ Yes              | ❌ No          | ✅ Yes                 |
| **License**        | MIT (Open Core)     | Closed         | Fair Code (Restricted) |
| **UI Experience**  | Simple, Linear      | Simple, Linear | Complex, Node graph    |
| **AI Integration** | First-class citizen | Add-on         | Add-on                 |
| **Pricing**        | Free (Self-hosted)  | $$$ per task   | Free\* (Self-hosted)   |

_Note: Activepieces is often preferred over n8n for simpler workflows due to its cleaner UI and truly open MIT license._

---

**Tags:** `#automation` `#opensource` `#selfhosted` `#workflow` `#activepieces`

**Related Docs:**

- `system-docker-compose-setup.md` (Coming soon)
- `process-automation-audit.md` (Coming soon)
