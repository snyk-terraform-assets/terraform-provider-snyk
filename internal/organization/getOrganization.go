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
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type OrganizationResponse struct {
	Jsonapi struct {
		Version string `json:"version"`
	} `json:"jsonapi"`
	Data struct {
		Type       string `json:"type"`
		ID         string `json:"id"`
		Attributes struct {
			Name       string `json:"name"`
			Slug       string `json:"slug"`
			IsPersonal string `json:"is_personal"`
			GroupID    string `json:"group_id"`
		} `json:"attributes"`
	} `json:"data"`
	Links struct {
		Self string `json:"self"`
	} `json:"links"`
}

type Organization struct {
	Name    string
	GroupId string
	ID      string
}

func (c *Client) GetOrganization(ctx context.Context, organizationID string) (org *Organization, e error) {

	url := fmt.Sprintf("%s/rest/orgs/%s", c.url, organizationID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
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

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid status code: %v", res.StatusCode)
	}

	var result OrganizationResponse

	body, _ := io.ReadAll(res.Body)
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	org = &Organization{Name: result.Data.Attributes.Name, GroupId: result.Data.Attributes.GroupID, ID: result.Data.ID}

	return org, nil
}
