package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"pet-service/config"
	"pet-service/pkg/logger"
)

var (
	client *redis.Client
)

// Init 初始化Redis连接
func Init(cfg *config.Config) error {
	client = redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
		PoolSize: cfg.Redis.PoolSize,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis连接失败: %w", err)
	}

	logger.Info(context.Background(), "Redis连接成功", logger.String("addr", cfg.Redis.Addr))
	return nil
}

// GetClient 获取Redis客户端
func GetClient() *redis.Client {
	return client
}

// Set 设置缓存
func Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	err := client.Set(ctx, key, value, expiration).Err()
	if err != nil {
		logger.Error(ctx, "Redis设置失败", logger.String("key", key), logger.ErrorField(err))
		return err
	}
	return nil
}

// Get 获取缓存
func Get(ctx context.Context, key string) (string, error) {
	val, err := client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil
		}
		logger.Error(ctx, "Redis获取失败", logger.String("key", key), logger.ErrorField(err))
		return "", err
	}
	return val, nil
}

// Del 删除缓存
func Del(ctx context.Context, keys ...string) error {
	err := client.Del(ctx, keys...).Err()
	if err != nil {
		logger.Error(ctx, "Redis删除失败", logger.Any("keys", keys), logger.ErrorField(err))
		return err
	}
	return nil
}

// Exists 检查key是否存在
func Exists(ctx context.Context, keys ...string) (int64, error) {
	count, err := client.Exists(ctx, keys...).Result()
	if err != nil {
		logger.Error(ctx, "Redis检查key失败", logger.Any("keys", keys), logger.ErrorField(err))
		return 0, err
	}
	return count, nil
}

// Expire 设置过期时间
func Expire(ctx context.Context, key string, expiration time.Duration) error {
	err := client.Expire(ctx, key, expiration).Err()
	if err != nil {
		logger.Error(ctx, "Redis设置过期时间失败", logger.String("key", key), logger.ErrorField(err))
		return err
	}
	return nil
}

// HSet 设置哈希
func HSet(ctx context.Context, key, field string, value interface{}) error {
	err := client.HSet(ctx, key, field, value).Err()
	if err != nil {
		logger.Error(ctx, "Redis HSet失败", logger.String("key", key), logger.ErrorField(err))
		return err
	}
	return nil
}

// HGet 获取哈希
func HGet(ctx context.Context, key, field string) (string, error) {
	val, err := client.HGet(ctx, key, field).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil
		}
		logger.Error(ctx, "Redis HGet失败", logger.String("key", key), logger.ErrorField(err))
		return "", err
	}
	return val, nil
}

// HGetAll 获取所有哈希
func HGetAll(ctx context.Context, key string) (map[string]string, error) {
	val, err := client.HGetAll(ctx, key).Result()
	if err != nil {
		logger.Error(ctx, "Redis HGetAll失败", logger.String("key", key), logger.ErrorField(err))
		return nil, err
	}
	return val, nil
}

// HDel 删除哈希字段
func HDel(ctx context.Context, key string, fields ...string) error {
	err := client.HDel(ctx, key, fields...).Err()
	if err != nil {
		logger.Error(ctx, "Redis HDel失败", logger.String("key", key), logger.ErrorField(err))
		return err
	}
	return nil
}

// Close 关闭Redis连接
func Close() error {
	if client != nil {
		return client.Close()
	}
	return nil
}
