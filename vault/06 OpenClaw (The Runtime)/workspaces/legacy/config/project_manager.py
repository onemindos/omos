"""
Project Manager - File/project organization for AI workflows
Extracted from ADA V2 for OneMind OS

Manages project workspaces with:
- CAD artifact storage
- Chat history logging
- Project context gathering for AI
- Multi-project switching

This is useful for organizing outputs from fabrication workflows.
"""

import os
import json
import shutil
import time
from pathlib import Path
from typing import List, Dict, Optional, Tuple


class ProjectManager:
    """
    Manages project workspaces for AI-assisted workflows.

    Structure:
        workspace_root/
        └── projects/
            ├── temp/          # Default temporary project
            │   ├── cad/       # Generated CAD files
            │   ├── browser/   # Browser screenshots/data
            │   └── chat_history.jsonl
            ├── my_project/
            │   ├── cad/
            │   ├── browser/
            │   └── chat_history.jsonl
            └── ...
    """

    def __init__(self, workspace_root: str):
        """
        Initialize ProjectManager.

        Args:
            workspace_root: Root directory for all projects
        """
        self.workspace_root = Path(workspace_root)
        self.projects_dir = self.workspace_root / "projects"
        self.current_project = "temp"

        # Ensure projects root exists
        if not self.projects_dir.exists():
            self.projects_dir.mkdir(parents=True)

        # Clear temp project on startup
        temp_path = self.projects_dir / "temp"
        if temp_path.exists():
            print("[ProjectManager] Clearing temp project...")
            shutil.rmtree(temp_path)

        # Create fresh temp project
        self.create_project("temp")

    def create_project(self, name: str) -> Tuple[bool, str]:
        """
        Create a new project directory with subfolders.

        Args:
            name: Project name

        Returns:
            Tuple of (success, message)
        """
        # Sanitize name
        safe_name = "".join(
            [c for c in name if c.isalnum() or c in (" ", "-", "_")]
        ).strip()
        project_path = self.projects_dir / safe_name

        if not project_path.exists():
            project_path.mkdir()
            (project_path / "cad").mkdir()
            (project_path / "browser").mkdir()
            print(f"[ProjectManager] Created project: {safe_name}")
            return True, f"Project '{safe_name}' created."
        return False, f"Project '{safe_name}' already exists."

    def switch_project(self, name: str) -> Tuple[bool, str]:
        """Switch the active project context."""
        safe_name = "".join(
            [c for c in name if c.isalnum() or c in (" ", "-", "_")]
        ).strip()
        project_path = self.projects_dir / safe_name

        if project_path.exists():
            self.current_project = safe_name
            print(f"[ProjectManager] Switched to project: {safe_name}")
            return True, f"Switched to project '{safe_name}'."
        return False, f"Project '{safe_name}' does not exist."

    def list_projects(self) -> List[str]:
        """Return list of available projects."""
        return [d.name for d in self.projects_dir.iterdir() if d.is_dir()]

    def get_current_project_path(self) -> Path:
        """Get path to current project directory."""
        return self.projects_dir / self.current_project

    def log_chat(self, sender: str, text: str):
        """Append a chat message to the current project's history."""
        log_file = self.get_current_project_path() / "chat_history.jsonl"
        entry = {"timestamp": time.time(), "sender": sender, "text": text}
        with open(log_file, "a", encoding="utf-8") as f:
            f.write(json.dumps(entry) + "\n")

    def save_cad_artifact(self, source_path: str, prompt: str) -> Optional[str]:
        """
        Copy a generated CAD file to the project's 'cad' folder.

        Args:
            source_path: Path to source STL file
            prompt: Description used to generate the file

        Returns:
            Destination path or None if failed
        """
        if not os.path.exists(source_path):
            print(f"[ProjectManager] Source file not found: {source_path}")
            return None

        # Create filename based on timestamp and prompt
        timestamp = int(time.time())
        safe_prompt = (
            "".join([c for c in prompt if c.isalnum() or c in (" ", "-", "_")])[:30]
            .strip()
            .replace(" ", "_")
        )
        filename = f"{timestamp}_{safe_prompt}.stl"

        dest_path = self.get_current_project_path() / "cad" / filename

        try:
            shutil.copy2(source_path, dest_path)
            print(f"[ProjectManager] Saved CAD artifact to: {dest_path}")
            return str(dest_path)
        except Exception as e:
            print(f"[ProjectManager] Failed to save artifact: {e}")
            return None

    def get_project_context(self, max_file_size: int = 10000) -> str:
        """
        Gather context about the current project for the AI.
        Lists all files and reads text file contents.

        Args:
            max_file_size: Maximum file size to read (bytes)

        Returns:
            Formatted string with project context
        """
        project_path = self.get_current_project_path()
        if not project_path.exists():
            return f"Project '{self.current_project}' does not exist."

        context_lines = [f"=== Project Context: '{self.current_project}' ==="]
        context_lines.append(f"Project directory: {project_path}")
        context_lines.append("")

        # List all files recursively
        all_files = []
        for root, dirs, files in os.walk(project_path):
            for f in files:
                rel_path = os.path.relpath(os.path.join(root, f), project_path)
                all_files.append(rel_path)

        if not all_files:
            context_lines.append("(No files in project yet)")
        else:
            context_lines.append(f"Files ({len(all_files)} total):")
            for f in all_files:
                context_lines.append(f"  - {f}")

        context_lines.append("")

        # Read text files (skip binary and large files)
        text_extensions = {
            ".txt",
            ".py",
            ".js",
            ".jsx",
            ".ts",
            ".tsx",
            ".json",
            ".md",
            ".html",
            ".css",
            ".jsonl",
        }
        for rel_path in all_files:
            ext = os.path.splitext(rel_path)[1].lower()
            if ext not in text_extensions:
                continue

            full_path = project_path / rel_path
            try:
                file_size = full_path.stat().st_size
                if file_size > max_file_size:
                    context_lines.append(
                        f"--- {rel_path} (too large: {file_size} bytes, skipped) ---"
                    )
                    continue

                with open(full_path, "r", encoding="utf-8", errors="ignore") as f:
                    content = f.read()
                context_lines.append(f"--- {rel_path} ---")
                context_lines.append(content)
                context_lines.append("")
            except Exception as e:
                context_lines.append(f"--- {rel_path} (error reading: {e}) ---")

        return "\n".join(context_lines)

    def get_recent_chat_history(self, limit: int = 10) -> List[Dict]:
        """Return the last N chat messages from history."""
        log_file = self.get_current_project_path() / "chat_history.jsonl"
        if not log_file.exists():
            return []

        try:
            with open(log_file, "r", encoding="utf-8") as f:
                lines = f.readlines()

            history = []
            for line in lines[-limit:]:
                try:
                    entry = json.loads(line)
                    history.append(entry)
                except json.JSONDecodeError:
                    continue
            return history
        except Exception as e:
            print(f"[ProjectManager] Failed to read chat history: {e}")
            return []


