package service

import (
	"gomoku/internal/common"
	"gomoku/internal/model"
	"strconv"
	"sync"
	"time"
)

type Room struct {
	mu   sync.Mutex
	mgr  *RoomManager
	ID   int64
	Seat [2]int64 //[0]=黑色 [1]=>白色用户uid

	Board     common.Board
	TurnSeat  int8 //当前轮到谁下棋 0 黑 1白
	Round     int32
	Moves     []model.Move
	State     int32            // RoomPlaying / RoomOver
	StepLimit int32            // 秒
	TotalTime int32            // 秒
	winLine   []model.WinPoint // 制胜连线（五连终局时填充，下发给前端高亮）

	blackNick string
	whiteNick string

	// 回合计时
	stepTimer    *time.Timer // 超时判负
	pushStop     chan struct{}
	stepDeadline time.Time
}

// startTurn 启动当前回合计时：超时判负 + 每秒倒计时推送（剩余≤5s 时加密到 200ms）。
func (r *Room) startTurn() {
	if r.stepTimer != nil {
		r.stepTimer.Stop()
	}
	//超时判负
	d := time.Duration(r.StepLimit) * time.Second
	r.stepDeadline = time.Now().Add(d)
	r.stepTimer = time.AfterFunc(d, func() { r.onTimeout() })

	//todo 这块还不是很理解，重做
	r.stopPush()
	r.pushStop = make(chan struct{})
	go r.pushTicker(r.pushStop)

}

// onTimeout 回合超时：超时方判负（服务器权威计时）
func (r *Room) onTimeout() {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.State != common.RoomPlaying {
		return
	}
	r.finish(1-r.TurnSeat, common.EndTimeout)
}

// finish 统一终局：标记结束、停表、广播、清除房间绑定、延迟清理
func (r *Room) finish(winner int8, reason int32) {
	r.State = common.RoomOver
	if r.stepTimer != nil {
		r.stepTimer.Stop()
	}
	r.stopPush()
	moves := r.Moves
	if moves == nil {
		moves = []model.Move{} // 空局（如超时）序列化为 [] 而非 null
	}
	r.pushBoth("room.onGameOver", &model.S2CGameOver{
		Winner:         int32(winner),
		Reason:         reason,
		EndAt:          model.UnixMS(time.Now()),
		Moves:          moves,
		WinLine:        r.winLine,
		AITakeoverSeat: -1,
	})
	// 立即清除双方 session 的房间room_id绑定，允许玩家马上进入下一局
	r.mgr.clearRoomID(r.Seat[0])
	r.mgr.clearRoomID(r.Seat[1])
	r.mgr.removeAfter(r.ID, 10*time.Second) // 延迟清理，保证双方收完推送
}

// stopPush 停止倒计时推送 goroutine（幂等）
func (r *Room) stopPush() {
	if r.pushStop != nil {
		select {
		case <-r.pushStop: // 已关闭
		default:
			close(r.pushStop)
		}
		r.pushStop = nil
	}
}

// pushTicker 周期性推送 room.onTimer（剩余时间）
func (r *Room) pushTicker(stop chan struct{}) {
	interval := time.Second
	t := time.NewTicker(interval)
	defer t.Stop()
	for {
		select {
		case <-stop:
			return
		case <-t.C:
			r.mu.Lock()
			if r.State != common.RoomPlaying {
				r.mu.Unlock()
				return
			}
			remain := time.Until(r.stepDeadline).Milliseconds()
			if remain < 0 {
				remain = 0
			}
			// 剩余 ≤5s 时加密推送周期
			if remain <= 5000 && interval != 200*time.Millisecond {
				interval = 200 * time.Millisecond
				t.Reset(interval)
			}
			seat := r.TurnSeat
			r.mu.Unlock()
			r.pushBoth("room.onTimer", &model.S2CTimer{
				Seat:          int32(seat),
				StepRemainMs:  remain,
				TotalRemainMs: remain,
			})
		}
	}
}

