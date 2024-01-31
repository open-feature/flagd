import type { Scenario } from "../types";
import { contextToPrettyJson, featureDefinitionToPrettyJson } from "../utils";

export const sharedEvaluators: Scenario = {
  description: [
    "In this scenario, we have two feature flags that share targeting rule logic.",
    "This is accomplished by defining a $evaluators object in the feature flag definition and referencing it by name in a targeting rule.",
    "Experiment with changing the email domain in the shared evaluator.",
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
          if: [{ $ref: "emailWithFaas" }, "true"],
        },
      },
      "feature-2": {
        state: "ENABLED",
        defaultVariant: "false",
        variants: {
          true: true,
          false: false,
        },
        targeting: {
          if: [{ $ref: "emailWithFaas" }, "true"],
        },
      },
    },
    $evaluators: {
      emailWithFaas: {
        ends_with: [{ var: "email" }, "@faas.com"],
      },
    },
  }),
  flagKey: "feature-1",
  returnType: "boolean",
  context: contextToPrettyJson({
    email: "example@faas.com",
  }),
};
