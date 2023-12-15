import { Scenario } from "../types";
import { contextToPrettyJson, featureDefinitionToPrettyJson } from "../utils";

export const basicNumber: Scenario = {
  description: [
    'In this scenario, we have a feature flag with the key "basic-number" that is enabled and has two variants: 1 and 2.',
    'The default variant is 1. Try changing the "defaultVariant" to "2" or add a targeting rule.',
  ].join(" "),
  flagDefinition: featureDefinitionToPrettyJson({
    flags: {
      "basic-number": {
        state: "ENABLED",
        defaultVariant: "1",
        variants: {
          "1": 1,
          "2": 2,
        },
        targeting: {},
      },
    },
  }),
  flagKey: "basic-number",
  returnType: "number",
  context: contextToPrettyJson({}),
};
