package kubernetes

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/open-feature/flagd/internal/pkg/logger"
	"github.com/open-feature/flagd/internal/pkg/sync"
	"k8s.io/client-go/tools/cache"

	"sigs.k8s.io/controller-runtime/pkg/controller/controllertest"

	"github.com/open-feature/open-feature-operator/apis/core/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

var Metadata = v1.TypeMeta{
	Kind:       "FeatureFlagConfiguration",
	APIVersion: apiVersion,
}

func Test_parseURI(t *testing.T) {
	tests := []struct {
		name     string
		uri      string
		ns       string
		resource string
		err      bool
	}{
		{
			name:     "simple success",
			uri:      "namespace/resource",
			ns:       "namespace",
			resource: "resource",
			err:      false,
		},
		{
			name: "simple error - no ns",
			uri:  "/resource",
			err:  true,
		},
		{
			name: "simple error - no resource",
			uri:  "resource/",
			err:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ns, rs, err := parseURI(tt.uri)
			if (err != nil) != tt.err {
				t.Errorf("parseURI() error = %v, wantErr %v", err, tt.err)
				return
			}
			if ns != tt.ns {
				t.Errorf("parseURI() got = %v, want %v", ns, tt.ns)
			}
			if rs != tt.resource {
				t.Errorf("parseURI() got1 = %v, want %v", rs, tt.resource)
			}
		})
	}
}

