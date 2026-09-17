# Solarpunk Theme - OneMindOS

A nature-inspired, optimistic color system reflecting regenerative technology and ecological harmony.

---

## Primary Palette

| Name | Hex | RGB | Use Case |
|------|-----|-----|----------|
| **Solar Gold** | `#F4A825` | 244, 168, 37 | Primary actions, energy, achievements |
| **Moss Green** | `#4A7C59` | 74, 124, 89 | Growth, habits, health metrics |
| **Sky Teal** | `#2D8A8A` | 45, 138, 138 | Water, flow states, calm actions |
| **Terracotta** | `#C1694F` | 193, 105, 79 | Earth, grounding, warnings |
| **Cream** | `#F5F1E6` | 245, 241, 230 | Background (light mode) |
| **Deep Loam** | `#1E2420` | 30, 36, 32 | Background (dark mode) |

---

## Secondary/Accent

| Name | Hex | RGB | Use Case |
|------|-----|-----|----------|
| **Dawn Pink** | `#E8B4B8` | 232, 180, 184 | Self-care, wellness, soft alerts |
| **Lichen** | `#8FBC8F` | 143, 188, 143 | Success states, completed items |
| **Amber Glow** | `#FFD93D` | 255, 217, 61 | Highlights, notifications |
| **Storm Cloud** | `#5D6D7E` | 93, 109, 126 | Secondary text, inactive states |
| **Copper** | `#B87333` | 184, 115, 51 | Premium/special features |

---

## Path-Specific Colors

| Path | Name | Hex | Symbolism |
|------|------|-----|-----------|
| **HP (Self)** | Forest Green | `#228B22` | Personal growth, inner garden |
| **LE (Kin)** | Warm Honey | `#EB9E3A` | Family warmth, heritage |
| **GE (Work)** | Deep Teal | `#1A5F5F` | Abundance, regenerative wealth |
| **Commons** | Sage | `#9CAF88` | Community, shared resources |

---

## Mode Configurations

### Light Mode
```
Background:     #F5F1E6  (Cream)
Surface:        #FFFFFF  (White)
Primary:        #4A7C59  (Moss Green)
Secondary:      #2D8A8A  (Sky Teal)
Accent:         #F4A825  (Solar Gold)
Text Primary:   #1E2420  (Deep Loam)
Text Secondary: #5D6D7E  (Storm Cloud)
Success:        #8FBC8F  (Lichen)
Warning:        #C1694F  (Terracotta)
Error:          #B85450  (Rust Red)
```

### Dark Mode
```
Background:     #1E2420  (Deep Loam)
Surface:        #2A3230  (Dark Moss)
Primary:        #8FBC8F  (Lichen)
Secondary:      #2D8A8A  (Sky Teal)
Accent:         #FFD93D  (Amber Glow)
Text Primary:   #F5F1E6  (Cream)
Text Secondary: #9CAF88  (Sage)
Success:        #8FBC8F  (Lichen)
Warning:        #EB9E3A  (Warm Honey)
Error:          #E07A5F  (Coral)
```

---

## Gradients

| Name | Start | End | Use |
|------|-------|-----|-----|
| **Sunrise** | `#F4A825` | `#E8B4B8` | Morning routines, energy states |
| **Forest Floor** | `#4A7C59` | `#1E2420` | Depth, navigation headers |
| **Golden Hour** | `#FFD93D` | `#C1694F` | Achievements, streak celebrations |
| **Ocean Breath** | `#2D8A8A` | `#9CAF88` | Calm states, meditation |
| **Dawn** | `#E8B4B8` | `#F5F1E6` | Onboarding, welcome screens |
| **Dusk** | `#5D6D7E` | `#1E2420` | Evening modes, wind-down |

---

## Semantic Colors

| Purpose | Light Mode | Dark Mode |
|---------|------------|-----------|
| **Growth/Progress** | `#4A7C59` | `#8FBC8F` |
| **Energy/Action** | `#F4A825` | `#FFD93D` |
| **Calm/Rest** | `#2D8A8A` | `#2D8A8A` |
| **Alert/Attention** | `#C1694F` | `#EB9E3A` |
| **Community** | `#9CAF88` | `#9CAF88` |
| **Wellness** | `#E8B4B8` | `#E8B4B8` |

