package strategy

import "github.com/jacksonCLyu/ridi-faces/pkg/configer"

// DefaultManagedReloadingStrategy is the default strategy for reloading managed
var DefaultManagedReloadingStrategy = NewManagedReloadingStrategy()

type managedReloadingStrategy struct {
	configuration configer.FileConfiguration
	// needReload returns true if the configuration should be reloaded.
	needReload bool
}

// NewManagedReloadingStrategy returns a new managed reloading strategy.
func NewManagedReloadingStrategy(opts ...ManagedReloadingOption) configer.ReloadingStrategy {
	options := &managedReloadingOptions{
		fileConfiguration: nil,
	}
	for _, opt := range opts {
		opt.apply(options)
	}
	return &managedReloadingStrategy{
		configuration: options.fileConfiguration,
	}
}

func (s *managedReloadingStrategy) SetConfiguration(fileConfig configer.FileConfiguration) {
	s.configuration = fileConfig
}

func (s managedReloadingStrategy) Init() error {
	return nil
}

func (s managedReloadingStrategy) NeedReloading() (bool, error) {
	return s.needReload, nil
}

func (s managedReloadingStrategy) ReloadingPerformed() error {
	s.needReload = false
	return nil
}

func (s managedReloadingStrategy) Refresh() {
	s.needReload = true
	s.configuration = nil
}
