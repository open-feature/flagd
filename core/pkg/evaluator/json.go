package evaluator

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	flagd_definitions "github.com/open-feature/flagd-schemas/json"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/diegoholiveira/jsonlogic/v3"
	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/model"
	"github.com/open-feature/flagd/core/pkg/store"
	"github.com/open-feature/flagd/core/pkg/sync"
	"github.com/santhosh-tekuri/jsonschema/v6"
	"github.com/xeipuuv/gojsonschema"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"golang.org/x/exp/maps"
)

const (
	SelectorMetadataKey = "scope"
	flagdPropertiesKey  = "$flagd"
	// targetingKeyKey is used to extract the targetingKey to bucket on in fractional
	// evaluation if the user did not supply the optional bucketing property.
	targetingKeyKey = "targetingKey"
	Disabled        = "DISABLED"
)

var compiledSchema *jsonschema.Schema
var regBrace *regexp.Regexp

func init() {
	regBrace = regexp.MustCompile("^[^{]*{|}[^}]*$")

	// Create a new JSON Schema compiler
	compiler := jsonschema.NewCompiler()

	// Add the Flagd Daemon schema
	flagdFile := strings.NewReader(flagd_definitions.FlagdSchema)
	schema, err := jsonschema.UnmarshalJSON(flagdFile)
	if err != nil {
		log.Fatalf("Failed to unmarshal targeting schema: %v", err)
	}
	if err := compiler.AddResource("https://flagd.dev/schema/v0/flagd.json", schema); err != nil {
		log.Fatalf("Failed to add flagd schema: %v", err)
	}

	// Add the Flag schema
	flagsFile := strings.NewReader(flagd_definitions.FlagSchema)
	schema, err = jsonschema.UnmarshalJSON(flagsFile)
	if err != nil {
		log.Fatalf("Failed to unmarshal targeting schema: %v", err)
	}
	if err := compiler.AddResource("https://flagd.dev/schema/v0/flags.json", schema); err != nil {
		log.Fatalf("Failed to add flags schema: %v", err)
	}

	// Add the Targeting schema
	targetingFile := strings.NewReader(flagd_definitions.TargetingSchema)
	schema, err = jsonschema.UnmarshalJSON(targetingFile)
	if err != nil {
		log.Fatalf("Failed to unmarshal targeting schema: %v", err)
	}
	if err := compiler.AddResource("https://flagd.dev/schema/v0/targeting.json", schema); err != nil {
		log.Fatalf("Failed to add targeting schema: %v", err)
	}
	compiledSchema, err = compiler.Compile("https://flagd.dev/schema/v0/flagd.json")
	if err != nil {
		log.Fatalf("Failed to compile flagd schema: %v", err)
	}

}

type constraints interface {
	bool | string | map[string]any | float64 | interface{}
}

type JSONEvaluatorOption func(je *JSON)

type flagdProperties struct {
	FlagKey   string `json:"flagKey"`
	Timestamp int64  `json:"timestamp"`
}

type variantEvaluator func(context.Context, string, string, map[string]any) (
	variant string, variants map[string]interface{}, reason string, metadata map[string]interface{}, error error)

// Deprecated - this will be remove in the next release
func WithEvaluator(name string, evalFunc func(interface{}, interface{}) interface{}) JSONEvaluatorOption {
	return func(_ *JSON) {
		jsonlogic.AddOperator(name, evalFunc)
	}
}

// JSON evaluator
type JSON struct {
	store          store.IStore
	Logger         *logger.Logger
	jsonEvalTracer trace.Tracer
	Resolver
}

func NewJSON(logger *logger.Logger, s store.IStore, opts ...JSONEvaluatorOption) *JSON {
	logger = logger.WithFields(
		zap.String("component", "evaluator"),
		zap.String("evaluator", "json"),
	)
	tracer := otel.Tracer("jsonEvaluator")

	ev := JSON{
		store:          s,
		Logger:         logger,
		jsonEvalTracer: tracer,
		Resolver:       NewResolver(s, logger, tracer),
	}

	for _, o := range opts {
		o(&ev)
	}

	return &ev
}

func (je *JSON) GetState() (string, error) {
	s, err := je.store.String()
	if err != nil {
		return "", fmt.Errorf("unable to fetch evaluator state: %w", err)
	}
	return s, nil
}

