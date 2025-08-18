<!-- DO NOT REMOVE THIS PART UNTIL COMMIT
	CONTEX:

	Bug #1481 https://github.com/open-feature/flagd/issues/1481
		The bug report highlights an inconsistency in error handling between the Resolve and ResolveAll methods in flagd when a feature flag's variants have different data types.

		Specifically, Resolve<T> (e.g., ResolveBoolean) validates the evaluated variant against the requested type T, whereas ResolveAll validates it against the type of the flag's defaultVariant.
		This leads to inconsistant results where one method can succeed and the other can fail for the same flag evaluation, depending on the context and which variant is selected. 

		Proposed Solutions.
		Two main solutions were discussed to address this inconsistency:

		1. 	Consistent defaultVariant Validation: The most immediate proposal is to align the behavior of Resolve<T> with ResolveAll. 
			This would involve modifying the Resolve<T> methods to also validate the type of the evaluated variant against the type of the defaultVariant.

		2. 	API Extension with Explicit Flag Type: A longer-term solution suggested is to enhance the flag schema by adding an optional type property to the flag definition itself. 
			This would allow for more reliable and explicit type validation during initialization, catching type mismatches early.


	The proposal is not to fix just a bug (#1) but extend API (#2) with a new field 'FlagType' and do all the validations against it instead of DefaultVariant.


	Example PRD 
		/flagd/blob/main/docs/architecture-decisions/duplicate-flag-keys.md 

	Important points:
		* 	We should think about the tasks required to support the new field in the SDKs of all languages.
			/flagd/docs/architecture-decisions/duplicate-flag-keys.md

		* 	Another thing to consider is to keep it consistent with the manifest schema 
			https://raw.githubusercontent.com/open-feature/cli/refs/heads/main/schema/v0/flag-manifest.json 
			Need to check if manifest has the same defs for types.
-->

---
# Valid statuses: draft | proposed | rejected | accepted | superseded
status: draft
author: @andreyturkov
created: 2025-08-14
updated: 2025-08-14
---

# Extending Flag Definition with a Type Property

## Background

Currently, `flagd` exhibits inconsistent behavior in type validation between its `Resolve<T>` and `ResolveAll` API methods. The `Resolve<T>` method validates the evaluated flag variant against the type `T` requested by the client, while `ResolveAll` validates it against the type of the `defaultVariant` specified in the flag's definition. This discrepancy can lead to situations where a flag evaluation succeeds with one method but fails with the other, depending on the evaluation context and the variant returned. This inconsistent behavior is further detailed in bug report #1481.

The root cause of this issue is the absence of a dedicated, authoritative type definition for the flag itself. Instead, the type is inferred from the `defaultVariant`, which is not always a reliable source of truth for all possible variants. This can lead to unexpected errors and make it difficult for developers to debug their feature flags.

To address this, we propose extending the flag definition with an optional `type` property. This would establish a clear and explicit type for each flag, ensuring that all variants conform to a single, authoritative type. This change aims to eliminate the current inconsistencies and provide a more robust and predictable type validation mechanism, ultimately improving the reliability and developer experience of `flagd`.

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

We propose to extend the flag definition with an optional `type` property. This approach is chosen over simply aligning the `Resolve<T>` and `ResolveAll` validation because it addresses the root cause of the type inconsistency and provides a more robust, long-term solution.

By introducing an explicit `type` field, we establish a single source of truth for the flag's type, independent of its variants. This allows for early and consistent type validation during flag definition parsing, preventing type-related errors at runtime.

The new `type` field will be optional to maintain backward compatibility with existing flag configurations. If the field is omitted, `flagd` will treat the flag as having `Object`, and no type validation will be performed against the `defaultVariant`. When the `type` field is present, `flagd` will enforce that all variants of the flag conform to the specified type.

This change will make the behavior of `flagd` more predictable and reliable, improving the overall developer experience. It also aligns `flagd` with best practices for API design, where explicit type definitions are preferred over implicit ones.

<!-- This is an optional element. Feel free to remove. -->
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

* Good, because it improves the reliability and predictability of flag evaluations.
* Good, because it allows for early error detection of type mismatches.
* Good, because it improves the developer experience by making the API more explicit.
* Bad, because it adds a new field to the flag definition, which developers need to be aware of.
* Bad, because it requires updating all `flagd` SDKs to support the new field.

### Timeline

* **Phase 1: Core Implementation (2-4 weeks)**
    *   Update the `flagd` core to support the new `type` field.
    *   Implement the type validation logic.
    *   Update the JSON schema.
    *   Add unit and integration tests.
* **Phase 2: SDK Updates (4-8 weeks)**
    *   Update all official `flagd` SDKs to support the new `type` field.
    *   This will involve coordination with the maintainers of each SDK.
* **Phase 3: Documentation and Communication (1-2 weeks)**
    *   Update the `flagd` documentation to reflect the changes.
    *   Communicate the changes to the community through blog posts, and other channels.

<!-- This is an optional element. Feel free to remove. -->
### Open questions

* How should `flagd` behave if the `type` field is specified but the `defaultVariant` has a different type? 
* What is the exact mapping between the `type` field and the protobuf `oneof` value types?
* Should we add a new `type` field to the `AnyFlag` message in the protobuf definition?
* How will this change be communicated to the community and what is the migration path for users with existing flag configurations?

<!-- This is an optional element. Feel free to remove. -->
## More Information

* **Bug Report:** [https://github.com/open-feature/flagd/issues/1481](https://github.com/open-feature/flagd/issues/1481)
* **OpenFeature Specification:** [https://openfeature.dev/docs/specification/](https://openfeature.dev/docs/specification/)
* **Flag Manifest Schema:** [https://raw.githubusercontent.com/open-feature/cli/refs/heads/main/schema/v0/flag-manifest.json](https://raw.githubusercontent.com/open-feature/cli/refs/heads/main/schema/v0/flag-manifest.json)