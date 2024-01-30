import { Scenario } from "../types";
import { basicBoolean } from "./basic-boolean";
import { basicNumber } from "./basic-number";
import { basicObject } from "./basic-object";
import { basicString } from "./basic-string";
import { booleanShorthand } from "./boolean-shorthand";
import { enableByDomain } from "./enable-by-domain";
import { enableByLocale } from "./enable-by-locale";
import { enableByTime } from "./enable-by-time";
import { enableByVersion } from "./enable-by-version";
import { pseudoRandomSplit } from "./fraction-string";
import { progressRollout } from "./progressive-rollout";
import { sharedEvaluator } from "./share-evaluators";

export const scenarios = {
  "Basic boolean flag": basicBoolean,
  "Basic numeric flag": basicNumber,
  "Basic string flag": basicString,
  "Basic object flag": basicObject,
  "Enable for a specific email domain": enableByDomain,
  "Enable based on users locale": enableByLocale,
  "Enable based on release version": enableByVersion,
  "Enable based on the current time": enableByTime,
  "Multi-variant experiment": pseudoRandomSplit,
  "Progressive rollout": progressRollout,
  "Shared evaluators": sharedEvaluator,
  "Boolean shorthand": booleanShorthand,
} satisfies { [name: string]: Scenario };

export type ScenarioName = keyof typeof scenarios;
