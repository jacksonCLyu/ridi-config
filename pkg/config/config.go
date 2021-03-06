package config

import (
	"errors"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/jacksonCLyu/ridi-config/pkg/config/encoding"
	"github.com/jacksonCLyu/ridi-config/pkg/config/filesystem"
	"github.com/jacksonCLyu/ridi-config/pkg/config/strategy"
	"github.com/jacksonCLyu/ridi-faces/pkg/configer"
	"github.com/jacksonCLyu/ridi-faces/pkg/env"
)

// Init init config
func Init(opts ...InitOption) error {
	var gErr error
	once.Do(func() {
		encoding.Init()
		initOpts := &initOptions{}
		for _, opt := range opts {
			opt.initApply(initOpts)
		}
		if initOpts.configurable != nil {
			DefaultConfig = initOpts.configurable
			return
		}
		DefaultConfig, gErr = NewConfig()
		if gErr != nil {
			return
		}
	})
	return gErr
}

// DefaultOptions returns the default configuration options.
func DefaultOptions() *options {
	return &options{
		filePath:          "./config.toml",
		reloadingStrategy: strategy.DefaultFileChangedReloadingStrategy,
		fileSystem:        filesystem.DefaultFileSystem,
		sourceURL:         &url.URL{Scheme: "file", Path: fixPath("./config.toml")},
		encoder:           encoding.DefaultCodec,
		decoder:           encoding.DefaultCodec,
	}
}

func fixPath(path string) string {
	if strings.HasPrefix(path, "."+string(filepath.Separator)) {
		rootPath := env.AppRootPath()
		return strings.Join([]string{rootPath, path[2:]}, string(filepath.Separator))
	}
	return path
}

// L returns the global default configuration
func L() configer.Configurable {
	if DefaultConfig == nil {
		_ = Init()
	}
	return DefaultConfig
}

var _ configer.Configurable = (*config)(nil)
var _ configer.FileConfiguration = (*config)(nil)

type config struct {
	// lock for syncing
	sync.RWMutex
	// FilePath file path
	FilePath string
	// ReloadStrategy reload strategy
	ReloadStrategy configer.ReloadingStrategy
	// SourceURL source url
	SourceURL *url.URL
	// fileSystem configuration file system
	fileSystem filesystem.FileSystem
	// configMap config map
	configMap map[string]configer.Field
	// codec codec
	encoder configer.Encoder
	decoder configer.Decoder
}

// NewConfig creates a new configuration
func NewConfig(opts ...Option) (configer.Configurable, error) {
	options := DefaultOptions()
	for _, opt := range opts {
		opt.apply(options)
	}
	if options.filePath == "" {
		return nil, errors.New("options config `filePath` is empty")
	}
	if filepath.Ext(options.filePath) == "" || filepath.Ext(options.filePath) == "." {
		return nil, errors.New("options config `filePath` file ext not found")
	}
	c := &config{
		FilePath:       options.filePath,
		ReloadStrategy: options.reloadingStrategy,
		SourceURL:      options.sourceURL,
		fileSystem:     options.fileSystem,
		configMap:      make(map[string]configer.Field),
	}
	// auto codec
	ext := filepath.Ext(c.FilePath)
	ext = ext[1:]
	if !encoding.IsSupport(ext) {
		if options.encoder == nil && options.decoder == nil {
			return nil, errors.New("options config `filePath` file ext not support")
		}
		// reset if given custom encoder or decoder
		if options.encoder != nil {
			c.SetEncoder(options.encoder)
		}
		if options.decoder != nil {
			c.SetDecoder(options.decoder)
		}
	}
	supportCodec := encoding.GetSupport(ext)
	c.encoder = supportCodec
	c.decoder = supportCodec
	err := c.Load(c.FilePath)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *config) GetEncoder() configer.Encoder {
	return c.encoder
}

func (c *config) SetEncoder(encoder configer.Encoder) {
	c.encoder = encoder
}

func (c *config) GetDecoder() configer.Decoder {
	return c.decoder
}

func (c *config) SetDecoder(decoder configer.Decoder) {
	c.decoder = decoder
}

//implements for FileConfiguration

func (c *config) Load(path string) error {
	reader, err := c.fileSystem.GetReader(path)
	if err != nil {
		return err
	}
	return c.LoadStream(reader)
}

