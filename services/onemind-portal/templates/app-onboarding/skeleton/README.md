# ${{ values.name }}

## Purpose
${{ values.reason }}

## Accountable reviewer
${{ values.reviewer }}

This is a review package. It does not execute a change. Never add credentials,
private backup contents, or temporary signed URLs to this directory.

## Proposed app
- URL: ${{ values.url }}
- Hosting: ${{ values.hosting }}

## Promotion checklist
- [ ] Confirm canonical URL, ownership, account context and duplicate entries.
- [ ] Confirm private/public access boundaries and authentication.
- [ ] Confirm category, user access, subscription and data handling.
- [ ] Confirm hosting fact; a URL alone does not prove self-hosting.
- [ ] In a separate reviewed PR, add the approved descriptor to BOTH catalog/apps.yaml copies.
- [ ] Add the matching card to packages/app/src/modules/home/directory.json.
- [ ] Build and promote the portal image using the release runbook.

Paths above are relative to services/onemind-portal; production catalog is also
platform/backstage/catalog/apps.yaml. This descriptor is not registered automatically.
