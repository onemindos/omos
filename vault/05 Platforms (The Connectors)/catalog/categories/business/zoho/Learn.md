# Zoho Learn

> **Learning Management System (LMS)**

## Overview

Zoho Learn is a knowledge management and learning management system for creating training courses, onboarding programs, and internal documentation for teams and customers.

## Status

- **Active**: ✅ Yes
- **Integration**: Embed + API
- **Last Updated**: January 25, 2026

## Use Cases

1. **Customer Onboarding** - New user training
2. **Product Documentation** - How-to guides
3. **Internal Training** - Team skill development
4. **Certification Programs** - Structured learning paths
5. **Knowledge Base** - Self-service support

## Key Features Used

- [ ] Courses & lessons
- [ ] Video content hosting
- [ ] Quizzes & assessments
- [ ] Certificates of completion
- [ ] Progress tracking
- [ ] Discussion forums
- [ ] Multi-format content (video, text, PDF, slides)
- [ ] Mobile learning

## Knowledge Base Structure

```
OneMind Learn:
├── Getting Started/
│   ├── What is OneMind?
│   ├── First Steps
│   ├── Interface Overview
│   └── Quick Start Guide
├── Core Features/
│   ├── Legacy AI Agent
│   ├── Memory System
│   ├── Task Management
│   ├── Calendar Integration
│   └── Voice Commands
├── Integrations/
│   ├── Todoist Setup
│   ├── Google Calendar Sync
│   ├── Home Assistant
│   └── API Access
├── Advanced Topics/
│   ├── Custom Workflows
│   ├── Automation Rules
│   ├── API Development
│   └── Self-Hosting
└── Troubleshooting/
    ├── Common Issues
    ├── Error Messages
    └── Contact Support
```

## Course Types

### Customer Courses (Public)

**1. OneMind Fundamentals** (30 min)

- Lesson 1: Introduction to OneMind (5 min video)
- Lesson 2: Setting Up Your Account (10 min)
- Lesson 3: Using Legacy AI (10 min)
- Lesson 4: Integrations Overview (5 min)
- Quiz: 10 questions
- Certificate: "OneMind Certified User"

**2. Advanced Workflows** (1 hour)

- Lesson 1: Automation Basics (15 min)
- Lesson 2: Custom Workflows (20 min)
- Lesson 3: API Integration (20 min)
- Lesson 4: Best Practices (5 min)
- Project: Build a workflow
- Certificate: "OneMind Advanced User"

### Internal Courses (Team Only)

**1. New Team Member Onboarding** (2 hours)

- Company overview
- Tools & systems
- Code standards
- Communication guidelines
- First week checklist

**2. Legacy AI Development** (4 hours)

- Agent architecture
- Tool creation
- Testing workflows
- Deployment process

## Content Formats

### Video Lessons

- Hosted on Zoho Learn (no YouTube needed)
- HD quality (1080p)
- Playback speed control
- Transcripts/captions
- Downloadable for offline viewing

### Text Lessons

- Rich text editor
- Markdown support
- Embedded images/GIFs
- Code blocks with syntax highlighting
- Collapsible sections

### Downloadable Resources

- PDF guides
- Templates
- Checklists
- Code samples
- Worksheets

### Interactive Elements

- Quizzes (multiple choice, true/false, fill-in-blank)
- Surveys (feedback collection)
- Polls (quick questions)
- Assignments (submit work)

## Progress Tracking

### Learner Dashboard

- Courses enrolled
- Progress per course (% complete)
- Certificates earned
- Time spent learning
- Upcoming assignments

### Admin Dashboard

- Total learners
- Course completion rates
- Average time per course
- Popular courses
- Learner engagement
- Quiz scores

## Quizzes & Assessments

### Quiz Types

**Knowledge Check** (Formative)

- During lesson
- Low stakes
- Immediate feedback
- Don't affect certificate

**Final Assessment** (Summative)

- End of course
- Required to pass
- Minimum score required (70%)
- Affects certificate eligibility

### Question Types

- Multiple choice (single answer)
- Multiple select (multiple answers)
- True/False
- Fill in the blank
- Matching
- Essay (manual grading)

### Example Quiz Questions

**Question 1 (Multiple Choice)**

```
What is the primary function of the Legacy AI agent?

A) Email management
B) Task automation and intelligent assistance
C) File storage
D) Social media posting

Correct Answer: B
Explanation: Legacy AI serves as your personal assistant...
```

**Question 2 (True/False)**

```
OneMind can integrate with Google Calendar.

True / False

Correct Answer: True
```

## Certificates

### Certificate Requirements

- Complete all lessons (100%)
- Pass final quiz (70%+ score)
- Complete assignments (if any)
- Minimum time spent (to prevent speed-clicking)

### Certificate Design

