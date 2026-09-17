# Claude Code

> **AI-Powered Coding Assistant**

## Overview

Claude Code (formerly Claude Dev) is an agentic AI coding assistant that can understand, modify, and create code across your entire project. It's the primary AI development tool for OneMind OS.

## Status

- **Active**: ✅ Yes
- **Model**: Claude Opus 4.5
- **Last Updated**: January 25, 2026

## Use Cases

1. **Code Generation** - Write new features
2. **Debugging** - Find and fix issues
3. **Refactoring** - Improve code structure
4. **Documentation** - Write docs and comments
5. **Testing** - Generate test cases

## Key Features

- Full codebase awareness
- File creation and editing
- Terminal command execution
- Multi-file changes
- Context-aware suggestions
- Conversation memory

## Capabilities

### Read & Understand

- Browse entire project structure
- Read any file in workspace
- Understand code relationships
- Analyze dependencies

### Write & Edit

- Create new files
- Modify existing files
- Multi-file edits in single action
- Preserve code style

### Execute

- Run terminal commands
- Execute tests
- Install packages
- Build and deploy

## Usage Patterns

### Feature Development

```
"Create a new tool for Legacy AI that integrates with Google Calendar.
It should:
- List upcoming events
- Create new events
- Update existing events
- Delete events

Use the existing tool patterns in agents/legacy/tools/"
```

### Bug Fixing

```
"The /agents/{id}/runs endpoint is returning a 500 error.
Check the error logs and fix the issue.
The error seems to be related to memory retrieval."
```

### Refactoring

```
"Refactor the consciousness module to use async/await
instead of callbacks. Maintain backward compatibility
with existing code."
```

### Documentation

```
"Add docstrings to all functions in agents/legacy/tools/
Follow the Google style guide for Python docstrings."
```

## Best Practices

### Clear Prompts

✅ Good:

```
"Create a new FastAPI endpoint at /health that returns
{status: 'ok', version: '1.0.0', timestamp: <current time>}"
```

❌ Vague:

```
"Add a health check"
```

### Context Provision

Reference specific files:

```
"Following the pattern in agents/legacy/tools/todoist.py,
create a new tool for Linear issue management."
```

### Iterative Development

Break large tasks into steps:

1. "First, create the data models..."
2. "Now, add the API endpoints..."
3. "Finally, write the tests..."

## Keyboard Shortcuts

| Shortcut      | Action                |
| ------------- | --------------------- |
| `Cmd+Shift+I` | Open Claude Code      |
| `Cmd+Enter`   | Submit message        |
| `Cmd+L`       | Clear conversation    |
| `Escape`      | Cancel current action |

## Configuration

### Settings

```json
{
  "claudeCode.model": "claude-opus-4-5",
  "claudeCode.autoApprove": ["read_file", "list_directory"],
  "claudeCode.showDiff": true,
  "claudeCode.maxTokens": 8192
}
```

### Permissions

Control which actions require approval:

- **Always Approve**: Read files, list directories
- **Ask First**: Edit files, run commands
- **Never Allow**: Delete files, force push

## Model Capabilities

### Claude Opus 4.5

| Capability        | Support      |
| ----------------- | ------------ |
| Code generation   | ✅ Excellent |
| Debugging         | ✅ Excellent |
| Refactoring       | ✅ Excellent |
| Documentation     | ✅ Excellent |
| Test generation   | ✅ Excellent |
| Multi-file edits  | ✅ Excellent |
| Terminal commands | ✅ Excellent |
| Web browsing      | ✅ Supported |

## Common Workflows

### 1. New Feature

```
1. Describe feature requirements
2. Claude generates implementation plan
3. Review and approve plan
4. Claude creates files and writes code
5. Claude runs tests
6. Review changes
7. Commit
```

### 2. Bug Fix

```
1. Describe the bug
2. Claude reads relevant files
3. Claude identifies the issue
4. Claude proposes fix
5. Review and approve
6. Claude applies fix
7. Claude runs tests
8. Commit
```

### 3. Codebase Questions

```
"How does the memory module work?"
"What's the flow from user input to agent response?"
"Where are the database models defined?"
```

## Integration with VS Code

### Copilot vs. Claude Code

| Feature               | Copilot    | Claude Code |
| --------------------- | ---------- | ----------- |
| Inline suggestions    | ✅ Best    | ❌ No       |
| Conversation          | ✅ Good    | ✅ Best     |
| Multi-file edits      | ⚠️ Limited | ✅ Best     |
| Terminal access       | ⚠️ Limited | ✅ Best     |
| Project understanding | ⚠️ Limited | ✅ Best     |

**Recommendation**: Use both!

- Copilot for quick inline completions
- Claude for complex multi-file tasks

## MCP Servers

### Connected Servers

- **GitHub** - Repository management
- **Context7** - Library documentation
- **Prisma** - Database operations
- **Docker** - Container management

### Example MCP Usage

```
"Use the GitHub MCP to create a new issue titled
'Implement Google Calendar integration' with the
appropriate labels and milestone."
```

## Security

- API key stored securely
- Actions require approval by default
- Audit log of all actions
- Never commits sensitive data

## Tips

- ✅ Start with clear, specific requests
- ✅ Review changes before approving
- ✅ Use file references for context
- ✅ Break large tasks into steps
- ✅ Ask Claude to explain its reasoning
- ✅ Iterate and refine

## Resources

- [Claude Documentation](https://docs.anthropic.com)
- [VS Code Extension](https://marketplace.visualstudio.com/items?itemName=anthropic.claude-code)

---

**Related**: [VS Code](./VS-Code.md) | [GitHub](./GitHub.md) | [Agno](./Agno.md)
