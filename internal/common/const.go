package common

const (
	BoardSize = 15
	WinCount  = 5
)

// 座位：黑方先手
const (
	SeatBlack int8 = 0
	SeatWhite int8 = 1
)

// 房间状态（
const (
	RoomPlaying int32 = 0 //游戏开始
	RoomOver    int32 = 1 //游戏结束
)

// 对局参数默认值（可被 config/game 覆盖）
const (
	DefaultStepTimeLimit   = 30 // 秒
	DefaultTotalTimeLimit  = 0  // 秒，0=不限
	DefaultTimerPushInterv = 1  // 秒
	DefaultHeartbeatTo     = 15 // 秒
	DefaultAITakeoverGrace = 10 // 秒
	DefaultAILevel         = 2
)
const (
	EndUnknown int32 = 0
	EndFive    int32 = 1 // 五连
	EndTimeout int32 = 2 // 超时判负
	EndLeave   int32 = 3 // 离开判负
	EndDraw    int32 = 4 // 和棋（预留）
)
