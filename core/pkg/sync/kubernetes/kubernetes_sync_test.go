package kubernetes

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	stdsync "sync"
	"testing"
	"time"

	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/sync"
	"github.com/open-feature/open-feature-operator/api/core/v1beta1"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic/fake"
	"k8s.io/client-go/kubernetes/scheme"
	testing2 "k8s.io/client-go/testing"
	"k8s.io/client-go/tools/cache"
)

var Metadata = v1.TypeMeta{
	Kind:       "FeatureFlag",
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
	validFFCfg := v1beta1.FeatureFlag{
		TypeMeta: Metadata,
	}

	tests := []struct {
		name    string
		input   interface{}
		want    *v1beta1.FeatureFlag
		wantErr bool
	}{
		{
			name:  "Simple success",
			input: toUnstructured(t, validFFCfg),
			want: &v1beta1.FeatureFlag{
				TypeMeta: Metadata,
			},
			wantErr: false,
		},
		{
			name: "Tombstone unwraps",
			input: cache.DeletedFinalStateUnknown{
				Obj: toUnstructured(t, validFFCfg),
			},
			want: &v1beta1.FeatureFlag{
				TypeMeta: Metadata,
			},
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

	validFFCfg := v1beta1.FeatureFlag{
		TypeMeta: Metadata,
		ObjectMeta: v1.ObjectMeta{
			Namespace: cfgNs,
			Name:      cfgName,
		},
	}

	type args struct {
		obj    interface{}
		object types.NamespacedName
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
				object: types.NamespacedName{
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
				object: types.NamespacedName{
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
				obj: toUnstructured(t, v1beta1.FeatureFlag{
					TypeMeta: v1.TypeMeta{
						Kind:       "FeatureFlag",
						APIVersion: "someAPIVersion",
					},
				}),
				object: types.NamespacedName{
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

	validFFCfgOld := v1beta1.FeatureFlag{
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
		object types.NamespacedName
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
				object: types.NamespacedName{
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
				object: types.NamespacedName{
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
				object: types.NamespacedName{
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
				newObj: toUnstructured(t, v1beta1.FeatureFlag{
					TypeMeta: v1.TypeMeta{
						Kind:       "FeatureFlag",
						APIVersion: "someAPIVersion",
					},
				}),
				object: types.NamespacedName{
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
				oldObj: toUnstructured(t, v1beta1.FeatureFlag{
					TypeMeta: v1.TypeMeta{
						Kind:       "FeatureFlag",
						APIVersion: "someAPIVersion",
					},
				}),
				newObj: toUnstructured(t, validFFCfgNew),
				object: types.NamespacedName{
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
	flagSpec := v1beta1.FlagSpec{
		Flags: v1beta1.Flags{FlagsMap: map[string]v1beta1.Flag{}},
	}

	validCfg := v1beta1.FeatureFlag{
		TypeMeta: Metadata,
		ObjectMeta: v1.ObjectMeta{
			Namespace:       "resourceNS",
			Name:            "resourceName",
			ResourceVersion: "v1",
		},
		Spec: v1beta1.FeatureFlagSpec{
			FlagSpec: flagSpec,
		},
	}

	type args struct {
		ClientResponse v1beta1.FeatureFlag
		ClientError    error
	}

	tests := []struct {
		name          string
		args          args
		injectionFunc func(s *Sync)
		want          string
		wantErr       bool
	}{
		{
			name: "Scenario - get from informer cache",
			injectionFunc: func(s *Sync) {
				s.URI = fmt.Sprintf("%s/%s", validCfg.Namespace, validCfg.Name)
				s.dynamicClient = fake.NewSimpleDynamicClient(scheme.Scheme, &validCfg)
				err := s.Init(context.Background())
				require.Nil(t, err)
			},
			wantErr: false,
			want:    `{"flags":{}}`,
		},
		{
			name: "Scenario - get from API if informer cache miss",
			injectionFunc: func(s *Sync) {
				s.URI = fmt.Sprintf("%s/%s", validCfg.Namespace, validCfg.Name)
				s.dynamicClient = fake.NewSimpleDynamicClient(scheme.Scheme, &validCfg)
				err := s.Init(context.Background())
				require.Nil(t, err)
			},
			wantErr: false,
			want:    `{"flags":{}}`,
		},
		{
			name: "Scenario - error for informer cache read error",
			injectionFunc: func(s *Sync) {
				s.URI = fmt.Sprintf("%s/%s", validCfg.Namespace, validCfg.Name)
				s.informer = &MockInformer{
					fakeStore: cache.FakeCustomStore{
						GetByKeyFunc: func(key string) (item interface{}, exists bool, err error) {
							return nil, false, errors.New("mock error")
						},
					},
				}
			},
			wantErr: true,
		},
		{
			name: "Scenario - error for API get error",
			injectionFunc: func(s *Sync) {
				s.URI = fmt.Sprintf("%s/%s", validCfg.Namespace, validCfg.Name)
				fakeClient := fake.NewSimpleDynamicClient(scheme.Scheme, &validCfg)
				fakeClient.PrependReactor("get", "featureflags", func(action testing2.Action) (handled bool, ret runtime.Object, err error) {
					return true, nil, errors.New("could not get ff config")
				})
				s.dynamicClient = fakeClient
				err := s.Init(context.Background())
				require.Nil(t, err)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup with args
			k := &Sync{
				logger: logger.NewLogger(nil, false),
			}

			if tt.injectionFunc != nil {
				tt.injectionFunc(k)
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
	flagSpec := v1beta1.FeatureFlagSpec{
		FlagSpec: v1beta1.FlagSpec{
			Flags: v1beta1.Flags{FlagsMap: map[string]v1beta1.Flag{}},
		},
	}

	validCfg := v1beta1.FeatureFlag{
		TypeMeta: Metadata,
		ObjectMeta: v1.ObjectMeta{
			Namespace:       "resourceNS",
			Name:            "resourceName",
			ResourceVersion: "v1",
		},
		Spec: flagSpec,
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
			want: `{"flags":{}}`,
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
			want: `{"flags":{}}`,
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
				logger: logger.NewLogger(nil, false),
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

func TestInit(t *testing.T) {
	t.Run("expect error with wrong URI format", func(t *testing.T) {
		k := Sync{URI: ""}
		e := k.Init(context.TODO())
		if e == nil {
			t.Errorf("Expected error but got none")
		}
		if k.IsReady() {
			t.Errorf("Expected NOT to be ready")
		}
	})
	t.Run("expect informer registration", func(t *testing.T) {
		const name = "myFF"
		const ns = "myNS"
		scheme := runtime.NewScheme()
		ff := &unstructured.Unstructured{}
		ff.SetUnstructuredContent(getCFG(name, ns))
		fakeClient := fake.NewSimpleDynamicClient(scheme, ff)
		k := Sync{
			URI:           fmt.Sprintf("%s/%s", ns, name),
			dynamicClient: fakeClient,
			namespace:     ns,
		}
		e := k.Init(context.TODO())
		if e != nil {
			t.Errorf("Unexpected error: %v", e)
		}
		if k.informer == nil {
			t.Errorf("Informer not initialized")
		}
		if k.IsReady() {
			t.Errorf("The Sync should not be ready")
		}
	})
}

func TestSync_ReSync(t *testing.T) {
	const name = "myFF"
	const ns = "myNS"
	const payload = "{\"flags\":null}"
	s := runtime.NewScheme()
	ff := &unstructured.Unstructured{}
	ff.SetUnstructuredContent(getCFG(name, ns))
	fakeDynamicClient := fake.NewSimpleDynamicClient(s, ff)

	l, err := logger.NewZapLogger(zapcore.FatalLevel, "console")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	tests := []struct {
		name     string
		k        *Sync
		countMsg int
		async    bool
	}{
		{
			name: "Happy Path",
			k: &Sync{
				URI:           fmt.Sprintf("%s/%s", ns, name),
				dynamicClient: fakeDynamicClient,
				namespace:     ns,
				logger:        logger.NewLogger(l, true),
			},
			countMsg: 2, // one for sync and one for resync
			async:    true,
		},
		{
			name: "CRD not found",
			k: &Sync{
				URI:           fmt.Sprintf("doesnt%s/exist%s", ns, name),
				dynamicClient: fakeDynamicClient,
				namespace:     ns,
				logger:        logger.NewLogger(l, true),
			},
			countMsg: 0,
			async:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := tt.k.Init(context.TODO())
			if e != nil {
				t.Errorf("Unexpected error: %v", e)
			}
			if tt.k.IsReady() {
				t.Errorf("The Sync should not be ready")
			}
			dataChannel := make(chan sync.DataSync, tt.countMsg)
			ctx, cancel := context.WithCancel(context.Background())
			if tt.async {
				var wg stdsync.WaitGroup
				wg.Add(1)
				go func() {
					defer wg.Done()
					if err := tt.k.Sync(ctx, dataChannel); err != nil {
						t.Errorf("Unexpected error: %v", err)
					}
				}()

				if err := tt.k.ReSync(ctx, dataChannel); err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				for i := tt.countMsg; i > 0; i-- {
					select {
					case d := <-dataChannel:
						if d.FlagData != payload {
							t.Errorf("Expected %v, got %v", payload, d.FlagData)
						}
					case <-time.After(time.Second):
						t.Fatalf("timeout waiting for data")
					}
				}

				cancel()
				wg.Wait()
			} else {
				defer cancel()
				if err := tt.k.Sync(context.TODO(), dataChannel); !strings.Contains(err.Error(), "not found") {
					t.Errorf("Unexpected error: %v", err)
				}
				if err := tt.k.ReSync(context.TODO(), dataChannel); !strings.Contains(err.Error(), "not found") {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

func TestNotify(t *testing.T) {
	const name = "myFF"
	const ns = "myNS"
	// Use scheme.Scheme so the fake client knows the list kind for featureflags.
	// Do NOT pre-populate the fake client: client-go v0.33+ uses RealFIFO which
	// only delivers events via the watch stream, not via the initial List/Replace.
	// Pre-populating causes the object to silently land in the store without
	// triggering AddFunc, so later UpdateStatus calls are seen as updates on an
	// object the informer has never "added" — leaving the test waiting forever.
	if err := v1beta1.AddToScheme(scheme.Scheme); err != nil {
		t.Fatalf("failed to add v1beta1 to scheme: %v", err)
	}
	fc := fake.NewSimpleDynamicClient(scheme.Scheme)
	l, err := logger.NewZapLogger(zapcore.FatalLevel, "console")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	k := Sync{
		URI:           fmt.Sprintf("%s/%s", ns, name),
		dynamicClient: fc,
		namespace:     ns,
		logger:        logger.NewLogger(l, true),
	}
	err = k.Init(context.TODO())
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if k.informer == nil {
		t.Errorf("Informer not initialized")
	}
	c := make(chan INotify)
	go func() { k.notify(context.TODO(), c) }()

	if k.IsReady() {
		t.Errorf("The Sync should not be ready")
	}

	// wait for informer callbacks to be set
	msg := <-c
	if msg.GetEvent().EventType != DefaultEventTypeReady {
		t.Errorf("Expected message %v, got %v", DefaultEventTypeReady, msg)
	}
	// create: explicitly create the object via the watch stream so AddFunc fires
	ff := &unstructured.Unstructured{}
	cfg := getCFG(name, ns)
	ff.SetUnstructuredContent(cfg)
	_, err = fc.Resource(featureFlagResource).Namespace(ns).Create(context.TODO(), ff, v1.CreateOptions{})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	msg = <-c
	if msg.GetEvent().EventType != DefaultEventTypeCreate {
		t.Errorf("Expected message %v, got %v", DefaultEventTypeCreate, msg)
	}
	// update
	old := cfg["metadata"].(map[string]interface{})
	old["resourceVersion"] = "newVersion"
	cfg["metadata"] = old
	ff.SetUnstructuredContent(cfg)
	_, err = fc.Resource(featureFlagResource).Namespace(ns).Update(context.TODO(), ff, v1.UpdateOptions{})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	msg = <-c
	if msg.GetEvent().EventType != DefaultEventTypeModify {
		t.Errorf("Expected message %v, got %v", DefaultEventTypeModify, msg)
	}
	// delete
	err = fc.Resource(featureFlagResource).Namespace(ns).Delete(context.TODO(), name, v1.DeleteOptions{})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	msg = <-c
	if msg.GetEvent().EventType != DefaultEventTypeDelete {
		t.Errorf("Expected message %v, got %v", DefaultEventTypeDelete, msg)
	}

	// validate we don't crash parsing wrong spec
	cfg["spec"] = map[string]interface{}{
		"featureFlagSpec": int64(12), // we expect string here
	}
	ff.SetUnstructuredContent(cfg)
	_, err = fc.Resource(featureFlagResource).Namespace(ns).Create(context.TODO(), ff, v1.CreateOptions{})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	cfg["status"] = map[string]interface{}{
		"bump": "1",
	}
	ff.SetUnstructuredContent(cfg)
	_, err = fc.Resource(featureFlagResource).Namespace(ns).Update(context.TODO(), ff, v1.UpdateOptions{})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	err = fc.Resource(featureFlagResource).Namespace(ns).Delete(context.TODO(), name, v1.DeleteOptions{})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func Test_NewK8sSync(t *testing.T) {
	l, err := logger.NewZapLogger(zapcore.FatalLevel, "console")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	const uri = "myURI"
	log := logger.NewLogger(l, true)
	dc := fake.NewSimpleDynamicClient(runtime.NewScheme())
	k := NewK8sSync(
		log,
		uri,
		dc,
	)
	if k == nil {
		t.Errorf("Object not initialized properly")
	}
	if k.URI != uri {
		t.Errorf("Object not initialized with the right URI")
	}
	if k.logger != log {
		t.Errorf("Object not initialized with the right logger")
	}
	if k.dynamicClient != dc {
		t.Errorf("Object not initialized with the right K8s dynamic client")
	}
}

func getCFG(name, namespace string) map[string]interface{} {
	return map[string]interface{}{
		"apiVersion": "core.openfeature.dev/v1beta1",
		"kind":       "FeatureFlag",
		"metadata": map[string]interface{}{
			"name":      name,
			"namespace": namespace,
		},
		"spec": map[string]interface{}{
			"flagSpec": map[string]interface{}{
				"flags": nil,
			},
		},
	}
}

// toUnstructured helper to convert an interface to unstructured.Unstructured
func toUnstructured(t *testing.T, obj interface{}) *unstructured.Unstructured {
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

// MockInformer contains an embedded controllertest.FakeInformer for desired method overriding
type MockInformer struct {
	cache.SharedInformer
	fake.FakeDynamicClient
	fakeStore cache.FakeCustomStore
}

func (m *MockInformer) GetStore() cache.Store {
	return &m.fakeStore
}

func TestMarshallFeatureFlagSpec(t *testing.T) {
	res, err := marshallFeatureFlagSpec(&v1beta1.FeatureFlag{
		Spec: v1beta1.FeatureFlagSpec{
			FlagSpec: v1beta1.FlagSpec{
				Flags: v1beta1.Flags{
					FlagsMap: map[string]v1beta1.Flag{
						"flag": {
							DefaultVariant: "kubernetes",
						},
					},
				},
			},
		},
	})

	require.NoError(t, err)
	require.JSONEq(t, `{"flags":{"flag":{"state":"","variants":null,"defaultVariant":"kubernetes"}}}`, res)
}

func TestMarshallFeatureFlagSpec_metadata(t *testing.T) {
	res, err := marshallFeatureFlagSpec(&v1beta1.FeatureFlag{
		Spec: v1beta1.FeatureFlagSpec{
			FlagSpec: v1beta1.FlagSpec{
				Metadata: json.RawMessage(`{"scope":"app-1"}`),
				Flags: v1beta1.Flags{
					FlagsMap: map[string]v1beta1.Flag{
						"flag-with-meta": {
							State:          "ENABLED",
							DefaultVariant: "on",
							Variants:       json.RawMessage(`{"on":true,"off":false}`),
							Metadata:       json.RawMessage(`{"owner":"team-x"}`),
						},
					},
				},
			},
		},
	})

	require.NoError(t, err)
	require.JSONEq(t, `{
		"metadata": {"scope": "app-1"},
		"flags": {
			"flag-with-meta": {
				"state": "ENABLED",
				"variants": {"on": true, "off": false},
				"defaultVariant": "on",
				"metadata": {"owner": "team-x"}
			}
		}
	}`, res)
}
