terraform {
  required_providers {
    sodium = {
      source  = "ijcd/sodium"
      version = ">= 0.0.1"
    }
  }
}

data "sodium_encrypted_item" "encrypted_key" {
  public_key_base64 = "qbQa1k8xeoBnKnzenV/QGTsNiCaGdDS0fjBpVuz1RFI="
  content_base64    = base64encode("the_secret")
}

output "encrypted_value" {
  value     = data.sodium_encrypted_item.encrypted_key.encrypted_value_base64
  sensitive = true
}
