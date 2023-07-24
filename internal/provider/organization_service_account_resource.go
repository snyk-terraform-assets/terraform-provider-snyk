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

package provider

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/snyk-terraform-assets/terraform-provider-snyk/internal/organization"
	"github.com/snyk-terraform-assets/terraform-provider-snyk/internal/snykclient"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &OrganizationServiceAccountResource{}
var _ resource.ResourceWithImportState = &OrganizationServiceAccountResource{}

func NewOrganizationServiceAccountResource() resource.Resource {
	return &OrganizationServiceAccountResource{}
}

// OrganizationServiceAccountResource defines the resource implementation.
type OrganizationServiceAccountResource struct {
	client snykclient.Client
}

// OrganizationServiceAccountResourceModel describes the resource data model.
type OrganizationServiceAccountResourceModel struct {
	Id                    types.String `tfsdk:"id"`
	AccessTokenTTLSeconds types.Int64  `tfsdk:"access_token_ttl_seconds"`
	AuthType              types.String `tfsdk:"auth_type"`
	Name                  types.String `tfsdk:"name"`
	JWKSUrl               types.String `tfsdk:"jwks_url"`
	RoleId                types.String `tfsdk:"role_id"`
	OrganizationId        types.String `tfsdk:"organization_id"`
	ApiKey                types.String `tfsdk:"api_key"`
}

func (r *OrganizationServiceAccountResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization_service_account"
}

func (r *OrganizationServiceAccountResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Provides Snyk Organization level [Service Account](https://docs.snyk.io/enterprise-setup/service-accounts)",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Snyk OrganizationServiceAccount id",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"api_key": schema.StringAttribute{
				Computed:            true,
				Sensitive:           true,
				MarkdownDescription: "service account api key",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"access_token_ttl_seconds": schema.Int64Attribute{
				MarkdownDescription: "The time, in seconds, that a generated access token will be valid for. Defaults to 1 hour if unset. Only used when auth_type is oauth_private_key_jwt. Constraints: Min 3600|Max 86400",
				Optional:            true,
			},
			"auth_type": schema.StringAttribute{
				MarkdownDescription: "Authentication strategy fro the service account: api_key - Regulare Snyk API Key. oauth_private_key_jwt - OAuth2 client_credentials grant, using private_key_jwt client_assertion as laid out in OIDC Connect Core 1.0, section 9. Allowed: api_key|oauth_private_key_jwt",
				Required:            true,
			},
			"jwks_url": schema.StringAttribute{
				MarkdownDescription: "A JWKs URL hosting your public keys, used to verify signed JWT requests. Must be https. Required only when auth_type is oauth_private_key_jwt",
				Optional:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "A human-friendly name for the service account.",
				Required:            true,
			},
			"role_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the role which the created service account should us. Obtained in the Snyk UI, via \"Group Page\" -> \"Settings\" -> \"Member Roles\" -> \"Create new Role\". Can be shared among multiple accounts.",
				Required:            true,
			},
			"organization_id": schema.StringAttribute{
				MarkdownDescription: "The id of the organization to create the service account in.",
				Required:            true,
			},
		},
	}
}

func (r *OrganizationServiceAccountResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*snykclient.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *snykclient.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = *client
}

func (r *OrganizationServiceAccountResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// 500s on an account already existing
	var plan *OrganizationServiceAccountResourceModel
	// Read Terraform plan into the model
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := uuid.Parse(plan.OrganizationId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to parse OrganizationServiceAccountDataSource Group Guid, got error: %s", err))
		return
	}

	res, err := r.client.OrgClient.CreateOrganizationServiceAccount(ctx, plan.OrganizationId.ValueString(), &organization.ServiceAccountRequest{AccessTokenTTLSeconds: int(plan.AccessTokenTTLSeconds.ValueInt64()), AuthType: plan.AuthType.ValueString(), JwksURL: plan.JWKSUrl.ValueString(), Name: plan.Name.ValueString(), RoleID: plan.RoleId.ValueString()})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create OrganizationServiceAccount: %s, got error: %s", plan.Name.ValueString(), err))
		return
	} else {
		plan.Id = types.StringValue(res.Data.ID)
		plan.ApiKey = types.StringValue(res.Data.Attributes.ApiKey)
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *OrganizationServiceAccountResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// There is no get just passthrough
	var data *OrganizationServiceAccountResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *OrganizationServiceAccountResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Best that can be done is delete and make new
}

func (r *OrganizationServiceAccountResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *OrganizationServiceAccountResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.OrgClient.DeleteOrganizationServiceAccount(ctx, data.OrganizationId.ValueString(), data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Could not delete OrganizationServiceAccount: %s, got error: %s", data.Name.ValueString(), err))
		return
	}
}

func (r *OrganizationServiceAccountResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// This shouldnt work, which is kind of subpar
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