# =============================================================================
# Agno-compatible Tool Wrappers
# =============================================================================

# Global instance
_project_manager: Optional[ProjectManager] = None


def _get_manager() -> ProjectManager:
    global _project_manager
    if _project_manager is None:
        # Default workspace root - should be configured
        workspace = os.environ.get("ONEMIND_WORKSPACE", "/tmp/onemind_workspace")
        _project_manager = ProjectManager(workspace)
    return _project_manager


def project_create(name: str) -> str:
    """
    Create a new project workspace.

    Args:
        name: Name for the new project

    Returns:
        JSON status message
    """
    manager = _get_manager()
    success, message = manager.create_project(name)
    return json.dumps({"status": "success" if success else "error", "message": message})


def project_switch(name: str) -> str:
    """
    Switch to a different project.

    Args:
        name: Name of the project to switch to

    Returns:
        JSON status message
    """
    manager = _get_manager()
    success, message = manager.switch_project(name)
    return json.dumps({"status": "success" if success else "error", "message": message})


def project_list() -> str:
    """
    List all available projects.

    Returns:
        JSON list of project names
    """
    manager = _get_manager()
    projects = manager.list_projects()
    return json.dumps({"projects": projects, "current": manager.current_project})


def project_context() -> str:
    """
    Get context about the current project for AI consumption.

    Returns:
        Text summary of project files and contents
    """
    manager = _get_manager()
    return manager.get_project_context()


def get_project_tools():
    """Get list of project management tools for Agno registration."""
    try:
        from agno.tools import tool

        return [
            tool(project_create),
            tool(project_switch),
            tool(project_list),
            tool(project_context),
        ]
    except ImportError:
        return [
            project_create,
            project_switch,
            project_list,
            project_context,
        ]


# =============================================================================
# Test
# =============================================================================

if __name__ == "__main__":
    # Test the project manager
    manager = ProjectManager("/tmp/test_workspace")
    print(f"Projects: {manager.list_projects()}")
    print(f"Current: {manager.current_project}")

    manager.create_project("test_project")
    manager.switch_project("test_project")
    manager.log_chat("user", "Hello, Legacy!")
    manager.log_chat("assistant", "Hello! How can I help?")

    print(manager.get_project_context())
