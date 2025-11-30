package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/radugheorghies/thetalake-terraform-provider/internal/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &UserResource{}
var _ resource.ResourceWithImportState = &UserResource{}

func NewUserResource() resource.Resource {
	return &UserResource{}
}

// UserResource defines the resource implementation.
type UserResource struct {
	client *client.Client
}

// UserResourceModel describes the resource data model.
type UserResourceModel struct {
	ID                   types.String `tfsdk:"id"`
	Name                 types.String `tfsdk:"name"`
	Email                types.String `tfsdk:"email"`
	Password             types.String `tfsdk:"password"`
	PasswordConfirmation types.String `tfsdk:"password_confirmation"`
	RoleID               types.Int64  `tfsdk:"role_id"`
	SearchID             types.Int64  `tfsdk:"search_id"`
}

func (r *UserResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (r *UserResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Theta Lake User Resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "User ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "User Name",
			},
			"email": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "User Email",
			},
			"password": schema.StringAttribute{
				Required:            true,
				Sensitive:           true,
				MarkdownDescription: "User Password",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(), // Force new if password changes as update might not support it
				},
			},
			"password_confirmation": schema.StringAttribute{
				Required:            true,
				Sensitive:           true,
				MarkdownDescription: "User Password Confirmation",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"role_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "User Role ID",
			},
			"search_id": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "User Search ID",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *UserResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *UserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data UserResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	userReq := client.User{
		Name:                 data.Name.ValueString(),
		Email:                data.Email.ValueString(),
		Password:             data.Password.ValueString(),
		PasswordConfirmation: data.PasswordConfirmation.ValueString(),
		RoleID:               int(data.RoleID.ValueInt64()),
	}

	if !data.SearchID.IsNull() && !data.SearchID.IsUnknown() {
		userReq.SearchID = int(data.SearchID.ValueInt64())
	}

	createdUser, err := r.client.CreateUser(userReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create user, got error: %s", err))
		return
	}

	data.ID = types.StringValue(strconv.Itoa(createdUser.ID))
	data.Name = types.StringValue(createdUser.Name)
	data.Email = types.StringValue(createdUser.Email)
	data.RoleID = types.Int64Value(int64(createdUser.RoleID))

	if createdUser.SearchID != 0 {
		data.SearchID = types.Int64Value(int64(createdUser.SearchID))
	} else {
		data.SearchID = types.Int64Null()
	}

	// Password fields are not returned by API, keep them from plan
	// data.Password and data.PasswordConfirmation are already set from plan

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *UserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data UserResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	createdUser, err := r.client.GetUser(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read user, got error: %s", err))
		return
	}

	data.Name = types.StringValue(createdUser.Name)
	data.Email = types.StringValue(createdUser.Email)
	data.RoleID = types.Int64Value(int64(createdUser.RoleID))

	if createdUser.SearchID != 0 {
		data.SearchID = types.Int64Value(int64(createdUser.SearchID))
	} else {
		data.SearchID = types.Int64Null()
	}

	// Passwords are not returned, so we don't update them in state from read

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *UserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data UserResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	userReq := client.User{
		Name:   data.Name.ValueString(),
		Email:  data.Email.ValueString(),
		RoleID: int(data.RoleID.ValueInt64()),
	}

	if !data.SearchID.IsNull() && !data.SearchID.IsUnknown() {
		userReq.SearchID = int(data.SearchID.ValueInt64())
	}

	updatedUser, err := r.client.UpdateUser(data.ID.ValueString(), userReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update user, got error: %s", err))
		return
	}

	data.Name = types.StringValue(updatedUser.Name)
	data.Email = types.StringValue(updatedUser.Email)
	data.RoleID = types.Int64Value(int64(updatedUser.RoleID))

	if updatedUser.SearchID != 0 {
		data.SearchID = types.Int64Value(int64(updatedUser.SearchID))
	} else {
		data.SearchID = types.Int64Null()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *UserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data UserResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteUser(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete user, got error: %s", err))
		return
	}
}

func (r *UserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
