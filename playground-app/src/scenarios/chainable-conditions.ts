import type { Scenario } from "../types";
import { contextToPrettyJson, featureDefinitionToPrettyJson } from "../utils";

export const chainableConditions: Scenario = {
  description: [
    "In this scenario, we have a feature flag with the key 'acceptable-feature-stability' with three variants: alpha, beta, and ga.",
    "The flag has a targeting rule that enables the flag based on the customer ID.",
    "The flag is enabled for customer-A in the alpha variant, for customer-B1 and customer-B2 in the beta variant, and for all other customers in the ga variant.",
    "Experiment by changing the 'customerId' in the context.",
  ].join(" "),
  flagDefinition: featureDefinitionToPrettyJson({
    flags: {
      "acceptable-feature-stability": {
        state: "ENABLED",
        defaultVariant: "ga",
        variants: {
          alpha: "alpha",
          beta: "beta",
          ga: "ga",
        },
        targeting: {
          if: [
            { "===": [{ var: "customerId" }, "customer-A"] },
            "alpha",
            { in: [{ var: "customerId" }, ["customer-B1", "customer-B2"]] },
            "beta",
            "ga",
          ],
        },
      },
    },
  }),
  flagKey: "acceptable-feature-stability",
  returnType: "string",
  context: contextToPrettyJson({
    targetingKey: "sessionId-123",
    customerId: "customer-A",
  }),
};
