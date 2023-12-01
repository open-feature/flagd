package eval

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/diegoholiveira/jsonlogic/v3"
	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/model"
	"github.com/open-feature/flagd/core/pkg/store"
	"github.com/open-feature/flagd/core/pkg/sync"
	schema "github.com/open-feature/schemas/json"
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

	flagdPropertiesKey = "$flagd"

	// targetingKeyKey is used to extract the targetingKey to bucket on in fractional
	// evaluation if the user did not supply the optional bucketing property.
	targetingKeyKey = "targetingKey"
)

var regBrace *regexp.Regexp

type flagdProperties struct {
	FlagKey   string `json:"flagKey"`
	Timestamp int64  `json:"timestamp"`
}

func init() {
	regBrace = regexp.MustCompile("^[^{]*{|}[^}]*$")
}

type variantEvaluator func(string, string, map[string]any) (
	variant string, variants map[string]interface{}, reason string, metadata map[string]interface{}, error error)

type JSONEvaluator struct {
	store          *store.Flags
	Logger         *logger.Logger
	jsonEvalTracer trace.Tracer
}

type constraints interface {
	bool | string | map[string]any | float64
}

const (
	Disabled = "DISABLED"
)

type JSONEvaluatorOption func(je *JSONEvaluator)

func WithEvaluator(name string, evalFunc func(interface{}, interface{}) interface{}) JSONEvaluatorOption {
	return func(_ *JSONEvaluator) {
		jsonlogic.AddOperator(name, evalFunc)
	}
}

func NewJSONEvaluator(logger *logger.Logger, s *store.Flags, opts ...JSONEvaluatorOption) *JSONEvaluator {
	ev := JSONEvaluator{
		Logger: logger.WithFields(
			zap.String("component", "evaluator"),
			zap.String("evaluator", "json"),
		),
		store:          s,
		jsonEvalTracer: otel.Tracer("jsonEvaluator"),
	}

	for _, o := range opts {
		o(&ev)
	}

	return &ev
}

func (je *JSONEvaluator) GetState() (string, error) {
	s, err := je.store.String()
	if err != nil {
		return "", fmt.Errorf("unable to fetch evaluator state: %w", err)
	}
	return s, nil
}

func (je *JSONEvaluator) SetState(payload sync.DataSync) (map[string]interface{}, bool, error) {
	_, span := je.jsonEvalTracer.Start(
		context.Background(),
		"flagSync",
		trace.WithAttributes(attribute.String("feature_flag.source", payload.Source)),
		trace.WithAttributes(attribute.String("feature_flag.sync_type", payload.Type.String())))
	defer span.End()

	var newFlags Flags

	err := je.configToFlags(payload.FlagData, &newFlags)
	if err != nil {
		span.SetStatus(codes.Error, "flagSync error")
		span.RecordError(err)
		return nil, false, err
	}

	var events map[string]interface{}
	var reSync bool

	switch payload.Type {
	case sync.ALL:
		events, reSync = je.store.Merge(je.Logger, payload.Source, newFlags.Flags)
	case sync.ADD:
		events = je.store.Add(je.Logger, payload.Source, newFlags.Flags)
	case sync.UPDATE:
		events = je.store.Update(je.Logger, payload.Source, newFlags.Flags)
	case sync.DELETE:
		events = je.store.DeleteFlags(je.Logger, payload.Source, newFlags.Flags)
	default:
		return nil, false, fmt.Errorf("unsupported sync type: %d", payload.Type)
	}

	// Number of events correlates to the number of flags changed through this sync, record it
	span.SetAttributes(attribute.Int("feature_flag.change_count", len(events)))

	return events, reSync, nil
}

