# OneMind Home launcher

Home adds direct app-launch favorites backed by the existing Backstage Software Catalog
star API. No second favorites database is introduced. Signed-in storage uses the installed
app-module-user-settings feature and existing user-settings backend; persistence must be
verified on the deployed portal before claiming successful account sync.

The app directory now has neutral local initial icons, star controls, a Favorites only filter,
and a searchable Manage favorites dialog. Existing catalog stars map to matching apps and
app stars appear in the native Starred Entities widget. Stars for entities outside the app
registry are preserved. Duplicate exact launch URLs appear once in the UI; their source
catalog identities and exact account/fragments remain unchanged. Removing a favorite clears
only starred aliases for that app URL. A failed save displays an error rather than claiming
success. No default favorites are silently selected; starter apps are suggestions.

The existing Directory widget includes Favorites, so saved Home layouts receive quick access
without a layout reset. A standalone Favorite apps widget is available in Edit > Add Widget.
New/default layouts emphasize the directory, native World Clocks (New York/UTC), Toolkit,
Starred Entities and Recently Visited. Today remains a dedicated operational page and an
optional Home widget. Saved user layout choices remain intact.

Visit tracking uses Backstage's native visit listener/storage API. It records OneMind portal
pathnames only, excluding URL query strings and fragments, and does not track external app
usage or send analytics to third parties. Recently/Most Visited widgets are available in Home
editing. Data belongs to the signed-in user's settings. No external favicon service is used.

Acceptance: TypeScript and production build; immutable CI image; pinned GitOps rollout;
star/unstar and search from the live Home dialog; persistence across reload; a matching native
catalog star; light/dark layout; keyboard-accessible named star controls; provider/auth settings
preserved; two ready replicas and Argo Healthy/Synced. Account persistence verified against
one user's settings is not an exhaustive cross-device/multi-user test.

References:
- https://backstage.io/docs/getting-started/homepage/
- https://gethomepage.dev/configs/bookmarks/
- https://homarr.dev/docs/widgets/bookmarks/

## Flow layout and color repair

Browser inspection found the full app directory in a fixed-size native grid card: its
MuiCardContent was 1,461 px tall with 11,485 px of content. This caused nested scrolling
and a cramped launcher. The default view now uses the supported HomePageLayoutBlueprint
with a direct, full-width Directory and naturally sized native utility widgets. App search
and category filters sit above Favorites; active filtering goes straight to matching apps.

Customize widgets retains Backstage's existing CustomHomepageGrid, saved user layouts and
all widget choices, with duplicate widget instances prevented. The launcher view does not
write or reset those saved layouts. Favorites continue to use the same catalog star API.

The palette now has teal/blue primary accents, purple secondary accents, category colors,
colored local app initials, warm gold active stars and gradient panels over the original
background assets. The OneMind logo remains unchanged. Star controls have 44 px hit targets.
No external image or favicon source, new package, backend permission or credential was added.

Verification requires actual natural scrolling and search/pinning in the default flow view,
readable light/dark rendering, access to the native widget view, healthy pinned rollout and
preserved preferences. Compilation is not visual acceptance evidence.