---

## Design Principles

1. **Warm over cold** — Colors feel alive and organic
2. **Muted over saturated** — Easy on eyes, inspired by natural materials
3. **Earth + sky balance** — Grounded yet hopeful
4. **High contrast for accessibility** — WCAG AA compliant combinations
5. **Seasonal awareness** — Subtle shifts possible for time of year

---

## Flutter Implementation

```dart
// lib/shared/theme/solarpunk_colors.dart

class SolarpunkColors {
  // Primary
  static const solarGold = Color(0xFFF4A825);
  static const mossGreen = Color(0xFF4A7C59);
  static const skyTeal = Color(0xFF2D8A8A);
  static const terracotta = Color(0xFFC1694F);
  static const cream = Color(0xFFF5F1E6);
  static const deepLoam = Color(0xFF1E2420);

  // Secondary
  static const dawnPink = Color(0xFFE8B4B8);
  static const lichen = Color(0xFF8FBC8F);
  static const amberGlow = Color(0xFFFFD93D);
  static const stormCloud = Color(0xFF5D6D7E);
  static const copper = Color(0xFFB87333);

  // Paths
  static const pathSelf = Color(0xFF228B22);
  static const pathKin = Color(0xFFEB9E3A);
  static const pathWork = Color(0xFF1A5F5F);
  static const pathCommons = Color(0xFF9CAF88);
}
```

---

## CSS Variables

```css
:root {
  /* Primary */
  --solar-gold: #F4A825;
  --moss-green: #4A7C59;
  --sky-teal: #2D8A8A;
  --terracotta: #C1694F;
  --cream: #F5F1E6;
  --deep-loam: #1E2420;

  /* Secondary */
  --dawn-pink: #E8B4B8;
  --lichen: #8FBC8F;
  --amber-glow: #FFD93D;
  --storm-cloud: #5D6D7E;
  --copper: #B87333;

  /* Paths */
  --path-self: #228B22;
  --path-kin: #EB9E3A;
  --path-work: #1A5F5F;
  --path-commons: #9CAF88;
}
```

---

## Inspiration Sources

- Morning sunlight through forest canopy
- Aged copper and terracotta planters
- Moss-covered stone
- Honey and amber
- Ocean at golden hour
- Lichen on bark
- Storm clouds before rain
- Dawn and dusk skies

---

## Navigation Structure (Solarpunk)

Transform the tactical menu into an organic, nature-inspired navigation system.

### Section Mapping

| Current | Solarpunk | Icon | Meaning |
|---------|-----------|------|---------|
| UIO | **CANOPY** | 🌿 | Overview & awareness (looking out from treetops) |
| PATHS | **ROOTS** | 🌱 | Personal cultivation & life management |
| SENSES | **MYCELIUM** | 🍄 | Interconnected awareness network |
| LEGACY | **SYMBIOSIS** | 🦋 | AI partnership & mutual benefit |
| HIVE | **HABITAT** | 🏡 | Physical world & environment |
| CORE | **SEED** | 🌰 | Foundation & settings |

---

### CANOPY (Command Center → Forest Overview)

| Current | Solarpunk | Icon | Color |
|---------|-----------|------|-------|
| Dashboard | **Grove** | `forest` | Moss Green |
| Command | **Commune** | `chat_bubble_outline` | Solar Gold |
| Eagle Eye | **Watchtower** | `visibility` | Sky Teal |
| Voice | **Whisper** | `mic` | Dawn Pink |
| Inbox | **Gather** | `inbox` | Amber Glow |
| Activity | **Pulse** | `favorite` | Terracotta |
| Overview | **Horizon** | `landscape` | Sky Teal |
| Edge | **Frontier** | `explore` | Lichen |
| Pipeline | **Stream** | `water_drop` | Sky Teal |
| Metrics | **Harvest** | `insights` | Solar Gold |

---

### ROOTS (Life Paths → Personal Cultivation)

