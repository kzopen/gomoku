package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/topfreegames/pitaya/v2"
	"github.com/topfreegames/pitaya/v2/component"
	"golang.org/x/crypto/bcrypt"
	"gomoku/internal/common"
	"gomoku/internal/model"
	"gorm.io/gorm"
)

type Component struct {
	component.Base
	app pitaya.Pitaya
	db  *gorm.DB
	rdb *redis.Client
}

func New(app pitaya.Pitaya, db *gorm.DB, rdb *redis.Client) *Component {
	return &Component{
		app: app,
		db:  db,
		rdb: rdb,
	}
}
func (a *Component) Login(ctx context.Context, in *model.C2SLogin) (*model.S2CLogin, error) {
	resp := &model.S2CLogin{}
	if in == nil || in.Username == "" || in.Password == "" {
		return common.ErrorResponse(resp, common.CodeBadParam, "用户名或密码不能为空"), nil
	}
	if a.db == nil {
		return common.ErrorResponse(resp, common.CodeInternalError, "数据库未连接"), nil
	}
	u, err := model.GetUserByUsername(ctx, a.db, in.Username)
	if err != nil {
		return common.ErrorResponse(resp, common.CodeInternalError, "服务器内部错误"), nil
	}
	//不存在用户
	if u == nil {
		// 登录即注册：账号不存在则自动建档（昵称取用户名）
		pwdHash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
		if err != nil {
			return common.ErrorResponse(resp, common.CodeInternalError, "服务器内部错误"), nil
		}
		newUserId, err := model.CreateUser(ctx, a.db, in.Username, string(pwdHash))
		//重复点击注册 直接登录
		if err != nil {
			// 并发注册撞唯一键：回查后按既有账号继续
			if u, err = model.GetUserByUsername(ctx, a.db, in.Username); err != nil || u == nil {
				return common.ErrorResponse(resp, common.CodeInternalError, "注册失败"), nil
			}
			if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(in.Password)) != nil {
				return common.ErrorResponse(resp, common.CodeBadParam, "用户名或密码错误"), nil
			}
		} else {
			//新用户
			_ = model.CreateUserStats(ctx, a.db, newUserId)
			u = &model.User{ID: newUserId, Username: in.Username, Nickname: in.Username}
		}

	} else if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(in.Password)) != nil {
		return common.ErrorResponse(resp, common.CodeBadParam, "用户名或密码错误"), nil
	}
	token := newToken()
	if a.rdb != nil {
		//session:$token=>{uid,nickname}
		if err := model.SetSession(ctx, a.rdb, token, model.SessionData{UID: u.ID, Nickname: u.Nickname}); err != nil {
			return common.ErrorResponse(resp, common.CodeInternalError, "登录态写入失败"), nil
		}
		//online:$uid=>1
		if err := model.SetOnline(ctx, a.rdb, u.ID); err != nil {
			return common.ErrorResponse(resp, common.CodeInternalError, "在线状态写入失败"), nil
		}
	}
	if s := a.app.GetSessionFromCtx(ctx); s != nil {
		uidStr := fmt.Sprintf("%d", u.ID)
		switch {
		case s.UID() == "":
			if err := s.Bind(ctx, uidStr); err != nil {
				return common.ErrorResponse(resp, common.CodeInternalError, "会话绑定失败"), nil
			}
		case s.UID() != uidStr:
			// 同一连接已绑定其他账号（换号登录），Pitaya 不允许改绑
			return common.ErrorResponse(resp, common.CodeBadParam, "该连接已登录其他账号，请刷新页面"), nil
		}
	}
	stats, err := model.GetUserStats(ctx, a.db, u.ID)
	if err != nil {
		return common.ErrorResponse(resp, common.CodeInternalError, "服务器内部错误"), nil
	}

	resp.UID = u.ID
	resp.Token = token
	resp.Nickname = u.Nickname
	resp.Avatar = u.Avatar
	resp.ELO = stats.ELO
	resp.Wins = stats.Wins
	resp.Losses = stats.Losses
	resp.Draws = stats.Draws
	return common.SuccessResponse(resp, ""), nil

}
func newToken() string {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
