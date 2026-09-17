# How to make Hermes + xurl Actually Work as a System

> **Title:** How to make Hermes + xurl Actually Work as a System
> **Author:** YanXbt ([@IBuzovskyi](https://x.com/IBuzovskyi))
> **Source:** [x.com/IBuzovskyi/status/2057114309616885997](https://x.com/IBuzovskyi/status/2057114309616885997?s=12)
> **Type:** X Article | **Posted:** 2026-05-20 | **Ingested:** 2026-06-07

---

xurl is a Hermes skill that gives your agent the ability to directly interact with X — searching posts, reading discussions, and publishing content.
On its own, xurl is mainly an execution tool. It becomes significantly more powerful when combined with other Hermes skills, especially /goal, research, and memory. This combination allows you to move from basic prompting to building structured, repeatable workflows.
In this guide, we’ll break down the most effective ways to combine xurl with other skills and show how to build a reliable system around it.
Why Basic xurl Usage Falls Short
When xurl is used without additional structure, most people face the same problems:
- Results are inconsistent because there’s no defined process
- The agent has no memory of what was already done
- Quality control is weak or missing
- The agent often repeats topics or produces shallow content
This keeps xurl at the level of a helpful feature rather than turning it into something you can build reliable workflows around.
Best Skill Combinations
Here are the combinations that give the most practical value:
xurl + /goal
The strongest combination. /goal lets you define a clear process (research → analysis → drafting → evaluation → publishing). xurl is responsible only for execution, while /goal handles structure and quality.
xurl + Research Skills
Useful when you need the agent to gather and process information before creating content.
xurl + Multi-step Reasoning
Helps the agent go through several rounds of thinking and improvement before publishing.
xurl + Memory / Planning
Important for recurring processes. Memory helps avoid repetition and maintain consistency over multiple runs.
How to Build a Recurring Content Workflow
One of the most useful setups is a recurring workflow that runs twice a day. Here’s how you can structure it:
bash
hermes goal "Every morning and evening, research the latest AI agent discussions, analyze them according to my content style, check against previously published topics, generate a high-quality thread, evaluate its potential virality, and publish it using xurl if it meets the quality bar."
To make this workflow more reliable, break it down into clear stages:
1. Data Collection — Use xurl to gather recent discussions.
2. Style Check — Make sure the content matches your defined style.
3. Repetition Check — Compare with previously published topics.
4. Draft Creation — Generate the thread.
5. Quality Evaluation — Score the draft on hook strength, insight, and engagement potential.
6. Publishing Decision — Publish only if the score is high enough.
Common Mistakes and How to Avoid Them
When building these workflows, people usually make the same errors:
- Trying to do everything in one goal — Large goals often lose context. Break them into smaller, clearly defined stages.
- Weak evaluation criteria — If the agent is too lenient when scoring its own drafts, quality drops. Make evaluation rules strict and specific.
- No memory tracking — Without tracking published topics, the agent starts repeating itself. Add explicit instructions to check previous outputs.
- Too much automation too soon — Don’t let the agent publish everything automatically at the beginning. Review outputs manually first.
How to Start (Step-by-Step)
If you want to implement this approach, here’s a practical and gradual way to begin:
Step 1: Create Your First Basic Goal
Start with a simple goal that combines xurl and /goal. Don’t overcomplicate it yet.
hermes goal "Research the latest AI agent discussions from the last 12 hours and publish a short thread using xurl"
Run this goal manually several times until the agent can reliably use xurl to post.
Step 2: Add Structure to the Process
Once the basic version works, expand the goal by breaking it into clear stages: research → analysis → drafting → evaluation → publishing.
Step 3: Introduce Memory and Repetition Control
Add instructions for the agent to track previously published topics. For example:
“Before generating a new thread, check the topics you have already posted about in the last 7 days and avoid repeating the same angles.”
Step 4: Add Quality Evaluation
Introduce clear evaluation criteria before publishing. Instruct the agent to score the draft on hook strength, insight depth, relevance, and style alignment. Only publish if the average score is high enough.
Step 5: Make the Workflow Recurring
Only after the goal works reliably in manual mode, set it up to run automatically (morning and evening). It’s recommended to run it manually for at least 5–7 days first.
Step 6: Review, Adjust, and Iterate
For the first 7–10 runs, review the output manually. Pay attention to content quality, style consistency, repetition, and evaluation accuracy. Adjust the goal instructions based on your observations.
What Improves When You Combine Skills
When xurl is used together with structure and other skills, several things improve:
- Output becomes more consistent
- You reduce the need for constant manual prompting
- Recurring processes become realistic to run
- The agent can handle multi-stage work with less supervision
Final Thoughts
xurl becomes much more powerful when combined with other Hermes skills. The most practical results right now come from pairing it with /goal and supporting capabilities like research and memory.
If you want to move beyond simple commands, building structured workflows that include xurl is one of the most effective next steps.
