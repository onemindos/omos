---
id: location-local
title: "Location: Local (In-Room)"
type: reference
status: active
created: 2026-03-29
updated: 2026-03-29
owner: zeus
tags: [type/reference, system/library, location/local]
ai_summary: "Dataview index of all books stored in-room / bedroom bookshelf"
---

# Local — In-Room Bookshelf

Books currently in the main living space / bedroom.

---

## Inventory

```dataview
TABLE title AS "Title", author AS "Author", genre AS "Genre", condition AS "Condition", shelf AS "Shelf", sku AS "SKU"
FROM "_library"
WHERE type = "book" AND location = "local"
SORT shelf ASC, title ASC
```

## Stats

```dataview
TABLE length(rows) AS "Count"
FROM "_library"
WHERE type = "book" AND location = "local"
GROUP BY condition
```


---

## Related

- [[25-49 HP (Holistic Performance)/27 Mind (Education & Learning)/27 Mind|27 Mind]]
