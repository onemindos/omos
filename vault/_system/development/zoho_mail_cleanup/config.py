"""
Zoho Mail Cleanup Toolkit — Config & Connection Module
======================================================
Handles IMAP connection to Zoho Mail with reusable connection context.

Setup:
  1. Enable IMAP in Zoho Mail: Settings → Mail Accounts → IMAP Access → Enable
  2. Generate App-Specific Password: Settings → Security → App Passwords → Generate
  3. Set environment variables or create .env file in this directory
"""

import imaplib
import os
import sys
from pathlib import Path
from contextlib import contextmanager

# ---------------------------------------------------------------------------
# Load .env if present (no external dependency needed)
# ---------------------------------------------------------------------------
ENV_FILE = Path(__file__).parent / ".env"

def load_dotenv():
    """Load .env file into os.environ (simple, no dependency)."""
    if ENV_FILE.exists():
        for line in ENV_FILE.read_text().splitlines():
            line = line.strip()
            if line and not line.startswith("#") and "=" in line:
                key, _, value = line.partition("=")
                key = key.strip()
                value = value.strip().strip('"').strip("'")
                os.environ.setdefault(key, value)

load_dotenv()

# ---------------------------------------------------------------------------
# Credentials
# ---------------------------------------------------------------------------
ZOHO_EMAIL    = os.environ.get("ZOHO_EMAIL", "zeus@onemindcodex.com")
ZOHO_APP_PASS = os.environ.get("ZOHO_APP_PASSWORD", "")

IMAP_HOST = "imappro.zoho.com"
IMAP_PORT = 993

# ---------------------------------------------------------------------------
# Connection helpers
# ---------------------------------------------------------------------------

def check_credentials():
    """Verify credentials are set before attempting connection."""
    if not ZOHO_APP_PASS:
        print("\n❌  ZOHO_APP_PASSWORD not set!")
        print()
        print("Setup steps:")
        print("  1. Go to https://mail.zoho.com → Settings → Mail Accounts → IMAP Access")
        print("     → Enable IMAP")
        print()
        print("  2. Go to https://accounts.zoho.com/home#security/security_pwd")
        print("     → App Passwords → Generate New Password")
        print("     → Name it 'onemind-cleanup' → Copy the password")
        print()
        print("  3. Create a .env file in this directory:")
        print(f"     {Path(__file__).parent / '.env'}")
        print()
        print('     ZOHO_EMAIL="zeus@onemindcodex.com"')
        print('     ZOHO_APP_PASSWORD="your-app-specific-password-here"')
        print()
        sys.exit(1)


def connect_imap() -> imaplib.IMAP4_SSL:
    """Create and return an authenticated IMAP connection."""
    check_credentials()
    try:
        mail = imaplib.IMAP4_SSL(IMAP_HOST, IMAP_PORT)
        mail.login(ZOHO_EMAIL, ZOHO_APP_PASS)
        return mail
    except imaplib.IMAP4.error as e:
        print(f"\n❌  IMAP login failed: {e}")
        print("    Check your email and app-specific password.")
        sys.exit(1)
    except Exception as e:
        print(f"\n❌  Connection failed: {e}")
        print(f"    Could not reach {IMAP_HOST}:{IMAP_PORT}")
        sys.exit(1)


@contextmanager
def imap_session():
    """Context manager for auto-closing IMAP sessions."""
    mail = connect_imap()
    try:
        yield mail
    finally:
        try:
            mail.logout()
        except Exception:
            pass


def fmt_count(n: int) -> str:
    """Format number with commas."""
    return f"{n:,}"


def confirm(prompt: str, default_yes: bool = False) -> bool:
    """Ask user for Y/N confirmation."""
    suffix = " [Y/n]: " if default_yes else " [y/N]: "
    answer = input(prompt + suffix).strip().lower()
    if not answer:
        return default_yes
    return answer in ("y", "yes")
