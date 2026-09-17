# CLI Guide

Run AI agents from the terminal using Agno's built-in CLI capabilities.

## Agno Native CLI

Every Agno agent and team has a built-in CLI method:

```python
from agno.agent import Agent
from agno.models.openai import OpenAIChat

agent = Agent(
    name="Assistant",
    model=OpenAIChat(id="gpt-4o"),
)

# Run as interactive CLI
agent.cli_app(markdown=True)
```

## Agent CLI

### Basic Usage

```python
from agno.agent import Agent
from agno.models.openai import OpenAIChat
from agno.tools.duckduckgo import DuckDuckGoTools

agent = Agent(
    name="Research Agent",
    model=OpenAIChat(id="gpt-4o"),
    tools=[DuckDuckGoTools()],
    markdown=True,
)

# Start interactive CLI
agent.cli_app(markdown=True)
```

### With Session Persistence

```python
from agno.agent import Agent
from agno.db.sqlite import SqliteDb
from agno.models.openai import OpenAIChat

db = SqliteDb(db_file="agent.db")

agent = Agent(
    name="Persistent Agent",
    model=OpenAIChat(id="gpt-4o"),
    db=db,
    session_id="my_session",  # Resume same conversation
    add_history_to_context=True,
)

agent.cli_app(markdown=True)
```

### Async CLI

```python
import asyncio

async def main():
    agent = Agent(
        model=OpenAIChat(id="gpt-4o"),
    )
    await agent.acli_app(markdown=True)

asyncio.run(main())
```

## Team CLI

Run multi-agent teams from the terminal:

```python
from agno.agent import Agent
from agno.team import Team
from agno.models.openai import OpenAIChat

researcher = Agent(name="Researcher", model=OpenAIChat(id="gpt-4o"))
writer = Agent(name="Writer", model=OpenAIChat(id="gpt-4o"))

team = Team(
    name="Content Team",
    members=[researcher, writer],
    model=OpenAIChat(id="gpt-4o"),
)

# Run team as CLI
team.cli_app(stream=True)
```

## Building Custom CLI with Typer

For more control, build a custom CLI using Typer:

```python
import typer
from rich.console import Console
from rich.panel import Panel
from rich.markdown import Markdown
from agno.agent import Agent
from agno.models.openai import OpenAIChat

app = typer.Typer(name="myagent", help="My AI Agent CLI")
console = Console()

@app.command()
def chat(
    message: str = typer.Argument(None, help="Message to send"),
    model: str = typer.Option("gpt-4o", help="Model to use"),
    stream: bool = typer.Option(True, help="Stream responses"),
):
    """Chat with the AI agent."""
    agent = Agent(
        name="Assistant",
        model=OpenAIChat(id=model),
    )

    if message:
        # Single message mode
        response = agent.run(message, stream=stream)
        console.print(Panel(Markdown(response.content), title="Response"))
    else:
        # Interactive mode
        agent.cli_app(markdown=True)

@app.command()
def search(query: str):
    """Search the web."""
    from agno.tools.duckduckgo import DuckDuckGoTools

    agent = Agent(
        model=OpenAIChat(id="gpt-4o"),
        tools=[DuckDuckGoTools()],
    )

    response = agent.run(f"Search for: {query}")
    console.print(Panel(Markdown(response.content), title="Search Results"))

if __name__ == "__main__":
    app()
```

### Running Custom CLI

```bash
# Install dependencies
pip install typer rich

# Run commands
python my_cli.py chat "Hello, how are you?"
python my_cli.py chat  # Interactive mode
python my_cli.py search "latest AI news"
```

## CLI Options

### cli_app() Parameters

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `markdown` | bool | False | Render responses as markdown |
| `stream` | bool | True | Stream responses in real-time |
| `exit_on` | list | ["exit", "quit", "bye"] | Commands to exit CLI |

### Example with All Options

```python
agent.cli_app(
    markdown=True,
    stream=True,
    exit_on=["exit", "quit", "bye", "goodbye"],
)
```

## Rich Terminal Output

Use Rich library for beautiful terminal output:

```python
from rich.console import Console
from rich.panel import Panel
from rich.table import Table
from rich.markdown import Markdown
from rich.progress import Progress

console = Console()

# Panels
console.print(Panel("Welcome to OneMind", title="Info", border_style="cyan"))

# Tables
table = Table(title="Available Commands")
table.add_column("Command", style="cyan")
table.add_column("Description")
table.add_row("chat", "Start interactive chat")
table.add_row("search", "Search the web")
console.print(table)

# Markdown
console.print(Markdown("# Hello\n\nThis is **markdown**"))

# Progress
with Progress() as progress:
    task = progress.add_task("Processing...", total=100)
    # Update progress
    progress.update(task, advance=50)
```

## Session Management

### Resuming Sessions

```python
import typer
from agno.agent import Agent
from agno.db.sqlite import SqliteDb
from agno.session import AgentSession

db = SqliteDb(db_file="sessions.db")

@app.command()
def chat(
    session_id: str = typer.Option(None, help="Session to resume"),
    new: bool = typer.Option(False, help="Start new session"),
):
    """Chat with session support."""

    if not new and session_id is None:
        # List existing sessions
        sessions = db.get_sessions(user_id="default")
        if sessions:
            console.print("Existing sessions:")
            for s in sessions[:5]:
                console.print(f"  - {s.session_id}")
            session_id = sessions[0].session_id

    agent = Agent(
        model=OpenAIChat(id="gpt-4o"),
        db=db,
        session_id=session_id,
        add_history_to_context=True,
    )

    agent.cli_app(markdown=True)
```

## Environment Variables

Common environment variables for CLI agents:

```bash
# Model API keys
export OPENAI_API_KEY="sk-..."
export ANTHROPIC_API_KEY="sk-ant-..."

# Database
export DATABASE_URL="postgresql://..."

# Optional
export AGNO_DEBUG="true"
export AGNO_LOG_LEVEL="INFO"
```

## Best Practices

### 1. Handle Errors Gracefully

```python
@app.command()
def chat(message: str):
    try:
        agent = Agent(model=OpenAIChat(id="gpt-4o"))
        response = agent.run(message)
        console.print(response.content)
    except Exception as e:
        console.print(f"[red]Error:[/red] {e}")
        raise typer.Exit(1)
```

### 2. Add Help Text

```python
@app.command()
def analyze(
    file: str = typer.Argument(..., help="File to analyze"),
    verbose: bool = typer.Option(False, "--verbose", "-v", help="Show details"),
):
    """Analyze a file using AI.

    Examples:
        myagent analyze report.pdf
        myagent analyze data.csv --verbose
    """
    pass
```

### 3. Use Configuration Files

```python
import os
from pathlib import Path

CONFIG_DIR = Path.home() / ".config" / "myagent"
CONFIG_FILE = CONFIG_DIR / "config.yaml"

def load_config():
    if CONFIG_FILE.exists():
        import yaml
        return yaml.safe_load(CONFIG_FILE.read_text())
    return {}
```

## References

- [Agno Agent CLI](https://docs.agno.com/basics/agents/usage/agent-cli)
- [Typer Documentation](https://typer.tiangolo.com/)
- [Rich Documentation](https://rich.readthedocs.io/)
