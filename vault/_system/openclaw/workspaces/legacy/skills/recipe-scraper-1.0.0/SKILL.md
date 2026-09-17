# recipe-scraper-1.0.0

Extracts structured recipe data from any public web URL using Schema.org JSON-LD markup.

**No API keys required. Zero npm dependencies.**

## When to use

Call `recipe_scrape` when Zeus shares a recipe URL and wants any of:

- Ingredient shopping list (grocery planning)
- Nutritional breakdown (HP tracking)
- Scaling for different servings
- Adding to a meal plan database
- Comparing cook time vs nutrition across recipes

## Invoke examples

```
"Zeus: Save this recipe for later" [pastes URL]
→ recipe_scrape({ url: "https://www.allrecipes.com/recipe/..." })

"How much protein is in that Serious Eats chicken recipe?"
→ recipe_scrape → return nutrition.protein per serving × servings

"Make a grocery list from this recipe"
→ recipe_scrape → extract ingredients → format as list

"How long does this take to cook?"
→ recipe_scrape → return total_time_minutes, prep + cook breakdown
```

## Returns

```json
{
  "name": "Crispy Roast Chicken",
  "description": "...",
  "source_url": "https://...",
  "image_url": "https://...",
  "prep_time_minutes": 15,
  "cook_time_minutes": 60,
  "total_time_minutes": 75,
  "servings": 4,
  "cuisine": "American",
  "tags": ["dinner", "chicken", "gluten-free"],
  "ingredients": [
    { "amount": 1, "unit": "lb", "name": "chicken thighs" },
    { "amount": 2, "unit": "tbsp", "name": "olive oil" }
  ],
  "instructions": [
    { "step": 1, "text": "Preheat oven to 425°F." },
    { "step": 2, "text": "Pat chicken dry with paper towels." }
  ],
  "nutrition": {
    "calories": 380,
    "protein": 34,
    "carbs": 2,
    "fat": 26,
    "fiber": 0
  },
  "source_method": "json-ld"
}
```

## Site compatibility

Works on **any site using Schema.org JSON-LD** (the vast majority of major recipe sites):

- AllRecipes ✅
- Serious Eats ✅
- NYT Cooking ✅
- Minimalist Baker ✅
- Budget Bytes ✅
- Simply Recipes ✅
- The Kitchn ✅
- Epicurious ✅

Does **not** work on:

- Paywalled sites (NYT Cooking requires login)
- JavaScript-only rendered sites (Instagram, TikTok)
- Pinterest (links to other sites — scrape the destination)

## Error handling

If no recipe data is found, throws a clear error explaining why. Gracefully handles:

- HTTP errors (4xx, 5xx)
- Redirects (up to 5)
- Missing optional fields (nutrition, timing, description)
- Malformed JSON-LD (skips, tries next script block)

## Security

- Only http/https public URLs accepted
- Localhost/private IP ranges explicitly blocked (SSRF prevention)
- No credentials stored or sent
- Redirect limit: 5 hops
- Timeout: 20 seconds
