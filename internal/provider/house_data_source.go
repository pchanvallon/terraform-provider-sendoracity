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

var _ datasource.DataSource = &HouseDataSource{}

func NewHouseDataSource() datasource.DataSource {
	return &HouseDataSource{}
}

type HouseDataSource struct {
	client *client.SendoraCityClient
	url    string
}

type HouseDataSourceModel struct {
	Id          types.String `tfsdk:"id"`
	CityId      types.String `tfsdk:"city_id"`
	Address     types.String `tfsdk:"address"`
	Inhabitants types.Int64  `tfsdk:"inhabitants"`
}

func (d *HouseDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_house"
}

func (d *HouseDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "House data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "House identifier",
				Required:            true,
			},
			"city_id": schema.StringAttribute{
				MarkdownDescription: "House city identifier",
				Computed:            true,
			},
			"address": schema.StringAttribute{
				MarkdownDescription: "House address",
				Computed:            true,
			},
			"inhabitants": schema.Int64Attribute{
				MarkdownDescription: "House inhabitants count",
				Computed:            true,
			},
		},
	}
}

func (d *HouseDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
	d.url = "houses"
}

func (d *HouseDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data HouseDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	res, err := d.client.DoRead(fmt.Sprintf("%s/%s", d.url, data.Id.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read house, got error: %s", err))
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

	tflog.Trace(ctx, "read house from a data source")
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
