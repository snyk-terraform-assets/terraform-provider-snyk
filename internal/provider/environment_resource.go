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
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/snyk-terraform-assets/terraform-provider-snyk/internal/cloudapi"
	"github.com/snyk-terraform-assets/terraform-provider-snyk/internal/snykclient"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &EnvironmentResource{}
var _ resource.ResourceWithImportState = &EnvironmentResource{}

func NewEnvironmentResource() resource.Resource {
	return &EnvironmentResource{}
}

// EnvironmentResource defines the resource implementation.
type EnvironmentResource struct {
	client snykclient.Client
}

// EnvironmentResourceModel describes the resource data model.
type EnvironmentResourceModel struct {
	Id             types.String                          `tfsdk:"id"`
	OrganizationId types.String                          `tfsdk:"organization_id"`
	Name           types.String                          `tfsdk:"name"`
	Kind           types.String                          `tfsdk:"kind"`
	Azure          *EnvironmentAzureConfigResourceModel  `tfsdk:"azure"`
	Google         *EnvironmentGoogleConfigResourceModel `tfsdk:"google"`
	Aws            *EnvironmentAWSConfigResourceModel    `tfsdk:"aws"`
}

type EnvironmentAzureConfigResourceModel struct {
	ApplicationId  types.String `tfsdk:"application_id"`
	SubscriptionId types.String `tfsdk:"subscription_id"`
	TenantId       types.String `tfsdk:"tenant_id"`
}

type EnvironmentGoogleConfigResourceModel struct {
	ProjectId           types.String `tfsdk:"project_id"`
	ServiceAccountEmail types.String `tfsdk:"service_account_email"`
}

type EnvironmentAWSConfigResourceModel struct {
	RoleArn types.String `tfsdk:"role_arn"`
}

func (r *EnvironmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_environment"
}

func (r *EnvironmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Provides Snyk [Cloud Environments](https://docs.snyk.io/scan-cloud-deployment/snyk-cloud/snyk-cloud-concepts#environments)",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "User assigned name",
				Optional:            true,
			},
			"kind": schema.StringAttribute{
				MarkdownDescription: "One of [aws,azure,google]",
				Required:            true,
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Snyk Environment ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"organization_id": schema.StringAttribute{
				MarkdownDescription: "Snyk Organization GUID",
				Required:            true,
			},
		},
		Blocks: map[string]schema.Block{
			"azure": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"application_id": schema.StringAttribute{
						Optional:    true,
						Description: "ID of the Azure app registration with permissions to scan",
					}, "subscription_id": schema.StringAttribute{
						Optional:    true,
						Description: "ID of the Azure subscription to be scanned",
					}, "tenant_id": schema.StringAttribute{
						Optional:    true,
						Description: "Azure Tenant (directory) ID",
					},
				},
			},
			"google": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"project_id": schema.StringAttribute{
						Optional:    true,
						Description: "Google project ID",
					}, "service_account_email": schema.StringAttribute{
						Optional:    true,
						Description: "Google service account email",
					},
				},
			},
			"aws": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"role_arn": schema.StringAttribute{
						Optional:    true,
						Description: "ARN of the AWS role created for Snyk Cloud",
					},
				},
			},
		},
	}
}

func (r *EnvironmentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *EnvironmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan *EnvironmentResourceModel
	// Read Terraform plan into the model
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := uuid.Parse(plan.OrganizationId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to parse EnvironmentDataSource Organization Guid, got error: %s", err))
		return
	}

	kind := plan.Kind.ValueString()
	kindDiags := r.checkKindConfiguration(kind, plan)
	resp.Diagnostics.Append(kindDiags...)
	if kindDiags.HasError() {
		return
	}

	request := r.prepareEnvironmentRequest(kind, plan)

	res, err := r.client.CloudapiClient.CreateEnvironment(ctx, plan.OrganizationId.ValueString(), request)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Environment, got error: %s", err))
		return
	} else {
		plan.Id = types.StringValue(res.Data.Id)
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *EnvironmentResource) prepareEnvironmentRequest(kind string, plan *EnvironmentResourceModel) *cloudapi.EnvironmentRequest {
	request := &cloudapi.EnvironmentRequest{Data: cloudapi.Data{Type: "",
		Attributes: cloudapi.Attributes{Kind: kind, Name: plan.Name.ValueString()}}}

	if kind == cloudapi.KIND_AWS {
		request.Data.Attributes.AwsOptions = &cloudapi.AwsOptions{}
		request.Data.Attributes.AwsOptions.RoleArn = plan.Aws.RoleArn.ValueString()
	} else if kind == cloudapi.KIND_GOOGLE {
		request.Data.Attributes.GoogleOptions = &cloudapi.GoogleOptions{}
		request.Data.Attributes.GoogleOptions.ServiceAccountEmail = plan.Google.ServiceAccountEmail.ValueString()
		request.Data.Attributes.GoogleOptions.ProjectId = plan.Google.ProjectId.ValueString()
	} else if kind == cloudapi.KIND_AZURE {
		request.Data.Attributes.AzureOptions = &cloudapi.AzureOptions{}
		request.Data.Attributes.AzureOptions.ApplicationId = plan.Azure.ApplicationId.ValueString()
		request.Data.Attributes.AzureOptions.TenantId = plan.Azure.TenantId.ValueString()
		request.Data.Attributes.AzureOptions.SubscriptionId = plan.Azure.SubscriptionId.ValueString()
	}
	return request
}

