// Package gen provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.11.0 DO NOT EDIT.
package gen

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/go-chi/chi/v5"
)

// A logical identifier for the subject (end-user, service) of this flag evaluation.
type Context struct {
	TargetingKey         *string                `json:"targetingKey,omitempty"`
	AdditionalProperties map[string]interface{} `json:"-"`
}

// The resolution details of the flag resolution operation.
type ResolutionDetails struct {
	Reason  *string `json:"reason,omitempty"`
	Variant *string `json:"variant,omitempty"`
}

// ResolutionDetailsBoolean defines model for resolutionDetailsBoolean.
type ResolutionDetailsBoolean struct {
	Reason  *string `json:"reason,omitempty"`
	Value   bool    `json:"value"`
	Variant *string `json:"variant,omitempty"`
}

// ResolutionDetailsNumber defines model for resolutionDetailsNumber.
type ResolutionDetailsNumber struct {
	Reason  *string `json:"reason,omitempty"`
	Value   float32 `json:"value"`
	Variant *string `json:"variant,omitempty"`
}

// ResolutionDetailsObject defines model for resolutionDetailsObject.
type ResolutionDetailsObject struct {
	Reason  *string                       `json:"reason,omitempty"`
	Value   ResolutionDetailsObject_Value `json:"value"`
	Variant *string                       `json:"variant,omitempty"`
}

// ResolutionDetailsObject_Value defines model for ResolutionDetailsObject.Value.
type ResolutionDetailsObject_Value struct {
	AdditionalProperties map[string]interface{} `json:"-"`
}

// ResolutionDetailsString defines model for resolutionDetailsString.
type ResolutionDetailsString struct {
	Reason  *string `json:"reason,omitempty"`
	Value   string  `json:"value"`
	Variant *string `json:"variant,omitempty"`
}

// ResolutionDetailsWithError defines model for resolutionDetailsWithError.
type ResolutionDetailsWithError struct {
	ErrorCode *string `json:"errorCode,omitempty"`
	Reason    *string `json:"reason,omitempty"`
	Variant   *string `json:"variant,omitempty"`
}

// FlagKey defines model for flagKey.
type FlagKey = string

// N400 defines model for 400.
type N400 = ResolutionDetailsWithError

// N404 defines model for 404.
type N404 = ResolutionDetailsWithError

// N500 defines model for 500.
type N500 = ResolutionDetailsWithError

// ResolveBooleanJSONBody defines parameters for ResolveBoolean.
type ResolveBooleanJSONBody = Context

// ResolveBooleanParams defines parameters for ResolveBoolean.
type ResolveBooleanParams struct {
	// The value that will be resolved in case of any error, or if the flag is not defined in the flag management system.
	DefaultValue bool `form:"default-value" json:"default-value"`
}

// ResolveNumberJSONBody defines parameters for ResolveNumber.
type ResolveNumberJSONBody = Context

// ResolveNumberParams defines parameters for ResolveNumber.
type ResolveNumberParams struct {
	// The value that will be resolved in case of any error, or if the flag is not defined in the flag management system.
	DefaultValue float32 `form:"default-value" json:"default-value"`
}

// ResolveObjectJSONBody defines parameters for ResolveObject.
type ResolveObjectJSONBody = Context

// ResolveObjectParams_DefaultValue defines parameters for ResolveObject.
type ResolveObjectParams_DefaultValue struct {
	AdditionalProperties map[string]interface{} `json:"-"`
}

// ResolveObjectParams defines parameters for ResolveObject.
type ResolveObjectParams struct {
	// The value that will be resolved in case of any error, or if the flag is not defined in the flag management system.
	DefaultValue ResolveObjectParams_DefaultValue `form:"default-value" json:"default-value"`
}

// ResolveStringJSONBody defines parameters for ResolveString.
type ResolveStringJSONBody = Context

// ResolveStringParams defines parameters for ResolveString.
type ResolveStringParams struct {
	// The value that will be resolved in case of any error, or if the flag is not defined in the flag management system.
	DefaultValue string `form:"default-value" json:"default-value"`
}

