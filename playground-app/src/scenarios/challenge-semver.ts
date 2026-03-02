import type { Scenario } from "../types";
import { contextToPrettyJson, featureDefinitionToPrettyJson } from "../utils";

/*
SOLUTION:
Replace targeting: {} with:

targeting: {
  "if": [
    { "sem_ver": [{ "var": "appVersion" }, ">=", "2.0.0"] },
    "modern",
    null
  ]
}
*/

export const challengeSemver: Scenario = {
  description: [
    "🏆 CHALLENGE: Version gate.",
    'The context contains {"appVersion": "2.3.1"}.',
    "Write a targeting rule that returns \"modern\" if appVersion >= \"2.0.0\", and null otherwise (falling back to the default \"legacy\" variant).",
    "Use the \"sem_ver\" custom operator: { \"sem_ver\": [{ \"var\": \"appVersion\" }, \">=\", \"2.0.0\"] }.",
    "✅ Expected result: variant \"modern\", value \"Modern UI\".",
    "Try changing appVersion to \"1.9.9\" — it should fall back to \"legacy\".",
  ].join(" "),
  flagDefinition: featureDefinitionToPrettyJson({
    flags: {
      "ui-generation": {
        state: "ENABLED",
        defaultVariant: "legacy",
        variants: {
          modern: "Modern UI",
          legacy: "Legacy UI",
        },
        targeting: {},
      },
    },
  }),
  flagKey: "ui-generation",
  returnType: "string",
  context: contextToPrettyJson({
    appVersion: "2.3.1",
  }),
};
