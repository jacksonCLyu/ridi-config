package config

import (
	"reflect"
	"testing"

	"github.com/jacksonCLyu/ridi-config/pkg/internal/encoding"
)

func TestContainsKey(t *testing.T) {
	encoding.Init()
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

func TestGet(t *testing.T) {
	encoding.Init()
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
