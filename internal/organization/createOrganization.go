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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const KIND_AWS = "aws"
const KIND_AZURE = "azure"
const KIND_GOOGLE = "google"

type OrganizationRequest struct {
	Name        string `json:"name"`
	GroupId     string `json:"groupId,omitempty"`
	SourceOrgId string `json:"sourceOrgId,omitempty"`
}

type OrganizationResponseV1 struct {
	ID      string    `json:"id"`
	Name    string    `json:"name"`
	Slug    string    `json:"slug"`
	URL     string    `json:"url"`
	Created time.Time `json:"created"`
	Group   struct {
		Name string `json:"name"`
		ID   string `json:"id"`
	} `json:"group"`
}

func (c *Client) CreateOrganization(ctx context.Context, request *OrganizationRequest) (or *OrganizationResponseV1, e error) {
	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(request); err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/v1/org", c.url)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, &body)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	query.Set("version", c.version)
	req.URL.RawQuery = query.Encode()

	req.Header.Set("Content-Type", "application/json")
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
		var resp OrganizationResponseV1
		err = json.Unmarshal([]byte(bodyBytes), &resp)
		if err != nil {
			return nil, err
		}
		return &resp, nil
	} else {
		body, _ := io.ReadAll(res.Body)
		bodyString := string(body)
		return nil, fmt.Errorf("invalid status code: %s", bodyString)
	}
}
