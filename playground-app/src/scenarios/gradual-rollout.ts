import type { Scenario } from "../types";
import { contextToPrettyJson, featureDefinitionToPrettyJson } from "../utils";

export const gradualRollout: Scenario = {
  description: [
    "This scenario demonstrates a gradual rollout using dynamic weights computed from timestamps.",
    "The feature flag 'gradual-rollout-feature' uses fractional evaluation with weights that change over time.",
    "The 'on' weight increases as time passes (current timestamp minus start time).",
    "The 'off' weight decreases as time passes (end time minus current timestamp).",
    "This creates a smooth linear transition of users from 0% to 100% rollout over 30 seconds.",
    "For example, user-1 will transition to true at ~t=+6s, while user-2 will transition at ~t=+13s",
    "By t=+30s, all users will see the flag as true, no matter their targeting key.",
  ].join(" "),
  flagDefinition: () => {
    const now = Math.floor(Date.now() / 1000);
    const startTime = now; // Start immediately
    const endTime = now + 30; // End 30 seconds from now

    return featureDefinitionToPrettyJson({
      flags: {
        "gradual-rollout-feature": {
          state: "ENABLED",
          defaultVariant: "off",
          variants: {
            on: true,
            off: false,
          },
          targeting: {
            // Gradual rollout: weights change based on current timestamp
            "fractional": [
              [
                "on",
                { "-": [{ "var": "$flagd.timestamp" }, startTime] },
              ],
              [
                "off",
                { "-": [endTime, { "var": "$flagd.timestamp" }] },
              ],
            ],
          },
        },
      },
    });
  },
  flagKey: "gradual-rollout-feature",
  returnType: "boolean",
  context: () =>
    contextToPrettyJson({
      targetingKey: "user-1",
    }),
  codeDefault: "false",
};
