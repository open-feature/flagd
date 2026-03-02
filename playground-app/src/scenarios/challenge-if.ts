import type { Scenario } from "../types";
import { contextToPrettyJson, featureDefinitionToPrettyJson } from "../utils";

/*
SOLUTION:
Replace targeting: {} with:

targeting: {
  "if": [
    { ">=": [{ "var": "age" }, 18] },
    "allowed",
    "denied"
  ]
}
*/

export const challengeIf: Scenario = {
  description: [
    "🏆 CHALLENGE: The bouncer.",
    'The context contains {"age": 21}.',
    "Write a targeting rule using \"if\" that returns \"allowed\" when age >= 18, and \"denied\" otherwise.",
    "Remember: { \"if\": [condition, value-if-true, value-if-false] }.",
    "✅ Expected result: variant \"allowed\", value \"Welcome in!\".",
    "Then try changing the age to 16 — it should return \"denied\".",
  ].join(" "),
  flagDefinition: featureDefinitionToPrettyJson({
    flags: {
      "door-policy": {
        state: "ENABLED",
        defaultVariant: "denied",
        variants: {
          allowed: "Welcome in!",
          denied: "Maybe next year.",
        },
        targeting: {},
      },
    },
  }),
  flagKey: "door-policy",
  returnType: "string",
  context: contextToPrettyJson({
    age: 21,
  }),
};
