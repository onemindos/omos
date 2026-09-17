# Coding Standards

Comprehensive coding standards for OneMind-OS development covering Python (backend) and Dart (frontend).

## Table of Contents

- [General Principles](#general-principles)
- [Python Standards (Backend)](#python-standards-backend)
- [Dart Standards (Frontend)](#dart-standards-frontend)
- [Git Workflow](#git-workflow)
- [Code Review Guidelines](#code-review-guidelines)

---

## General Principles

### Code Quality
- **Readability over cleverness**: Write code that others can understand
- **Self-documenting code**: Choose descriptive names over comments
- **DRY (Don't Repeat Yourself)**: Extract common patterns
- **YAGNI (You Aren't Gonna Need It)**: Don't add features speculatively
- **KISS (Keep It Simple, Stupid)**: Prefer simple solutions

### Testing
- **Minimum 40% coverage** for new code, 60% for services
- **Test behavior, not implementation**: Focus on what, not how
- **AAA pattern**: Arrange, Act, Assert
- **One assertion per test** (guideline, not rule)

### Documentation
- **Public APIs must have docstrings/documentation**
- **Complex logic needs explanatory comments**
- **Architecture decisions documented** in ADRs or docs/
- **README for each major module**

---

## Python Standards (Backend)

### Style Guide

Follow **PEP 8** with these specifics:

```python
# Line length: 100 characters (Black formatter default)
# Indentation: 4 spaces
# Quotes: Double quotes for strings
```

### Formatting

Use **Black** formatter:

```bash
black backend/
```

Configuration in `pyproject.toml`:

```toml
[tool.black]
line-length = 100
target-version = ['py311']
```

### Type Hints

**Required for all public APIs**:

```python
# ✅ Good
from typing import Optional, List
from uuid import UUID

async def get_task(
    task_id: UUID,
    user_id: str = "zeus"
) -> Optional[Task]:
    """Get task by ID."""
    pass

# ❌ Bad - no type hints
async def get_task(task_id, user_id="zeus"):
    pass
```

### Docstrings

Use **Google-style docstrings**:

```python
def complex_operation(param1: str, param2: int) -> dict:
    """
    Short one-line description.

    Longer description explaining the operation, business logic,
    and any important details.

    Args:
        param1: Description of first parameter
        param2: Description of second parameter

    Returns:
        Dictionary containing result with keys:
        - 'status': Operation status
        - 'data': Result data

    Raises:
        DatabaseError: If database operation fails
        ValueError: If parameters are invalid

    Example:
        ```python
        result = complex_operation("test", 42)
        print(result['status'])
        ```
    """
    pass
```

### Service Architecture

**Extend BaseService** for all new services:

```python
from backend.shared.base_service import BaseService
from backend.shared.exceptions import DatabaseError

class MyService(BaseService[MyModel, MyCreateSchema, MyUpdateSchema]):
    """
    Service for managing my domain.

    Extends BaseService with domain-specific business logic.
    """

    table_name = "my_table"
    model_class = MyModel

    def __init__(self, db_pool=None, nats_client=None, ws_manager=None):
        super().__init__(db_pool, nats_client, ws_manager)

    async def custom_operation(self, param: str) -> MyModel:
        """Custom business logic."""
        self._check_db_available()

        try:
            row = await self._fetch_one(
                "SELECT * FROM my_table WHERE custom_field = $1",
                param,
                error_message="Failed to fetch by custom field"
            )
            return self._row_to_model(row) if row else None
        except DatabaseError:
            raise
        except Exception as e:
            self.logger.error(f"Unexpected error: {e}")
            raise DatabaseError(f"Operation failed: {e}")
```

### Error Handling

**Use typed exceptions from `backend.shared.exceptions`**:

```python
from backend.shared.exceptions import DatabaseError, NotFoundError, ConflictError

# ✅ Good - typed exceptions
try:
    result = await service.get_item(item_id)
    if not result:
        raise NotFoundError(f"Item {item_id} not found")
except DatabaseError as e:
    logger.error(f"Database error: {e}")
    raise
except Exception as e:
    logger.error(f"Unexpected error: {e}")
    raise DatabaseError(f"Operation failed: {e}")

# ❌ Bad - bare exceptions
try:
    result = await service.get_item(item_id)
except:  # Too broad
    return None  # Silent failure
```

**Never use bare `except:`**:

```python
# ❌ Bad
try:
    risky_operation()
except:  # Catches everything including KeyboardInterrupt
    pass

# ✅ Good
try:
    risky_operation()
except (SpecificError, AnotherError) as e:
    logger.error(f"Operation failed: {e}")
    raise
```

### Async/Await

**Use async/await consistently**:

```python
# ✅ Good
async def fetch_data(user_id: str) -> List[Item]:
    async with db_pool.acquire() as conn:
        rows = await conn.fetch("SELECT * FROM items WHERE user_id = $1", user_id)
        return [Item(**dict(row)) for row in rows]

# ❌ Bad - mixing sync/async
def fetch_data(user_id: str):  # Should be async
    # ... async operations without await
```

### Database Operations

**Use parameterized queries** (never string formatting):

```python
# ✅ Good - parameterized
await conn.execute(
    "SELECT * FROM users WHERE id = $1 AND status = $2",
    user_id, status
)

# ❌ Bad - SQL injection risk
await conn.execute(
    f"SELECT * FROM users WHERE id = '{user_id}' AND status = '{status}'"
)
```

### Imports

**Order imports** per PEP 8:

```python
# Standard library
import os
import logging
from datetime import datetime
from typing import Optional, List
from uuid import UUID

# Third-party
import asyncpg
from fastapi import APIRouter, Depends

# Local application
from backend.shared.base_service import BaseService
from backend.shared.exceptions import DatabaseError
from backend.lifeos.tasks.models import Task, TaskCreate
```

Use **absolute imports**:

```python
# ✅ Good
from backend.lifeos.tasks.service import TaskService

# ❌ Bad
from ..tasks.service import TaskService
```

### Naming Conventions

```python
# Classes: PascalCase
class TaskService:
    pass

# Functions/methods: snake_case
def get_active_tasks():
    pass

# Constants: UPPER_SNAKE_CASE
MAX_RETRIES = 3
DEFAULT_TIMEOUT = 30

# Private methods: leading underscore
def _internal_helper():
    pass

# Type variables: PascalCase
T = TypeVar('T')
TCreate = TypeVar('TCreate')
```

### Testing

**Use pytest with async support**:

```python
import pytest
from uuid import uuid4

@pytest.mark.asyncio
class TestTaskService:
    async def test_create_task_success(self, task_service, test_db_pool):
        """Test successful task creation."""
        # Arrange
        data = TaskCreate(
            title="Review PRs",
            priority="high"
        )

        # Act
        task = await task_service.create(data)

        # Assert
        assert task is not None
        assert task.title == "Review PRs"
        assert task.status == "new"

    async def test_create_task_no_database(self, task_service):
        """Test error when database unavailable."""
        # Arrange
        task_service.db_pool = None

        # Act & Assert
        with pytest.raises(DatabaseError):
            await task_service.create(TaskCreate(title="Test"))
```

---

## Dart Standards (Frontend)

### Style Guide

Follow [Effective Dart](https://dart.dev/guides/language/effective-dart):

```dart
// Use dart format before committing
dart format lib/

// Analyze code
flutter analyze
```

### Line Length

**80 characters** (Dart convention):

```dart
// ✅ Good
final task = Task(
  id: id,
  title: title,
  status: status,
);

// ❌ Bad
final task = Task(id: id, title: title, status: status, description: description);
```

### Data Models

**Use Freezed** for all data models:

```dart
import 'package:freezed_annotation/freezed_annotation.dart';

part 'task.freezed.dart';
part 'task.g.dart';

@freezed
class Task with _$Task {
  const factory Task({
    required String id,
    required String title,
    String? description,
    @Default('active') String status,
    DateTime? dueDate,
    @JsonKey(name: 'created_at') DateTime? createdAt,
  }) = _Task;

  factory Task.fromJson(Map<String, dynamic> json) => _$TaskFromJson(json);
}
```

### State Management

**Use Riverpod** for all state:

```dart
import 'package:riverpod_annotation/riverpod_annotation.dart';

part 'tasks_provider.g.dart';

@riverpod
class TasksNotifier extends _$TasksNotifier {
  @override
  Future<List<Task>> build() async {
    // Load initial state
    final client = ref.watch(taskClientProvider);
    return await client.listTasks();
  }

  Future<void> createTask(TaskCreate data) async {
    final client = ref.read(taskClientProvider);
    final newTask = await client.createTask(data);

    // Update state optimistically
    state = AsyncData([...state.value!, newTask]);
  }
}

// Widget usage
class TaskList extends ConsumerWidget {
  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final tasksAsync = ref.watch(tasksNotifierProvider);

    return tasksAsync.when(
      data: (tasks) => ListView.builder(
        itemCount: tasks.length,
        itemBuilder: (context, index) => TaskTile(task: tasks[index]),
      ),
      loading: () => CircularProgressIndicator(),
      error: (error, stack) => ErrorWidget(error),
    );
  }
}
```

### API Clients

**Extend ApiClientBase**:

```dart
import 'package:dio/dio.dart';
import '../base/api_client_base.dart';

class TaskClient extends ApiClientBase {
  TaskClient({required String baseUrl}) : super(baseUrl: baseUrl);

  Future<List<Task>> listTasks({String? status}) async {
    try {
      final response = await dio.get(
        '/api/tasks',
        queryParameters: status != null ? {'status': status} : null,
      );

      return (response.data as List)
          .map((json) => Task.fromJson(json))
          .toList();
    } on DioException catch (e) {
      logger.e('Failed to list tasks: ${e.message}');
      rethrow;
    }
  }

  Future<Task> getTask(String id) async {
    try {
      final response = await dio.get('/api/tasks/$id');
      return Task.fromJson(response.data);
    } on DioException catch (e) {
      logger.e('Failed to get task $id: ${e.message}');
      rethrow;
    }
  }
}
```

### Error Handling

**Use try-catch with DioException**:

```dart
// ✅ Good
try {
  final result = await client.getData();
  return result;
} on DioException catch (e) {
  if (e.type == DioExceptionType.connectionTimeout) {
    logger.e('Connection timeout: ${e.message}');
    throw TimeoutException('Request timed out');
  } else if (e.response?.statusCode == 404) {
    logger.e('Not found: ${e.requestOptions.path}');
    throw NotFoundException('Resource not found');
  } else {
    logger.e('API error: ${e.message}');
    throw ApiException('Failed to fetch data');
  }
}

// ❌ Bad - catching everything
try {
  final result = await client.getData();
  return result;
} catch (e) {
  return null;  // Silent failure
}
```

### Naming Conventions

```dart
// Classes: PascalCase
class TaskService {}

// Methods/variables: camelCase
void createTask() {}
final taskList = [];

// Constants: lowerCamelCase
const apiUrl = 'http://localhost:8000';

// Private: leading underscore
String _internalState;
void _helperMethod() {}

// Enums: PascalCase values
enum TaskStatus {
  active,
  completed,
  archived;
}
```

### Widget Structure

**Separate widgets into small, focused components**:

```dart
// ✅ Good - focused widget
class TaskTile extends StatelessWidget {
  final Task task;
  final VoidCallback onTap;

  const TaskTile({
    Key? key,
    required this.task,
    required this.onTap,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return ListTile(
      title: Text(task.title),
      subtitle: Text(task.status),
      onTap: onTap,
    );
  }
}

// ❌ Bad - monolithic widget with too many responsibilities
class TaskScreen extends StatefulWidget {
  // 500+ lines of code handling everything
}
```

### Testing

**Widget tests with Riverpod**:

```dart
testWidgets('TaskList displays tasks', (tester) async {
  // Arrange
  final tasks = [
    Task(id: '1', title: 'Task 1', status: 'active'),
    Task(id: '2', title: 'Task 2', status: 'active'),
  ];

  // Act
  await tester.pumpWidget(
    ProviderScope(
      overrides: [
        tasksNotifierProvider.overrideWith((ref) => tasks),
      ],
      child: MaterialApp(home: TaskList()),
    ),
  );

  // Assert
  expect(find.text('Task 1'), findsOneWidget);
  expect(find.text('Task 2'), findsOneWidget);
});
```

**Unit tests for providers**:

```dart
test('TasksNotifier loads tasks', () async {
  // Arrange
  final container = ProviderContainer(
    overrides: [
      taskClientProvider.overrideWithValue(mockTaskClient),
    ],
  );

  // Act
  final tasks = await container.read(tasksNotifierProvider.future);

  // Assert
  expect(tasks, isNotEmpty);
  expect(tasks.first.title, equals('Test Task'));
});
```

---

## Git Workflow

### Commit Messages

Use **Conventional Commits** format:

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

**Types**:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation only
- `refactor`: Code refactoring (no functionality change)
- `test`: Adding/updating tests
- `chore`: Maintenance (dependencies, tooling)
- `perf`: Performance improvement
- `style`: Code style changes (formatting, no logic change)

**Examples**:

```
feat(tasks): add recurring task support

Implement RRule-based recurrence for tasks. Users can now create
tasks that repeat daily, weekly, or monthly.

Closes #42

---

fix(api): resolve SSE connection timeout

Increase SSE connection timeout from 30s to 2 minutes to prevent
premature disconnections during long agent runs.

---

docs(frontend): add architecture guide

Create comprehensive architecture documentation covering state
management, API clients, and widget patterns.

---

refactor(services): extract BaseService

Create shared base service class to eliminate duplicate CRUD code
across 50+ services. Reduces codebase by ~800 lines.

---

test(habits): add streak calculation tests

Add 15 unit tests covering streak tracking logic including:
- Consecutive check-ins
- Missed days
- Frequency-based streaks
```

### Branch Naming

```
<type>/<description>

Examples:
- feature/recurring-tasks
- fix/sse-timeout
- refactor/base-service
- docs/architecture-guide
```

### Pull Requests

**Title**: Use conventional commit format

**Description template**:

```markdown
## Summary
Brief description of changes

## Changes
- Bullet list of changes
- Focus on what and why, not how

## Testing
- How to test the changes
- Any special setup required

## Screenshots (if UI changes)
[Add screenshots]

## Checklist
- [ ] Code follows style guidelines
- [ ] Tests added/updated
- [ ] Documentation updated
- [ ] All tests passing
- [ ] Reviewed own code
```

---

## Code Review Guidelines

### As Author

**Before requesting review**:
- [ ] Run tests locally (`pytest`, `flutter test`)
- [ ] Run linters (`black`, `flutter analyze`)
- [ ] Write descriptive commit messages
- [ ] Update relevant documentation
- [ ] Self-review your own code
- [ ] Ensure CI passes

**PR best practices**:
- Keep PRs small (< 400 lines ideal)
- One feature/fix per PR
- Link related issues
- Provide context in description
- Respond to feedback promptly

### As Reviewer

**Focus areas**:
- **Correctness**: Does code do what it claims?
- **Tests**: Are there sufficient tests?
- **Security**: Any vulnerabilities?
- **Performance**: Any obvious performance issues?
- **Maintainability**: Is code readable and well-structured?
- **Documentation**: Are complex parts documented?

**Feedback style**:
- Be specific and constructive
- Suggest alternatives
- Distinguish between blocking and non-blocking comments
- Praise good patterns

**Review template**:

```markdown
## Summary
[Overall impression - approve, request changes, or comment]

## Strengths
- Good patterns or improvements worth highlighting

## Issues
### Blocking
- Critical issues that must be addressed

### Non-blocking
- Suggestions for improvement
- Nitpicks and style issues

## Questions
- Clarification questions
```

---

## Style Configuration Files

### Python: pyproject.toml

```toml
[tool.black]
line-length = 100
target-version = ['py311']

[tool.pytest.ini_options]
asyncio_mode = "auto"
testpaths = ["tests"]
python_files = "test_*.py"
python_classes = "Test*"
python_functions = "test_*"

[tool.coverage.run]
source = ["backend"]
omit = ["*/tests/*", "*/migrations/*"]

[tool.coverage.report]
fail_under = 40
```

### Dart: analysis_options.yaml

```yaml
include: package:flutter_lints/flutter.yaml

analyzer:
  exclude:
    - "**/*.g.dart"
    - "**/*.freezed.dart"
  errors:
    invalid_annotation_target: ignore

linter:
  rules:
    - prefer_const_constructors
    - prefer_const_literals_to_create_immutables
    - prefer_final_fields
    - prefer_final_locals
    - sort_child_properties_last
    - use_key_in_widget_constructors
```

---

## Resources

### Python
- [PEP 8](https://pep8.org/)
- [Black Formatter](https://github.com/psf/black)
- [Type Hints (PEP 484)](https://peps.python.org/pep-0484/)
- [pytest Documentation](https://docs.pytest.org/)

### Dart/Flutter
- [Effective Dart](https://dart.dev/guides/language/effective-dart)
- [Flutter Style Guide](https://github.com/flutter/flutter/wiki/Style-guide-for-Flutter-repo)
- [Riverpod Documentation](https://riverpod.dev/)
- [Freezed Documentation](https://pub.dev/packages/freezed)

### Git
- [Conventional Commits](https://www.conventionalcommits.org/)
- [How to Write a Git Commit Message](https://chris.beams.io/posts/git-commit/)


---

## Related

- [[00-24 SO (Sovereign Operations)/05 Development (Code & Testing) (The Forge)/05 Development|05 Development]]
