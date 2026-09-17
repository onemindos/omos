"""
Legacy Consciousness - Deep Link Generator

Generates deep links to open content in native apps on iOS/macOS/Android.
Used by Canvas and Notifier to let users view content directly.

Supported Apps:
- Email: Gmail, Apple Mail, Outlook
- GitHub: GitHub app
- Tasks: Todoist
- Documents: Google Docs, Office 365
- Calendar: Apple Calendar, Google Calendar
- CRM: Zoho CRM
- Chat: Slack, Discord
- Web: Safari, Chrome

Usage:
    from backend.consciousness.deep_links import DeepLinks

    # Generate email deep link
    link = DeepLinks.email("gmail", message_id="123abc")
    # Returns: "googlegmail:///co?messageId=123abc"

    # Generate GitHub PR link
    link = DeepLinks.github("pr", "owner", "repo", 42)
    # Returns: "github://pull-request/owner/repo/42"
"""

from enum import Enum
from typing import Optional
from urllib.parse import quote, urlencode


class Platform(str, Enum):
    """Target platform for deep links."""
    IOS = "ios"
    MACOS = "macos"
    ANDROID = "android"
    WEB = "web"


class DeepLinks:
    """
    Static methods to generate deep links for various apps.

    All methods return a deep link string that can be used to open
    content directly in the native app.
    """

    # =========================================================================
    # Email
    # =========================================================================

    @staticmethod
    def email(
        provider: str,
        message_id: str = "",
        thread_id: str = "",
        to: str = "",
        subject: str = "",
        body: str = "",
    ) -> str:
        """
        Generate email deep link.

        Args:
            provider: gmail, outlook, apple, or generic
            message_id: Specific message ID
            thread_id: Thread/conversation ID
            to: Compose new email to this address
            subject: Email subject (for compose)
            body: Email body (for compose)
        """
        if provider == "gmail":
            if message_id:
                return f"googlegmail:///co?messageId={message_id}"
            elif thread_id:
                return f"googlegmail:///th?threadId={thread_id}"
            elif to:
                params = {"to": to}
                if subject:
                    params["subject"] = subject
                if body:
                    params["body"] = body
                return f"googlegmail:///co?{urlencode(params)}"
            return "googlegmail:///"

        elif provider == "outlook":
            if message_id:
                return f"ms-outlook://emails/{message_id}"
            elif to:
                params = {"to": to}
                if subject:
                    params["subject"] = subject
                return f"ms-outlook://compose?{urlencode(params)}"
            return "ms-outlook://emails"

        elif provider == "apple":
            if message_id:
                return f"message://{message_id}"
            elif to:
                params = {}
                if subject:
                    params["subject"] = subject
                if body:
                    params["body"] = body
                query = f"?{urlencode(params)}" if params else ""
                return f"mailto:{to}{query}"
            return "message://"

        else:
            # Generic mailto
            params = {}
            if subject:
                params["subject"] = subject
            if body:
                params["body"] = body
            query = f"?{urlencode(params)}" if params else ""
            return f"mailto:{to or ''}{query}"

    # =========================================================================
    # GitHub
    # =========================================================================

    @staticmethod
    def github(
        type: str,
        owner: str,
        repo: str,
        number: int = 0,
        file_path: str = "",
        line: int = 0,
    ) -> str:
        """
        Generate GitHub deep link.

        Args:
            type: pr, issue, repo, file, commit
            owner: Repository owner
            repo: Repository name
            number: PR/Issue number
            file_path: File path (for file links)
            line: Line number (for file links)
        """
        base_web = f"https://github.com/{owner}/{repo}"

        if type == "pr":
            # GitHub app deep link (iOS)
            app_link = f"github://pull-request/{owner}/{repo}/{number}"
            web_link = f"{base_web}/pull/{number}"
            return app_link  # Could return web_link as fallback

        elif type == "issue":
            app_link = f"github://issue/{owner}/{repo}/{number}"
            web_link = f"{base_web}/issues/{number}"
            return app_link

        elif type == "repo":
            app_link = f"github://repo/{owner}/{repo}"
            return app_link

        elif type == "file":
            # Files don't have good app deep links, use web
            line_anchor = f"#L{line}" if line else ""
            return f"{base_web}/blob/main/{file_path}{line_anchor}"

        elif type == "commit":
            return f"github://commit/{owner}/{repo}/{file_path}"  # file_path = commit SHA

        return f"github://repo/{owner}/{repo}"

    @staticmethod
    def github_web(
        type: str,
        owner: str,
        repo: str,
        number: int = 0,
    ) -> str:
        """Generate GitHub web URL (fallback)."""
        base = f"https://github.com/{owner}/{repo}"
        if type == "pr":
            return f"{base}/pull/{number}"
        elif type == "issue":
            return f"{base}/issues/{number}"
        return base

    # =========================================================================
    # Todoist
    # =========================================================================

    @staticmethod
    def todoist_task(task_id: str) -> str:
        """Generate Todoist task deep link."""
        return f"todoist://task?id={task_id}"

    @staticmethod
    def todoist_project(project_id: str) -> str:
        """Generate Todoist project deep link."""
        return f"todoist://project?id={project_id}"

    @staticmethod
    def todoist_filter(filter_name: str) -> str:
        """Generate Todoist filter deep link."""
        return f"todoist://filter?name={quote(filter_name)}"

    @staticmethod
    def todoist_add_task(
        content: str = "",
        due: str = "",
        priority: int = 1,
    ) -> str:
        """Generate Todoist add task deep link."""
        params = {}
        if content:
            params["content"] = content
        if due:
            params["due"] = due
        if priority > 1:
            params["priority"] = str(priority)
        query = f"?{urlencode(params)}" if params else ""
        return f"todoist://addtask{query}"

    # =========================================================================
    # Google Docs / Drive
    # =========================================================================

    @staticmethod
    def google_doc(doc_id: str) -> str:
        """Generate Google Docs deep link."""
        return f"googledocs://docs.google.com/document/d/{doc_id}"

    @staticmethod
    def google_sheet(sheet_id: str) -> str:
        """Generate Google Sheets deep link."""
        return f"googlesheets://docs.google.com/spreadsheets/d/{sheet_id}"

    @staticmethod
    def google_slides(presentation_id: str) -> str:
        """Generate Google Slides deep link."""
        return f"googleslides://docs.google.com/presentation/d/{presentation_id}"

    @staticmethod
    def google_drive(file_id: str = "") -> str:
        """Generate Google Drive deep link."""
        if file_id:
            return f"googledrive://drive.google.com/file/d/{file_id}"
        return "googledrive://"

    # =========================================================================
    # Calendar
    # =========================================================================

    @staticmethod
    def calendar_apple(event_id: str = "") -> str:
        """Generate Apple Calendar deep link."""
        if event_id:
            return f"calshow:{event_id}"
        return "calshow:"

    @staticmethod
    def calendar_google(event_id: str = "", calendar_id: str = "") -> str:
        """Generate Google Calendar deep link."""
        if event_id:
            return f"googlecalendar://event?eventId={event_id}"
        return "googlecalendar://"

    @staticmethod
    def calendar_add_event(
        title: str,
        start: str = "",
        end: str = "",
        location: str = "",
    ) -> str:
        """Generate calendar add event link (Apple Calendar)."""
        params = {"title": title}
        if start:
            params["startDate"] = start
        if end:
            params["endDate"] = end
        if location:
            params["location"] = location
        return f"calshow://new?{urlencode(params)}"

    # =========================================================================
    # Zoho CRM
    # =========================================================================

    @staticmethod
    def zoho_crm(module: str, record_id: str) -> str:
        """Generate Zoho CRM deep link."""
        # Zoho CRM uses web links primarily
        return f"https://crm.zoho.com/crm/tab/{module}/{record_id}"

    @staticmethod
    def zoho_contact(contact_id: str) -> str:
        """Generate Zoho CRM contact link."""
        return DeepLinks.zoho_crm("Contacts", contact_id)

    @staticmethod
    def zoho_deal(deal_id: str) -> str:
        """Generate Zoho CRM deal link."""
        return DeepLinks.zoho_crm("Deals", deal_id)

    # =========================================================================
    # Slack
    # =========================================================================

    @staticmethod
    def slack_channel(
        team_id: str,
        channel_id: str,
    ) -> str:
        """Generate Slack channel deep link."""
        return f"slack://channel?team={team_id}&id={channel_id}"

    @staticmethod
    def slack_message(
        team_id: str,
        channel_id: str,
        message_ts: str,
    ) -> str:
        """Generate Slack message deep link."""
        return f"slack://channel?team={team_id}&id={channel_id}&message={message_ts}"

    @staticmethod
    def slack_dm(
        team_id: str,
        user_id: str,
    ) -> str:
        """Generate Slack DM deep link."""
        return f"slack://user?team={team_id}&id={user_id}"

    # =========================================================================
    # Discord
    # =========================================================================

    @staticmethod
    def discord_channel(
        guild_id: str,
        channel_id: str,
    ) -> str:
        """Generate Discord channel deep link."""
        return f"discord://discord.com/channels/{guild_id}/{channel_id}"

    @staticmethod
    def discord_message(
        guild_id: str,
        channel_id: str,
        message_id: str,
    ) -> str:
        """Generate Discord message deep link."""
        return f"discord://discord.com/channels/{guild_id}/{channel_id}/{message_id}"

    # =========================================================================
    # Home Assistant
    # =========================================================================

    @staticmethod
    def home_assistant(
        entity_id: str = "",
        path: str = "",
    ) -> str:
        """Generate Home Assistant deep link."""
        if entity_id:
            return f"homeassistant://navigate/lovelace/default_view?entity_id={entity_id}"
        elif path:
            return f"homeassistant://navigate{path}"
        return "homeassistant://"

    # =========================================================================
    # Web / Browser
    # =========================================================================

    @staticmethod
    def safari(url: str) -> str:
        """Generate Safari deep link (iOS/macOS)."""
        return f"x-safari-{url}"

    @staticmethod
    def chrome(url: str) -> str:
        """Generate Chrome deep link."""
        return f"googlechrome://{url.replace('https://', '').replace('http://', '')}"

    @staticmethod
    def firefox(url: str) -> str:
        """Generate Firefox deep link."""
        return f"firefox://open-url?url={quote(url)}"

    # =========================================================================
    # Universal Links (OneMind App)
    # =========================================================================

    @staticmethod
    def onemind(
        path: str,
        **params,
    ) -> str:
        """
        Generate OneMind app deep link.

        For when we have a native OneMind app.
        """
        query = f"?{urlencode(params)}" if params else ""
        return f"onemind://{path}{query}"

    @staticmethod
    def onemind_canvas(item_id: str) -> str:
        """Open specific canvas item in OneMind app."""
        return f"onemind://canvas/{item_id}"

    @staticmethod
    def onemind_chat(session_id: str = "") -> str:
        """Open chat in OneMind app."""
        if session_id:
            return f"onemind://chat/{session_id}"
        return "onemind://chat"

    @staticmethod
    def onemind_input(request_id: str) -> str:
        """Open input request dialog in OneMind app."""
        return f"onemind://input/{request_id}"

    # =========================================================================
    # Utility
    # =========================================================================

    @staticmethod
    def with_fallback(app_link: str, web_link: str) -> dict:
        """Return both app and web links for fallback handling."""
        return {
            "app": app_link,
            "web": web_link,
        }

    @staticmethod
    def detect_platform_from_user_agent(user_agent: str) -> Platform:
        """Detect platform from User-Agent string."""
        ua = user_agent.lower()
        if "iphone" in ua or "ipad" in ua:
            return Platform.IOS
        elif "macintosh" in ua:
            return Platform.MACOS
        elif "android" in ua:
            return Platform.ANDROID
        else:
            return Platform.WEB
