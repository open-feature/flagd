import type { Scenario } from "../types";
import { contextToPrettyJson, featureDefinitionToPrettyJson } from "../utils";

export const flagMetadata: Scenario = {
  description: [
    "In this scenario, we have a feature flag with metadata about the flag.",
    "There is top-level metadata for the flag set and metadata specific to the flag.",
    "These values are merged together, with the flag metadata taking precedence.",
  ].join(" "),
  flagDefinition: featureDefinitionToPrettyJson({
    flags: {
      "flag-with-metadata": {
        state: "ENABLED",
        variants: {
          on: true,
          off: false,
        },
        defaultVariant: "on",
        metadata: {
          version: "1",
        },
      },
    },
    metadata: {
      flagSetId: "playground/dev",
    },
  }),
  flagKey: "flag-with-metadata",
  returnType: "boolean",
  context: contextToPrettyJson({}),
};
