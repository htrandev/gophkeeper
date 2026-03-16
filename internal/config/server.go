package config

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// ServerConfig хранти настройки сервера.
type ServerConfig struct {
	Addr        string        `mapstructure:"ADDRESS"`
	LogLvl      string        `mapstructure:"LOG_LEVEL"`
	DatabaseDsn string        `mapstructure:"DATABASE_DSN"`
	Signature   string        `mapstructure:"SIGNATURE"`
	MaxRetry    int           `mapstructure:"MAX_RETRY"`
	TokenTTL    time.Duration `mapstructure:"TOKEN_TTL"`
}

// InitServerConfig загружает конфигураяю сервера из разных источников:
// файлов, переданных флагов, переменных окружения.
// Приоритет: переменные окруженя -> флаги -> файл.
func InitServerConfig() (ServerConfig, error) {
	v := viper.New()

	filepath := getConfigFilePath()
	if filepath != "" {
		v.SetConfigFile(filepath)

		if err := v.ReadInConfig(); err != nil {
			return ServerConfig{}, fmt.Errorf("load config file: %w", err)
		}
	}

	flagVals := parseServerFlags(v)

	v.AutomaticEnv()

	for key := range flagVals {
		if envVal, exists := os.LookupEnv(key); exists {
			v.Set(key, envVal)
		}
	}

	var s ServerConfig
	if err := v.Unmarshal(&s); err != nil {
		return ServerConfig{}, fmt.Errorf("unmarshal config: %w", err)
	}

	return s, nil
}

// parseServerFlags парсит переданные серверу флаги.
func parseServerFlags(v *viper.Viper) map[string]any {
	var (
		addr        = pflag.String("a", "localhost:8090", "address to run grpc server")
		logLvl      = pflag.String("lvl", "debug", "log level")
		databaseDsn = pflag.String("d", "", "db dsn")
		signature   = pflag.String("k", "", "secret key")
		maxRetry    = pflag.Int("maxRetry", 3, "pg max retry")
		// privateKeyFile = pflag.String("crypto-key", "", "path to private key file")
		tokenTTL = pflag.Duration("ttl", 1*time.Hour, "token ttl in seconds")
	)
	pflag.Parse()

	// создаем мапу со значениями флагов, для дальнейшего мерджа с переменными окружения
	flagVals := map[string]any{
		"ADDRESS":      *addr,
		"LOG_LEVEL":    *logLvl,
		"DATABASE_DSN": *databaseDsn,
		"SIGNATURE":    *signature,
		"MAX_RETRY":    *maxRetry,
		// "CRYPTO_KEY":   *privateKeyFile,
		"TOKEN_TTL": *tokenTTL,
	}

	for key, val := range flagVals {
		if val != nil {
			v.Set(key, val)
		}
	}

	return flagVals
}
