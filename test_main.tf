terraform {
  required_providers {
    sodium = {
      source  = "ijcd/sodium"
      version = ">= 0.0.1"
    }
  }
}

resource "sodium_encrypted_item" "encrypted_key" {
  public_key_base64 = "qbQa1k8xeoBnKnzenV/QGTsNiCaGdDS0fjBpVuz1RFI="
  content_base64    = base64encode("the_secret")
}
