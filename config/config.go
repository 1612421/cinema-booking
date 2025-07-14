package config

import (
	"bytes"
	_ "embed"
	"fmt"
	"strings"
	"time"

	cacheredis "github.com/1612421/cinema-booking/pkg/go-kit/cache/redis"
	"github.com/1612421/cinema-booking/pkg/go-kit/database/mysql"
	"github.com/1612421/cinema-booking/pkg/go-kit/log"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Config struct {
	Service        *Service           `json:"service" mapstructure:"service"`
	MySQL          *mysql.Config      `json:"mysql" mapstructure:"mysql"`
	Redis          *cacheredis.Config `json:"redis" mapstructure:"redis"`
	Log            *log.Config        `json:"log,omitempty" mapstructure:"log"`
	SessionService *ExternalService   `json:"session_service,omitempty" mapstructure:"session_service"`
	Auth           AuthConfig         `json:"auth" mapstructure:"auth"`
}

type AuthConfig struct {
	Secret   string `json:"secret" mapstructure:"secret"`
	ExpireIn int    `json:"expire_in" mapstructure:"expire_in"`
}

type Service struct {
	Name        string  `json:"name" mapstructure:"name"`
	Version     string  `json:"version,omitempty" mapstructure:"version"`
	HTTP        *Listen `json:"http,omitempty" mapstructure:"http"`
	Environment string  `json:"env,omitempty" mapstructure:"env"`
}

// Listen represents a network end point address.
type Listen struct {
	Host string `json:"host,omitempty" mapstructure:"host"`
	Port int    `json:"port,omitempty" mapstructure:"port"`
}

type ExternalService struct {
	Address       string        `json:"address,omitempty" mapstructure:"address"`
	ClientID      string        `json:"client_id,omitempty" mapstructure:"client_id"`
	ClientKey     string        `json:"client_key,omitempty" mapstructure:"client_key"`
	HashKey       string        `json:"hash_key,omitempty" mapstructure:"hash_key"`
	Timeout       time.Duration `json:"timeout,omitempty" mapstructure:"timeout"`
	Retry         uint          `json:"retry,omitempty" mapstructure:"retry"`
	MaxConnection int32         `json:"max_connection,omitempty" mapstructure:"max_connection"`
	UseTLS        bool          `json:"use_tls,omitempty" mapstructure:"use_tls"`
	DNSResolver   bool          `json:"dns_resolver,omitempty" mapstructure:"dns_resolver"`
}

const (
	defaultPort = 8080
	defaultHost = "0.0.0.0"
)

var (
	//go:embed local.yml
	defaultConfig []byte
	cfg           = &Config{}
)

func GetConfig() *Config {
	return cfg
}

func init() {
	viper.AutomaticEnv()
	configFile := viper.GetString("CONFIG_PATH")

	viper.SetConfigType("yaml")
	viper.AllowEmptyEnv(false)

	if configFile == "" {
		err := viper.ReadConfig(bytes.NewBuffer(defaultConfig))
		if err != nil {
			panic(err)
		}
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigFile(configFile)
		err := viper.ReadInConfig()
		if err != nil {
			panic(err)
		}
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}

	logger := log.New(cfg.Log)

	logger.Info("Config loaded", zap.Reflect("config", cfg))
}

func (l *Listen) Address() string {
	if l == nil {
		return fmt.Sprintf("%s:%d", defaultHost, defaultPort)
	}

	if l.Host == "" {
		l.Host = "0.0.0.0"
	}
	return fmt.Sprintf("%s:%d", l.Host, l.Port)
}
