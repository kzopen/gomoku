package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/topfreegames/pitaya/v2"
	"github.com/topfreegames/pitaya/v2/acceptor"
	"github.com/topfreegames/pitaya/v2/component"
	"github.com/topfreegames/pitaya/v2/config"
	logruswrapper "github.com/topfreegames/pitaya/v2/logger/logrus"
	"github.com/topfreegames/pitaya/v2/serialize/json"
	"gomoku/internal/app/auth"
	"gomoku/internal/app/match"
	"gomoku/internal/app/room"
	"gomoku/internal/model"
	"gomoku/internal/service"
	"os"
	"strings"
	"time"
)

func main() {
	var cfgPath string
	flag.StringVar(&cfgPath, "c", "config/config.yaml", "配置文件路径") //设置配置文件地址
	flag.Parse()
	cfg, v := loadConfig(cfgPath)
	//日志
	l := logrus.New()
	l.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	lvl, err := logrus.ParseLevel(cfg.App.LogLevel)
	if err != nil {
		lvl = logrus.InfoLevel
	}
	l.SetLevel(lvl)
	pitaya.SetLogger(logruswrapper.NewWithFieldLogger(l))
	//db
	// 数据层（连接失败不阻塞启动，组件内按 nil 降级返回错误码）
	db := openMySQL(cfg, l)
	rdb := openRedis(cfg, l)

	// Pitaya builder：单机 Standalone / 集群 Cluster
	mode := pitaya.Standalone
	if cfg.App.Cluster {
		mode = pitaya.Cluster
	}
	pitayaConfig := config.NewConfig(v)
	builder := pitaya.NewBuilderWithConfigs(cfg.App.Frontend, cfg.App.ServerType, mode, map[string]string{}, pitayaConfig)
	builder.Serializer = json.NewSerializer() // JSON 序列化（《开发文档》§7.5 路径 A）
	// WebSocket 接入
	builder.AddAcceptor(acceptor.NewWSAcceptor(fmt.Sprintf(":%d", cfg.App.Port)))

	app := builder.Build()

	// 全局只有一个房间管理器（对局核心：撮合成功后由 match 调用创建房间，room.place 定位对局）
	roomMgr := service.NewRoomManager(app, builder.SessionPool, db, cfg.Game)

	// 注册组件（组件名与方法名统一小写，与《开发文档》§6.2/§7.3 路由一致）
	register := func(c component.Component, name string) {
		app.Register(c, component.WithName(name), component.WithNameFunc(strings.ToLower))
	}
	register(auth.New(app, db, rdb), "auth")
	register(match.New(app, rdb, roomMgr), "match")
	register(room.New(app, db, rdb, roomMgr), "room")
	//register(rank.New(app), "rank")

	l.Infof("starting gomoku server: type=%s frontend=%v port=%d mode=%v",
		cfg.App.ServerType, cfg.App.Frontend, cfg.App.Port, mode)
	app.Start()
}

// openMySQL 打开 MySQL；失败仅告警（auth 组件按 nil 降级）
func openMySQL(cfg *model.Config, l *logrus.Logger) *sql.DB {
	db, err := model.OpenMySQL(cfg.MySQL)
	if err != nil {
		l.Warnf("MySQL 连接失败（登录功能将不可用）: %v", err)
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		l.Warnf("MySQL Ping 失败（登录功能将不可用）: %v", err)
		return nil
	}
	return db
}

// openRedis 打开 Redis；失败仅告警（match/room 组件按 nil 降级）
func openRedis(cfg *model.Config, l *logrus.Logger) *redis.Client {
	c, err := model.OpenRedis(cfg.Redis)
	if err != nil {
		l.Warnf("Redis 连接失败（匹配/在线功能将不可用）: %v", err)
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := c.Ping(ctx).Err(); err != nil {
		l.Warnf("Redis Ping 失败（匹配/在线功能将不可用）: %v", err)
		return nil
	}
	return c
}

func loadConfig(path string) (*model.Config, *viper.Viper) {
	v := viper.New()
	v.SetConfigFile(path)
	if err := v.ReadInConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "读取配置失败: %v\n", err)
		os.Exit(1)
	}
	cfg := &model.Config{}
	if err := v.Unmarshal(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "解析配置失败: %v\n", err)
		os.Exit(1)
	}
	if cfg.App.ServerType == "" {
		cfg.App.ServerType = "room"
	}
	if cfg.App.Port == 0 {
		cfg.App.Port = 3250
	}
	return cfg, v

}
