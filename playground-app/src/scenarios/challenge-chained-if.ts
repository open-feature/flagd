import type { Scenario } from "../types";
import { contextToPrettyJson, featureDefinitionToPrettyJson } from "../utils";

/*
SOLUTION:
Replace targeting: {} with:

targeting: {
  "if": [
    { "<=": [{ "var": "speed" }, 30] },
    "green",
    { "<=": [{ "var": "speed" }, 60] },
    "yellow",
    "red"
  ]
}
*/

export const challengeChainedIf: Scenario = {
  description: [
    "🏆 CHALLENGE: Traffic light.",
    'The context contains {"speed": 45}.',
    "Write a targeting rule using chained \"if\" that returns: \"green\" if speed <= 30, \"yellow\" if speed <= 60, \"red\" otherwise.",
    "Remember chained if: { \"if\": [cond1, val1, cond2, val2, fallback] }.",
    "✅ Expected result with speed 45: variant \"yellow\", value \"Slow down\".",
    "Try speed 25 (→ \"green\") and speed 80 (→ \"red\") to verify.",
  ].join(" "),
  flagDefinition: featureDefinitionToPrettyJson({
    flags: {
      "speed-warning": {
        state: "ENABLED",
        defaultVariant: "green",
        variants: {
          green: "All clear",
          yellow: "Slow down",
          red: "Too fast!",
        },
        targeting: {},
      },
    },
  }),
  flagKey: "speed-warning",
  returnType: "string",
  context: contextToPrettyJson({
    speed: 45,
  }),
};
