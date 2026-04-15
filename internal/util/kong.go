package util

import (
	"encoding/json"
	"io"
	"strings"

	"github.com/alecthomas/kong"
	"gopkg.in/yaml.v3"
)

func YamlConfigLoader(r io.Reader) (kong.Resolver, error) {
	values := map[string]any{}
	if err := yaml.NewDecoder(r).Decode(&values); err != nil {
		return nil, err
	}
	normalizeKeys(values)
	payload, err := json.Marshal(values)
	if err != nil {
		return nil, err
	}
	return kong.JSON(BytesReader(payload))
}

// normalizeKeys replaces hyphens with underscores in map keys so that
// kong.JSON can match them via its snake_case variant (e.g. "region-map" -> "region_map").
func normalizeKeys(m map[string]any) {
	for key, val := range m {
		newKey := strings.ReplaceAll(key, "-", "_")
		if newKey != key {
			delete(m, key)
			m[newKey] = val
			key = newKey
		}
		if nested, ok := m[key].(map[string]any); ok {
			normalizeKeys(nested)
		}
	}
}
