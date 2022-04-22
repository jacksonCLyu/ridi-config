package strategy

import (
	"os"
	"time"

	"github.com/jacksonCLyu/ridi-faces/pkg/configer"
	"github.com/jacksonCLyu/ridi-utils/utils/fileutil"
	"github.com/pkg/errors"
)

// DefaultTriggerInterval default reloading tirgger interval
const DefaultTriggerInterval = 5000 * time.Millisecond

// DefaultFileChangedReloadingStrategy is a strategy that reloads the configuration
var DefaultFileChangedReloadingStrategy = NewFileChangedReloadingStrategy()

// FileChangedReloadingStrategy file change reloading strategy
type FileChangedReloadingStrategy struct {
	configuration   configer.FileConfiguration
	lastModified    time.Duration
	lastChecked     time.Duration
	triggerInterval time.Duration
	reloading       bool
}

// NewFileChangedReloadingStrategy creates a new FileChangedReloadingStrategy
func NewFileChangedReloadingStrategy(opts ...FileChangedReloadingOption) configer.ReloadingStrategy {
	options := &fileChangedReloadingOptions{
		fileConfiguration: nil,
		triggerInterval:   0,
	}
	for _, opt := range opts {
		opt.apply(options)
	}
	if options.triggerInterval == 0 {
		options.triggerInterval = DefaultTriggerInterval
	}
	return &FileChangedReloadingStrategy{
		configuration:   options.fileConfiguration,
		lastModified:    0,
		lastChecked:     0,
		triggerInterval: options.triggerInterval,
		reloading:       false,
	}
}

// SetConfiguration set configuration
func (s *FileChangedReloadingStrategy) SetConfiguration(configuration configer.FileConfiguration) {
	s.configuration = configuration
}

// Init init fileConfiguration
func (s *FileChangedReloadingStrategy) Init() error {
	return s.updateLastModified()
}

// NeedReloading judge if need reloading the configuration
func (s *FileChangedReloadingStrategy) NeedReloading() (bool, error) {
	if !s.reloading {
		now := time.Now().Local().UnixMilli()
		if now > s.lastChecked.Milliseconds()+s.triggerInterval.Milliseconds() {
			s.lastChecked = time.Duration(now)
			var err error
			if s.reloading, err = s.hasChanged(); err != nil {
				return s.reloading, err
			}
		}
	}
	return s.reloading, nil
}

// ReloadingPerformed the callback of reloading configuration performed
func (s *FileChangedReloadingStrategy) ReloadingPerformed() error {
	return s.updateLastModified()
}

func (s *FileChangedReloadingStrategy) updateLastModified() error {
	defer func() {
		s.reloading = false
	}()
	var file *os.File
	var fileInfo os.FileInfo
	var gErr error
	file, gErr = s.getFile()
	if gErr != nil {
		return gErr
	}
	fileInfo, gErr = file.Stat()
	if gErr != nil {
		return gErr
	}
	modTime := fileInfo.ModTime()
	s.lastModified = time.Duration(modTime.Local().UnixMilli())
	return gErr
}

func (s *FileChangedReloadingStrategy) getFile() (*os.File, error) {
	if s.configuration.GetURL() != nil {
		return s.getFileFromURL()
	}
	return nil, errors.New("file configuration doesn't have `URL` property")
}

func (s *FileChangedReloadingStrategy) getFileFromURL() (*os.File, error) {
	url, err := fileutil.GetFileFromURL(s.configuration.GetURL())
	if err != nil {
		return nil, err
	}
	return url, err
}

func (s *FileChangedReloadingStrategy) hasChanged() (bool, error) {
	var file *os.File
	var fileInfo os.FileInfo
	var gErr error
	file, gErr = s.getFile()
	if gErr != nil {
		return false, gErr
	}
	fileInfo, gErr = file.Stat()
	if gErr != nil {
		return false, gErr
	}
	modTime := fileInfo.ModTime()
	return modTime.Local().UnixMilli() > s.lastModified.Milliseconds(), nil
}
