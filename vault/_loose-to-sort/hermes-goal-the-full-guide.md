# Hermes /goal — THE FULL GUIDE

> **Title:** Hermes /goal — THE FULL GUIDE
> **Author:** YanXbt ([@IBuzovskyi](https://x.com/IBuzovskyi))
> **Source:** [x.com/IBuzovskyi/status/2056764150936748082](https://x.com/IBuzovskyi/status/2056764150936748082?s=12)
> **Type:** X Article | **Posted:** 2026-05-19 | **Ingested:** 2026-06-07

---

1. What is /goal and why it matters
/goal is one of the most powerful features introduced in Hermes v0.14 (Foundation Release).
Unlike normal chat interactions where you give a task and get an immediate response, /goal turns Hermes into an autonomous agent. You set a long-term objective, and Hermes will break it down into smaller tasks, use tools, write and run code, iterate, and continue working until the goal is completed (or until you stop it).
In simple terms, it transforms Hermes from a reactive chatbot into a background worker that can handle complex, multi-step tasks with minimal supervision.
2. Main Commands
Here are the key commands you’ll use:
/goal <description>
• Command: /goal <description>
• Function: Starts working on a long-term goal
• When to use: Main command to begin
/goal or /goal status
• Command: /goal or /goal status
• Function: Shows current progress
• When to use: Check how the task is going
/goal pause
• Command: /goal pause
• Function: Pauses the current goal
• When to use: Temporarily stop execution
/goal resume
• Command: /goal resume
• Function: Resumes a paused goal
• When to use: Continue after pausing
/goal clear
• Command: /goal clear
• Function: Clears the current goal
• When to use: Start fresh
/subgoal <text>
• Command: /subgoal <text>
• Function: Adds extra conditions or sub-objectives
• When to use: Refine requirements during execution
/handoff <platform>
• Command: /handoff <platform>
• Function: Transfers the session to Telegram, Discord, etc.
• When to use: Continue work in another app
3. How to write strong goals
This is the most important part. The quality of your results depends heavily on how well you define the goal.
Good goals are:
- Specific and measurable
- Have clear success criteria
- Well-scoped (not too broad)
Examples of strong goals:
- “Create a fully functional Flappy Bird clone in HTML5 with physics, keyboard and mouse controls, scoring system, and collision detection. The game must run on localhost and all core mechanics must work.”
- “Build a clean multi-page website with a homepage, features, and pricing. Make it responsive and pass basic Lighthouse checks.”
- “Refactor the main processing module, improve performance by 30%, add proper error handling, and ensure all tests pass.”
Weak goals (avoid these):
- “Make something cool”
- “Improve my code”
- “Work on the project”
Rule: The more clearly you can describe the final outcome and how to verify it, the better Hermes will perform.
4. Recommended Workflow
1. Provide context first
Give Hermes information about your project, tech stack, folder structure, and previous decisions.
2. (Optional) Generate good goals
Use this meta-prompt:
> “Based on what you know about me and my projects, suggest 3 strong /goal ideas that would run for a long time and create the most value.”
3. Launch the goal
Write /goal followed by a detailed description.
4. Manage the process
- Use /goal status to check progress
- Add /subgoal if you need to adjust direction
- Pause or resume when needed
5. Review the result
Hermes will return completed work, a summary, or an explanation if the goal couldn’t be fully achieved.
5. Best Practices & Common Mistakes
Best practices:
- Always make goals measurable with clear success criteria
- Use /subgoal actively to steer the agent
- Increase max_turns for long-running tasks (hermes config set goals.max_turns 500)
- Combine Hermes with Codex or Claude Code for better results
- Run complex goals overnight
Common mistakes:
- Using vague goals without success criteria
- Intervening too often instead of letting the agent work
- Starting without enough project context
- Making goals too broad or open-ended
6. When to use /goal
Good for: (1/2)- Complex, multi-step tasks (building apps, refactoring, research)
- Work that benefits from iteration and self-correction
- Tasks you can leave running for hours
- Situations where you want to delegate rather than micromanage
Not ideal for:
- Simple or quick tasks
- Situations where you need full control over every step
- When you haven’t clearly defined the desired outcome yet
7. Examples of strong /goal prompts
Here are practical examples of well-written goals:
Example 1 – Game
/goal Create a fully functional Flappy Bird clone in HTML5. Include physics, keyboard and mouse controls, scoring system, and collision detection. The game must run on localhost and all core mechanics must work without bugs.
Example 2 – Web Project
/goal Build a clean multi-page website for a productivity tool. Include homepage, features page, and pricing section. Use modern design, responsive layout, and smooth animations. All pages must pass basic Lighthouse checks.
Example 3 – Code Refactoring
/goal Refactor the main processing module in my repository. Improve performance by at least 30%, add proper error handling, write unit tests for all functions, and ensure all existing tests still pass.
Example 4 – Research
/goal Research 5 competitors in the AI productivity space. Create a structured comparison table with pricing, key features, strengths and weaknesses. Save the final report as a markdown file.
Tip: Start with a clear end result + verifiable criteria. You can always add more details later using /subgoal.
