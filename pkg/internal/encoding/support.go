package encoding

import "github.com/jacksonCLyu/ridi-faces/pkg/configer"

type EncType string

const (
	Toml EncType = "toml"
	Yml  EncType = "yml"
	Yaml EncType = "yaml"
)

func (e EncType) String() string {
	return string(e)
}

var SupportSet = make(map[EncType]configer.Codec)

// Init support set
func Init() {
	SupportSet[Toml] = &tomlCodec{}
	SupportSet[Yml] = &ymlCodec{}
	SupportSet[Yaml] = &ymlCodec{}
}

// IsSupport check if support
func IsSupport(t string) bool {
	_, ok := SupportSet[EncType(t)]
	return ok
}

// GetSupport get support codec
func GetSupport(t string) configer.Codec {
	return SupportSet[EncType(t)]
}

// AddSupport add support
func AddSupport(t string, c configer.Codec) {
	SupportSet[EncType(t)] = c
}