func (je *JSONEvaluator) ResolveAllValues(ctx context.Context, reqID string, context map[string]any) []AnyValue {
	_, span := je.jsonEvalTracer.Start(ctx, "resolveAll")
	defer span.End()

	values := []AnyValue{}
	var value interface{}
	var variant string
	var reason string
	var metadata map[string]interface{}
	var err error
	allFlags := je.store.GetAll()
	for flagKey, flag := range allFlags {
		if flag.State == Disabled {
			// ignore evaluation of disabled flag
			continue
		}

		defaultValue := flag.Variants[flag.DefaultVariant]
		switch defaultValue.(type) {
		case bool:
			value, variant, reason, metadata, err = resolve[bool](
				reqID,
				flagKey,
				context,
				je.evaluateVariant,
			)
		case string:
			value, variant, reason, metadata, err = resolve[string](
				reqID,
				flagKey,
				context,
				je.evaluateVariant,
			)
		case float64:
			value, variant, reason, metadata, err = resolve[float64](
				reqID,
				flagKey,
				context,
				je.evaluateVariant,
			)
		case map[string]any:
			value, variant, reason, metadata, err = resolve[map[string]any](
				reqID,
				flagKey,
				context,
				je.evaluateVariant,
			)
		}
		if err != nil {
			je.Logger.ErrorWithID(reqID, fmt.Sprintf("bulk evaluation: key: %s returned error: %s", flagKey, err.Error()))
		}
		values = append(values, NewAnyValue(value, variant, reason, flagKey, metadata, err))
	}
	return values
}

func (je *JSONEvaluator) ResolveBooleanValue(
	ctx context.Context, reqID string, flagKey string, context map[string]any) (
	value bool,
	variant string,
	reason string,
	metadata map[string]interface{},
	err error,
) {
	_, span := je.jsonEvalTracer.Start(ctx, "resolveBoolean")
	defer span.End()

	je.Logger.DebugWithID(reqID, fmt.Sprintf("evaluating boolean flag: %s", flagKey))
	return resolve[bool](reqID, flagKey, context, je.evaluateVariant)
}

func (je *JSONEvaluator) ResolveStringValue(
	ctx context.Context, reqID string, flagKey string, context map[string]any) (
	value string,
	variant string,
	reason string,
	metadata map[string]interface{},
	err error,
) {
	_, span := je.jsonEvalTracer.Start(ctx, "resolveString")
	defer span.End()

	je.Logger.DebugWithID(reqID, fmt.Sprintf("evaluating string flag: %s", flagKey))
	return resolve[string](reqID, flagKey, context, je.evaluateVariant)
}

func (je *JSONEvaluator) ResolveFloatValue(
	ctx context.Context, reqID string, flagKey string, context map[string]any) (
	value float64,
	variant string,
	reason string,
	metadata map[string]interface{},
	err error,
) {
	_, span := je.jsonEvalTracer.Start(ctx, "resolveFloat")
	defer span.End()

	je.Logger.DebugWithID(reqID, fmt.Sprintf("evaluating float flag: %s", flagKey))
	value, variant, reason, metadata, err = resolve[float64](reqID, flagKey, context, je.evaluateVariant)
	return
}

func (je *JSONEvaluator) ResolveIntValue(ctx context.Context, reqID string, flagKey string, context map[string]any) (
	value int64,
	variant string,
	reason string,
	metadata map[string]interface{},
	err error,
) {
	_, span := je.jsonEvalTracer.Start(ctx, "resolveInt")
	defer span.End()

	je.Logger.DebugWithID(reqID, fmt.Sprintf("evaluating int flag: %s", flagKey))
	var val float64
	val, variant, reason, metadata, err = resolve[float64](reqID, flagKey, context, je.evaluateVariant)
	value = int64(val)
	return
}

func (je *JSONEvaluator) ResolveObjectValue(
	ctx context.Context, reqID string, flagKey string, context map[string]any) (
	value map[string]any,
	variant string,
	reason string,
	metadata map[string]interface{},
	err error,
) {
	_, span := je.jsonEvalTracer.Start(ctx, "resolveObject")
	defer span.End()

	je.Logger.DebugWithID(reqID, fmt.Sprintf("evaluating object flag: %s", flagKey))
	return resolve[map[string]any](reqID, flagKey, context, je.evaluateVariant)
}