func (je *JSON) SetState(payload sync.DataSync) (map[string]interface{}, bool, error) {
	_, span := je.jsonEvalTracer.Start(
		context.Background(),
		"flagSync",
		trace.WithAttributes(attribute.String("feature_flag.source", payload.Source)))
	defer span.End()

	var definition Definition

	err := configToFlagDefinition(je.Logger, payload.FlagData, &definition)
	if err != nil {
		span.SetStatus(codes.Error, "flagSync error")
		span.RecordError(err)
		return nil, false, err
	}

	var events map[string]interface{}
	var reSync bool

	events, reSync = je.store.Update(payload.Source, definition.Flags, definition.Metadata)

	// Number of events correlates to the number of flags changed through this sync, record it
	span.SetAttributes(attribute.Int("feature_flag.change_count", len(events)))

	return events, reSync, nil
}

// Resolver implementation for flagd flags. This resolver should be kept reusable, hence must interact with interfaces.
type Resolver struct {
	store  store.IStore
	Logger *logger.Logger
	tracer trace.Tracer
}

func NewResolver(store store.IStore, logger *logger.Logger, jsonEvalTracer trace.Tracer) Resolver {
	// register supported json logic custom operator implementations
	jsonlogic.AddOperator(FractionEvaluationName, NewFractional(logger).Evaluate)
	jsonlogic.AddOperator(StartsWithEvaluationName, NewStringComparisonEvaluator(logger).StartsWithEvaluation)
	jsonlogic.AddOperator(EndsWithEvaluationName, NewStringComparisonEvaluator(logger).EndsWithEvaluation)
	jsonlogic.AddOperator(SemVerEvaluationName, NewSemVerComparison(logger).SemVerEvaluation)

	return Resolver{store: store, Logger: logger, tracer: jsonEvalTracer}
}

func (je *Resolver) ResolveAllValues(ctx context.Context, reqID string, context map[string]any) ([]AnyValue,
	model.Metadata, error,
) {
	_, span := je.tracer.Start(ctx, "resolveAll")
	defer span.End()

	var selector store.Selector
	s := ctx.Value(store.SelectorContextKey{})
	if s != nil {
		selector = s.(store.Selector)
	}
	allFlags, flagSetMetadata, err := je.store.GetAll(ctx, &selector)
	if err != nil {
		return nil, flagSetMetadata, fmt.Errorf("error retreiving flags from the store: %w", err)
	}

	values := []AnyValue{}
	var value interface{}
	var variant string
	var reason string
	var metadata map[string]interface{}

	for _, flag := range allFlags {
		if flag.State == Disabled {
			// ignore evaluation of disabled flag
			continue
		}

		defaultValue := flag.Variants[flag.DefaultVariant]
		switch defaultValue.(type) {
		case bool:
			value, variant, reason, metadata, err = resolve[bool](ctx, reqID, flag.Key, context, je.evaluateVariant)
		case string:
			value, variant, reason, metadata, err = resolve[string](ctx, reqID, flag.Key, context, je.evaluateVariant)
		case float64:
			value, variant, reason, metadata, err = resolve[float64](ctx, reqID, flag.Key, context, je.evaluateVariant)
		case map[string]any:
			value, variant, reason, metadata, err = resolve[map[string]any](ctx, reqID, flag.Key, context, je.evaluateVariant)
		}
		if err != nil {
			je.Logger.ErrorWithID(reqID, fmt.Sprintf("bulk evaluation: key: %s returned error: %s", flag.Key, err.Error()))
		}
		values = append(values, NewAnyValue(value, variant, reason, flag.Key, metadata, err))
	}

	return values, flagSetMetadata, nil
}

func (je *Resolver) ResolveBooleanValue(
	ctx context.Context, reqID string, flagKey string, context map[string]any) (
	value bool,
	variant string,
	reason string,
	metadata map[string]interface{},
	err error,
) {
	_, span := je.tracer.Start(ctx, "resolveBoolean")
	defer span.End()

	je.Logger.DebugWithID(reqID, fmt.Sprintf("evaluating boolean flag: %s", flagKey))
	return resolve[bool](ctx, reqID, flagKey, context, je.evaluateVariant)
}

func (je *Resolver) ResolveStringValue(
	ctx context.Context, reqID string, flagKey string, context map[string]any) (
	value string,
	variant string,
	reason string,
	metadata map[string]interface{},
	err error,
) {
	_, span := je.tracer.Start(ctx, "resolveString")
	defer span.End()

	je.Logger.DebugWithID(reqID, fmt.Sprintf("evaluating string flag: %s", flagKey))
	return resolve[string](ctx, reqID, flagKey, context, je.evaluateVariant)
}