```
╔══════════════════════════════════════════════╗
║                                              ║
║         CERTIFICATE OF COMPLETION            ║
║                                              ║
║  This certifies that                         ║
║  [STUDENT NAME]                              ║
║  has successfully completed                  ║
║                                              ║
║  OneMind Fundamentals Course                 ║
║                                              ║
║  Date: [Date]                                ║
║  Certificate ID: [Unique ID]                 ║
║                                              ║
║  [Signature]                                 ║
║  Zeus Delacruz, Founder                      ║
╚══════════════════════════════════════════════╝
```

### Certificate Verification

- Unique certificate ID
- Public verification page
- QR code for quick verification
- LinkedIn integration (add to profile)

## Discussion Forums

### Forum Categories

- **General Discussion**: Open topics
- **Questions & Answers**: Help from peers
- **Feature Requests**: Suggest improvements
- **Bug Reports**: Report issues
- **Show & Tell**: Share your projects

### Forum Features

- Threaded discussions
- Upvoting (most helpful responses)
- Mark as "solved"
- Moderator tools
- Email notifications
- Search forums

## Gamification

### Engagement Features

- **Points**: Earn for completing lessons, quizzes, forum participation
- **Badges**: Unlock for milestones (First Course, Master Learner, etc.)
- **Leaderboards**: Top learners (optional, can be hidden)
- **Streaks**: Consecutive days learning

### Example Badges

- 🥇 **First Course**: Complete your first course
- 🎓 **Certified Pro**: Earn 3 certificates
- 💬 **Community Helper**: Answer 10 forum questions
- 📚 **Bookworm**: Complete 5 courses
- ⚡ **Speed Learner**: Complete course in < 1 day

## Mobile Learning

### Mobile App Features

- Download lessons for offline viewing
- Continue from where you left off
- Push notifications for new content
- Forum access
- Progress syncs with web

## Integrations

### Zoho CRM Integration

- Assign courses to leads/contacts
- Track completion in CRM
- Trigger courses based on deal stage
- Report on training completion

### Zoho Desk Integration

- Link courses to support tickets
- "Learn More" button in ticket responses
- Reduce ticket volume with self-service

### Zoho People Integration

- Onboarding courses for new hires
- Required training tracking
- Compliance training

## API Integration

### Enroll User in Course

```python
POST https://learn.zoho.com/api/v1/courses/{course_id}/enroll
{
  "email": "user@example.com",
  "send_notification": true
}
```

### Get User Progress

```python
GET https://learn.zoho.com/api/v1/users/{user_id}/progress
```

### Webhooks

- `course.completed` - User finished course
- `quiz.passed` - User passed quiz
- `certificate.earned` - Certificate issued
- `lesson.viewed` - Lesson watched

## Analytics & Reporting

### Course Analytics

- Total enrollments
- Completion rate
- Average completion time
- Drop-off points (which lessons lose people)
- Quiz scores (average, distribution)
- Satisfaction ratings

### Learner Analytics

- Progress per learner
- Time spent learning
- Quiz performance
- Engagement level (active, inactive)
- Certificates earned

## Content Creation Workflow

### Creating a New Course

1. **Plan Content**
   - Define learning objectives
   - Outline lessons
   - Identify prerequisites
   - Estimate duration

2. **Create Lessons**
   - Record videos (Loom, OBS)
   - Write text content
   - Design quizzes
   - Gather resources

3. **Build in Zoho Learn**
   - Upload videos
   - Format text
   - Add quizzes
   - Set passing criteria

4. **Review & Test**
   - Take course yourself
   - Fix errors/typos
   - Time each lesson
   - Adjust difficulty

5. **Publish & Promote**
   - Set visibility (public/private)
   - Send announcement email
   - Post in community
   - Update website

## Best Practices

- ✅ Keep lessons short (5-10 min each)
- ✅ Use videos for demos, text for reference
- ✅ Add quizzes every 2-3 lessons
- ✅ Provide downloadable resources
- ✅ Include real-world examples
- ✅ Encourage forum participation
- ✅ Update content regularly (quarterly review)
- ✅ Get learner feedback (satisfaction surveys)

## Legacy AI Integration

### Automated Actions

1. **New User Sign-Up** → Auto-enroll in "Getting Started" course
2. **Course Completed** → Send congratulations email, update CRM
3. **Quiz Failed** → Suggest relevant forum threads or support articles
4. **Inactive Learner** (no progress in 7 days) → Send reminder email
5. **Certificate Earned** → Post to social media (with permission)

### AI-Powered Features

- Course recommendations based on user role
- Smart search in knowledge base
- Auto-generate quiz questions from content
- Identify struggling learners (offer 1:1 help)

## Notes

- Free tier: 1 course, 5 users
- Paid plans: $1/user/month (Standard), $3/user/month (Professional)
- Unlimited courses on paid plans
- Video hosting included (no size limits on Professional)
- White-label option (remove Zoho branding)
- SCORM support (import courses from other LMS)
- Mobile apps free for all users

---

**Related**: [Desk](00-24%20SO%20(Sovereign%20Operations)/02%20Infrastructure%20(VPS,%20Servers%20&%20Networking)%20(The%20Body)/Tech-Stack/platforms/business/Zoho/Desk.md) | [CRM](CRM.md) | [Forms](Forms.md)


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
