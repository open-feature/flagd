package integration_test

import "flag"

var (
	tls      string
	certPath string
)

func init() {
	flag.StringVar(&tls, "tls", "false", "tls enabled for testing")
	flag.StringVar(&certPath, "cert-path", "", "path to cert to use in tls tests")
}
