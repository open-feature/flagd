package nameresolvers

import (
	"fmt"
	"strings"

	"google.golang.org/grpc/resolver"
)

const scheme = "envoy"

type envoyBuilder struct{}

// Build A custom NameResolver to resolve gRPC target uri for envoy in the
// format of.
//
// Custom URI Scheme:
//
// envoy://[proxy-agent-host]:[proxy-agent-port]/[service-name]
func (*envoyBuilder) Build(target resolver.Target,
	cc resolver.ClientConn, _ resolver.BuildOptions,
) (resolver.Resolver, error) {
	_, err := isValidTarget(target)
	if err != nil {
		return nil, err
	}

	r := &envoyResolver{
		target: target,
		cc:     cc,
	}
	r.start()
	return r, nil
}

func (*envoyBuilder) Scheme() string {
	return scheme
}

type envoyResolver struct {
	target resolver.Target
	cc     resolver.ClientConn
}

// Envoy NameResolver, will always override the authority with the specified authority i.e. URL.path and
// use the socketAddress i.e. Host:Port to connect.
func (r *envoyResolver) start() {
	addr := fmt.Sprintf("%s:%s", r.target.URL.Hostname(), r.target.URL.Port())
	err := r.cc.UpdateState(resolver.State{Addresses: []resolver.Address{{Addr: addr}}})
	if err != nil {
		return
	}
}

func (*envoyResolver) ResolveNow(resolver.ResolveNowOptions) {}

func (*envoyResolver) Close() {}

// Validate user specified target
//
// Sample target string:  envoy://localhost:9211/test.service
//
// return `true` if the target string used match the scheme and format
func isValidTarget(target resolver.Target) (bool, error) {
	// make sure and host and port not empty
	// used as resolver.Address
	if target.URL.Scheme != "envoy" || target.URL.Hostname() == "" || target.URL.Port() == "" {
		return false, fmt.Errorf("envoy-resolver: invalid scheme or missing host/port, target: %s",
			target)
	}

	// make sure the path is valid
	// used as :authority e.g. test.service
	path := target.Endpoint()
	if path == "" || strings.Contains(path, "/") {
		return false, fmt.Errorf("envoy-resolver: invalid path %s", path)
	}

	return true, nil
}

func init() {
	resolver.Register(&envoyBuilder{})
}
