---
name: code-review
description: Code review with style checking, security analysis, and best practices
license: Apache-2.0
metadata:
  version: "1.0.0"
  author: zeus
  tags: ["python", "typescript", "code-quality", "security"]
  pillar: "IT"
---

# Code Review Skill

Use this skill when reviewing code for quality, security, and best practices.

## When to Use

- User asks for code review or feedback
- User wants to improve code quality
- User needs help with refactoring
- Reviewing PRs or commits
- Security audit of code

## Review Process

### 1. Understand Context
- What does this code do?
- What's the broader system context?
- What are the requirements?

### 2. Check Structure
- File organization
- Module boundaries
- Dependency management
- Naming conventions

### 3. Analyze Logic
- Algorithm correctness
- Edge cases handled
- Error handling
- Performance considerations

### 4. Security Review
- Input validation
- Authentication/authorization
- Data sanitization
- Secrets management
- OWASP Top 10 vulnerabilities

### 5. Style & Readability
- Consistent formatting
- Clear naming
- Appropriate comments
- Documentation

### 6. Testing
- Test coverage
- Test quality
- Edge case tests
- Integration tests

## Output Format

```markdown
## Code Review Summary

**Overall:** ✅ Approved / ⚠️ Needs Changes / ❌ Requires Rework

### Critical Issues
- [ ] Issue 1 (security/correctness)

### Improvements
- [ ] Suggestion 1
- [ ] Suggestion 2

### Positive Notes
- Good pattern used here
- Clean implementation of X

### Detailed Feedback
[Line-by-line comments]
```

## Language-Specific Guidelines

### Python
- Follow PEP 8
- Use type hints
- Prefer pathlib over os.path
- Use context managers

### TypeScript
- Strict mode enabled
- Proper typing (avoid `any`)
- Use const/let, not var
- Async/await over callbacks

### Dart/Flutter
- Follow effective_dart
- Use const constructors
- Proper widget composition
- State management patterns

## Best Practices

- Be constructive, not critical
- Explain the "why" behind suggestions
- Provide code examples for fixes
- Prioritize security and correctness
- Acknowledge good patterns
