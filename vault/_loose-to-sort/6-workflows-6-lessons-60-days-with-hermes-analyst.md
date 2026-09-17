# 6 Workflows, 6 Lessons, 60 Days with Hermes Analyst

> **Title:** 6 Workflows, 6 Lessons, 60 Days with Hermes Analyst
> **Author:** 0xJeff ([@0xJeff](https://x.com/0xJeff))
> **Source:** [x.com/0xJeff/status/2061361760968560887](https://x.com/0xJeff/status/2061361760968560887?s=12)
> **Type:** X Article | **Posted:** 2026-06-01 | **Ingested:** 2026-06-07

---

It’s been a little more than 2 months now and Hermes has taught me so many different things.
One of the key things is that Agents like Hermes or OpenClaw don’t fail on intelligence, they tend to fail on architecture. Most bug was tools fighting each other, not the model being dumb (although some models are actually dumb lol).
In this article, we’ll cover key lessons I’ve learnt across major investment/data analyst workflows so you don’t have to repeat them.
Let’s dig in 👇
1. Model Lesson - Shifting Provider to Provider can take a lot of time, efforts, and money
In 60 days, we ran through 5-6 provider setups for Hermes + Hindsight:
OpenRouter → Kimi-k2.6, Qwen 3.6 27B, DeepSeek v4 Flash & Pro, Gemma, and more
Opencode Go subscription + Opencode Zen (testing Mini Max 2.5-2.7, GLM5.1, and Kimi-k2.6)
DeepSeek direct API
Venice AI for private inference stack (using DIEM credits)
Grok subscription for x_search
Each swap cost 2-3 sessions of debugging API compatibility, auth flows, and timeout settings.
What I learned: Model started to matter less and less. Open-weight labs are measuring up to frontier labs level in term of intelligence while keeping the cost low.
The best way to go is to pick 1 provider and stick with them. Going direct tend to net you with better discount/connectivity (e.g. DeepSeek provided 75% discount on direct API first before making it permanent and available on Openrouter).
I opted for Opencode Go $5/month package when I started out and the experience wasn’t that great compared to using direct DeepSeek API. The inference goes through multi-hops which add 5-10 sec wait time on top which makes it very slow.
2. What matters more is the tools and skills
Nous Research recently got called out by tech youtuber Theo for having many useless skills like Polymarket skill. Some say having too many skills bloat up the context. I say having too many skills = a lot better than not having enough skills for Hermes.
Hermes auto-creates a skill every time he realizes that you’ll be running the same workflow again in the future. This saves a lot of time and inference cost - 3 mins of manual workflow takes 10 s the second time around.
Having many natively-integrated skill would allow Hermes to easily test & auto set up relevant skills for users (Instead of finding external tools, GitHub repos (that may or may not be safe) yourself).
Most of your job as a user is to point the agent to the right tool(s) for the right job. Course correct Hermes to prevent him from using the wrong tool.
Example:
Using Exa for structured web search and Firecrawl for JS-heavy sites is a lot better than manual web search
Browser CDP is better than headless browser solution like Playwright that Cloudflare tend to block
Direct tool access via API, MCP, markdown skill file tend to yield far better outputs + save time & cost compared to manual web search or Browser CDP
For any recurring workflows/cron jobs, having the right tool is best.
For one-off scraping or interactive exploration, using Browser CDP and/or agentic search is great.
3. Hermes key differentiation stems from understanding user’s preferences
While native memory/preference via User(.)md, Memory(.)md, Soul(.)md help Hermes understands you + the world better, having an external memory provider enables Hermes to recall past learnings and draw relationships between facts & events better.
I’ve used Hindsight since the beginning because of its high accuracy in recalling BUT I often make the mistake of spending too much inference on Hindsight (unknowingly) when using reflect.
Made a mistake of hooking Hindsight up with OpenRouter and tied it with Kimi-k2.6 and plugged “Hindsight reflect” to multiple cron jobs which ended up burning $20-30 for multiple days straight.
Reflect often take 240s and would timeout on me so I learned to use “Recall” for time-sensitive cron jobs like reporting/synthesizing morning briefs
The combination of having rich native memory + external memory provider result in highly relevant outputs that keep me up-to-date on the market.
4. The feedback loop = the easiest way to train your personal agent
#1 Hermes workflow that I look forward to everyday at 10AM is the morning briefs which consist of all X accounts & news sources that we track + portco developments during the day (if any) + Top 10 synthesis of all the insights
We continue to iterate on the outputs to get rid of unnecessary stuff, keeping only the relevant insights.
Feedback loop follows 6 steps
Hermes produces something
I read it immediately and flag what’s wrong
I give a specific correction/next steps\
Hermes encode the correction as a permanent rule
Next output becomes tighter
Repeat the loop
This helps me clean up the outputs to the format that’s easy for me to digest.
The biggest problem right now is echo chamber + self-reinforcing loop where “Why it matters” tend to gravitate towards existing holdings. Sources/analysts tend to mention the same big-cap names (NVIDIA, TSMC, MU, VRT, SIVE) which exacerbate the chamber.
Haven’t solved this yet. If you know any reliable solution, would appreciate it ahaha.
5. x402 solves a lot of problems (for real)
The last 2 weeks was an eye opener for me — x402-integrated tools + x402 tools aggregator is one of the biggest innovations of crypto
Before x402, I found myself spending a lot of time trying to find the right tool for the right job. Trying to find if an MCP exist. If not, sign up for a free API. If not free, ponder whether it’s worth it to pay.
After x402, I can just set up an agentic wallet with 1 install command. Seed it with $5-10 in USDC and start exploring hundreds of premium tools, each of which cost a couple of cents.
Tools like Nansen, Exa/Firecrawl, Surf, BlockRun, and more that can help research and access premium data in a relatively short amount of time.
I tested x402 out once and since then integrated it as part of my usual due diligence workflows now.
Covered part of this last week for Onchain Investigation workflow (if you haven’t seen it yet)
0xJeff
@0xJeff
·
May 25
 Article
Hermes as an Onchain Analyst
Every single week, Hermes gets an upgrade.
Earlier last week it was x_search where Hermes could directly tap into X rich real-time news/insights without having to use X API (saving lots of $ on...
19
70
633
112K
6. Skill bundling is a must
Early on, every workflow was one giant prompt. The onchain dump investigation skill was 2000+ words covering Dexscreener, Nansen, RPC endpoints, Cookie MCP, output synthesis format. Every time a failure or a bug was found, the prompt got longer.
Each new session would start fresh, asking “Check why [x] dump” would prompt Hermes to re-derive it from the growing spaghetti wall of text.
The fix is pretty simple → treat skill as a directory instead of prompt.
onchain-dump-investigation/
├── SKILL(.)md # Pipeline logic (stays at ~100 lines)
├── references/
│ ├── nansen-agentcash(.)md # Endpoint shapes + field quirks
│ ├── base-rpc-endpoints(.)md # Working RPCs, rate limit patterns
│ └── cookie-mcp-queries(.)md # Query templates + cross-reference matrix
└── scripts/
└── check_wallets(.)sh
When the skill loads, the agent gets the full context — pipeline instructions from SKILL(.)md, API specifics from references, and executable scripts. The prompt stays lean because the noise (endpoint URLs, field names, error codes) lives in reference files that only load when relevant.
A well-bundled skill costs ~500 tokens to load (SKILL(.)md + 2-3 reference files) but saves 5000+ tokens of re-explaining context in every session. Over 60+ days and hundreds of sessions, the cost saving & the efficiency compounds.
Wrapping it up
If there’s one thread tying everything together, it’s building an agent is 90% architecture, 10% AI. Everyone has access to the same models (most of which are highly intelligent).
What separates useful agent to useless agent is tool/skill design, memory persistence, feedback/learning loop, and the unit economics that make it sustainable to run.
On the side note, I’ve just started giving Hermes a small bankroll to run autonomous prediction market trading experiment. Running into so many diff bugs rn, makes me appreciate builders who are building reliable agents for users lol.
^Will let you know how this goes if the results are great. As always, thanks for reading. Let me know your experience with Hermes.
Rec setting it up if you haven’t yet. Super simple hermes install → configure your model provider → start feeding your preference & needs.
Personal Note: We have over 10+ tools now (most of which are free or pretty low cost). If you're interested to know which tool(s) I'm using for my Hermes analyst stack. Do check out the Tool Shed
Thanks for checking this episode out and enjoy working with your Hermes! See you in the next one.
