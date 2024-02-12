<!-- markdownlint-disable first-line-h1 no-inline-html -->
<a href="https://terraform.io">
    <img src="https://github.com/hashicorp/terraform-provider-aws/raw/main/.github/terraform_logo_light.svg" alt="Terraform logo" title="Terraform" align="right" height="50" />
</a>

# Terraform Chronicle

[![CI](https://github.com/form3tech-oss/terraform-provider-chronicle/actions/workflows/ci.yaml/badge.svg)](https://github.com/form3tech-oss/terraform-provider-chronicle/actions/workflows/ci.yaml)
[![release](https://github.com/form3tech-oss/terraform-provider-chronicle/actions/workflows/release.yaml/badge.svg)](https://github.com/form3tech-oss/terraform-provider-chronicle/actions/workflows/release.yaml)

Terraform provider for Chronicle

# Documentation

Find available docs [here](docs/).

## Building The Provider

```bash
make build
```

## Debugging

TF examples can be found under `examples/` directory, and can be used for debugging purposes.

- run `./debug.sh` script in this repo root and init target workspace (`terraform init`)
- attach with your debugger to port `2345`, this will also print required gRPC debug info to stdout
- note the env variable from the above output (you'll need it in the next step)
- set breakpoints in code and run a plan (`TF_REATTACH_PROVIDERS='{...}' terraform plan`)
- re-run the above command as many times as needed (provider process keeps running after the plan has finished)

## Release

After merging to master, GitHub Action will push new tag, which in turn will trigger `goreleaser`.

`goreleaser` action will build and upload to GitHub release.

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.21+ is *required*).

To compile the provider, run `make build`.

```sh
make build
```

In order to test the provider, you can simply run `make test`.

```sh
make test
```

In order to run the full suite of Acceptance tests, set the environment variables listed below and run `make testacc`.

The order of precedence for chronicle's API configuration is the following: `Credential file through TF > Access Token through TF > Environment Variable`.
Environment variables always take the lowest precedence

*Note:* Acceptance tests create real resources
| Environment variables            | Description                                |
|:---------------------------------|:-------------------------------------------|
| CHRONICLE_BACKSTORY_CREDENTIALS  | backstory base64 credentials               |
| CHRONICLE_BIGQUERY_CREDENTIALS   | bigquery base64 credentials                |
| CHRONICLE_INGESTION_CREDENTIALS  | ingestion base64 credentials               |
| CHRONICLE_FORWARDER_CREDENTIALS  | forwarder base64 credentials               |
| CHRONICLE_REGION                 | API region                                 |

## Using a local version of the provider
Firstly install the provider by running:

```sh
make install
```
Modify your `.terraformrc` file as the following:

```tf
provider_installation {
  filesystem_mirror {
    path = "<HOME_PATH>/.terraform.d/plugins"
    include = ["github.com/form3tech-oss/chronicle"]
  }
  direct {
    exclude = ["github.com/form3tech-oss/chronicle"]
  }
}
```

Now leverage the provider by using:
```tf
terraform {
  required_providers {
    chronicle = {
      source  = "github.com/form3tech-oss/chronicle"
    }
  }
}
```