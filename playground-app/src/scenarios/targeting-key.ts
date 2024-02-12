import type { Scenario } from "../types";
import { contextToPrettyJson, featureDefinitionToPrettyJson } from "../utils";

export const targetingKey: Scenario = {
  description: [
    "In this scenario, we have a feature flag that is evaluated based on its targeting key.",
    "The targeting key is contain a string uniquely identifying the subject of the flag evaluation, such as a user's email, or a session identifier.",
    "In this case, null is returned from targeting if the targeting key doesn't match; this results in a reason of \"DEFAULT\", since no variant was matched by the targeting rule.",
  ].join(" "),
  flagDefinition: featureDefinitionToPrettyJson({
    flags: {
      "targeting-key-flag": {
        state: "ENABLED",
        variants: {
          miss: "miss",
          hit: "hit"
        },
        defaultVariant: "miss",
        targeting: {
          if: [
            {
              "==": [ { var: "targetingKey" }, "5c3d8535-f81a-4478-a6d3-afaa4d51199e" ]
            },
            "hit",
            null
          ]
        }
      }
    },
  }),
  flagKey: "targeting-key-flag",
  returnType: "string",
  context: contextToPrettyJson({
    targetingKey: "5c3d8535-f81a-4478-a6d3-afaa4d51199e",
  }),
};
