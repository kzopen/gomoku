package model

type C2SLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type S2CLogin struct {
	Code     int32  `json:"code"`
	Msg      string `json:"msg"`
	UID      int64  `json:"uid"`
	Token    string `json:"token"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	ELO      int32  `json:"elo"`
	Wins     int32  `json:"wins"`
	Losses   int32  `json:"losses"`
	Draws    int32  `json:"draws"`
}

type C2SJoinMatch struct {
}

type S2CJoinMatch struct {
	Code int32  `json:"code"`
	Msg  string `json:"msg"`
}
type C2SCancelMatch struct{}

type S2CCancelMatch struct {
	Code int32  `json:"code"`
	Msg  string `json:"msg"`
}

// S2CMatchSuccess 对应 proto S2CMatchSuccess（push room.onMatchSuccess）
type S2CMatchSuccess struct {
	RoomID         int64  `json:"room_id"`
	Seat           int32  `json:"seat"` // Seat
	BlackUID       int64  `json:"black_uid"`
	WhiteUID       int64  `json:"white_uid"`
	BlackNickname  string `json:"black_nickname"`
	WhiteNickname  string `json:"white_nickname"`
	StepTimeLimit  int32  `json:"step_time_limit"`  // 秒
	TotalTimeLimit int32  `json:"total_time_limit"` // 秒
	StartAt        int64  `json:"start_at"`         // unix ms
}

// Move 对应 proto Move
type Move struct {
	Seat int32 `json:"seat"`
	X    int32 `json:"x"`
	Y    int32 `json:"y"`
	TS   int64 `json:"ts"` // unix ms
}

type S2CTurnChange struct {
	Seat          int32 `json:"seat"`            // 当前行动方
	StepRemainMs  int64 `json:"step_remain_ms"`  // 当前方剩余步时
	TotalRemainMs int64 `json:"total_remain_ms"` // 当前方剩余总时
}

// WinPoint 制胜连线上的一个坐标（S2CGameOver.win_line）
type WinPoint struct {
	X int32 `json:"x"`
	Y int32 `json:"y"`
}

// S2CGameOver 对应 proto S2CGameOver（push room.onGameOver）
type S2CGameOver struct {
	Winner         int32      `json:"winner"`           // Seat
	Reason         int32      `json:"reason"`           // EndReason
	EndAt          int64      `json:"end_at"`           // unix ms
	Moves          []Move     `json:"moves"`            // 完整走子记录（回放）
	WinLine        []WinPoint `json:"win_line"`         // 制胜连线坐标（前端高亮；非五连结束为空）
	AITakeoverSeat int32      `json:"ai_takeover_seat"` // 被 AI 接管的座位，-1=无
}

// S2CTimer 对应 proto S2CTimer（push room.onTimer）
type S2CTimer struct {
	Seat          int32 `json:"seat"`
	StepRemainMs  int64 `json:"step_remain_ms"`
	TotalRemainMs int64 `json:"total_remain_ms"`
}
type C2SPlacePiece struct {
	X int32 `json:"x"`
	Y int32 `json:"y"`
}

type S2CPlacePiece struct {
	Code  int32  `json:"code"`
	Msg   string `json:"msg"`
	Seat  int32  `json:"seat"`
	X     int32  `json:"x"`
	Y     int32  `json:"y"`
	Round int32  `json:"round"` // 第几手，从 1 开始
}

// ==================== 重连 / 离开 ====================

// C2SReconnect room.reconnect 请求（断线恢复）
type C2SReconnect struct {
	RoomID int64  `json:"room_id"`
	Token  string `json:"token"`
}

// S2CReconnect room.reconnect 响应：全量对局状态（前端据此恢复棋盘/计时）
type S2CReconnect struct {
	Code           int32  `json:"code"`
	Msg            string `json:"msg"`
	State          int32  `json:"state"` // RoomPlaying / RoomOver
	Seat           int32  `json:"seat"`  // 本玩家座位（观战 -1）
	Moves          []Move `json:"moves"`
	TurnSeat       int32  `json:"turn_seat"`
	StepRemainMs   int64  `json:"step_remain_ms"`
	TotalRemainMs  int64  `json:"total_remain_ms"`
	AITakeover     bool   `json:"ai_takeover"`
	Spectator      bool   `json:"spectator"`
	BlackUID       int64  `json:"black_uid"`
	WhiteUID       int64  `json:"white_uid"`
	BlackNickname  string `json:"black_nickname"`
	WhiteNickname  string `json:"white_nickname"`
	StepTimeLimit  int32  `json:"step_time_limit"`
	TotalTimeLimit int32  `json:"total_time_limit"`
}

// C2SLeaveGame room.leave 请求（主动离开判负）
type C2SLeaveGame struct{}

// S2CLeaveGame room.leave 响应
type S2CLeaveGame struct {
	Code int32  `json:"code"`
	Msg  string `json:"msg"`
}

// S2CPlayerBack push room.onPlayerBack（对方重连恢复）
type S2CPlayerBack struct {
	Seat int32 `json:"seat"`
}

// ==================== 排行榜 / 在线 ====================

// C2SOnlineCount rank.online_count 请求（在线人数）
type C2SOnlineCount struct{}

// S2COnlineCount rank.online_count 响应
type S2COnlineCount struct {
	Code  int32  `json:"code"`
	Msg   string `json:"msg"`
	Count int64  `json:"count"`
}

// C2SPing rank.ping 请求（心跳续期，前端 notify 发送）
type C2SPing struct{}

// S2CPing rank.ping 响应（notify 无响应，占位）
type S2CPing struct {
	Code int32 `json:"code"`
}
