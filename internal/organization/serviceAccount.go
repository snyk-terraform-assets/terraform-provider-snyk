package organization

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const SERVICEACCOUNTVERSION = "2022-09-15~experimental"

type ServiceAccountRequest struct {
	AccessTokenTTLSeconds int    `json:"access_token_ttl_seconds,omitempty"`
	AuthType              string `json:"auth_type"`
	JwksURL               string `json:"jwks_url,omitempty"`
	Name                  string `json:"name"`
	RoleID                string `json:"role_id"`
}

type ServiceAccountResponse struct {
	Data struct {
		Attributes struct {
			AccessTokenTTLSeconds int    `json:"access_token_ttl_seconds,omitempty"`
			AuthType              string `json:"auth_type"`
			ClientId              string `json:"client_id"`
			ApiKey                string `json:"api_key"`
			JwksURL               string `json:"jwks_url"`
			Name                  string `json:"name"`
			RoleID                string `json:"role_id"`
		} `json:"attributes"`
		ID    string `json:"id"`
		Links struct {
			First   string `json:"first,omitempty"`
			Last    string `json:"last,omitempty"`
			Next    string `json:"next,omitempty"`
			Prev    string `json:"prev,omitempty"`
			Related string `json:"related,omitempty"`
			Self    string `json:"self,omitempty"`
		} `json:"links,omitempty"`
		Type string `json:"type"`
	} `json:"data"`
	JsonAPI struct {
		Version string `json:"version,omitempty"`
	} `json:"jsonapi,omitempty"`
	Links struct {
		First   string `json:"first,omitempty"`
		Last    string `json:"last,omitempty"`
		Next    string `json:"next,omitempty"`
		Prev    string `json:"prev,omitempty"`
		Related string `json:"related,omitempty"`
		Self    string `json:"self,omitempty"`
	} `json:"links,omitempty"`
}

func (c *Client) CreateOrganizationServiceAccount(ctx context.Context, orgID string, request *ServiceAccountRequest) (sar *ServiceAccountResponse, e error) {
	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(request); err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/rest/orgs/%s/service_accounts", c.url, orgID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, &body)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	query.Set("version", SERVICEACCOUNTVERSION)
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
		var resp ServiceAccountResponse
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

func (c *Client) DeleteOrganizationServiceAccount(ctx context.Context, orgID, saID string) (e error) {
	url := fmt.Sprintf("%s/v1/orgs/%s/service_accounts/%s", c.url, orgID, saID)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	query := req.URL.Query()
	req.URL.RawQuery = query.Encode()

	req.Header.Set("Content-Type", "application/vnd.api+json")
	req.Header.Set("Authorization", SERVICEACCOUNTVERSION)

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
