---
status: proposed
author: @tangenti
created: 2025-06-27
updated: 2025-06-27
---

# Support for Duplicate Flag Keys

This ADR proposes allowing a single sync source to provide multiple flags that share the same key. This enables greater flexibility for modularizing flag configurations.

## Background

Currently, the `flagd` [flag configuration](https://flagd.dev/schema/v0/flags.json) stores flags in a JSON object (a map), where each key must be unique. While the JSON specification technically allows duplicate keys, it's not recommended and not well-supported in the implementations.

This limitation prevents use cases for flag modularization and multi-tenancy, such as:

- **Component-based Flags:** Two different services, each with its own in-process provider, cannot independently define a flag with the same key when communicating with the same `flagd` daemon.
- **Multi-Tenant Targeting:** A single flagd daemon uses the same flag key with different targeting rules for different tenants

## Requirements

- Allow a single sync source to define multiple flags that have the same key.
- Flags from a sync source with the same keys can have different types and targeting rules.
- No breaking changes for the current flagd flag configuration schema or flagd sync services.

## Proposal

We will update the `flagd` flag configuration schema to support receiving flags as an **array of flag objects**. The existing schema will remain fully supported.

### API Change

#### Flag Configuration Schema

We'll add a new schema as a [subschema](https://json-schema.org/learn/glossary#subschema) to the existing flagd flag configuration schema. It will be a composite of the original schema except `flags` (`#/definitions/base`), with a new schema for `flags` that allows flags array in addition to the currently supported flags object. The existing main schema will be the composite of the

```json
...
"flagsArray": {
    "type": "array",
    "items": {
        "allOf": [
            {
                "$ref": "#/definitions/flag"
            },
            {
                "type": "object",
                "properties": {
                    "key": {
                        "description": "Key of the flag",
                        "type": "string",
                        "minLength": 1
                    }
                },
                "required": [
                    "key"
                ]
            }
        ]
    }
},
"flagsArraySchema": {
    "$id": "https://flagd.dev/schema/v0/flags.json#flagsarray",
    "type": "object",
    "allOf": [
        {
            "$ref": "#/definitions/base"
        },
        {
            "properties": {
                "flags": {
                    "oneOf": [
                        {
                            "$ref": "#/definitions/flagsArray"
                        },
                        {
                            "$ref": "#/definitions/flagsMap"
                        }
                    ]
                }
            },
            "required": [
                "flags"
            ]
        }
    ]
}
...
```

If the config level flag set ID is not specified, `metadata.flagSetID` of each flag will be interpreted as its flag set ID.

A flag will be uniquely identified by the composite key `(flagKey, flagSetID)`. The following three flags will be considered as three different flags.

1. `{"flagKey": "enable-feature", "flagSetID": ""}`
2. `{"flagKey": "enable-feature", "flagSetID": "default"}`
3. `{"flagKey": "enable-feature", "flagSetID": "beta"}`

### Flagd daemon

Flagd daemon will perform the JSON schema checks with the reference to `https://flagd.dev/schema/v0/flags.json#flagsarray`, allowing both flags as an object and as an array.

If the flag array contains two or more flags with the same composite key, the config will be considered as invalid.

If the request from in-process flagd providers result in a config that has duplicate flag keys, the flagd daemon will only keep one of them in the response.

### Flagd Daemon Storage

1. Flagd will have separate stores for `flags` and `sources`.

1. The `flags` store will use the composite key for flags.

1. `selector` will be removed from the store

1. `flagSetID` will be moved from `source` metadata to `flag` metadata.

### Flags Lifecycle

Currently, the flags configurations from the latest update of a source will override the existing ones. If a flag was presented in the previous configuration but not in the current configuration, it will **NOT** get removed. Flags removals can only be triggerred when a source is removed, or a full resync is triggered.

We'll keep the same behaviors with this proposal:

1. If two sources provide the flags with the same composite key, the latest one will be stored.

1. If a flag from a source no longer presents in the latest configuration of the same source, it will be kept.

This behavior is not ideal and should be addressed in a separate ADR.

### Consequences

#### The good

- One source can provide flags with the same keys.
- Flag set ID no longer bound to a source, so one source can have multiple flag sets.
- No breaking change of the API definition and the API behaviors.
- No significant change on the flagd stores and how selections work.

#### The bad

- The proposal still leverages the concept of flag set in the flagd storage.

- The schema does not guarantee that flags of the same flag set from the same source will not have the same keys. This is guaranteed in the proposal of #1634.

- The flag array is less readable compared to the flag sets object proposed in #1634.
