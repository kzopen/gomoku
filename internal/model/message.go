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

// S2CGameOver 对应 proto S2CGameOver（push room.onGameOver）
type S2CGameOver struct {
	Winner         int32  `json:"winner"`           // Seat
	Reason         int32  `json:"reason"`           // EndReason
	EndAt          int64  `json:"end_at"`           // unix ms
	Moves          []Move `json:"moves"`            // 完整走子记录（回放）
	AITakeoverSeat int32  `json:"ai_takeover_seat"` // 被 AI 接管的座位，-1=无
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
	Code int32  `json:"code"`
	Msg  string `json:"msg"`
	X    int32  `json:"x"`
	Y    int32  `json:"y"`
}
