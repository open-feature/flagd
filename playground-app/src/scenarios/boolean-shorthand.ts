import type { Scenario } from "../types";
import { contextToPrettyJson, featureDefinitionToPrettyJson } from "../utils";

export const booleanShorthand: Scenario = {
  description: [
    "In this scenario, we have a feature flag with a targeting rule that returns true when the age is 18 or greater.",
    "This targeting rule leverages the boolean shorthand syntax, which converts a boolean to its string equivalent.",
    "The converted value is then used as the variant key.",
    "Try changing the value of the context attribute 'age'.",
  ].join(" "),
  flagDefinition: featureDefinitionToPrettyJson({
    flags: {
      "feature-1": {
        state: "ENABLED",
        defaultVariant: "false",
        variants: {
          true: true,
          false: false,
        },
        targeting: {
          ">=": [{ var: "age" }, 18],
        },
      },
    },
  }),
  flagKey: "feature-1",
  returnType: "boolean",
  context: contextToPrettyJson({
    age: 20,
  }),
};
