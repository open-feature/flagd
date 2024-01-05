import type { Scenario } from "../types";
import { contextToPrettyJson, featureDefinitionToPrettyJson } from "../utils";

export const enableByTime: Scenario = {
  description: [
    'In this scenario, we have a feature flag with the key "enable-announcement-banner" that is enabled and has two variants: true and false.',
    "This flag has a targeting rule defined that enables the flag after a specified time.",
    'The current time (epoch) can be accessed using "$flagd.timestamp" which is automatically provided by flagd.',
    'Five seconds after loading this scenario, the response will change to "true".',
  ].join(" "),
  flagDefinition: () =>
    featureDefinitionToPrettyJson({
      flags: {
        "enable-announcement-banner": {
          state: "ENABLED",
          defaultVariant: "false",
          variants: {
            true: true,
            false: false,
          },
          targeting: {
            if: [
              {
                ">": [
                  { var: "$flagd.timestamp" },
                  Math.floor(Date.now() / 1000) + 5,
                ],
              },
              "true",
            ],
          },
        },
      },
    }),
  flagKey: "enable-announcement-banner",
  returnType: "boolean",
  context: () => contextToPrettyJson({}),
};
