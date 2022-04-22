package toml

import (
	"github.com/jacksonCLyu/ridi-faces/pkg/configer"
	toml2 "github.com/pelletier/go-toml/v2"
)

// Decode to decode the given toml bytes to config map
func Decode(bytes []byte) (map[string]configer.Field, error) {
	decodeMap := make(map[string]any)
	err := toml2.Unmarshal(bytes, &decodeMap)
	if err != nil {
		return nil, err
	}
	configMap := make(map[string]configer.Field, len(decodeMap))
	for key, value := range decodeMap {
		configMap[key] = configer.Atof(value)
	}
	return configMap, nil
}

// Encode to encode the given config map to toml bytes
func Encode(configMap map[string]configer.Field) ([]byte, error) {
	return toml2.Marshal(configMap)
}
