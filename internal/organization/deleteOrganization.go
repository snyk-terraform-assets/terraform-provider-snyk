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
	"fmt"
	"net/http"
)

func (c *Client) DeleteOrganization(ctx context.Context, orgID string) (e error) {
	url := fmt.Sprintf("%s/v1/org/%s", c.url, orgID)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	query := req.URL.Query()
	req.URL.RawQuery = query.Encode()

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", c.authorization)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		if err := res.Body.Close(); err != nil && e == nil {
			e = err
		}
	}()

	if res.StatusCode != http.StatusNoContent {
		return fmt.Errorf("invalid status code: %v", res.StatusCode)
	}

	return nil
}