func (r *EnvironmentResource) checkKindConfiguration(kind string, plan *EnvironmentResourceModel) (diags diag.Diagnostics) {
	if kind == cloudapi.KIND_AWS {
		if strings.TrimSpace(plan.Aws.RoleArn.ValueString()) == "" {
			diags.AddError("Configuration Error", "Unable to read AWS role_arn. A valid AWS Role arn should be provided.")
			return
		}

		if plan.Google != nil || plan.Azure != nil {
			diags.AddError("Configuration Error", "Invalid configuration, Google and AWS configurations should be empty when using AWS")
			return
		}

	} else if kind == cloudapi.KIND_GOOGLE {
		if plan.Aws != nil || plan.Azure != nil {
			diags.AddError("Configuration Error", "Invalid configuration, Azure and AWS configurations should be empty when using Google")
			return
		}
		if strings.TrimSpace(plan.Google.ProjectId.ValueString()) == "" {
			diags.AddError("Configuration Error", "Unable to read Google project_id. A valid Google project_id should be provided.")
			return
		}
		if strings.TrimSpace(plan.Google.ServiceAccountEmail.ValueString()) == "" {
			diags.AddError("Configuration Error", "Unable to read Google service_account_email. A valid Google service_account_email  should be provided.")
			return
		}

	} else if kind == cloudapi.KIND_AZURE {
		if plan.Aws != nil || plan.Google != nil {
			diags.AddError("Configuration Error", "Invalid configuration, Google and AWS configurations should be empty when using Azure")
			return
		}

		if strings.TrimSpace(plan.Azure.ApplicationId.ValueString()) == "" {
			diags.AddError("Configuration Error", "Unable to read Azure application_id. A valid Azure application_id should be provided.")
			return
		}
		if strings.TrimSpace(plan.Azure.TenantId.ValueString()) == "" {
			diags.AddError("Configuration Error", "Unable to read Azure tenant_id. A valid Azure tenant_id  should be provided.")
			return
		}
		if strings.TrimSpace(plan.Azure.SubscriptionId.ValueString()) == "" {
			diags.AddError("Configuration Error", "Unable to read Azure subscription_id. A valid Azure subscription_id  should be provided.")
			return
		}
	} else {
		diags.AddError("Configuration Error", "Unable to parse Environment kind. Kind should be one of [aws,azure,google]")
		return
	}

	return
}

func (r *EnvironmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *EnvironmentResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	res, err := r.client.CloudapiClient.GetEnvironment(ctx, data.OrganizationId.ValueString(), data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get Environment, got error: %s", err))
		return
	}
	if !r.convertRemoteData2Local(data, res, resp.Diagnostics) {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *EnvironmentResource) convertRemoteData2Local(data *EnvironmentResourceModel, res *cloudapi.EnvironmentObject,
	diags diag.Diagnostics) bool {
	data.Name = types.StringValue(res.Attributes.Name)
	data.Kind = types.StringValue(res.Attributes.Kind)

	if res.Attributes.Kind == cloudapi.KIND_AWS {
		data.Aws = &EnvironmentAWSConfigResourceModel{}
		data.Aws.RoleArn = types.StringValue(res.Attributes.AwsOptions.RoleArn)
	} else if res.Attributes.Kind == cloudapi.KIND_GOOGLE {
		data.Google = &EnvironmentGoogleConfigResourceModel{}
		data.Google.ServiceAccountEmail = types.StringValue(res.Attributes.GoogleOptions.ServiceAccountEmail)
		data.Google.ProjectId = types.StringValue(res.Attributes.GoogleOptions.ProjectId)
	} else if res.Attributes.Kind == cloudapi.KIND_AZURE {
		data.Azure = &EnvironmentAzureConfigResourceModel{}
		data.Azure.TenantId = types.StringValue(res.Attributes.AzureOptions.TenantId)
		data.Azure.ApplicationId = types.StringValue(res.Attributes.AzureOptions.ApplicationId)
		data.Azure.SubscriptionId = types.StringValue(res.Attributes.AzureOptions.SubscriptionId)
	} else {
		diags.AddError("Update reading remote state", "Invalid kind, known kinds are [aws,azure,google]")
		return false
	}
	return true
}

func (r *EnvironmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan *EnvironmentResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	kind := plan.Kind.ValueString()
	resp.Diagnostics.Append(r.checkKindConfiguration(kind, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	request := r.prepareEnvironmentRequest(kind, plan)

	err := r.client.CloudapiClient.UpdateEnvironment(ctx, plan.OrganizationId.ValueString(), plan.Id.ValueString(), request)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update Environment, got error: %s", err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *EnvironmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *EnvironmentResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.CloudapiClient.DeleteEnvironment(ctx, data.OrganizationId.ValueString(), data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Could not delete Environment, got error: %s", err))
		return
	}
}

func (r *EnvironmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
