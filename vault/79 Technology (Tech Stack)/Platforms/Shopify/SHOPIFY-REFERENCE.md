# Shopify

> E-commerce platform for online stores, product management, and order fulfillment.

---

## Overview

| Attribute         | Value                               |
| ----------------- | ----------------------------------- |
| **Type**          | E-commerce Platform / API           |
| **APIs**          | Admin API, Storefront API, Webhooks |
| **Auth**          | OAuth 2.0, Private App Tokens       |
| **License**       | Commercial (SaaS)                   |
| **Documentation** | [shopify.dev](https://shopify.dev)  |

---

## What is Shopify?

Shopify provides e-commerce infrastructure:

- **Online stores** — Full storefront hosting
- **Product catalog** — Inventory and variant management
- **Orders** — Checkout, payments, fulfillment
- **Customers** — CRM and segmentation
- **Analytics** — Sales and traffic data
- **Apps** — Extensible via custom apps

---

## OneMind Usage

### Primary Use Cases

```text
┌─────────────────────────────────────────────────────────────────┐
│                      OneMind Shopify                             │
│                                                                  │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │                   Store Management                        │   │
│  │  Products • Inventory • Orders • Customers                │   │
│  └──────────────────────────────────────────────────────────┘   │
│                                                                  │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │                   Agent Automation                        │   │
│  │  Order alerts • Stock updates • Customer queries          │   │
│  └──────────────────────────────────────────────────────────┘   │
│                                                                  │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │                   Analytics & Reporting                   │   │
│  │  Sales metrics • Revenue tracking • Performance           │   │
│  └──────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
```

| Use Case                   | Description                      |
| -------------------------- | -------------------------------- |
| **Order Notifications**    | Alert Zeus on new orders         |
| **Inventory Sync**         | Track stock levels automatically |
| **Customer Support**       | Answer order status queries      |
| **Product Updates**        | Manage catalog via agent         |
| **Sales Dashboard**        | Revenue and performance metrics  |
| **Fulfillment Automation** | Trigger shipping workflows       |

---

## Authentication

### Private App Token (Recommended for OneMind)

```python
import os

SHOPIFY_STORE = "your-store"  # your-store.myshopify.com
SHOPIFY_ACCESS_TOKEN = os.getenv("SHOPIFY_ACCESS_TOKEN")
SHOPIFY_API_VERSION = "2024-01"

BASE_URL = f"https://{SHOPIFY_STORE}.myshopify.com/admin/api/{SHOPIFY_API_VERSION}"

headers = {
    "X-Shopify-Access-Token": SHOPIFY_ACCESS_TOKEN,
    "Content-Type": "application/json"
}
```

### OAuth Flow (For Public Apps)

```python
from fastapi import FastAPI, Request
from fastapi.responses import RedirectResponse
import httpx

app = FastAPI()

SHOPIFY_API_KEY = os.getenv("SHOPIFY_API_KEY")
SHOPIFY_API_SECRET = os.getenv("SHOPIFY_API_SECRET")
REDIRECT_URI = "https://api.onemind.cloud/shopify/callback"
SCOPES = "read_products,write_products,read_orders,read_customers"


@app.get("/shopify/install")
async def install(shop: str):
    """Initiate OAuth flow."""
    auth_url = (
        f"https://{shop}/admin/oauth/authorize?"
        f"client_id={SHOPIFY_API_KEY}&"
        f"scope={SCOPES}&"
        f"redirect_uri={REDIRECT_URI}"
    )
    return RedirectResponse(auth_url)


@app.get("/shopify/callback")
async def callback(request: Request, shop: str, code: str):
    """Handle OAuth callback."""
    async with httpx.AsyncClient() as client:
        response = await client.post(
            f"https://{shop}/admin/oauth/access_token",
            json={
                "client_id": SHOPIFY_API_KEY,
                "client_secret": SHOPIFY_API_SECRET,
                "code": code
            }
        )
        data = response.json()
        access_token = data["access_token"]

        # Store token securely
        await store_shop_token(shop, access_token)

    return {"status": "connected", "shop": shop}
```

---

## Admin API

### Python Client

```python
import httpx
from typing import Optional
from dataclasses import dataclass


@dataclass
class ShopifyClient:
    store: str
    access_token: str
    api_version: str = "2024-01"

    @property
    def base_url(self) -> str:
        return f"https://{self.store}.myshopify.com/admin/api/{self.api_version}"

    @property
    def headers(self) -> dict:
        return {
            "X-Shopify-Access-Token": self.access_token,
            "Content-Type": "application/json"
        }

    async def get(self, endpoint: str, params: dict = None) -> dict:
        async with httpx.AsyncClient() as client:
            response = await client.get(
                f"{self.base_url}/{endpoint}",
                headers=self.headers,
                params=params
            )
            response.raise_for_status()
            return response.json()

    async def post(self, endpoint: str, data: dict) -> dict:
        async with httpx.AsyncClient() as client:
            response = await client.post(
                f"{self.base_url}/{endpoint}",
                headers=self.headers,
                json=data
            )
            response.raise_for_status()
            return response.json()

    async def put(self, endpoint: str, data: dict) -> dict:
        async with httpx.AsyncClient() as client:
            response = await client.put(
                f"{self.base_url}/{endpoint}",
                headers=self.headers,
                json=data
            )
            response.raise_for_status()
            return response.json()


# Initialize client
shopify = ShopifyClient(
    store="your-store",
    access_token=os.getenv("SHOPIFY_ACCESS_TOKEN")
)
```

---

## Products

### List Products

```python
async def get_products(limit: int = 50, status: str = "active") -> list[dict]:
    """Get all products."""
    response = await shopify.get(
        "products.json",
        params={"limit": limit, "status": status}
    )
    return response["products"]


# Get specific product
async def get_product(product_id: int) -> dict:
    response = await shopify.get(f"products/{product_id}.json")
    return response["product"]
```

### Create Product

```python
async def create_product(
    title: str,
    body_html: str,
    vendor: str,
    product_type: str,
    variants: list[dict],
    images: list[dict] = None
) -> dict:
    """Create a new product."""
    product_data = {
        "product": {
            "title": title,
            "body_html": body_html,
            "vendor": vendor,
            "product_type": product_type,
            "variants": variants
        }
    }

    if images:
        product_data["product"]["images"] = images

    response = await shopify.post("products.json", product_data)
    return response["product"]


# Example usage
new_product = await create_product(
    title="OneMind T-Shirt",
    body_html="<p>Premium OneMind branded t-shirt</p>",
    vendor="OneMind",
    product_type="Apparel",
    variants=[
        {"option1": "Small", "price": "29.99", "sku": "OM-TS-S"},
        {"option1": "Medium", "price": "29.99", "sku": "OM-TS-M"},
        {"option1": "Large", "price": "29.99", "sku": "OM-TS-L"}
    ]
)
```

### Update Product

```python
async def update_product(product_id: int, updates: dict) -> dict:
    """Update product details."""
    response = await shopify.put(
        f"products/{product_id}.json",
        {"product": {"id": product_id, **updates}}
    )
    return response["product"]


# Update price and inventory
await update_product(product_id, {
    "title": "OneMind Premium T-Shirt",
    "variants": [{"id": variant_id, "price": "34.99"}]
})
```

---

## Inventory

### Get Inventory Levels

```python
async def get_inventory_levels(location_id: int = None) -> list[dict]:
    """Get inventory levels for all products."""
    params = {}
    if location_id:
        params["location_ids"] = location_id

    response = await shopify.get("inventory_levels.json", params)
    return response["inventory_levels"]
```

### Update Inventory

```python
async def set_inventory(
    inventory_item_id: int,
    location_id: int,
    available: int
) -> dict:
    """Set inventory quantity."""
    response = await shopify.post(
        "inventory_levels/set.json",
        {
            "inventory_item_id": inventory_item_id,
            "location_id": location_id,
            "available": available
        }
    )
    return response["inventory_level"]


async def adjust_inventory(
    inventory_item_id: int,
    location_id: int,
    adjustment: int
) -> dict:
    """Adjust inventory by amount (positive or negative)."""
    response = await shopify.post(
        "inventory_levels/adjust.json",
        {
            "inventory_item_id": inventory_item_id,
            "location_id": location_id,
            "available_adjustment": adjustment
        }
    )
    return response["inventory_level"]
```

---

## Orders

### List Orders

```python
async def get_orders(
    status: str = "any",
    limit: int = 50,
    since_id: int = None
) -> list[dict]:
    """Get orders."""
    params = {"status": status, "limit": limit}
    if since_id:
        params["since_id"] = since_id

    response = await shopify.get("orders.json", params)
    return response["orders"]


# Get recent unfulfilled orders
unfulfilled = await get_orders(status="unfulfilled", limit=10)
```

### Get Order Details

```python
async def get_order(order_id: int) -> dict:
    """Get order by ID."""
    response = await shopify.get(f"orders/{order_id}.json")
    return response["order"]
```

### Fulfill Order

```python
async def create_fulfillment(
    order_id: int,
    line_items: list[dict],
    tracking_number: str = None,
    tracking_company: str = None
) -> dict:
    """Create fulfillment for order."""
    fulfillment_data = {
        "fulfillment": {
            "line_items_by_fulfillment_order": line_items,
            "notify_customer": True
        }
    }

    if tracking_number:
        fulfillment_data["fulfillment"]["tracking_info"] = {
            "number": tracking_number,
            "company": tracking_company or "Other"
        }

    response = await shopify.post("fulfillments.json", fulfillment_data)
    return response["fulfillment"]
```

### Cancel Order

```python
async def cancel_order(order_id: int, reason: str = "customer") -> dict:
    """Cancel an order."""
    response = await shopify.post(
        f"orders/{order_id}/cancel.json",
        {"reason": reason}
    )
    return response["order"]
```

---

## Customers

### List Customers

```python
async def get_customers(limit: int = 50, query: str = None) -> list[dict]:
    """Get customers."""
    params = {"limit": limit}
    if query:
        params["query"] = query

    response = await shopify.get("customers.json", params)
    return response["customers"]


# Search customers
vip_customers = await get_customers(query="orders_count:>10")
```

### Get Customer

```python
async def get_customer(customer_id: int) -> dict:
    """Get customer by ID."""
    response = await shopify.get(f"customers/{customer_id}.json")
    return response["customer"]


async def get_customer_orders(customer_id: int) -> list[dict]:
    """Get all orders for a customer."""
    response = await shopify.get(f"customers/{customer_id}/orders.json")
    return response["orders"]
```

---

## Webhooks

### Register Webhooks

```python
async def create_webhook(topic: str, address: str) -> dict:
    """Create a webhook subscription."""
    response = await shopify.post(
        "webhooks.json",
        {
            "webhook": {
                "topic": topic,
                "address": address,
                "format": "json"
            }
        }
    )
    return response["webhook"]


# Common webhooks for OneMind
WEBHOOK_TOPICS = [
    "orders/create",
    "orders/paid",
    "orders/fulfilled",
    "orders/cancelled",
    "products/update",
    "inventory_levels/update",
    "customers/create"
]

for topic in WEBHOOK_TOPICS:
    await create_webhook(
        topic=topic,
        address=f"https://api.onemind.cloud/webhooks/shopify/{topic.replace('/', '_')}"
    )
```

### Handle Webhooks

```python
from fastapi import FastAPI, Request, HTTPException
import hmac
import hashlib
import base64

app = FastAPI()

SHOPIFY_WEBHOOK_SECRET = os.getenv("SHOPIFY_WEBHOOK_SECRET")


def verify_webhook(data: bytes, hmac_header: str) -> bool:
    """Verify webhook signature."""
    digest = hmac.new(
        SHOPIFY_WEBHOOK_SECRET.encode(),
        data,
        hashlib.sha256
    ).digest()
    computed_hmac = base64.b64encode(digest).decode()
    return hmac.compare_digest(computed_hmac, hmac_header)


@app.post("/webhooks/shopify/orders_create")
async def handle_order_created(request: Request):
    """Handle new order webhook."""
    body = await request.body()
    hmac_header = request.headers.get("X-Shopify-Hmac-SHA256")

    if not verify_webhook(body, hmac_header):
        raise HTTPException(status_code=401, detail="Invalid signature")

    order = await request.json()

    # Process new order
    await notify_new_order(order)
    await create_fulfillment_task(order)

    return {"status": "ok"}


@app.post("/webhooks/shopify/inventory_levels_update")
async def handle_inventory_update(request: Request):
    """Handle inventory update webhook."""
    body = await request.body()
    hmac_header = request.headers.get("X-Shopify-Hmac-SHA256")

    if not verify_webhook(body, hmac_header):
        raise HTTPException(status_code=401, detail="Invalid signature")

    inventory = await request.json()

    # Check for low stock
    if inventory["available"] < 10:
        await alert_low_stock(inventory)

    return {"status": "ok"}
```

---

## Agent Integration

### Shopify Tools for Legacy Agent

```python
from agno import Tool, Toolkit


class ShopifyTools(Toolkit):
    """Shopify store management tools for Legacy agent."""

    def __init__(self):
        self.client = ShopifyClient(
            store=os.getenv("SHOPIFY_STORE"),
            access_token=os.getenv("SHOPIFY_ACCESS_TOKEN")
        )

    def get_tools(self) -> list[Tool]:
        return [
            Tool(
                name="get_recent_orders",
                description="Get recent Shopify orders",
                function=self.get_recent_orders
            ),
            Tool(
                name="get_order_status",
                description="Get status of a specific order by order number",
                function=self.get_order_status
            ),
            Tool(
                name="get_inventory_status",
                description="Check inventory levels for products",
                function=self.get_inventory_status
            ),
            Tool(
                name="get_sales_summary",
                description="Get sales summary for a time period",
                function=self.get_sales_summary
            )
        ]

    async def get_recent_orders(self, limit: int = 10) -> list[dict]:
        """Get recent orders."""
        orders = await self.client.get(
            "orders.json",
            params={"limit": limit, "status": "any"}
        )

        return [
            {
                "order_number": o["order_number"],
                "created_at": o["created_at"],
                "total": o["total_price"],
                "status": o["financial_status"],
                "fulfillment": o["fulfillment_status"] or "unfulfilled",
                "customer": o.get("customer", {}).get("email")
            }
            for o in orders["orders"]
        ]

    async def get_order_status(self, order_number: str) -> dict:
        """Get order by order number."""
        orders = await self.client.get(
            "orders.json",
            params={"name": order_number}
        )

        if not orders["orders"]:
            return {"error": f"Order {order_number} not found"}

        order = orders["orders"][0]
        return {
            "order_number": order["order_number"],
            "status": order["financial_status"],
            "fulfillment": order["fulfillment_status"] or "unfulfilled",
            "total": order["total_price"],
            "items": [
                {"title": item["title"], "quantity": item["quantity"]}
                for item in order["line_items"]
            ],
            "shipping_address": order.get("shipping_address", {}).get("city")
        }

    async def get_inventory_status(self, product_title: str = None) -> list[dict]:
        """Get inventory levels."""
        products = await self.client.get("products.json", params={"limit": 50})

        results = []
        for product in products["products"]:
            if product_title and product_title.lower() not in product["title"].lower():
                continue

            for variant in product["variants"]:
                results.append({
                    "product": product["title"],
                    "variant": variant["title"],
                    "sku": variant["sku"],
                    "inventory": variant["inventory_quantity"],
                    "low_stock": variant["inventory_quantity"] < 10
                })

        return results

    async def get_sales_summary(self, days: int = 7) -> dict:
        """Get sales summary."""
        from datetime import datetime, timedelta

        since = (datetime.utcnow() - timedelta(days=days)).isoformat()

        orders = await self.client.get(
            "orders.json",
            params={
                "status": "any",
                "financial_status": "paid",
                "created_at_min": since,
                "limit": 250
            }
        )

        total_revenue = sum(float(o["total_price"]) for o in orders["orders"])
        order_count = len(orders["orders"])

        return {
            "period_days": days,
            "total_orders": order_count,
            "total_revenue": f"${total_revenue:.2f}",
            "average_order": f"${total_revenue / order_count:.2f}" if order_count else "$0.00"
        }
```

### Add to Agent

```python
from agno import Agent
from tools.shopify_tools import ShopifyTools

agent = Agent(
    name="Legacy",
    model="anthropic.claude-3-5-sonnet",
    tools=[ShopifyTools()],
    instructions="You can check orders, inventory, and sales from the Shopify store."
)
```

---

## GraphQL API

### Setup

```python
GRAPHQL_URL = f"https://{SHOPIFY_STORE}.myshopify.com/admin/api/{API_VERSION}/graphql.json"


async def graphql_query(query: str, variables: dict = None) -> dict:
    """Execute GraphQL query."""
    async with httpx.AsyncClient() as client:
        response = await client.post(
            GRAPHQL_URL,
            headers=headers,
            json={"query": query, "variables": variables or {}}
        )
        return response.json()
```

### Example Queries

```python
# Get products with inventory
PRODUCTS_QUERY = """
query GetProducts($first: Int!) {
    products(first: $first) {
        edges {
            node {
                id
                title
                status
                totalInventory
                variants(first: 10) {
                    edges {
                        node {
                            id
                            title
                            price
                            inventoryQuantity
                        }
                    }
                }
            }
        }
    }
}
"""

products = await graphql_query(PRODUCTS_QUERY, {"first": 50})


# Get order details
ORDER_QUERY = """
query GetOrder($id: ID!) {
    order(id: $id) {
        id
        name
        totalPriceSet {
            shopMoney {
                amount
                currencyCode
            }
        }
        displayFulfillmentStatus
        displayFinancialStatus
        customer {
            email
            displayName
        }
        lineItems(first: 50) {
            edges {
                node {
                    title
                    quantity
                }
            }
        }
    }
}
"""

order = await graphql_query(ORDER_QUERY, {"id": "gid://shopify/Order/123456"})
```

---

## Environment Variables

```bash
# Shopify Store
SHOPIFY_STORE=your-store
SHOPIFY_ACCESS_TOKEN=shpat_xxxxxxxxxxxxx
SHOPIFY_API_VERSION=2024-01

# OAuth (for public apps)
SHOPIFY_API_KEY=your_api_key
SHOPIFY_API_SECRET=your_api_secret

# Webhooks
SHOPIFY_WEBHOOK_SECRET=your_webhook_secret
```

---

## Rate Limiting

Shopify has API rate limits:

| Plan     | REST Limit | GraphQL Cost    |
| -------- | ---------- | --------------- |
| Basic    | 2 req/sec  | 1000 points/sec |
| Shopify  | 4 req/sec  | 1000 points/sec |
| Advanced | 4 req/sec  | 1000 points/sec |
| Plus     | 20 req/sec | 2000 points/sec |

### Handle Rate Limits

```python
import asyncio
from httpx import HTTPStatusError


async def rate_limited_request(func, *args, max_retries=3, **kwargs):
    """Execute request with rate limit handling."""
    for attempt in range(max_retries):
        try:
            return await func(*args, **kwargs)
        except HTTPStatusError as e:
            if e.response.status_code == 429:
                retry_after = float(e.response.headers.get("Retry-After", 2))
                await asyncio.sleep(retry_after)
            else:
                raise

    raise Exception("Max retries exceeded")
```

---

## Troubleshooting

### Authentication Errors

```python
# Verify token is valid
async def verify_shop_access():
    try:
        response = await shopify.get("shop.json")
        print(f"Connected to: {response['shop']['name']}")
        return True
    except Exception as e:
        print(f"Auth error: {e}")
        return False
```

### Webhook Verification Failing

```python
# Debug webhook signature
@app.post("/webhooks/debug")
async def debug_webhook(request: Request):
    body = await request.body()
    hmac_header = request.headers.get("X-Shopify-Hmac-SHA256")

    print(f"Body length: {len(body)}")
    print(f"HMAC header: {hmac_header}")
    print(f"Secret configured: {bool(SHOPIFY_WEBHOOK_SECRET)}")

    is_valid = verify_webhook(body, hmac_header)
    print(f"Signature valid: {is_valid}")

    return {"valid": is_valid}
```

### Missing Scopes

```bash
# Required scopes for OneMind
read_products
write_products
read_orders
write_orders
read_customers
read_inventory
write_inventory
read_fulfillments
write_fulfillments
```

---

## Resources

- [Shopify Developer Docs](https://shopify.dev)
- [Admin API Reference](https://shopify.dev/docs/api/admin-rest)
- [GraphQL Admin API](https://shopify.dev/docs/api/admin-graphql)
- [Webhooks Guide](https://shopify.dev/docs/apps/webhooks)
- [Rate Limits](https://shopify.dev/docs/api/usage/rate-limits)

---

_Last updated: January 2026_
