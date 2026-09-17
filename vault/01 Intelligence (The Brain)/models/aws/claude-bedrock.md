# claude-bedrock

Claude models accessed via AWS Bedrock — same models as direct Anthropic API, billed through AWS with IAM auth.

| Model ID | Context | Notes |
|----------|---------|-------|
| `anthropic.claude-3-7-sonnet-20250219-v1:0` | 200K | Same as direct Anthropic |
| `anthropic.claude-3-5-haiku-20241022-v1:0` | 200K | Fast/cheap |
| `anthropic.claude-3-5-sonnet-20241022-v2:0` | 200K | Stable release |

**When to use Bedrock Claude vs Direct Anthropic:**
- Bedrock: when S3/Lambda is the trigger, when using AWS VPC, when consolidated AWS billing is needed
- Direct: for all other agent use — simpler, same price per token


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
