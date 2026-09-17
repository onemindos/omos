---
id: location-attic
title: "Location: Attic Storage"
type: reference
status: active
created: 2026-03-29
updated: 2026-03-29
owner: zeus
tags: [type/reference, system/library, location/attic]
ai_summary: "Dataview index of all books stored in attic — long-term storage, less accessible"
---

# Attic Storage

Books in long-term attic storage. Less accessible — used for overflow, seasonal, or archive.

---

## Inventory

```dataview
TABLE title AS "Title", author AS "Author", genre AS "Genre", condition AS "Condition", sku AS "SKU"
FROM "_library"
WHERE type = "book" AND location = "attic"
SORT title ASC
```

## Stats

```dataview
TABLE length(rows) AS "Count"
FROM "_library"
WHERE type = "book" AND location = "attic"
GROUP BY category
```


---

## Related

- [[25-49 HP (Holistic Performance)/27 Mind (Education & Learning)/27 Mind|27 Mind]]
