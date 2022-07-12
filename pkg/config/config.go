package config

import (
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/jacksonCLyu/ridi-faces/pkg/configer"
	"github.com/jacksonCLyu/ridi-faces/pkg/env"
	"github.com/jacksonCLyu/ridi-utils/utils/assignutil"
	"github.com/jacksonCLyu/ridi-utils/utils/errcheck"
	"github.com/jacksonCLyu/ridi-utils/utils/rescueutil"
	"github.com/pkg/errors"

	"github.com/jacksonCLyu/ridi-config/pkg/internal/encoding"
	"github.com/jacksonCLyu/ridi-config/pkg/internal/filesystem"
	"github.com/jacksonCLyu/ridi-config/pkg/internal/strategy"
)

func init() {
	// init encoding module
	encoding.Init()
}

// SetDefaultConfig sets the global default configuration
func SetDefaultConfig(configurable configer.Configurable) {
	defaultConfig = configurable
}

// Init init config
func Init(opts ...InitOption) (gErr error) {
	defer rescueutil.Recover(func(err any) {
		if err != nil {
			gErr = err.(error)
		}
	})
	initOpts := &initOptions{}
	for _, opt := range opts {
		opt.initApply(initOpts)
	}
	if initOpts.configurable != nil {
		defaultConfig = initOpts.configurable
	} else {
		defaultConfig = assignutil.Assign(NewConfig())
	}
	return
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
		return filepath.Join(env.AppRootPath(), path[2:])
	}
	return path
}

// L returns the global default configuration
func L() configer.Configurable {
	if defaultConfig == nil {
		_ = Init()
	}
	return defaultConfig
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
	// deep the section deep
	deep int
	// codec codec
	encoder configer.Encoder
	decoder configer.Decoder
}

// NewConfig creates a new configuration
func NewConfig(opts ...Option) (configurable configer.Configurable, err error) {
	defer rescueutil.Recover(func(e any) {
		if e != nil {
			err = e.(error)
		}
	})
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
	// give `this` to reloading strategy
	c.ReloadStrategy.SetConfiguration(c)
	c.ReloadStrategy.Init()
	// auto codec
	ext := filepath.Ext(c.FilePath)[1:]
	if !encoding.IsSupport(ext) {
		if options.encoder == nil && options.decoder == nil {
			err = errors.New("options config `filePath` file ext not support and config `encoder` and `decoder` is empty")
			return
		}
		// reset if given custom encoder or decoder
		if options.encoder != nil {
			c.SetEncoder(options.encoder)
		}
		if options.decoder != nil {
			c.SetDecoder(options.decoder)
		}
	} else {
		supportCodec := encoding.GetSupport(ext)
		if supportCodec == nil {
			err = errors.New("options config `filePath` file ext codec not found")
			return
		}
		c.SetEncoder(supportCodec)
		c.SetDecoder(supportCodec)
	}
	errcheck.CheckAndPanic(c.Load(c.FilePath))
	configurable = c
	return
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
	// require wirte lock
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
	base := filepath.Base(c.FilePath)
	if strings.HasSuffix(base, filepath.Ext(base)) {
		return base[:len(base)-len(filepath.Ext(base))]
	}
	return base
}

func (c *config) GetFilePath() string {
	return c.FilePath
}

func (c *config) GetURL() *url.URL {
	return c.SourceURL
}

func (c *config) SetURL(url *url.URL) {
	c.SourceURL = url
}

func (c *config) Merge(config configer.FileConfiguration) error {
	// TODO merge implementation
	return nil
}

func (c *config) Sync() error {
	// TODO sync implementation
	return nil
}

func (c *config) Watch(paths ...string) error {
	// TODO watch implementation
	return nil
}

func (c *config) ReloadIfNeeded() error {
	if !c.NeedReload() {
		return nil
	}
	return c.Reload()
}

func (c *config) Reload() error {
	defer c.ReloadStrategy.ReloadingPerformed()
	reader, err := c.fileSystem.GetReader(c.FilePath)
	if err != nil {
		return err
	}
	return c.LoadStream(reader)
}

func (c *config) NeedReload() (needReloading bool) {
	defer rescueutil.Recover(func(err any) {
		log.Printf("config need reloading error: %+v\n", errors.WithStack(err.(error)))
		needReloading = false
	})
	c.RLock()
	defer c.RUnlock()
	// if deep has value, don't need reload
	if c.deep > 0 {
		return false
	}
	needReloading = assignutil.Assign(c.ReloadStrategy.NeedReloading())
	return
}

