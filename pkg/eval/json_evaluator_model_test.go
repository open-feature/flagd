package eval

import (
	"testing"

	"github.com/open-feature/flagd/pkg/logger"
	"github.com/stretchr/testify/require"
)

func TestMergeFlags(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		current    Flags
		new        Flags
		newSource  string
		want       Flags
		wantNotifs map[string]interface{}
	}{
		{
			name:       "both nil",
			current:    Flags{Flags: nil},
			new:        Flags{Flags: nil},
			want:       Flags{Flags: map[string]Flag{}},
			wantNotifs: map[string]interface{}{},
		},
		{
			name:       "both empty flags",
			current:    Flags{Flags: map[string]Flag{}},
			new:        Flags{Flags: map[string]Flag{}},
			want:       Flags{Flags: map[string]Flag{}},
			wantNotifs: map[string]interface{}{},
		},
		{
			name:       "empty current",
			current:    Flags{Flags: nil},
			new:        Flags{Flags: map[string]Flag{}},
			want:       Flags{Flags: map[string]Flag{}},
			wantNotifs: map[string]interface{}{},
		},
		{
			name:       "empty new",
			current:    Flags{Flags: map[string]Flag{}},
			new:        Flags{Flags: nil},
			want:       Flags{Flags: map[string]Flag{}},
			wantNotifs: map[string]interface{}{},
		},
		{
			name: "extra fields on each",
			current: Flags{Flags: map[string]Flag{
				"waka": {
					DefaultVariant: "off",
					Source:         "1",
				},
			}},
			new: Flags{Flags: map[string]Flag{
				"paka": {
					DefaultVariant: "on",
				},
			}},
			newSource: "2",
			want: Flags{Flags: map[string]Flag{
				"waka": {
					DefaultVariant: "off",
					Source:         "1",
				},
				"paka": {
					DefaultVariant: "on",
					Source:         "2",
				},
			}},
			wantNotifs: map[string]interface{}{
				"paka": map[string]interface{}{"type": "write", "source": "2"},
			},
		},
		{
			name: "override",
			current: Flags{Flags: map[string]Flag{
				"waka": {DefaultVariant: "off"},
			}},
			new: Flags{Flags: map[string]Flag{
				"waka": {DefaultVariant: "on"},
				"paka": {DefaultVariant: "on"},
			}},
			want: Flags{Flags: map[string]Flag{
				"waka": {DefaultVariant: "on"},
				"paka": {DefaultVariant: "on"},
			}},
			wantNotifs: map[string]interface{}{
				"waka": map[string]interface{}{"type": "update", "source": ""},
				"paka": map[string]interface{}{"type": "write", "source": ""},
			},
		},
		{
			name: "identical",
			current: Flags{Flags: map[string]Flag{
				"hello": {DefaultVariant: "off"},
			}},
			new: Flags{Flags: map[string]Flag{
				"hello": {DefaultVariant: "off"},
			}},
			want: Flags{Flags: map[string]Flag{
				"hello": {DefaultVariant: "off"},
			}},
			wantNotifs: map[string]interface{}{},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gotNotifs := tt.current.Merge(logger.NewLogger(nil, false), tt.newSource, tt.new)
			require.Equal(t, tt.want, tt.want)
			require.Equal(t, tt.wantNotifs, gotNotifs)
		})
	}
}

func TestFlags_Add(t *testing.T) {
	mockLogger := logger.NewLogger(nil, false)
	mockSource := "source"
	mockOverrideSource := "source-2"

	type request struct {
		source string
		flags  Flags
	}

	tests := []struct {
		name                     string
		storedState              Flags
		addRequest               request
		expectedState            Flags
		expectedNotificationKeys []string
	}{
		{
			name: "Add success",
			storedState: Flags{
				Flags: map[string]Flag{
					"A": {Source: mockSource},
				},
			},
			addRequest: request{
				source: mockSource,
				flags: Flags{
					Flags: map[string]Flag{
						"B": {Source: mockSource},
					},
				},
			},
			expectedState: Flags{
				Flags: map[string]Flag{
					"A": {Source: mockSource},
					"B": {Source: mockSource},
				},
			},
			expectedNotificationKeys: []string{"B"},
		},
		{
			name: "Add multiple success",
			storedState: Flags{
				Flags: map[string]Flag{
					"A": {Source: mockSource},
				},
			},
			addRequest: request{
				source: mockSource,
				flags: Flags{
					Flags: map[string]Flag{
						"B": {Source: mockSource},
						"C": {Source: mockSource},
					},
				},
			},
			expectedState: Flags{
				Flags: map[string]Flag{
					"A": {Source: mockSource},
					"B": {Source: mockSource},
					"C": {Source: mockSource},
				},
			},
			expectedNotificationKeys: []string{"B", "C"},
		},
		{
			name: "Add success - conflict and override",
			storedState: Flags{
				Flags: map[string]Flag{
					"A": {Source: mockSource},
				},
			},
			addRequest: request{
				source: mockOverrideSource,
				flags: Flags{
					Flags: map[string]Flag{
						"A": {Source: mockOverrideSource},
					},
				},
			},
			expectedState: Flags{
				Flags: map[string]Flag{
					"A": {Source: mockOverrideSource},
				},
			},
			expectedNotificationKeys: []string{"A"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			messages := tt.storedState.Add(mockLogger, tt.addRequest.source, tt.addRequest.flags)

			require.Equal(t, tt.storedState, tt.expectedState)

			for k := range messages {
				require.Containsf(t, tt.expectedNotificationKeys, k,
					"Message key %s not present in the expected key list", k)
			}
		})
	}
}

