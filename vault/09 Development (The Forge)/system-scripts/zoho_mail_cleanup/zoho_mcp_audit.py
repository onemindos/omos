#!/usr/bin/env python3
"""Full Zoho Mail audit via MCP."""
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

# --- Folders ---
print("=== ALL FOLDERS ===")
folders_data = tool_call("ZohoMail_getAllFolders", {
    "path_variables": {"accountId": ACCOUNT_ID}
})
folders = folders_data.get("data", [])
for f in folders:
    print(f"  [{f.get('folderId')}] {f.get('path', f.get('folderName','?'))}"
          f"  ({f.get('messageCount',0)} msgs, {f.get('unreadCount',0)} unread)")

# --- Domains ---
print("\n=== DOMAINS ===")
d = tool_call("ZohoMail_fetchAllDomains", {})
for domain in d.get("data", []):
    print(f"  {domain.get('domainName','?')}  status={domain.get('status','?')}  mx={domain.get('mxStatus','?')}")

# --- Org Users ---
print("\n=== ORG USERS ===")
users_data = tool_call("ZohoMail_fetchOrgUsersDetails", {})
for u in users_data.get("data", []):
    print(f"  {u.get('displayName','?')} <{u.get('primaryEmailAddress','?')}>  role={u.get('role','?')}  status={u.get('status','?')}")

# --- Groups ---
print("\n=== GROUPS ===")
groups = tool_call("ZohoMail_fetchAllGroups", {})
for g in groups.get("data", []):
    print(f"  {g.get('groupName','?')} <{g.get('groupMailId','?')}>")

# --- Policies ---
print("\n=== POLICIES (email filters/rules) ===")
policies = tool_call("ZohoMail_getAllPolicies", {})
pol_list = policies.get("data", [])
if pol_list:
    for p in pol_list:
        print(f"  [{p.get('policyId','?')}] {p.get('policyName','?')}  active={p.get('isActive','?')}")
        print(f"    criteria: {p.get('criteria','?')}")
        print(f"    action:   {p.get('action','?')}")
else:
    print("  None found or not accessible")
    print(f"  Raw: {str(policies)[:300]}")