| Current | Solarpunk | Icon | Color |
|---------|-----------|------|-------|
| Tasks | **Tend** | `eco` | Moss Green |
| Projects | **Groves** | `park` | Forest Green |
| Planner | **Seasons** | `calendar_month` | Warm Honey |
| Calendar | **Cycles** | `event` | Sky Teal |
| Habits | **Rituals** | `self_improvement` | Lichen |
| Routines | **Rhythms** | `repeat` | Dawn Pink |
| Scheduling | **Tides** | `schedule` | Sky Teal |
| Pomodoro | **Focus** | `local_florist` | Terracotta |
| Journal | **Reflect** | `auto_stories` | Cream |
| Goals | **Aspirations** | `flag` | Solar Gold |
| Meal Plan | **Nourish** | `restaurant` | Warm Honey |
| Recipes | **Alchemy** | `menu_book` | Terracotta |
| Shopping | **Forage** | `shopping_basket` | Lichen |
| Pantry | **Stores** | `kitchen` | Copper |

---

### MYCELIUM (Senses → Interconnected Awareness)

| Current | Solarpunk | Icon | Color |
|---------|-----------|------|-------|
| Awareness | **Sensing** | `sensors` | Sky Teal |
| Consciousness | **Presence** | `psychology` | Dawn Pink |
| Presence | **Grounding** | `place` | Terracotta |
| Health Score | **Vitality** | `favorite` | Lichen |
| Trends | **Patterns** | `show_chart` | Moss Green |
| Correlations | **Threads** | `hub` | Amber Glow |
| Alerts | **Signals** | `notifications` | Terracotta |
| Log | **Chronicle** | `history_edu` | Storm Cloud |

---

### SYMBIOSIS (Legacy AI → Mutual Partnership)

| Subgroup | Solarpunk Name |
|----------|----------------|
| CORE AGENTS | **COMPANIONS** |
| AUTOMATION | **FLOWS** |
| CAPABILITIES | **GIFTS** |
| MEMORY & LEARNING | **WISDOM** |
| MONITORING | **REFLECTION** |

| Current | Solarpunk | Icon | Color |
|---------|-----------|------|-------|
| Agents | **Companions** | `diversity_3` | Solar Gold |
| Teams | **Circles** | `groups` | Warm Honey |
| Models | **Minds** | `psychology` | Sky Teal |
| Workflows | **Flows** | `account_tree` | Lichen |
| Approvals | **Consent** | `handshake` | Moss Green |
| Guardrails | **Boundaries** | `fence` | Terracotta |
| Tools | **Instruments** | `construction` | Copper |
| Skills | **Craft** | `architecture` | Forest Green |
| Knowledge | **Library** | `local_library` | Amber Glow |
| MCP | **Bridges** | `cable` | Sky Teal |
| Memory | **Remember** | `history` | Dawn Pink |
| Sessions | **Conversations** | `forum` | Lichen |
| Entities | **Beings** | `category` | Moss Green |
| Culture | **Values** | `volunteer_activism` | Warm Honey |
| Traces | **Paths** | `route` | Storm Cloud |
| Evals | **Growth** | `trending_up` | Lichen |

---

### HABITAT (Hive → Living Environment)

| Current | Solarpunk | Icon | Color |
|---------|-----------|------|-------|
| Home Assistant | **Dwelling** | `cottage` | Warm Honey |
| Vehicles | **Journeys** | `electric_car` | Sky Teal |
| Robotics | **Helpers** | `smart_toy` | Lichen |
| Drones | **Scouts** | `flight` | Storm Cloud |
| Glasses | **Vision** | `visibility` | Dawn Pink |
| Watch | **Pulse** | `watch` | Terracotta |

---

### SEED (Core → Foundation)

| Current | Solarpunk | Icon | Color |
|---------|-----------|------|-------|
| Profile | **Identity** | `person` | Warm Honey |
| Databases | **Archives** | `storage` | Storm Cloud |
| Integrations | **Connections** | `hub` | Sky Teal |
| Settings | **Configure** | `tune` | Moss Green |

---

## Navigation Visual Design

### Header Transformation

