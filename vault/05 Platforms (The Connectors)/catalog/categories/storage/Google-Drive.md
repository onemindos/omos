# Google Drive

> **Google's cloud storage and document collaboration platform**

## Overview

Google Drive is Google's cloud storage integrated with Google Workspace (Docs, Sheets, Slides). Used for document collaboration, sharing files with external parties, and backup of specific GE/HP materials that need non-Apple cloud access.

## Status

| Attribute   | Value                            |
| ----------- | -------------------------------- |
| Status      | ✅ Active                         |
| Category    | Storage                          |
| Type        | Cloud Storage / Collaboration    |
| Integration | Web, iOS, macOS, Drive API       |

## Use Cases in OneMind

1. **Document Collaboration** — Shared docs with clients, partners, collaborators
2. **GE Content** — Marketing assets, brand decks, business documents
3. **Backup Layer** — Secondary backup for critical HP/GE files
4. **Integration** — N8N Drive node for file-based automations

## Key Features

- 15 GB free (Google One upgrades available)
- Real-time collaboration (Docs, Sheets, Slides)
- Google Drive API v3 for programmatic access
- N8N native node
- File sharing with permission levels
- Folder sync via Google Drive for Desktop

## Quick Start (API)

```python
from google.oauth2 import service_account
from googleapiclient.discovery import build

creds = service_account.Credentials.from_service_account_file('creds.json')
service = build('drive', 'v3', credentials=creds)
results = service.files().list().execute()
```

## Links

- Website: [drive.google.com](https://drive.google.com)
- API Docs: [developers.google.com/drive](https://developers.google.com/drive)


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
