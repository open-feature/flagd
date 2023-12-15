import type { Scenario } from "../types";
import { contextToPrettyJson, featureDefinitionToPrettyJson } from "../utils";

export const basicObject: Scenario = {
  description: [
    'In this scenario, we have a feature flag with the key "basic-object" that is enabled and has two variants: foo and bar.',
    'The default variant is foo. Try changing the "defaultVariant" to "bar" or add a targeting rule.',
  ].join(" "),
  flagDefinition: featureDefinitionToPrettyJson({
    flags: {
      "basic-object": {
        state: "ENABLED",
        defaultVariant: "foo",
        variants: {
          foo: {
            foo: "foo",
          },
          bar: {
            bar: "bar",
          },
        },
      },
    },
  }),
  flagKey: "basic-object",
  returnType: "object",
  context: contextToPrettyJson({}),
};
