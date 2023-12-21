import type { FlagValueType, JsonObject } from "@openfeature/core";

type StringVariants = {
  [key: string]: string;
};

type NumberVariants = {
  [key: string]: number;
};

type BooleanVariants = {
  [key: string]: boolean;
};

type ObjectVariants = {
  [key: string]: JsonObject;
};

export type FeatureDefinition = {
  flags: {
    [key: string]: {
      state: "ENABLED" | "DISABLED";
      defaultVariant: string;
      variants:
        | StringVariants
        | NumberVariants
        | BooleanVariants
        | ObjectVariants;
      targeting?: JsonObject;
    };
  };
  $evaluators?: JsonObject;
};

export type Scenario = {
  /**
   * A description of the scenario.
   */
  description: string;
  /**
   * A stringify version of the flag definition.
   */
  flagDefinition: string | (() => string);
  /**
   * The flag key that should be used as the default value in the playground.
   */
  flagKey: string;
  /**
   * The expected return type of the flag.
   */
  returnType: FlagValueType;
  /**
   * A string or function that returns a string that represents evaluation context.
   */
  context: string | (() => string);
};
