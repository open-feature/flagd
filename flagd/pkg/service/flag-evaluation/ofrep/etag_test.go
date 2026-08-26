package ofrep

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func serveWithETag(t *testing.T, status int, body, ifNoneMatch string) *httptest.ResponseRecorder {
	t.Helper()

	var fields []string
	if ifNoneMatch != "" {
		fields = []string{ifNoneMatch}
	}

	return serveWithETagFields(t, status, body, fields)
}

// serveWithETagFields sends one If-None-Match header line per field, which is how a client spreads
// its validators over repeated fields instead of one list.
func serveWithETagFields(t *testing.T, status int, body string, ifNoneMatch []string) *httptest.ResponseRecorder {
	t.Helper()

	handler := conditionalETag(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))

	req := httptest.NewRequest(http.MethodPost, "/ofrep/v1/evaluate/flags", nil)
	for _, field := range ifNoneMatch {
		req.Header.Add("If-None-Match", field)
	}

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)
	return recorder
}

func TestConditionalETag_HashesWrittenBytes(t *testing.T) {
	const body = `{"flags":[{"key":"a","value":true}],"metadata":{}}`

	first := serveWithETag(t, http.StatusOK, body, "")
	require.Equal(t, http.StatusOK, first.Code)
	assert.Equal(t, body, first.Body.String(), "the body must pass through untouched")

	etag := first.Header().Get("ETag")
	require.NotEmpty(t, etag)
	assert.Equal(t, quoteETag(bodyETag([]byte(body))), etag, "the tag must be the hash of the written bytes")

	t.Run("identical bytes yield the same tag", func(t *testing.T) {
		assert.Equal(t, etag, serveWithETag(t, http.StatusOK, body, "").Header().Get("ETag"))
	})

	t.Run("any byte difference yields a new tag", func(t *testing.T) {
		changed := serveWithETag(t, http.StatusOK, `{"flags":[{"key":"a","value":false}],"metadata":{}}`, "")
		assert.NotEqual(t, etag, changed.Header().Get("ETag"))
	})
}

func TestConditionalETag_NotModified(t *testing.T) {
	const body = `{"flags":[]}`
	etag := serveWithETag(t, http.StatusOK, body, "").Header().Get("ETag")
	require.NotEmpty(t, etag)

	recorder := serveWithETag(t, http.StatusOK, body, etag)

	assert.Equal(t, http.StatusNotModified, recorder.Code)
	assert.Equal(t, etag, recorder.Header().Get("ETag"), "a 304 must still name the representation")
	assert.Empty(t, recorder.Body.String(), "a 304 must not carry a body")
	assert.Empty(t, recorder.Header().Get("Content-Type"), "a 304 must not describe a body it does not have")
}

func TestConditionalETag_UnquotedValidatorMatches(t *testing.T) {
	const body = `{"flags":[]}`
	etag := serveWithETag(t, http.StatusOK, body, "").Header().Get("ETag")

	recorder := serveWithETag(t, http.StatusOK, body, normalizeETag(etag))
	assert.Equal(t, http.StatusNotModified, recorder.Code)
}

func TestConditionalETag_StaleValidatorIsServedFresh(t *testing.T) {
	recorder := serveWithETag(t, http.StatusOK, `{"flags":[]}`, `"stale"`)

	require.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, `{"flags":[]}`, recorder.Body.String())
	assert.NotEqual(t, `"stale"`, recorder.Header().Get("ETag"))
}

