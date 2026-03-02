import type { Scenario } from "../types";
import { contextToPrettyJson, featureDefinitionToPrettyJson } from "../utils";

/*
SOLUTION:
Replace targeting: {} with:

targeting: {
  "fractional": [
    { "cat": [
      { "var": "$flagd.flagKey" },
      { "var": "email" }
    ]},
    ["control", 80],
    ["experiment", 20]
  ]
}
*/

export const challengeFractional: Scenario = {
  description: [
    "🏆 CHALLENGE: A/B experiment.",
    'The context contains {"targetingKey": "user-42", "email": "tester@example.com"}.',
    "Write a targeting rule using the \"fractional\" operator that splits users 80/20 between \"control\" and \"experiment\".",
    "Use a bucketing value that combines $flagd.flagKey and email with \"cat\".",
    "Hint: { \"fractional\": [{ \"cat\": [{ \"var\": \"$flagd.flagKey\" }, { \"var\": \"email\" }] }, [\"control\", 80], [\"experiment\", 20]] }.",
    "✅ Expected result: one of the two variants (the split is deterministic for a given email).",
    "Try different email addresses to see some users get \"experiment\".",
  ].join(" "),
  flagDefinition: featureDefinitionToPrettyJson({
    flags: {
      "checkout-experiment": {
        state: "ENABLED",
        defaultVariant: "control",
        variants: {
          control: "Classic Checkout",
          experiment: "New Checkout",
        },
        targeting: {},
      },
    },
  }),
  flagKey: "checkout-experiment",
  returnType: "string",
  context: contextToPrettyJson({
    targetingKey: "user-42",
    email: "tester@example.com",
  }),
};
