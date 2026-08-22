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
