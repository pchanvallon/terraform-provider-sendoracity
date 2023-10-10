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

var _ resource.Resource = &CityResource{}
var _ resource.ResourceWithImportState = &CityResource{}

func NewCityResource() resource.Resource {
	return &CityResource{}
}

type CityResource struct {
	client *client.SendoraCityClient
	url    string
}

type CityResourceModel struct {
	Id        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	Touristic types.Bool   `tfsdk:"touristic"`
}

func (r *CityResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_city"
}

func (r *CityResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "City resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "City identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "City name",
			},
			"touristic": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Whether the city is touristic or not",
			},
		},
	}
}

func (r *CityResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
	r.url = "cities"
}

func (r *CityResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *CityResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	touristic := data.Touristic.ValueBool()
	body := &City{
		Name:      data.Name.ValueString(),
		Touristic: &touristic,
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
			fmt.Sprintf("Unable to create city, got error: %s", err))
		return
	}

	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read response body, got error: %s", err))
		return
	}
	defer res.Body.Close()

	city := &City{}
	if err = json.Unmarshal(responseBody, &city); err != nil {
		resp.Diagnostics.AddError("JSON parser Error",
			fmt.Sprintf("Unable to unmarshal JSON body for city, got error: %s", err))
		return
	}

	data.Id = types.StringValue(strconv.Itoa(city.Id))

	tflog.Trace(ctx, "created a city resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CityResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *CityResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	res, err := r.client.DoRead(fmt.Sprintf("%s/%s", r.url, data.Id.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read city, got error: %s", err))
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

	city := &City{}
	if err = json.Unmarshal(responseBody, &city); err != nil {
		resp.Diagnostics.AddError("JSON parser Error",
			fmt.Sprintf("Unable to unmarshal JSON body for city, got error: %s", err))
		return
	}

	data.Name = types.StringValue(city.Name)
	data.Touristic = types.BoolValue(*city.Touristic)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CityResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *CityResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	touristic := data.Touristic.ValueBool()
	body := &City{
		Name:      data.Name.ValueString(),
		Touristic: &touristic,
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
			fmt.Sprintf("Unable to update city with id %s, got error: %s", id, err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CityResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *CityResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	id := data.Id.ValueString()
	if _, err := r.client.DoDelete(fmt.Sprintf("%s/%s", r.url, id)); err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to delete city with id %s, got error: %s", id, err))
		return
	}
}

func (r *CityResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
