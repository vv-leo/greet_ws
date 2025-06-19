package resp

type WsMessageModel struct {
	CgiBaseResponse struct {
		Ret    int    `json:"Ret"`
		ErrMsg string `json:"ErrMsg"`
	} `json:"CgiBaseResponse"`
	ResponseData interface{} `json:"ResponseData"`
	Data         interface{} `json:"Data"`
}
