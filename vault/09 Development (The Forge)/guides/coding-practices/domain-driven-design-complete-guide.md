# 📐 Domain-Driven Design (DDD)

## 🎯 What is Domain-Driven Design?

Domain-Driven Design (DDD) is a software development approach introduced by Eric Evans in his 2003 book "Domain-Driven Design: Tackling Complexity in the Heart of Software." DDD focuses on creating software that matches the business domain by deeply understanding and modeling the core business logic.

### Core Philosophy

DDD emphasizes collaboration between technical and domain experts to create a shared understanding of the problem space and build software that accurately reflects the business domain.
* * *

## 🌟 Key Principles

### 1\. Ubiquitous Language

A common language shared by developers and domain experts:

*   Used in code, documentation, and conversations
*   Reduces translation errors
*   Evolves with understanding of the domain
*   Lives in the code itself

### 2\. Focus on the Core Domain

*   Identify what makes your business unique
*   Invest most effort in the core domain
*   Delegate or simplify supporting domains
*   Recognize generic subdomains

### 3\. Bounded Contexts

*   Define clear boundaries between models
*   Each context has its own ubiquitous language
*   Different contexts can model same concepts differently
*   Explicit relationships between contexts
* * *

## 📚 DDD Structure

DDD is organized into two main areas:

### Strategic Design (High-Level)

Focuses on breaking down large systems:

*   **Bounded Contexts**: Define model boundaries
*   **Context Mapping**: Relationships between contexts
*   **Subdomains**: Core, Supporting, Generic
*   **Domain Vision Statement**: Strategic direction

### Tactical Design (Implementation)

Focuses on building blocks within a context:

*   **Entities**: Objects with identity
*   **Value Objects**: Immutable objects
*   **Aggregates**: Consistency boundaries
*   **Domain Events**: Something that happened
*   **Services**: Operations that don't belong to entities
*   **Repositories**: Persistence abstraction
*   **Factories**: Complex object creation
* * *

## 💡 When to Use DDD

### Good Fit ✅

*   Complex business logic
*   Long-lived applications
*   Large teams
*   Evolving requirements
*   High business value
*   Need for business alignment

### Poor Fit ❌

*   Simple CRUD applications
*   Short-term projects
*   Small teams with technical focus
*   Well-understood, stable domains
*   Data-centric applications
* * *

## 🎯 Benefits of DDD

1. **Better Communication**: Ubiquitous language bridges business-tech gap
2. **Business Alignment**: Software reflects business reality
3. **Flexibility**: Changes in business easier to implement
4. **Maintainability**: Clear structure and boundaries
5. **Knowledge Sharing**: Domain knowledge lives in code
6. **Reduced Complexity**: Strategic patterns manage complexity
* * *

## ⚠️ Challenges

1. **Learning Curve**: Significant investment to learn
2. **Time Investment**: Requires deep domain understanding
3. **Team Buy-in**: Needs commitment from everyone
4. **Overengineering Risk**: Can be overkill for simple domains
5. **Continuous Learning**: Domain understanding evolves
* * *

## 📖 Document Navigation

### Strategic Design

Learn how to break down complex domains into manageable contexts with clear boundaries and relationships.

### Tactical Design

Discover the building blocks for implementing domain models within bounded contexts.

### Building Blocks

Deep dive into Entities, Value Objects, Aggregates, Services, Repositories, and more.

### Implementation Guide

Practical examples and patterns for implementing DDD in real projects.
* * *

## 📚 Essential Reading

### Books

*   **"Domain-Driven Design" by Eric Evans** (2003) - The original blue book
*   **"Implementing Domain-Driven Design" by Vaughn Vernon** (2013) - Practical red book
*   **"Domain-Driven Design Distilled" by Vaughn Vernon** (2016) - Quick overview

### Online Resources

