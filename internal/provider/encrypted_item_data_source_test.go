// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSodiumEncryptedItemDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccSodiumEncryptedItemDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.sodium_encrypted_item.test", "id", "encrypted_item-id"),
				),
			},
		},
	})
}

const testAccSodiumEncryptedItemDataSourceConfig = `
data "sodium_encrypted_item" "encrypted_key" {
	public_key_base64 = "qbQa1k8xeoBnKnzenV/QGTsNiCaGdDS0fjBpVuz1RFI="
	content_base64    = base64encode("the_secret")
}
`
