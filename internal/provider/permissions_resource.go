package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/snyk-terraform-assets/terraform-provider-snyk/internal/cloudapi"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &PermissionsResource{}

func NewPermissionsResource() resource.Resource {
	return &PermissionsResource{}
}

// PermissionsResource defines the data source implementation.
type PermissionsResource struct {
	client cloudapi.Client
}

// PermissionsResourceModel describes the data source data model.
type PermissionsResourceModel struct {
	OrganizationId types.String `tfsdk:"organization_id"`
	Platform       types.String `tfsdk:"platform"`
	Type           types.String `tfsdk:"type"`
	Data           types.String `tfsdk:"data"`
}

func (d *PermissionsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_permissions"
}

func (d *PermissionsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Permissions data source",

		Attributes: map[string]schema.Attribute{
			"organization_id": schema.StringAttribute{
				MarkdownDescription: "Snyk organization ID",
				Required:            true,
			},
			"platform": schema.StringAttribute{
				MarkdownDescription: "Platform, should be `aws`",
				Required:            true,
				Validators:          []validator.String{stringvalidator.OneOf("aws")},
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Type, should be `cf`",
				Required:            true,
				Validators:          []validator.String{stringvalidator.OneOf("cf")},
			},
			"data": schema.StringAttribute{
				MarkdownDescription: "Permission data, e.g. template body to create an AWS CloudFormation Stack",
				Computed:            true,
				Optional:            true,
			},
		},
	}
}

func (d *PermissionsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*cloudapi.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = *client
}

func (d *PermissionsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data PermissionsResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := d.recomputeData(ctx, &data); err != nil {
		resp.Diagnostics.AddError("client error", fmt.Sprintf("%s", err))
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PermissionsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

func (d *PermissionsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data PermissionsResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := d.recomputeData(ctx, &data); err != nil {
		resp.Diagnostics.AddError("client error", fmt.Sprintf("%s", err))
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PermissionsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (d *PermissionsResource) recomputeData(ctx context.Context, permissions *PermissionsResourceModel) error {
	organizationId := permissions.OrganizationId.ValueString()
	platform := permissions.Platform.ValueString()
	type_ := permissions.Type.ValueString()
	result, err := d.client.GetPermissions(ctx, organizationId, &cloudapi.PermissionsRequest{
		Platform: platform,
		Type:     type_,
	})
	if err != nil {
		return fmt.Errorf("unable to get permissions: %w", err)
	}
	permissions.Data = types.StringValue(result.Attributes.Data)
	return nil
}
