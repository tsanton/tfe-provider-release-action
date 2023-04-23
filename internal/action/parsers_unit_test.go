package action_test

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	r "github.com/tsanton/tfe-provider-release-action/action"
)

func Test_read_file_to_byte_array(t *testing.T) {
	fullPath := path.Join("/workspace/testing-assets/", "/terraform-provider-tfepatch_0.1.8_manifest.json")
	content, err := r.ReadFileToByteArray(fullPath)
	assert.Nil(t, err)
	assert.Greater(t, len(content), 0)
}
