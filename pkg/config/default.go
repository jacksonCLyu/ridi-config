package config

import (
	"time"

	"github.com/jacksonCLyu/ridi-faces/pkg/configer"
)

// DefaultConfig default config for devup
var DefaultConfig configer.Configurable

// ContainsKey returns true if the key is in the config
func ContainsKey(key string) bool {
	return L().ContainsKey(key)
}

// GetString returns the string value of the key
func GetString(key string) (string, error) {
	return L().GetString(key)
}

// GetStringSlice returns the string value of the key
func GetStringSlice(key string) ([]string, error) {
	return L().GetStringSlice(key)
}

// GetBool returns the bool value of the key
func GetBool(key string) (bool, error) {
	return L().GetBool(key)
}

// GetBoolSlice returns the bool slice value of the key
func GetBoolSlice(key string) ([]bool, error) {
	return L().GetBoolSlice(key)
}

// GetInt returns the int value of the key
func GetInt(key string) (int, error) {
	return L().GetInt(key)
}

// GetIntSlice returns the int slice value of the key
func GetIntSlice(key string) ([]int, error) {
	return L().GetIntSlice(key)
}

// GetInt32 returns the int32 value of the key
func GetInt32(key string) (int32, error) {
	return L().GetInt32(key)
}

// GetInt32Slice returns the int32 slice value of the key
func GetInt32Slice(key string) ([]int32, error) {
	return L().GetInt32Slice(key)
}

// GetInt64 returns the int64 value of the key
func GetInt64(key string) (int64, error) {
	return L().GetInt64(key)
}

// GetInt64Slice returns the int64 slice value of the key
func GetInt64Slice(key string) ([]int64, error) {
	return L().GetInt64Slice(key)
}

// GetUint returns the uint value of the key
func GetUint(key string) (uint, error) {
	return L().GetUint(key)
}

// GetUint32 returns the uint32 value of the key
func GetUint32(key string) (uint32, error) {
	return L().GetUint32(key)
}

// GetUint32Slice returns the uint32 slice value of the key
func GetUint32Slice(key string) ([]uint32, error) {
	return L().GetUint32Slice(key)
}

// GetUint64 returns the uint64 value of the key
func GetUint64(key string) (uint64, error) {
	return L().GetUint64(key)
}

// GetUint64Slice returns the uint64 slice value of the key
func GetUint64Slice(key string) ([]uint64, error) {
	return L().GetUint64Slice(key)
}

// GetFloat32 returns the float32 value of the key
func GetFloat32(key string) (float32, error) {
	return L().GetFloat32(key)
}

// GetFloat32Slice returns the float32 slice value of the key
func GetFloat32Slice(key string) ([]float32, error) {
	return L().GetFloat32Slice(key)
}

// GetFloat64 returns the float64 value of the key
func GetFloat64(key string) (float64, error) {
	return L().GetFloat64(key)
}

// GetFloat64Slice returns the float64 slice value of the key
func GetFloat64Slice(key string) ([]float64, error) {
	return L().GetFloat64Slice(key)
}

// GetDuration returns the duration value of the key
func GetDuration(key string) (time.Duration, error) {
	return L().GetDuration(key)
}

// GetTime returns the time value of the key
func GetTime(key string) (time.Time, error) {
	return L().GetTime(key)
}

// GetSection returns the section value of the key
func GetSection(key string) (configer.Configurable, error) {
	return L().GetSection(key)
}

// Get returns the value of the key
func Get(key string) (any, error) {
	return L().Get(key)
}

// Set sets the value of the key
func Set(key string, value interface{}) error {
	return L().Set(key, value)
}
