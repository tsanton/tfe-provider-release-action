package action_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	r "github.com/tsanton/tfe-provider-release-action/action"
)

func Test_live_run(t *testing.T) {
	/* Arrange */
	host := "https://app.terraform.io"
	orgName, tfeToken := runnerValidator(t)
	namespace := orgName
	providerName := "terraform-provider-tfepatch"
	cli := liveClientSetup(t, host, tfeToken)
	ctx := context.Background()

	presp := bootstrapProvider(t, cli, orgName, providerName)
	gresp := boostrapGpgKey(t, cli, orgName, providerName)
	defer func() {
		perr := cli.ProviderService.Delete(ctx, orgName, string(presp.Data.Attributes.RegistryName), presp.Data.Attributes.Namespace, presp.Data.Attributes.Name)
		gerr := cli.GpgService.Delete(ctx, gresp.Data.Attributes.Namespace, gresp.Data.Attributes.KeyId)
		if perr != nil {
			panic("unable to cleanup provider")
		}
		if gerr != nil {
			panic("unable to cleanup gpg key")
		}
	}()

	config := r.NewRunConfig("/workspace/testing-assets/", orgName, namespace, providerName, gresp.Data.Attributes.KeyId)
	err := config.ParseGoreleaseArtifacts(logger, getGoreleaserArtifactString())
	assert.Nil(t, err)
	err = config.ParseGoreleaserMetadata(logger, getGoreleaserMetadataString())
	assert.Nil(t, err)

	/* Act */
	err = r.Run(cli, logger, config)
	assert.Nil(t, err)

	/* Assert */
	providerVersion, _ := cli.ProviderVersionService.Read(ctx, config.Organization, config.Namespace, config.ProviderName, config.GoreleaserMetadata.Version)
	providerVersionPlatforms, _ := cli.ProviderVersionPlatformService.List(ctx, orgName, namespace, providerName, config.GoreleaserMetadata.Version)
	assert.NotNil(t, providerVersion)
	assert.NotNil(t, providerVersionPlatforms)
}