// ResolveBooleanJSONRequestBody defines body for ResolveBoolean for application/json ContentType.
type ResolveBooleanJSONRequestBody = ResolveBooleanJSONBody

// ResolveNumberJSONRequestBody defines body for ResolveNumber for application/json ContentType.
type ResolveNumberJSONRequestBody = ResolveNumberJSONBody

// ResolveObjectJSONRequestBody defines body for ResolveObject for application/json ContentType.
type ResolveObjectJSONRequestBody = ResolveObjectJSONBody

// ResolveStringJSONRequestBody defines body for ResolveString for application/json ContentType.
type ResolveStringJSONRequestBody = ResolveStringJSONBody

// Getter for additional properties for ResolveObjectParams_DefaultValue. Returns the specified
// element and whether it was found
func (a ResolveObjectParams_DefaultValue) Get(fieldName string) (value interface{}, found bool) {
	if a.AdditionalProperties != nil {
		value, found = a.AdditionalProperties[fieldName]
	}
	return
}

// Setter for additional properties for ResolveObjectParams_DefaultValue
func (a *ResolveObjectParams_DefaultValue) Set(fieldName string, value interface{}) {
	if a.AdditionalProperties == nil {
		a.AdditionalProperties = make(map[string]interface{})
	}
	a.AdditionalProperties[fieldName] = value
}

// Override default JSON handling for ResolveObjectParams_DefaultValue to handle AdditionalProperties
func (a *ResolveObjectParams_DefaultValue) UnmarshalJSON(b []byte) error {
	object := make(map[string]json.RawMessage)
	err := json.Unmarshal(b, &object)
	if err != nil {
		return err
	}

	if len(object) != 0 {
		a.AdditionalProperties = make(map[string]interface{})
		for fieldName, fieldBuf := range object {
			var fieldVal interface{}
			err := json.Unmarshal(fieldBuf, &fieldVal)
			if err != nil {
				return fmt.Errorf("error unmarshaling field %s: %w", fieldName, err)
			}
			a.AdditionalProperties[fieldName] = fieldVal
		}
	}
	return nil
}

// Override default JSON handling for ResolveObjectParams_DefaultValue to handle AdditionalProperties
func (a ResolveObjectParams_DefaultValue) MarshalJSON() ([]byte, error) {
	var err error
	object := make(map[string]json.RawMessage)

	for fieldName, field := range a.AdditionalProperties {
		object[fieldName], err = json.Marshal(field)
		if err != nil {
			return nil, fmt.Errorf("error marshaling '%s': %w", fieldName, err)
		}
	}
	return json.Marshal(object)
}

// Getter for additional properties for Context. Returns the specified
// element and whether it was found
func (a Context) Get(fieldName string) (value interface{}, found bool) {
	if a.AdditionalProperties != nil {
		value, found = a.AdditionalProperties[fieldName]
	}
	return
}

// Setter for additional properties for Context
func (a *Context) Set(fieldName string, value interface{}) {
	if a.AdditionalProperties == nil {
		a.AdditionalProperties = make(map[string]interface{})
	}
	a.AdditionalProperties[fieldName] = value
}

// Override default JSON handling for Context to handle AdditionalProperties
func (a *Context) UnmarshalJSON(b []byte) error {
	object := make(map[string]json.RawMessage)
	err := json.Unmarshal(b, &object)
	if err != nil {
		return err
	}

	if raw, found := object["targetingKey"]; found {
		err = json.Unmarshal(raw, &a.TargetingKey)
		if err != nil {
			return fmt.Errorf("error reading 'targetingKey': %w", err)
		}
		delete(object, "targetingKey")
	}

	if len(object) != 0 {
		a.AdditionalProperties = make(map[string]interface{})
		for fieldName, fieldBuf := range object {
			var fieldVal interface{}
			err := json.Unmarshal(fieldBuf, &fieldVal)
			if err != nil {
				return fmt.Errorf("error unmarshaling field %s: %w", fieldName, err)
			}
			a.AdditionalProperties[fieldName] = fieldVal
		}
	}
	return nil
}

