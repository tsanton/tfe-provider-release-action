package models

/*
	### goreleaser artifact: ###
	[
		{
			"name":"terraform-provider-tfepatch_0.1.8_SHA256SUMS",
			"path":"dist/terraform-provider-tfepatch_0.1.8_SHA256SUMS",
			"internal_type":12,
			"type":"Checksum",
			"extra":{ }
   		},
		{
			"name":"terraform-provider-tfepatch_0.1.8_SHA256SUMS.sig",
			"path":"dist/terraform-provider-tfepatch_0.1.8_SHA256SUMS.sig",
			"internal_type":13,
			"type":"Signature",
			"extra":{
				"ID":"default"
			}
		},
		{
			"name":"terraform-provider-tfepatch_0.1.8_manifest.json",
			"path":"terraform-registry-manifest.json",
			"internal_type":3,
			"type":"File"
		}
	]
*/

type GoreleaserArtifactFile struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Type string `json:"type"`
}
