package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pchanvallon/terraform-provider-sendoracity/internal/client"
)

var _ provider.Provider = &SendoraCityProvider{}

type SendoraCityProvider struct {
	version string
}

type SendoraCityProviderModel struct {
	BaseUri types.String `tfsdk:"base_uri"`
}

func (p *SendoraCityProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "sendoracity"
	resp.Version = p.version
}

func (p *SendoraCityProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"base_uri": schema.StringAttribute{
				MarkdownDescription: "City API base URI",
				Optional:            true,
			},
		},
	}
}

func (p *SendoraCityProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data SendoraCityProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.BaseUri.IsNull() {
		data.BaseUri = types.StringValue(os.Getenv("BASE_URI"))
	}

	if data.BaseUri.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Missing URL Configuration",
			"While configuring the provider, the patrol URL was not found in "+
				"the BASE_URI environment variable or provider configuration block url attribute.",
		)
	}

	client := client.NewClient(data.BaseUri.ValueString())
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *SendoraCityProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewCityResource,
		NewHouseResource,
		NewStoreResource,
	}
}

func (p *SendoraCityProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewCityDataSource,
		NewHouseDataSource,
		NewStoreDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &SendoraCityProvider{
			version: version,
		}
	}
}
