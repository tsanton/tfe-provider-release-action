package action_test

import (
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	r "github.com/tsanton/tfe-provider-release-action/action"
	m "github.com/tsanton/tfe-provider-release-action/action/models"
)

func Test_parse_goreleaser_artifact(t *testing.T) {
	config := r.RunConfig{}
	err := config.ParseGoreleaseArtifacts(logger, getGoreleaserArtifactString())
	assert.Nil(t, err)
	/* ProviderVersionPlatforms */
	assert.Equal(t, 2, len(config.ProviderVersionPlatforms))

	linux, ok := lo.Find(config.ProviderVersionPlatforms, func(p m.GoreleaserArtifactArchive) bool { return p.Os == "linux" && p.Goarch == "amd64" })
	assert.True(t, ok)
	assert.Equal(t, "terraform-provider-tfepatch_0.1.8_linux_amd64.zip", linux.Name)
	assert.Equal(t, "dist/terraform-provider-tfepatch_0.1.8_linux_amd64.zip", linux.Path)
	assert.Equal(t, "sha256:f79a37934c2e9d793ee874eeef921213e72069480c9df3953c907d79a0f5f034", linux.ShaSum)

	darwin, ok := lo.Find(config.ProviderVersionPlatforms, func(p m.GoreleaserArtifactArchive) bool { return p.Os == "darwin" && p.Goarch == "amd64" })
	assert.True(t, ok)
	assert.Equal(t, "terraform-provider-tfepatch_0.1.8_darwin_amd64.zip", darwin.Name)
	assert.Equal(t, "dist/terraform-provider-tfepatch_0.1.8_darwin_amd64.zip", darwin.Path)
	assert.Equal(t, "sha256:f2991fc425fbaa4033f10cf7962243832056762c9fbecb656db11e429ea9551e", darwin.ShaSum)

	/* Checksum */
	assert.NotEqual(t, m.GoreleaserArtifactFile{}, config.Checksum)
	assert.Equal(t, "terraform-provider-tfepatch_0.1.8_SHA256SUMS", config.Checksum.Name)
	assert.Equal(t, "Checksum", config.Checksum.Type)

	/* Signature */
	assert.NotEqual(t, m.GoreleaserArtifactFile{}, config.Signature)
	assert.Equal(t, "terraform-provider-tfepatch_0.1.8_SHA256SUMS.sig", config.Signature.Name)
	assert.Equal(t, "Signature", config.Signature.Type)
}

func Test_parse_gorelease_manifest(t *testing.T) {
	config := r.RunConfig{}
	err := config.ParseGoreleaserMetadata(logger, getGoreleaserMetadataString())
	assert.Nil(t, err)
	assert.Equal(t, "terraform-provider-tfepatch", config.GoreleaserMetadata.ProviderName)
	assert.Equal(t, "0.1.8", config.GoreleaserMetadata.Version)
}