func TestFlags_Update(t *testing.T) {
	mockLogger := logger.NewLogger(nil, false)
	mockSource := "source"
	mockOverrideSource := "source-2"

	type request struct {
		source string
		flags  Flags
	}

	tests := []struct {
		name                     string
		storedState              Flags
		UpdateRequest            request
		expectedState            Flags
		expectedNotificationKeys []string
	}{
		{
			name: "Update success",
			storedState: Flags{
				Flags: map[string]Flag{
					"A": {Source: mockSource, DefaultVariant: "True"},
				},
			},
			UpdateRequest: request{
				source: mockSource,
				flags: Flags{
					Flags: map[string]Flag{
						"A": {Source: mockSource, DefaultVariant: "False"},
					},
				},
			},
			expectedState: Flags{
				Flags: map[string]Flag{
					"A": {Source: mockSource, DefaultVariant: "False"},
				},
			},
			expectedNotificationKeys: []string{"A"},
		},
		{
			name: "Update multiple success",
			storedState: Flags{
				Flags: map[string]Flag{
					"A": {Source: mockSource, DefaultVariant: "True"},
					"B": {Source: mockSource, DefaultVariant: "True"},
				},
			},
			UpdateRequest: request{
				source: mockSource,
				flags: Flags{
					Flags: map[string]Flag{
						"A": {Source: mockSource, DefaultVariant: "False"},
						"B": {Source: mockSource, DefaultVariant: "False"},
					},
				},
			},
			expectedState: Flags{
				Flags: map[string]Flag{
					"A": {Source: mockSource, DefaultVariant: "False"},
					"B": {Source: mockSource, DefaultVariant: "False"},
				},
			},
			expectedNotificationKeys: []string{"A", "B"},
		},
		{
			name: "Update success - conflict and override",
			storedState: Flags{
				Flags: map[string]Flag{
					"A": {Source: mockSource, DefaultVariant: "True"},
				},
			},
			UpdateRequest: request{
				source: mockOverrideSource,
				flags: Flags{
					Flags: map[string]Flag{
						"A": {Source: mockOverrideSource, DefaultVariant: "True"},
					},
				},
			},
			expectedState: Flags{
				Flags: map[string]Flag{
					"A": {Source: mockOverrideSource, DefaultVariant: "True"},
				},
			},
			expectedNotificationKeys: []string{"A"},
		},
		{
			name: "Update fail",
			storedState: Flags{
				Flags: map[string]Flag{
					"A": {Source: mockSource},
				},
			},
			UpdateRequest: request{
				source: mockSource,
				flags: Flags{
					Flags: map[string]Flag{
						"B": {Source: mockSource},
					},
				},
			},
			expectedState: Flags{
				Flags: map[string]Flag{
					"A": {Source: mockSource},
				},
			},
			expectedNotificationKeys: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			messages := tt.storedState.Update(mockLogger, tt.UpdateRequest.source, tt.UpdateRequest.flags)

			require.Equal(t, tt.storedState, tt.expectedState)

			for k := range messages {
				require.Containsf(t, tt.expectedNotificationKeys, k,
					"Message key %s not present in the expected key list", k)
			}
		})
	}
}

func TestFlags_Delete(t *testing.T) {
	mockLogger := logger.NewLogger(nil, false)
	mockSource := "source"
	mockSource2 := "source2"

	tests := []struct {
		name                     string
		storedState              Flags
		deleteRequest            Flags
		expectedState            Flags
		expectedNotificationKeys []string
	}{
		{
			name: "Remove success",
			storedState: Flags{
				Flags: map[string]Flag{
					"A": {Source: mockSource},
					"B": {Source: mockSource},
					"C": {Source: mockSource2},
				},
			},
			deleteRequest: Flags{
				Flags: map[string]Flag{
					"A": {Source: mockSource},
				},
			},
			expectedState: Flags{
				Flags: map[string]Flag{
					"B": {Source: mockSource},
					"C": {Source: mockSource2},
				},
			},
			expectedNotificationKeys: []string{"A"},
		},
		{
			name: "Nothing to remove",
			storedState: Flags{
				Flags: map[string]Flag{
					"A": {Source: mockSource},
					"B": {Source: mockSource},
					"C": {Source: mockSource2},
				},
			},
			deleteRequest: Flags{
				Flags: map[string]Flag{
					"D	": {Source: mockSource},
				},
			},
			expectedState: Flags{
				Flags: map[string]Flag{
					"A": {Source: mockSource},
					"B": {Source: mockSource},
					"C": {Source: mockSource2},
				},
			},
			expectedNotificationKeys: []string{},
		},
		{
			name: "Remove all",
			storedState: Flags{
				Flags: map[string]Flag{
					"A": {Source: mockSource},
					"B": {Source: mockSource},
					"C": {Source: mockSource2},
				},
			},
			deleteRequest: Flags{
				Flags: map[string]Flag{},
			},
			expectedState: Flags{
				Flags: map[string]Flag{
					"C": {Source: mockSource2},
				},
			},
			expectedNotificationKeys: []string{"A", "B"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			messages := tt.storedState.Delete(mockLogger, mockSource, tt.deleteRequest)

			require.Equal(t, tt.storedState, tt.expectedState)

			for k := range messages {
				require.Containsf(t, tt.expectedNotificationKeys, k,
					"Message key %s not present in the expected key list", k)
			}
		})
	}
}
