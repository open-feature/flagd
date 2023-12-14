import { useState, useEffect, useMemo, useCallback } from "react";
import { FlagdCore, MemoryStorage } from "@openfeature/flagd-core";

type Template = {
  [template: string]: {
    flagDefinition: string;
    flagKey: string;
    returnType: string;
    context: string;
  };
};

const templates = {
  "Email ends with": {
    flagDefinition: JSON.stringify(
      JSON.parse(
        '{"flags":{"new-welcome-banner":{"state":"ENABLED","variants":{"on":true,"off":false},"defaultVariant":"off","targeting":{"if":[{"ends_with":[{"var":"email"},"@example.com"]},"on","off"]}}}}'
      ),
      null,
      2
    ),
    flagKey: "new-welcome-banner",
    returnType: "boolean",
    context: JSON.stringify(
      JSON.parse('{ "email": "test@example.com" }'),
      null,
      2
    ),
  },
  "Shared evaluators": {
    flagDefinition: JSON.stringify(
      JSON.parse(
        '{"flags":{"fibAlgo":{"variants":{"recursive":"recursive","memo":"memo","loop":"loop","binet":"binet"},"defaultVariant":"recursive","state":"ENABLED","targeting":{"if":[{"$ref":"emailWithFaas"},"binet",null]}},"headerColor":{"variants":{"red":"#FF0000","blue":"#0000FF","green":"#00FF00","yellow":"#FFFF00"},"defaultVariant":"red","state":"ENABLED","targeting":{"if":[{"$ref":"emailWithFaas"},{"fractional":[{"var":"email"},["red",25],["blue",25],["green",25],["yellow",25]]},null]}}},"$evaluators":{"emailWithFaas":{"ends_with":[{"var": "email"}, "@faas.com"]}}}'
      ),
      null,
      2
    ),
    flagKey: "fibAlgo",
    returnType: "string",
    context: JSON.stringify(
      JSON.parse('{ "email": "mike@faas.com" }'),
      null,
      2
    ),
  },
  Empty: {
    flagDefinition: "",
    flagKey: "",
    returnType: "boolean",
    context: "",
  },
  // Add more templates here...
} satisfies Template;

type TemplateName = keyof typeof templates;

// const flagKeys = ["new-welcome-banner", "test", "test2", "test3"];

