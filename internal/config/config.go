package config

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

const (
	formatJSON = "json"
	envFile    = "./internal/config/.env"
)

type Config struct {
	Server struct {
		Host        string `envconfig:"SERVER_HOST" default:":9000"`
		MetricsBind string `envconfig:"BIND_METRICS" default:":9090"`
		HealthHost  string `envconfig:"BIND_HEALTH" default:":9091"`
	}

	Service struct {
		LogLevel  string `envconfig:"LOGGER_LEVEL" default:"debug"`
		LogFormat string `envconfig:"LOGGER_FORMAT" default:"console"`
	}

	DB struct {
		Address  string `envconfig:"DB_ADDRESS"`
		Name     string `envconfig:"DB_NAME"`
		User     string `envconfig:"DB_USER"`
		Password string `envconfig:"DB_PASSWORD"`
		Port     int    `envconfig:"DB_PORT"`
		MaxConn  int    `envconfig:"DB_MAX_CONN"`
	}
}

func Parse() (*Config, error) {
	var cfg = new(Config)

	// Загружаем в переменные окружения из .env
	err := godotenv.Load(envFile)
	if err != nil {
		return nil, errors.Wrap(err, "error loading .env file")
	}

	// Загружаем в envconfig
	err = envconfig.Process("", cfg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to process env vars")
	}

	return cfg, nil
}

func (cfg Config) Logger() (logger zerolog.Logger) {
	level := zerolog.InfoLevel
	if newLevel, err := zerolog.ParseLevel(cfg.Service.LogLevel); err == nil {
		level = newLevel
	}

	var out io.Writer = os.Stdout
	if cfg.Service.LogFormat != formatJSON {
		out = zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.StampMicro}
	}
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	return zerolog.New(out).Level(level).With().Caller().Timestamp().Logger()
}

// Получаем адрес в БД
func (cfg Config) GetDBConnString() string {
	return fmt.Sprintf(
		"host=%s port=%d dbname=%s sslmode=disable user=%s password=%s",
		cfg.DB.Address, cfg.DB.Port, cfg.DB.Name, cfg.DB.User, cfg.DB.Password,
	)
}

func (cfg Config) PgPoolConfig() (*pgxpool.Config, error) {
	poolCfg, err := pgxpool.ParseConfig(fmt.Sprintf("%s pool_max_conns=%d", cfg.GetDBConnString(), cfg.DB.MaxConn))
	if err != nil {
		return nil, err
	}

	return poolCfg, nil
}
