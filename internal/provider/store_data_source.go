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

var _ datasource.DataSource = &StoreDataSource{}

func NewStoreDataSource() datasource.DataSource {
	return &StoreDataSource{}
}

type StoreDataSource struct {
	client *client.SendoraCityClient
	url    string
}

type StoreDataSourceModel struct {
	Id      types.String `tfsdk:"id"`
	CityId  types.String `tfsdk:"city_id"`
	Address types.String `tfsdk:"address"`
	Name    types.String `tfsdk:"name"`
	Type    types.String `tfsdk:"type"`
}

func (d *StoreDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_store"
}

func (d *StoreDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Store data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Store identifier",
				Optional:            true,
			},
			"city_id": schema.StringAttribute{
				MarkdownDescription: "Store city identifier",
				Computed:            true,
			},
			"address": schema.StringAttribute{
				MarkdownDescription: "Store address",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Store name",
				Computed:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Store type",
				Computed:            true,
			},
		},
	}
}

func (d *StoreDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
	d.url = "stores"
}

func (d *StoreDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data StoreDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	res, err := d.client.DoRead(fmt.Sprintf("%s/%s", d.url, data.Id.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read store, got error: %s", err))
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

	tflog.Trace(ctx, "read store from a data source")
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
