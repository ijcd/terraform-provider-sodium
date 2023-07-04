// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure SodiumProvider satisfies various provider interfaces.
var _ provider.Provider = &SodiumProvider{}

// SodiumProvider defines the provider implementation.
type SodiumProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// SodiumProviderModel describes the provider data model.
type SodiumProviderModel struct {
	// Endpoint types.String `tfsdk:"endpoint"`
}

func (p *SodiumProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "sodium"
	resp.Version = p.version
}

func (p *SodiumProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			// "endpoint": schema.StringAttribute{
			// 	MarkdownDescription: "Example provider attribute",
			// 	Optional:            true,
			// },
		},
	}
}

func (p *SodiumProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data SodiumProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available.
	// if data.Endpoint.IsNull() { /* ... */ }

	// Example client configuration for data sources and resources
	client := http.DefaultClient
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *SodiumProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		// NewExampleResource,
	}
}

func (p *SodiumProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewEncryptedItemDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &SodiumProvider{
			version: version,
		}
	}
}
