import type { Scenario } from "../types";
import { contextToPrettyJson, featureDefinitionToPrettyJson } from "../utils";

export const enableByVersion: Scenario = {
  description: [
    'In this scenario, we have a feature flag with the key "enable-performance-mode" that is enabled and has two variants: on and off.',
    'This rule looks for the evaluation context "version". If the version is greater or equal to "1.7.0" the feature is enabled.',
    'Otherwise, the "defaultVariant" is return. Experiment by changing the version in the context.',
  ].join(" "),
  flagDefinition: featureDefinitionToPrettyJson({
    flags: {
      "enable-performance-mode": {
        state: "ENABLED",
        defaultVariant: "off",
        variants: {
          on: true,
          off: false,
        },
        targeting: {
          if: [{ sem_ver: [{ var: "version" }, ">=", "1.7.0"] }, "on"],
        },
      },
    },
  }),
  flagKey: "enable-performance-mode",
  returnType: "boolean",
  context: contextToPrettyJson({
    version: "1.6.0",
  }),
  codeDefault: "false",
};
