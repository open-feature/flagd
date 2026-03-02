---
description: Basic JsonLogic operators used in flagd targeting
---

# Basic Operators

Now that you know [the shape](./shape-of-operations.md), let's look at the operators you'll use most.

## `if` — Conditional logic

The `if` operator is the backbone of targeting rules. It takes parameters in groups of three: **condition**, **value-if-true**, **fallback**.

```json title="Basic if" hl_lines="1"
{ "if": [ true, "yes", "no" ] }
```

Result: `"yes"`

&nbsp;

In flagd targeting, the result is a **variant name**:

```json title="if in a targeting rule" hl_lines="2 3 4 5 6"
{
  "if": [
    true,
    "on",
    "off"
  ]
}
```

| Parameter | Role |
|-----------|------|
| `true` | The **condition** to evaluate |
| `"on"` | Returned when the condition is **truthy** |
| `"off"` | Returned when the condition is **falsy** |

## Chaining conditions

You can chain multiple conditions by adding more parameters. Think of it like `if / else if / else`:

```json title="Chained if" hl_lines="2 4 6"
{ "if": [
    false,
    "first",
    false,
    "second",
    "fallback"
]}
```

Result: `"fallback"` (both conditions were false)

The pattern is:

```none
condition1, value1, condition2, value2, ..., fallback
```

## Comparison operators

These all follow the same shape and return `true` or `false`:

```json title="Equals" hl_lines="1"
{ "==": [ "a", "a" ] }
```

```json title="Not equals" hl_lines="1"
{ "!=": [ "a", "b" ] }
```

```json title="Greater than" hl_lines="1"
{ ">": [ 10, 5 ] }
```

```json title="Less than or equal" hl_lines="1"
{ "<=": [ 3, 3 ] }
```

## Combining conditions

Use `and` / `or` to combine multiple checks:

```json title="and" hl_lines="1"
{ "and": [
    { ">": [ 10, 5 ] },
    { "<": [ 10, 20 ] }
]}
```

Result: `true` (both conditions met)

&nbsp;

```json title="or" hl_lines="1"
{ "or": [
    { "==": [ 1, 2 ] },
    { "==": [ 3, 3 ] }
]}
```

Result: `true` (second condition met)

## Negation

```json title="not" hl_lines="1"
{ "!": [ true ] }
```

Result: `false`

## Putting it together

A real-world-ish targeting rule that combines `if` with a comparison:

```json title="Complete targeting rule" hl_lines="3 4 5"
{
  "if": [
    { ">=": [ 18, 13 ] },
    "teen-feature",
    "default-feature"
  ]
}
```

Of course, hard-coded values aren't very useful. Next we'll learn how to pull dynamic values from the evaluation context.

Next: [Variables with `var`](./variables.md)
