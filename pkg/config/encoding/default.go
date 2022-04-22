package encoding

import (
	"github.com/jacksonCLyu/ridi-config/pkg/config/encoding/toml"
	"github.com/jacksonCLyu/ridi-config/pkg/config/encoding/yml"
	"github.com/jacksonCLyu/ridi-faces/pkg/configer"
)

// DefaultCodec default toml configuration encoder
var DefaultCodec = &tomlCodec{}

// TomlCodec is a codec for the default language.
type tomlCodec struct{}

// Decode the impl of lang.ConfigDecoder
func (c *tomlCodec) Decode(b []byte) (map[string]configer.Field, error) {
	return toml.Decode(b)
}

// Encode the impl of lang.ConfigEncoder
func (c *tomlCodec) Encode(m map[string]configer.Field) ([]byte, error) {
	return toml.Encode(m)
}

type ymlCodec struct{}

func (c *ymlCodec) Decode(b []byte) (map[string]configer.Field, error) {
	return yml.Decode(b)
}

func (c *ymlCodec) Encode(m map[string]configer.Field) ([]byte, error) {
	return yml.Encode(m)
}
