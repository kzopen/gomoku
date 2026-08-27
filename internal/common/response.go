package common

// Response 是所有带 code/msg 的业务响应所实现的最小接口。
type Response interface {
	SetCode(int32)
	SetMsg(string)
}

// SuccessResponse 构造成功响应并统一设置提示信息。
func SuccessResponse[T Response](resp T, msg string) T {
	resp.SetCode(CodeOK)
	resp.SetMsg(msg)
	return resp
}

// ErrorResponse 构造错误响应并统一设置错误码和提示信息。
func ErrorResponse[T Response](resp T, code int32, msg string) T {
	resp.SetCode(code)
	resp.SetMsg(msg)
	return resp
}
