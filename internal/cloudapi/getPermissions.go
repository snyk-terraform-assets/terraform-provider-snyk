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

type (
	PermissionsRequest struct {
		Platform string `json:"platform"`
		Type     string `json:"type"`
	}

	DocumentPermissions struct {
		Data PermissionsObject `json:"data"`
	}

	PermissionsAttributes struct {
		Data string `json:"data"`
	}

	PermissionsObject struct {
		Type       string                 `json:"type"`
		Attributes *PermissionsAttributes `json:"attributes,omitempty"`
	}
)

func (c *Client) GetPermissions(ctx context.Context, orgID string, request *PermissionsRequest) (env *PermissionsObject, e error) {
	url := fmt.Sprintf("%s/rest/orgs/%s/cloud/permissions", c.url, orgID)
	version := "2022-04-13~experimental"
	requestJson := map[string]interface{}{
		"data": map[string]interface{}{
			"attributes": request,
			"type":       "permission",
		},
	}

	requestBytes, err := json.Marshal(requestJson)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(requestBytes))
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	query.Set("version", version)
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

	var result DocumentPermissions

	body, _ := io.ReadAll(res.Body)
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}
