---
title: "DuckDB + Parquet"
type: reference
status: active
tags: [platform, infrastructure, database, analytics, storage]
---
# DuckDB + Parquet

> OneMind's long-term data storage stack. No Postgres, no cloud DB — just files + fast analytics.

## Why This Stack
- **DuckDB** — embedded analytical database (like SQLite but for analytics)
- **Parquet** — columnar file format (compressed, fast reads, portable)
- Together: query terabytes of data from flat files, zero server, zero dependencies

## DuckDB
```bash
# Install
brew install duckdb

# CLI
duckdb mydb.duckdb

# Query Parquet files directly (no import needed)
SELECT * FROM 'data/telemetry/*.parquet' WHERE timestamp > '2026-01-01';

# Query CSV
SELECT * FROM read_csv_auto('data.csv');

# Export to Parquet
COPY (SELECT * FROM table) TO 'output.parquet' (FORMAT PARQUET);
```

## Python
```python
import duckdb

con = duckdb.connect('onemind.duckdb')

# Query Parquet directly
df = con.execute("SELECT * FROM 'telemetry/*.parquet' LIMIT 100").fetchdf()

# Write Parquet
con.execute("COPY table TO 'output.parquet' (FORMAT PARQUET)")

# Aggregate across all files
con.execute("""
    SELECT date_trunc('hour', timestamp) as hour, avg(heart_rate)
    FROM 'health/*.parquet'
    GROUP BY 1 ORDER BY 1
""")
```

## Storage Architecture
```
~/.onemind/data/
├── telemetry/          ← drone, robot, sensor data (Parquet)
├── health/             ← Oura, Strava, Apple Health (Parquet)
├── fleet/              ← agent execution logs (Parquet)
├── financials/         ← transactions, invoices (Parquet)
└── onemind.duckdb      ← active queries, views, indexes
```

## NATS → Parquet Pipeline
```
NATS JetStream (7d hot) → nats-parquet-sink → Parquet files → DuckDB queries
```

## Why Not Postgres
- No server to maintain
- Files are portable (copy, backup, share)
- 10-100x faster for analytical queries
- Works on any node (Pi, laptop, server)
- JetStream handles the "hot" event store (7d)
- Parquet handles the "cold" analytical store (forever)

## Key Features
| Feature | Notes |
|---------|-------|
| Zero server | Embedded — just a library |
| Parquet native | Query files without import |
| Columnar | Fast aggregations on specific columns |
| Compression | Parquet files are 2-10x smaller than CSV |
| Portable | Copy files anywhere, query anywhere |
| Python | First-class pandas/polars integration |
| SQL | Full SQL (joins, window functions, CTEs) |