func (je *Resolver) ResolveFloatValue(
	ctx context.Context, reqID string, flagKey string, context map[string]any) (
	value float64,
	variant string,
	reason string,
	metadata map[string]interface{},
	err error,
) {
	_, span := je.tracer.Start(ctx, "resolveFloat")
	defer span.End()

	je.Logger.DebugWithID(reqID, fmt.Sprintf("evaluating float flag: %s", flagKey))
	value, variant, reason, metadata, err = resolve[float64](ctx, reqID, flagKey, context, je.evaluateVariant)
	return
}

func (je *Resolver) ResolveIntValue(ctx context.Context, reqID string, flagKey string, context map[string]any) (
	value int64,
	variant string,
	reason string,
	metadata map[string]interface{},
	err error,
) {
	_, span := je.tracer.Start(ctx, "resolveInt")
	defer span.End()

	je.Logger.DebugWithID(reqID, fmt.Sprintf("evaluating int flag: %s", flagKey))
	var val float64
	val, variant, reason, metadata, err = resolve[float64](ctx, reqID, flagKey, context, je.evaluateVariant)
	value = int64(val)
	return
}

func (je *Resolver) ResolveObjectValue(
	ctx context.Context, reqID string, flagKey string, context map[string]any) (
	value map[string]any,
	variant string,
	reason string,
	metadata map[string]interface{},
	err error,
) {
	_, span := je.tracer.Start(ctx, "resolveObject")
	defer span.End()

	je.Logger.DebugWithID(reqID, fmt.Sprintf("evaluating object flag: %s", flagKey))
	return resolve[map[string]any](ctx, reqID, flagKey, context, je.evaluateVariant)
}

func (je *Resolver) ResolveAsAnyValue(
	ctx context.Context,
	reqID string,
	flagKey string,
	context map[string]any,
) AnyValue {
	_, span := je.tracer.Start(ctx, "resolveAnyValue")
	defer span.End()

	je.Logger.DebugWithID(reqID, fmt.Sprintf("evaluating flag `%s` as a generic flag", flagKey))
	value, variant, reason, meta, err := resolve[interface{}](ctx, reqID, flagKey, context, je.evaluateVariant)
	return NewAnyValue(value, variant, reason, flagKey, meta, err)
}

// resolve is a helper for generic flag resolving
func resolve[T constraints](ctx context.Context, reqID string, key string, context map[string]any,
	variantEval variantEvaluator) (value T, variant string, reason string, metadata map[string]interface{}, err error,
) {
	variant, variants, reason, metadata, err := variantEval(ctx, reqID, key, context)
	if err != nil {
		return value, variant, reason, metadata, err
	}

	var ok bool
	value, ok = variants[variant].(T)
	if !ok {
		return value, variant, model.ErrorReason, metadata, errors.New(model.TypeMismatchErrorCode)
	}

	return value, variant, reason, metadata, nil
}

