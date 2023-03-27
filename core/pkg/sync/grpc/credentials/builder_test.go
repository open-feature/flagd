package credentials

import (
	"os"
	"testing"
)

const sampleCert = `-----BEGIN CERTIFICATE-----
MIIEnDCCAoQCCQCHcl3hGXwRQzANBgkqhkiG9w0BAQsFADAQMQ4wDAYDVQQDDAVm
bGFnZDAeFw0yMzAyMTAxODM1NDVaFw0zMzAyMDcxODM1NDVaMBAxDjAMBgNVBAMM
BWZsYWdkMIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAwDLEAUti/kG9
MhJLtO7oAy7diHxWKDFmsIHrE+z2IzTxjXxVHQLv1HiYB/UN75y7qlb3MwvzSc+C
BoLuoiM0PDiMio9/o9X5j0U+v3H1JpUU5LardkvsprFqJWmHF+D7aRdM0LBLn2X6
HQOhSnPyH9Qjl2l2tyPiPTZ6g0i2+rXZsNUoTs4fm6ThhZ0LeXR8KDmCTun3ze1d
hXA7ydxwILH2OVc+Wnzl30+BRvOiLQbc9nYnwSREFeIy8sFbhrTHqSNn3eY79ssZ
T6f4tN3jEV1d7NqoFk9KFLJKJhMt7smMB9NLwVWi581Zj1krYirNlP6mtmPrn3kJ
lsgT15kFftShMVcYFSHqOSLiy4SspHGK8KJaFoEVx0wp/weRwrWXi6vWg7tuHATH
fw7gW/9CyV+ylc0pJ002wtPAgzJYUaOrna0R2r3yQsSzRcDnqsm4FLkPHLoyjrwQ
vshKcEqjhGml1M+lTDEo3RO5ZoQ3ZN2AZKPDrK2zGG4wFJjHRu9FtutOEZkYYOzA
emTQWW8US3q8WVQqGl/EwQqzXk9Lco7uhLdXmqVOvAi6z01gehQJPnjhH7iqAPVp
1tlOBHit1F3sTAQIO/2zff3LCKiD2d27KINh4aFEyDbDmglPA8VPO3BMQVSjFlxj
K1s2G1IDBixXK76VmBP+ZpvxOaQtYIUCAwEAATANBgkqhkiG9w0BAQsFAAOCAgEA
K9+wnl5gpkfNBa+OSxlhOn3CKhcaW/SWZ4aLw2yK1NZNnNjpwUcLQScUDBKDoJJR
5roc3PIImX7hdnobZWqFhD23laaAlu5XLk9P7n51uMEiNjQQc2WaaBZDTRJfki1C
MvPskXqptgPsVyuPJc0DxfaCz7pDYjq/CtJ+osaj404P5mlO1QJ8W91QSx+aq2x4
uUTUWuyr/8flIcxiX0o8VTb2LcUvWZBMGa3CdeLnPHrOjovfjJFy0Ysk3SGEACLL
9mpbNbv23v9UXVfyFffHpyzvyUJIOsNXG0O1AYf5t9bukqHolGR/RQUN4yGd3M62
mFR5bOST36DjNSzTrx1eyCLv22+h9VVlWFPrebFnq1W5SSi8PtsGSMjhvX7dB1kS
t0yJtlj2HwBAvI1zVKG76q6neSU51UXFQUbO0OA0sxjicEOlNfXnShM/kY2lobpX
hrCysWpqoSS0S3UBvmuRiraLWkP1KueC0XHoAi8yuwMAdM6Y+h2OJpnO0PdpUmrp
lAqdxbyICnB1Nsm5QGGm6Pxd8lEbQ9ZSwFjgqApjT2zVhuaaUC7jdlEP1H5snt9n
8FQR06lrzGyW04ud9pd6MXJup1oghAlvnzXioAH2Az0IXcHvqUGZQattFv27OXqj
QZ6ayNO119SNscvC6Qe9GLlbBEHDQWKPiftnS2Mh6Do=
-----END CERTIFICATE-----`

func TestCredentialBuilder_Build(t *testing.T) {
	// "insecure" is a hardcoded term at insecure.NewCredentials
	const insecure = "insecure"
	// "tls" is a hardcoded term at tlsCreds.Info
	const tls = "tls"
	// local test file with valid certificate
	const validCertFile = "valid.cert"
	// local test file with invalid certificate
	const invalidCertFile = "invalid.cert"

	// init cert files for tests & cleanup with a deffer
	err := os.WriteFile(validCertFile, []byte(sampleCert), 0o600)
	if err != nil {
		t.Errorf("error creating valid certificate file: %s", err)
	}

	err = os.WriteFile(invalidCertFile, []byte("--certificate--"), 0o600)
	if err != nil {
		t.Errorf("error creating invalid certificate file: %s", err)
	}

	defer func() {
		errV := os.Remove(validCertFile)
		errI := os.Remove(invalidCertFile)
		if errV != nil || errI != nil {
			t.Errorf("error removing cerificate files: %v, %v", errV, errI)
		}
	}()

	tests := []struct {
		name           string
		certPath       string
		secure         bool
		expectSecProto string
		error          bool
	}{
		{
			name:           "Insecure source results in insecure connection",
			secure:         false,
			certPath:       "",
			expectSecProto: insecure,
		},
		{
			name:           "Secure source results in secure connection",
			certPath:       validCertFile,
			secure:         true,
			expectSecProto: tls,
		},
		{
			name:           "Secure source with no certificate results in a secure connection",
			secure:         true,
			expectSecProto: tls,
		},
		{
			name:     "Invalid cert path results in an error",
			secure:   true,
			certPath: "invalid/path",
			error:    true,
		},
		{
			name:     "Invalid certificate results in an error",
			secure:   true,
			certPath: invalidCertFile,
			error:    true,
		},
		{
			name:     "Prevent insecure if certificate path is set - configuration check",
			secure:   false,
			certPath: validCertFile,
			error:    true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			builder := CredentialBuilder{}
			tCred, err := builder.Build(test.secure, test.certPath)

			if test.error {
				if err == nil {
					t.Errorf("test expected non error execution. But resulted in an error: %s", err.Error())
				}

				// Test expected an error. Nothing to validate further
				return
			}

			// check for errors to be certain
			if err != nil {
				t.Errorf("unexpected error: %s", err.Error())
			}

			protoc := tCred.Info().SecurityProtocol
			if protoc != test.expectSecProto {
				t.Errorf("buildTransportCredentials() returned protocol= %v, want %v", protoc, test.expectSecProto)
			}
		})
	}
}
