---
status: proposed
author: @andreyturkov
created: 2025-08-14
updated: 2025-08-14
---

# Extending Flag Definition with a Type Property

## Background

Currently, `flagd` has inconsistent behavior in type validation between its `Resolve<T>` and `ResolveAll` API methods.
The `Resolve<T>` method validates the evaluated flag variant against the type `T` requested by the client, while `ResolveAll` validates it against the type of the `defaultVariant` specified in the flag's definition.
This discrepancy can lead to situations where a flag evaluation succeeds with one method but fails with the other, depending on the evaluation context and the variant returned.
This inconsistent behavior is further detailed in bug report #1481.

The root cause of this issue is the absence of a dedicated, authoritative type definition for the flag itself.
Instead, the type is inferred from the `defaultVariant` or API itself (`T` from `Resolve<T>`), which is not always a reliable source of truth for all possible variants.
This problem is getting worse by the planned support for code-defined defaults (as detailed in the [Support Code Default ADR](https://github.com/open-feature/flagd/blob/main/docs/architecture-decisions/support-code-default.md)), which allows the `defaultVariant` to be `null`.
This makes it impossible to resolve the flag's type from the `defaultVariant`, increasing the risk of runtime errors.

## Requirements

* The new `flagType`field in the flag definition must be optional to ensure backward compatibility.
* If the `flagType`field is present, `flagd` must validate that all variant values of the flag conform to this type during initialization.
* Type mismatches found during initialization must be reported as errors.
* The `Resolve<T>` and `ResolveAll` methods must use the `flagType`field for validation when it is available.
* The implementation must be consistent with the OpenFeature specification and the flag manifest schema.

## Considered Options

* **Consistent `defaultVariant` Validation:** Align the behavior of `Resolve<T>` with `ResolveAll` by making `Resolve<T>` validate the evaluated variant against the type of the `defaultVariant`.
* **API Extension with Explicit Flag Type:** Introduce an optional `flagType`property to the flag definition to serve as the single source for type validation.

## Proposal

This proposal is to extend the flag definition with an optional `flagType`property. This approach is chosen over simply aligning the `Resolve<T>` and `ResolveAll` validation because it addresses the root cause of the type inconsistency and provides a more robust, long-term solution.

By introducing an explicit `flagType`field, it establishes a single source of truth for the flag's type, independent of its variants. This allows for early and consistent type validation during flag definition parsing, preventing type-related errors at runtime.

The new `flagdType`field will be optional to maintain backward compatibility with existing flag configurations. If the field is omitted, `flagd` will fall back to infer the flag type from its variants. As the flag schema enforces that all variants are of the same type, the type of the first variant will be used. When the `flagdType`field is present, `flagd` will enforce that all variants of the flag conform to the specified type.

This change will make the behavior of `flagd` more predictable and reliable.

### API changes

The `flagd` flag definition will be updated to include an optional `flagType`property. This property will be a string enum with the following possible values: `"boolean"`, `"string"`, `"integer"`, `"float"`, and `"object"`.
This aligns with the OpenFeature CLI and the flag manifest schema.

#### Handling of numeric types

A known challenge with this approach is the differentiation between `integer` and `float` types, as JSON does not natively distinguish between them. However, maintaining this distinction is important for several reasons:

* **Precision**: For certain use cases, such as when a flag's value represents a project number or other identifier, using floating-point numbers can lead to precision loss and unexpected behavior. Enforcing an integer type ensures that the value remains consistent and accurate.
* **Alignment with OpenFeature**: The OpenFeature specification includes both `integer` and `float` types. By supporting both, `flagd` remains consistent with the broader OpenFeature ecosystem, including the flag manifest used in the OpenFeature CLI.

To address this, `flagd` will implement the following validation logic:

* If `flagType` is set to `integer`, `flagd` will validate that all variants of the flag are whole numbers.
* If `flagType` is set to `float`, `flagd` will accept any numeric value.

This validation will be performed during the initialization of the flag definition, allowing for early detection of type mismatches.

#### JSON Schema

The following changes will be made to the `schemas/json/flags.json` file:

1. A new `flagType`property will be added to the `flag` definition:

```json
"flag": {
  "type": "object",
  "properties": {
    "flagType": {
      "title": "Flag Type",
      "description": "The type of the flag. If specified, all variants must conform to this type.",
      "type": "string",
      "enum": [
        "boolean",
        "string",
        "integer",
        "float",
        "object"
      ]
    },
    "state": {
      ...
    },
    ...
  }
}
```

1. The `booleanFlag`, `stringFlag`, `integerFlag`, `floatFlag`, and `objectFlag` definitions will be updated to enforce the `flagType`property:

```json
"booleanFlag": {
  "allOf": [
    {
      "$ref": "#/definitions/flag"
    },
    {
      "$ref": "#/definitions/booleanVariants"
    },
    {
      "properties": {
        "flagType": {
          "const": "boolean"
        }
      }
    }
  ]
}
```

Similar changes will be made to `stringFlag`, `integerFlag`, `floatFlag`, and `objectFlag` to enforce their respective types.

### Consequences

#### The good

* It improves the reliability and predictability of flag evaluations.
* It allows for early error detection of type mismatches.
* It improves the developer experience by making the API more explicit.

#### The bad

* It adds a new field to the flag definition, which developers need to be aware of.

### Timeline

* **Phase 1: Core Implementation**
    * Update the JSON schema.
    * Update the `flagd` core to support the new `flagType`field.
    * Implement the type validation logic.
    * Add unit and integration tests.
* **Phase 2: Documentation**
    * Update the `flagd` documentation to reflect the changes.

## More Information

* **Bug Report:** [https://github.com/open-feature/flagd/issues/1481](https://github.com/open-feature/flagd/issues/1481)
* **Flag schema** [https://flagd.dev/schema/v0/flags.json](https://flagd.dev/schema/v0/flags.json)
* **Flag Manifest Schema:** [https://raw.githubusercontent.com/open-feature/cli/refs/heads/main/schema/v0/flag-manifest.json](https://raw.githubusercontent.com/open-feature/cli/refs/heads/main/schema/v0/flag-manifest.json)
