import { Scenario } from "../types";
import { basicBoolean } from "./basic-boolean";
import { basicNumber } from "./basic-number";
import { basicObject } from "./basic-object";
import { basicString } from "./basic-string";
import { booleanShorthand } from "./boolean-shorthand";
import { chainableConditions } from "./chainable-conditions";
import { enableByDomain } from "./enable-by-domain";
import { enableByLocale } from "./enable-by-locale";
import { enableByTime } from "./enable-by-time";
import { enableByVersion } from "./enable-by-version";
import { pseudoRandomSplit } from "./fraction-string";
import { progressRollout } from "./progressive-rollout";
import { sharedEvaluators } from "./share-evaluators";
import { targetingKey } from "./targeting-key";
import { flagMetadata } from "./flag-metadata";
import { challengeShape } from "./challenge-shape";
import { challengeIf } from "./challenge-if";
import { challengeChainedIf } from "./challenge-chained-if";
import { challengeVar } from "./challenge-var";
import { challengeStringMatch } from "./challenge-string-match";
import { challengeFractional } from "./challenge-fractional";
import { challengeSemver } from "./challenge-semver";
import { challengeTimestamp } from "./challenge-timestamp";
import { challengeReverse } from "./challenge-reverse";

export const scenarios = {
  "Basic boolean flag": basicBoolean,
  "Basic numeric flag": basicNumber,
  "Basic string flag": basicString,
  "Basic object flag": basicObject,
  "Enable for a specific email domain": enableByDomain,
  "Enable based on users locale": enableByLocale,
  "Enable based on release version": enableByVersion,
  "Enable based on the current time": enableByTime,
  "Chainable if/else/then": chainableConditions,
  "Multi-variant experiment": pseudoRandomSplit,
  "Progressive rollout": progressRollout,
  "Shared evaluators": sharedEvaluators,
  "Boolean variant shorthand": booleanShorthand,
  "Targeting key": targetingKey,
  "Flag metadata": flagMetadata,
  "CHALLENGE: Write your first operation": challengeShape,
  "CHALLENGE: The bouncer": challengeIf,
  "CHALLENGE: Traffic light": challengeChainedIf,
  "CHALLENGE: Deep dive": challengeVar,
  "CHALLENGE: VIP access": challengeStringMatch,
  "CHALLENGE: A/B experiment": challengeFractional,
  "CHALLENGE: Version gate": challengeSemver,
  "CHALLENGE: Launch day": challengeTimestamp,
  "CHALLENGE: Crack the context": challengeReverse,
} satisfies { [name: string]: Scenario };

export type ScenarioName = keyof typeof scenarios;