// nolint: funlen
func (je *Resolver) evaluateVariant(ctx context.Context, reqID string, flagKey string, evalCtx map[string]any) (
	variant string, variants map[string]interface{}, reason string, metadata map[string]interface{}, err error,
) {

	var selector store.Selector
	s := ctx.Value(store.SelectorContextKey{})
	if s != nil {
		selector = s.(store.Selector)
	}
	flag, metadata, err := je.store.Get(ctx, flagKey, &selector)
	if err != nil {
		// flag not found
		je.Logger.DebugWithID(reqID, fmt.Sprintf("requested flag could not be found: %s", flagKey))
		return "", map[string]interface{}{}, model.ErrorReason, metadata, errors.New(model.FlagNotFoundErrorCode)
	}

	for key, value := range flag.Metadata {
		// If value is not nil or empty, copy to metadata
		if value != nil {
			metadata[key] = value
		}
	}

	if flag.State == Disabled {
		je.Logger.DebugWithID(reqID, fmt.Sprintf("requested flag is disabled: %s", flagKey))
		return "", flag.Variants, model.ErrorReason, metadata, errors.New(model.FlagDisabledErrorCode)
	}

	// get the targeting logic, if any
	targeting := flag.Targeting

	if targeting != nil && string(targeting) != "{}" {
		targetingBytes, err := targeting.MarshalJSON()
		if err != nil {
			je.Logger.ErrorWithID(reqID, fmt.Sprintf("Error parsing rules for flag: %s, %s", flagKey, err))
			return "", flag.Variants, model.ErrorReason, metadata, errors.New(model.ParseErrorCode)
		}

		evalCtx = setFlagdProperties(je.Logger, evalCtx, flagdProperties{
			FlagKey:   flagKey,
			Timestamp: time.Now().Unix(),
		})

		b, err := json.Marshal(evalCtx)
		if err != nil {
			je.Logger.ErrorWithID(reqID, fmt.Sprintf("error parsing context for flag: %s, %s, %v", flagKey, err, evalCtx))

			return "", flag.Variants, model.ErrorReason, metadata, errors.New(model.ErrorReason)
		}

		var result bytes.Buffer
		// evaluate JsonLogic rules to determine the variant
		err = jsonlogic.Apply(bytes.NewReader(targetingBytes), bytes.NewReader(b), &result)
		if err != nil {
			je.Logger.ErrorWithID(reqID, fmt.Sprintf("error applying targeting rules: %s", err))
			return "", flag.Variants, model.ErrorReason, metadata, errors.New(model.ParseErrorCode)
		}

		// check if string is "null" before we strip quotes, so we can differentiate between JSON null and "null"
		trimmed := strings.TrimSpace(result.String())

		if trimmed == "null" {
			if flag.DefaultVariant == "" {
				return "", flag.Variants, model.ErrorReason, metadata, errors.New(model.FlagNotFoundErrorCode)
			}

			return flag.DefaultVariant, flag.Variants, model.DefaultReason, metadata, nil
		}

		// strip whitespace and quotes from the variant
		variant = strings.ReplaceAll(trimmed, "\"", "")

		// if this is a valid variant, return it
		if _, ok := flag.Variants[variant]; ok {
			return variant, flag.Variants, model.TargetingMatchReason, metadata, nil
		}
		je.Logger.ErrorWithID(reqID,
			fmt.Sprintf("invalid or missing variant: %s for flagKey: %s, variant is not valid", variant, flagKey))
		return "", flag.Variants, model.ErrorReason, metadata, errors.New(model.GeneralErrorCode)
	}

	if flag.DefaultVariant == "" {
		return "", flag.Variants, model.ErrorReason, metadata, errors.New(model.FlagNotFoundErrorCode)
	}

	return flag.DefaultVariant, flag.Variants, model.StaticReason, metadata, nil
}

func setFlagdProperties(
	log *logger.Logger,
	context map[string]any,
	properties flagdProperties,
) map[string]any {
	if context == nil {
		context = map[string]any{}
	}

	newContext := maps.Clone(context)

	if _, ok := newContext[flagdPropertiesKey]; ok {
		log.Warn("overwriting $flagd properties in the context")
	}

	newContext[flagdPropertiesKey] = properties

	return newContext
}

func getFlagdProperties(context map[string]any) (flagdProperties, bool) {
	properties, ok := context[flagdPropertiesKey]
	if !ok {
		return flagdProperties{}, false
	}

	b, err := json.Marshal(properties)
	if err != nil {
		return flagdProperties{}, false
	}

	var p flagdProperties
	if err := json.Unmarshal(b, &p); err != nil {
		return flagdProperties{}, false
	}

	return p, true
}

type JsonDef struct {
	Flags    interface{}            `json:"flags"`
	Metadata map[string]interface{} `json:"metadata"`
}

type Flag struct {
	model.Flag
	Key       string `json:"key,omitempty"`
	FlagSetId string `json:"flagSetId,omitempty"`
}

