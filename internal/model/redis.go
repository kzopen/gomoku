package model

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

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

// OnlineCount 统计在线人数：60s 窗口内发过心跳的用户数（ZSET 时间戳，O(log N) 不阻塞）。
// 先惰性清理过期成员（ZREMRANGEBYSCORE），再 ZCOUNT 统计，无需单独 TTL 和定时任务。
func OnlineCount(ctx context.Context, c *redis.Client) (int64, error) {
	if c == nil {
		return 0, nil
	}
	now := time.Now().Unix()
	min := strconv.FormatInt(now-OnlineTTL, 10)
	_ = c.ZRemRangeByScore(ctx, KeyOnlineZSet, "-inf", min).Err() // 踢出 60s 前的心跳
	return c.ZCount(ctx, KeyOnlineZSet, min, "+inf").Result()
}

// Heartbeat 心跳续期：刷新 online:{uid} 的 TTL（快速判在线）+ 更新 ZSET 时间戳（统计人数）
func Heartbeat(ctx context.Context, c *redis.Client, uid int64) error {
	if c == nil {
		return nil
	}
	if err := SetOnline(ctx, c, uid); err != nil {
		return err
	}
	z := &redis.Z{Score: float64(time.Now().Unix()), Member: uid}
	return c.ZAdd(ctx, KeyOnlineZSet, *z).Err()
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