*   DDD Community: [domaindrivendesign.org](http://domaindrivendesign.org)
*   Vaughn Vernon's Blog
*   Martin Fowler's articles on DDD
* * *

## 🔄 Last Updated

February 2026 - Comprehensive DDD guide for software architecture



---


# 🗺️ Strategic Design Patterns

## 🎯 Overview

Strategic Design focuses on the high-level structure of your system, helping you break down complex domains into manageable pieces with clear boundaries and relationships.
* * *

## 🏰 Bounded Contexts

### What is a Bounded Context?

A bounded context is a clear boundary within which a particular domain model is defined and applicable. Each context has its own ubiquitous language.

### Key Concepts

*   **Model Boundary**: Clear limits where the model applies
*   **Linguistic Boundary**: Each context has its own language
*   **Team Ownership**: Typically owned by one team
*   **Technical Boundary**: Often a separate service/module

### Example

```markdown
E-commerce System:
├── Sales Context
│   └── Product (price, inventory, promotions)
├── Shipping Context
│   └── Product (weight, dimensions, fragility)
└── Catalog Context
    └── Product (description, images, categories)
```

The same "Product" means different things in each context!

### Identifying Bounded Contexts

1. Look for linguistic boundaries (same word, different meanings)
2. Find natural team boundaries
3. Identify different rates of change
4. Look for different business processes
5. Find areas with different technical needs
* * *

## 🔗 Context Mapping

Context mapping defines the relationships between bounded contexts.

### Partnership

Two contexts work together closely, teams coordinate:

```css
[Sales Context] <--Partnership--> [Inventory Context]
```

*   Joint planning
*   Coordinated releases
*   Shared success

### Shared Kernel

Two contexts share a common model subset:

```css
[Context A] <--Shared Kernel--> [Context B]
        ↓
   [Common Model]
```

*   Requires high coordination
*   Changes need consensus
*   Use sparingly

### Customer-Supplier

Upstream (supplier) serves downstream (customer):

```css
[Supplier Context] --→ [Customer Context]
```

*   Customer defines needs
*   Supplier prioritizes customer needs
*   Formal relationship

### Conformist

Downstream conforms to upstream without influence:

```css
[Upstream] --→ [Downstream (Conformist)]
```

*   No negotiation power
*   Accept upstream model as-is
*   Common with external systems

### Anticorruption Layer (ACL)

Downstream protects itself from upstream changes:

```css
[Upstream] --→ [ACL] --→ [Downstream]
```

*   Translation layer
*   Protects domain model
*   Adapts external model
*   Essential for legacy integration

### Open Host Service

Upstream provides a public API:

```css
[Context A] ←-→ [Open Host Service]
              ↑
[Context B] ←-┘
```

*   Well-defined protocol
*   Serves multiple clients
*   Versioned API

### Published Language

Standardized language for integration:

```css
[Context A] <--JSON API--> [Context B]
```

*   Standard format (JSON, XML, Protocol Buffers)
*   Documented schema
*   Common in microservices

### Separate Ways

Contexts have no connection:

```css
[Context A]    [Context B]
```

*   No integration needed
*   Duplicate functionality acceptable
*   Reduces coupling
* * *

## 🎯 Subdomains

### Core Domain

*   **What it is**: Your competitive advantage
*   **Investment**: Maximum effort and best people
*   **Example**: Amazon's recommendation engine
*   **Characteristics**:
    *   Makes your business unique
    *   High complexity
    *   Frequently changes
    *   In-house development

### Supporting Subdomain

*   **What it is**: Necessary but not differentiating
*   **Investment**: Adequate resources
*   **Example**: Order management system
*   **Characteristics**:
    *   Supports core domain
    *   Simpler than core
    *   Custom but not unique
    *   Could be outsourced

### Generic Subdomain

*   **What it is**: Common solutions
*   **Investment**: Buy or use existing solutions
*   **Example**: Authentication, payment processing
*   **Characteristics**:
    *   No competitive advantage
    *   Well-understood problem
    *   Off-the-shelf solutions available
    *   Minimize custom development

### Subdomain Classification Matrix

| Type | Complexity | Competitive Value | Approach |
| ---| ---| ---| --- |
| Core | High | High | In-house, best team |
| Supporting | Medium | Low-Medium | In-house or outsource |
| Generic | Low-Medium | None | Buy/Open source |

* * *

## 📍 Context Map Example

```sql
E-Commerce System Context Map:

┌─────────────────┐
│  User Interface │ (Conformist)
└────────┬────────┘
         │
    ┌────┴─────────────────────┐
    ↓                          ↓
┌────────────┐         ┌──────────────┐
│   Sales    │←-ACL-→│  Inventory   │ (Core Domain)
│  (Core)    │         │ (Supporting) │
└─────┬──────┘         └──────────────┘
      │
      │ Partnership
      ↓
┌────────────┐
│  Shipping  │ (Supporting)
└─────┬──────┘
      │
      │ Open Host
      ↓
┌────────────────┐
│ 3rd Party APIs │
│  (External)    │
└────────────────┘
```

* * *

## 🎯 Strategic Design Decisions

### Questions to Ask

1. **What are our core domains?**
    *   What makes us competitive?
    *   Where should we invest most?
2. **What are the bounded contexts?**
    *   Where are the model boundaries?
    *   What are the linguistic boundaries?
3. **How do contexts relate?**
    *   Which pattern fits each relationship?
    *   Do we need anticorruption layers?
4. **What can we buy vs build?**
    *   Which subdomains are generic?
    *   What's our core competency?

### Common Mistakes

❌ **Too Many Bounded Contexts**

*   Creates excessive integration complexity
*   Solution: Start with fewer, larger contexts

❌ **Ignoring Context Boundaries**

*   Leads to big ball of mud
*   Solution: Enforce clear boundaries

❌ **Treating Everything as Core**

*   Wastes resources on generic problems
*   Solution: Identify true differentiators

❌ **No Anticorruption Layers**

*   External changes break your model
*   Solution: Always protect from external systems
* * *

## 📚 Strategic Design Checklist

- [ ] Core domain identified
- [ ] Bounded contexts defined
- [ ] Context map created
- [ ] Relationships documented
- [ ] Subdomain classification done
- [ ] ACLs planned for external systems
- [ ] Team ownership assigned
- [ ] Integration patterns chosen



---


# 🔨 Tactical Design Patterns

## 🎯 Overview

Tactical Design provides the building blocks for implementing domain models within bounded contexts. These patterns help create expressive, maintainable domain logic.
* * *

## 👤 Entities

### What is an Entity?

An object defined by its identity, not its attributes. It has continuity through time and different states.

### Characteristics

*   **Identity**: Unique identifier (ID)
*   **Mutable**: State can change over time
*   **Lifecycle**: Created, modified, persisted
*   **Equality**: Compared by ID, not attributes

### Example

```typescript
class Customer {
  private id: CustomerId;
  private email: Email;
  private name: string;
  private status: CustomerStatus;

  constructor(id: CustomerId, email: Email, name: string) {
    this.id = id;
    this.email = email;
    this.name = name;
    this.status = CustomerStatus.Active;
  }

  // Identity-based equality
  equals(other: Customer): boolean {
    return this.id.equals(other.id);
  }

  // Business operations
  deactivate(): void {
    if (this.status === CustomerStatus.Active) {
      this.status = CustomerStatus.Inactive;
      // Emit domain event
      this.addEvent(new CustomerDeactivated(this.id));
    }
  }
}
```

### When to Use

*   Object needs tracking over time
*   Object has lifecycle transitions
*   Identity matters more than attributes
*   Examples: User, Order, Product, Account
* * *

## 💎 Value Objects

### What is a Value Object?

An immutable object defined by its attributes, with no conceptual identity.

### Characteristics

*   **Immutable**: Cannot change after creation
*   **No Identity**: Defined by values
*   **Replaceable**: Create new instead of modifying
*   **Equality**: Compared by all attributes

### Example

```typescript
class Money {
  private readonly amount: number;
  private readonly currency: Currency;

  constructor(amount: number, currency: Currency) {
    if (amount < 0) throw new Error("Amount cannot be negative");
    this.amount = amount;
    this.currency = currency;
  }

  // Value-based equality
  equals(other: Money): boolean {
    return this.amount === other.amount &&
           this.currency === other.currency;
  }

  // Operations return new instances
  add(other: Money): Money {
    if (!this.currency.equals(other.currency)) {
      throw new Error("Cannot add different currencies");
    }
    return new Money(this.amount + other.amount, this.currency);
  }

  // No setters - immutable!
}

class Email {
  private readonly value: string;

  constructor(value: string) {
    if (!this.isValid(value)) {
      throw new Error("Invalid email format");
    }
    this.value = value;
  }

  private isValid(email: string): boolean {
    return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email);
  }

  toString(): string {
    return this.value;
  }
}
```

### Benefits

*   **Thread-safe**: Immutability prevents race conditions
*   **Validation**: Encapsulates business rules
*   **Expressiveness**: Makes code more readable
*   **No side effects**: Safe to pass around

### When to Use

*   Measure, quantity, or description
*   No tracking needed
*   Examples: Money, Address, DateRange, Email, Color
* * *

## 📦 Aggregates

### What is an Aggregate?

A cluster of entities and value objects treated as a single unit, with one entity as the Aggregate Root.

### Characteristics

*   **Aggregate Root**: Single entry point
*   **Consistency Boundary**: Enforces invariants
*   **Transaction Boundary**: One aggregate per transaction
*   **References**: External objects referenced by ID only

### Structure

```typescript
// Aggregate Root
class Order {
  private id: OrderId;
  private customerId: CustomerId; // Reference by ID!
  private items: OrderItem[]; // Internal entities
  private total: Money;
  private status: OrderStatus;

  // Only way to add items
  addItem(product: Product, quantity: number): void {
    if (this.status !== OrderStatus.Draft) {
      throw new Error("Cannot modify submitted order");
    }

    const item = new OrderItem(product, quantity);
    this.items.push(item);
    this.recalculateTotal(); // Maintain invariant
  }

  // Invariant enforcement
  private recalculateTotal(): void {
    this.total = this.items
      .map(item => item.getSubtotal())
      .reduce((a, b) => a.add(b), Money.zero());
  }

  submit(): void {
    if (this.items.length === 0) {
      throw new Error("Cannot submit empty order");
    }
    this.status = OrderStatus.Submitted;
    this.addEvent(new OrderSubmitted(this.id, this.total));
  }
}

// Internal entity
class OrderItem {
  private productId: ProductId;
  private quantity: number;
  private price: Money;

  constructor(product: Product, quantity: number) {
    if (quantity <= 0) throw new Error("Invalid quantity");
    this.productId = product.id;
    this.quantity = quantity;
    this.price = product.price;
  }

  getSubtotal(): Money {
    return this.price.multiply(this.quantity);
  }
}
```

### Aggregate Design Rules

1. **Keep Small**: Only what needs to change together
2. **Reference by ID**: No direct object references across aggregates
3. **One Transaction**: Modify one aggregate per transaction
4. **Eventual Consistency**: Between aggregates

### Common Mistakes

❌ **Too Large**: Including too much

```typescript
// Bad: Order contains full Customer object
class Order {
  customer: Customer; // Entire customer!
}
```

✅ **Right Size**: Reference by ID

```typescript
// Good: Order references Customer by ID
class Order {
  customerId: CustomerId; // Just the ID
}
```

❌ **Multiple Aggregates in Transaction**

```typescript
// Bad
order.submit();
inventory.reserve(order.items);
await save(order, inventory); // Two aggregates!
```

✅ **Single Aggregate + Events**

```typescript
// Good
order.submit(); // Emits OrderSubmitted event
await save(order);
// Event handler separately updates inventory
```

* * *

## 📢 Domain Events

### What is a Domain Event?

Something that happened in the domain that domain experts care about.

### Characteristics

*   **Past tense**: OrderPlaced, UserRegistered
*   **Immutable**: Cannot change after creation
*   **Time-stamped**: When it occurred
*   **Rich**: Contains relevant data

### Example

```typescript
class OrderPlaced {
  readonly occurredAt: Date;
  readonly orderId: OrderId;
  readonly customerId: CustomerId;
  readonly total: Money;

  constructor(orderId: OrderId, customerId: CustomerId, total: Money) {
    this.occurredAt = new Date();
    this.orderId = orderId;
    this.customerId = customerId;
    this.total = total;
  }
}

// Entity publishes events
class Order {
  private events: DomainEvent[] = [];

  submit(): void {
    // Business logic
    this.status = OrderStatus.Submitted;

    // Publish event
    this.addEvent(new OrderPlaced(
      this.id,
      this.customerId,
      this.total
    ));
  }

  private addEvent(event: DomainEvent): void {
    this.events.push(event);
  }

  getEvents(): DomainEvent[] {
    return this.events;
  }
}

// Event handler in another context
class InventoryService {
  @EventHandler(OrderPlaced)
  async handle(event: OrderPlaced): Promise<void> {
    // React to event
    await this.reserveInventory(event.orderId);
  }
}
```

### Benefits

*   **Decoupling**: Loose coupling between contexts
*   **Audit Trail**: History of what happened
*   **Event Sourcing**: Store events as source of truth
*   **Integration**: Easy integration between contexts
* * *

## ⚙️ Domain Services

### What is a Domain Service?

An operation that doesn't naturally belong to an entity or value object.

### When to Use

*   Operation involves multiple aggregates
*   Operation is a significant business concept
*   Operation is stateless

### Example

```typescript
class MoneyTransferService {
  transfer(
    from: Account,
    to: Account,
    amount: Money
  ): TransferResult {
    // Doesn't belong to either account
    if (!from.canWithdraw(amount)) {
      return TransferResult.insufficientFunds();
    }

    from.withdraw(amount);
    to.deposit(amount);

    return TransferResult.success();
  }
}

// NOT a domain service (belongs to entity):
class Account {
  withdraw(amount: Money): void {
    // This belongs here!
  }
}
```

### Types of Services

**Domain Service**: Pure domain logic

```typescript
class PricingService {
  calculateDiscount(customer: Customer, order: Order): Money {
    // Complex pricing rules
  }
}
```

**Application Service**: Orchestration

```typescript
class OrderApplicationService {
  async placeOrder(command: PlaceOrderCommand): Promise<void> {
    // Orchestrates workflow
    const order = await this.orderRepo.find(command.orderId);
    order.submit();
    await this.orderRepo.save(order);
    await this.eventBus.publish(order.getEvents());
  }
}
```

* * *

## 📚 Repositories

### What is a Repository?

An abstraction for accessing aggregates, providing the illusion of an in-memory collection.

### Characteristics

*   **Collection Interface**: add(), find(), remove()
*   **Aggregate Focus**: One repository per aggregate root
*   **Persistence Abstraction**: Hides data access details

### Example

```typescript
interface OrderRepository {
  find(id: OrderId): Promise<Order | null>;
  findByCustomer(customerId: CustomerId): Promise<Order[]>;
  save(order: Order): Promise<void>;
  remove(order: Order): Promise<void>;
}

class PostgresOrderRepository implements OrderRepository {
  async find(id: OrderId): Promise<Order | null> {
    const row = await this.db.query(
      'SELECT * FROM orders WHERE id = $1',
      [id.value]
    );
    return row ? this.toDomain(row) : null;
  }

  async save(order: Order): Promise<void> {
    const data = this.toData(order);
    await this.db.query(
      'INSERT INTO orders (...) VALUES (...) ON CONFLICT UPDATE...',
      data
    );
  }

  private toDomain(row: any): Order {
    // Map database row to domain object
  }

  private toData(order: Order): any {
    // Map domain object to database row
  }
}
```

### Repository Guidelines

✅ **Do:**

*   One repository per aggregate root
*   Collection-like interface
*   Hide persistence details
*   Return domain objects

❌ **Don't:**

*   Expose IQueryable/database leakage
*   Create repositories for every entity
*   Put business logic in repositories
*   Return DTOs or data objects
* * *

## 🏭 Factories

### What is a Factory?

Encapsulates complex object creation logic.

### When to Use

*   Complex creation logic
*   Multiple creation methods
*   Invariants to enforce at creation

### Example

```typescript
class OrderFactory {
  createFromCart(cart: ShoppingCart, customer: Customer): Order {
    // Complex creation logic
    const order = new Order(
      OrderId.generate(),
      customer.id
    );

    cart.items.forEach(item => {
      order.addItem(item.product, item.quantity);
    });

    if (customer.isPremium()) {
      order.applyDiscount(new Percentage(10));
    }

    return order;
  }

  reconstitute(data: OrderData): Order {
    // Recreate from database
    const order = new Order(data.id, data.customerId);
    // ... restore state
    return order;
  }
}
```

* * *

## 🎯 Tactical Pattern Summary

| Pattern | Purpose | Identity | Mutable |
| ---| ---| ---| --- |
| Entity | Tracked objects | Yes | Yes |
| Value Object | Descriptive concepts | No | No |
| Aggregate | Consistency boundary | Root has ID | Yes |
| Domain Event | Something happened | No | No |
| Service | Operations | No | Stateless |
| Repository | Persistence | N/A | N/A |
| Factory | Complex creation | N/A | N/A |



---


# 💻 Implementation Guide

## 🎯 Getting Started with DDD

### Step 1: Understand the Domain

Before writing code, deeply understand the business:

**Activities:**

*   Interview domain experts
*   Observe business processes
*   Read existing documentation
*   Attend business meetings
*   Create domain glossary

**Output:**

*   Ubiquitous language document
*   Domain model sketches
*   Key business rules identified
* * *

### Step 2: Define Bounded Contexts

Identify clear boundaries in your domain:

**Questions to Ask:**

*   Where does language change meaning?
*   Which teams own which areas?
*   What changes together?
*   What are the natural seams?

**Example Contexts:**

```java
E-Commerce Platform:
├── Catalog Context (Product browsing)
├── Shopping Context (Cart, checkout)
├── Order Context (Order processing)
├── Shipping Context (Fulfillment)
├── Billing Context (Payments, invoices)
└── Customer Context (Profiles, preferences)
```

* * *

### Step 3: Model the Core Domain

Focus on what makes your business unique:

**Core Domain Checklist:**

- [ ] Identified core vs. supporting vs. generic
- [ ] Main aggregates defined
- [ ] Key entities modeled
- [ ] Value objects extracted
- [ ] Domain events identified
- [ ] Business rules captured
* * *

## 🏗️ Project Structure

### Layered Architecture

```bash
src/
├── domain/                  # Pure domain logic
│   ├── model/
│   │   ├── Order.ts        # Entities
│   │   ├── Money.ts        # Value Objects
│   │   └── OrderItem.ts
│   ├── events/
│   │   └── OrderPlaced.ts  # Domain Events
│   ├── services/
│   │   └── PricingService.ts
│   └── repositories/       # Interfaces only
│       └── OrderRepository.ts
│
├── application/            # Use cases
│   ├── commands/
│   │   └── PlaceOrder.ts
│   ├── queries/
│   │   └── GetOrderDetails.ts
│   └── services/
│       └── OrderService.ts
│
├── infrastructure/         # Technical details
│   ├── persistence/
│   │   ├── OrderRepositoryImpl.ts
│   │   └── migrations/
│   ├── messaging/
│   │   └── EventBus.ts
│   └── web/
│       └── OrderController.ts
│
└── shared/                 # Shared kernel
    └── types/
        └── Result.ts
```

* * *

## 🔧 Implementation Patterns

### Pattern 1: Rich Domain Model

**Bad: Anemic Domain Model**

```typescript
// Just data, no behavior
class Order {
  id: string;
  items: OrderItem[];
  total: number;
  status: string;
}

// Logic in service (procedural)
class OrderService {
  submit(order: Order) {
    if (order.items.length === 0) {
      throw new Error("Empty order");
    }
    order.status = "submitted";
    order.total = this.calculateTotal(order);
  }
}
```

**Good: Rich Domain Model**

```typescript
// Behavior in domain
class Order {
  private id: OrderId;
  private items: OrderItem[];
  private total: Money;
  private status: OrderStatus;

  submit(): void {
    this.validateCanSubmit();
    this.status = OrderStatus.Submitted;
    this.addEvent(new OrderSubmitted(this.id));
  }

  private validateCanSubmit(): void {
    if (this.items.length === 0) {
      throw new DomainError("Cannot submit empty order");
    }
  }

  addItem(product: Product, quantity: number): void {
    const item = new OrderItem(product, quantity);
    this.items.push(item);
    this.recalculateTotal();
  }
}
```

* * *

### Pattern 2: Value Object Validation

```typescript
class Email {
  private readonly value: string;

  private constructor(value: string) {
    this.value = value;
  }

  // Factory method with validation
  static create(value: string): Result<Email> {
    if (!value) {
      return Result.fail("Email is required");
    }

    if (!this.isValid(value)) {
      return Result.fail("Invalid email format");
    }

    return Result.ok(new Email(value));
  }

  private static isValid(email: string): boolean {
    return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email);
  }
}

// Usage
const emailResult = Email.create("user@example.com");
if (emailResult.isFailure) {
  return Response.badRequest(emailResult.error);
}
const email = emailResult.value;
```

* * *

### Pattern 3: Aggregate Persistence

```typescript
class OrderRepository {
  async save(order: Order): Promise<Result<void>> {
    const transaction = await this.db.beginTransaction();

    try {
      // Save aggregate root
      await this.saveOrder(order, transaction);

      // Save internal entities
      await this.saveOrderItems(order.items, transaction);

      // Publish domain events
      const events = order.getEvents();
      await this.eventBus.publish(events);

      await transaction.commit();
      return Result.ok();

    } catch (error) {
      await transaction.rollback();
      return Result.fail(error.message);
    }
  }
}
```

* * *

### Pattern 4: Domain Event Handling

```typescript
// Event publisher
class Order {
  private events: DomainEvent[] = [];

  submit(): void {
    this.status = OrderStatus.Submitted;
    this.addEvent(new OrderSubmitted(this.id, this.total));
  }

  private addEvent(event: DomainEvent): void {
    this.events.push(event);
  }

  getEvents(): DomainEvent[] {
    return [...this.events];
  }

  clearEvents(): void {
    this.events = [];
  }
}

// Event handler
class SendOrderConfirmationEmail {
  @EventHandler(OrderSubmitted)
  async handle(event: OrderSubmitted): Promise<void> {
    const order = await this.orderRepo.find(event.orderId);
    const customer = await this.customerRepo.find(order.customerId);

    await this.emailService.send({
      to: customer.email,
      subject: "Order Confirmation",
      template: "order-confirmation",
      data: { order }
    });
  }
}
```

* * *

### Pattern 5: Application Service

```typescript
class PlaceOrderService {
  constructor(
    private orderRepo: OrderRepository,
    private productRepo: ProductRepository,
    private eventBus: EventBus
  ) {}

  async execute(command: PlaceOrderCommand): Promise<Result<OrderId>> {
    // 1. Load aggregates
    const products = await this.productRepo.findByIds(
      command.items.map(i => i.productId)
    );

    // 2. Create aggregate
    const orderResult = Order.create(command.customerId);
    if (orderResult.isFailure) {
      return Result.fail(orderResult.error);
    }

    const order = orderResult.value;

    // 3. Execute domain logic
    command.items.forEach(item => {
      const product = products.find(p => p.id.equals(item.productId));
      order.addItem(product, item.quantity);
    });

    order.submit();

    // 4. Persist
    await this.orderRepo.save(order);

    // 5. Publish events
    await this.eventBus.publish(order.getEvents());

    return Result.ok(order.id);
  }
}
```

* * *

## 🧪 Testing DDD

### Unit Tests: Domain Logic

```typescript
describe('Order', () => {
  it('should not allow submitting empty order', () => {
    const order = new Order(orderId, customerId);

    expect(() => order.submit())
      .toThrow('Cannot submit empty order');
  });

  it('should calculate total correctly', () => {
    const order = new Order(orderId, customerId);
    const product = createTestProduct(Money.dollars(10));

    order.addItem(product, 3);

    expect(order.total).toEqual(Money.dollars(30));
  });

  it('should emit OrderSubmitted event', () => {
    const order = new Order(orderId, customerId);
    order.addItem(createTestProduct(), 1);

    order.submit();

    const events = order.getEvents();
    expect(events).toContainInstanceOf(OrderSubmitted);
  });
});
```

### Integration Tests: Repositories

```typescript
describe('OrderRepository', () => {
  it('should persist and retrieve order', async () => {
    const order = createTestOrder();
    order.addItem(createTestProduct(), 2);

    await orderRepo.save(order);

    const retrieved = await orderRepo.find(order.id);
    expect(retrieved.id).toEqual(order.id);
    expect(retrieved.items.length).toBe(1);
  });
});
```

* * *

## 🚀 Migration Strategy

### Strangler Fig Pattern

Gradually replace legacy system:

```markdown
1. Identify a subdomain
2. Build new bounded context
3. Implement anticorruption layer
4. Route new features to new context
5. Migrate old features incrementally
6. Retire legacy context
```

### Parallel Run

Run old and new side-by-side:

```sql
User Request
     ↓
  ┌──┴──┐
  ↓     ↓
Old   New
System System
  ↓     ↓
Compare results
```

* * *

## ⚠️ Common Pitfalls

### 1\. Big Ball of Mud

**Problem**: No clear boundaries
**Solution**: Start with larger contexts, split when needed

### 2\. Anemic Domain Model

**Problem**: All logic in services
**Solution**: Push behavior into entities

### 3\. Over-Engineering

**Problem**: DDD everywhere
**Solution**: Use only where complexity justifies it

### 4\. Ignoring Events

**Problem**: Tight coupling between contexts
**Solution**: Communicate via domain events

### 5\. Wrong Aggregate Boundaries

**Problem**: Performance or consistency issues
**Solution**: Keep aggregates small, design around transactions
* * *

## 📚 Tools & Libraries

### TypeScript/JavaScript

*   **TypeORM**: Object-relational mapping
*   **MediatR**: Mediator pattern
*   **NestJS**: DDD-friendly framework

### C#

*   **MediatR**: CQRS and mediator
*   **FluentValidation**: Validation
*   **Entity Framework**: OR/M

### Java

*   **Axon Framework**: DDD + Event Sourcing
*   **Spring Data**: Repositories
*   **JPA/Hibernate**: OR/M

### Python

*   **SQLAlchemy**: OR/M
*   **Pydantic**: Value object validation
*   **FastAPI**: Modern API framework
* * *

## 🎯 DDD Checklist

**Strategic Design**

- [ ] Core domain identified
- [ ] Bounded contexts defined
- [ ] Context map created
- [ ] Ubiquitous language documented

**Tactical Design**

- [ ] Aggregates designed
- [ ] Entities vs value objects clear
- [ ] Domain events identified
- [ ] Repositories defined
- [ ] Services extracted

**Implementation**

- [ ] Layered architecture in place
- [ ] Domain logic in domain layer
- [ ] Tests covering business rules
- [ ] Events published and handled

**Ongoing**

- [ ] Regular domain expert collaboration
- [ ] Refine ubiquitous language
- [ ] Refactor as understanding grows
- [ ] Document architectural decisions

---

## Related

- [[00-24 SO (Sovereign Operations)/05 Development (Code & Testing) (The Forge)/05 Development|05 Development]]
