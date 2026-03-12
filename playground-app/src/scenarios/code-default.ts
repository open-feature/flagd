import type { Scenario } from "../types";
import { contextToPrettyJson, featureDefinitionToPrettyJson } from "../utils";

export const codeDefault: Scenario = {
  description: [
    "This scenario demonstrates code defaults in flagd. When defaultVariant is omitted or set to null, flagd falls back to the code-defined default value.",
    "In this example, the flag has two variants (on and off) with defaultVariant set to null.",
    "When evaluated, flagd returns reason=DEFAULT and omits the value and variant fields, allowing the client to use its own code default.",
    "Note: Omitting defaultVariant entirely has the same effect as setting it to null - both trigger code default behavior.",
    "Compare this with the basic boolean flag which has an explicit defaultVariant set to false.",
  ].join(" "),
  flagDefinition: featureDefinitionToPrettyJson({
    flags: {
      "code-default-flag": {
        state: "ENABLED",
        variants: {
          on: true,
          off: false,
        },
        defaultVariant: null,
      },
    },
  }),
  flagKey: "code-default-flag",
  returnType: "boolean",
  context: contextToPrettyJson({}),
  codeDefault: "true",
};
