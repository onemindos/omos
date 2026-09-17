---
title: "Dog"
type: pet
status: active
created: 2026-04-12
updated: 2026-04-12
owner: zeus
access_level: L0
species: dog
tags: [type/pet, domain/LE, status/active]
ai_summary: "The Delacruz family dog. Registry-only entity — no chat access. Fill in name, breed, and vet info."
---

# Dog

> **Action required:** Fill in name, breed, DOB, vet details.

## Identity

| Field | Value |
|-------|-------|
| Name | {{fill in}} |
| Species | Dog |
| Breed | {{fill in}} |
| Date of birth | {{fill in}} |
| Color / markings | {{fill in}} |
| Microchip ID | {{fill in}} |

## Health

| Field | Value |
|-------|-------|
| Vet | {{name + clinic + address}} |
| Vet phone | {{phone}} |
| Last checkup | {{date}} |
| Vaccinations current | {{yes/no}} |
| Medications | {{any ongoing}} |
| Allergies | {{any known}} |
| Spayed/Neutered | {{yes/no}} |

## Care Routine

- **Feeding**: {{times, amount, brand}}
- **Walk schedule**: {{frequency, who walks}}
- **Grooming**: {{frequency, groomer}}
- **Boarding/sitter**: {{contact for travel}}

## Emergency

| Contact | Info |
|---------|------|
| Emergency vet | {{name + number}} |
| Pet sitter | {{name + number}} |

## Home Assistant

- Entity: {{sensor.dog_location or similar if tracked}}
- Smart collar: {{GPS collar model if applicable}}

## Notes

{{Personality, quirks, things guests should know}}


---

## Related

- [[00-24 SO (Sovereign Operations)/08 OM Assets (Entity Registry) (The Registry) (Hive Mind)/08 OM Assets|08 OM Assets]]