```
CURRENT (Tactical):
┌─────────────────────────────────┐
│  [LOGO]  ONEMIND               │
│          OPERATING SYSTEM       │
│  ● SYS  ● AWR    [AWARE]       │
└─────────────────────────────────┘

SOLARPUNK (Organic):
┌─────────────────────────────────┐
│  [🌿]  OneMind                  │
│        living system            │
│  ○ rooted  ○ present  [dawn]   │
└─────────────────────────────────┘
```

### Awareness Modes → Natural States

| Current | Solarpunk | Color | Meaning |
|---------|-----------|-------|---------|
| Dormant | **Resting** | Storm Cloud | Regenerating |
| Aware | **Dawn** | Amber Glow | Awakening |
| Present | **Day** | Solar Gold | Active growth |
| Omnipresent | **Zenith** | Lichen | Full presence |

### Section Accent Colors

| Section | Current | Solarpunk |
|---------|---------|-----------|
| CANOPY | Red bar | Moss Green gradient |
| ROOTS | Red bar | Warm Honey gradient |
| MYCELIUM | Red bar | Sky Teal gradient |
| SYMBIOSIS | Red bar | Solar Gold gradient |
| HABITAT | Red bar | Terracotta gradient |
| SEED | Red bar | Lichen gradient |

---

## Icon Replacements

Replace tactical/military icons with organic alternatives:

| Context | Current | Solarpunk |
|---------|---------|-----------|
| Active | `radio_button_checked` | `local_florist` |
| Alert | `warning` | `eco` |
| Success | `check_circle` | `park` |
| Error | `error` | `water_drop` (empty) |
| Loading | `hourglass` | `autorenew` (leaf spin) |
| Search | `search` | `explore` |
| Menu | `menu` | `forest` |
| Back | `arrow_back` | `undo` |
| Close | `close` | `close` (softened) |
| Add | `add` | `add_circle` |
| Delete | `delete` | `compost` |
| Edit | `edit` | `edit_note` |
| Save | `save` | `bookmark` |
| Share | `share` | `share` |
| Favorite | `star` | `favorite` |

---

## Status Bar Evolution

```
CURRENT:
┌─────────────────────────────────┐
│ ● CONNECTED           v2.0     │
└─────────────────────────────────┘

SOLARPUNK:
┌─────────────────────────────────┐
│ ○ rooted & growing    v2.0 🌱  │
└─────────────────────────────────┘
```

| Status | Current | Solarpunk |
|--------|---------|-----------|
| Connected | `CONNECTED` | `rooted` |
| Disconnected | `OFFLINE` | `dormant` |
| Syncing | `SYNCING` | `growing` |
| Error | `ERROR` | `wilting` |

---

## Animation Concepts

| Action | Current | Solarpunk |
|--------|---------|-----------|
| Loading | Spinner | Unfurling leaf |
| Success | Checkmark pop | Flower bloom |
| Error | Shake | Gentle wilt |
| Transition | Slide | Fade like fog |
| Expand | Accordion | Branch growth |
| Collapse | Accordion | Leaf fold |

---

## Typography Suggestions

| Element | Current | Solarpunk |
|---------|---------|-----------|
| Headers | Monospace, ALL CAPS | Rounded sans, Title Case |
| Body | System default | Warm humanist sans (e.g., Nunito, Quicksand) |
| Labels | Monospace | Rounded mono (e.g., JetBrains Mono) |
| Numbers | Monospace | Tabular figures, same as body |

---

## Sample Navigation Code

```dart
// Solarpunk section colors
Map<String, Color> sectionColors = {
  'CANOPY': SolarpunkColors.mossGreen,
  'ROOTS': SolarpunkColors.pathKin,  // Warm Honey
  'MYCELIUM': SolarpunkColors.skyTeal,
  'SYMBIOSIS': SolarpunkColors.solarGold,
  'HABITAT': SolarpunkColors.terracotta,
  'SEED': SolarpunkColors.lichen,
};

// Natural state colors (awareness modes)
Map<String, Color> stateColors = {
  'resting': SolarpunkColors.stormCloud,
  'dawn': SolarpunkColors.amberGlow,
  'day': SolarpunkColors.solarGold,
  'zenith': SolarpunkColors.lichen,
};
```
