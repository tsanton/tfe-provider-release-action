# **Terraform Cloud Provider Release Action**

This action is intended to release your private terraform providers to Terraform Cloud/Terraform Enterprise.

To utilise this action required you have created a [provider platform]() and uploaded a [GPG key]().

## **Dependencies**

This actions in made to run in lockstep with the [goreleaser-action](https://github.com/goreleaser/goreleaser-action) as it is dependent on the output that action.

In short `goreleaser` outputs the following fields:

- `metadata`: see the [getGoreleaserMetadataString()](./internal/action/main_test.go#L138) function for data format example
- `artifacts`: see the [getGoreleaserArtifactString()](./internal/action/main_test.go#L154) function for data format example

These outputs from goreleaser informs this action about the required metadata for the

## **Limitations**

You can configure the build of multiple packages in goreleaser

This project does not support the release of multiple builds as it grabs the first checksum and signature objects in the `artifact` output from the goreleaser action and upload that to your [provider version](https://developer.hashicorp.com/terraform/cloud-docs/registry/publish-providers#create-a-version-and-platform)
If you try, it might work, but for two releases you'll end up with all artifacts for one release having the wrong checksums submitted

## **Usage**

As mentioned above, this action depends on you having created a provider platform and uploaded a GPG key.
As of the time of writing, it's not possible to do any of this in Terraform Cloud. Therefor I've created a demo project to release your first private provider.
Ironically it's a terraform provider to patch terraform cloud, i.e. create and manage the resources required by this actions.

All you have to do is to follow the steps in [this repository]()!

## **Sauce**

The action is made by following the steps in this [official terraform cloud doc](https://developer.hashicorp.com/terraform/cloud-docs/registry/publish-providers) for how publish private providers.
