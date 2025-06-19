package resp

type WsMessageModel struct {
	CgiBaseResponse struct {
		Ret    int    `json:"Ret"`
		ErrMsg string `json:"ErrMsg"`
	} `json:"CgiBaseResponse"`
	ResponseData interface{} `json:"ResponseData"`
	Data         interface{} `json:"Data"`
}

/*type MessageSuccessModel struct {
	CgiBaseResponse struct {
		Ret    int    `json:"Ret"`
		ErrMsg string `json:"ErrMsg"`
	} `json:"CgiBaseResponse"`
	ResponseData struct {
		Timestamp    time.Time `json:"Timestamp"`
		ID           string    `json:"ID"`
		ServerID     int       `json:"ServerID"`
		DebugTimings struct {
			Queue           int `json:"Queue"`
			Marshal         int `json:"Marshal"`
			GetParticipants int `json:"GetParticipants"`
			GetDevices      int `json:"GetDevices"`
			GroupEncrypt    int `json:"GroupEncrypt"`
			PeerEncrypt     int `json:"PeerEncrypt"`
			Send            int `json:"Send"`
			Resp            int `json:"Resp"`
			Retry           int `json:"Retry"`
		} `json:"DebugTimings"`
	} `json:"ResponseData"`
	Data interface{} `json:"Data"`
}*/
