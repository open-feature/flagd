import type { Scenario } from "../types";
import { contextToPrettyJson, featureDefinitionToPrettyJson } from "../utils";

export const basicString: Scenario = {
  description: [
    'In this scenario, we have a feature flag with the key "basic-string" that is enabled and has two variants: foo and bar.',
    'The default variant is foo. Try changing the "defaultVariant" to "bar" or add a targeting rule.',
  ].join(" "),
  flagDefinition: featureDefinitionToPrettyJson({
    flags: {
      "basic-string": {
        state: "ENABLED",
        defaultVariant: "foo",
        variants: {
          foo: "foo",
          bar: "bar",
        },
        targeting: {},
      },
    },
  }),
  flagKey: "basic-string",
  returnType: "string",
  context: contextToPrettyJson({}),
};
