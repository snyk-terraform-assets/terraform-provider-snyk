// Copyright 2023 Snyk Limited All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package http

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"

	"github.com/hashicorp/go-cleanhttp"
)

type config struct {
	tlsSkipVerify bool
	certificates  []string
}

// Option is a configuration option for the HTTP client.
type Option func(c *config)

// WithTLSSkipVerify enables or disables the verification of the server TLS
// certificates. This option defaults to false.
func WithTLSSkipVerify(tlsSkipVerify bool) Option {
	return func(c *config) {
		c.tlsSkipVerify = tlsSkipVerify
	}
}

// WithExtraCertificates adds more certificates to the pool of certificates
// trusted by this client. path is the path to a PEM file containing one or more
// certificates. If path is empty or the file at path does not exist, no
// certificates are added to the pool.
func WithExtraCertificates(path string) Option {
	return func(c *config) {
		c.certificates = append(c.certificates, path)
	}
}

func NewClient(options ...Option) (*http.Client, error) {
	var c config

	for _, o := range options {
		o(&c)
	}

	pool, err := x509.SystemCertPool()
	if err != nil {
		return nil, fmt.Errorf("read system cert pool: %v", err)
	}
	if pool == nil {
		pool = x509.NewCertPool()
	}

	if err := loadCertificates(pool, c.certificates); err != nil {
		return nil, fmt.Errorf("load certificates: %v", err)
	}

	client := cleanhttp.DefaultClient()

	if transport, ok := client.Transport.(*http.Transport); ok {
		transport.TLSClientConfig = &tls.Config{
			MinVersion:         tls.VersionTLS12,
			InsecureSkipVerify: c.tlsSkipVerify,
			RootCAs:            pool,
		}
	}

	return client, nil
}

func loadCertificates(pool *x509.CertPool, certificates []string) error {
	for _, certificate := range certificates {
		if err := loadCertificate(pool, certificate); err != nil {
			return fmt.Errorf("load certificate %v: %v", certificate, err)
		}
	}

	return nil
}

func loadCertificate(pool *x509.CertPool, certPath string) error {
	certData, err := os.ReadFile(certPath)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("read certificates: %v", err)
	}

	if ok := pool.AppendCertsFromPEM(certData); !ok {
		return fmt.Errorf("no certificates found from NODE_EXTRA_CA_CERTS")
	}

	return nil
}
