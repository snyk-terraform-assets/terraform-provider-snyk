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

package organization

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	snyk_http "github.com/snyk-terraform-assets/terraform-provider-snyk/internal/http"
)

const VERSION = "2023-09-20"

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type ClientConfig struct {
	HTTPClient  HTTPClient
	URL         string
	Token       string
	BearerToken string
	Version     string
}

type Client struct {
	httpClient    HTTPClient
	url           string
	authorization string
	version       string
}

func NewClient(url string, token string) (*Client, error) {

	httpClient, err := snyk_http.NewClient(
		snyk_http.WithExtraCertificates(os.Getenv("NODE_EXTRA_CA_CERTS")),
	)
	if err != nil {
		return nil, err
	}

	client, err := newClient(ClientConfig{
		HTTPClient: httpClient,
		URL:        url,
		Token:      token,
		Version:    VERSION,
	})

	if err != nil {
		return nil, err
	} else {
		return client, nil
	}
}

func newClient(config ClientConfig) (*Client, error) {
	httpClient := config.HTTPClient

	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	if config.URL == "" {
		if env, ok := os.LookupEnv("SNYK_API"); ok {
			config.URL = env
		} else {
			return nil, fmt.Errorf("no URL provided")
		}
	}

	if config.Token == "" && config.BearerToken == "" {
		if env, ok := os.LookupEnv("SNYK_TOKEN"); ok {
			config.Token = env
		} else {
			return nil, fmt.Errorf("no token provided")
		}
	}

	if config.Version == "" {
		return nil, fmt.Errorf("no version provided")
	}

	parsedURL, err := url.Parse(config.URL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %v", err)
	}

	sanitizedURL := url.URL{
		Scheme: parsedURL.Scheme,
		Host:   parsedURL.Host,
	}

	var authzHeader string
	if config.BearerToken != "" {
		authzHeader = fmt.Sprintf("Bearer %s", config.BearerToken)
	} else {
		authzHeader = fmt.Sprintf("token %s", config.Token)
	}

	client := Client{
		httpClient:    httpClient,
		url:           sanitizedURL.String(),
		authorization: authzHeader,
		version:       config.Version,
	}

	return &client, nil
}
