import type { Scenario } from "../types";
import { contextToPrettyJson, featureDefinitionToPrettyJson } from "../utils";

export const enableByDomain: Scenario = {
  description: [
    'In this scenario, we have a feature flag with the key "enable-mainframe-access" that is enabled and has two variants: true and false.',
    'This flag has a targeting rule defined that enables the flag for users with an email address that ends with "@ingen.com".',
    "Try changing the email address in the context to something else and see what happens.",
  ].join(" "),
  flagDefinition: featureDefinitionToPrettyJson({
    flags: {
      "enable-mainframe-access": {
        state: "ENABLED",
        defaultVariant: "true",
        variants: {
          true: true,
          false: false,
        },
        targeting: {
          if: [{ ends_with: [{ var: "email" }, "@ingen.com"] }, "true"],
        },
      },
    },
  }),
  flagKey: "enable-mainframe-access",
  returnType: "boolean",
  context: contextToPrettyJson({
    email: "john.arnold@ingen.com",
  }),
};
