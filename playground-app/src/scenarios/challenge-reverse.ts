import type { Scenario } from "../types";
import { contextToPrettyJson, featureDefinitionToPrettyJson } from "../utils";

/*
SOLUTION:
Replace the empty context with:

{
  "email": "anyone@premium.io",
  "appVersion": "3.0.0",
  "plan": "enterprise"
}
*/

export const challengeReverse: Scenario = {
  description: [
    "🏆 CHALLENGE: Crack the context!",
    "This one works in reverse — the flag definition and targeting rules are complete, but the context is empty.",
    "Read the targeting rule carefully and figure out what context values you need to provide to get the \"gold\" variant.",
    "The rule checks THREE conditions with \"and\": an email domain, a minimum version, and a plan type.",
    "✅ Expected result when you supply the right context: variant \"gold\", value \"🥇 Gold Tier\".",
    "If any condition fails, you'll get \"bronze\".",
  ].join(" "),
  flagDefinition: featureDefinitionToPrettyJson({
    flags: {
      "reward-tier": {
        state: "ENABLED",
        defaultVariant: "bronze",
        variants: {
          gold: "🥇 Gold Tier",
          bronze: "🥉 Bronze Tier",
        },
        targeting: {
          if: [
            {
              and: [
                { ends_with: [{ var: "email" }, "@premium.io"] },
                { sem_ver: [{ var: "appVersion" }, ">=", "3.0.0"] },
                { "==": [{ var: "plan" }, "enterprise"] },
              ],
            },
            "gold",
            "bronze",
          ],
        },
      },
    },
  }),
  flagKey: "reward-tier",
  returnType: "string",
  context: contextToPrettyJson({}),
};
