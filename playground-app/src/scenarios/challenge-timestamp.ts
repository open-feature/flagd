import type { Scenario } from "../types";
import { contextToPrettyJson, featureDefinitionToPrettyJson } from "../utils";

/*
SOLUTION:
Replace targeting: {} with:

targeting: {
  "if": [
    { ">=": [{ "var": "$flagd.timestamp" }, 1773792000] },
    "launched",
    null
  ]
}
*/

export const challengeTimestamp: Scenario = {
  description: [
    "🏆 CHALLENGE: Launch day.",
    "The flag below should show a \"launched\" banner, but only after March 15, 2026 00:00 UTC (timestamp: 1773792000).",
    "Before that date, it should return null (falling back to the \"coming-soon\" default).",
    "Write a targeting rule using \"if\", \">=\", \"var\" with \"$flagd.timestamp\", and the timestamp 1773792000.",
    "✅ Expected result RIGHT NOW (March 2, 2026): variant \"coming-soon\", value \"Stay tuned!\".",
    "Change the timestamp to a date in the past (e.g. 1704067200 = Jan 1, 2024) to verify it switches to \"launched\".",
  ].join(" "),
  flagDefinition: featureDefinitionToPrettyJson({
    flags: {
      "launch-banner": {
        state: "ENABLED",
        defaultVariant: "coming-soon",
        variants: {
          launched: "🚀 We're live!",
          "coming-soon": "Stay tuned!",
        },
        targeting: {},
      },
    },
  }),
  flagKey: "launch-banner",
  returnType: "string",
  context: contextToPrettyJson({}),
};
