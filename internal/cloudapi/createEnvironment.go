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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const KIND_AWS = "aws"
const KIND_AZURE = "azure"
const KIND_GOOGLE = "google"

type EnvironmentRequest struct {
	Data Data `json:"data"`
}

type Data struct {
	Attributes Attributes `json:"attributes"`
	Type       string     `json:"type"`
	Id         string     `json:"id,omitempty"`
}

type Attributes struct {
	Kind          string         `json:"kind"`
	Name          string         `json:"name"`
	AwsOptions    *AwsOptions    `json:"-"`
	AzureOptions  *AzureOptions  `json:"-"`
	GoogleOptions *GoogleOptions `json:"-"`
	Options       interface{}    `json:"options,omitempty"`
}

type EnvironmentResponse = EnvironmentRequest

type AwsOptions struct {
	RoleArn string `json:"role_arn,omitempty"`
}
type AzureOptions struct {
	ApplicationId  string `json:"application_id,omitempty"`
	SubscriptionId string `json:"subscription_id,omitempty"`
	TenantId       string `json:"tenant_id,omitempty"`
}

type GoogleOptions struct {
	ProjectId           string `json:"project_id,omitempty"`
	ServiceAccountEmail string `json:"service_account_email,omitempty"`
}

func convertEnvRequestOptionsForMarshal(req *EnvironmentRequest) *EnvironmentRequest {
	if req.Data.Attributes.AwsOptions != nil {
		req.Data.Attributes.Options = req.Data.Attributes.AwsOptions
	} else if req.Data.Attributes.AzureOptions != nil {
		req.Data.Attributes.Options = req.Data.Attributes.AzureOptions
	} else if req.Data.Attributes.GoogleOptions != nil {
		req.Data.Attributes.Options = req.Data.Attributes.GoogleOptions
	}
	return req
}

func convertEnvRequestOptionsForUnMarshal(req *EnvironmentRequest) (*EnvironmentRequest, error) {
	data, err := json.Marshal(req.Data.Attributes.Options)
	if err != nil {
		return req, err
	}
	if req.Data.Attributes.Kind == KIND_AWS {
		options := AwsOptions{}
		err = json.Unmarshal(data, &options)
		if err != nil {
			return req, err
		}
		req.Data.Attributes.AwsOptions = &options
	} else if req.Data.Attributes.Kind == KIND_GOOGLE {
		options := GoogleOptions{}
		err = json.Unmarshal(data, &options)
		if err != nil {
			return req, err
		}
		req.Data.Attributes.GoogleOptions = &options
	} else if req.Data.Attributes.Kind == KIND_AZURE {
		options := AzureOptions{}
		err = json.Unmarshal(data, &options)
		if err != nil {
			return req, err
		}
		req.Data.Attributes.AzureOptions = &options
	}
	return req, nil
}

func (c *Client) CreateEnvironment(ctx context.Context, orgID string, request *EnvironmentRequest) (er *EnvironmentResponse, e error) {
	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(convertEnvRequestOptionsForMarshal(request)); err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/rest/orgs/%s/cloud/environments", c.url, orgID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, &body)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
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

	if res.StatusCode == http.StatusCreated { // what about http.StatusOK?
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		var resp EnvironmentResponse
		err = json.Unmarshal([]byte(bodyBytes), &resp)
		if err != nil {
			return nil, err
		}
		return convertEnvRequestOptionsForUnMarshal(&resp)
	} else {
		return nil, fmt.Errorf("invalid status code: %v", res.StatusCode)
	}

	return nil, err
}
