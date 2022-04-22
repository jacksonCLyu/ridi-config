package filesystem

import (
	"io"
	"net/url"
	"os"
	"path/filepath"
)

// ConfigFileSystem config file system struct
type ConfigFileSystem struct {
	opts *options
}

// DefaultFileSystem default file system
var DefaultFileSystem = NewFileSystem()

// NewFileSystem new file system
func NewFileSystem(opts ...Option) *ConfigFileSystem {
	options := &options{
		currentUser:         "",
		versioning:          0,
		proxyHost:           "",
		proxyPort:           0,
		maxHostConnections:  0,
		maxTotalConnections: 0,
	}
	for _, opt := range opts {
		opt.apply(options)
	}
	return &ConfigFileSystem{opts: options}
}

// GetReader returns the reader
func (fs *ConfigFileSystem) GetReader(filePath string) (io.Reader, error) {
	return os.Open(filePath)
}

// GetReaderFromURL returns the reader from URL
func (fs *ConfigFileSystem) GetReaderFromURL(url *url.URL) (io.Reader, error) {
	return os.Open(url.Path)
}

// GetWriter returns the writer
func (fs *ConfigFileSystem) GetWriter(file *os.File) (io.Writer, error) {
	return file, nil
}

// GetWriterFromURL returns the writer from URL
func (fs *ConfigFileSystem) GetWriterFromURL(url *url.URL) (io.Writer, error) {
	return os.Create(url.Path)
}

// GetPath returns the file path
func (fs *ConfigFileSystem) GetPath(file *os.File, url *url.URL, basePath string, fileName string) string {
	if url != nil {
		return url.Path
	}
	if file != nil {
		return file.Name()
	}
	return basePath + fileName
}

// GetBasePath returns the base file path
func (fs *ConfigFileSystem) GetBasePath(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		return path
	}
	return filepath.Dir(abs)
}

// GetFileName returns the file name
func (fs *ConfigFileSystem) GetFileName(path string) string {
	base := filepath.Base(path)
	if filepath.Ext(base) != "" {
		return base[0 : len(filepath.Base(base))-len(filepath.Ext(base))]
	}
	return base
}

// LocateFromURL returns the file infos from URL
func (fs *ConfigFileSystem) LocateFromURL(basePath string, fileName string) *url.URL {
	abs, err := filepath.Abs(basePath + fileName)
	if err != nil {
		return nil
	}
	return &url.URL{
		Scheme: "file",
		Path:   abs,
	}
}

// GetURL return the file URL
func (fs *ConfigFileSystem) GetURL(basePath string, fileName string) *url.URL {
	return fs.LocateFromURL(basePath, fileName)
}
