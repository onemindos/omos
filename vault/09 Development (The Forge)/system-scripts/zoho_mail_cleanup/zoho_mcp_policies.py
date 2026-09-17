#!/usr/bin/env python3
"""Zoho MCP - get org details and policies."""
import urllib.request, json

MCP_URL = "https://onemind-897614417.zohomcp.com/mcp/message?key=6aac3d71867d9953cea0e9c1ba41cfb7"
ACCOUNT_ID = "8698543000000009002"

def tool_call(name, args):
    body = json.dumps({
        "jsonrpc": "2.0", "id": 1, "method": "tools/call",
        "params": {"name": name, "arguments": args}
    }).encode()
    req = urllib.request.Request(
        MCP_URL, data=body,
        headers={"Content-Type": "application/json"}, method="POST"
    )
    with urllib.request.urlopen(req, timeout=20) as resp:
        r = json.loads(resp.read())
    content = r.get("result", {}).get("content", [{}])
    raw = content[0].get("text", "{}") if content else "{}"
    try:
        return json.loads(raw)
    except Exception:
        return {"raw": raw}

# Get full account data to find zoid
print("=== FULL ACCOUNT DATA ===")
accounts = tool_call("ZohoMail_getMailAccounts", {})
for acc in accounts.get("data", []):
    print(f"  accountId: {acc.get('accountId')}")
    print(f"  accountName: {acc.get('accountName')}")
    print(f"  zuid: {acc.get('zuid')}")
    print(f"  role: {acc.get('role')}")
    print(f"  primaryEmail: {acc.get('primaryEmailAddress')}")
    print(f"  emailAddresses:")
    for e in acc.get("emailAddress", []):
        print(f"    {'PRIMARY' if e.get('isPrimary') else 'ALIAS  '}: {e.get('mailId')}  alias={e.get('isAlias')}  confirmed={e.get('isConfirmed')}")
    send_details = acc.get("sendMailDetails", [])
    print(f"  sendMailDetails ({len(send_details)} entries):")
    for s in send_details:
        print(f"    [{s.get('sendMailId')}] name={s.get('displayName')}  from={s.get('userName')}")
    print()

# Try policies with different zoid guesses
# zoid is often the numeric org ID — try the account "accountId" prefix pattern
# From accountId 8698543000000009002, org part might be 8698543
zoid_guesses = ["8698543", "856177564", ACCOUNT_ID, "8698543000000009002"]
print("=== POLICIES — trying different zoid values ===")
for zoid in zoid_guesses:
    result = tool_call("ZohoMail_getAllPolicies", {
        "path_variables": {"getAllPoliciesZoid": zoid}
    })
    raw = str(result)
    if "error" not in raw.lower() and "not present" not in raw.lower() and "Mandatory" not in raw:
        print(f"  zoid={zoid} -> SUCCESS")
        print(json.dumps(result, indent=2)[:3000])
        break
    else:
        print(f"  zoid={zoid} -> {raw[:120]}")
