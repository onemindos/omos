# Code Style Guide Reference

## Python Style (PEP 8)

### Naming Conventions
- `snake_case` for functions and variables
- `PascalCase` for classes
- `UPPER_CASE` for constants
- `_private` prefix for internal use

### Formatting
- 4 spaces for indentation (no tabs)
- Max line length: 88 characters (Black default)
- Two blank lines before top-level definitions
- One blank line between methods

### Imports
```python
# Standard library
import os
import sys

# Third party
import requests
from fastapi import FastAPI

# Local
from myapp.utils import helper
```

## TypeScript/JavaScript Style

### Naming
- `camelCase` for functions and variables
- `PascalCase` for classes and components
- `UPPER_CASE` for constants

### Best Practices
- Use `const` by default, `let` when needed
- Prefer arrow functions for callbacks
- Use async/await over .then() chains
- Destructure when possible

## Common Issues to Flag

1. **Security**: Hardcoded secrets, SQL injection, XSS
2. **Performance**: N+1 queries, unnecessary re-renders
3. **Maintainability**: Functions > 50 lines, deep nesting
4. **Testing**: Missing error handling, untested edge cases


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
