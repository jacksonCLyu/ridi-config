package config

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/jacksonCLyu/ridi-faces/pkg/env"
)

func TestInit(t *testing.T) {
	newConfig, err := NewConfig(WithFilePath("./testdata/test.toml"))
	if err != nil {
		t.Error(err)
	}
	err = Init(WithConfigurable(newConfig))
	if err != nil {
		t.Error(err)
	}
	type args struct {
		key string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "containsKey1",
			args: args{
				key: "servers",
			},
			want: true,
		},
		{
			name: "containsKey2",
			args: args{
				key: "servers.alpha",
			},
			want: true,
		},
		{
			name: "containsKey3",
			args: args{
				key: "servers.ip",
			},
			want: false,
		},
		{
			name: "containsKey4",
			args: args{
				key: "servers.alpha.ip",
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ContainsKey(tt.args.key); got != tt.want {
				t.Errorf("ContainsKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContainsKey(t *testing.T) {
	newConfig, err := NewConfig(WithFilePath("./testdata/test.toml"))
	if err != nil {
		t.Error(err)
	}
	type args struct {
		key string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "containsKey1",
			args: args{
				key: "servers",
			},
			want: true,
		},
		{
			name: "containsKey2",
			args: args{
				key: "servers.alpha",
			},
			want: true,
		},
		{
			name: "containsKey3",
			args: args{
				key: "servers.ip",
			},
			want: false,
		},
		{
			name: "containsKey4",
			args: args{
				key: "servers.alpha.ip",
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newConfig.ContainsKey(tt.args.key); got != tt.want {
				t.Errorf("ContainsKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContainsKeyYaml(t *testing.T) {
	configPath := filepath.Join(env.AppRootPath(), "testdata", "test.yaml")
	newConfig, err := NewConfig(WithFilePath(configPath))
	if err != nil {
		t.Error(err)
	}
	type args struct {
		key string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "containsKey1",
			args: args{
				key: "servers",
			},
			want: true,
		},
		{
			name: "containsKey2",
			args: args{
				key: "servers.alpha",
			},
			want: true,
		},
		{
			name: "containsKey3",
			args: args{
				key: "servers.ip",
			},
			want: false,
		},
		{
			name: "containsKey4",
			args: args{
				key: "servers.alpha.ip",
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newConfig.ContainsKey(tt.args.key); got != tt.want {
				t.Errorf("ContainsKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContainsKeyYml(t *testing.T) {
	configPath := filepath.Join(env.AppRootPath(), "testdata", "test.yml")
	newConfig, err := NewConfig(WithFilePath(configPath))
	if err != nil {
		t.Error(err)
	}
	type args struct {
		key string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "containsKey1",
			args: args{
				key: "servers",
			},
			want: true,
		},
		{
			name: "containsKey2",
			args: args{
				key: "servers.alpha",
			},
			want: true,
		},
		{
			name: "containsKey3",
			args: args{
				key: "servers.ip",
			},
			want: false,
		},
		{
			name: "containsKey4",
			args: args{
				key: "servers.alpha.ip",
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newConfig.ContainsKey(tt.args.key); got != tt.want {
				t.Errorf("ContainsKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGet(t *testing.T) {
	newConfig, err := NewConfig(WithFilePath("./testdata/test.toml"))
	if err != nil {
		t.Errorf("load config failed: %v", err)
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    any
		wantErr bool
	}{
		{
			name: "GetString",
			args: args{
				key: "servers.alpha.ip",
			},
			want:    "10.0.0.1",
			wantErr: false,
		},
		{
			name: "GetString2",
			args: args{
				key: "servers.abc",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "GetInt64",
			args: args{
				key: "database.connection_max",
			},
			want:    int64(5000),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newConfig.Get(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Get() got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestGetYaml(t *testing.T) {
	newConfig, err := NewConfig(WithFilePath("./testdata/test.yaml"))
	if err != nil {
		t.Errorf("load config failed: %v", err)
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    any
		wantErr bool
	}{
		{
			name: "GetString",
			args: args{
				key: "servers.alpha.ip",
			},
			want:    "10.0.0.1",
			wantErr: false,
		},
		{
			name: "GetString2",
			args: args{
				key: "servers.abc",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "GetInt",
			args: args{
				key: "database.connection_max",
			},
			want:    5000,
			wantErr: false,
		},
		{
			name: "GetSlice1",
			args: args{
				key: "database.ports",
			},
			want:    []int64{8001, 8002},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newConfig.Get(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Get() got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestGetIntSlice(t *testing.T) {
	newConfig, err := NewConfig(WithFilePath("./testdata/test.yaml"))
	if err != nil {
		t.Errorf("load config failed: %v", err)
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    any
		wantErr bool
	}{
		{
			name: "GetIntSlice",
			args: args{
				key: "database.ports",
			},
			want:    []int{8001, 8002},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newConfig.GetIntSlice(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Get() got = %v, want %v", got, tt.want)
				}
			}
		})
	}

	tests2 := []struct {
		name    string
		args    args
		want    any
		wantErr bool
	}{
		{
			name: "GetInt32Slice",
			args: args{
				key: "database.ports",
			},
			want:    []int32{8001, 8002},
			wantErr: false,
		},
	}

	for _, tt := range tests2 {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newConfig.GetInt32Slice(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Get() got = %v, want %v", got, tt.want)
				}
			}
		})
	}

	tests3 := []struct {
		name    string
		args    args
		want    any
		wantErr bool
	}{
		{
			name: "GetInt64Slice",
			args: args{
				key: "database.ports",
			},
			want:    []int64{8001, 8002},
			wantErr: false,
		},
	}

	for _, tt := range tests3 {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newConfig.GetInt64Slice(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Get() got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestGetYml(t *testing.T) {
	newConfig, err := NewConfig(WithFilePath("./testdata/test.yml"))
	if err != nil {
		t.Errorf("load config failed: %v", err)
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    any
		wantErr bool
	}{
		{
			name: "GetString",
			args: args{
				key: "servers.alpha.ip",
			},
			want:    "10.0.0.1",
			wantErr: false,
		},
		{
			name: "GetString2",
			args: args{
				key: "servers.abc",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "GetInt",
			args: args{
				key: "database.connection_max",
			},
			want:    5000,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newConfig.Get(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Get() got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestGetSection(t *testing.T) {
	newConfig, err := NewConfig(WithFilePath("./testdata/test.toml"))
	if err != nil {
		t.Errorf("load config failed: %v", err)
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "GetSection",
			args: args{
				key: "servers",
			},
			want:    "10.0.0.1",
			wantErr: false,
		},
		{
			name: "GetSection2",
			args: args{
				key: "servers.abc",
			},
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newConfig.GetSection(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSection() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got == nil {
					t.Errorf("GetSection() got = %v, want %v", got, tt.want)
					return
				}
				ip, _ := got.GetString("alpha.ip")
				if !reflect.DeepEqual(ip, tt.want) {
					t.Errorf("GetSection() got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestGetSectionYaml(t *testing.T) {
	newConfig, err := NewConfig(WithFilePath("./testdata/test.yaml"))
	if err != nil {
		t.Errorf("load config failed: %v", err)
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "GetSection",
			args: args{
				key: "servers",
			},
			want:    "10.0.0.1",
			wantErr: false,
		},
		{
			name: "GetSection2",
			args: args{
				key: "servers.abc",
			},
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newConfig.GetSection(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSection() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got == nil {
					t.Errorf("GetSection() got = %v, want %v", got, tt.want)
					return
				}
				ip, _ := got.GetString("alpha.ip")
				if !reflect.DeepEqual(ip, tt.want) {
					t.Errorf("GetSection() got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestGetSectionYml(t *testing.T) {
	newConfig, err := NewConfig(WithFilePath("./testdata/test.yml"))
	if err != nil {
		t.Errorf("load config failed: %v", err)
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "GetSection",
			args: args{
				key: "servers",
			},
			want:    "10.0.0.1",
			wantErr: false,
		},
		{
			name: "GetSection2",
			args: args{
				key: "servers.abc",
			},
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newConfig.GetSection(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSection() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got == nil {
					t.Errorf("GetSection() got = %v, want %v", got, tt.want)
					return
				}
				ip, _ := got.GetString("alpha.ip")
				if !reflect.DeepEqual(ip, tt.want) {
					t.Errorf("GetSection() got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
