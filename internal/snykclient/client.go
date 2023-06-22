package snykclient

import (
	"github.com/snyk-terraform-assets/terraform-provider-snyk/internal/cloudapi"
	"github.com/snyk-terraform-assets/terraform-provider-snyk/internal/organization"
)

type Client struct {
	CloudapiClient *cloudapi.Client
	OrgClient      *organization.Client
}

func NewClient(url string, token string) (*Client, error) {
	cloudapiClient, err := cloudapi.NewClient(url, token)
	if err != nil {
		return nil, err
	}
	orgClient, err := organization.NewClient(url, token)
	if err != nil {
		return nil, err
	}

	return &Client{CloudapiClient: cloudapiClient, OrgClient: orgClient}, nil
}
