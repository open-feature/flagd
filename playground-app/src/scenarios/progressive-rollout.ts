import type { Scenario } from "../types";
import { contextToPrettyJson, featureDefinitionToPrettyJson } from "../utils";

export const progressRollout: Scenario = {
  description: [
    'In this scenario, we have a feature flag with the key "enable-new-llm-model" with multiple variant for illustrative purposes.',
    "This flag has a targeting rule defined that enables the flag for a percentage of users based on the release phase.",
    'The "targetingKey" ensures that the user always sees the same results during a each phase of the rollout process.',
  ].join(" "),
  flagDefinition: () => {
    const phase1 = Math.floor(Date.now() / 1000) + 5;
    const phase2 = Math.floor(Date.now() / 1000) + 10;
    const phase3 = Math.floor(Date.now() / 1000) + 15;
    const enabled = Math.floor(Date.now() / 1000) + 20;
    return featureDefinitionToPrettyJson({
      flags: {
        "enable-new-llm-model": {
          state: "ENABLED",
          defaultVariant: "disabled",
          variants: {
            disabled: false,
            phase1Enabled: true,
            phase1Disabled: false,
            phase2Enabled: true,
            phase2Disabled: false,
            phase3Enabled: true,
            phase3Disabled: false,
            enabled: true,
          },
          targeting: {
            if: [
              {
                ">": [{ var: "$flagd.timestamp" }, enabled],
              },
              "enabled",
              {
                if: [
                  { ">": [{ var: "$flagd.timestamp" }, phase3] },
                  {
                    fractional: [
                      ["phase3Enabled", 50],
                      ["phase3Disabled", 50],
                    ],
                  },
                  {
                    if: [
                      { ">": [{ var: "$flagd.timestamp" }, phase2] },
                      {
                        fractional: [
                          ["phase2Enabled", 25],
                          ["phase2Disabled", 75],
                        ],
                      },
                      { ">": [{ var: "$flagd.timestamp" }, phase1] },
                      {
                        fractional: [
                          ["phase1Enabled", 10],
                          ["phase1Disabled", 90],
                        ],
                      },
                    ],
                  },
                ],
              },
            ],
          },
        },
      },
    });
  },
  flagKey: "enable-new-llm-model",
  returnType: "boolean",
  context: () =>
    contextToPrettyJson({
      targetingKey: "sessionId-12345",
    }),
};