// LoadRemote load configuration from url
func (c *config) LoadRemote(url *url.URL) error {
	if c.SourceURL == nil {
		if strings.EqualFold(c.FilePath, "") {
			c.FilePath = url.String()
		}
		c.SourceURL = url
	}
	is, err := c.fileSystem.GetReaderFromURL(url)
	if err != nil {
		return err
	}
	return c.LoadStream(is)
}

func (c *config) LoadStream(r io.Reader) error {
	c.Lock()
	defer c.Unlock()
	c.configMap = make(map[string]configer.Field)
	all, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	c.configMap, err = c.decoder.Decode(all)
	return err
}

func (c *config) Save(path string) error {
	open, err := os.Open(path)
	if err != nil {
		return err
	}
	if writer, err := c.fileSystem.GetWriter(open); err != nil {
		return err
	} else {
		return c.SaveStream(writer)
	}
}

func (c *config) SaveStream(writer io.Writer) error {
	c.Lock()
	defer c.Unlock()
	if all, err := c.encoder.Encode(c.configMap); err != nil {
		return err
	} else {
		_, err := writer.Write(all)
		return err
	}
}

func (c *config) SaveRemote(url *url.URL) error {
	fromURL, err := c.fileSystem.GetWriterFromURL(url)
	if err != nil {
		return err
	}
	return c.SaveStream(fromURL)
}

func (c *config) GetFileName() string {
	c.RLock()
	defer c.RUnlock()
	base := filepath.Base(c.FilePath)
	if strings.HasSuffix(base, filepath.Ext(base)) {
		return base[:len(base)-len(filepath.Ext(base))]
	}
	return base
}

func (c *config) GetFilePath() string {
	c.RLock()
	defer c.RUnlock()
	return c.FilePath
}

func (c *config) GetURL() *url.URL {
	c.RLock()
	defer c.RUnlock()
	return c.SourceURL
}

func (c *config) SetURL(url *url.URL) {
	c.Lock()
	defer c.Unlock()
	c.SourceURL = url
}

func (c *config) Merge(config configer.FileConfiguration) error {
	c.Lock()
	defer c.Unlock()
	// TODO merge implementation
	return nil
}

func (c *config) Sync() error {
	c.Lock()
	defer c.Unlock()
	// TODO sync implementation
	return nil
}

func (c *config) Watch(paths ...string) error {
	c.Lock()
	defer c.Unlock()
	// TODO watch implementation
	return nil
}

func (c *config) Reload() error {
	c.RLock()
	defer c.RUnlock()
	var needReload bool
	var err error
	if needReload, err = c.ReloadStrategy.NeedReloading(); err != nil {
		return err
	}
	if !needReload {
		return nil
	}
	c.Lock()
	defer c.Unlock()
	reader, err := c.fileSystem.GetReader(c.FilePath)
	if err != nil {
		return err
	}
	return c.LoadStream(reader)
}

func (c *config) GetReloadStrategy() configer.ReloadingStrategy {
	c.RLock()
	defer c.RUnlock()
	return c.ReloadStrategy
}

func (c *config) SetReloadStrategy(strategy configer.ReloadingStrategy) {
	c.Lock()
	defer c.Unlock()
	c.ReloadStrategy = strategy
}

func (c *config) ContainsKey(key string) bool {
	c.RLock()
	defer c.RUnlock()
	return containsKey(c.configMap, key)
}

func (c *config) GetString(key string) (string, error) {
	c.RLock()
	defer c.RUnlock()
	field, err := c.get(key)
	if err != nil {
		return "", err
	}
	t := configer.FieldTypeString
	if field.Type != t {
		return "", errors.New("field type is not " + t.String())
	}
	return field.Value.(string), nil
}

func (c *config) GetInt(key string) (int, error) {
	c.RLock()
	defer c.RUnlock()
	field, err := c.get(key)
	if err != nil {
		return 0, err
	}
	if field.Type != configer.FieldTypeInt {
		return 0, errors.New("field type is not int")
	}
	return field.Value.(int), nil
}

func (c *config) GetBool(key string) (bool, error) {
	c.RLock()
	defer c.RUnlock()
	field, err := c.get(key)
	if err != nil {
		return false, err
	}
	if field.Type != configer.FieldTypeBool {
		return false, errors.New("field type is not bool")
	}
	return field.Value.(bool), nil
}

