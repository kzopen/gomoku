package room

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/topfreegames/pitaya/v2"
	"github.com/topfreegames/pitaya/v2/component"
	"gomoku/internal/common"
	"gomoku/internal/model"
	"gomoku/internal/service"
)

type Component struct {
	component.Base
	app pitaya.Pitaya
	rdb *redis.Client
	mgr *service.RoomManager
	db  *sql.DB
}

func New(app pitaya.Pitaya, db *sql.DB, rdb *redis.Client, mgr *service.RoomManager) *Component {
	return &Component{
		app: app,
		db:  db,
		rdb: rdb,
		mgr: mgr,
	}
}

// Reconnect 处理 room.reconnect（断线重连恢复）：
// token 校验 → 定位房间 → 校验玩家在房 → 新连接绑定 uid/room_id → 返回全量状态 → 通知对方。
func (c *Component) Reconnect(ctx context.Context, in *model.C2SReconnect) (*model.S2CReconnect, error) {
	bad := func(code int32, msg string) (*model.S2CReconnect, error) {
		return &model.S2CReconnect{Code: code, Msg: msg}, nil
	}
	if in == nil || in.Token == "" || in.RoomID <= 0 {
		return bad(common.CodeBadParam, "参数错误")
	}
	if c.rdb == nil {
		return bad(common.CodeInternalError, "会话服务不可用")
	}
	if c.mgr == nil {
		return bad(common.CodeInternalError, "房间服务不可用")
	}
	// 1) token 校验（Redis session:{token}）
	sd, err := model.GetSession(ctx, c.rdb, in.Token)
	if err != nil || sd == nil {
		return bad(common.CodeNotLoggedIn, "会话失效")
	}
	// 2) 定位房间
	r := c.mgr.Get(in.RoomID)
	if r == nil {
		return bad(common.CodeRoomNotFound, "房间不存在或已销毁")
	}
	// 3) 校验玩家在本房间并构造全量状态
	uid := fmt.Sprintf("%d", sd.UID)
	resp := r.Snapshot(uid)
	if resp == nil {
		return bad(common.CodeRoomNotFound, "不在该房间")
	}
	// 4) 新连接绑定 uid + room_id（此后 room.place 经新 session 定位房间）
	s := c.app.GetSessionFromCtx(ctx)
	if s != nil {
		if err := s.Bind(ctx, uid); err != nil {
			// 旧 session 可能尚未清理导致绑定冲突；忽略，房间定位已由 room_id 保证
		}
		_ = s.Set("room_id", in.RoomID)
	}
	// 5) 通知对方玩家已回来
	r.NotifyBack(resp.Seat)
	return resp, nil
}

func (c *Component) Place(ctx context.Context, in *model.C2SPlacePiece) (*model.S2CPlacePiece, error) {
	s := c.app.GetSessionFromCtx(ctx)

	if s == nil || s.UID() == "" {
		return &model.S2CPlacePiece{Code: common.CodeNotLoggedIn, Msg: "未登录"}, nil
	}
	roomID, ok := s.Get("room_id").(int64)
	if !ok || roomID <= 0 {
		return &model.S2CPlacePiece{Code: common.CodeRoomNotFound, Msg: "房间不存在或已销毁"}, nil
	}
	if c.mgr == nil {
		return &model.S2CPlacePiece{Code: common.CodeInternalError, Msg: "房间服务不可用"}, nil
	}
	room := c.mgr.Get(roomID)
	if room == nil {
		return &model.S2CPlacePiece{Code: common.CodeRoomNotFound, Msg: "房间不存在或已销毁"}, nil
	}
	if in == nil {
		return &model.S2CPlacePiece{Code: common.CodeBadParam, Msg: "参数错误"}, nil
	}
	return room.Place(s.UID(), in.X, in.Y), nil
}
