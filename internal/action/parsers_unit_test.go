package action_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	r "github.com/tsanton/tfe-provider-release-action/action"
)

func Test_read_file_to_byte_array(t *testing.T) {
	content, err := r.ReadFileToByteArray("/workspace/testing-assets/terraform-provider-tfepatch_0.1.8_manifest.json")
	assert.Nil(t, err)
	assert.Greater(t, len(content), 0)
}
