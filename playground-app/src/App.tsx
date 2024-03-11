import { useState, useEffect, useMemo, useCallback } from "react";
import { useMedia } from "react-use";
import { FlagdCore, MemoryStorage } from "@openfeature/flagd-core";
import { ScenarioName, scenarios } from "./scenarios";
import type { FlagValueType } from "@openfeature/core";
import { getString, isValidJson } from "./utils";
import { Editor } from "@monaco-editor/react";
import { Observable } from "react-use/lib/useObservable";

declare global {
  var component$: Observable<{ ref: HTMLElement}>;
}

// see: https://github.com/squidfunk/mkdocs-material/discussions/3429
const BODY_COLOR_SCHEME_ATTR = 'data-md-color-scheme';
const PALETTE_SWITCH_SELECTOR = '[data-md-component=palette]';
const getPalette = () => document.body.getAttribute(BODY_COLOR_SCHEME_ATTR) &&
  document.body.getAttribute(BODY_COLOR_SCHEME_ATTR) !== 'default' ? 'custom-dark' : 'custom';

function App() {
  const [selectedTemplate, setSelectedTemplate] =
    useState<ScenarioName>("Basic boolean flag");
  const [featureDefinition, setFeatureDefinition] = useState(
    scenarios[selectedTemplate].flagDefinition
  );
  const [flagKey, setFlagKey] = useState(scenarios[selectedTemplate].flagKey);
  const [returnType, setReturnType] = useState(
    scenarios[selectedTemplate].returnType
  );
  const [evaluationContext, setEvaluationContext] = useState(
    getString(scenarios[selectedTemplate].context)
  );
  const [showOutput, setShowOutput] = useState(false);
  const [output, setOutput] = useState("");
  const [autocompleteFlagKeys, setAutocompleteFlagKeys] = useState<string[]>(
    []
  );
  const [description, setDescription] = useState(
    scenarios[selectedTemplate].description
  );
  const [validFeatureDefinition, setValidFeatureDefinition] = useState(true);
  const [validEvaluationContext, setValidEvaluationContext] = useState(true);
  const [status, setStatus] = useState<"success" | "failure">("success");
  const [editorTheme, updateEditorTheme] = useState<"custom" | "custom-dark">(getPalette());

  const resetInputs = useCallback(() => {
    setOutput("");
    setShowOutput(false);
    const template = scenarios[selectedTemplate];
    setFeatureDefinition(template.flagDefinition);
    setFlagKey(template.flagKey);
    setReturnType(template.returnType);
    setEvaluationContext(getString(template.context));
    setDescription(template.description);
    setValidFeatureDefinition(true);
    setValidEvaluationContext(true);
    setStatus("success");
  }, [selectedTemplate]);

  useEffect(() => {
    resetInputs();
  }, [selectedTemplate, resetInputs]);

  const flagStorage = useMemo(() => new MemoryStorage(), []);
  const flagdCore = useMemo(
    () => new FlagdCore(flagStorage, console),
    [flagStorage]
  );

  useEffect(() => {
    if (isValidJson(featureDefinition)) {
      try {
        flagdCore.setConfigurations(featureDefinition);
        setAutocompleteFlagKeys(Array.from(flagdCore.getFlags().keys()));
        setValidFeatureDefinition(true);
      } catch (err) {
        console.error("Invalid flagd configuration", err);
        setValidFeatureDefinition(false);
      }
    }
  }, [featureDefinition, flagdCore]);

  useEffect(() => {
    try {
      JSON.parse(evaluationContext);
      setValidEvaluationContext(true);
    } catch (err) {
      console.error("Invalid JSON input", err);
      setValidEvaluationContext(false);
    }
  }, [evaluationContext]);

  useEffect(() => {
    // update the monaco theme based on the mkdocs theme, see: https://github.com/squidfunk/mkdocs-material/discussions/3429
    const ref = document.querySelector(PALETTE_SWITCH_SELECTOR);
    const subscription = window.component$?.subscribe((component) => {
      if (component?.ref === ref) {
        updateEditorTheme(getPalette());
      }
    });
    return () => {
      subscription?.unsubscribe();
    };
  });

  const evaluate = () => {
    setShowOutput(true);
    try {
      let result;
      const context = JSON.parse(evaluationContext);
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
      setStatus("success");
      setOutput(JSON.stringify(result, null, 2));
    } catch (error) {
      console.error("Invalid JSON input", error);
      setStatus("failure");
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

  const isCompact = useMedia("(max-width: 1220px)");

  const codeStyle = {
    border: 'none',
    backgroundColor: 'var(--md-code-bg-color)',
    color: 'var(--md-code-fg-color)',
    fontFeatureSettings: 'kern',
    fontFamily: 'var(--md-code-font-family)'
  };

  return (
    <div
      style={{
        maxWidth: "825px",
      }}
    >
      <div>
        <p
          style={{
            // Moves content closer to the page header for more screen real estate
            margin: "-32px 0 0 0",
            lineHeight: "1.4",
            fontSize: "medium",
          }}
        >
          Explore flagd flag definitions in your browser. Begin by selecting an
          example below; these are merely starting points, so customize the flag
          definition as you wish. Find an overview of the flag definition
          structure <a href="/reference/flag-definitions/">here</a>.
        </p>
      </div>
      <div>
        <h4>Select a scenario</h4>
        <div
          style={{
            display: "flex",
            flexDirection: isCompact ? "column" : "row",
            textAlign: "left",
            gap: "16px",
            height: "100%",
          }}
        >
          <div style={{ flex: "2" }}>
            <select
              style={{
                width: "100%",
                minWidth: "250px",
                padding: "8px",
                ...codeStyle,
              }}
              value={selectedTemplate}
              onChange={(e) =>
                setSelectedTemplate(e.target.value as ScenarioName)
              }
            >
              {Object.keys(scenarios).map((templateName) => (
                <option key={templateName} value={templateName}>
                  {templateName}
                </option>
              ))}
            </select>
          </div>
          <div style={{ flex: "3" }}>
            <p
              style={{
                lineHeight: "1.4",
                margin: "-4px 0 0 0",
                fontSize: "small",
              }}
            >
              {description}
            </p>
          </div>
        </div>
        <div
          style={{
            display: "flex",
            flexDirection: isCompact ? "column" : "row",
            textAlign: "left",
            gap: "16px",
            height: "100%",
          }}
        >
          <div
            style={{
              flex: "3",
            }}
          >
            <h4>Feature definition</h4>
            <div style={{ backgroundColor: codeStyle.backgroundColor }} >
              <Editor
                theme={editorTheme}
                width="100%"
                height="500px"
                defaultLanguage="json"
                value={featureDefinition}
                options={{
                  minimap: { enabled: false }
                }}
                beforeMount={(monaco) => {
                  // inherent from the normal vs/vs-dark themes, but with transparent backgrounds
                  // so our CSS variables can be used for that (css vars cant be used in the editor theme directly)
                  monaco?.editor.defineTheme('custom-dark', {
                    base: 'vs-dark', // inherent all the normal "dark" syntax highlighting
                    inherit: true,
                    rules: [],
                    colors: {
                      "editor.background": "#00000000",
                    }
                  });
                  monaco?.editor.defineTheme('custom', {
                    base: 'vs', // inherent all the normal "light" syntax highlighting
                    inherit: true,
                    rules: [],
                    colors: {
                      "editor.background": "#00000000",
                    }
                  });
                  monaco?.languages.json.jsonDefaults.setDiagnosticsOptions({
                    enableSchemaRequest: true,
                    allowComments: true,
                  });
                }}
                onChange={(value) => {
                  if (value) {
                    setFeatureDefinition(value);
                  }
                }}
              />
            </div>
          </div>
          <div
            style={{
              flex: "2",
            }}
          >
            <div>
              <h4>Flag key</h4>
              <input
                style={{
                  width: "100%",
                  maxWidth: "800px",
                  padding: "8px",
                  boxSizing: "border-box",
                  ...codeStyle,
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
            </div>
            <div>
              <h4>Return type</h4>
              <select
                style={{
                  width: "100%",
                  padding: "8px 0 8px 0",
                  ...codeStyle,
                }}
                value={returnType}
                onChange={(e) => setReturnType(e.target.value as FlagValueType)}
              >
                <option value="boolean">boolean</option>
                <option value="string">string</option>
                <option value="number">number</option>
                <option value="object">object</option>
              </select>
            </div>
            <div>
              <h4>Evaluation context</h4>
              <textarea
                style={{
                  width: "100%",
                  minHeight: "80px",
                  boxSizing: "border-box",
                  resize: "vertical",
                  padding: "8px",
                  ...codeStyle,
                }}
                name="evaluation-context"
                value={evaluationContext}
                onChange={(e) => setEvaluationContext(e.target.value)}
              />
            </div>
            <div style={{ display: "flex", gap: "8px" }}>
              <button
                className="md-button md-button--primary"
                onClick={evaluate}
                disabled={!validFeatureDefinition || !validEvaluationContext}
              >
                Evaluate
              </button>
              <button className="md-button" onClick={resetInputs}>
                Reset
              </button>
            </div>
            <div
              className={`output ${showOutput ? "visible" : ""} admonition ${
                status === "success" ? "success" : "failure"
              }`}
            >
              <p className="admonition-title">
                {status === "success" ? "Success" : "Failure"}
              </p>
              {typeof parsedOutput === "object" ? (
                <div style={{ margin: '0.6rem 0 0.6rem 0' }} >
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
    </div>
  );
}

export default App;
