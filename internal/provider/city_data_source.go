package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/pchanvallon/terraform-provider-sendoracity/internal/client"
)

var _ datasource.DataSource = &CityDataSource{}

func NewCityDataSource() datasource.DataSource {
	return &CityDataSource{}
}

type CityDataSource struct {
	client *client.SendoraCityClient
	url    string
}

type CityDataSourceModel struct {
	Id        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	Touristic types.Bool   `tfsdk:"touristic"`
}

func (d *CityDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_city"
}

func (d *CityDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "City data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "City identifier",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "City name",
				Optional:            true,
				Computed:            true,
			},
			"touristic": schema.BoolAttribute{
				MarkdownDescription: "Whether the city is touristic or not",
				Computed:            true,
			},
		},
	}
}

func (d *CityDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = client
	d.url = "cities"
}

func (d *CityDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CityDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var city City
	if !data.Id.IsNull() {
		res, err := d.client.DoRead(fmt.Sprintf("%s/%s", d.url, data.Id.ValueString()))
		if err != nil {
			resp.Diagnostics.AddError("Client Error",
				fmt.Sprintf("Unable to read city, got error: %s", err))
			return
		}

		responseBody, err := io.ReadAll(res.Body)
		if err != nil {
			resp.Diagnostics.AddError("Client Error",
				fmt.Sprintf("Unable to read response body, got error: %s", err))
			return
		}
		defer res.Body.Close()

		if err = json.Unmarshal(responseBody, &city); err != nil {
			resp.Diagnostics.AddError("JSON parser Error",
				fmt.Sprintf("Unable to unmarshal JSON body for city, got error: %s", err))
			return
		}
	} else if !data.Name.IsNull() {
		res, err := d.client.DoList(d.url, map[string]string{"name": data.Name.ValueString()})
		if err != nil {
			resp.Diagnostics.AddError("Client Error",
				fmt.Sprintf("Unable to read city, got error: %s", err))
			return
		}

		responseBody, err := io.ReadAll(res.Body)
		if err != nil {
			resp.Diagnostics.AddError("Client Error",
				fmt.Sprintf("Unable to read response body, got error: %s", err))
			return
		}
		defer res.Body.Close()

		cities := []City{}
		if err = json.Unmarshal(responseBody, &cities); err != nil {
			resp.Diagnostics.AddError("JSON parser Error",
				fmt.Sprintf("Unable to unmarshal JSON body for city, got error: %s", err))
			return
		}
		city = cities[0]
	} else {
		resp.Diagnostics.AddError("Invalid configuration",
			"Either id or name must be set")
		return
	}

	data.Id = types.StringValue(strconv.Itoa(city.Id))
	data.Name = types.StringValue(city.Name)
	data.Touristic = types.BoolValue(*city.Touristic)

	tflog.Trace(ctx, "read city from a data source")
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
