package model

// Endpoint store data for endpoint.
type Endpoint struct {
	Type    string                 `yaml:"type" msgpack:"t"`
	Path    string                 `yaml:"path" msgpack:"p"`
	Options map[string]interface{} `yaml:"options" msgpack:"o"`
}
