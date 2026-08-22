package model

import (
	"time"
)

// GameConfig 对应 config.yaml 的 game 段
type GameConfig struct {
	BoardSize               int  `mapstructure:"board_size"`
	StepTimeLimit           int  `mapstructure:"step_time_limit"`
	TotalTimeLimit          int  `mapstructure:"total_time_limit"`
	TimerPushInterval       int  `mapstructure:"timer_push_interval"`
	HeartbeatTimeout        int  `mapstructure:"heartbeat_timeout"`
	AITakeoverGrace         int  `mapstructure:"ai_takeover_grace"`
	AITakeoverReturnControl bool `mapstructure:"ai_takeover_return_control"`
	AILevel                 int  `mapstructure:"ai_level"`
	ForbidRules             bool `mapstructure:"forbid_rules"`
}

// MySQLConfig 对应 config.yaml 的 mysql 段
type MySQLConfig struct {
	DSN     string `mapstructure:"dsn"`
	MaxOpen int    `mapstructure:"max_open"`
	MaxIdle int    `mapstructure:"max_idle"`
}

// RedisConfig 对应 config.yaml 的 redis 段
type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// AppConfig 对应 config.yaml 的 app 段
type AppConfig struct {
	Name       string `mapstructure:"name"`
	ServerType string `mapstructure:"server_type"`
	Frontend   bool   `mapstructure:"frontend"`
	Port       int    `mapstructure:"port"`
	GRPCPort   int    `mapstructure:"grpc_port"`
	Cluster    bool   `mapstructure:"cluster"`
	LogLevel   string `mapstructure:"log_level"`
}

// Config 为服务端全部配置
type Config struct {
	App   AppConfig   `mapstructure:"app"`
	MySQL MySQLConfig `mapstructure:"mysql"`
	Redis RedisConfig `mapstructure:"redis"`
	Game  GameConfig  `mapstructure:"game"`
}

// Redis key 模板（与《开发文档》§8.2 一致）
const (
	KeyOnline    = "online:%d"     // string；在线状态（1），TTL 60s 心跳续期
	KeySession   = "session:%s"    // string；登录态（uid、昵称），TTL 7d
	KeyRoom      = "room:%d"       // hash；房间摘要
	KeyRoomMoves = "room:%d:moves" // list；走子序列
	KeyMatchQ    = "match:queue"   // list；匹配等待队列
	KeyRankElo   = "rank:elo"      // zset；member=uid, score=elo
	KeyStatUser  = "stat:user:%d"  // hash；战绩缓存，TTL 1h
)

// SessionTTL / OnlineTTL / StatTTL 单位秒
const (
	SessionTTL = 7 * 24 * 3600
	OnlineTTL  = 60
	StatTTL    = 3600
)

// UnixMS 返回毫秒时间戳（对局时间字段统一用 ms）
func UnixMS(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}
