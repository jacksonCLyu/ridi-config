package filesystem

import (
	"io"
	"net/url"
	"os"
)

// FileSystem file system interface
type FileSystem interface {
	// GetReader get file reader
	GetReader(filePath string) (io.Reader, error)
	// GetReaderFromURL get file reader from URL
	GetReaderFromURL(url *url.URL) (io.Reader, error)
	// GetWriter get file writer
	GetWriter(file *os.File) (io.Writer, error)
	// GetWriterFromURL get file writer from URL
	GetWriterFromURL(url *url.URL) (io.Writer, error)
	// GetPath get file path
	GetPath(file *os.File, url *url.URL, basePath string, fileName string) string
	// GetBasePath get file base path
	GetBasePath(filepath string) string
	// GetFileName get file name
	GetFileName(filepath string) string
	// LocateFromURL locate file from URL
	LocateFromURL(basePath string, fileName string) *url.URL
	// GetURL get file URL
	GetURL(basePath string, fileName string) *url.URL
}
