import type { Scenario } from "../types";
import { contextToPrettyJson, featureDefinitionToPrettyJson } from "../utils";

/*
SOLUTION:
Replace targeting: {} with:

targeting: {
  "if": [
    { "==": [{ "var": "tier" }, "premium"] },
    "premium-ui",
    "basic-ui"
  ]
}
*/

export const challengeShape: Scenario = {
  description: [
    "🏆 CHALLENGE: Write your first operation!",
    'The context contains {"tier": "premium"}.',
    "Write a targeting rule that returns variant \"premium-ui\" when tier equals \"premium\", and \"basic-ui\" otherwise.",
    "Hint: every JsonLogic operation is shaped like { \"operator\": [param1, param2] }.",
    "You'll need \"if\", \"==\", and \"var\".",
    "✅ Expected result: variant \"premium-ui\", value \"Premium Experience\".",
  ].join(" "),
  flagDefinition: featureDefinitionToPrettyJson({
    flags: {
      "ui-mode": {
        state: "ENABLED",
        defaultVariant: "basic-ui",
        variants: {
          "premium-ui": "Premium Experience",
          "basic-ui": "Basic Experience",
        },
        targeting: {},
      },
    },
  }),
  flagKey: "ui-mode",
  returnType: "string",
  context: contextToPrettyJson({
    tier: "premium",
  }),
};
