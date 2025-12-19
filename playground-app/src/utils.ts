import { FeatureDefinition } from "./types";
import type { EvaluationContext } from "@openfeature/core";
import yaml from "js-yaml";

const schemaMixin = {
  $schema: "https://flagd.dev/schema/v0/flags.json",
};

export function prettyPrintJson(json: string): string {
  return JSON.stringify(JSON.parse(json), null, 2);
}

export function featureDefinitionToPrettyJson(
  definition: FeatureDefinition
): string {
  return prettyPrintJson(JSON.stringify({ ...schemaMixin, ...definition }));
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

export function parseYaml(input: string): unknown {
  return yaml.load(input);
}

export function yamlToCompactJson(input: string): string {
  const parsed = parseYaml(input);
  return JSON.stringify(parsed);
}

export function isValidYaml(input: string): boolean {
  try {
    parseYaml(input);
    return true;
  } catch {
    return false;
  }
}