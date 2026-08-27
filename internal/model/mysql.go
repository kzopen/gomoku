package model

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// User 对应 user 表
type User struct {
	ID           int64  `gorm:"column:id;primaryKey"`
	Username     string `gorm:"column:username"`
	PasswordHash string `gorm:"column:password_hash"`
	Nickname     string `gorm:"column:nickname"`
	Avatar       string `gorm:"column:avatar"`
}

func (User) TableName() string { return "user" }

// UserStats 对应 player_stats 表
type UserStats struct {
	UserID   int64 `gorm:"column:user_id;primaryKey"`
	Wins     int32 `gorm:"column:wins"`
	Losses   int32 `gorm:"column:losses"`
	Draws    int32 `gorm:"column:draws"`
	Runaways int32 `gorm:"column:runaways"`
	ELO      int32 `gorm:"column:elo"`
}

func (UserStats) TableName() string { return "player_stats" }

// OpenMySQL 建立 MySQL 连接池（不 Ping；连通性由调用方按需处理）
func OpenMySQL(cfg MySQLConfig) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open mysql: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get mysql pool: %w", err)
	}
	if cfg.MaxOpen > 0 {
		sqlDB.SetMaxOpenConns(cfg.MaxOpen)
	}
	if cfg.MaxIdle > 0 {
		sqlDB.SetMaxIdleConns(cfg.MaxIdle)
	}
	sqlDB.SetConnMaxLifetime(3 * time.Minute)
	return db, nil
}

// Ping 健康检查
func Ping(ctx context.Context, db *gorm.DB) error {
	if db == nil {
		return errors.New("mysql db is nil")
	}
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("get mysql pool: %w", err)
	}
	return sqlDB.PingContext(ctx)
}

// GetUserByUsername 按登录名查询用户
func GetUserByUsername(ctx context.Context, db *gorm.DB, username string) (*User, error) {
	u := &User{}
	if err := db.WithContext(ctx).Where("username = ?", username).First(u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // 未找到
		}
		return nil, fmt.Errorf("query user: %w", err)
	}
	return u, nil
}

// GetUserByID 按 uid 查询用户
func GetUserByID(ctx context.Context, db *gorm.DB, uid int64) (*User, error) {
	u := &User{}
	if err := db.WithContext(ctx).First(u, uid).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("query user by id: %w", err)
	}
	return u, nil
}

// CreateUser 注册新用户（登录即注册）；昵称默认取用户名
func CreateUser(ctx context.Context, db *gorm.DB, username, passwordHash string) (int64, error) {
	u := &User{Username: username, PasswordHash: passwordHash, Nickname: username}
	if err := db.WithContext(ctx).Create(u).Error; err != nil {
		return 0, fmt.Errorf("create user: %w", err)
	}
	return u.ID, nil
}

// CreateUserStats 初始化战绩（幂等）
func CreateUserStats(ctx context.Context, db *gorm.DB, uid int64) error {
	stats := &UserStats{UserID: uid, ELO: 1200}
	if err := db.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(stats).Error; err != nil {
		return fmt.Errorf("create user stats: %w", err)
	}
	return nil
}

// GetUserStats 查询玩家战绩
func GetUserStats(ctx context.Context, db *gorm.DB, uid int64) (*UserStats, error) {
	s := &UserStats{}
	if err := db.WithContext(ctx).Where("user_id = ?", uid).First(s).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &UserStats{UserID: uid, ELO: 1200}, nil // 未建档按初始分
		}
		return nil, fmt.Errorf("query user stats: %w", err)
	}
	return s, nil
}
