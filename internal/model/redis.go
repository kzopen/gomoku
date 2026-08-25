package model

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

// OpenRedis 建立 Redis 客户端（不 Ping；连通性由调用方按需处理）
func OpenRedis(cfg RedisConfig) (*redis.Client, error) {
	c := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	return c, nil
}

// ==================== 登录态 ====================

// SessionData 为 session:{token} 的 value
type SessionData struct {
	UID      int64  `json:"uid"`
	Nickname string `json:"nickname"`
}

// SetSession 写入登录态，TTL 7d
func SetSession(ctx context.Context, c *redis.Client, token string, sd SessionData) error {
	b, _ := json.Marshal(sd)
	return c.Set(ctx, fmt.Sprintf(KeySession, token), b, SessionTTL*time.Second).Err()
}

// GetSession 读取登录态；不存在返回 (nil, nil)
func GetSession(ctx context.Context, c *redis.Client, token string) (*SessionData, error) {
	b, err := c.Get(ctx, fmt.Sprintf(KeySession, token)).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	sd := &SessionData{}
	if err := json.Unmarshal(b, sd); err != nil {
		return nil, err
	}
	return sd, nil
}

// DelSession 删除登录态（登出/踢线）
func DelSession(ctx context.Context, c *redis.Client, token string) error {
	return c.Del(ctx, fmt.Sprintf(KeySession, token)).Err()
}

// ==================== 在线状态 ====================

// SetOnline 标记在线，TTL 60s（心跳续期）
func SetOnline(ctx context.Context, c *redis.Client, uid int64) error {
	return c.Set(ctx, fmt.Sprintf(KeyOnline, uid), "1", OnlineTTL*time.Second).Err()
}

// IsOnline 判断在线
func IsOnline(ctx context.Context, c *redis.Client, uid int64) (bool, error) {
	n, err := c.Exists(ctx, fmt.Sprintf(KeyOnline, uid)).Result()
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

// OnlineCount 统计在线人数：online:* key 的数量（TTL 未过期即在线）。
// 教学版用 Keys 简单直观；生产量大时应改用 SCAN 游标遍历避免阻塞 Redis。
func OnlineCount(ctx context.Context, c *redis.Client) (int64, error) {
	if c == nil {
		return 0, nil
	}
	keys, err := c.Keys(ctx, "online:*").Result()
	if err != nil {
		return 0, err
	}
	return int64(len(keys)), nil
}

// ==================== 匹配队列 ====================

// MatchEnqueue 入队（LPUSH）
func MatchEnqueue(ctx context.Context, c *redis.Client, uid int64) error {
	return c.LPush(ctx, KeyMatchQ, uid).Err()
}

// MatchDequeue 出队（RPOP，与 LPUSH 配对取最早入队者）
func MatchDequeue(ctx context.Context, c *redis.Client) (int64, error) {
	return c.RPop(ctx, KeyMatchQ).Int64()
}

// MatchRemove 从队列移除（取消匹配，LREM）
func MatchRemove(ctx context.Context, c *redis.Client, uid int64) error {
	return c.LRem(ctx, KeyMatchQ, 0, uid).Err()
}
