import type { Scenario } from "../types";
import { contextToPrettyJson, featureDefinitionToPrettyJson } from "../utils";

export const pseudoRandomSplit: Scenario = {
  description: [
    'In this scenario, we have a feature flag with the key "pseudo-random-split" that is enabled and has four variants: red, blue, green, and grey.',
    'The targeting rule uses the "fractional" operator, which deterministically splits the traffic based on the configuration.',
    "This configuration splits the traffic evenly between the four variants but notice that the targeting rule is based on the sessionId, which means that the same user will always get the same variant.",
    'Try changing the "targetingKey" to another value and see what happens.',
  ].join(" "),
  flagDefinition: featureDefinitionToPrettyJson({
    flags: {
      "color-palette-experiment": {
        state: "ENABLED",
        defaultVariant: "grey",
        variants: {
          red: "#b91c1c",
          blue: "#0284c7",
          green: "#16a34a",
          grey: "#4b5563",
        },
        targeting: {
          fractional: [
            ["red", 25],
            ["blue", 25],
            ["green", 25],
            ["grey", 25],
          ],
        },
      },
    },
  }),
  flagKey: "color-palette-experiment",
  returnType: "string",
  context: contextToPrettyJson({
    targetingKey: "sessionId-123",
  }),
};