func resolve[T constraints](reqID string, key string, context map[string]any, variantEval variantEvaluator) (
	value T, variant string, reason string, metadata map[string]interface{}, err error,
) {
	variant, variants, reason, metadata, err := variantEval(reqID, key, context)
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

// runs the rules (if defined) to determine the variant, otherwise falling through to the default
// nolint: funlen
func (je *JSONEvaluator) evaluateVariant(reqID string, flagKey string, context map[string]any) (
	variant string, variants map[string]interface{}, reason string, metadata map[string]interface{}, err error,
) {
	metadata = map[string]interface{}{}

	flag, ok := je.store.Get(flagKey)
	if !ok {
		// flag not found
		je.Logger.DebugWithID(reqID, fmt.Sprintf("requested flag could not be found: %s", flagKey))
		return "", map[string]interface{}{}, model.ErrorReason, metadata, errors.New(model.FlagNotFoundErrorCode)
	}

	// add selector to evaluation metadata
	selector := je.store.SelectorForFlag(flag)
	if selector != "" {
		metadata[SelectorMetadataKey] = selector
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

		context = je.setFlagdProperties(context, flagdProperties{
			FlagKey:   flagKey,
			Timestamp: time.Now().Unix(),
		})

		b, err := json.Marshal(context)
		if err != nil {
			je.Logger.ErrorWithID(reqID, fmt.Sprintf("error parsing context for flag: %s, %s, %v", flagKey, err, context))

			return "", flag.Variants, model.ErrorReason, metadata, errors.New(model.ErrorReason)
		}

		var result bytes.Buffer
		// evaluate JsonLogic rules to determine the variant
		err = jsonlogic.Apply(bytes.NewReader(targetingBytes), bytes.NewReader(b), &result)
		if err != nil {
			je.Logger.ErrorWithID(reqID, fmt.Sprintf("error applying rules: %s", err))
			return "", flag.Variants, model.ErrorReason, metadata, errors.New(model.ParseErrorCode)
		}

		// string here are still
		trimmed := strings.TrimSpace(result.String())

		// strip whitespace and quotes from the variant
		variant = strings.ReplaceAll(trimmed, "\"", "")

		// if this is a valid variant, return it
		// TODO: explicit null handling?
		if _, ok := flag.Variants[variant]; ok {
			return variant, flag.Variants, model.TargetingMatchReason, metadata, nil
		} else {
			je.Logger.ErrorWithID(reqID, fmt.Sprintf("invalid or missing variant: %s for flagKey: %s, variant is not valid", variant, flagKey))
			return "", flag.Variants, model.ErrorReason, metadata, errors.New(model.ParseErrorCode)
		}
	} else {
		reason = model.StaticReason
	}

	return flag.DefaultVariant, flag.Variants, reason, metadata, nil
}

func (je *JSONEvaluator) setFlagdProperties(
	context map[string]any,
	properties flagdProperties,
) map[string]any {
	if context == nil {
		context = map[string]any{}
	}

	newContext := maps.Clone(context)

	if _, ok := newContext[flagdPropertiesKey]; ok {
		je.Logger.Warn("overwriting $flagd properties in the context")
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

// configToFlags convert string configurations to flags and store them to pointer newFlags
func (je *JSONEvaluator) configToFlags(config string, newFlags *Flags) error {
	schemaLoader := gojsonschema.NewStringLoader(schema.FlagdDefinitions)
	flagStringLoader := gojsonschema.NewStringLoader(config)

	result, err := gojsonschema.Validate(schemaLoader, flagStringLoader)
	if err != nil {
		return fmt.Errorf("error validating json schema: %w", err)
	} else if !result.Valid() {
		return fmt.Errorf("JSON schema validation failed: %s", buildErrorString(result.Errors()))
	}

	transposedConfig, err := je.transposeEvaluators(config)
	if err != nil {
		return fmt.Errorf("transposing evaluators: %w", err)
	}

	err = json.Unmarshal([]byte(transposedConfig), &newFlags)
	if err != nil {
		return fmt.Errorf("unmarshalling provided configurations: %w", err)
	}

	return validateDefaultVariants(newFlags)
}

// validateDefaultVariants returns an error if any of the default variants aren't valid
func validateDefaultVariants(flags *Flags) error {
	for name, flag := range flags.Flags {
		if _, ok := flag.Variants[flag.DefaultVariant]; !ok {
			return fmt.Errorf(
				"default variant: '%s' isn't a valid variant of flag: '%s'", flag.DefaultVariant, name,
			)
		}
	}

	return nil
}

func (je *JSONEvaluator) transposeEvaluators(state string) (string, error) {
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
