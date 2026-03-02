import type { Scenario } from "../types";
import { contextToPrettyJson, featureDefinitionToPrettyJson } from "../utils";

/*
SOLUTION:
Replace targeting: {} with:

targeting: {
  "if": [
    { "ends_with": [{ "var": "email" }, ".vip"] },
    "vip",
    "standard"
  ]
}
*/

export const challengeStringMatch: Scenario = {
  description: [
    "🏆 CHALLENGE: VIP access.",
    'The context contains {"email": "ceo@bigcorp.vip"}.',
    "Write a targeting rule that returns \"vip\" if the email ends with \".vip\", and \"standard\" otherwise.",
    "Use the custom \"ends_with\" operator: { \"ends_with\": [{ \"var\": \"email\" }, \".vip\"] }.",
    "✅ Expected result: variant \"vip\", value true.",
    "Try changing the email to \"user@example.com\" — it should return \"standard\" (false).",
  ].join(" "),
  flagDefinition: featureDefinitionToPrettyJson({
    flags: {
      "vip-access": {
        state: "ENABLED",
        defaultVariant: "standard",
        variants: {
          vip: true,
          standard: false,
        },
        targeting: {},
      },
    },
  }),
  flagKey: "vip-access",
  returnType: "boolean",
  context: contextToPrettyJson({
    email: "ceo@bigcorp.vip",
  }),
};
