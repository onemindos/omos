# ${{ values.name }}

## Purpose
${{ values.reason }}

## Accountable reviewer
${{ values.reviewer }}

This is a review package. It does not execute a change. Never add credentials,
private backup contents, or temporary signed URLs to this directory.

## Recovery preparation
- Target: ${{ values.entity }}
- Scenario: ${{ values.scenario }}
- Evidence reference: ${{ values.evidence }}

## Objective and acceptable data loss
${{ values.recoveryObjective }}

## Preparation and execution gates
- [ ] Identify the owning runbook, accountable operator, storage and dependencies.
- [ ] Verify backup age, integrity, encryption key availability and compatible versions.
- [ ] Identify a separate recovery environment; preserve original production volumes.
- [ ] Define restore validation, failure criteria and a rollback/cutover decision.
- [ ] Review identity, private routing, sealed-secret keys and offsite dependencies.
- [ ] Obtain explicit approval for live restore or cutover after an isolated rehearsal.
- [ ] Record isolated restoration evidence and then review a separate operational change.

Portal recovery: /docs/default/component/onemind-runbooks/resilience/
Identity recovery requires its own verified procedure; the Backstage backup does not
prove Authentik recovery. Nothing in this proposal executes a restore or deletes data.
