# moonshot-v1

Moonshot's stable long-context general-purpose models — pick tier by context need.

| Model ID | Context | Input | Output |
|----------|---------|-------|--------|
| `moonshot-v1-128k` | 128K | $2.60/1M | $10.40/1M |
| `moonshot-v1-32k` | 32K | $0.85/1M | $2.50/1M |
| `moonshot-v1-8k` | 8K | $0.17/1M | $0.50/1M |

**Best for:** Context-tiered pricing — use the smallest window that fits the task. moonshot-v1-8k is a very cheap general model for short tasks.

API: `https://api.moonshot.cn/v1` (OpenAI-compatible), key: `MOONSHOT_API_KEY`.

> Note: Servers in China — higher latency from US/EU. Use async/batch workloads.


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
