import type { Scenario } from "../types";
import { contextToPrettyJson, featureDefinitionToPrettyJson } from "../utils";

/*
SOLUTION:
Replace targeting: {} with:

targeting: {
  "if": [
    { "and": [
      { "==": [{ "var": "user.plan" }, "enterprise"] },
      { "==": [{ "var": "user.region" }, "eu"] }
    ]},
    "eu-enterprise",
    null
  ]
}
*/

export const challengeVar: Scenario = {
  description: [
    "🏆 CHALLENGE: Deep dive.",
    'The context contains a nested object: {"user": {"plan": "enterprise", "region": "eu"}}.',
    "Write a targeting rule that returns \"eu-enterprise\" when user.plan equals \"enterprise\" AND user.region equals \"eu\". Return null otherwise (to fall back to the default).",
    "Hint: use dot notation with var — { \"var\": \"user.plan\" } — and combine checks with \"and\".",
    "✅ Expected result: variant \"eu-enterprise\", value \"EU Enterprise Dashboard\".",
  ].join(" "),
  flagDefinition: featureDefinitionToPrettyJson({
    flags: {
      "dashboard-mode": {
        state: "ENABLED",
        defaultVariant: "standard",
        variants: {
          "eu-enterprise": "EU Enterprise Dashboard",
          standard: "Standard Dashboard",
        },
        targeting: {},
      },
    },
  }),
  flagKey: "dashboard-mode",
  returnType: "string",
  context: contextToPrettyJson({
    user: {
      plan: "enterprise",
      region: "eu",
    },
  }),
};
