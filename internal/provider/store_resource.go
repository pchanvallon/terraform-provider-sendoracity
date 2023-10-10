package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/pchanvallon/terraform-provider-sendoracity/internal/client"
)

var _ resource.Resource = &StoreResource{}
var _ resource.ResourceWithImportState = &StoreResource{}

func NewStoreResource() resource.Resource {
	return &StoreResource{}
}

type StoreResource struct {
	client *client.SendoraCityClient
	url    string
}

type StoreResourceModel struct {
	Id      types.String `tfsdk:"id"`
	CityId  types.String `tfsdk:"city_id"`
	Address types.String `tfsdk:"address"`
	Name    types.String `tfsdk:"name"`
	Type    types.String `tfsdk:"type"`
}

func (r *StoreResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_store"
}

func (r *StoreResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Store resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Store identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"city_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Store city identifier",
			},
			"address": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Store address",
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Store name",
			},
			"type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Store type",
				Validators: []validator.String{
					stringvalidator.OneOf("Food", "Sports", "Clothes", "Electronics", "Other"),
				},
			},
		},
	}
}

func (r *StoreResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.SendoraCityClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.SendoraCityClient, got: %T."+
				"Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
	r.url = "stores"
}

func (r *StoreResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *StoreResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	cityId, err := strconv.Atoi(data.CityId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to convert city_id to int, got error: %s", err))
		return
	}

	body := &Store{
		CityId:  cityId,
		Address: data.Address.ValueString(),
		Name:    data.Name.ValueString(),
		Type:    data.Type.ValueString(),
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		resp.Diagnostics.AddError("JSON parser Error",
			fmt.Sprintf("Unable to marshal JSON body for host, got error: %s", err))
		return
	}

	res, err := r.client.DoCreate(r.url, jsonBody)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to create store, got error: %s", err))
		return
	}

	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read response body, got error: %s", err))
		return
	}
	defer res.Body.Close()

	store := &Store{}
	if err = json.Unmarshal(responseBody, &store); err != nil {
		resp.Diagnostics.AddError("JSON parser Error",
			fmt.Sprintf("Unable to unmarshal JSON body for store, got error: %s", err))
		return
	}

	data.Id = types.StringValue(strconv.Itoa(store.Id))

	tflog.Trace(ctx, "created a store resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *StoreResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *StoreResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	res, err := r.client.DoRead(fmt.Sprintf("%s/%s", r.url, data.Id.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read store, got error: %s", err))
		return
	}

	if res.StatusCode == http.StatusNotFound {
		resp.State.RemoveResource(ctx)
		return
	}

	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read response body, got error: %s", err))
		return
	}
	defer res.Body.Close()

	store := &Store{}
	if err = json.Unmarshal(responseBody, &store); err != nil {
		resp.Diagnostics.AddError("JSON parser Error",
			fmt.Sprintf("Unable to unmarshal JSON body for store, got error: %s", err))
		return
	}

	data.CityId = types.StringValue(strconv.Itoa(store.CityId))
	data.Address = types.StringValue(store.Address)
	data.Name = types.StringValue(store.Name)
	data.Type = types.StringValue(store.Type)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *StoreResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *StoreResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	cityId, err := strconv.Atoi(data.CityId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to convert city_id to int, got error: %s", err))
		return
	}

	body := &Store{
		CityId:  cityId,
		Address: data.Address.ValueString(),
		Name:    data.Name.ValueString(),
		Type:    data.Type.ValueString(),
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		resp.Diagnostics.AddError("JSON parser Error",
			fmt.Sprintf("Unable to marshal JSON body for host, got error: %s", err))
		return
	}

	id := data.Id.ValueString()
	if _, err = r.client.DoUpdate(fmt.Sprintf("%s/%s", r.url, id), jsonBody); err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to update store with id %s, got error: %s", id, err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *StoreResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *StoreResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	id := data.Id.ValueString()
	if _, err := r.client.DoDelete(fmt.Sprintf("%s/%s", r.url, id)); err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to delete store with id %s, got error: %s", id, err))
		return
	}
}

func (r *StoreResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
