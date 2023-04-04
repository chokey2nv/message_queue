package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/chokey2nv/obiex.finance/database"
)

type Config struct {
	DBClient  *database.DbClient
	RMQClient *RabbitClient
}
type RabbitClient struct {
	Host string
	Port int
	User string
	Pass string
}

func Load() (*Config, error) {
	cfg := Config{}
	cfg.DBClient = database.NewDBClient(
		getEnv("MYSQL_HOST", "localhost"), //127.0.0.1
		getEnvInt("MYSQL_PORT", 3306),
		getEnv("MYSQL_USER", "root"),
		getEnv("MYSQL_PASSWORD", ""),
		getEnv("MYSQL_DATABASE", "obiex_finance"),
	)
	cfg.RMQClient = NewRabbitMQClient(
		getEnv("RABBITMQ_HOST", "localhost"),
		getEnvInt("RABBITMQ_PORT", 5672),
		getEnv("RABBITMQ_USER", "guest"),
		getEnv("RABBITMQ_PASSWORD", "guest"),
	)
	return &cfg, nil
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func getEnvInt(key string, fallback int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return fallback
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		fmt.Printf("Failed to parse %s: %v\n", key, err)
		return fallback
	}
	return value
}
func NewRabbitMQClient(host string, port int, user string, pass string) *RabbitClient {
	return &RabbitClient{
		Host: host, Port: port, User: user, Pass: pass,
	}
}
func (c *RabbitClient) RabbitMQURL() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d/",
		c.User,
		c.Pass,
		c.Host,
		c.Port)
}
