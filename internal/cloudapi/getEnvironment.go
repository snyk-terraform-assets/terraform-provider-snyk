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

package cloudapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type (
	DocumentEnv struct {
		Data EnvironmentObject `json:"data"`
	}
	CollectionDocumentEnvs struct {
		Data []EnvironmentObject `json:"data"`
	}
	EnvironmentAttributes struct {
		Name          string          `json:"name"`
		Options       json.RawMessage `json:"options,omitempty"`
		NativeID      string          `json:"native_id"`
		Properties    json.RawMessage `json:"properties,omitempty"`
		Kind          string          `json:"kind"`
		Revision      int             `json:"revision"`
		CreatedAt     string          `json:"created_at"`
		Status        string          `json:"status"`
		Error         string          `json:"error,omitempty"`
		UpdatedAt     string          `json:"updated_at,omitempty"`
		UpdatedBy     string          `json:"updated_by,omitempty"`
		AwsOptions    *AwsOptions     `json:"-"`
		AzureOptions  *AzureOptions   `json:"-"`
		GoogleOptions *GoogleOptions  `json:"-"`
	}

	EnvironmentObject struct {
		ID         string                 `json:"id,omitempty"`
		Type       string                 `json:"type"`
		Attributes *EnvironmentAttributes `json:"attributes,omitempty"`
	}
)

func prepareOptionsForUnMarshal(env *EnvironmentAttributes) (*EnvironmentAttributes, error) {
	data, err := json.Marshal(env.Options)
	if err != nil {
		return env, err
	}
	if env.Kind == KIND_AWS {
		options := AwsOptions{}
		err = json.Unmarshal(data, &options)
		if err != nil {
			return env, err
		}
		env.AwsOptions = &options
	} else if env.Kind == KIND_GOOGLE {
		options := GoogleOptions{}
		err = json.Unmarshal(data, &options)
		if err != nil {
			return env, err
		}
		env.GoogleOptions = &options
	} else if env.Kind == KIND_AZURE {
		options := AzureOptions{}
		err = json.Unmarshal(data, &options)
		if err != nil {
			return env, err
		}
		env.AzureOptions = &options
	}
	return env, nil
}

func (c *Client) GetEnvironment(ctx context.Context, orgID string, environmentID string) (env *EnvironmentObject, e error) {

	url := fmt.Sprintf("%s/rest/orgs/%s/cloud/environments?id=%s", c.url, orgID, environmentID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	//query.Set("id", snykCloudEnvironmentID)
	query.Set("version", c.version)
	req.URL.RawQuery = query.Encode()

	req.Header.Set("Content-Type", "application/vnd.api+json")
	req.Header.Set("Authorization", c.authorization)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := res.Body.Close(); err != nil && e == nil {
			e = err
		}
	}()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid status code: %v", res.StatusCode)
	}

	var result CollectionDocumentEnvs

	body, _ := io.ReadAll(res.Body)
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	envObject := &result.Data[0]

	_, err = prepareOptionsForUnMarshal(envObject.Attributes)
	if err != nil {
		return nil, err
	}

	return envObject, nil
}
