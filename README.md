# Terraform Provider Snyk

Terraform Provider Snyk allows Terraform to manage [Snyk](https://snyk.io) resources.

_This template repository is built on the [Terraform Plugin Framework](https://github.com/hashicorp/terraform-plugin-framework). The template repository built on the [Terraform Plugin SDK](https://github.com/hashicorp/terraform-plugin-sdk) can be found at [terraform-provider-scaffolding](https://github.com/hashicorp/terraform-provider-scaffolding). See [Which SDK Should I Use?](https://www.terraform.io/docs/plugin/which-sdk.html) in the Terraform documentation for additional information._

## Using the provider

See examples folder for instructions on how to configured the provider and resources.

## Development
### Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.18

### Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install
```

### Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

### Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```shell
make testacc
```

### Testing Locally
Make sure Terraform is configured to point out to the local installation of the provider by modifying ```~/.terraformrc```, adjust source code location accordingly. This configuration is based on [this tutorial](https://developer.hashicorp.com/terraform/tutorials/providers-plugin-framework/providers-plugin-framework-provider).

```bash
➜  snyk-terraform-provider git:(main) ✗ cat ~/.terraformrc
plugin_cache_dir   = "$HOME/.terraform.d/plugin-cache"
provider_installation {

  dev_overrides {
     "registry.terraform.io/snyk-terraform-assets/snyk" = "/Users/muratcelep/git/terraform-provider-snyk"
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
```


## Releases

Releasing a new version of the policy engine is highly automated.  A new
version can be released by running the following command, from the root of
the repository:

```bash
VERSION=v1.2.3 make release
```

In order to determine the `VERSION`, we use [semantic versioning].

The make target will print out a link to create a PR from a release branch.

Once this PR is merged to `main`, the release is created automatically.

### How this works

-  We use [changie] to add changes entries on each PR.  These are batched
   together when a version is released and added to CHANGELOG.md.

-  When we open a new PR, we kick off the
   [rc.yml workflow](../.github/workflows/rc.yml) that tests the release build.

-  When we merge a `release/*` PR, the
   [release_workflow.yml](../.github/workflows/release_workflow.yml)
   tags the release, and runs [goreleaser] to build the executables and
   upload them to the releases page on GitHub.

There is also a
[release_manual.yml](../.github/workflows/release_manual.yml) workflow that
can be triggered by manually pushing a tag.

[changie]: https://changie.dev/
[semantic versioning]: https://semver.org/
[goreleaser]: https://goreleaser.com/
