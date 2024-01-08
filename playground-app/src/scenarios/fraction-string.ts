import type { Scenario } from "../types";
import { contextToPrettyJson, featureDefinitionToPrettyJson } from "../utils";

export const pseudoRandomSplit: Scenario = {
  description: [
    'In this scenario, we have a feature flag with the key "color-palette-experiment" that is enabled and has four variants: red, blue, green, and grey.',
    'The targeting rule uses the "fractional" operator, which deterministically splits the traffic based on the configuration.',
    'This configuration splits the traffic evenly between the four variants by bucketing evaluations pseudorandomly using the "targetingKey" and feature flag key.',
    'Experiment by changing the "targetingKey" to another value.',
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
