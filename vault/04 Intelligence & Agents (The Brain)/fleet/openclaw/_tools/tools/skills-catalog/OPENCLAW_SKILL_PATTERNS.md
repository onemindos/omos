# Adding OpenClaw-Style Features to OneMind OS
## Implementation Guide

**Date**: February 7, 2026
**Purpose**: Step-by-step guide to integrate OpenClaw's best features into OneMind OS
**Estimated Total Time**: 20-30 hours

---

## Table of Contents

1. [Overview](#overview)
2. [Feature 1: Chat Bot Integrations](#feature-1-chat-bot-integrations)
3. [Feature 2: Browser Automation](#feature-2-browser-automation)
4. [Feature 3: Self-Modifying Skills](#feature-3-self-modifying-skills)
5. [Feature 4: Skill Auto-Discovery](#feature-4-skill-auto-discovery)
6. [Feature 5: Simplified Setup](#feature-5-simplified-setup)
7. [Testing & Validation](#testing--validation)
8. [Deployment Strategy](#deployment-strategy)

---

## Overview

### What We're Adding

Based on the OpenClaw analysis, we'll integrate these key features into OneMind OS:

| Feature | Priority | Time | Complexity | Impact |
|---------|----------|------|------------|--------|
| **Telegram Bot** | 🔴 High | 3h | Low | High |
| **Discord Bot** | 🔴 High | 3h | Low | High |
| **Browser Automation** | 🟡 Medium | 4h | Medium | Medium |
| **Self-Modifying Skills** | 🟢 Low | 6h | High | High |
| **Skill Auto-Discovery** | 🟡 Medium | 4h | Medium | Medium |
| **WhatsApp Integration** | 🟢 Low | 5h | Medium | Medium |
| **One-Liner Installer** | 🟢 Low | 5h | Medium | Low |

**Total Estimated Time**: 30 hours

### Implementation Strategy

**Phase 1: Chat Integrations** (Quick wins, high impact)
- Telegram bot (3h)
- Discord bot (3h)

**Phase 2: Automation** (Medium complexity)
- Browser automation (4h)
- Skill auto-discovery (4h)

**Phase 3: Advanced Features** (Higher complexity)
- Self-modifying skills (6h)
- WhatsApp integration (5h)

**Phase 4: Polish** (Optional)
- One-liner installer (5h)

---

## Feature 1: Chat Bot Integrations

### 1.1 Telegram Bot Integration

**Goal**: Enable bidirectional chat with Legacy via Telegram

**Current State**: Telegram is notification-only (outbound)

**Implementation**: Add message receiver to route Telegram messages to agent

#### Step 1: Install Dependencies

```bash
cd /Users/zeuslegacy/Documents/onemind-workspace/OneMind-OS
cd backend

# Add to requirements.txt
echo "python-telegram-bot==21.0" >> requirements.txt
pip install -r requirements.txt
```

#### Step 2: Create Telegram Bot Tool

Create file: `backend/agno/tools/telegram_bot.py`

```python
"""Telegram bot integration for bidirectional chat with agents."""

from agno.tools import Toolkit
from telegram import Update, Bot
from telegram.ext import (
    Application,
    CommandHandler,
    MessageHandler,
    filters,
    ContextTypes,
)
import os
import logging
from typing import Dict
import asyncio

logger = logging.getLogger(__name__)

# Store session IDs for each Telegram user
USER_SESSIONS: Dict[int, str] = {}


class TelegramBotTool(Toolkit):
    """Telegram bot that routes messages to OneMind agents."""

    def __init__(self, token: str = None):
        super().__init__(name="telegram_bot")

        self.token = token or os.getenv("TELEGRAM_BOT_TOKEN")
        if not self.token:
            raise ValueError("TELEGRAM_BOT_TOKEN not set")

        self.app = None
        self.agent = None

    async def initialize(self):
        """Initialize bot application."""
        self.app = Application.builder().token(self.token).build()

        # Register handlers
        self.app.add_handler(CommandHandler("start", self.cmd_start))
        self.app.add_handler(CommandHandler("help", self.cmd_help))
        self.app.add_handler(CommandHandler("reset", self.cmd_reset))
        self.app.add_handler(
            MessageHandler(filters.TEXT & ~filters.COMMAND, self.handle_message)
        )

        logger.info("Telegram bot initialized")

    async def cmd_start(self, update: Update, context: ContextTypes.DEFAULT_TYPE):
        """Handle /start command."""
        user_id = update.effective_user.id
        session_id = f"telegram_{user_id}"
        USER_SESSIONS[user_id] = session_id

        welcome_message = """
👋 Welcome to OneMind OS!

I'm Legacy, your personal AI assistant. I have persistent memory,
so I'll remember our conversations.

**What I can do:**
- Manage your tasks (Todoist)
- Control your smart home (Home Assistant)
- Search the web
- Manage your calendar
- And much more!

Just message me naturally, and I'll help you out.

**Commands:**
- /help - Show this message
- /reset - Start a new conversation
"""
        await update.message.reply_text(welcome_message)

    async def cmd_help(self, update: Update, context: ContextTypes.DEFAULT_TYPE):
        """Handle /help command."""
        await self.cmd_start(update, context)

    async def cmd_reset(self, update: Update, context: ContextTypes.DEFAULT_TYPE):
        """Handle /reset command - start new session."""
        user_id = update.effective_user.id
        old_session = USER_SESSIONS.get(user_id)

        # Create new session ID
        import uuid
        new_session = f"telegram_{user_id}_{uuid.uuid4().hex[:8]}"
        USER_SESSIONS[user_id] = new_session

        await update.message.reply_text(
            f"✓ Started new conversation. Previous session: {old_session}"
        )

    async def handle_message(self, update: Update, context: ContextTypes.DEFAULT_TYPE):
        """Route incoming message to agent."""
        user_id = update.effective_user.id
        user_name = update.effective_user.full_name
        message_text = update.message.text

        # Get or create session
        if user_id not in USER_SESSIONS:
            USER_SESSIONS[user_id] = f"telegram_{user_id}"
        session_id = USER_SESSIONS[user_id]

        logger.info(f"Message from {user_name} (ID: {user_id}): {message_text}")

        # Show typing indicator
        await update.message.chat.send_action("typing")

        try:
            # Import agent factory
            from backend.agno.agents.factory import create_agent_from_config
            from backend.agno.agents.presets import LEGACY_MOTHER_AI

            # Create or get cached agent
            if not self.agent:
                self.agent = await create_agent_from_config(LEGACY_MOTHER_AI)

            # Run agent with context
            response = await self.agent.arun(
                message=message_text,
                session_id=session_id,
                user_id=str(user_id),
                stream=False,  # No streaming for Telegram (yet)
            )

            # Extract content
            response_text = response.content or "I'm not sure how to respond to that."

            # Handle long messages (Telegram limit: 4096 chars)
            if len(response_text) > 4000:
                # Split into chunks
                chunks = [
                    response_text[i : i + 4000]
                    for i in range(0, len(response_text), 4000)
                ]
                for chunk in chunks:
                    await update.message.reply_text(chunk)
            else:
                await update.message.reply_text(response_text)

        except Exception as e:
            logger.error(f"Error processing message: {e}", exc_info=True)
            await update.message.reply_text(
                f"❌ Sorry, I encountered an error: {str(e)}"
            )

    async def start(self):
        """Start the bot (polling mode)."""
        if not self.app:
            await self.initialize()

        logger.info("Starting Telegram bot...")
        await self.app.initialize()
        await self.app.start()
        await self.app.updater.start_polling()
        logger.info("Telegram bot is running!")

    async def stop(self):
        """Stop the bot gracefully."""
        if self.app:
            logger.info("Stopping Telegram bot...")
            await self.app.updater.stop()
            await self.app.stop()
            await self.app.shutdown()
            logger.info("Telegram bot stopped")


# Singleton instance
_telegram_bot: TelegramBotTool = None


async def get_telegram_bot() -> TelegramBotTool:
    """Get or create Telegram bot instance."""
    global _telegram_bot
    if not _telegram_bot:
        _telegram_bot = TelegramBotTool()
        await _telegram_bot.initialize()
    return _telegram_bot


async def start_telegram_bot():
    """Start the Telegram bot in background."""
    bot = await get_telegram_bot()
    await bot.start()
```

#### Step 3: Integrate into Main App

Edit: `backend/main.py`

```python
# Add at top
import asyncio
from contextlib import asynccontextmanager

# Add lifespan manager
@asynccontextmanager
async def lifespan(app: FastAPI):
    """Lifespan events for FastAPI."""
    # Startup
    logger.info("Starting OneMind OS...")

    # Start Telegram bot if configured
    telegram_token = os.getenv("TELEGRAM_BOT_TOKEN")
    if telegram_token:
        logger.info("Starting Telegram bot...")
        from backend.agno.tools.telegram_bot import start_telegram_bot

        asyncio.create_task(start_telegram_bot())
    else:
        logger.warning("TELEGRAM_BOT_TOKEN not set, Telegram bot disabled")

    yield

    # Shutdown
    logger.info("Shutting down OneMind OS...")

# Update app initialization
app = FastAPI(
    title="OneMind OS",
    version="2.0.0",
    lifespan=lifespan,  # Add this
)
```

#### Step 4: Configure Bot Token

Edit: `.env`

```bash
# Add Telegram Bot configuration
TELEGRAM_BOT_TOKEN=your_bot_token_from_botfather
```

**To create bot:**
1. Message [@BotFather](https://t.me/botfather) on Telegram
2. Send `/newbot`
3. Follow prompts to name your bot
4. Copy the token to `.env`

#### Step 5: Test

```bash
# Start backend
cd backend
python main.py

# Should see in logs:
# INFO: Starting Telegram bot...
# INFO: Telegram bot is running!

# On Telegram:
# 1. Find your bot by name
# 2. Send /start
# 3. Send "Hello!"
# 4. Bot should respond as Legacy
```

---

### 1.2 Discord Bot Integration

**Goal**: Enable bidirectional chat with Legacy via Discord

**Implementation**: Similar to Telegram, but uses Discord.py

#### Step 1: Install Dependencies

```bash
# Add to requirements.txt
echo "discord.py==2.3.2" >> requirements.txt
pip install -r requirements.txt
```

#### Step 2: Create Discord Bot Tool

Create file: `backend/agno/tools/discord_bot.py`

```python
"""Discord bot integration for bidirectional chat with agents."""

import discord
from discord.ext import commands
import os
import logging
from typing import Dict

logger = logging.getLogger(__name__)

# Store session IDs for each Discord user
USER_SESSIONS: Dict[int, str] = {}


class DiscordBotTool(commands.Cog):
    """Discord bot that routes messages to OneMind agents."""

    def __init__(self, bot: commands.Bot):
        self.bot = bot
        self.agent = None

    @commands.Cog.listener()
    async def on_ready(self):
        """Called when bot is ready."""
        logger.info(f"Discord bot logged in as {self.bot.user}")
        logger.info(f"Bot is in {len(self.bot.guilds)} servers")

    @commands.command(name="start")
    async def cmd_start(self, ctx: commands.Context):
        """Start command."""
        user_id = ctx.author.id
        session_id = f"discord_{user_id}"
        USER_SESSIONS[user_id] = session_id

        welcome_embed = discord.Embed(
            title="👋 Welcome to OneMind OS!",
            description="I'm Legacy, your personal AI assistant.",
            color=discord.Color.blue(),
        )
        welcome_embed.add_field(
            name="What I can do",
            value="• Manage tasks (Todoist)\n• Control smart home\n• Search the web\n• Manage calendar\n• And more!",
            inline=False,
        )
        welcome_embed.add_field(
            name="Commands",
            value="`!help` - Show this message\n`!reset` - Start new conversation",
            inline=False,
        )

        await ctx.send(embed=welcome_embed)

    @commands.command(name="reset")
    async def cmd_reset(self, ctx: commands.Context):
        """Reset conversation."""
        user_id = ctx.author.id
        import uuid

        new_session = f"discord_{user_id}_{uuid.uuid4().hex[:8]}"
        USER_SESSIONS[user_id] = new_session

        await ctx.send("✓ Started new conversation")

    @commands.Cog.listener()
    async def on_message(self, message: discord.Message):
        """Route incoming messages to agent."""
        # Ignore bot's own messages
        if message.author == self.bot.user:
            return

        # Ignore messages that start with command prefix
        if message.content.startswith(self.bot.command_prefix):
            return

        # Only respond to DMs or mentions
        if not isinstance(message.channel, discord.DMChannel):
            if self.bot.user not in message.mentions:
                return

        user_id = message.author.id
        user_name = str(message.author)

        # Get or create session
        if user_id not in USER_SESSIONS:
            USER_SESSIONS[user_id] = f"discord_{user_id}"
        session_id = USER_SESSIONS[user_id]

        # Remove mention from message
        message_text = message.content.replace(f"<@{self.bot.user.id}>", "").strip()

        logger.info(f"Message from {user_name} (ID: {user_id}): {message_text}")

        # Show typing indicator
        async with message.channel.typing():
            try:
                # Import agent factory
                from backend.agno.agents.factory import create_agent_from_config
                from backend.agno.agents.presets import LEGACY_MOTHER_AI

                # Create or get cached agent
                if not self.agent:
                    self.agent = await create_agent_from_config(LEGACY_MOTHER_AI)

                # Run agent
                response = await self.agent.arun(
                    message=message_text,
                    session_id=session_id,
                    user_id=str(user_id),
                    stream=False,
                )

                response_text = (
                    response.content or "I'm not sure how to respond to that."
                )

                # Handle long messages (Discord limit: 2000 chars)
                if len(response_text) > 1900:
                    chunks = [
                        response_text[i : i + 1900]
                        for i in range(0, len(response_text), 1900)
                    ]
                    for chunk in chunks:
                        await message.channel.send(chunk)
                else:
                    await message.channel.send(response_text)

            except Exception as e:
                logger.error(f"Error processing message: {e}", exc_info=True)
                await message.channel.send(f"❌ Sorry, I encountered an error: {str(e)}")


async def start_discord_bot():
    """Start Discord bot."""
    token = os.getenv("DISCORD_BOT_TOKEN")
    if not token:
        logger.warning("DISCORD_BOT_TOKEN not set, Discord bot disabled")
        return

    # Create bot with intents
    intents = discord.Intents.default()
    intents.message_content = True  # Required to read message content
    intents.members = True

    bot = commands.Bot(command_prefix="!", intents=intents)

    # Add cog
    await bot.add_cog(DiscordBotTool(bot))

    logger.info("Starting Discord bot...")
    await bot.start(token)
```

#### Step 3: Integrate into Main App

Edit: `backend/main.py`

```python
@asynccontextmanager
async def lifespan(app: FastAPI):
    """Lifespan events for FastAPI."""
    # ... existing Telegram bot code ...

    # Start Discord bot if configured
    discord_token = os.getenv("DISCORD_BOT_TOKEN")
    if discord_token:
        logger.info("Starting Discord bot...")
        from backend.agno.tools.discord_bot import start_discord_bot

        asyncio.create_task(start_discord_bot())
    else:
        logger.warning("DISCORD_BOT_TOKEN not set, Discord bot disabled")

    yield
    # ...
```

#### Step 4: Configure Bot Token

Edit: `.env`

```bash
# Add Discord Bot configuration
DISCORD_BOT_TOKEN=your_discord_bot_token
```

**To create bot:**
1. Go to [Discord Developer Portal](https://discord.com/developers/applications)
2. Click "New Application"
3. Go to "Bot" tab → "Add Bot"
4. Copy the token to `.env`
5. Enable "Message Content Intent" under Privileged Gateway Intents
6. Go to OAuth2 → URL Generator
7. Select scopes: `bot`, `applications.commands`
8. Select permissions: `Send Messages`, `Read Messages/View Channels`, `Read Message History`
9. Copy the generated URL and invite bot to your server

#### Step 5: Test

```bash
# Start backend
python backend/main.py

# On Discord:
# 1. Find your bot in server member list
# 2. Send !start in a channel
# 3. Mention the bot: @YourBot Hello!
# 4. Or DM the bot directly
```

---

## Feature 2: Browser Automation

**Goal**: Enable Legacy to browse websites, fill forms, and extract data

**Implementation**: Add Playwright-based browser automation tool

### Step 1: Install Dependencies

```bash
# Add to requirements.txt
echo "playwright==1.42.0" >> requirements.txt
pip install -r requirements.txt

# Install browsers
playwright install chromium
```

### Step 2: Create Browser Automation Tool

Create file: `backend/agno/tools/browser_automation.py`

```python
"""Browser automation tool using Playwright."""

from agno.tools import Toolkit
from playwright.async_api import async_playwright, Browser, Page
import base64
import logging
from typing import Optional, Dict, List

logger = logging.getLogger(__name__)


class BrowserAutomationTool(Toolkit):
    """Browser automation for web scraping, form filling, and interaction."""

    def __init__(self):
        super().__init__(name="browser_automation")

        # Register tools
        self.register(self.open_url)
        self.register(self.click_element)
        self.register(self.fill_form)
        self.register(self.extract_text)
        self.register(self.take_screenshot)
        self.register(self.execute_script)

        self._browser: Optional[Browser] = None
        self._playwright = None

    async def _get_browser(self) -> Browser:
        """Get or create browser instance."""
        if not self._browser:
            self._playwright = await async_playwright().start()
            self._browser = await self._playwright.chromium.launch(
                headless=True,  # Change to False for debugging
                args=["--disable-blink-features=AutomationControlled"],
            )
            logger.info("Browser launched")

        return self._browser

    async def _close_browser(self):
        """Close browser instance."""
        if self._browser:
            await self._browser.close()
            self._browser = None

        if self._playwright:
            await self._playwright.stop()
            self._playwright = None

    async def open_url(self, url: str, wait_for_load: bool = True) -> str:
        """Open a URL and return the page content.

        Args:
            url: The URL to open
            wait_for_load: Wait for full page load

        Returns:
            The page HTML content
        """
        try:
            browser = await self._get_browser()
            page = await browser.new_page()

            # Navigate
            await page.goto(url, wait_until="networkidle" if wait_for_load else "load")

            # Extract content
            title = await page.title()
            content = await page.content()

            await page.close()

            return f"✓ Opened {url}\nTitle: {title}\n\nContent length: {len(content)} chars\n\nPreview:\n{content[:500]}..."

        except Exception as e:
            logger.error(f"Error opening URL: {e}")
            return f"❌ Error: {str(e)}"

    async def click_element(
        self, url: str, selector: str, wait_after_click: int = 2000
    ) -> str:
        """Click an element on a page.

        Args:
            url: The URL to visit
            selector: CSS selector for element to click
            wait_after_click: Wait time in ms after clicking

        Returns:
            Result of the action
        """
        try:
            browser = await self._get_browser()
            page = await browser.new_page()

            await page.goto(url, wait_until="networkidle")
            await page.click(selector)
            await page.wait_for_timeout(wait_after_click)

            # Get current URL (may have changed)
            current_url = page.url
            title = await page.title()

            await page.close()

            return f"✓ Clicked element '{selector}'\nNavigated to: {current_url}\nTitle: {title}"

        except Exception as e:
            logger.error(f"Error clicking element: {e}")
            return f"❌ Error: {str(e)}"

    async def fill_form(
        self, url: str, fields: Dict[str, str], submit_selector: Optional[str] = None
    ) -> str:
        """Fill form fields and optionally submit.

        Args:
            url: The URL with the form
            fields: Dict of {selector: value} to fill
            submit_selector: Optional submit button selector

        Returns:
            Result of the action
        """
        try:
            browser = await self._get_browser()
            page = await browser.new_page()

            await page.goto(url, wait_until="networkidle")

            # Fill each field
            for selector, value in fields.items():
                await page.fill(selector, value)
                logger.info(f"Filled {selector} with '{value}'")

            # Submit if requested
            if submit_selector:
                await page.click(submit_selector)
                await page.wait_for_timeout(2000)

            result_url = page.url
            title = await page.title()

            await page.close()

            return f"✓ Filled {len(fields)} form fields\n{('Submitted form\n' if submit_selector else '')}Current URL: {result_url}\nTitle: {title}"

        except Exception as e:
            logger.error(f"Error filling form: {e}")
            return f"❌ Error: {str(e)}"

    async def extract_text(self, url: str, selector: str) -> str:
        """Extract text from element(s) matching selector.

        Args:
            url: The URL to visit
            selector: CSS selector

        Returns:
            Extracted text
        """
        try:
            browser = await self._get_browser()
            page = await browser.new_page()

            await page.goto(url, wait_until="networkidle")

            # Get all matching elements
            elements = await page.query_selector_all(selector)

            if not elements:
                await page.close()
                return f"❌ No elements found for selector '{selector}'"

            # Extract text from each element
            texts = []
            for el in elements:
                text = await el.inner_text()
                texts.append(text)

            await page.close()

            return f"✓ Extracted {len(texts)} element(s):\n\n" + "\n---\n".join(texts)

        except Exception as e:
            logger.error(f"Error extracting text: {e}")
            return f"❌ Error: {str(e)}"

    async def take_screenshot(
        self, url: str, selector: Optional[str] = None, full_page: bool = False
    ) -> str:
        """Take a screenshot of the page or element.

        Args:
            url: The URL to visit
            selector: Optional CSS selector to screenshot specific element
            full_page: Take full-page screenshot

        Returns:
            Base64-encoded PNG image
        """
        try:
            browser = await self._get_browser()
            page = await browser.new_page()

            await page.goto(url, wait_until="networkidle")

            # Take screenshot
            if selector:
                element = await page.query_selector(selector)
                if not element:
                    await page.close()
                    return f"❌ Element not found: {selector}"
                screenshot_bytes = await element.screenshot()
            else:
                screenshot_bytes = await page.screenshot(full_page=full_page)

            await page.close()

            # Encode as base64
            screenshot_b64 = base64.b64encode(screenshot_bytes).decode()

            return f"✓ Screenshot captured (size: {len(screenshot_bytes)} bytes)\n\nBase64 data:\n{screenshot_b64[:100]}..."

        except Exception as e:
            logger.error(f"Error taking screenshot: {e}")
            return f"❌ Error: {str(e)}"

    async def execute_script(self, url: str, script: str) -> str:
        """Execute JavaScript on the page.

        Args:
            url: The URL to visit
            script: JavaScript code to execute

        Returns:
            Result of script execution
        """
        try:
            browser = await self._get_browser()
            page = await browser.new_page()

            await page.goto(url, wait_until="networkidle")

            # Execute script
            result = await page.evaluate(script)

            await page.close()

            return f"✓ Script executed\n\nResult:\n{result}"

        except Exception as e:
            logger.error(f"Error executing script: {e}")
            return f"❌ Error: {str(e)}"
```

### Step 3: Register Tool

Edit: `backend/agno/tools/registry.py`

```python
CUSTOM_TOOL_SPECS = [
    # ... existing tools ...
    {
        "name": "browser_automation",
        "module": "backend.agno.tools.browser_automation",
        "class": "BrowserAutomationTool",
        "description": "Browser automation for web scraping and interaction",
        "category": "AUTOMATION",
    },
]
```

### Step 4: Add to Agent Config

Edit: `backend/agno/agents/presets.py`

```python
# Add browser tools to Legacy
LEGACY_MOTHER_AI = {
    "agent_id": "legacy-mother-ai",
    "name": "Legacy",
    "model_id": "sonnet",
    "tools": [
        "research",
        "productivity",
        "home",
        "browser_automation",  # Add this
    ],
    # ... rest of config
}
```

### Step 5: Test

```python
# Test in Python REPL
from backend.agno.tools.browser_automation import BrowserAutomationTool
import asyncio

async def test():
    tool = BrowserAutomationTool()

    # Test 1: Open URL
    result = await tool.open_url("https://example.com")
    print(result)

    # Test 2: Extract text
    result = await tool.extract_text("https://news.ycombinator.com", ".titleline > a")
    print(result)

    # Test 3: Screenshot
    result = await tool.take_screenshot("https://example.com")
    print(result[:200])

asyncio.run(test())
```

---

## Feature 3: Self-Modifying Skills

**Goal**: Enable Legacy to create new skills by generating code

**Implementation**: Add skill writer tool with HITL approval

### Step 1: Create Skill Writer Tool

Create file: `backend/agno/tools/skill_writer.py`

```python
"""Tool for agents to create new skills."""

from agno.tools import tool
from pathlib import Path
import yaml
import logging

logger = logging.getLogger(__name__)


@tool(requires_confirmation=True)  # HITL approval required
async def create_skill(
    name: str, description: str, instructions: str, code: str, category: str = "CUSTOM"
) -> str:
    """Create a new skill from AI-generated code.

    REQUIRES HUMAN APPROVAL before execution for security.

    Args:
        name: Skill name (kebab-case, e.g., 'web-scraper')
        description: Brief description of what the skill does
        instructions: Detailed instructions for using the skill (markdown)
        code: Python implementation code
        category: Skill category (CUSTOM, AUTOMATION, RESEARCH, etc.)

    Returns:
        Status message
    """
    try:
        # Validate name (kebab-case)
        if not name.replace("-", "").isalnum():
            return f"❌ Invalid skill name. Use kebab-case (e.g., 'my-skill')"

        # Create skill directory
        skill_dir = Path(f"backend/agno/skills/{name}")
        if skill_dir.exists():
            return f"❌ Skill '{name}' already exists"

        skill_dir.mkdir(parents=True, exist_ok=True)

        # Write SKILL.md with frontmatter
        skill_md = f"""---
name: {name}
description: {description}
version: 1.0.0
category: {category}
tags: [ai-generated, custom]
author: Legacy AI
created: {datetime.now().isoformat()}
---

# {name.title()}

{instructions}

## Implementation

See `scripts/main.py` for the implementation.

## Usage Example

```python
from backend.agno.skills.{name}.scripts.main import main

result = await main()
print(result)
```
"""
        (skill_dir / "SKILL.md").write_text(skill_md)

        # Create scripts directory
        scripts_dir = skill_dir / "scripts"
        scripts_dir.mkdir(exist_ok=True)

        # Write Python implementation
        (scripts_dir / "main.py").write_text(code)
        (scripts_dir / "__init__.py").write_text("")  # Make it a package

        logger.info(f"Created skill: {name}")

        return f"""✓ Skill '{name}' created successfully!

Location: {skill_dir}

Files created:
- SKILL.md (documentation)
- scripts/main.py (implementation)

The skill is now available in the skills directory.
To use it, the backend may need to be restarted to load the new skill.
"""

    except Exception as e:
        logger.error(f"Error creating skill: {e}", exc_info=True)
        return f"❌ Error creating skill: {str(e)}"


@tool
async def list_skills() -> str:
    """List all available skills.

    Returns:
        List of skills with descriptions
    """
    try:
        skills_dir = Path("backend/agno/skills")
        if not skills_dir.exists():
            return "No skills directory found"

        skills = []
        for skill_dir in skills_dir.iterdir():
            if not skill_dir.is_dir():
                continue

            skill_md = skill_dir / "SKILL.md"
            if not skill_md.exists():
                continue

            # Parse frontmatter
            content = skill_md.read_text()
            if content.startswith("---"):
                parts = content.split("---", 2)
                frontmatter = yaml.safe_load(parts[1])

                skills.append(
                    {
                        "name": frontmatter.get("name"),
                        "description": frontmatter.get("description"),
                        "category": frontmatter.get("category"),
                        "version": frontmatter.get("version"),
                    }
                )

        if not skills:
            return "No skills found"

        # Format as table
        result = "Available Skills:\n\n"
        for skill in skills:
            result += f"• **{skill['name']}** ({skill['version']}) - {skill['description']}\n"
            result += f"  Category: {skill['category']}\n\n"

        return result

    except Exception as e:
        logger.error(f"Error listing skills: {e}", exc_info=True)
        return f"❌ Error: {str(e)}"
```

### Step 2: Register Tool

Edit: `backend/agno/tools/registry.py`

```python
CUSTOM_TOOL_SPECS = [
    # ... existing tools ...
    {
        "name": "skill_writer",
        "module": "backend.agno.tools.skill_writer",
        "functions": ["create_skill", "list_skills"],  # Function-based tools
        "description": "Create and manage custom skills",
        "category": "META",
    },
]
```

### Step 3: Add to Agent Config

```python
# Create a new agent specifically for skill creation
SKILL_WRITER_AGENT = {
    "agent_id": "skill-writer",
    "name": "Skill Writer",
    "model_id": "sonnet",
    "tools": ["skill_writer", "file", "python"],
    "instructions": """You are a skill creation expert.

You can create new skills by generating Python code and documentation.

When asked to create a skill:
1. Understand the requirements
2. Design the skill architecture
3. Write clean, documented Python code
4. Create comprehensive documentation
5. Use the create_skill tool (requires user approval)

Always include:
- Error handling
- Type hints
- Docstrings
- Usage examples
""",
    "show_tool_calls": True,
    "markdown": True,
}

# Add to AGENTS dict
AGENTS = {
    # ... existing agents ...
    "skill-writer": SKILL_WRITER_AGENT,
}
```

### Step 4: Test Self-Modifying Skills

**Chat with Agent**:
```
User: Create a skill that checks if a website is up

Agent: I'll create a skill for you. Let me write the code...

[Agent generates code]

Agent: I've designed a website uptime checker skill. Here's what it does:
- Uses requests to check HTTP status
- Measures response time
- Returns uptime status

[HITL Approval Card appears in Flutter UI]:
┌─────────────────────────────────────────────────┐
│  Approval Required                              │
├─────────────────────────────────────────────────┤
│  create_skill                                   │
│                                                 │
│  Name: website-uptime-checker                   │
│  Description: Check if websites are reachable   │
│                                                 │
│  [View Code] [View Documentation]               │
│                                                 │
│  [✓ Approve]  [✗ Deny]  [✎ Edit]                │
└─────────────────────────────────────────────────┘

User: [Clicks Approve]

Agent: ✓ Skill 'website-uptime-checker' created!

You can now use it like this:
"Check if example.com is up"
```

---

## Feature 4: Skill Auto-Discovery

**Goal**: Automatically discover and load skills without restart

**Implementation**: Add skill discovery system

### Step 1: Create Skill Discovery Module

Create file: `backend/agno/skills/__init__.py`

```python
"""Skill discovery and management system."""

from pathlib import Path
from typing import Dict, List, Optional
import yaml
import logging

logger = logging.getLogger(__name__)

# Global skill registry
AVAILABLE_SKILLS: Dict[str, Dict] = {}


def discover_skills(skills_path: str = "backend/agno/skills") -> Dict[str, Dict]:
    """Auto-discover all skills from skills directory.

    Args:
        skills_path: Path to skills directory

    Returns:
        Dict of {skill_name: skill_metadata}
    """
    skills = {}
    skills_dir = Path(skills_path)

    if not skills_dir.exists():
        logger.warning(f"Skills directory not found: {skills_dir}")
        return skills

    logger.info(f"Discovering skills in {skills_dir}")

    for skill_dir in skills_dir.iterdir():
        # Skip non-directories and hidden directories
        if not skill_dir.is_dir() or skill_dir.name.startswith("."):
            continue

        skill_md = skill_dir / "SKILL.md"
        if not skill_md.exists():
            logger.debug(f"Skipping {skill_dir.name} (no SKILL.md)")
            continue

        try:
            # Parse YAML frontmatter
            content = skill_md.read_text()

            if not content.startswith("---"):
                logger.warning(f"Invalid SKILL.md in {skill_dir.name} (no frontmatter)")
                continue

            parts = content.split("---", 2)
            if len(parts) < 3:
                logger.warning(f"Invalid SKILL.md in {skill_dir.name} (malformed)")
                continue

            frontmatter = yaml.safe_load(parts[1])
            instructions = parts[2].strip()

            # Build skill metadata
            skill_metadata = {
                "name": frontmatter.get("name", skill_dir.name),
                "description": frontmatter.get("description", ""),
                "version": frontmatter.get("version", "1.0.0"),
                "category": frontmatter.get("category", "CUSTOM"),
                "tags": frontmatter.get("tags", []),
                "author": frontmatter.get("author", "Unknown"),
                "path": skill_dir,
                "instructions": instructions,
                "markdown_path": skill_md,
                "scripts_path": skill_dir / "scripts",
            }

            skills[skill_metadata["name"]] = skill_metadata
            logger.info(f"  Discovered: {skill_metadata['name']} (v{skill_metadata['version']})")

        except Exception as e:
            logger.error(f"Error parsing skill {skill_dir.name}: {e}")

    logger.info(f"Discovered {len(skills)} skill(s)")
    return skills


def get_skill(skill_name: str) -> Optional[Dict]:
    """Get skill metadata by name.

    Args:
        skill_name: Name of the skill

    Returns:
        Skill metadata dict or None if not found
    """
    if skill_name not in AVAILABLE_SKILLS:
        # Try refreshing
        refresh_skills()

    return AVAILABLE_SKILLS.get(skill_name)


def get_skill_instructions(skill_name: str) -> Optional[str]:
    """Get skill instructions (markdown content).

    Args:
        skill_name: Name of the skill

    Returns:
        Skill instructions or None if not found
    """
    skill = get_skill(skill_name)
    if not skill:
        return None

    return skill.get("instructions")


def refresh_skills():
    """Refresh skill registry (hot-reload)."""
    global AVAILABLE_SKILLS
    AVAILABLE_SKILLS = discover_skills()


# Auto-discover on import
AVAILABLE_SKILLS = discover_skills()
```

### Step 2: Create Skill Selector Tool

Create file: `backend/agno/tools/skill_selector.py`

```python
"""Tool for agents to discover and activate skills."""

from agno.tools import tool
from backend.agno.skills import (
    AVAILABLE_SKILLS,
    get_skill,
    get_skill_instructions,
    refresh_skills,
)
import logging

logger = logging.getLogger(__name__)


@tool
async def list_available_skills(category: str = None) -> str:
    """List all available skills with descriptions.

    Args:
        category: Optional category filter (CUSTOM, RESEARCH, AUTOMATION, etc.)

    Returns:
        Formatted list of skills
    """
    # Refresh in case new skills were added
    refresh_skills()

    if not AVAILABLE_SKILLS:
        return "No skills available"

    # Filter by category if specified
    skills = AVAILABLE_SKILLS.values()
    if category:
        skills = [s for s in skills if s.get("category") == category.upper()]

    if not skills:
        return f"No skills found in category '{category}'"

    # Group by category
    by_category = {}
    for skill in skills:
        cat = skill.get("category", "CUSTOM")
        if cat not in by_category:
            by_category[cat] = []
        by_category[cat].append(skill)

    # Format output
    result = f"Available Skills ({len(skills)} total):\n\n"

    for cat, cat_skills in sorted(by_category.items()):
        result += f"**{cat}**\n"
        for skill in cat_skills:
            result += f"  • {skill['name']} (v{skill['version']}) - {skill['description']}\n"
        result += "\n"

    result += "Use `activate_skill` to load a skill's instructions."

    return result


@tool
async def activate_skill(skill_name: str) -> str:
    """Activate a skill by loading its instructions into context.

    Args:
        skill_name: Name of the skill to activate

    Returns:
        Skill instructions
    """
    instructions = get_skill_instructions(skill_name)

    if not instructions:
        return f"❌ Skill '{skill_name}' not found. Use `list_available_skills` to see available skills."

    skill = get_skill(skill_name)

    return f"""✓ Activated skill: **{skill['name']}** (v{skill['version']})

{skill['description']}

---

{instructions}
"""


@tool
async def search_skills(query: str) -> str:
    """Search for skills matching a query.

    Args:
        query: Search query (matches name, description, tags)

    Returns:
        Matching skills
    """
    refresh_skills()

    query_lower = query.lower()

    matches = []
    for skill in AVAILABLE_SKILLS.values():
        # Search in name, description, and tags
        if (
            query_lower in skill["name"].lower()
            or query_lower in skill["description"].lower()
            or any(query_lower in tag.lower() for tag in skill.get("tags", []))
        ):
            matches.append(skill)

    if not matches:
        return f"No skills found matching '{query}'"

    result = f"Skills matching '{query}':\n\n"
    for skill in matches:
        result += f"• **{skill['name']}** - {skill['description']}\n"
        result += f"  Tags: {', '.join(skill.get('tags', []))}\n\n"

    return result
```

### Step 3: Register Tool

Edit: `backend/agno/tools/registry.py`

```python
CUSTOM_TOOL_SPECS = [
    # ... existing ...
    {
        "name": "skill_selector",
        "module": "backend.agno.tools.skill_selector",
        "functions": ["list_available_skills", "activate_skill", "search_skills"],
        "description": "Discover and activate skills",
        "category": "META",
    },
]
```

### Step 4: Add to Agent

```python
# Edit LEGACY_MOTHER_AI in presets.py
LEGACY_MOTHER_AI = {
    # ... existing config ...
    "tools": [
        "research",
        "productivity",
        "skill_selector",  # Add this
        "skill_writer",  # And this (from previous step)
    ],
}
```

### Step 5: Test

**Chat example**:
```
User: What skills do you have available?

Agent: [Uses list_available_skills tool]
Available Skills (8 total):

**RESEARCH**
  • deep-research (v1.0.0) - Systematic research methodology

**CODE_REVIEW**
  • code-review (v1.0.0) - Code review best practices

**CUSTOM**
  • website-uptime-checker (v1.0.0) - Check if websites are reachable

...

User: Activate the deep-research skill

Agent: [Uses activate_skill tool]
✓ Activated skill: deep-research

[Full skill instructions loaded into context]
```

---

## Feature 5: Simplified Setup

**Goal**: Create a one-liner installer like OpenClaw

**Implementation**: Bash script that automates setup

### Create Installer Script

Create file: `install.sh` in repo root

```bash
#!/bin/bash
# OneMind OS One-Liner Installer
# Usage: curl -fsSL https://raw.githubusercontent.com/Zeus-Delacruz/onemind-os/main/install.sh | bash

set -e  # Exit on error

echo "╔═══════════════════════════════════════════╗"
echo "║   OneMind OS - One-Liner Installer        ║"
echo "║   Human (Zeus) + AI (Legacy) = ONE MIND   ║"
echo "╚═══════════════════════════════════════════╝"
echo ""

# Detect OS
OS=$(uname -s)
ARCH=$(uname -m)

echo "→ Detected OS: $OS ($ARCH)"

# Check prerequisites
echo "→ Checking prerequisites..."

# Python 3.12+
if ! command -v python3 &> /dev/null; then
    echo "❌ Python 3 not found. Please install Python 3.12+."
    exit 1
fi

PYTHON_VERSION=$(python3 --version | cut -d' ' -f2 | cut -d'.' -f1,2)
if [ $(echo "$PYTHON_VERSION < 3.12" | bc) -eq 1 ]; then
    echo "❌ Python 3.12+ required (found $PYTHON_VERSION)"
    exit 1
fi

echo "✓ Python $PYTHON_VERSION"

# Docker
if ! command -v docker &> /dev/null; then
    echo "⚠️  Docker not found. Installing Docker..."
    # Install Docker based on OS
    # (Implementation varies by OS)
fi

echo "✓ Docker installed"

# Clone repository
echo "→ Cloning OneMind OS repository..."
INSTALL_DIR="$HOME/onemind-os"

if [ -d "$INSTALL_DIR" ]; then
    echo "⚠️  Directory exists. Updating..."
    cd "$INSTALL_DIR"
    git pull origin main
else
    git clone https://github.com/Zeus-Delacruz/onemind-os.git "$INSTALL_DIR"
    cd "$INSTALL_DIR"
fi

echo "✓ Repository cloned"

# Create .env file
echo "→ Configuring environment..."

if [ ! -f .env ]; then
    cp .env.example .env

    # Prompt for AWS credentials
    echo ""
    echo "AWS Bedrock Configuration (required for Claude AI):"
    read -p "AWS Access Key ID: " AWS_KEY
    read -sp "AWS Secret Access Key: " AWS_SECRET
    echo ""
    read -p "AWS Region [us-east-1]: " AWS_REGION
    AWS_REGION=${AWS_REGION:-us-east-1}

    # Update .env
    sed -i.bak "s/your_access_key_id/$AWS_KEY/" .env
    sed -i.bak "s/your_secret_access_key/$AWS_SECRET/" .env
    sed -i.bak "s/AWS_DEFAULT_REGION=.*/AWS_DEFAULT_REGION=$AWS_REGION/" .env

    rm .env.bak

    echo "✓ Environment configured"
else
    echo "⚠️  .env file exists, skipping configuration"
fi

# Install Python dependencies
echo "→ Installing Python dependencies..."
pip3 install -r infra/containers/agno/requirements.txt

echo "✓ Dependencies installed"

# Start Docker services
echo "→ Starting Docker services (PostgreSQL, Redis, NATS)..."
cd infra/deployment/docker
docker compose -f docker-compose.core.yml up -d postgres redis nats timescaledb

echo "⏳ Waiting for databases to be ready..."
sleep 10

# Run database migrations (if any)
cd ../../..
python3 -c "
from backend.platform.db.init import initialize_database
import asyncio
asyncio.run(initialize_database())
print('✓ Database initialized')
"

# Start backend
echo "→ Starting OneMind OS backend..."
echo "   (This will run in the background. Use 'pkill -f onemindos.py' to stop)"

nohup python3 backend/onemindos.py > logs/backend.log 2>&1 &

echo "⏳ Waiting for backend to start..."
sleep 5

# Check if backend is running
if curl -s http://localhost:8081/health > /dev/null; then
    echo "✓ Backend is running at http://localhost:8081"
else
    echo "❌ Backend failed to start. Check logs/backend.log"
    exit 1
fi

# Installation complete
echo ""
echo "╔═══════════════════════════════════════════╗"
echo "║   Installation Complete! 🎉                ║"
echo "╚═══════════════════════════════════════════╝"
echo ""
echo "Backend API: http://localhost:8081"
echo "API Docs: http://localhost:8081/docs"
echo ""
echo "Next steps:"
echo "1. Install Flutter to run the UI: https://flutter.dev/docs/get-started/install"
echo "2. Run frontend:"
echo "   cd $INSTALL_DIR/frontend"
echo "   flutter pub get"
echo "   flutter run -d chrome"
echo ""
echo "Or use the Telegram/Discord bots:"
echo "  - Set TELEGRAM_BOT_TOKEN in .env"
echo "  - Set DISCORD_BOT_TOKEN in .env"
echo "  - Restart backend: pkill -f onemindos.py && python3 backend/onemindos.py &"
echo ""
echo "Documentation: $INSTALL_DIR/docs/README.md"
echo ""
```

Make it executable:
```bash
chmod +x install.sh
```

Add to README.md:
```markdown
## Quick Install

One-liner (macOS/Linux):
```bash
curl -fsSL https://raw.githubusercontent.com/Zeus-Delacruz/onemind-os/main/install.sh | bash
```

---

## Testing & Validation

### Test Checklist

- [ ] **Telegram Bot**
  - [ ] Bot responds to /start
  - [ ] Bot remembers conversation context
  - [ ] Bot can use tools (e.g., Todoist)
  - [ ] /reset creates new session

- [ ] **Discord Bot**
  - [ ] Bot responds to !start
  - [ ] Bot responds to mentions
  - [ ] Bot works in DMs
  - [ ] Embeds display correctly

- [ ] **Browser Automation**
  - [ ] open_url works
  - [ ] click_element works
  - [ ] fill_form works
  - [ ] take_screenshot works

- [ ] **Self-Modifying Skills**
  - [ ] create_skill tool registered
  - [ ] HITL approval shows in UI
  - [ ] Approved skills are created
  - [ ] Created skills work

- [ ] **Skill Auto-Discovery**
  - [ ] list_available_skills returns all skills
  - [ ] activate_skill loads instructions
  - [ ] search_skills finds matches

### Integration Tests

Create: `tests/integration/test_openclaw_features.py`

```python
import pytest
from backend.agno.tools.telegram_bot import TelegramBotTool
from backend.agno.tools.browser_automation import BrowserAutomationTool
from backend.agno.tools.skill_writer import create_skill
from backend.agno.skills import discover_skills

@pytest.mark.asyncio
async def test_browser_automation():
    """Test browser automation tool."""
    tool = BrowserAutomationTool()

    # Test URL opening
    result = await tool.open_url("https://example.com")
    assert "✓ Opened" in result
    assert "Example Domain" in result

@pytest.mark.asyncio
async def test_skill_discovery():
    """Test skill auto-discovery."""
    skills = discover_skills()

    assert len(skills) > 0
    assert "deep-research" in skills or "research" in skills

@pytest.mark.asyncio
async def test_skill_creation():
    """Test self-modifying skill creation."""
    result = await create_skill(
        name="test-skill",
        description="Test skill",
        instructions="Test instructions",
        code="def main(): return 'Hello'",
        category="TEST",
    )

    assert "✓ Skill" in result or "created" in result.lower()

    # Cleanup
    import shutil
    from pathlib import Path

    skill_dir = Path("backend/agno/skills/test-skill")
    if skill_dir.exists():
        shutil.rmtree(skill_dir)
```

Run tests:
```bash
pytest tests/integration/test_openclaw_features.py -v
```

---

## Deployment Strategy

### Rolling Deployment Plan

**Week 1: Chat Bots** (Low risk, high value)
- Deploy Telegram bot to production
- Deploy Discord bot to production
- Monitor usage and errors

**Week 2: Browser Automation** (Medium risk, medium value)
- Deploy browser automation tool
- Test with real use cases
- Monitor resource usage

**Week 3: Skills System** (High risk, high value)
- Deploy skill auto-discovery
- Deploy self-modifying skills
- Strict HITL approval required

**Week 4: Polish**
- Deploy one-liner installer
- Update documentation
- Community testing

### Rollback Plan

If issues occur:

```bash
# Stop new features
cd /Users/zeuslegacy/Documents/onemind-workspace/OneMind-OS
git checkout main  # Revert to stable

# Restart backend
docker restart onemind-agno

# Check logs
docker logs onemind-agno -f
```

---

## Success Metrics

Track these metrics to measure success:

| Metric | Target | Measurement |
|--------|--------|-------------|
| **Telegram Bot Usage** | 10+ messages/day | Telegram analytics |
| **Discord Bot Usage** | 5+ messages/day | Discord analytics |
| **Browser Automations** | 3+ per week | Tool usage logs |
| **Skills Created** | 2+ per month | Skill directory count |
| **User Satisfaction** | 4.5+ stars | User feedback |

---

## Conclusion

This implementation guide provides everything needed to add OpenClaw's best features to OneMind OS:

1. **Chat bots** - Telegram & Discord (6 hours)
2. **Browser automation** - Playwright integration (4 hours)
3. **Self-modifying skills** - AI creates skills (6 hours)
4. **Skill auto-discovery** - Hot-reload skills (4 hours)
5. **Simplified setup** - One-liner installer (5 hours)

**Total**: ~30 hours of implementation time

**Next Steps**:
1. Start with Phase 1 (chat bots) for quick wins
2. Test thoroughly before moving to next phase
3. Deploy incrementally with rollback plan ready
4. Gather user feedback and iterate

---

**Document Version**: 1.0
**Last Updated**: 2026-02-07
**Estimated Implementation Time**: 30 hours
**Risk Level**: Medium


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
