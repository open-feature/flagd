import { useState, useEffect, useMemo, useCallback } from "react";
import { useMedia } from "react-use";
import { FlagdCore, MemoryStorage } from "@openfeature/flagd-core";
import { ScenarioName, scenarios } from "./scenarios";
import type { FlagValueType } from "@openfeature/core";
import { getString, isValidYaml, yamlToCompactJson } from "./utils";
import { BeforeMount, Editor } from "@monaco-editor/react";
import { Observable } from "react-use/lib/useObservable";

declare global {
  var component$: Observable<{ ref: HTMLElement }>;
}

// see: https://github.com/squidfunk/mkdocs-material/discussions/3429
const BODY_COLOR_SCHEME_ATTR = "data-md-color-scheme";
const PALETTE_SWITCH_SELECTOR = "[data-md-component=palette]";
const getPalette = () =>
  document.body.getAttribute(BODY_COLOR_SCHEME_ATTR) &&
    document.body.getAttribute(BODY_COLOR_SCHEME_ATTR) !== "default"
    ? "custom-dark"
    : "custom";
const monacoBeforeMount: BeforeMount = (monaco) => {
  // inherent from the normal vs/vs-dark themes, but with transparent backgrounds
  // so our CSS variables can be used for that (css vars cant be used in the editor theme directly)
  monaco?.editor.defineTheme("custom-dark", {
    base: "vs-dark", // inherent all the normal "dark" syntax highlighting
    inherit: true,
    rules: [],
    colors: {
      "editor.background": "#00000000",
    },
  });
  monaco?.editor.defineTheme("custom", {
    base: "vs", // inherent all the normal "light" syntax highlighting
    inherit: true,
    rules: [],
    colors: {
      "editor.background": "#00000000",
    },
  });
  monaco?.languages.json.jsonDefaults.setDiagnosticsOptions({
    enableSchemaRequest: true,
    allowComments: false, // we don't support JSON comments in flagd
  });
};

function formatJson(shortenedString: string) {
  const object = JSON.parse(shortenedString);
  return JSON.stringify(object, null, 2);
};

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
  const [showCopyNotification, setShowCopyNotification] = useState(false);
  const [status, setStatus] = useState<"success" | "failure">("success");
  const [editorTheme, updateEditorTheme] = useState<"custom" | "custom-dark">(
    getPalette()
  );

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
    setShowCopyNotification(false)
    setStatus("success");
  }, [selectedTemplate]);

  useEffect(() => {
    resetInputs();
  }, [selectedTemplate, resetInputs]);

  const flagStorage = useMemo(() => new MemoryStorage(console), []);
  const flagdCore = useMemo(
    () => new FlagdCore(flagStorage, console),
    [flagStorage]
  );

  useEffect(() => {
    if (isValidYaml(featureDefinition)) {
      try {
        const jsonConfig = yamlToCompactJson(featureDefinition);
        flagdCore.setConfigurations(jsonConfig);
        setAutocompleteFlagKeys(Array.from(flagdCore.getFlags().keys()));
        setValidFeatureDefinition(true);
      } catch (err) {
        console.error("Invalid flagd configuration", err);
        setValidFeatureDefinition(false);
      }
    } else {
      setValidFeatureDefinition(false);
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

  useEffect(() => {
    const urlParams = new URLSearchParams(window.location.search);
    const flagsParam = urlParams.get('flags');
    const flagKeyParam = urlParams.get('flag-key');
    const returnTypeParam = urlParams.get('return-type');
    const evalContextParam = urlParams.get('eval-context');
    const scenarioParam = urlParams.get('scenario-name');
    if (flagsParam) {
      try {
        const formattedFeatureDefinition = formatJson(flagsParam);
        setFeatureDefinition(formattedFeatureDefinition);
        if (flagKeyParam) setFlagKey(flagKeyParam);
        if (returnTypeParam) setReturnType(returnTypeParam as FlagValueType);
        if (evalContextParam) {
          const formattedEvaluationContext = formatJson(evalContextParam);
          setEvaluationContext(formattedEvaluationContext);
        }
      } catch (error) {
        console.error("Error decoding URL parameters: ", error);
      }
    } else if (scenarioParam && scenarios[scenarioParam as keyof typeof scenarios]) {
      setSelectedTemplate(scenarioParam as keyof typeof scenarios);
      setFeatureDefinition(scenarios[scenarioParam as keyof typeof scenarios].flagDefinition);
    }
  }, []);

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
    border: "none",
    backgroundColor: "var(--md-code-bg-color)",
    color: "var(--md-code-fg-color)",
    fontFeatureSettings: "kern",
    fontFamily: "var(--md-code-font-family)",
  };

  const copyUrl = () => {
    const baseUrl = window.location.origin + window.location.pathname;
    const newUrl = new URL(baseUrl)
    const encodedFeatureDefinition = yamlToCompactJson(featureDefinition);
    const encodedEvaluationContext = yamlToCompactJson(evaluationContext);

    if (Object.keys(scenarios).includes(selectedTemplate) &&
      scenarios[selectedTemplate].flagDefinition === featureDefinition) {
      newUrl.searchParams.set('scenario-name', selectedTemplate);
    } else {
      newUrl.searchParams.delete('scenario-name');
      newUrl.searchParams.set('flags', encodedFeatureDefinition);
      newUrl.searchParams.set('flag-key', flagKey);
      newUrl.searchParams.set('return-type', returnType);
      newUrl.searchParams.set('eval-context', encodedEvaluationContext);
    }
    window.history.pushState({}, '', newUrl.href);

    navigator.clipboard.writeText(newUrl.href).then(() => {
      console.log('URL copied to clipboard');
      setShowCopyNotification(true)
      setTimeout(() => {
        setShowCopyNotification(false)
      }, 5000);
    }).catch(err => {
      console.error('Failed to copy URL: ', err);
    });
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
            <div style={{ backgroundColor: codeStyle.backgroundColor }}>
              <Editor
                theme={editorTheme}
                width="100%"
                height="500px"
                language="yaml"
                value={featureDefinition}
                options={{
                  minimap: { enabled: false },
                  lineNumbers: "off",
                }}
                beforeMount={monacoBeforeMount}
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
              <div style={{ backgroundColor: codeStyle.backgroundColor }}>
                <Editor
                  theme={editorTheme}
                  width="100%"
                  height="80px"
                  language="yaml"
                  options={{
                    minimap: { enabled: false },
                    lineNumbers: "off",
                    folding: false,
                  }}
                  beforeMount={monacoBeforeMount}
                  value={evaluationContext}
                  onChange={(value) => {
                    if (value) {
                      setEvaluationContext(value);
                    }
                  }}
                />
              </div>
            </div>
            <div style={{ display: "flex", gap: "8px", paddingTop: "8px" }}>
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
              <button
                className="md-button"
                onClick={copyUrl}
                disabled={!validFeatureDefinition || !validEvaluationContext}
              >
                Share
              </button>
            </div>
            <div
              className={`output ${showOutput ? "visible" : ""} admonition ${status === "success" ? "success" : "failure"
                }`}
            >
              <p className="admonition-title">
                {status === "success" ? "Success" : "Failure"}
              </p>
              {typeof parsedOutput === "object" ? (
                <div style={{ margin: "0.6rem 0 0.6rem 0" }}>
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
            {showCopyNotification && (
              <h4 className="admonition-title" style={{
                paddingLeft: "45px",
                borderLeftWidth: "0rem",
                borderLeftStyle: "solid",
                left: "15px"
              }}>URL copied to clipboard</h4>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}

export default App;