func (c *config) GetFloat64(key string) (float64, error) {
	c.RLock()
	defer c.RUnlock()
	field, err := c.get(key)
	if err != nil {
		return 0.0, err
	}
	if field.Type != configer.FieldTypeFloat64 {
		return 0.0, errors.New("field type is not float64")
	}
	return field.Value.(float64), nil
}

func (c *config) GetStringSlice(key string) ([]string, error) {
	c.RLock()
	defer c.RUnlock()
	field, err := c.get(key)
	if err != nil {
		return []string{}, err
	}
	if field.Type != configer.FieldTypeStringSlice {
		return []string{}, errors.New("field type is not string slice")
	}
	return field.Value.([]string), nil
}

func (c *config) GetIntSlice(key string) ([]int, error) {
	c.RLock()
	defer c.RUnlock()
	field, err := c.get(key)
	if err != nil {
		return []int{}, err
	}
	if field.Type != configer.FieldTypeIntSlice {
		return []int{}, errors.New("field type is not int slice")
	}
	return field.Value.([]int), nil
}

func (c *config) GetBoolSlice(key string) ([]bool, error) {
	c.RLock()
	defer c.RUnlock()
	field, err := c.get(key)
	if err != nil {
		return []bool{}, err
	}
	if field.Type != configer.FieldTypeBoolSlice {
		return []bool{}, errors.New("field type is not bool slice")
	}
	return field.Value.([]bool), nil
}

func (c *config) GetFloat64Slice(key string) ([]float64, error) {
	c.RLock()
	defer c.RUnlock()
	field, err := c.get(key)
	if err != nil {
		return []float64{}, err
	}
	if field.Type != configer.FieldTypeFloat64Slice {
		return []float64{}, errors.New("field type is not float64 slice")
	}
	return field.Value.([]float64), nil
}

func (c *config) GetSection(key string) (configer.Configurable, error) {
	c.RLock()
	defer c.RUnlock()
	field, err := c.get(key)
	if err != nil {
		return nil, err
	}
	if field.Type != configer.FieldTypeSection {
		return nil, errors.New("field type is not Configurable")
	}
	return field.Value.(*config), nil
}

func (c *config) GetInt32(key string) (int32, error) {
	c.RLock()
	defer c.RUnlock()
	field, err := c.get(key)
	if err != nil {
		return 0, err
	}
	t := configer.FieldTypeInt32
	if field.Type != t {
		return 0, errors.New("field type is not " + t.String())
	}
	return field.Value.(int32), nil
}

func (c *config) GetInt32Slice(key string) ([]int32, error) {
	c.RLock()
	defer c.RUnlock()
	field, err := c.get(key)
	if err != nil {
		return []int32{}, err
	}
	t := configer.FieldTypeInt32Slice
	if field.Type != t {
		return []int32{}, errors.New("field type is not " + t.String())
	}
	return field.Value.([]int32), nil
}

func (c *config) GetInt64(key string) (int64, error) {
	c.RLock()
	defer c.RUnlock()
	field, err := c.get(key)
	if err != nil {
		return 0, err
	}
	t := configer.FieldTypeInt64
	if field.Type != t {
		return 0, errors.New("field type is not " + t.String())
	}
	return field.Value.(int64), nil
}

func (c *config) GetInt64Slice(key string) ([]int64, error) {
	c.RLock()
	defer c.RUnlock()
	field, err := c.get(key)
	if err != nil {
		return []int64{}, err
	}
	t := configer.FieldTypeInt64Slice
	if field.Type != t {
		return []int64{}, errors.New("field type is not " + t.String())
	}
	return field.Value.([]int64), nil
}

func (c *config) GetUint(key string) (uint, error) {
	c.RLock()
	defer c.RUnlock()
	field, err := c.get(key)
	if err != nil {
		return 0, err
	}
	t := configer.FieldTypeUint
	if field.Type != t {
		return 0, errors.New("field type is not " + t.String())
	}
	return field.Value.(uint), nil
}

func (c *config) GetUintSlice(key string) ([]uint, error) {
	c.RLock()
	defer c.RUnlock()
	field, err := c.get(key)
	if err != nil {
		return []uint{}, err
	}
	t := configer.FieldTypeUintSlice
	if field.Type != t {
		return []uint{}, errors.New("field type is not " + t.String())
	}
	return field.Value.([]uint), nil
}

func (c *config) GetUint32(key string) (uint32, error) {
	c.RLock()
	defer c.RUnlock()
	field, err := c.get(key)
	if err != nil {
		return 0, err
	}
	t := configer.FieldTypeUint32
	if field.Type != t {
		return 0, errors.New("field type is not " + t.String())
	}
	return field.Value.(uint32), nil
}