function App() {
  const [selectedTemplate, setSelectedTemplate] =
    useState<TemplateName>("Email ends with");
  const [featureDefinition, setFeatureDefinition] = useState(
    templates[selectedTemplate].flagDefinition
  );
  const [flagKey, setFlagKey] = useState(templates[selectedTemplate].flagKey);
  const [returnType, setReturnType] = useState(
    templates[selectedTemplate].returnType
  );
  const [evaluationContext, setEvaluationContext] = useState(
    templates[selectedTemplate].context
  );
  const [showOutput, setShowOutput] = useState(false);
  const [output, setOutput] = useState("");
  const [autocompleteFlagKeys, setAutocompleteFlagKeys] = useState<string[]>(
    []
  );

  const resetInputs = useCallback(() => {
    setOutput("");
    setShowOutput(false);
    const template = templates[selectedTemplate];
    setFeatureDefinition(template.flagDefinition);
    setFlagKey(template.flagKey);
    setReturnType(template.returnType);
    setEvaluationContext(template.context);
  }, [
    selectedTemplate,
    setOutput,
    setShowOutput,
    setFeatureDefinition,
    setFlagKey,
    setReturnType,
    setEvaluationContext,
  ]);

  useEffect(() => {
    resetInputs();
  }, [selectedTemplate, resetInputs]);

  const flagStorage = useMemo(() => new MemoryStorage(), []);
  const flagdCore = useMemo(() => new FlagdCore(flagStorage), [flagStorage]);

  const isValidJson = (jsonString: string) => {
    try {
      JSON.parse(jsonString);
      return true;
    } catch {
      return false;
    }
  };

  useEffect(() => {
    if (isValidJson(featureDefinition)) {
      const flagDefinition = JSON.parse(featureDefinition);
      if (typeof flagDefinition.flags === "object") {
        setAutocompleteFlagKeys(Object.keys(flagDefinition.flags));
      }
      flagdCore.setConfigurations(featureDefinition);
    }
  }, [featureDefinition, flagdCore]);

  const evaluate = () => {
    setShowOutput(true);
    try {
      const context = evaluationContext ? JSON.parse(evaluationContext) : {};
      let result;
      switch (returnType) {
        case "boolean":
          result = flagdCore.resolveBooleanEvaluation(
            flagKey,
            false,
            context,
            console
          );
          break;
        case "string":
          result = flagdCore.resolveStringEvaluation(
            flagKey,
            "",
            context,
            console
          );
          break;
        case "number":
          result = flagdCore.resolveNumberEvaluation(
            flagKey,
            0,
            context,
            console
          );
          break;
        case "object":
          result = flagdCore.resolveObjectEvaluation(
            flagKey,
            {},
            context,
            console
          );
          break;
      }
      setOutput(JSON.stringify(result, null, 2));
    } catch (error) {
      console.error("Invalid JSON input", error);
      setOutput((error as Error).message);
    }
  };

  const parsedOutput = useMemo(() => {
    try {
      return JSON.parse(output);
    } catch {
      return output;
    }
  }, [output]);

  const admonitionTitle =
    typeof parsedOutput === "object" ? "Success" : "Failure";
  const admonitionClass =
    typeof parsedOutput === "object" ? "success" : "failure";

  return (
    <div>
      <div style={{ marginBlock: "10px" }}>
        <label style={{ display: "block", marginBottom: "5px" }}>
          Select a template
        </label>
        <select
          style={{
            width: "200px",
            padding: "10px",
            boxSizing: "border-box",
          }}
          value={selectedTemplate}
          onChange={(e) => setSelectedTemplate(e.target.value as TemplateName)}
        >
          {Object.keys(templates).map((templateName) => (
            <option key={templateName} value={templateName}>
              {templateName}
            </option>
          ))}
        </select>
      </div>

      <div style={{ display: "flex", marginBottom: "10px" }}>
        <div style={{ flex: "1", marginRight: "20px", textAlign: "left" }}>
          <label style={{ display: "block", marginBottom: "5px" }}>
            Feature Definition
          </label>
          <textarea
            style={{
              width: "500px",
              height: "500px",
              minWidth: "450px",
              maxWidth: "800px",
              minHeight: "400px",
              padding: "10px",
            }}
            name="feature-definition"
            value={featureDefinition}
            onChange={(e) => setFeatureDefinition(e.target.value)}
          />
        </div>
        <div style={{ flex: "1" }}>
          <div style={{ marginBottom: "10px", textAlign: "left" }}>
            <label style={{ display: "block", marginBottom: "5px" }}>
              Flag Key
              <input
                style={{
                  width: "100%",
                  padding: "10px",
                  boxSizing: "border-box",
                  border: "1px solid",
                }}
                name="flag-key"
                list="flag-keys"
                value={flagKey}
                onChange={(e) => setFlagKey(e.target.value)}
              />
              <datalist id="flag-keys">
                {autocompleteFlagKeys.map((key, index) => (
                  <option key={index} value={key} />
                ))}
              </datalist>
            </label>
          </div>
          <div style={{ marginBottom: "10px", textAlign: "left" }}>
            <label style={{ display: "block", marginBottom: "5px" }}>
              Return Type
              <select
                style={{
                  width: "100%",
                  padding: "10px",
                  boxSizing: "border-box",
                }}
                value={returnType}
                onChange={(e) => setReturnType(e.target.value)}
              >
                <option value="boolean">boolean</option>
                <option value="string">string</option>
                <option value="number">number</option>
                <option value="object">object</option>
              </select>
            </label>
          </div>
          <div style={{ marginBottom: "10px", textAlign: "left" }}>
            <label style={{ display: "block", marginBottom: "5px" }}>
              Evaluation Context
              <textarea
                style={{
                  width: "100%",
                  minHeight: "100px",
                  maxHeight: "300px",
                  padding: "10px",
                  boxSizing: "border-box",
                  resize: "vertical",
                }}
                name="evaluation-context"
                value={evaluationContext}
                onChange={(e) => setEvaluationContext(e.target.value)}
              />
            </label>
          </div>
          <div style={{ display: "flex", gap: "10px" }}>
            <button
              className="md-button md-button--primary"
              onClick={evaluate}
              disabled={
                !isValidJson(featureDefinition) ||
                (!isValidJson(evaluationContext) && evaluationContext !== "")
              }
            >
              Evaluate
            </button>
            <button className="md-button" onClick={resetInputs}>
              Reset
            </button>
          </div>
          <div
            className={`output ${
              showOutput ? "visible" : ""
            } admonition ${admonitionClass}`}
          >
            <p className="admonition-title">{admonitionTitle}</p>
            {typeof parsedOutput === "object" ? (
              <div>
                {Object.entries(parsedOutput).map(([key, value]) => (
                  <div key={key}>
                    <strong>{key}:</strong> {JSON.stringify(value)}
                  </div>
                ))}
              </div>
            ) : (
              <p>{parsedOutput}</p>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}

export default App;