func Test_toFFCfg(t *testing.T) {
	validFFCfg := v1alpha1.FeatureFlagConfiguration{
		TypeMeta: Metadata,
	}

	tests := []struct {
		name    string
		input   interface{}
		want    *v1alpha1.FeatureFlagConfiguration
		wantErr bool
	}{
		{
			name:    "Simple success",
			input:   toUnstructured(t, validFFCfg),
			want:    &validFFCfg,
			wantErr: false,
		},
		{
			name: "Simple error",
			input: struct {
				flag string
			}{
				flag: "test",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := toFFCfg(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("toFFCfg() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toFFCfg() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_commonHandler(t *testing.T) {
	cfgNs := "resourceNS"
	cfgName := "resourceName"

	validFFCfg := v1alpha1.FeatureFlagConfiguration{
		TypeMeta: Metadata,
		ObjectMeta: v1.ObjectMeta{
			Namespace: cfgNs,
			Name:      cfgName,
		},
	}

	type args struct {
		obj    interface{}
		object client.ObjectKey
	}
	tests := []struct {
		name      string
		args      args
		wantErr   bool
		wantEvent bool
		eventType DefaultEventType
	}{
		{
			name: "simple success",
			args: args{
				obj: toUnstructured(t, validFFCfg),
				object: client.ObjectKey{
					Namespace: cfgNs,
					Name:      cfgName,
				},
			},
			wantEvent: true,
			wantErr:   false,
		},
		{
			name: "simple scenario - only notify if resource name matches",
			args: args{
				obj: toUnstructured(t, validFFCfg),
				object: client.ObjectKey{
					Namespace: cfgNs,
					Name:      "SomeOtherResource",
				},
			},
			wantEvent: false,
			wantErr:   false,
		},
		{
			name: "simple error - API mismatch",
			args: args{
				obj: toUnstructured(t, v1alpha1.FeatureFlagConfiguration{
					TypeMeta: v1.TypeMeta{
						Kind:       "FeatureFlagConfiguration",
						APIVersion: "someAPIVersion",
					},
				}),
				object: client.ObjectKey{
					Namespace: cfgNs,
					Name:      cfgName,
				},
			},
			wantErr:   true,
			wantEvent: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			syncChan := make(chan INotify, 1)

			err := commonHandler(tt.args.obj, tt.args.object, tt.eventType, syncChan)
			if err != nil && !tt.wantErr {
				t.Errorf("commonHandler() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				if err == nil {
					t.Errorf("commonHandler() expected error but received none.")
				}

				// Expected error occurred, hence continue
				return
			}

			if tt.wantEvent != true {
				// Not interested in the event, hence ignore notification check. But check for chan writes
				if len(syncChan) != 0 {
					t.Errorf("commonHandler() expected no events, but events are available: %d", len(syncChan))
				}

				return
			}

			// watch events with a timeout
			var notify INotify
			select {
			case notify = <-syncChan:
			case <-time.After(2 * time.Second):
				t.Errorf("timedout waiting for events from commonHandler()")
			}

			if notify.GetEvent().EventType != tt.eventType {
				t.Errorf("commonHandler() event = %v, wanted %v", notify.GetEvent().EventType, DefaultEventTypeDelete)
			}
		})
	}
}

func Test_updateFuncHandler(t *testing.T) {
	cfgNs := "resourceNS"
	cfgName := "resourceName"

	validFFCfgOld := v1alpha1.FeatureFlagConfiguration{
		TypeMeta: Metadata,
		ObjectMeta: v1.ObjectMeta{
			Namespace:       cfgNs,
			Name:            cfgName,
			ResourceVersion: "v1",
		},
	}

	validFFCfgNew := validFFCfgOld
	validFFCfgNew.ResourceVersion = "v2"

	type args struct {
		oldObj interface{}
		newObj interface{}
		object client.ObjectKey
	}
	tests := []struct {
		name      string
		args      args
		wantErr   bool
		wantEvent bool
	}{
		{
			name: "Simple success",
			args: args{
				oldObj: toUnstructured(t, validFFCfgOld),
				newObj: toUnstructured(t, validFFCfgNew),
				object: client.ObjectKey{
					Namespace: cfgNs,
					Name:      cfgName,
				},
			},
			wantErr:   false,
			wantEvent: true,
		},
		{
			name: "Simple scenario - notify only if resource name match",
			args: args{
				oldObj: toUnstructured(t, validFFCfgOld),
				newObj: toUnstructured(t, validFFCfgNew),
				object: client.ObjectKey{
					Namespace: cfgNs,
					Name:      "SomeOtherResource",
				},
			},
			wantErr:   false,
			wantEvent: false,
		},
		{
			name: "Simple scenario - notify only if resource version is new",
			args: args{
				oldObj: toUnstructured(t, validFFCfgOld),
				newObj: toUnstructured(t, validFFCfgOld),
				object: client.ObjectKey{
					Namespace: cfgNs,
					Name:      "SomeOtherResource",
				},
			},
			wantErr:   false,
			wantEvent: false,
		},
		{
			name: "Simple error - API version mismatch new object",
			args: args{
				oldObj: toUnstructured(t, validFFCfgOld),
				newObj: toUnstructured(t, v1alpha1.FeatureFlagConfiguration{
					TypeMeta: v1.TypeMeta{
						Kind:       "FeatureFlagConfiguration",
						APIVersion: "someAPIVersion",
					},
				}),
				object: client.ObjectKey{
					Namespace: cfgNs,
					Name:      cfgName,
				},
			},
			wantErr:   true,
			wantEvent: false,
		},
		{
			name: "Simple error - API version mismatch old object",
			args: args{
				oldObj: toUnstructured(t, v1alpha1.FeatureFlagConfiguration{
					TypeMeta: v1.TypeMeta{
						Kind:       "FeatureFlagConfiguration",
						APIVersion: "someAPIVersion",
					},
				}),
				newObj: toUnstructured(t, validFFCfgNew),
				object: client.ObjectKey{
					Namespace: cfgNs,
					Name:      cfgName,
				},
			},
			wantErr:   true,
			wantEvent: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			syncChan := make(chan INotify, 1)

			err := updateFuncHandler(tt.args.oldObj, tt.args.newObj, tt.args.object, syncChan)
			if err != nil && !tt.wantErr {
				t.Errorf("updateFuncHandler() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				if err == nil {
					t.Errorf("updateFuncHandler() expected error but received none.")
				}

				// Expected error occurred, hence continue
				return
			}

			if tt.wantEvent != true {
				// Not interested in the event, hence ignore notification check. But check for chan writes
				if len(syncChan) != 0 {
					t.Errorf("updateFuncHandler() expected no events, but events are available: %d", len(syncChan))
				}

				return
			}

			// watch events with a timeout
			var notify INotify
			select {
			case notify = <-syncChan:
			case <-time.After(2 * time.Second):
				t.Errorf("timedout waiting for events from updateFuncHandler()")
			}

			if notify.GetEvent().EventType != DefaultEventTypeModify {
				t.Errorf("updateFuncHandler() event = %v, wanted %v", notify.GetEvent().EventType, DefaultEventTypeModify)
			}
		})
	}
}

func TestSync_fetch(t *testing.T) {
	flagSpec := "fakeFlagSpec"

	validCfg := v1alpha1.FeatureFlagConfiguration{
		TypeMeta: Metadata,
		ObjectMeta: v1.ObjectMeta{
			Namespace:       "resourceNS",
			Name:            "resourceName",
			ResourceVersion: "v1",
		},
		Spec: v1alpha1.FeatureFlagConfigurationSpec{
			FeatureFlagSpec: flagSpec,
		},
	}

	type args struct {
		InformerGetFunc func(key string) (item interface{}, exists bool, err error)
		ClientResponse  v1alpha1.FeatureFlagConfiguration
		ClientError     error
	}

	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Scenario - get from informer cache",
			args: args{
				InformerGetFunc: func(key string) (item interface{}, exists bool, err error) {
					return toUnstructured(t, validCfg), true, nil
				},
			},
			wantErr: false,
			want:    flagSpec,
		},
		{
			name: "Scenario - get from API if informer cache miss",
			args: args{
				InformerGetFunc: func(key string) (item interface{}, exists bool, err error) {
					return nil, false, nil
				},
				ClientResponse: validCfg,
			},
			wantErr: false,
			want:    flagSpec,
		},
		{
			name: "Scenario - error for informer cache read error",
			args: args{
				InformerGetFunc: func(key string) (item interface{}, exists bool, err error) {
					return nil, false, errors.New("mock error")
				},
			},
			wantErr: true,
		},
		{
			name: "Scenario - error for API get error",
			args: args{
				InformerGetFunc: func(key string) (item interface{}, exists bool, err error) {
					return nil, false, nil
				},
				ClientError: errors.New("mock error"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup with args
			k := &Sync{
				informer: &MockInformer{
					fakeStore: cache.FakeCustomStore{
						GetByKeyFunc: tt.args.InformerGetFunc,
					},
				},
				readClient: &MockClient{
					getResponse: tt.args.ClientResponse,
					clientErr:   tt.args.ClientError,
				},
				Logger: logger.NewLogger(nil, false),
			}

			// Test fetch
			got, err := k.fetch(context.Background())

			if (err != nil) != tt.wantErr {
				t.Errorf("fetch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("fetch() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSync_watcher(t *testing.T) {
	flagSpec := "fakeFlagSpec"

	validCfg := v1alpha1.FeatureFlagConfiguration{
		TypeMeta: Metadata,
		ObjectMeta: v1.ObjectMeta{
			Namespace:       "resourceNS",
			Name:            "resourceName",
			ResourceVersion: "v1",
		},
		Spec: v1alpha1.FeatureFlagConfigurationSpec{
			FeatureFlagSpec: flagSpec,
		},
	}

	type args struct {
		InformerGetFunc func(key string) (item interface{}, exists bool, err error)
		notification    INotify
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "scenario - create event",
			want: flagSpec,
			args: args{
				InformerGetFunc: func(key string) (item interface{}, exists bool, err error) {
					return toUnstructured(t, validCfg), true, nil
				},
				notification: &Notifier{
					Event: Event[DefaultEventType]{
						EventType: DefaultEventTypeCreate,
					},
				},
			},
		},
		{
			name: "scenario - modify event",
			want: flagSpec,
			args: args{
				InformerGetFunc: func(key string) (item interface{}, exists bool, err error) {
					return toUnstructured(t, validCfg), true, nil
				},
				notification: &Notifier{
					Event: Event[DefaultEventType]{
						EventType: DefaultEventTypeModify,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// setup sync
			k := &Sync{
				informer: &MockInformer{
					fakeStore: cache.FakeCustomStore{
						GetByKeyFunc: tt.args.InformerGetFunc,
					},
				},
				Logger: logger.NewLogger(nil, false),
			}

			// create communication channels with buffer to so that calls are non-blocking
			notifies := make(chan INotify, 1)
			dataSyncs := make(chan sync.DataSync, 1)

			// emit event
			notifies <- tt.args.notification

			tCtx, cFunc := context.WithTimeout(context.Background(), 2*time.Second)
			defer cFunc()

			// start watcher
			go k.watcher(tCtx, notifies, dataSyncs)

			// wait for data sync
			select {
			case <-tCtx.Done():
				t.Errorf("timeout waiting for the results")
			case dataSyncs := <-dataSyncs:
				if dataSyncs.FlagData != tt.want {
					t.Errorf("fetch() got = %v, want %v", dataSyncs.FlagData, tt.want)
				}
			}
		})
	}
}

// toUnstructured helper to convert an interface to unstructured.Unstructured
func toUnstructured(t *testing.T, obj interface{}) interface{} {
	bytes, err := json.Marshal(obj)
	if err != nil {
		t.Errorf("test setup faulure: %s", err.Error())
	}

	var res map[string]interface{}

	err = json.Unmarshal(bytes, &res)
	if err != nil {
		t.Errorf("test setup faulure: %s", err.Error())
	}

	return &unstructured.Unstructured{Object: res}
}

// Mock implementations

// MockClient contains an embedded client.Reader for desired method overriding
type MockClient struct {
	client.Reader
	clientErr error

	getResponse v1alpha1.FeatureFlagConfiguration
}

func (m MockClient) Get(ctx context.Context, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
	// return error if error is set
	if m.clientErr != nil {
		return m.clientErr
	}

	// else try returning response
	cfg, ok := obj.(*v1alpha1.FeatureFlagConfiguration)
	if !ok {
		return errors.New("must contain a pointer typed v1alpha1.FeatureFlagConfiguration")
	}

	*cfg = m.getResponse
	return nil
}

// MockInformer contains an embedded controllertest.FakeInformer for desired method overriding
type MockInformer struct {
	controllertest.FakeInformer

	fakeStore cache.FakeCustomStore
}

func (m MockInformer) GetStore() cache.Store {
	return &m.fakeStore
}
