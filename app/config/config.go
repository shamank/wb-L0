package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"os"
	"time"
)

type (
	Config struct {
		HTTP     HTTPConfig
		Postgres PostgresConfig
		Nats     NatsConfig
		Cache    CacheConfig
	}

	HTTPConfig struct {
		Host               string
		Port               string
		ReadTimeout        time.Duration
		WriteTimeout       time.Duration
		MaxHeaderMegabytes int
	}

	PostgresConfig struct {
		Host     string
		Port     string
		User     string
		Password string
		DBName   string
		SSLMode  string
	}

	NatsConfig struct {
		URL       string
		ClusterID string
		ClientID  string
		Channel   string
	}

	CacheConfig struct {
		TTL     time.Duration
		CleanUp time.Duration
	}
)

func ConfigInit(configDir string) (*Config, error) {
	if err := parseConfig(configDir); err != nil {
		return nil, err
	}
	readTimeoutDuration, err := time.ParseDuration(viper.GetString("http.readTimeout"))
	if err != nil {
		return nil, err
	}

	writeTimeoutDuration, err := time.ParseDuration(viper.GetString("http.writeTimeout"))
	if err != nil {
		return nil, err
	}

	ttl, err := time.ParseDuration(viper.GetString("cache.ttl"))
	if err != nil {
		return nil, err
	}

	cleanUp, err := time.ParseDuration(viper.GetString("cache.cleanUp"))
	if err != nil {
		return nil, err
	}

	return &Config{
		HTTP: HTTPConfig{
			Host:               viper.GetString("http.host"),
			Port:               viper.GetString("http.port"),
			ReadTimeout:        readTimeoutDuration,
			WriteTimeout:       writeTimeoutDuration,
			MaxHeaderMegabytes: viper.GetInt("http.maxHeaderBytes"),
		},
		Postgres: PostgresConfig{
			Host:     viper.GetString("pg.host"),
			Port:     viper.GetString("pg.port"),
			User:     viper.GetString("pg.user"),
			Password: os.Getenv("DB_PASSWORD"),
			DBName:   viper.GetString("pg.dbname"),
			SSLMode:  viper.GetString("pg.sslmode"),
		},
		Nats: NatsConfig{
			URL:       viper.GetString("nats.url"),
			ClusterID: viper.GetString("nats.clusterID"),
			ClientID:  viper.GetString("nats.clientID"),
			Channel:   viper.GetString("nats.channel"),
		},
		Cache: CacheConfig{
			TTL:     ttl,
			CleanUp: cleanUp,
		},
	}, nil
}

func parseConfig(configDir string) error {
	viper.AddConfigPath(configDir)
	viper.SetConfigName("main")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if err := godotenv.Load(); err != nil {
		return err
	}

	return viper.MergeInConfig()
}
