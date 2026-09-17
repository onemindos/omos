# VS Code

> **Primary Code Editor**

## Overview

Visual Studio Code is the primary code editor for OneMind OS development, enhanced with AI capabilities through GitHub Copilot and Claude integration.

## Status

- **Active**: ✅ Yes
- **Theme**: Dark+ (default)
- **Last Updated**: January 25, 2026

## Use Cases

1. **Code Development** - Python, TypeScript, Dart
2. **AI Pair Programming** - Copilot & Claude
3. **Debugging** - Integrated debugger
4. **Terminal** - Integrated terminal
5. **Git Integration** - Version control

## Essential Extensions

### AI & Coding Assistance

| Extension               | Purpose            |
| ----------------------- | ------------------ |
| **GitHub Copilot**      | AI code completion |
| **GitHub Copilot Chat** | AI conversation    |
| **Claude Dev**          | Claude integration |

### Language Support

| Extension    | Purpose                 |
| ------------ | ----------------------- |
| **Python**   | Python language support |
| **Pylance**  | Python IntelliSense     |
| **ESLint**   | JavaScript linting      |
| **Prettier** | Code formatting         |
| **Dart**     | Dart language support   |
| **Flutter**  | Flutter development     |

### Productivity

| Extension        | Purpose              |
| ---------------- | -------------------- |
| **GitLens**      | Enhanced Git         |
| **Docker**       | Container management |
| **Remote - SSH** | Remote development   |
| **Live Share**   | Collaboration        |
| **Error Lens**   | Inline errors        |

### Tools

| Extension               | Purpose          |
| ----------------------- | ---------------- |
| **REST Client**         | API testing      |
| **Thunder Client**      | API testing      |
| **YAML**                | YAML support     |
| **Markdown All in One** | Markdown editing |

## Settings

### settings.json

```json
{
  "editor.fontSize": 14,
  "editor.fontFamily": "JetBrains Mono, Menlo, monospace",
  "editor.fontLigatures": true,
  "editor.tabSize": 2,
  "editor.formatOnSave": true,
  "editor.minimap.enabled": false,
  "editor.wordWrap": "on",
  "editor.bracketPairColorization.enabled": true,
  "editor.guides.bracketPairs": true,

  "workbench.colorTheme": "One Dark Pro",
  "workbench.iconTheme": "material-icon-theme",

  "terminal.integrated.fontSize": 13,
  "terminal.integrated.defaultProfile.osx": "zsh",

  "files.autoSave": "afterDelay",
  "files.autoSaveDelay": 1000,
  "files.exclude": {
    "**/.git": true,
    "**/__pycache__": true,
    "**/.pytest_cache": true,
    "**/node_modules": true
  },

  "python.defaultInterpreterPath": "${workspaceFolder}/venv/bin/python",
  "python.formatting.provider": "black",

  "[python]": {
    "editor.defaultFormatter": "ms-python.black-formatter",
    "editor.formatOnSave": true
  },

  "[typescript]": {
    "editor.defaultFormatter": "esbenp.prettier-vscode"
  },

  "github.copilot.enable": {
    "*": true,
    "plaintext": false,
    "markdown": true
  }
}
```

## Keybindings

### Custom Keybindings

```json
[
  {
    "key": "cmd+shift+t",
    "command": "workbench.action.terminal.toggleTerminal"
  },
  {
    "key": "cmd+shift+e",
    "command": "workbench.files.action.focusFilesExplorer"
  },
  {
    "key": "cmd+k cmd+s",
    "command": "workbench.action.files.saveAll"
  }
]
```

### Essential Shortcuts

| Shortcut      | Action                 |
| ------------- | ---------------------- |
| `Cmd+P`       | Quick open file        |
| `Cmd+Shift+P` | Command palette        |
| `Cmd+B`       | Toggle sidebar         |
| `Cmd+J`       | Toggle panel           |
| `Cmd+\``      | Toggle terminal        |
| `Cmd+/`       | Toggle comment         |
| `F5`          | Start debugging        |
| `Cmd+Shift+F` | Search in files        |
| `Cmd+D`       | Select next occurrence |
| `Alt+Click`   | Add cursor             |
| `Cmd+Shift+K` | Delete line            |