func (r *Room) Place(uid string, x, y int32) *model.S2CPlacePiece {
	r.mu.Lock()
	defer r.mu.Unlock()
	resp := &model.S2CPlacePiece{}
	if r.State != common.RoomPlaying {
		resp.Code = common.CodeRoomState
		resp.Msg = "房间状态不允许（未开始/已结束）"
		return resp
	}
	seat := r.seatOf(uid)
	if seat < 0 {
		resp.Code = common.CodeRoomNotFound
		resp.Msg = "不在本房间"
		return resp
	}
	if seat != r.TurnSeat {
		resp.Code = common.CodeNotYourTurn
		resp.Msg = "非己方回合"
		return resp
	}
	if x < 0 || x >= common.BoardSize || y < 0 || y >= common.BoardSize || r.Board[x][y] != 0 {
		resp.Code = common.CodeIllegalMove
		resp.Msg = "坐标越界或已有棋子"
		return resp
	}
	color := seat + 1 //color 1=黑 2=白
	r.Board[x][y] = int(color)
	r.Round++
	r.Moves = append(r.Moves, model.Move{
		Seat: int32(seat),
		X:    x,
		Y:    y,
		TS:   model.UnixMS(time.Now()),
	})
	resp.Code = common.CodeOK
	resp.Seat = int32(seat)
	resp.X = x
	resp.Y = y
	resp.Round = r.Round
	r.pushBoth("room.onPlace", resp)
	// 胜负判定：一次遍历同时得判胜结果和制胜连线
	if pts := common.FindWinLine(&r.Board, int(x), int(y), int(color)); len(pts) >= common.WinCount {
		r.winLine = make([]model.WinPoint, 0, len(pts))
		for _, p := range pts {
			r.winLine = append(r.winLine, model.WinPoint{X: int32(p.X), Y: int32(p.Y)})
		}
		r.finish(seat, common.EndFive)
		return resp
	}
	r.TurnSeat = 1 - seat
	r.startTurn()
	r.pushBoth("room.onTurn", &model.S2CTurnChange{
		Seat:          int32(r.TurnSeat),
		StepRemainMs:  int64(r.StepLimit) * 1000,
		TotalRemainMs: int64(r.TotalTime) * 1000,
	})
	return resp
}

// seatOf 返回 uid 的座位；不在房间返回 -1
func (r *Room) seatOf(uid string) int8 {
	for i := 0; i < 2; i++ {
		if strconv.FormatInt(r.Seat[i], 10) == uid {
			return int8(i)
		}
	}
	return -1
}

func (r *Room) pushBoth(route string, v interface{}) {
	r.mgr.pushUsers(route, v, []string{strconv.FormatInt(r.Seat[0], 10), strconv.FormatInt(r.Seat[1], 10)})
}

// Snapshot 构造重连恢复的全量状态（room.reconnect）。
// uid 不在本房间时返回 nil。
func (r *Room) Snapshot(uid string) *model.S2CReconnect {
	r.mu.Lock()
	defer r.mu.Unlock()
	seat := r.seatOf(uid)
	if seat < 0 {
		return nil
	}
	moves := r.Moves
	if moves == nil {
		moves = []model.Move{}
	}
	remain := int64(0)
	if r.State == common.RoomPlaying && !r.stepDeadline.IsZero() {
		remain = time.Until(r.stepDeadline).Milliseconds()
		if remain < 0 {
			remain = 0
		}
	}
	return &model.S2CReconnect{
		Code:           common.CodeOK,
		State:          r.State,
		Seat:           int32(seat),
		Moves:          moves,
		TurnSeat:       int32(r.TurnSeat),
		StepRemainMs:   remain,
		TotalRemainMs:  0,
		AITakeover:     false,
		Spectator:      false,
		BlackUID:       r.Seat[0],
		WhiteUID:       r.Seat[1],
		BlackNickname:  r.blackNick,
		WhiteNickname:  r.whiteNick,
		StepTimeLimit:  r.StepLimit,
		TotalTimeLimit: r.TotalTime,
	}
}

// NotifyBack 广播玩家重连恢复（room.onPlayerBack，通知对方）
func (r *Room) NotifyBack(seat int32) {
	r.pushBoth("room.onPlayerBack", &model.S2CPlayerBack{Seat: seat})
}
