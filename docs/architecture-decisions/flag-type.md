---
# Valid statuses: draft | proposed | rejected | accepted | superseded
status: proposed
author: @andreyturkov
created: 2025-08-14
updated: 2025-08-14
---

# Extending Flag Definition with a Type Property


## Background

Currently, `flagd` has inconsistent behavior in type validation between its `Resolve<T>` and `ResolveAll` API methods. The `Resolve<T>` method validates the evaluated flag variant against the type `T` requested by the client, while `ResolveAll` validates it against the type of the `defaultVariant` specified in the flag's definition. This discrepancy can lead to situations where a flag evaluation succeeds with one method but fails with the other, depending on the evaluation context and the variant returned. This inconsistent behavior is further detailed in bug report #1481.

The root cause of this issue is the absence of a dedicated, authoritative type definition for the flag itself. Instead, the type is inferred from the `defaultVariant` or API itself (`T` from `Resolve<T>`) , which is not always a reliable source of truth for all possible variants. This can lead to unexpected errors and make it difficult for developers to debug their feature flags.


## Requirements

* The new `type` field in the flag definition must be optional to ensure backward compatibility.
* If the `type` field is present, `flagd` must validate that all variants of the flag conform to this type during initialization.
* Type mismatches found during initialization must be reported as errors.
* The `Resolve<T>` and `ResolveAll` methods must use the `type` field for validation when it is available.
* The implementation must be consistent with the OpenFeature specification and the flag manifest schema.


## Considered Options

* **Consistent `defaultVariant` Validation:** Align the behavior of `Resolve<T>` with `ResolveAll` by making `Resolve<T>` validate the evaluated variant against the type of the `defaultVariant`.
* **API Extension with Explicit Flag Type:** Introduce an optional `type` property to the flag definition to serve as the authoritative source for type validation.


## Proposal

This proposal is to extend the flag definition with an optional `type` property. This approach is chosen over simply aligning the `Resolve<T>` and `ResolveAll` validation because it addresses the root cause of the type inconsistency and provides a more robust, long-term solution.

By introducing an explicit `type` field, it establishes a single source of truth for the flag's type, independent of its variants. This allows for early and consistent type validation during flag definition parsing, preventing type-related errors at runtime.

The new `type` field will be optional to maintain backward compatibility with existing flag configurations. If the field is omitted, `flagd` will treat the flag as having `object`, and no type validation will be performed against the `defaultVariant`. When the `type` field is present, `flagd` will enforce that all variants of the flag conform to the specified type.

This change will make the behavior of `flagd` more predictable and reliable.


### API changes

The `flagd` flag definition will be updated to include an optional `type` property. This property will be a string enum with the following possible values: `"boolean"`, `"string"`, `"number"`, and `"object"`.


#### JSON Schema

The following changes will be made to the `schemas/json/flags.json` file:

1.  A new `type` property will be added to the `flag` definition:

```json
"flag": {
  "type": "object",
  "properties": {
    "type": {
      "title": "Flag Type",
      "description": "The type of the flag. If specified, all variants must conform to this type.",
      "type": "string",
      "enum": [
        "boolean",
        "string",
        "number",
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

2.  The `booleanFlag`, `stringFlag`, `numberFlag`, and `objectFlag` definitions will be updated to enforce the `type` property:

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
        "type": {
          "const": "boolean"
        }
      }
    }
  ]
}
```

Similar changes will be made to `stringFlag`, `numberFlag`, and `objectFlag` to enforce their respective types.

### Consequences

#### The good
* It improves the reliability and predictability of flag evaluations.
* It allows for early error detection of type mismatches.
* It improves the developer experience by making the API more explicit.

#### The bad
* It adds a new field to the flag definition, which developers need to be aware of.
* It requires updating all `flagd` SDKs to support the new field.
* It requires updating flag manifest schema

### Timeline

* **Phase 1: Core Implementation**
    * Update the `flagd` core to support the new `type` field.
    * Implement the type validation logic.
    * Update the JSON schema.
    * Add unit and integration tests.
* **Phase 2: SDK Updates**
    * Update all `flagd` SDKs to support the new `type` field.
    * Update flag manifest
* **Phase 3: Documentation**
    * Update the `flagd` documentation to reflect the changes.



## More Information

* **Bug Report:** [https://github.com/open-feature/flagd/issues/1481](https://github.com/open-feature/flagd/issues/1481)
* **Flag schema** [https://flagd.dev/schema/v0/flags.json](https://flagd.dev/schema/v0/flags.json)
* **Flag Manifest Schema:** [https://raw.githubusercontent.com/open-feature/cli/refs/heads/main/schema/v0/flag-manifest.json](https://raw.githubusercontent.com/open-feature/cli/refs/heads/main/schema/v0/flag-manifest.json)