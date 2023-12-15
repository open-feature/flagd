import type { Scenario } from "../types";
import { contextToPrettyJson, featureDefinitionToPrettyJson } from "../utils";

export const enableByLocale: Scenario = {
  description: [
    'In this scenario, we have a feature flag with the key "supports-one-hour-delivery" that is enabled and has two variants: true and false.',
    'This flag has a targeting rule defined that enables the flag for users with a locale of "us" or "ca".',
    "Try changing the locale in the context to something else and see what happens.",
  ].join(" "),
  flagDefinition: featureDefinitionToPrettyJson({
    flags: {
      "supports-one-hour-delivery": {
        state: "ENABLED",
        defaultVariant: "false",
        variants: {
          true: true,
          false: false,
        },
        targeting: {
          if: [{ in: [{ var: "locale" }, ["us", "ca"]] }, "true"],
        },
      },
    },
  }),
  context: contextToPrettyJson({
    locale: "us",
  }),
  flagKey: "supports-one-hour-delivery",
  returnType: "boolean",
};
