package config

import (
	"net/url"

	"github.com/jacksonCLyu/ridi-config/pkg/internal/filesystem"
	"github.com/jacksonCLyu/ridi-faces/pkg/configer"
)

type InitOption interface {
	initApply(opts *initOptions)
}

type initOptions struct {
	configurable configer.Configurable
}

type initOptionFunc func(opts *initOptions)

func (f initOptionFunc) initApply(opts *initOptions) {
	f(opts)
}

func WithConfigurable(configurable configer.Configurable) InitOption {
	return initOptionFunc(func(opts *initOptions) {
		opts.configurable = configurable
	})
}

// Option is a struct that holds the options for the config package.
type Option interface {
	apply(opts *options)
}

type options struct {
	filePath          string
	reloadingStrategy configer.ReloadingStrategy
	fileSystem        filesystem.FileSystem
	sourceURL         *url.URL
	encoder           configer.Encoder
	decoder           configer.Decoder
}

// WithReloadingStrategy sets the reloading strategy for the config package.
func WithReloadingStrategy(strategy configer.ReloadingStrategy) Option {
	return reloadingOption{strategy: strategy}
}

// WithSourceURL sets the source URL for the config package.
func WithSourceURL(url *url.URL) Option {
	return sourceOption{url: url}
}

// WithFileSystem sets configuration file system
func WithFileSystem(fileSystem filesystem.FileSystem) Option {
	return fsOption{fileSystem: fileSystem}
}

// WithFilePath sets the file path for the config package.
func WithFilePath(filePath string) Option {
	return filePathOption(filePath)
}

// WithEncoder sets the encoder for the config package.
func WithEncoder(encoder configer.Encoder) Option {
	return encoderOption{encoder: encoder}
}

// WithDecoder sets the decoder for the config package.
func WithDecoder(decoder configer.Decoder) Option {
	return decoderOption{decoder: decoder}
}

type filePathOption string

func (o filePathOption) apply(opts *options) {
	opts.filePath = string(o)
	opts.sourceURL = &url.URL{Scheme: "file", Path: fixPath(opts.filePath)}
}

type reloadingOption struct {
	strategy configer.ReloadingStrategy
}

type sourceOption struct {
	url *url.URL
}

func (o reloadingOption) apply(opts *options) {
	opts.reloadingStrategy = o.strategy
}

func (o sourceOption) apply(opts *options) {
	opts.sourceURL = o.url
}

type fsOption struct {
	fileSystem filesystem.FileSystem
}

func (o fsOption) apply(opts *options) {
	opts.fileSystem = o.fileSystem
}

type encoderOption struct {
	encoder configer.Encoder
}

func (o encoderOption) apply(opts *options) {
	opts.encoder = o.encoder
}

type decoderOption struct {
	decoder configer.Decoder
}

func (o decoderOption) apply(opts *options) {
	opts.decoder = o.decoder
}
