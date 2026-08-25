package service

import (
	"context"
	"database/sql"
	"github.com/topfreegames/pitaya/v2"
	"github.com/topfreegames/pitaya/v2/session"
	"gomoku/internal/common"
	"gomoku/internal/model"
	"strconv"
	"sync"
	"time"
)

type RoomManager struct {
	mu        sync.Mutex
	rooms     map[int64]*Room
	nextID    int64
	app       pitaya.Pitaya
	pool      session.SessionPool
	db        *sql.DB
	game      model.GameConfig
	frontType string
}

func (m *RoomManager) nickname(ctx context.Context, uid int64) string {
	if m.db == nil {
		return ""
	}
	u, err := model.GetUserByID(ctx, m.db, uid)
	if err != nil || u == nil {
		return ""
	}
	return u.Nickname
}

// 创建游戏房间
func (m *RoomManager) CreateGameRoom(blackUID, whiteUID int64) (int64, error) {
	m.mu.Lock()
	id := m.nextID
	m.nextID++
	m.mu.Unlock()

	r := &Room{
		mgr:       m,
		ID:        id,
		Seat:      [2]int64{blackUID, whiteUID},
		TurnSeat:  common.SeatBlack,
		State:     common.RoomPlaying,
		StepLimit: int32(m.game.StepTimeLimit),  //回合限时
		TotalTime: int32(m.game.TotalTimeLimit), //整体限时
	}
	if r.StepLimit == 0 {
		r.StepLimit = common.DefaultStepTimeLimit
	}
	if r.TotalTime == 0 {
		r.TotalTime = common.DefaultTotalTimeLimit
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	r.blackNick = m.nickname(ctx, blackUID)
	r.whiteNick = m.nickname(ctx, whiteUID)

	m.mu.Lock()
	m.rooms[id] = r
	m.mu.Unlock()
	//写入双方session
	m.setRoomID(blackUID, id) //uid->房间号关联了
	m.setRoomID(whiteUID, id)
	//推送消息
	seats := []struct {
		uid  int64
		seat int32
	}{
		{blackUID, 0},
		{whiteUID, 1},
	}
	for _, p := range seats {
		m.pushUsers("room.onMatchSuccess", &model.S2CMatchSuccess{
			RoomID:         id,
			Seat:           p.seat,
			BlackUID:       blackUID,
			WhiteUID:       whiteUID,
			BlackNickname:  r.blackNick,
			WhiteNickname:  r.whiteNick,
			StepTimeLimit:  r.StepLimit,
			TotalTimeLimit: r.TotalTime,
			StartAt:        model.UnixMS(time.Now()),
		}, []string{strconv.FormatInt(p.uid, 10)})
	}
	r.startTurn()
	r.pushBoth("room.onTurn", &model.S2CTurnChange{
		Seat:          0,
		StepRemainMs:  int64(r.StepLimit) * 1000,
		TotalRemainMs: int64(r.TotalTime) * 1000,
	})
	return id, nil
}

// clearRoomID 清除 session 的房间绑定（终局后立即执行，允许玩家进入下一局）
func (m *RoomManager) clearRoomID(uid int64) {
	if m.pool == nil {
		return
	}
	if s := m.pool.GetSessionByUID(strconv.FormatInt(uid, 10)); s != nil {
		_ = s.Set("room_id", int64(0))
	}
}

// 延迟清理房间管理器的房间
func (m *RoomManager) removeAfter(id int64, d time.Duration) {
	time.AfterFunc(d, func() {
		m.mu.Lock()
		delete(m.rooms, id)
		m.mu.Unlock()
	})
}
func (m *RoomManager) pushUsers(route string, v interface{}, uids []string) {
	if m.app == nil {
		return
	}
	if err, _ := m.app.SendPushToUsers(route, v, uids, m.frontType); err != nil {
		//推送失败
	}
}

func NewRoomManager(app pitaya.Pitaya, pool session.SessionPool, db *sql.DB, game model.GameConfig) *RoomManager {
	return &RoomManager{
		rooms:     make(map[int64]*Room),
		nextID:    1,
		app:       app,
		pool:      pool,
		db:        db,
		game:      game,
		frontType: app.GetServer().Type,
	}
}
func (m *RoomManager) setRoomID(uid int64, roomID int64) {
	if m.pool == nil {
		return
	}
	if s := m.pool.GetSessionByUID(strconv.FormatInt(uid, 10)); s != nil {
		_ = s.Set("room_id", roomID)
	}
}

// Get 按房间号取房间
func (m *RoomManager) Get(id int64) *Room {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.rooms[id]
}