// Override default JSON handling for Context to handle AdditionalProperties
func (a Context) MarshalJSON() ([]byte, error) {
	var err error
	object := make(map[string]json.RawMessage)

	if a.TargetingKey != nil {
		object["targetingKey"], err = json.Marshal(a.TargetingKey)
		if err != nil {
			return nil, fmt.Errorf("error marshaling 'targetingKey': %w", err)
		}
	}

	for fieldName, field := range a.AdditionalProperties {
		object[fieldName], err = json.Marshal(field)
		if err != nil {
			return nil, fmt.Errorf("error marshaling '%s': %w", fieldName, err)
		}
	}
	return json.Marshal(object)
}

// Getter for additional properties for ResolutionDetailsObject_Value. Returns the specified
// element and whether it was found
func (a ResolutionDetailsObject_Value) Get(fieldName string) (value interface{}, found bool) {
	if a.AdditionalProperties != nil {
		value, found = a.AdditionalProperties[fieldName]
	}
	return
}

// Setter for additional properties for ResolutionDetailsObject_Value
func (a *ResolutionDetailsObject_Value) Set(fieldName string, value interface{}) {
	if a.AdditionalProperties == nil {
		a.AdditionalProperties = make(map[string]interface{})
	}
	a.AdditionalProperties[fieldName] = value
}

// Override default JSON handling for ResolutionDetailsObject_Value to handle AdditionalProperties
func (a *ResolutionDetailsObject_Value) UnmarshalJSON(b []byte) error {
	object := make(map[string]json.RawMessage)
	err := json.Unmarshal(b, &object)
	if err != nil {
		return err
	}

	if len(object) != 0 {
		a.AdditionalProperties = make(map[string]interface{})
		for fieldName, fieldBuf := range object {
			var fieldVal interface{}
			err := json.Unmarshal(fieldBuf, &fieldVal)
			if err != nil {
				return fmt.Errorf("error unmarshaling field %s: %w", fieldName, err)
			}
			a.AdditionalProperties[fieldName] = fieldVal
		}
	}
	return nil
}

// Override default JSON handling for ResolutionDetailsObject_Value to handle AdditionalProperties
func (a ResolutionDetailsObject_Value) MarshalJSON() ([]byte, error) {
	var err error
	object := make(map[string]json.RawMessage)

	for fieldName, field := range a.AdditionalProperties {
		object[fieldName], err = json.Marshal(field)
		if err != nil {
			return nil, fmt.Errorf("error marshaling '%s': %w", fieldName, err)
		}
	}
	return json.Marshal(object)
}

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (POST /flags/{flag-key}/resolve/boolean)
	ResolveBoolean(w http.ResponseWriter, r *http.Request, flagKey FlagKey, params ResolveBooleanParams)

	// (POST /flags/{flag-key}/resolve/number)
	ResolveNumber(w http.ResponseWriter, r *http.Request, flagKey FlagKey, params ResolveNumberParams)

	// (POST /flags/{flag-key}/resolve/object)
	ResolveObject(w http.ResponseWriter, r *http.Request, flagKey FlagKey, params ResolveObjectParams)

	// (POST /flags/{flag-key}/resolve/string)
	ResolveString(w http.ResponseWriter, r *http.Request, flagKey FlagKey, params ResolveStringParams)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.HandlerFunc) http.HandlerFunc

// ResolveBoolean operation middleware
func (siw *ServerInterfaceWrapper) ResolveBoolean(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "flag-key" -------------
	var flagKey FlagKey

	err = runtime.BindStyledParameter("simple", false, "flag-key", chi.URLParam(r, "flag-key"), &flagKey)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "flag-key", Err: err})
		return
	}

	// Parameter object where we will unmarshal all parameters from the context
	var params ResolveBooleanParams

	// ------------- Required query parameter "default-value" -------------
	if paramValue := r.URL.Query().Get("default-value"); paramValue != "" {

	} else {
		siw.ErrorHandlerFunc(w, r, &RequiredParamError{ParamName: "default-value"})
		return
	}

	err = runtime.BindQueryParameter("form", true, true, "default-value", r.URL.Query(), &params.DefaultValue)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "default-value", Err: err})
		return
	}

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.ResolveBoolean(w, r, flagKey, params)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