func (c *config) GetReloadStrategy() configer.ReloadingStrategy {
	return c.ReloadStrategy
}

func (c *config) SetReloadStrategy(strategy configer.ReloadingStrategy) {
	c.ReloadStrategy = strategy
}

func (c *config) ContainsKey(key string) bool {
	if err := c.ReloadIfNeeded(); err != nil {
		log.Printf("config contains key error: %+v\n", errors.WithStack(err))
		return false
	}
	return containsKey(c.configMap, key)
}

func (c *config) GetString(key string) (string, error) {
	if !ContainsKey(key) {
		return "", nil
	}
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
	if !ContainsKey(key) {
		return 0, nil
	}
	field, err := c.get(key)
	if err != nil {
		return 0, err
	}
	if field.Type != configer.FieldTypeInt {
		return 0, errors.New("field type is not int")
	}
	return int(field.Value.(int64)), nil
}

func (c *config) GetBool(key string) (bool, error) {
	if !ContainsKey(key) {
		return false, nil
	}
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
	if !ContainsKey(key) {
		return float64(0), nil
	}
	field, err := c.get(key)
	if err != nil {
		return float64(0), err
	}
	if field.Type != configer.FieldTypeFloat {
		return float64(0), errors.New("field type is not float64")
	}
	return field.Value.(float64), nil
}

func (c *config) GetStringSlice(key string) ([]string, error) {
	if !ContainsKey(key) {
		return []string{}, nil
	}
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
	if !ContainsKey(key) {
		return []int{}, nil
	}
	field, err := c.get(key)
	if err != nil {
		return []int{}, err
	}
	if field.Type != configer.FieldTypeIntSlice {
		return []int{}, errors.New("field type is not int slice")
	}
	fv := field.Value.([]int64)
	r := make([]int, len(fv))
	for i, v := range fv {
		r[i] = int(v)
	}
	return r, nil
}

func (c *config) GetBoolSlice(key string) ([]bool, error) {
	if !ContainsKey(key) {
		return []bool{}, nil
	}
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
	if !ContainsKey(key) {
		return []float64{}, nil
	}
	field, err := c.get(key)
	if err != nil {
		return []float64{}, err
	}
	if field.Type != configer.FieldTypeFloatSlice {
		return []float64{}, errors.New("field type is not float64 slice")
	}
	return field.Value.([]float64), nil
}

func (c *config) GetSection(key string) (configer.Configurable, error) {

	field, err := c.get(key)
	if err != nil {
		return nil, err
	}
	if field.Type != configer.FieldTypeSection {
		return nil, errors.New("field type is not Configurable")
	}
	return &config{
		FilePath:       c.FilePath,
		SourceURL:      c.SourceURL,
		ReloadStrategy: c.ReloadStrategy,
		fileSystem:     c.fileSystem,
		configMap:      field.Value.(map[string]configer.Field),
		deep:           c.deep + 1,
		encoder:        c.encoder,
		decoder:        c.decoder,
	}, nil
}

func (c *config) GetInt32(key string) (int32, error) {
	if !ContainsKey(key) {
		return int32(0), nil
	}
	field, err := c.get(key)
	if err != nil {
		return int32(0), err
	}
	t := configer.FieldTypeInt
	if field.Type != t {
		return int32(0), errors.New("field type is not " + t.String())
	}
	return int32(field.Value.(int64)), nil
}

func (c *config) GetInt32Slice(key string) ([]int32, error) {
	if !ContainsKey(key) {
		return []int32{}, nil
	}
	field, err := c.get(key)
	if err != nil {
		return []int32{}, err
	}
	t := configer.FieldTypeIntSlice
	if field.Type != t {
		return []int32{}, errors.New("field type is not " + t.String())
	}
	v := field.Value.([]int64)
	r := make([]int32, len(v))
	for i, v := range v {
		r[i] = int32(v)
	}
	return r, nil
}

func (c *config) GetInt64(key string) (int64, error) {
	if !ContainsKey(key) {
		return int64(0), nil
	}
	field, err := c.get(key)
	if err != nil {
		return int64(0), err
	}
	t := configer.FieldTypeInt
	if field.Type != t {
		return int64(0), errors.New("field type is not " + t.String())
	}
	return field.Value.(int64), nil
}

func (c *config) GetInt64Slice(key string) ([]int64, error) {
	if !ContainsKey(key) {
		return []int64{}, nil
	}
	field, err := c.get(key)
	if err != nil {
		return []int64{}, err
	}
	t := configer.FieldTypeIntSlice
	if field.Type != t {
		return []int64{}, errors.New("field type is not " + t.String())
	}
	return field.Value.([]int64), nil
}

