---
id: library-dashboard
title: "Book Library Dashboard"
type: reference
status: active
created: 2026-03-29
updated: 2026-03-29
owner: zeus
tags: [type/reference, system/library, status/active]
ai_summary: "Master dashboard for the book library — collection stats, recent additions, reading progress"
---

# Book Library Dashboard

---

## Collection Overview

```dataview
TABLE length(rows) AS "Count"
FROM "_library"
WHERE type = "book"
GROUP BY format
```

## By Category

```dataview
TABLE length(rows) AS "Count"
FROM "_library"
WHERE type = "book"
GROUP BY category
```

## By Location

```dataview
TABLE length(rows) AS "Count"
FROM "_library"
WHERE type = "book"
GROUP BY location
```

## Reading Progress

```dataview
TABLE length(rows) AS "Count"
FROM "_library"
WHERE type = "book"
GROUP BY read_status
```

## Recently Added (Last 20)

```dataview
TABLE title AS "Title", author AS "Author", sku AS "SKU", format AS "Format", location AS "Location"
FROM "_library"
WHERE type = "book"
SORT created DESC
LIMIT 20
```

## Top Rated

```dataview
TABLE title AS "Title", author AS "Author", rating AS "Rating", genre AS "Genre"
FROM "_library"
WHERE type = "book" AND rating > 0
SORT rating DESC
LIMIT 10
```

## Currently Reading

```dataview
TABLE title AS "Title", author AS "Author", date_started AS "Started", location AS "Location"
FROM "_library"
WHERE type = "book" AND read_status = "reading"
SORT date_started ASC
```

## Books Lent Out

```dataview
TABLE title AS "Title", author AS "Author", location AS "Lent To"
FROM "_library"
WHERE type = "book" AND status = "lent"
SORT title ASC
```


---

## Related

- [[25-49 HP (Holistic Performance)/27 Mind (Education & Learning)/27 Mind|27 Mind]]
