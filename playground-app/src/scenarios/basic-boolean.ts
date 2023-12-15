import type { Scenario } from "../types";
import { contextToPrettyJson, featureDefinitionToPrettyJson } from "../utils";

export const basicBoolean: Scenario = {
  description: [
    "In this scenario, we have a feature flag with the key 'basic-boolean' that is enabled and has two variants: true and false.",
    "The default variant is false. Try changing the 'defaultVariant' to 'true' or add a targeting rule.",
  ].join(" "),
  flagDefinition: featureDefinitionToPrettyJson({
    flags: {
      "basic-boolean": {
        state: "ENABLED",
        defaultVariant: "false",
        variants: {
          true: true,
          false: false,
        },
      },
    },
  }),
  flagKey: "basic-boolean",
  returnType: "boolean",
  context: contextToPrettyJson({}),
};