// ResolveNumber operation middleware
func (siw *ServerInterfaceWrapper) ResolveNumber(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "flag-key" -------------
	var flagKey FlagKey

	err = runtime.BindStyledParameter("simple", false, "flag-key", chi.URLParam(r, "flag-key"), &flagKey)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "flag-key", Err: err})
		return
	}

	// Parameter object where we will unmarshal all parameters from the context
	var params ResolveNumberParams

	// ------------- Required query parameter "default-value" -------------
	if paramValue := r.URL.Query().Get("default-value"); paramValue != "" {

	} else {
		siw.ErrorHandlerFunc(w, r, &RequiredParamError{ParamName: "default-value"})
		return
	}

	err = runtime.BindQueryParameter("form", true, true, "default-value", r.URL.Query(), &params.DefaultValue)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "default-value", Err: err})
		return
	}

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.ResolveNumber(w, r, flagKey, params)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

// ResolveObject operation middleware
func (siw *ServerInterfaceWrapper) ResolveObject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "flag-key" -------------
	var flagKey FlagKey

	err = runtime.BindStyledParameter("simple", false, "flag-key", chi.URLParam(r, "flag-key"), &flagKey)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "flag-key", Err: err})
		return
	}

	// Parameter object where we will unmarshal all parameters from the context
	var params ResolveObjectParams

	// ------------- Required query parameter "default-value" -------------
	if paramValue := r.URL.Query().Get("default-value"); paramValue != "" {

	} else {
		siw.ErrorHandlerFunc(w, r, &RequiredParamError{ParamName: "default-value"})
		return
	}

	err = runtime.BindQueryParameter("form", false, true, "default-value", r.URL.Query(), &params.DefaultValue)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "default-value", Err: err})
		return
	}

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.ResolveObject(w, r, flagKey, params)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

// ResolveString operation middleware
func (siw *ServerInterfaceWrapper) ResolveString(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "flag-key" -------------
	var flagKey FlagKey

	err = runtime.BindStyledParameter("simple", false, "flag-key", chi.URLParam(r, "flag-key"), &flagKey)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "flag-key", Err: err})
		return
	}

	// Parameter object where we will unmarshal all parameters from the context
	var params ResolveStringParams

	// ------------- Required query parameter "default-value" -------------
	if paramValue := r.URL.Query().Get("default-value"); paramValue != "" {

	} else {
		siw.ErrorHandlerFunc(w, r, &RequiredParamError{ParamName: "default-value"})
		return
	}

	err = runtime.BindQueryParameter("form", true, true, "default-value", r.URL.Query(), &params.DefaultValue)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "default-value", Err: err})
		return
	}

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.ResolveString(w, r, flagKey, params)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshalingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshalingParamError) Error() string {
	return fmt.Sprintf("Error unmarshaling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshalingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{})
}

type ChiServerOptions struct {
	BaseURL          string
	BaseRouter       chi.Router
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, r chi.Router) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseRouter: r,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, r chi.Router, baseURL string) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseURL:    baseURL,
		BaseRouter: r,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options ChiServerOptions) http.Handler {
	r := options.BaseRouter

	if r == nil {
		r = chi.NewRouter()
	}
	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}

	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/flags/{flag-key}/resolve/boolean", wrapper.ResolveBoolean)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/flags/{flag-key}/resolve/number", wrapper.ResolveNumber)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/flags/{flag-key}/resolve/object", wrapper.ResolveObject)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/flags/{flag-key}/resolve/string", wrapper.ResolveString)
	})

	return r
}
