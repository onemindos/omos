

---


# 🦨 Code Smells - What They Are

## 🎯 Definition

A **code smell** is a surface indication that usually corresponds to a deeper problem in the system. The term was popularized by Martin Fowler and Kent Beck.

### Key Points

*   **Not a bug**: Code smells are not errors - the code still works
*   **Design issues**: Indicate poor design or implementation choices
*   **Refactoring targets**: Suggest where refactoring is needed
*   **Subjective**: What smells to one developer might not to another
*   **Warning signs**: Early indicators of technical debt

### The Metaphor

Like actual smells, code smells:

*   Are easy to spot once you know what to look for
*   Indicate something might be wrong
*   Don't tell you exactly what the problem is
*   Get worse if ignored
* * *

## 🎨 Why Code Smells Matter

### Impact on Development

1. **Reduced Maintainability**: Harder to understand and modify
2. **Increased Bugs**: Complex code hides defects
3. **Slower Development**: More time spent understanding code
4. **Higher Costs**: Technical debt accumulates
5. **Team Friction**: Poor code quality frustrates developers

### Benefits of Addressing Smells

*   ✅ Easier to read and understand
*   ✅ Simpler to modify and extend
*   ✅ Fewer bugs and defects
*   ✅ Faster development velocity
*   ✅ Better team morale
* * *

## 📊 Categories of Code Smells

Code smells are organized into several categories:

### 1\. **Bloaters**

Code that has grown too large or complex:

*   Long Method
*   Large Class
*   Primitive Obsession
*   Long Parameter List
*   Data Clumps

### 2\. **Object-Orientation Abusers**

Incorrect or incomplete application of OOP principles:

*   Switch Statements
*   Temporary Field
*   Refused Bequest
*   Alternative Classes with Different Interfaces

### 3\. **Change Preventers**

Structures that make changes difficult:

*   Divergent Change
*   Shotgun Surgery
*   Parallel Inheritance Hierarchies

### 4\. **Dispensables**

Unnecessary code that should be removed:

*   Comments (excessive)
*   Duplicate Code
*   Lazy Class
*   Data Class
*   Dead Code
*   Speculative Generality

### 5\. **Couplers**

Excessive coupling between classes:

*   Feature Envy
*   Inappropriate Intimacy
*   Message Chains
*   Middle Man
* * *

## 🔍 How to Detect Code Smells

### Manual Review

*   Code reviews with team
*   Pair programming sessions
*   Regular refactoring time

### Automated Tools

**JavaScript/TypeScript:**

*   ESLint
*   SonarQube
*   CodeClimate

**Python:**

*   Pylint
*   Radon (complexity)
*   SonarQube

**Java:**

*   PMD
*   Checkstyle
*   SonarQube

**C#:**

*   ReSharper
*   NDepend
*   SonarQube

### Metrics to Watch

*   **Cyclomatic Complexity**: > 10 is concerning
*   **Lines of Code**: Methods > 20 lines, Classes > 200 lines
*   **Coupling**: Classes depending on too many others
*   **Cohesion**: Low cohesion indicates scattered responsibility
* * *

## ⚖️ When to Refactor

### Refactor When:

*   ✅ Adding new features (make it easy to add)
*   ✅ Fixing bugs (prevent similar bugs)
*   ✅ Code review reveals issues
*   ✅ You have good test coverage
*   ✅ Team has bandwidth

### Don't Refactor When:

*   ❌ No tests exist (write tests first)
*   ❌ Under tight deadline (technical debt is intentional)
*   ❌ Code works and won't change
*   ❌ Complete rewrite would be better
* * *

## 📈 The Boy Scout Rule

> "Leave the code better than you found it."

### Practice:

*   Fix small smells when you see them
*   Improve code you're already touching
*   Don't need permission for small improvements
*   Accumulates to significant improvement
* * *

## 🎯 Prioritizing Smells

### High Priority (Fix Soon)

*   Code that changes frequently
*   Core business logic
*   Code causing bugs
*   Code blocking features

### Medium Priority (Plan to Fix)

*   Stable but smelly code
*   Minor duplication
*   Slightly long methods

### Low Priority (Maybe Never)

*   Code that rarely changes
*   Generated code
*   Third-party code
*   Prototypes/experiments
* * *

## 📚 Document Navigation

### Code Smell Categories

Detailed catalog of all code smells with examples and detection methods.

### Refactoring Techniques

