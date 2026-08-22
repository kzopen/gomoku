package common

const (
	CodeOK            = 0    // 成功
	CodeBadParam      = 1001 // 参数错误
	CodeNotLoggedIn   = 1002 // 未登录 / 会话失效
	CodeAlreadyIn     = 2001 // 已在匹配中 / 已在房间中
	CodeMatchFailed   = 2002 // 匹配队列已满或撮合失败
	CodeNotYourTurn   = 3001 // 非己方回合
	CodeIllegalMove   = 3002 // 坐标越界或已有棋子
	CodeRoomState     = 3003 // 房间状态不允许（未开始/已结束）
	CodeAITakeovered  = 3004 // 已被 AI 接管，无权落子
	CodeRoomNotFound  = 4001 // 房间不存在或已销毁
	CodeInternalError = 5001 // 服务器内部错误
	CodeNotImpl       = 5002 // 功能尚未实现（骨架占位）
)

type Resp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}
