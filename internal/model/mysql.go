package model

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// User 对应 user 表
type User struct {
	ID           int64
	Username     string
	PasswordHash string
	Nickname     string
	Avatar       string
}

// UserStats 对应 player_stats 表
type UserStats struct {
	Wins     int32
	Losses   int32
	Draws    int32
	Runaways int32
	ELO      int32
}

// OpenMySQL 建立 MySQL 连接池（不 Ping；连通性由调用方按需处理）
func OpenMySQL(cfg MySQLConfig) (*sql.DB, error) {
	db, err := sql.Open("mysql", cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("open mysql: %w", err)
	}
	db.SetMaxOpenConns(cfg.MaxOpen)
	db.SetMaxIdleConns(cfg.MaxIdle)
	db.SetConnMaxLifetime(3 * time.Minute)
	return db, nil
}

// Ping 健康检查
func Ping(ctx context.Context, db *sql.DB) error {
	return db.PingContext(ctx)
}

// GetUserByUsername 按登录名查询用户
func GetUserByUsername(ctx context.Context, db *sql.DB, username string) (*User, error) {
	row := db.QueryRowContext(ctx,
		`SELECT id, username, password_hash, nickname, avatar FROM user WHERE username = ?`, username)
	u := &User{}
	if err := row.Scan(&u.ID, &u.Username, &u.PasswordHash, &u.Nickname, &u.Avatar); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // 未找到
		}
		return nil, fmt.Errorf("query user: %w", err)
	}
	return u, nil
}

// GetUserByID 按 uid 查询用户
func GetUserByID(ctx context.Context, db *sql.DB, uid int64) (*User, error) {
	row := db.QueryRowContext(ctx,
		`SELECT id, username, password_hash, nickname, avatar FROM user WHERE id = ?`, uid)
	u := &User{}
	if err := row.Scan(&u.ID, &u.Username, &u.PasswordHash, &u.Nickname, &u.Avatar); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("query user by id: %w", err)
	}
	return u, nil
}

// CreateUser 注册新用户（登录即注册）；昵称默认取用户名
func CreateUser(ctx context.Context, db *sql.DB, username, passwordHash string) (int64, error) {
	res, err := db.ExecContext(ctx,
		`INSERT INTO user (username, password_hash, nickname) VALUES (?, ?, ?)`,
		username, passwordHash, username)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// CreateUserStats 初始化战绩（幂等）
func CreateUserStats(ctx context.Context, db *sql.DB, uid int64) error {
	_, err := db.ExecContext(ctx,
		`INSERT IGNORE INTO player_stats (user_id, wins, losses, draws, runaways, elo) VALUES (?, 0, 0, 0, 0, 1200)`, uid)
	return err
}

// GetUserStats 查询玩家战绩
func GetUserStats(ctx context.Context, db *sql.DB, uid int64) (*UserStats, error) {
	row := db.QueryRowContext(ctx,
		`SELECT wins, losses, draws, runaways, elo FROM player_stats WHERE user_id = ?`, uid)
	s := &UserStats{}
	if err := row.Scan(&s.Wins, &s.Losses, &s.Draws, &s.Runaways, &s.ELO); err != nil {
		if err == sql.ErrNoRows {
			return &UserStats{ELO: 1200}, nil // 未建档按初始分
		}
		return nil, fmt.Errorf("query user stats: %w", err)
	}
	return s, nil
}