Specific techniques to fix each type of code smell.

### Detection Tools

Automated tools and practices for identifying code smells.
* * *

## 💡 Key Principles

### 1\. SOLID Principles

Following SOLID reduces smells:

*   **S**ingle Responsibility
*   **O**pen/Closed
*   **L**iskov Substitution
*   **I**nterface Segregation
*   **D**ependency Inversion

### 2\. DRY (Don't Repeat Yourself)

Eliminates duplication smells

### 3\. KISS (Keep It Simple)

Prevents over-engineering

### 4\. YAGNI (You Aren't Gonna Need It)

Avoids speculative generality
* * *

## 🔄 Refactoring Cycle

```
1. Identify Smell
       ↓
2. Write Tests (if missing)
       ↓
3. Refactor
       ↓
4. Run Tests
       ↓
5. Commit
       ↓
6. Repeat
```

* * *

## 📖 Essential Reading

### Books

*   **"Refactoring" by Martin Fowler** - The definitive guide
*   **"Clean Code" by Robert C. Martin** - Best practices
*   **"Working Effectively with Legacy Code" by Michael Feathers**

### Online Resources

*   [refactoring.guru](http://refactoring.guru) - Interactive catalog
*   [sourcemaking.com](http://sourcemaking.com) - Patterns and smells
*   Martin Fowler's blog
* * *

## 🔄 Last Updated

February 2026 - Comprehensive guide to code smells and refactoring



---


# 🗂️ Code Smell Categories

## 🎈 Bloaters

Bloaters are code elements that have grown so large they're hard to work with.

### 1\. Long Method

**Problem**: Method does too much

**Example:**

```typescript
// Bad: 150 lines doing everything
function processOrder(order) {
  // validate order
  // calculate totals
  // apply discounts
  // check inventory
  // process payment
  // send confirmation
  // update analytics
  // ... 100+ more lines
}

// Good: Extract to smaller methods
function processOrder(order) {
  validateOrder(order);
  const total = calculateTotal(order);
  processPayment(order, total);
  sendConfirmation(order);
}
```

**Symptoms:**

*   Method > 20-30 lines
*   Hard to name clearly
*   Many levels of nesting
*   Multiple responsibilities

**Fix:** Extract Method, Replace Method with Method Object
* * *

### 2\. Large Class (God Class)

**Problem**: Class knows/does too much

**Example:**

```typescript
// Bad: 2000 line class
class UserManager {
  validateUser() {...}
  authenticateUser() {...}
  sendEmail() {...}
  generateReport() {...}
  processPayment() {...}
  updateDatabase() {...}
  // ... 50 more methods
}

// Good: Split responsibilities
class UserValidator {...}
class UserAuthenticator {...}
class EmailService {...}
class ReportGenerator {...}
```

**Symptoms:**

*   Class > 200-300 lines
*   Many instance variables
*   Dozens of methods
*   Hard to test

**Fix:** Extract Class, Extract Subclass, Extract Interface
* * *

### 3\. Primitive Obsession

**Problem**: Using primitives instead of small objects

**Example:**

```typescript
// Bad: Primitives everywhere
function createUser(
  name: string,
  email: string,
  phone: string,
  street: string,
  city: string,
  zip: string,
  country: string
) {...}

// Good: Value objects
class Address {
  constructor(
    public street: string,
    public city: string,
    public zip: ZipCode,
    public country: Country
  ) {}
}

class Email {
  private constructor(private value: string) {}

  static create(email: string): Result<Email> {
    if (!isValid(email)) {
      return Result.fail("Invalid email");
    }
    return Result.ok(new Email(email));
  }
}

function createUser(
  name: Name,
  email: Email,
  phone: Phone,
  address: Address
) {...}
```

**Fix:** Replace Data Value with Object, Introduce Parameter Object
* * *

### 4\. Long Parameter List

**Problem**: Methods with many parameters

**Example:**

```typescript
// Bad
function createReport(
  title: string,
  author: string,
  date: Date,
  format: string,
  includeCharts: boolean,
  includeTable: boolean,
  fontSize: number,
  pageSize: string
) {...}

// Good: Parameter object
interface ReportConfig {
  metadata: ReportMetadata;
  formatting: ReportFormat;
  content: ContentOptions;
}

function createReport(config: ReportConfig) {...}
```

**Fix:** Introduce Parameter Object, Preserve Whole Object
* * *

## 🎭 Object-Orientation Abusers

### 1\. Switch Statements

**Problem**: Type checking instead of polymorphism

**Example:**

```typescript
// Bad: Switch on type
function getSpeed(animal: Animal) {
  switch(animal.type) {
    case 'bird': return animal.flySpeed;
    case 'fish': return animal.swimSpeed;
    case 'dog': return animal.runSpeed;
  }
}

// Good: Polymorphism
interface Animal {
  getSpeed(): number;
}

class Bird implements Animal {
  getSpeed() { return this.flySpeed; }
}

class Fish implements Animal {
  getSpeed() { return this.swimSpeed; }
}
```

**Fix:** Replace Conditional with Polymorphism
* * *

### 2\. Temporary Field

**Problem**: Fields only used in certain circumstances

**Example:**

```typescript
// Bad
class Calculator {
  private tempResult: number; // Only used in calculate()

  calculate(a: number, b: number) {
    this.tempResult = a + b; // Why is this a field?
    return this.processResult(this.tempResult);
  }
}

// Good
class Calculator {
  calculate(a: number, b: number) {
    const result = a + b; // Local variable
    return this.processResult(result);
  }
}
```

**Fix:** Extract Class, Introduce Null Object
* * *

## 🚧 Change Preventers

### 1\. Divergent Change

**Problem**: One class changes for many reasons

**Example:**

```typescript
// Bad: User class changes for many reasons
class User {
  // Changes when user model changes
  updateProfile() {...}

  // Changes when database changes
  saveToDatabase() {...}

  // Changes when email service changes
  sendWelcomeEmail() {...}

  // Changes when auth logic changes
  authenticate() {...}
}

// Good: Single Responsibility
class User { updateProfile() {...} }
class UserRepository { save() {...} }
class UserEmailService { sendWelcome() {...} }
class UserAuthenticator { authenticate() {...} }
```

**Fix:** Extract Class
* * *

### 2\. Shotgun Surgery

**Problem**: One change requires many small edits everywhere

**Example:**

```typescript
// Bad: Adding logging requires changes in 20 files
class OrderService {
  create() {
    console.log("Creating order"); // Logging scattered
    // ...
  }
}

class PaymentService {
  process() {
    console.log("Processing payment"); // Everywhere!
    // ...
  }
}

// Good: Centralized logging
class OrderService {
  constructor(private logger: Logger) {}

  create() {
    this.logger.info("Creating order");
    // ...
  }
}
```

**Fix:** Move Method, Move Field, Inline Class
* * *

## 🗑️ Dispensables

### 1\. Duplicate Code

**Problem**: Same code in multiple places

**Example:**

```typescript
// Bad
class OrderReport {
  generate() {
    const header = "Report Generated: " + new Date();
    // ... generate report
  }
}

class InvoiceReport {
  generate() {
    const header = "Report Generated: " + new Date();
    // ... generate invoice
  }
}

// Good: Extract common code
class Report {
  protected getHeader() {
    return "Report Generated: " + new Date();
  }
}

class OrderReport extends Report {
  generate() {
    const header = this.getHeader();
    // ...
  }
}
```

**Fix:** Extract Method, Extract Class, Pull Up Method
* * *

### 2\. Dead Code

**Problem**: Unused code

**Example:**

```typescript
// Bad
function oldCalculation() { // Never called!
  // ...
}

if (false) { // Never executes!
  // ...
}

class UnusedClass { // Not imported anywhere
  // ...
}
```

**Fix:** Delete it! Version control keeps history
* * *

### 3\. Speculative Generality

**Problem**: "We might need it someday"

**Example:**

```typescript
// Bad: Over-engineered for future
abstract class AbstractFactoryProducerFactoryBean {
  abstract createFactory(): AbstractFactory;
}

// We only have one implementation!
class ConcreteFactory extends AbstractFactoryProducerFactoryBean {
  // ...
}

// Good: YAGNI - Keep it simple
class Factory {
  create() {...}
}
```

**Fix:** Collapse Hierarchy, Inline Class, Remove Parameter
* * *

## 🔗 Couplers

### 1\. Feature Envy

**Problem**: Method more interested in another class

**Example:**

```typescript
// Bad: Method envies Order class
class ReportGenerator {
  generate(order: Order) {
    const total = order.items
      .map(i => i.price * i.quantity)
      .reduce((a, b) => a + b, 0);
    // Using order data more than own data
  }
}

// Good: Move to Order
class Order {
  getTotal() {
    return this.items
      .map(i => i.price * i.quantity)
      .reduce((a, b) => a + b, 0);
  }
}

class ReportGenerator {
  generate(order: Order) {
    const total = order.getTotal();
  }
}
```

**Fix:** Move Method, Extract Method
* * *

### 2\. Inappropriate Intimacy

**Problem**: Classes too coupled

**Example:**

```typescript
// Bad: Classes know too much about each other
class Order {
  items: OrderItem[];

  getTotal() {
    // Directly accessing customer's private field!
    if (this.customer.discountRate > 0.1) {...}
  }
}

class Customer {
  discountRate: number; // Should be private!
}

// Good: Use proper encapsulation
class Order {
  getTotal() {
    const discount = this.customer.getDiscount();
    // ...
  }
}

class Customer {
  private discountRate: number;

  getDiscount() {
    return this.discountRate;
  }
}
```

**Fix:** Move Method, Extract Class, Change Bidirectional to Unidirectional
* * *

### 3\. Message Chains

**Problem**: client.getA().getB().getC().doSomething()

**Example:**

```typescript
// Bad: Law of Demeter violation
const street = customer
  .getAddress()
  .getLocation()
  .getStreet()
  .getName();

// Good: Hide delegate
class Customer {
  getStreetName() {
    return this.address.getStreetName();
  }
}

const street = customer.getStreetName();
```

**Fix:** Hide Delegate, Extract Method
* * *

## 📊 Quick Reference

| Smell | Symptom | Quick Fix |
| ---| ---| --- |
| Long Method | \> 20 lines | Extract Method |
| Large Class | \> 200 lines | Extract Class |
| Long Parameter List | \> 3-4 params | Parameter Object |
| Duplicate Code | Copy-paste | Extract Method |
| Switch Statement | Type checking | Polymorphism |
| Feature Envy | Uses other class | Move Method |
| Dead Code | Never used | Delete |
| Comments | Explain bad code | Refactor code |



---


# 🔧 Refactoring Techniques

## 🎯 Core Refactorings

### Extract Method

**When**: Long method, duplicate code
**How**: Pull code into new method

```typescript
// Before
function printOwing() {
  printBanner();

  // print details
  console.log("name: " + name);
  console.log("amount: " + getOutstanding());
}

// After
function printOwing() {
  printBanner();
  printDetails();
}

function printDetails() {
  console.log("name: " + name);
  console.log("amount: " + getOutstanding());
}
```

* * *

### Extract Class

**When**: Class doing work of two classes
**How**: Create new class, move fields/methods

```typescript
// Before
class Person {
  name: string;
  officeAreaCode: string;
  officeNumber: string;

  getTelephoneNumber() {
    return `(${this.officeAreaCode}) ${this.officeNumber}`;
  }
}

// After
class Person {
  name: string;
  officeTelephone: TelephoneNumber;

  getTelephoneNumber() {
    return this.officeTelephone.toString();
  }
}

class TelephoneNumber {
  areaCode: string;
  number: string;

  toString() {
    return `(${this.areaCode}) ${this.number}`;
  }
}
```

* * *

### Replace Conditional with Polymorphism

**When**: Switch/if-else on type code
**How**: Create subclasses

```typescript
// Before
class Bird {
  type: string;

  getSpeed() {
    switch(this.type) {
      case 'european': return 35;
      case 'african': return 40;
      case 'norwegian': return 24;
    }
  }
}

// After
abstract class Bird {
  abstract getSpeed(): number;
}

class EuropeanBird extends Bird {
  getSpeed() { return 35; }
}

class AfricanBird extends Bird {
  getSpeed() { return 40; }
}
```

* * *

## 🛠️ Tools & IDE Support

### VS Code

*   Built-in refactoring (F2, Ctrl+Shift+R)
*   Extensions: Better Comments, SonarLint

### IntelliJ/WebStorm

*   Extensive refactoring menu
*   Safe rename, extract, inline

### Visual Studio

*   ReSharper
*   Built-in refactorings
* * *

## 📚 Refactoring Catalog

Full catalog at: [refactoring.guru](http://refactoring.guru), [refactoring.com](http://refactoring.com)

---

## Related

- [[00-24 SO (Sovereign Operations)/05 Development (Code & Testing) (The Forge)/05 Development|05 Development]]
