# iMessage

> **Apple's messaging platform for iOS/macOS integration**

## Overview

iMessage provides native Apple ecosystem messaging. Integration is limited but possible via AppleScript on macOS or third-party bridges.

## Status

| Attribute   | Value                  |
| ----------- | ---------------------- |
| Status      | 🔄 Limited             |
| Category    | Communication          |
| Type        | Messaging Platform     |
| Integration | AppleScript, Shortcuts |

## Use Cases in OneMind

1. **Personal Notifications** - Direct messages to iPhone
2. **Quick Responses** - Reply from watch/phone
3. **Siri Integration** - Voice-triggered messages
4. **Family Communication** - Reach family members

## Limitations

- No official API
- macOS required for automation
- AppleScript-based (fragile)
- Cannot receive messages programmatically easily

## Integration Methods

### Method 1: AppleScript (macOS)

```python
import subprocess

def send_imessage(recipient: str, message: str) -> bool:
    """Send iMessage via AppleScript on macOS."""
    script = f'''
    tell application "Messages"
        set targetService to 1st account whose service type = iMessage
        set targetBuddy to participant "{recipient}" of targetService
        send "{message}" to targetBuddy
    end tell
    '''

    try:
        subprocess.run(
            ["osascript", "-e", script],
            check=True,
            capture_output=True
        )
        return True
    except subprocess.CalledProcessError:
        return False
```

### Method 2: Shortcuts Automation

Create a Shortcut on iOS/macOS:

1. Input: Text (message), Text (recipient)
2. Action: Send Message
3. Trigger via URL scheme or command line

```python
import subprocess

def trigger_shortcut(shortcut_name: str, input_text: str):
    """Run Apple Shortcut with input."""
    subprocess.run([
        "shortcuts", "run", shortcut_name,
        "-i", input_text
    ])
```

### Method 3: BlueBubbles Bridge

Self-hosted solution for iMessage API access:

```python
import httpx

BLUEBUBBLES_URL = "http://localhost:1234"
BLUEBUBBLES_PASSWORD = "your_password"

async def send_via_bluebubbles(
    chat_guid: str,
    message: str
):
    """Send iMessage via BlueBubbles server."""
    async with httpx.AsyncClient() as client:
        response = await client.post(
            f"{BLUEBUBBLES_URL}/api/v1/message/text",
            params={"password": BLUEBUBBLES_PASSWORD},
            json={
                "chatGuid": chat_guid,
                "message": message
            }
        )
        return response.json()
```

## Legacy AI Integration

```python
from agno import Tool
import subprocess

@Tool
async def send_imessage(
    recipient: str,
    message: str
) -> str:
    """Send iMessage to a contact.

    Args:
        recipient: Phone number or email
        message: Message to send

    Note: Only works on macOS with Messages app configured.
    """
    script = f'''
    tell application "Messages"
        set targetService to 1st account whose service type = iMessage
        set targetBuddy to participant "{recipient}" of targetService
        send "{message}" to targetBuddy
    end tell
    '''

    try:
        subprocess.run(
            ["osascript", "-e", script],
            check=True,
            capture_output=True
        )
        return f"iMessage sent to {recipient}"
    except Exception as e:
        return f"Failed to send iMessage: {e}"
```

## macOS Setup

1. **Enable Messages App** - Sign in with Apple ID
2. **Grant Permissions** - Allow automation in System Preferences
3. **Test Manually** - Send a message first via Messages app

```bash
# Test via terminal
osascript -e 'tell application "Messages" to send "Test" to participant "+1234567890"'
```

## Alternative: Use Other Platforms

For reliable messaging automation, consider:

- **Telegram** - Full API, great bot support
- **Discord** - Webhooks, bots
- **SMS via Twilio** - Reaches any phone

## Environment Variables

```bash
# If using BlueBubbles
BLUEBUBBLES_URL=http://mac-mini.local:1234
BLUEBUBBLES_PASSWORD=your_password

# Contact mappings
IMESSAGE_ZEUS=+1234567890
```

## Best Practices

1. **Use as fallback** - Primary should be Telegram/Discord
2. **macOS Mac Mini** - Dedicated server for iMessage
3. **Error handling** - AppleScript can fail silently
4. **Rate limiting** - Apple may flag automated messages

## Resources

- [BlueBubbles](https://bluebubbles.app/)
- [Apple Shortcuts](https://support.apple.com/guide/shortcuts/welcome/ios)
- [Messages AppleScript](https://developer.apple.com/library/archive/documentation/AppleScript/Conceptual/AppleScriptX/)

---

Last Updated: January 25, 2026
