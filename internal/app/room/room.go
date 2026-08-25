package room

import (
	"context"
	"database/sql"
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
