package yml

import (
	"github.com/jacksonCLyu/ridi-faces/pkg/configer"
	"gopkg.in/yaml.v2"
)

// Decode decodes the given YAML document bytes to the config map.
func Decode(b []byte) (map[string]configer.Field, error) {
	decodeMap := make(map[string]any)
	if err := yaml.Unmarshal(b, &decodeMap); err != nil {
		return nil, err
	}
	configMap := make(map[string]configer.Field, len(decodeMap))
	for k, v := range decodeMap {
		configMap[k] = configer.Atof(v)
	}
	return configMap, nil
}

// Encode encodes the given config map to YAML document bytes.
func Encode(m map[string]configer.Field) ([]byte, error) {
	return yaml.Marshal(m)
}
