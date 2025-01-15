package api

type Resp struct {
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
	Code int         `json:"code"`
}

type CheckStatusResp struct {
	Status bool `json:"status"`
}
