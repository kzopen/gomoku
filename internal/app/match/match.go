package match

import (
	"context"
	"database/sql"
	"github.com/redis/go-redis/v9"
	"github.com/topfreegames/pitaya/v2"
	"github.com/topfreegames/pitaya/v2/component"
	"gomoku/internal/common"
	"gomoku/internal/model"
	"gomoku/internal/service"
	"strconv"
	"time"
)

type Component struct {
	component.Base
	app pitaya.Pitaya
	rdb *redis.Client
	mgr *service.RoomManager
	db  *sql.DB
}

// New 创建 match 组件
func New(app pitaya.Pitaya, rdb *redis.Client, mgr *service.RoomManager) *Component {
	return &Component{
		app: app,
		rdb: rdb,
		mgr: mgr,
	}
}
func (c *Component) Join(ctx context.Context, _ *model.C2SJoinMatch) (*model.S2CJoinMatch, error) {
	resp := &model.S2CJoinMatch{}
	uid, ok := c.uidFromCtx(ctx)
	if !ok {
		resp.Code = common.CodeNotLoggedIn
		resp.Msg = "未登录"
		return resp, nil
	}
	if c.rdb == nil {
		resp.Code = common.CodeInternalError
		resp.Msg = "匹配服务不可用"
		return resp, nil
	}
	uidInt, _ := strconv.ParseInt(uid, 10, 64)
	// 已在队列中
	if c.inQueue(ctx, uidInt) {
		resp.Code = common.CodeAlreadyIn
		resp.Msg = "已在匹配中"
		return resp, nil
	}
	if err := model.MatchEnqueue(ctx, c.rdb, uidInt); err != nil {
		resp.Code = common.CodeMatchFailed
		resp.Msg = "匹配队列已满或撮合失败"
		return resp, nil
	}
	resp.Code = common.CodeOK
	resp.Msg = "已进入匹配队列"
	return resp, nil
}

func (c *Component) uidFromCtx(ctx context.Context) (string, bool) {
	s := c.app.GetSessionFromCtx(ctx)
	if s == nil || s.UID() == "" {
		return "", false
	}
	return s.UID(), true
}
func (c *Component) inQueue(ctx context.Context, id int64) bool {
	queue, err := c.rdb.LRange(ctx, model.KeyMatchQ, 0, -1).Result()
	if err != nil {
		return false
	}
	for _, v := range queue {
		if v == strconv.FormatInt(id, 10) {
			return true
		}
	}
	return false
}

func (c *Component) Init() {
	go c.matchLoop()
}
func (c *Component) matchLoop() {
	t := time.NewTicker(500 * time.Millisecond)
	defer t.Stop()
	for range t.C {
		c.tryMatch()
	}
}
func (c *Component) tryMatch() {
	if c.rdb == nil || c.mgr == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	len, _ := c.rdb.LLen(ctx, model.KeyMatchQ).Result()
	if len >= 2 {
		a, err := model.MatchDequeue(ctx, c.rdb)
		if err != nil {
			return // 队列空或不可用
		}
		b, err := model.MatchDequeue(ctx, c.rdb)
		if err != nil {
			err := model.MatchEnqueue(ctx, c.rdb, a)
			if err != nil {
				return
			}
			return // 队列空或不可用
		}
		if _, err := c.mgr.CreateGameRoom(a, b); err != nil {
			// 创建房间失败：放回队列，避免丢玩家
			_ = model.MatchEnqueue(ctx, c.rdb, a)
			_ = model.MatchEnqueue(ctx, c.rdb, b)
		}
	}
}
