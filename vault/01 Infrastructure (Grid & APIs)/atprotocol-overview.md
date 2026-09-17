# AT Protocol (Authenticated Transfer Protocol)

**Type:** System Architecture - Decentralized Social Protocol  
**Category:** Web3, Decentralization, Social Networks  
**Created:** 2026-01-20  
**Status:** Production (Powers Bluesky)  
**Website:** https://atproto.com/

---

## Overview

**AT Protocol** is an open, decentralized network protocol for building social applications. Created by Bluesky (founded by Twitter's Jack Dorsey), it aims to give users control over their data, identity, and social graph while maintaining the performance and user experience of centralized platforms.

**Core Philosophy:**

> "Speech and reach should be separate layers. Speech should be permissive and distributed. Reach should be flexible and scalable."

---

## The Problem AT Protocol Solves

### Traditional Social Media Issues

**Centralized Platforms:**

- ❌ Platform lock-in (your data trapped)
- ❌ Account suspension = lose everything
- ❌ No algorithmic choice (forced timeline)
- ❌ Corporate control over reach
- ❌ Vendor lock-in for developers

**Previous Decentralization Attempts:**

- **Mastodon/ActivityPub:** Server-dependent identity (hard to migrate)
- **Peer-to-peer (P2P):** Poor mobile experience, availability issues
- **Blockchain:** Slow, expensive, complexity overhead

### AT Protocol's Solution

✅ **Account portability** - Take your identity anywhere  
✅ **Data ownership** - You control your content  
✅ **Algorithmic choice** - Choose your own feeds  
✅ **Federation** - Multiple providers, one network  
✅ **Developer freedom** - Build apps on top

---

## Core Concepts

### 1. Identity Layer - DIDs (Decentralized Identifiers)

**What is a DID?**
A permanent, cryptographically-verifiable identifier for your account.

```
Example DID: did:plc:z72i7hdynmk6r22z27h6tvur
Human handle: @alice.bsky.social
```

**Key Features:**

- **Portable** - Works across any PDS (Personal Data Server)
- **Cryptographic** - Public/private key pairs
- **Verifiable** - Proves you own your account
- **Persistent** - Never changes, even if you switch servers

**Two Types of Keys:**

1. **Signing Key** - Validates your data repository (managed by PDS)
2. **Rotation Keys** - Updates your DID document (you control)

**Why This Matters:**
If your server goes down or kicks you off, you keep your identity, followers, and content by moving to a new server.

---

### 2. Data Repositories - Your Personal Database

**What is it?**
A signed, portable collection of all your social data (posts, likes, follows, blocks, etc.).

**Think of it as:**

```
Git repository + Database records + Cryptographic signatures
```

**Structure:**

```
User Repository (signed by DID)
├── Posts (app.bsky.feed.post)
├── Likes (app.bsky.feed.like)
├── Follows (app.bsky.graph.follow)
├── Blocks (app.bsky.graph.block)
├── Profile (app.bsky.actor.profile)
└── Custom records (your app's schema)
```

**Key Properties:**

- **Signed** - Cryptographically verified as yours
- **Versioned** - Complete history like Git commits
- **Portable** - Can be exported and moved
- **Interoperable** - Other apps can read your data

---

### 3. Network Architecture - Three Core Services

#### A) Personal Data Server (PDS)

**Your home in the cloud.**

**What it does:**

- Hosts your data repository
- Manages your identity (DID)
- Handles authentication
- Distributes your content
- Serves as your outbox/inbox

**Think of it as:**
Your personal email server, but for social media.

**Options:**

- Self-host (full control)
- Use Bluesky's PDS (bsky.social)
- Use any third-party provider
- Switch providers anytime

**Example:**

```
Alice's PDS: alice.pds-provider.com
Bob's PDS: bob.selfhosted.org
```

---

#### B) Relay (Firehose Aggregator)

**Collects updates from across the network.**

**What it does:**

- Subscribes to many PDSs
- Aggregates all public data into a "firehose"
- Provides real-time event stream
- No user data stored (just indexes)

**Use Case:**
App developers subscribe to the relay to get all network activity without crawling thousands of PDSs.

**Flow:**

```
PDS 1 ──┐
PDS 2 ──┤
PDS 3 ──┼──> RELAY ──> Firehose Stream
PDS 4 ──┤              (all posts/events)
PDS 5 ──┘
```

---

#### C) App View (Application Backend)

**Provides the "Twitter-like" experience.**

**What it does:**

- Indexes data from the relay
- Calculates metrics (likes, reposts, followers)
- Powers search
- Generates feeds (algorithmic timelines)
- Handles recommendations

**Examples:**

- `bsky.app` - Bluesky's official App View
- Custom App Views for specific communities
- Specialized views (music-focused, news-focused, etc.)

**Why Separate?**
You can switch App Views without losing your data or identity.

---

### 4. Lexicon - Global Schema Language

**What is Lexicon?**
A schema language that defines the structure of records and API calls across the network.

**Think of it as:**

```
JSON Schema + OpenAPI + Type Definitions
```

**Purpose:**
Ensures all apps speak the same language.

**Example Lexicon:**

```json
{
  "lexicon": 1,
  "id": "app.bsky.feed.post",
  "defs": {
    "main": {
      "type": "record",
      "key": "tid",
      "record": {
        "type": "object",
        "required": ["text", "createdAt"],
        "properties": {
          "text": { "type": "string", "maxLength": 3000 },
          "createdAt": { "type": "string", "format": "datetime" }
        }
      }
    }
  }
}
```

**Standard Namespaces:**

- `com.atproto.*` - Core protocol functions (repos, identity, sync)
- `app.bsky.*` - Bluesky social features (posts, likes, follows)
- `your.app.*` - Your custom application schemas

---

## How It All Works Together

### Publishing a Post

```
1. USER writes post
   ↓
2. PDS creates record in user's repository
   ↓
3. PDS signs record with user's DID
   ↓
4. PDS broadcasts update to Relay
   ↓
5. Relay adds to firehose stream
   ↓
6. App Views consume firehose
   ↓
7. App Views index post, update metrics
   ↓
8. OTHER USERS see post in their feeds
```

### Following Someone

```
1. USER clicks "Follow @alice"
   ↓
2. PDS creates follow record in user's repo
   ↓
3. Signed and distributed via Relay
   ↓
4. App View indexes the relationship
   ↓
5. Alice appears in user's following list
   ↓
6. Alice's posts now in user's timeline
```

---

## Key Features

### 1. Account Portability

**The Problem:**
On Twitter/X, if your account is suspended, you lose:

- Your followers
- Your content
- Your username
- Your identity

**AT Protocol Solution:**

```
Old Server Dies/Bans You
    ↓
Update DID document (rotation key)
    ↓
Point to New PDS
    ↓
Upload data backup
    ↓
All followers/content intact!
```

**How to Migrate:**

1. Export your repository (like Git clone)
2. Sign up at new PDS
3. Update your DID to point to new PDS
4. Import your repository
5. Done! Followers automatically see new location

---

### 2. Algorithmic Choice

**Traditional:** Platform decides what you see (engagement bait, ads)

**AT Protocol:** You choose your algorithm

**Feed Generators:**
Third-party services that create custom timelines.

**Examples:**

- Chronological feed
- "What's Hot" (trending)
- Topic-based (just tech news)
- Community curated (trusted moderators)
- AI-powered (personalized interests)
- Quiet mode (friends only)

**How It Works:**

```
App View provides data
    ↓
Feed Generator algorithms
    ↓
User subscribes to feeds
    ↓
Custom timeline experience
```

---

### 3. Speech vs Reach

**Two-Layer Model:**

**Layer 1: Speech (Permissionless)**

- Anyone can post
- Data stored in repositories
- Distributed via federation
- Like "anyone can make a website"

**Layer 2: Reach (Algorithmic)**

- App Views decide what to amplify
- Feed generators create timelines
- Labelers add context/warnings
- Like "search engines decide ranking"

**Why Split?**

- Preserves free expression (can't be silenced)
- Allows community curation (harmful content de-amplified)
- User choice (switch to different App View)

---

### 4. Federated Architecture

**Not Peer-to-Peer:**

- Servers host your data (not your phone)
- Always available (no need for devices online)
- Fast performance (server-side processing)

**Not Fully Centralized:**

- Multiple PDS providers
- Competitive App Views
- Open protocol (anyone can participate)

**Best of Both Worlds:**

```
Centralized UX + Decentralized Control
```

---

## Technical Specifications

### Core Specs

**1. Authenticated Transfer Protocol (ATP)**

- Core protocol for data sync
- Repository format
- Cryptographic verification

**2. DIDs & Handles**

- `did:plc` - Placeholder DID (Bluesky's method)
- `did:web` - Web-based DID (DNS)
- Handle resolution (username → DID)

**3. Repository & Data Model**

- MST (Merkle Search Tree) structure
- CAR (Content Addressable aRchive) format
- Record types and schemas

**4. Lexicon**

- Schema definition language
- Validation rules
- API contracts

**5. HTTP API (XRPC)**

- RPC over HTTP
- WebSocket event streams
- Query/procedure calls

---

## Building on AT Protocol

### SDK Options

**Official SDKs:**

- `@atproto/api` (JavaScript/TypeScript)
- Python client library
- Go libraries

**What You Can Build:**

**1. Custom Clients**

```javascript
import { BskyAgent } from "@atproto/api";

const agent = new BskyAgent({ service: "https://bsky.social" });
await agent.login({ identifier: "alice.bsky.social", password: "***" });
await agent.post({ text: "Hello ATProto!" });
```

**2. Feed Generators**
Custom algorithmic timelines.

**3. Labelers**
Content moderation/categorization services.

**4. App Views**
Entirely new social apps using the same data layer.

**5. Bots & Automation**
Scheduled posting, analytics, moderation tools.

---

## Comparisons

### vs ActivityPub (Mastodon)

| Feature             | AT Protocol           | ActivityPub              |
| ------------------- | --------------------- | ------------------------ |
| **Identity**        | Portable DIDs         | Server-bound             |
| **Migration**       | Easy (keep followers) | Hard (lose followers)    |
| **Data Format**     | Signed repositories   | JSON activities          |
| **Algorithm**       | User choice           | Server timeline          |
| **Scale**           | Relay aggregation     | Server-to-server         |
| **App flexibility** | High (decouple view)  | Medium (server controls) |

### vs Nostr

| Feature               | AT Protocol        | Nostr                               |
| --------------------- | ------------------ | ----------------------------------- |
| **Architecture**      | Federated servers  | Relays (P2P-like)                   |
| **Identity**          | DID (recoverable)  | Private key (lose key = locked out) |
| **Data availability** | Always (servers)   | Depends on relays                   |
| **Performance**       | Fast (server-side) | Variable (client-dependent)         |
| **Moderation**        | App View layer     | Client-side filtering               |

### vs Web3 Social (Lens, Farcaster)

| Feature        | AT Protocol   | Lens/Farcaster        |
| -------------- | ------------- | --------------------- |
| **Blockchain** | No blockchain | Yes (Ethereum L2)     |
| **Speed**      | Instant posts | Block time delays     |
| **Cost**       | Free posting  | Gas fees              |
| **Complexity** | Moderate      | High (crypto wallets) |
| **UX**         | Web2-like     | Web3 friction         |

---

## Use Cases & Applications

### 1. Bluesky (Primary App)

**What:** Twitter-like microblogging  
**URL:** https://bsky.app  
**Status:** 20M+ users (as of Jan 2026)

**Features:**

- 300-character posts
- Images, videos, links
- Likes, reposts, quotes
- Custom feeds
- Moderation tools

---

### 2. Potential Applications

**Social Networks:**

- Long-form blogging (like Medium)
- Photo sharing (like Instagram)
- Video platform (like YouTube)
- Professional networking (like LinkedIn)

**Community Tools:**

- Forums/discussion boards
- Event coordination
- Group chat integration
- Knowledge bases

**Content Distribution:**

- Podcasting
- Newsletter distribution
- Streaming integration
- Collaborative writing

**Developer Tools:**

- Analytics dashboards
- Moderation interfaces
- Bot frameworks
- Data export/backup tools

---

## Security & Privacy

### Cryptographic Security

**Data Integrity:**

- All records signed with DID keys
- Merkle tree verification
- Immutable history (like Git)

**Identity Security:**

- Rotation keys for account recovery
- Multi-factor authentication support
- Key compromise detection

### Privacy Features

**User Controls:**

- Public posts (default)
- Block/mute functionality
- Account reporting
- Content warnings

**Data Protection:**

- HTTPS/TLS for transport
- Optional end-to-end encryption (future)
- Selective data sharing

**Moderation:**

- User-controlled blocklists
- Shared moderation lists
- Labeler services (third-party)
- App View filtering

---

## Challenges & Limitations

### Current Challenges

**1. PDS Hosting Complexity**

- Self-hosting requires technical knowledge
- Most users rely on Bluesky's server
- Need more user-friendly PDS options

**2. Relay Centralization Risk**

- Currently few relays (Bluesky runs main one)
- High bandwidth/storage requirements
- Need incentive models for relay operators

**3. Discoverability**

- Finding users across PDSs
- Search across federated data
- Content recommendation cold start

**4. Moderation at Scale**

- Balancing speech vs safety
- Coordinating across App Views
- Handling illegal content

### Future Improvements

- [ ] End-to-end encryption for DMs
- [ ] Payments/micropayments integration
- [ ] Video/livestreaming support
- [ ] Mobile-optimized PDS hosting
- [ ] Better migration UX
- [ ] Cross-protocol bridges (ActivityPub interop)

---

## Getting Started

### As a User

**1. Join Bluesky:**

```
Visit: https://bsky.app
Sign up for account
Start posting!
```

**2. Explore Feeds:**

- Discover custom feeds in "Feeds" tab
- Subscribe to algorithmic choices
- Create your own curation

**3. Understand Your Identity:**

- Your handle (human-readable)
- Your DID (permanent ID)
- Your PDS (where data lives)

### As a Developer

**1. Set Up SDK:**

```bash
npm install @atproto/api
```

**2. Quick Start:**

```javascript
import { BskyAgent } from "@atproto/api";

const agent = new BskyAgent({
  service: "https://bsky.social",
});

await agent.login({
  identifier: "yourhandle.bsky.social",
  password: "your-password",
});

// Post
await agent.post({
  text: "Hello from AT Protocol!",
  createdAt: new Date().toISOString(),
});

// Get timeline
const timeline = await agent.getTimeline({ limit: 10 });
```

**3. Read the Docs:**
https://atproto.com/guides/applications

---

## Resources

### Official Documentation

- **Main site:** https://atproto.com/
- **Protocol specs:** https://atproto.com/specs
- **Developer guides:** https://atproto.com/guides
- **Lexicon explorer:** https://atproto.com/lexicons

### Community

- **Bluesky on Bluesky:** @atproto.com
- **GitHub:** https://github.com/bluesky-social
- **Discord:** (Check official site for link)

### Articles & Videos

- [ATProto for Distributed Systems Engineers](https://atproto.com/articles/atproto-for-distsys-engineers)
- [ATProto Ethos](https://atproto.com/articles/atproto-ethos)
- [Web Without Walls - React Universe 2024](https://www.youtube.com/watch?v=F1sJW6nTP6E)

---

## Why This Matters for Legacy Codex

### Potential Integration Points

**1. Identity Layer**

- Use DIDs for user identity
- Portable reputation/credentials
- Cross-platform proof of ownership

**2. Data Portability**

- Export knowledge graph as AT repo
- Sync notes/tasks across devices
- Backup to personal PDS

**3. Social Features**

- Share knowledge publicly
- Collaborative note-taking
- Discover related content

**4. Decentralization Philosophy**

- Aligns with data sovereignty goals
- User-controlled infrastructure
- Open protocol principles

**Future Exploration:**

- Build Lexicon for knowledge graphs
- PDS integration for backup
- Feed generator for curated learning

---

## Key Takeaways

✅ **Portable Identity** - Your account is truly yours  
✅ **Data Ownership** - Export and move anytime  
✅ **Algorithmic Choice** - Pick your timeline algorithm  
✅ **Federation** - Multiple providers, one network  
✅ **Developer Freedom** - Build apps on open protocol  
✅ **Practical Decentralization** - Good UX + user control

**Philosophy:**

> "AT Protocol proves decentralization doesn't have to sacrifice user experience."

---

**Tags:** `#decentralization` `#social-protocol` `#web3` `#identity` `#federation` `#bluesky`

**Related Docs:**

- `reference-decentralization-landscape.md` (Coming soon)
- `system-identity-management.md` (Coming soon)
- `process-data-portability-strategy.md` (Coming soon)
