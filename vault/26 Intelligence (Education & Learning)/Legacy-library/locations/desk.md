---
id: location-desk
title: "Location: Desk / Workspace"
type: reference
status: active
created: 2026-03-29
updated: 2026-03-29
owner: zeus
tags: [type/reference, system/library, location/desk]
ai_summary: "Dataview index of all books on desk — active reading and current reference"
---

# Desk / Workspace

Books currently on or near the desk. Typically active reading or frequently referenced material.

---

## Inventory

```dataview
TABLE title AS "Title", author AS "Author", read_status AS "Reading", genre AS "Genre", sku AS "SKU"
FROM "_library"
WHERE type = "book" AND location = "desk"
SORT title ASC
```

## Stats

```dataview
TABLE length(rows) AS "Count"
FROM "_library"
WHERE type = "book" AND location = "desk"
GROUP BY read_status
```
