import { FeatureDefinition } from "./types";
import type { EvaluationContext } from "@openfeature/core";

export function prettyPrintJson(json: string): string {
  return JSON.stringify(JSON.parse(json), null, 2);
}

export function featureDefinitionToPrettyJson(
  definition: FeatureDefinition
): string {
  return prettyPrintJson(JSON.stringify(definition));
}

export function contextToPrettyJson(context: EvaluationContext) {
  return prettyPrintJson(JSON.stringify(context));
}

/**
 * Returns a string from a string or a function that returns a string.
 */
export function getString(input: string | (() => string)): string {
  if (typeof input === "function") {
    return input();
  }
  return input;
}

export function isValidJson(input: string) {
  try {
    JSON.parse(input);
    return true;
  } catch {
    return false;
  }
}