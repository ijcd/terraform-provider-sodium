// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golang.org/x/crypto/nacl/box"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &EncryptedItemResource{}
var _ resource.ResourceWithImportState = &EncryptedItemResource{}

func NewEncryptedItemResource() resource.Resource {
	return &EncryptedItemResource{}
}

// EncryptedItemResource defines the resource implementation.
type EncryptedItemResource struct {
	// client *http.Client
}

// EncryptedItemResourceModel describes the resource data model.
type EncryptedItemResourceModel struct {
	PublicKeyBase64      types.String `tfsdk:"public_key_base64"`
	ContentBase64        types.String `tfsdk:"content_base64"`
	EncryptedValueBase64 types.String `tfsdk:"encrypted_value_base64"`
	Id                   types.String `tfsdk:"id"`
}

func (r *EncryptedItemResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_encrypted_item"
}

func (r *EncryptedItemResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "EncryptedItem resource",

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
			"id": schema.StringAttribute{
				MarkdownDescription: "EncryptedItem identifier",
				Computed:            true,
			},
		},
	}
}

func (r *EncryptedItemResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	// client, ok := req.ProviderData.(*http.Client)

	// if !ok {
	// 	resp.Diagnostics.AddError(
	// 		"Unexpected Resource Configure Type",
	// 		fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
	// 	)

	// 	return
	// }

	// r.client = client
}

func (r *EncryptedItemResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *EncryptedItemResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// compute/update the state in place
	computeState(data, &resp.Diagnostics)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func computeState(data *EncryptedItemResourceModel, diag *diag.Diagnostics) {
	secret, err := base64.StdEncoding.DecodeString(data.ContentBase64.ValueString())
	if err != nil {
		diag.AddError("Decode Error", fmt.Sprintf("Unable to decode content_base64, got error: %s", err))
	}
	secretBytes := []byte(secret)

	// Getting public key from input
	var pkBytes [32]byte
	key, err := base64.StdEncoding.DecodeString(data.PublicKeyBase64.ValueString())
	if err != nil {
		diag.AddError("Decode Error", fmt.Sprintf("Failed to decode base64 public_key_base64, got error: %s", err))
	}
	copy(pkBytes[:], key)

	// Encrypting string with given pubKey
	enc, err := box.SealAnonymous(nil, secretBytes, &pkBytes, nil)
	if err != nil {
		diag.AddError("Encrypt Error", fmt.Sprintf("Failed to encrypt, got error: %s", err))
	}

	// Encoding result to base64
	encEnc := base64.StdEncoding.EncodeToString(enc)
	data.EncryptedValueBase64 = types.StringValue(encEnc)

	// The ID is a hash of the input value
	inChecksum := sha1.Sum([]byte(data.ContentBase64.ValueString()))
	data.Id = types.StringValue(hex.EncodeToString(inChecksum[:]))
}

func (r *EncryptedItemResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *EncryptedItemResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *EncryptedItemResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *EncryptedItemResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// compute/update the state in place
	computeState(data, &resp.Diagnostics)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *EncryptedItemResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *EncryptedItemResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete example, got error: %s", err))
	//     return
	// }
}

func (r *EncryptedItemResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
