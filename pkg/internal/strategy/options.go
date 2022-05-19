package strategy

import (
	"github.com/jacksonCLyu/ridi-faces/pkg/configer"
)

type fileChangedReloadingOptions struct {
	fileConfiguration configer.FileConfiguration
	triggerInterval   int64
}

type managedReloadingOptions struct {
	fileConfiguration configer.FileConfiguration
}

// FileChangedReloadingOption option interface for file changed reloading strategy
type FileChangedReloadingOption interface {
	apply(opts *fileChangedReloadingOptions)
}

// ManagedReloadingOption option interface for managed reloading streategy
type ManagedReloadingOption interface {
	apply(opts *managedReloadingOptions)
}

// WithFileConfiguration sets the file configuration
func WithFileConfiguration(configuration configer.FileConfiguration) FileChangedReloadingOption {
	return fileChangedConfigurationOption{configuration: configuration}
}

// WithTriggerInterval sets the trigger interval
func WithTriggerInterval(triggerInterval int64) FileChangedReloadingOption {
	return triggerIntervalOption(triggerInterval)
}

// WithManagedConfiguration sets the file configuration to use for the reloading strategy.
func WithManagedConfiguration(configuration configer.FileConfiguration) ManagedReloadingOption {
	return managedConfigurationOption{configuration: configuration}
}

type managedConfigurationOption struct {
	configuration configer.FileConfiguration
}

func (o managedConfigurationOption) apply(opts *managedReloadingOptions) {
	opts.fileConfiguration = o.configuration
}

type fileChangedConfigurationOption struct {
	configuration configer.FileConfiguration
}

func (o fileChangedConfigurationOption) apply(opts *fileChangedReloadingOptions) {
	opts.fileConfiguration = o.configuration
}

type triggerIntervalOption int64

func (o triggerIntervalOption) apply(opts *fileChangedReloadingOptions) {
	opts.triggerInterval = int64(o)
}
