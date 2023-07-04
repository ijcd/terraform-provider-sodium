// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golang.org/x/crypto/nacl/box"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &EncryptedItemDataSource{}

func NewEncryptedItemDataSource() datasource.DataSource {
	return &EncryptedItemDataSource{}
}

// EncryptedItemDataSource defines the data source implementation.
type EncryptedItemDataSource struct {
	// client *http.Client
}

// EncryptedItemDataSourceModel describes the data source data model.
type EncryptedItemDataSourceModel struct {
	PublicKeyBase64      types.String `tfsdk:"public_key_base64"`
	ContentBase64        types.String `tfsdk:"content_base64"`
	EncryptedValueBase64 types.String `tfsdk:"encrypted_value_base64"`
	ContentChecksum      types.String `tfsdk:"content_checksum"`
	Id                   types.String `tfsdk:"id"`
}

// EncryptedItemDataSourceModel describes the data source data model.
type EncryptedItemDataSourceStateModel struct {
}

func (d *EncryptedItemDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_encrypted_item"
	tflog.Info(ctx, "EncryptedItem data source metadata", map[string]interface{}{
		"typename": resp.TypeName,
	})
}

func (d *EncryptedItemDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "EncryptedItem data source",

		Attributes: map[string]schema.Attribute{
			"public_key_base64": schema.StringAttribute{
				MarkdownDescription: "Public key to use when encrypting base64 encoded",
				Required:            true,
			},
			"content_base64": schema.StringAttribute{
				MarkdownDescription: "Base64 encoded version of the raw string to encrypt.",
				Sensitive:           true,
				Required:            true,
			},
			"encrypted_value_base64": schema.StringAttribute{
				MarkdownDescription: "Base64 encoded version of the encrypted result .",
				Sensitive:           true,
				Computed:            true,
			},
			"content_checksum": schema.StringAttribute{
				MarkdownDescription: "Base64 encoded version of the encrypted result .",
				Computed:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "EncryptedItem identifier",
				Computed:            true,
			},
		},
	}
}

func (d *EncryptedItemDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	// client, ok := req.ProviderData.(*http.Client)

	// if !ok {
	// 	resp.Diagnostics.AddError(
	// 		"Unexpected Data Source Configure Type",
	// 		fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
	// 	)

	// 	return
	// }

	// d.client = client
}

func (d *EncryptedItemDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data EncryptedItemDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	secret, err := base64.StdEncoding.DecodeString(data.ContentBase64.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Decode Error", fmt.Sprintf("Unable to decode content_base64, got error: %s", err))
	}
	secretBytes := []byte(secret)

	// Getting public key from input
	var pkBytes [32]byte
	key, err := base64.StdEncoding.DecodeString(data.PublicKeyBase64.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Decode Error", fmt.Sprintf("Failed to decode base64 public_key_base64, got error: %s", err))
	}
	copy(pkBytes[:], key)

	// Encrypting string with given pubKey
	enc, err := box.SealAnonymous(nil, secretBytes, &pkBytes, nil)
	if err != nil {
		resp.Diagnostics.AddError("Encrypt Error", fmt.Sprintf("Failed to encrypt, got error: %s", err))
	}

	// Encoding result to base64
	encEnc := base64.StdEncoding.EncodeToString(enc)
	data.EncryptedValueBase64 = types.StringValue(encEnc)

	// The ID is an encoding of the output value
	outChecksum := sha1.Sum([]byte(encEnc))
	data.Id = types.StringValue(hex.EncodeToString(outChecksum[:]))

	// The checksum is an encoding of the input value
	inChecksum := sha1.Sum([]byte(data.ContentBase64.ValueString()))
	data.ContentChecksum = types.StringValue(hex.EncodeToString(inChecksum[:]))

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
