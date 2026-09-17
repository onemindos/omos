---
id: ui-06-20260111060003
title: Tpl Project With Todoist
type: reference
pillar: UI
domain: 00-24
folder: 00-24 UI (Unified Intelligence)
status: active
created: '2026-01-11'
modified: '2026-01-11'
owner: zeus
tags:
- domain/ui
- folder/06
todoist_project: 00-24 UI (Unified Intelligence)
ai_summary: Tpl Project With Todoist
---

# 🎯 PROJECT: {{title}}

> **Pillar**: {{HP|LE|GE|UI}}
> **Area**: {{area-number}} {{area-name}}
> **Status**: 🟡 Planning | 🟢 Active | 🔵 Complete | 🔴 Blocked | ⏸️ Paused

---

## 📋 Project Tasks

### All Project Tasks

```todoist
name: {{title}} - All Tasks
filter: "##{{PILLAR}} & ##{{area-number}}-{{area-name}} & ##{{project-tag}}"
sorting:
  - priority
  - due
group: true
```

### This Week

```todoist
name: {{title}} - This Week
filter: "7 days & ##{{PILLAR}} & ##{{area-number}}-{{area-name}} & ##{{project-tag}}"
groupBy: due
```

### Overdue & Urgent

```todoist
name: {{title}} - Urgent
filter: "(overdue | today) & ##{{PILLAR}} & ##{{area-number}}-{{area-name}} & ##{{project-tag}} & p1"
sorting: due
```

---

## 📖 Overview

### What is this project?

{{Describe the project in 2-3 sentences}}

### Why does it matter?

{{Explain the significance and impact}}

### What does success look like?

{{Define clear success criteria}}

---

## 🎯 Goals & Metrics

### Primary Goal

{{Main objective}}

### Success Metrics

- **Lead Metric**: {{What you control}} - Target: {{number}}
- **Lag Metric**: {{What you measure}} - Target: {{number}}
- **Timeline**: {{Start date}} → {{End date}}

### Key Results (OKRs)

1. {{Objective 1}} - {{Measurable result}}
2. {{Objective 2}} - {{Measurable result}}
3. {{Objective 3}} - {{Measurable result}}

---

## 📦 Scope

### In Scope

- {{Item 1}}
- {{Item 2}}
- {{Item 3}}

### Out of Scope

- {{Item 1}}
- {{Item 2}}

### Assumptions

- {{Assumption 1}}
- {{Assumption 2}}

---

## 🗺️ Milestones

- [ ] **Milestone 1** - {{Description}} ({{Date}})
- [ ] **Milestone 2** - {{Description}} ({{Date}})
- [ ] **Milestone 3** - {{Description}} ({{Date}})
- [ ] **Milestone 4** - {{Description}} ({{Date}})

---

## ✅ Action Items

### To Create in Todoist

**Project Setup:**
```
Project: {{pillar-range}} {{pillar-name}} > {{area-number}} {{area-name}}
Section: {{project-name}}
Tags: ##{{PILLAR}} ##{{area-number}}-{{area-name}} ##{{project-tag}}
Labels: @{{pillar}}, @project
```

**Sample Tasks:**

1. **Task:** {{Task title}}
   - **Due:** {{Date}}
   - **Priority:** p1/p2/p3/p4
   - **Description:**
     ```
     📄 [[{{vault-path}}/{{filename}}.md]]
     🎯 Goal: {{What success looks like}}
     📊 Metrics: {{How to measure}}
     🏷️ ##{{PILLAR}} ##{{area}} ##{{project-tag}}
     ```
   - **Labels:** @{{pillar}}, @{{context}}, @{{type}}

2. **Task:** {{Task title}}
   - _[Repeat pattern above]_

---

## 📚 Resources

### Documentation

- [[Link to related doc]]
- [[Link to research]]

### Tools & Systems

- {{Tool 1}} - {{Purpose}}
- {{Tool 2}} - {{Purpose}}

### References

- {{External link}}
- {{External resource}}

---

## 👥 Stakeholders

| Role | Person | Responsibility |
|------|--------|----------------|
| Owner | {{Name}} | Overall success |
| Contributor | {{Name}} | {{Specific area}} |
| Reviewer | {{Name}} | {{Review type}} |

---

## 🚧 Risks & Blockers

### Current Blockers

- {{Blocker 1}} - _Status: {{Active|Resolved}}_
- {{Blocker 2}} - _Status: {{Active|Resolved}}_

### Potential Risks

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| {{Risk 1}} | High/Med/Low | High/Med/Low | {{Strategy}} |
| {{Risk 2}} | High/Med/Low | High/Med/Low | {{Strategy}} |

---

## 📊 Progress Tracking

### Weekly Check-in

**Current Week:** _Week of {{date}}_

- **Progress:** {{%}}
- **Velocity:** {{tasks completed}} tasks/week
- **Blockers:** {{None | List}}
- **Next Week Focus:** {{Key priorities}}

### Status History

| Date | Status | Progress | Notes |
|------|--------|----------|-------|
| {{date}} | 🟡 Planning | 0% | Project created |
| {{date}} | 🟢 Active | 25% | Started execution |

---

## 📝 Project Log

### {{date:YYYY-MM-DD}}

- Created project
- Defined scope and goals
- Set up Todoist structure

### {{date:YYYY-MM-DD}}

- {{What happened}}
- {{Decisions made}}
- {{Next actions}}

---

## 🔗 Related Projects

- [[Link to related project 1]]
- [[Link to related project 2]]

---

## 📎 Attachments

_Links to relevant files, designs, documents_

---

## ✨ Post-Project Review

### What Worked Well

- {{Success 1}}
- {{Success 2}}

### What Could Be Improved

- {{Improvement 1}}
- {{Improvement 2}}

### Lessons Learned

- {{Lesson 1}}
- {{Lesson 2}}

### Next Actions

- [ ] Archive project in Todoist
- [ ] Update project status to 🔵 Complete
- [ ] Document lessons in [[Knowledge Base]]
- [ ] Share wins with team/family

---

**Tags:** #project #{{pillar}} #{{area}}

---

*Project managed via ONE MIND CODEX + Todoist Integration*
*Last Updated: {{date:YYYY-MM-DD}}*
