package common

import (
	"testing"

	"gomoku/internal/model"
)

func TestSuccessResponse_setsCodeAndMessage(t *testing.T) {
	resp := SuccessResponse(&model.S2CJoinMatch{}, "已进入匹配队列")
	if resp.Code != CodeOK || resp.Msg != "已进入匹配队列" {
		t.Fatalf("unexpected response: code=%d msg=%q", resp.Code, resp.Msg)
	}
}

func TestErrorResponse_setsCodeAndMessage(t *testing.T) {
	resp := ErrorResponse(&model.S2CPlacePiece{}, CodeBadParam, "参数错误")
	if resp.Code != CodeBadParam || resp.Msg != "参数错误" {
		t.Fatalf("unexpected response: code=%d msg=%q", resp.Code, resp.Msg)
	}
}

func TestSuccessResponse_supportsCommonResponse(t *testing.T) {
	resp := SuccessResponse(&Resp{}, "ok")
	if resp.Code != CodeOK || resp.Msg != "ok" {
		t.Fatalf("unexpected response: code=%d msg=%q", resp.Code, resp.Msg)
	}
}
