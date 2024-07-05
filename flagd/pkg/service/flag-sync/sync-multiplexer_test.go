package sync

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRegistration(t *testing.T) {
	// given
	mux, err := NewMux(getSimpleFlagStore())
	if err != nil {
		t.Fatal("error during flag extraction")
		return
	}

	tests := []struct {
		testName            string
		id                  interface{}
		source              string
		flagStringValidator func(flagString string, testSource string, testName string)
		connection          chan payload
		expectError         bool
	}{
		{
			testName:   "subscribe to all flags",
			id:         context.Background(),
			connection: make(chan payload, 1),
		},
		{
			testName: "subscribe to source A",
			id:       context.Background(),
			source:   "A",
			flagStringValidator: func(flagString string, testSource string, testName string) {
				assert.Contains(t, flagString, fmt.Sprintf("\"source\":\"%s\"", testSource))
			},
			connection: make(chan payload, 1),
		},
		{
			testName: "subscribe to source B",
			id:       context.Background(),
			source:   "B",
			flagStringValidator: func(flagString string, testSource string, testName string) {
				assert.Contains(t, flagString, fmt.Sprintf("\"source\":\"%s\"", testSource))
			},
			connection: make(chan payload, 1),
		},

		{
			testName:   "subscribe to empty",
			id:         context.Background(),
			source:     "C",
			connection: make(chan payload, 1),
			flagStringValidator: func(flagString string, testSource string, testName string) {
				assert.Equal(t, flagString, string(emptyConfigBytes))
			},
			expectError: false,
		},
		{
			testName:    "subscribe to non-existing",
			id:          context.Background(),
			source:      "D",
			connection:  make(chan payload, 1),
			expectError: true,
		},
	}

	// validate registration
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			// when
			err := mux.Register(test.id, test.source, test.connection)

			// then
			if !test.expectError && err != nil {
				t.Fatal("expected no errors, but got error")
			}

			if test.expectError && err != nil {
				// pass
				return
			}

			// validate subscription
			var initSync payload
			select {
			case <-time.After(2 * time.Second):
				t.Fatal("data sync did not complete for initial sync within an acceptable timeframe")

			case initSync = <-test.connection:
				break
			}

			if test.flagStringValidator != nil {
				test.flagStringValidator(initSync.flags, test.source, test.testName)
			}
		})
	}
}

func TestUpdateAndRemoval(t *testing.T) {
	// given
	mux, err := NewMux(getSimpleFlagStore())
	if err != nil {
		t.Fatal("error during flag extraction")
		return
	}

	identifier := context.Background()
	channel := make(chan payload, 1)
	err = mux.Register(identifier, "", channel)
	if err != nil {
		t.Fatal("error during subscription registration")
		return
	}

	select {
	case <-time.After(2 * time.Second):
		t.Fatal("data sync did not complete for initial sync within an acceptable timeframe")
	case <-channel:
		break
	}

	// when - updates are triggered
	err = mux.Publish()
	if err != nil {
		t.Fatal("failure to trigger update request on multiplexer")
		return
	}

	// then
	select {
	case <-time.After(2 * time.Second):
		t.Fatal("data sync did not complete for initial sync within an acceptable timeframe")
	case <-channel:
		break
	}

	// when - subscription removed & update triggered
	mux.Unregister(identifier, "")
	err = mux.Publish()
	if err != nil {
		t.Fatal("failure to trigger update request on multiplexer")
		return
	}

	// then
	select {
	case <-time.After(2 * time.Second):
		break
	case <-channel:
		t.Fatal("expected no sync but got an update as removal was not performed")
	}
}

func TestGetAllFlags(t *testing.T) {
	// given
	mux, err := NewMux(getSimpleFlagStore())
	if err != nil {
		t.Fatal("error during flag extraction")
		return
	}

	// when - get all with open scope
	flags, err := mux.GetAllFlags("")
	if err != nil {
		t.Fatal("error when retrieving all flags")
		return
	}

	if len(flags) == 0 {
		t.Fatal("expected no empty flags")
		return
	}

	// when - get all with a scope
	flags, err = mux.GetAllFlags("A")
	if err != nil {
		t.Fatal("error when retrieving all flags")
		return
	}

	if len(flags) == 0 || !strings.Contains(flags, fmt.Sprintf("\"source\":\"%s\"", "A")) {
		t.Fatal("expected flags to be scoped")
		return
	}

	// when - get all for a flagless-scope
	flags, err = mux.GetAllFlags("C")
	if err != nil {
		t.Fatal("error when retrieving all flags")
		return
	}

	assert.Equal(t, flags, string(emptyConfigBytes))
}