func (c *config) GetUint(key string) (uint, error) {
	if !ContainsKey(key) {
		return uint(0), nil
	}
	field, err := c.get(key)
	if err != nil {
		return uint(0), err
	}
	t := configer.FieldTypeInt
	if field.Type != t {
		return uint(0), errors.New("field type is not " + t.String())
	}
	return uint(field.Value.(int64)), nil
}

func (c *config) GetUintSlice(key string) ([]uint, error) {
	if !ContainsKey(key) {
		return []uint{}, nil
	}
	field, err := c.get(key)
	if err != nil {
		return []uint{}, err
	}
	t := configer.FieldTypeIntSlice
	if field.Type != t {
		return []uint{}, errors.New("field type is not " + t.String())
	}
	v := field.Value.([]int64)
	r := make([]uint, len(v))
	for i, v := range v {
		r[i] = uint(v)
	}
	return r, nil
}

func (c *config) GetUint32(key string) (uint32, error) {
	if !ContainsKey(key) {
		return uint32(0), nil
	}
	field, err := c.get(key)
	if err != nil {
		return uint32(0), err
	}
	t := configer.FieldTypeInt
	if field.Type != t {
		return uint32(0), errors.New("field type is not " + t.String())
	}
	return uint32(field.Value.(int64)), nil
}

func (c *config) GetUint32Slice(key string) ([]uint32, error) {
	if !ContainsKey(key) {
		return []uint32{}, nil
	}
	field, err := c.get(key)
	if err != nil {
		return []uint32{}, err
	}
	t := configer.FieldTypeIntSlice
	if field.Type != t {
		return []uint32{}, errors.New("field type is not " + t.String())
	}
	v := field.Value.([]int64)
	r := make([]uint32, len(v))
	for i, v := range v {
		r[i] = uint32(v)
	}
	return r, nil
}

func (c *config) GetUint64(key string) (uint64, error) {
	if !ContainsKey(key) {
		return uint64(0), nil
	}
	field, err := c.get(key)
	if err != nil {
		return uint64(0), err
	}
	t := configer.FieldTypeInt
	if field.Type != t {
		return uint64(0), errors.New("field type is not " + t.String())
	}
	return uint64(field.Value.(int64)), nil
}

func (c *config) GetUint64Slice(key string) ([]uint64, error) {
	if !ContainsKey(key) {
		return []uint64{}, nil
	}
	field, err := c.get(key)
	if err != nil {
		return []uint64{}, err
	}
	t := configer.FieldTypeIntSlice
	if field.Type != t {
		return []uint64{}, errors.New("field type is not " + t.String())
	}
	v := field.Value.([]int64)
	r := make([]uint64, len(v))
	for i, v := range v {
		r[i] = uint64(v)
	}
	return r, nil
}

func (c *config) GetFloat32(key string) (float32, error) {
	if !ContainsKey(key) {
		return float32(0), nil
	}
	field, err := c.get(key)
	if err != nil {
		return float32(0), err
	}
	t := configer.FieldTypeFloat
	if field.Type != t {
		return float32(0), errors.New("field type is not " + t.String())
	}
	return float32(field.Value.(float64)), nil
}

func (c *config) GetFloat32Slice(key string) ([]float32, error) {
	if !ContainsKey(key) {
		return []float32{}, nil
	}
	field, err := c.get(key)
	if err != nil {
		return []float32{}, err
	}
	t := configer.FieldTypeFloatSlice
	if field.Type != t {
		return []float32{}, errors.New("field type is not " + t.String())
	}
	v := field.Value.([]float64)
	r := make([]float32, len(v))
	for i, v := range v {
		r[i] = float32(v)
	}
	return r, nil
}

func (c *config) GetDuration(key string) (time.Duration, error) {
	if !ContainsKey(key) {
		return time.Duration(0), nil
	}
	field, err := c.get(key)
	if err != nil {
		return time.Duration(0), err
	}
	t := configer.FieldTypeDuration
	if field.Type != t {
		return time.Duration(0), errors.New("field type is not " + t.String())
	}
	return field.Value.(time.Duration), nil
}

func (c *config) GetTime(key string) (time.Time, error) {
	if !ContainsKey(key) {
		return time.Now().Local(), nil
	}
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
	if err := c.ReloadIfNeeded(); err != nil {
		return configer.Field{}, err
	}
	c.RLock()
	defer c.RUnlock()
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