// If-None-Match is a list, not a single value: any member matching under weak comparison is a hit.
func TestConditionalETag_ValidatorListForms(t *testing.T) {
	const body = `{"flags":[]}`
	current := serveWithETag(t, http.StatusOK, body, "").Header().Get("ETag")
	require.NotEmpty(t, current)

	for name, tests := range map[string]struct {
		ifNoneMatch string
		want        int
	}{
		"wildcard":                     {`*`, http.StatusNotModified},
		"current tag last in list":     {`"stale", ` + current, http.StatusNotModified},
		"current tag first in list":    {current + `, "stale"`, http.StatusNotModified},
		"weak current tag":             {`W/` + current, http.StatusNotModified},
		"weak current tag in list":     {`W/"stale", W/` + current, http.StatusNotModified},
		"unquoted current tag in list": {`"stale",` + normalizeETag(current), http.StatusNotModified},
		"empty list entries":           {`, ,` + current, http.StatusNotModified},
		"only stale tags":              {`"stale", W/"staler"`, http.StatusOK},
		"empty list":                   {`,`, http.StatusOK},
		// net/http stops scanning at the first entry that is not a valid entity tag, which would
		// let junk hide the tag beside it. This scan compares the junk and carries on.
		"malformed entry beside a valid tag": {`bogus, ` + current, http.StatusNotModified},
		"unterminated quoted tag":            {`"` + normalizeETag(current), http.StatusNotModified},
		// The wildcard has to be the whole entry, which is stricter than net/http.
		"wildcard with trailing junk": {`*junk`, http.StatusOK},
		// A comma inside a quoted tag belongs to the tag, so this is one stale validator rather
		// than a list whose members happen to bracket the current tag.
		"comma inside a quoted tag": {`"` + normalizeETag(current) + `,` + normalizeETag(current) + `"`, http.StatusOK},
	} {
		t.Run(name, func(t *testing.T) {
			recorder := serveWithETag(t, http.StatusOK, body, tests.ifNoneMatch)

			assert.Equal(t, tests.want, recorder.Code)
			assert.Equal(t, current, recorder.Header().Get("ETag"), "the tag names the representation either way")
		})
	}
}

// Repeated header fields carry the same meaning as one comma-joined list.
func TestConditionalETag_RepeatedValidatorFields(t *testing.T) {
	const body = `{"flags":[]}`
	current := serveWithETag(t, http.StatusOK, body, "").Header().Get("ETag")
	require.NotEmpty(t, current)

	recorder := serveWithETagFields(t, http.StatusOK, body, []string{`"stale"`, current})
	assert.Equal(t, http.StatusNotModified, recorder.Code)

	stale := serveWithETagFields(t, http.StatusOK, body, []string{`"stale"`, `"staler"`})
	assert.Equal(t, http.StatusOK, stale.Code)
}

// A stale validator must not 304 an error body.
func TestConditionalETag_NonOKResponsesAreNotTagged(t *testing.T) {
	for _, status := range []int{
		http.StatusBadRequest,
		http.StatusInternalServerError,
		http.StatusRequestEntityTooLarge,
	} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			body := `{"errorDetails":"nope"}`
			recorder := serveWithETag(t, status, body, quoteETag(bodyETag([]byte(body))))

			assert.Equal(t, status, recorder.Code)
			assert.Empty(t, recorder.Header().Get("ETag"))
			assert.Equal(t, body, recorder.Body.String())
		})
	}
}

func TestConditionalETag_ImplicitOK(t *testing.T) {
	handler := conditionalETag(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"flags":[]}`))
	}))

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodPost, "/ofrep/v1/evaluate/flags", nil))

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.NotEmpty(t, recorder.Header().Get("ETag"))
	assert.Equal(t, `{"flags":[]}`, recorder.Body.String())
}

func TestSplitETagList(t *testing.T) {
	for name, tests := range map[string]struct {
		list string
		want []string
	}{
		"empty":               {"", nil},
		"single tag":          {`"abc"`, []string{`"abc"`}},
		"list":                {`"abc", "def"`, []string{`"abc"`, `"def"`}},
		"list without spaces": {`"abc","def"`, []string{`"abc"`, `"def"`}},
		"weak tags":           {`W/"abc", W/"def"`, []string{`W/"abc"`, `W/"def"`}},
		"wildcard":            {`*`, []string{`*`}},
		"unquoted tag":        {`abc`, []string{`abc`}},
		"empty entries":       {`, "abc" , ,`, []string{`"abc"`}},
		"only separators":     {` , `, nil},
		"comma inside a tag":  {`"a,b"`, []string{`"a,b"`}},
		"comma inside a tag in a list": {
			`"a,b", "c"`, []string{`"a,b"`, `"c"`},
		},
	} {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tests.want, splitETagList(tests.list))
		})
	}
}

// An empty 200 is still a representation, so it gets a validator.
func TestConditionalETag_EmptyBody(t *testing.T) {
	recorder := serveWithETag(t, http.StatusOK, "", "")

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, quoteETag(bodyETag(nil)), recorder.Header().Get("ETag"))
}
