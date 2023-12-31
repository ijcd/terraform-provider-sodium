## Using the provider

```terraform
terraform {
  required_providers {
    sodium = {
      source = "github.com/ijcd/sodium"
    }
    github = {
      source  = "integrations/github"
      version = ">= 4.5.2"
    }
  }
}

provider "github" {
  owner = "org_name"
  token = "github_token"
}


# To make sure the repository exists and the correct permissions are set.
data "github_repository" "main" {
  full_name = "org_name/repo_name"
}

data "github_actions_public_key" "gh_actions_public_key" {
  repository = "repo_name"
}

resource "sodium_encrypted_item" "foo" {
  public_key     = data.github_actions_public_key.gh_actions_public_key.key
  content_base64 = base64encode("SuperSecretPassword")
}

resource "github_actions_secret" "gh_actions_secret" {
  repository      = "repo_name"
  secret_name     = "SECRET_FOO"
  encrypted_value = sodium_encrypted_item.foo.encrypted_value_base64
}
```

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```shell
make testacc
```
