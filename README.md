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

You can configure the build of multiple packages in goreleaser which this project has not been designed to support. \
The reason this might break is because this release action select the first occurrence of a `checksum` and `signature` objects in the `artifact` output from the goreleaser action. \
These files are then upload that to your [provider version](https://developer.hashicorp.com/terraform/cloud-docs/registry/publish-providers#create-a-version-and-platform). \
When you then proceed to upload your provider platform version the checksums will (probably) not be included in the `checksum` file and you might end up with signature mismatches. \

## **Usage**

As mentioned above, this action depends on you having created a provider platform and uploaded a GPG key. \
At the time of writing it's not possible to do any of this in Terraform Cloud, therefore I've created a step through guide to release your first private provider. \
Ironically it's a terraform provider to patch terraform cloud, i.e. create and manage the resources required by this actions. \

All you have to do is to follow the steps in [this repository](https://github.com/Tsanton/terraform-provider-tfepatch)!

## **Sauce**

The action is made by following the steps in this [official terraform cloud doc](https://developer.hashicorp.com/terraform/cloud-docs/registry/publish-providers) for how publish private providers.
