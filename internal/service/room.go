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
	State     int32 // RoomPlaying / RoomOver
	StepLimit int32 // 秒
	TotalTime int32 // 秒

	blackNick string
	whiteNick string

	// 回合计时
	stepTimer    *time.Timer // 超时判负
	pushStop     chan struct{}
	stepDeadline time.Time
}

// startTurn 启动当前回合计时：超时判负 + 每秒倒计时推送（剩余≤5s 时加密到 200ms）。
func (r *Room) startGame() {
	if r.stepTimer != nil {
		r.stepTimer.Stop()
	}
	//超时判负
	d := time.Duration(r.StepLimit) * time.Second
	r.stepDeadline = time.Now().Add(d)
	r.stepTimer = time.AfterFunc(d, func() { r.onTimeout() })

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

func (r *Room) pushBoth(route string, v interface{}) {
	r.mgr.pushUsers(route, v, []string{strconv.FormatInt(r.Seat[0], 10), strconv.FormatInt(r.Seat[1], 10)})
}
