package config

import (
	"os"
	"strconv"
	"time"
)

// Config 应用配置
type Config struct {
	Server   ServerConfig
	Redis    RedisConfig
	Database DatabaseConfig
	Log      LogConfig
	JWT      JWTConfig
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret        string
	TokenDuration int
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// RedisConfig Redis配置
type RedisConfig struct {
	Addr     string
	Password string
	DB       int
	PoolSize int
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	DSN             string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
	OutputPath string
}

// Load 加载配置
func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Addr:         getEnv("SERVER_ADDR", ":8888"),
			ReadTimeout:  60 * time.Second,
			WriteTimeout: 60 * time.Second,
			IdleTimeout:  120 * time.Second,
		},
		Redis: RedisConfig{
			Addr:     getEnv("REDIS_ADDR", "localhost:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
			PoolSize: getEnvInt("REDIS_POOL_SIZE", 10),
		},
		Database: DatabaseConfig{
			DSN:             getEnv("DATABASE_DSN", "root:123456@tcp(localhost:3306)/pet_service?charset=utf8mb4&parseTime=True&loc=Local"),
			MaxIdleConns:    getEnvInt("DB_MAX_IDLE_CONNS", 10),
			MaxOpenConns:    getEnvInt("DB_MAX_OPEN_CONNS", 100),
			ConnMaxLifetime: 1 * time.Hour,
		},
		Log: LogConfig{
			Level:      getEnv("LOG_LEVEL", "info"),
			MaxSize:    getEnvInt("LOG_MAX_SIZE", 100),
			MaxBackups: getEnvInt("LOG_MAX_BACKUPS", 3),
			MaxAge:     getEnvInt("LOG_MAX_AGE", 28),
			Compress:   getEnvBool("LOG_COMPRESS", true),
			OutputPath: getEnv("LOG_OUTPUT_PATH", "./logs/app.log"),
		},
		JWT: JWTConfig{
			Secret:        getEnv("JWT_SECRET", "pet-service-secret-key-2024"),
			TokenDuration: getEnvInt("JWT_TOKEN_DURATION", 24),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return defaultValue
}
