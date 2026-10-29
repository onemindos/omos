# Local model candidates for the offline fallback (D10)

agentops calls a local model through any OpenAI-compatible server
(`MODEL_LOCAL=openai:http://127.0.0.1:11434/v1|<model>`). Nothing here is installed or
downloaded by the PR; the owner picks one. Licence facts are from the Hugging Face model API
(`https://huggingface.co/api/models/<repo>`, field `cardData.license`) and the GitHub API for
the runtimes, fetched 2026-10-29 (onemind-xr `docs/evidence/M4/I1/local-model-licences.txt`).
Check the licence again at download time; model cards change.

## Runtimes

| Runtime | Licence | Notes |
|---|---|---|
| Ollama (`ollama/ollama`) | MIT | OpenAI-compatible API at `/v1`; one binary on macOS, Linux, Windows |
| llama.cpp `llama-server` (`ggml-org/llama.cpp`) | MIT | OpenAI-compatible API; GGUF model files |

## Models

| Model (Hugging Face repo) | Licence | Gated | Size class | Fit |
|---|---|---|---|---|
| `Qwen/Qwen2.5-7B-Instruct` | Apache-2.0 | no | 7B | **Suggested first try on the Mac or home server** (open licence, 7B fits the Mac's memory class; quality for the JSON explain step unverified) |
| `Qwen/Qwen3-8B` | Apache-2.0 | no | 8B | Newer alternative to Qwen2.5-7B |
| `Qwen/Qwen3-4B` | Apache-2.0 | no | 4B | **Suggested first try for a field kit** (laptop, Jetson); quality unverified |
| `mistralai/Mistral-7B-Instruct-v0.3` | Apache-2.0 | no | 7B | Alternative 7B |
| `ibm-granite/granite-3.3-8b-instruct` | Apache-2.0 | no | 8B | Alternative 8B |
| `ibm-granite/granite-3.3-2b-instruct` | Apache-2.0 | no | 2B | Small devices |
| `HuggingFaceTB/SmolLM2-1.7B-Instruct` | Apache-2.0 | no | 1.7B | Smallest; quality unverified |
| `microsoft/Phi-4-mini-instruct` | MIT | no | 3.8B | Field kit alternative |
| `microsoft/Phi-3.5-mini-instruct` | MIT | no | 3.8B | Older Phi |
| `Qwen/Qwen2.5-3B-Instruct` | **other** (a custom licence per its model card, not Apache-2.0) | no | 3B | Avoid: not an open licence, unlike the 7B |
| `meta-llama/Llama-3.2-3B-Instruct` | **Llama 3.2 Community Licence** | yes (manual) | 3B | Avoid by default: custom licence with an acceptable-use policy, gated download |
| `google/gemma-3-4b-it` | **Gemma terms of use** | yes (manual) | 4B | Avoid by default: custom terms, gated download |

Size classes are parameter counts from the repo names; memory and speed on the owner's
hardware were **not measured** in this run (no model was downloaded here: disk headroom and
the no-download rule for this lane).

## Whatever model is chosen

The model only words facts the code gathered. Its output is checked: unknown ids are removed,
missing gaps and citations are added back, and a watch proposal outside the watch's list is not
sent. The built-in stub is always the last link, so a brief is written even with no model.
