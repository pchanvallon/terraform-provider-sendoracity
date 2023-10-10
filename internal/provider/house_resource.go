package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/pchanvallon/terraform-provider-sendoracity/internal/client"
)

var _ resource.Resource = &HouseResource{}
var _ resource.ResourceWithImportState = &HouseResource{}

func NewHouseResource() resource.Resource {
	return &HouseResource{}
}

type HouseResource struct {
	client *client.SendoraCityClient
	url    string
}

type HouseResourceModel struct {
	Id          types.String `tfsdk:"id"`
	CityId      types.String `tfsdk:"city_id"`
	Address     types.String `tfsdk:"address"`
	Inhabitants types.Int64  `tfsdk:"inhabitants"`
}

func (r *HouseResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_house"
}

func (r *HouseResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "House resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "House identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"city_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "House city identifier",
			},
			"address": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "House address",
			},
			"inhabitants": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "House inhabitants count",
			},
		},
	}
}

func (r *HouseResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
	r.url = "houses"
}

func (r *HouseResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *HouseResourceModel

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

	body := &House{
		CityId:      cityId,
		Address:     data.Address.ValueString(),
		Inhabitants: int(data.Inhabitants.ValueInt64()),
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
			fmt.Sprintf("Unable to create house, got error: %s", err))
		return
	}

	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read response body, got error: %s", err))
		return
	}
	defer res.Body.Close()

	house := &House{}
	if err = json.Unmarshal(responseBody, &house); err != nil {
		resp.Diagnostics.AddError("JSON parser Error",
			fmt.Sprintf("Unable to unmarshal JSON body for house, got error: %s", err))
		return
	}

	data.Id = types.StringValue(strconv.Itoa(house.Id))

	tflog.Trace(ctx, "created a house resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *HouseResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *HouseResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	res, err := r.client.DoRead(fmt.Sprintf("%s/%s", r.url, data.Id.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read house, got error: %s", err))
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

	house := &House{}
	if err = json.Unmarshal(responseBody, &house); err != nil {
		resp.Diagnostics.AddError("JSON parser Error",
			fmt.Sprintf("Unable to unmarshal JSON body for house, got error: %s", err))
		return
	}

	data.CityId = types.StringValue(strconv.Itoa(house.CityId))
	data.Address = types.StringValue(house.Address)
	data.Inhabitants = types.Int64Value(int64(house.Inhabitants))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *HouseResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *HouseResourceModel

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

	body := &House{
		CityId:      cityId,
		Address:     data.Address.ValueString(),
		Inhabitants: int(data.Inhabitants.ValueInt64()),
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
			fmt.Sprintf("Unable to update house with id %s, got error: %s", id, err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *HouseResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *HouseResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	id := data.Id.ValueString()
	if _, err := r.client.DoDelete(fmt.Sprintf("%s/%s", r.url, id)); err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to delete house with id %s, got error: %s", id, err))
		return
	}
}

func (r *HouseResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
