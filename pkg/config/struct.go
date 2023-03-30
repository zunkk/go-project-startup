package config

import (
	"reflect"
	"time"

	"github.com/mitchellh/mapstructure"
)

type Duration time.Duration

func (d *Duration) MarshalText() (text []byte, err error) {
	return []byte(time.Duration(*d).String()), nil
}

func (d *Duration) UnmarshalText(b []byte) error {
	x, err := time.ParseDuration(string(b))
	if err != nil {
		return err
	}
	*d = Duration(x)
	return nil
}

func (d *Duration) ToDuration() time.Duration {
	return time.Duration(*d)
}

func (d *Duration) String() string {
	return time.Duration(*d).String()
}

func StringToTimeDurationHookFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		if t != reflect.TypeOf(Duration(5)) {
			return data, nil
		}

		d, err := time.ParseDuration(data.(string))
		if err != nil {
			return nil, err
		}
		return Duration(d), nil
	}
}

type HTTP struct {
	Port                  int      `mapstructure:"port" toml:"port"`
	MultipartMemory       int64    `mapstructure:"-" toml:"multipart_memory"`
	ReadTimeout           Duration `mapstructure:"read_timeout" toml:"read_timeout"`
	WriteTimeout          Duration `mapstructure:"write_timeout" toml:"write_timeout"`
	TLSEnable             bool     `mapstructure:"tls_enable" toml:"tls_enable"`
	TLSCertFilePath       string   `mapstructure:"tls_cert_file_path" toml:"tls_cert_file_path"`
	TLSKeyFilePath        string   `mapstructure:"tls_key_file_path" toml:"tls_key_file_path"`
	JWTTokenValidDuration Duration `mapstructure:"jwt_token_valid_duration" toml:"jwt_token_valid_duration"`
	JWTTokenHMACKey       string   `mapstructure:"jwt_token_hmac_key" toml:"jwt_token_hmac_key"`
}

type DBInfo struct {
	IP       string `mapstructure:"ip" toml:"ip"`
	Port     uint32 `mapstructure:"port" toml:"port"`
	Username string `mapstructure:"username" toml:"username"`
	Password string `mapstructure:"password" toml:"password"`
	DBName   string `mapstructure:"db_name" toml:"db_name"`
}

type Mongodb struct {
	DBInfo          `mapstructure:",squash"`
	ConnectTimeout  Duration `mapstructure:"connect_timeout" toml:"connect_timeout"`
	MaxPoolSize     int      `mapstructure:"max_pool_size" toml:"max_pool_size"`
	MaxConnIdleTime Duration `mapstructure:"max_conn_idle_time" toml:"max_conn_idle_time"`
}

type Log struct {
	Level        string   `mapstructure:"level" toml:"level"`
	Filename     string   `mapstructure:"file_name" toml:"file_name"`
	MaxAge       Duration `mapstructure:"max_age" toml:"max_age"`
	MaxSizeStr   string   `mapstructure:"max_size" toml:"max_size"`
	MaxSize      int64    `mapstructure:"-" toml:"-"`
	RotationTime Duration `mapstructure:"rotation_time" toml:"rotation_time"`
}
