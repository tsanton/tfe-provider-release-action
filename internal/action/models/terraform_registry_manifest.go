package models

/*
	{
		"version": 1,
		"metadata": {
			"protocol_versions": [
			"6.0"
		]
	}
*/

type RegistryManifest struct {
	Version  int                  `json:"version"`
	Metadata RegistryProtocolData `json:"metadata"`
}

type RegistryProtocolData struct {
	ProtocolVersions []string `json:"protocol_versions"`
}