func (c *config) GetUint32Slice(key string) ([]uint32, error) {
	c.RLock()
	defer c.RUnlock()
	field, err := c.get(key)
	if err != nil {
		return []uint32{}, err
	}
	t := configer.FieldTypeUint32Slice
	if field.Type != t {
		return []uint32{}, errors.New("field type is not " + t.String())
	}
	return field.Value.([]uint32), nil
}

func (c *config) GetUint64(key string) (uint64, error) {
	c.RLock()
	defer c.RUnlock()
	field, err := c.get(key)
	if err != nil {
		return 0, err
	}
	t := configer.FieldTypeUint64
	if field.Type != t {
		return 0, errors.New("field type is not " + t.String())
	}
	return field.Value.(uint64), nil
}

func (c *config) GetUint64Slice(key string) ([]uint64, error) {
	c.RLock()
	defer c.RUnlock()
	field, err := c.get(key)
	if err != nil {
		return []uint64{}, err
	}
	t := configer.FieldTypeUint64Slice
	if field.Type != t {
		return []uint64{}, errors.New("field type is not " + t.String())
	}
	return field.Value.([]uint64), nil
}

func (c *config) GetFloat32(key string) (float32, error) {
	c.RLock()
	defer c.RUnlock()
	field, err := c.get(key)
	if err != nil {
		return 0.0, err
	}
	t := configer.FieldTypeFloat32
	if field.Type != t {
		return 0.0, errors.New("field type is not " + t.String())
	}
	return field.Value.(float32), nil
}

func (c *config) GetFloat32Slice(key string) ([]float32, error) {
	c.RLock()
	defer c.RUnlock()
	field, err := c.get(key)
	if err != nil {
		return []float32{}, err
	}
	t := configer.FieldTypeFloat32Slice
	if field.Type != t {
		return []float32{}, errors.New("field type is not " + t.String())
	}
	return field.Value.([]float32), nil
}

func (c *config) GetDuration(key string) (time.Duration, error) {
	c.RLock()
	defer c.RUnlock()
	field, err := c.get(key)
	if err != nil {
		return 0, err
	}
	t := configer.FieldTypeDuration
	if field.Type != t {
		return 0, errors.New("field type is not " + t.String())
	}
	return field.Value.(time.Duration), nil
}

func (c *config) GetTime(key string) (time.Time, error) {
	c.RLock()
	defer c.RUnlock()
	field, err := c.get(key)
	if err != nil {
		return time.Now().Local(), err
	}
	t := configer.FieldTypeTime
	if field.Type != t {
		return time.Now().Local(), errors.New("field type is not " + t.String())
	}
	return field.Value.(time.Time), nil
}

func (c *config) Get(key string) (any, error) {
	c.RLock()
	defer c.RUnlock()
	field, err := c.get(key)
	if err != nil {
		return nil, err
	}
	return field.Value, nil
}

func (c *config) Set(key string, value any) error {
	c.Lock()
	defer c.Unlock()
	c.configMap[key] = configer.Atof(value)
	return c.Save(c.FilePath)
}

func containsKey(configMap map[string]configer.Field, key string) bool {
	if strings.Contains(key, ".") {
		index := strings.Index(key, ".")
		parentKey := key[:index]
		v := configMap[parentKey]
		if v.Type != configer.FieldTypeSection {
			return false
		}
		subKey := key[index+1:]
		return containsKey(v.Value.(map[string]configer.Field), subKey)
	}
	_, ok := configMap[key]
	return ok
}

// get This method acquires the lock by default
func (c *config) get(key string) (configer.Field, error) {
	return getRecursive(c.configMap, key)
}

func getRecursive(configMap map[string]configer.Field, key string) (configer.Field, error) {
	if strings.Contains(key, ".") {
		index := strings.Index(key, ".")
		parentKey := key[:index]
		v := configMap[parentKey]
		if v.Type != configer.FieldTypeSection {
			return configer.Field{}, errors.New("config not found for key:`" + key + "`")
		}
		subKey := key[index+1:]
		return getRecursive(v.Value.(map[string]configer.Field), subKey)
	}
	if !containsKey(configMap, key) {
		return configer.Field{}, errors.New("config not found for key:`" + key + "`")
	}
	return configMap[key], nil
}
