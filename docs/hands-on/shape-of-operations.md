---
description: Understanding the shape of JsonLogic operations
---

# The Shape of Operations

Every JsonLogic expression follows **one simple pattern**: an object with a single key (the operator) whose value is an array of parameters.

## The pattern

```json
{ "operator": [ param1, param2, ... ] }
```

That's it. Every operation in JsonLogic looks like this.

!!! tip

    Think of it as a function call written inside-out:
    `operator(param1, param2)` becomes `{ "operator": [param1, param2] }`

## A concrete example

Here's the `==` (equals) operator:

```json title="Equality check" hl_lines="1"
{ "==": [ 1, 1 ] }
```

Breaking it down:

| Part | Meaning |
|------|---------|
| `"=="` | The **operator** (equals) |
| `[ 1, 1 ]` | The **parameters** (two values to compare) |

Result: `true`

## More examples

**Greater than:**

```json title="Greater than" hl_lines="1"
{ ">": [ 10, 5 ] }
```

Result: `true`

&nbsp;

**String concatenation:**

```json title="Concatenate strings" hl_lines="1"
{ "cat": [ "hello", " ", "world" ] }
```

Result: `"hello world"`

## Nesting

Parameters can themselves be operations. This is how you compose logic:

```json title="Nested operations" hl_lines="1 2 3"
{ ">": [
    { "+": [ 3, 2 ] },
    4
]}
```

Here, `{ "+": [3, 2] }` evaluates to `5`, then `{ ">": [5, 4] }` evaluates to `true`.

## Key takeaways

- **One key** per object — the operator name
- **One array** as the value — the parameters
- **Parameters can be nested operations**, plain values, or variable references (covered later)

Next: [Basic Operators](./basic-operators.md)
