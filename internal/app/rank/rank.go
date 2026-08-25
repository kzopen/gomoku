// Package rank 实现排行榜/在线统计组件（在线人数、心跳续期）。
// 客户端路由：rank.online_count / rank.ping。
package rank

import (
	"context"
	"strconv"

	"github.com/redis/go-redis/v9"
	"github.com/topfreegames/pitaya/v2"
	"github.com/topfreegames/pitaya/v2/component"

	"gomoku/internal/common"
	"gomoku/internal/model"
)

// Component 为排行榜/在线组件
type Component struct {
	component.Base
	app pitaya.Pitaya
	rdb *redis.Client
}

// New 创建 rank 组件
func New(app pitaya.Pitaya, rdb *redis.Client) *Component {
	return &Component{app: app, rdb: rdb}
}

// OnlineCount 处理 rank.online_count：返回当前在线人数
func (c *Component) OnlineCount(ctx context.Context, _ *model.C2SOnlineCount) (*model.S2COnlineCount, error) {
	resp := &model.S2COnlineCount{}
	if c.rdb == nil {
		resp.Code = common.CodeInternalError
		resp.Msg = "在线服务不可用"
		return resp, nil
	}
	n, err := model.OnlineCount(ctx, c.rdb)
	if err != nil {
		resp.Code = common.CodeInternalError
		resp.Msg = "查询失败"
		return resp, nil
	}
	resp.Code = common.CodeOK
	resp.Count = n
	return resp, nil
}

// Ping 处理 rank.ping：心跳续期（前端每 30s notify 一次，刷新 online:{uid} 的 TTL）
func (c *Component) Ping(ctx context.Context, _ *model.C2SPing) (*model.S2CPing, error) {
	resp := &model.S2CPing{Code: common.CodeOK}
	if c.rdb == nil {
		return resp, nil
	}
	s := c.app.GetSessionFromCtx(ctx)
	if s == nil || s.UID() == "" {
		return resp, nil // 未登录不续期
	}
	uid, err := strconv.ParseInt(s.UID(), 10, 64)
	if err != nil {
		return resp, nil
	}
	_ = model.SetOnline(ctx, c.rdb, uid) // TTL 60s 续期
	return resp, nil
}