// configToFlagDefinition convert string configurations to flags and store them to pointer newFlags
func configToFlagDefinition(log *logger.Logger, config string, definition *Definition) error {
	// Transpose evaluators and unmarshal directly into JsonDef
	transposedConfig, err := transposeEvaluators(config)
	if err != nil {
		return fmt.Errorf("transposing evaluators: %w", err)
	}

	var intermediateConfig JsonDef
	err = json.Unmarshal([]byte(transposedConfig), &intermediateConfig)
	if err != nil {
		return fmt.Errorf("unmarshalling provided configurations: %w", err)
	}

	definition.Metadata = intermediateConfig.Metadata
	definition.Flags = []model.Flag{}

	// Process flags directly
	switch v := intermediateConfig.Flags.(type) {
	case map[string]interface{}: // Handle ValidFlags format
		for k, e := range v {
			flag, err := convertToModelFlag(e)
			if err != nil {
				return fmt.Errorf("failed to process flag for key %s: %w", k, err)
			}
			flag.Key = k // Populate the `Key` field explicitly
			definition.Flags = append(definition.Flags, flag)
		}

	case []interface{}: // Handle ValidMapFlags format
		for _, value := range v {
			flag, err := convertToModelFlag(value)
			if err != nil {
				return fmt.Errorf("failed to process flag: %w", err)
			}
			definition.Flags = append(definition.Flags, flag)
		}

	default:
		return fmt.Errorf("unexpected type for flags property: %T", v)
	}

	return validateDefaultVariants(definition)
}

// Helper function to convert a generic interface{} to model.Flag
func convertToModelFlag(data interface{}) (model.Flag, error) {
	flag := model.Flag{}

	// Assert the data is a map[string]interface{}
	flagData, ok := data.(map[string]interface{})
	if !ok {
		return flag, fmt.Errorf("unexpected type for flag data: %T", data)
	}

	// Populate fields of model.Flag
	if key, ok := flagData["key"].(string); ok {
		flag.Key = key
	}

	if state, ok := flagData["state"].(string); ok {
		flag.State = state
	}

	if defaultVariant, ok := flagData["defaultVariant"].(string); ok {
		flag.DefaultVariant = defaultVariant
	}

	if variants, ok := flagData["variants"].(map[string]interface{}); ok {
		flag.Variants = variants
	}

	if targeting, ok := flagData["targeting"].(json.RawMessage); ok {
		flag.Targeting = targeting
	} else if targetingRaw, ok := flagData["targeting"].(string); ok {
		flag.Targeting = json.RawMessage(targetingRaw)
	} else if flagData["targeting"] != nil {
		marshal, err := json.Marshal(flagData["targeting"])
		if err != nil {
			return flag, fmt.Errorf("marshalling targeting: %w", err)
		}

		flag.Targeting = json.RawMessage(marshal)
	}

	if source, ok := flagData["source"].(string); ok {
		flag.Source = source
	}

	if metadata, ok := flagData["metadata"].(map[string]interface{}); ok {
		// Assuming Metadata is a struct that can be populated directly
		flag.Metadata = model.Metadata(metadata)
	}

	return flag, nil
}

// validateDefaultVariants returns an error if any of the default variants aren't valid
func validateDefaultVariants(flags *Definition) error {
	for _, flag := range flags.Flags {
		// Default Variant is not provided in the config
		if flag.DefaultVariant == "" {
			continue
		}

		if _, ok := flag.Variants[flag.DefaultVariant]; !ok {
			return fmt.Errorf(
				"default variant: '%s' isn't a valid variant of flag: '%s'", flag.DefaultVariant, flag.Key,
			)
		}
	}

	return nil
}

func transposeEvaluators(state string) (string, error) {
	var evaluators Evaluators
	if err := json.Unmarshal([]byte(state), &evaluators); err != nil {
		return "", fmt.Errorf("unmarshal: %w", err)
	}

	for evalName, evalRaw := range evaluators.Evaluators {
		// replace any occurrences of "evaluator": "evalName"
		regex, err := regexp.Compile(fmt.Sprintf(`"\$ref":(\s)*"%s"`, evalName))
		if err != nil {
			return "", fmt.Errorf("compile regex: %w", err)
		}

		marshalledEval, err := evalRaw.MarshalJSON()
		if err != nil {
			return "", fmt.Errorf("marshal evaluator: %w", err)
		}

		evalValue := string(marshalledEval)
		if len(evalValue) < 3 {
			return "", errors.New("evaluator object is empty")
		}
		evalValue = regBrace.ReplaceAllString(evalValue, "")
		state = regex.ReplaceAllString(state, evalValue)
	}

	return state, nil
}

// buildErrorString efficiently converts json schema errors to a formatted string, usable for logging
func buildErrorString(errors []gojsonschema.ResultError) string {
	var builder strings.Builder

	for i, err := range errors {
		builder.WriteByte(' ')
		builder.WriteString(strconv.Itoa(i + 1))
		builder.WriteByte(':')
		builder.WriteString(err.String())
		builder.WriteByte(' ')
	}

	return builder.String()
}