## Workspace Configuration

### .vscode/settings.json

```json
{
  "python.pythonPath": "${workspaceFolder}/venv/bin/python",
  "python.testing.pytestEnabled": true,
  "python.testing.pytestArgs": ["tests"],
  "python.linting.enabled": true,
  "python.linting.ruffEnabled": true
}
```

### .vscode/launch.json

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Python: FastAPI",
      "type": "python",
      "request": "launch",
      "module": "uvicorn",
      "args": ["main:app", "--reload"],
      "cwd": "${workspaceFolder}/agents/legacy"
    },
    {
      "name": "Python: Current File",
      "type": "python",
      "request": "launch",
      "program": "${file}"
    }
  ]
}
```

### .vscode/tasks.json

```json
{
  "version": "2.0.0",
  "tasks": [
    {
      "label": "Start Legacy Agent",
      "type": "shell",
      "command": "cd agents/legacy && source ../../venv/bin/activate && python main.py",
      "group": "build"
    },
    {
      "label": "Start UI (Dev)",
      "type": "shell",
      "command": "cd ui && pnpm dev",
      "group": "build"
    },
    {
      "label": "Run Tests",
      "type": "shell",
      "command": "pytest -v",
      "group": {
        "kind": "test",
        "isDefault": true
      }
    }
  ]
}
```

## GitHub Copilot

### Copilot Features

- **Code Completion** - Inline suggestions
- **Copilot Chat** - Ask questions about code
- **Explain Code** - Understand existing code
- **Generate Tests** - Create unit tests
- **Fix Code** - Debug and fix issues

### Copilot Chat Commands

```
/explain - Explain selected code
/fix - Fix errors in selected code
/tests - Generate tests
/doc - Add documentation
@workspace - Reference entire workspace
@terminal - Reference terminal output
```

## Debugging

### Python Debugging

1. Set breakpoints (click gutter)
2. Start debugging (F5)
3. Use Debug Console for evaluation
4. Step through code (F10, F11)
5. Inspect variables in sidebar

### Debug Configuration

```json
{
  "name": "Python: Debug Tests",
  "type": "python",
  "request": "launch",
  "module": "pytest",
  "args": ["-v", "-s", "${file}"],
  "console": "integratedTerminal"
}
```

## Multi-Root Workspaces

### OneMind.code-workspace

```json
{
  "folders": [
    {
      "name": "OneMind-OS (Main Project)",
      "path": "../onemind-os"
    },
    {
      "name": "OneMind-Codex (Knowledge Vault)",
      "path": "../OneMind-Codex"
    },
    {
      "name": "OneMind-Codex-Live (Public Docs)",
      "path": "../OneMind-Codex-Live"
    }
  ],
  "settings": {
    "files.exclude": {
      "**/node_modules": true
    }
  }
}
```

## Snippets

### python.json

```json
{
  "Agent Tool": {
    "prefix": "tool",
    "body": [
      "@Tool",
      "def ${1:tool_name}(${2:param}: str) -> str:",
      "    \"\"\"${3:Tool description.}\"\"\"",
      "    $0",
      "    return result"
    ]
  },
  "FastAPI Endpoint": {
    "prefix": "endpoint",
    "body": [
      "@app.${1|get,post,put,delete|}(\"/${2:path}\")",
      "async def ${3:handler}(${4}):",
      "    $0"
    ]
  }
}
```

## Best Practices

- ✅ Use workspace settings for project-specific config
- ✅ Configure Python interpreter per project
- ✅ Enable format on save
- ✅ Use integrated terminal
- ✅ Learn keyboard shortcuts
- ✅ Use Copilot for boilerplate
- ✅ Keep extensions minimal (for performance)

## Resources

- [VS Code Documentation](https://code.visualstudio.com/docs)
- [Keyboard Shortcuts](https://code.visualstudio.com/docs/getstarted/keybindings)
- [Extension Marketplace](https://marketplace.visualstudio.com)

---

**Related**: [GitHub](00-24%20SO%20(Sovereign%20Operations)/02%20Infrastructure%20(VPS,%20Servers%20&%20Networking)%20(The%20Body)/Tech-Stack/platforms/development/GitHub.md) | [Claude Code](./Claude-Code.md)


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
