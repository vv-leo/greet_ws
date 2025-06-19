package req

// CgiRequest 请求ws结构体
type CgiRequest struct {
	ToUserId  string `json:"ToUserId"`
	IsGroup   bool   `json:"IsGroup"`
	IsFakeMsg bool   `json:"IsFakeMsg"`
}
type WithCgiRequest struct {
	CgiRequest
	StanzaId    string `json:"StanzaId"`
	Participant string `json:"Participant"`
	Content     string `json:"Content"`
}

// WsTextModel 文本消息
type WsSendMessageTextModel struct {
	WithCgiRequest WithCgiRequest `json:"CgiRequest"`
}

type WsSendMessageRichTextModel struct {
	WithCgiRequest RichTextCgiRequest `json:"CgiRequest"`
}

type ImageCgiRequest struct {
	WithCgiRequest
	ImageUrl string `json:"ImageUrl"`
}

// richText
type RichTextCgiRequest struct {
	ToUserId  string `json:"ToUserId"`
	IsGroup   bool   `json:"IsGroup"`
	Text      string `json:"Text"`
	Title     string `json:"Title"`
	Body      string `json:"Body"`
	ImageUrl  string `json:"ImageUrl"`
	SourceUrl string `json:"SourceUrl"`
}

// WsSendMessageImageModel 图片消息Url
type WsSendMessageImageUrlModel struct {
	ImageCgiRequest ImageCgiRequest `json:"CgiRequest"`
}

type WsDownLoadImgModel struct {
	CgiRequest struct {
		Message struct {
			ImageMessage struct {
				URL                  string `json:"url"`
				Mimetype             string `json:"mimetype"`
				FileSha256           string `json:"fileSha256"`
				FileLength           int64  `json:"fileLength"`
				Height               int    `json:"height"`
				Width                int    `json:"width"`
				MediaKey             string `json:"mediaKey"`
				FileEncSha256        string `json:"fileEncSha256"`
				DirectPath           string `json:"directPath"`
				MediaKeyTimestamp    int64  `json:"mediaKeyTimestamp"`
				JpegThumbnail        string `json:"jpegThumbnail"`
				FirstScanSidecar     string `json:"firstScanSidecar"`
				FirstScanLength      int    `json:"firstScanLength"`
				ScansSidecar         string `json:"scansSidecar"`
				ScanLengths          []int  `json:"scanLengths"`
				MidQualityFileSha256 string `json:"midQualityFileSha256"`
				ImageSourceType      int    `json:"imageSourceType"`
			} `json:"imageMessage"`
			MessageContextInfo struct {
				DeviceListMetadata struct {
					RecipientTimestamp int64 `json:"recipientTimestamp"`
				} `json:"deviceListMetadata"`
				DeviceListMetadataVersion int    `json:"deviceListMetadataVersion"`
				MessageSecret             string `json:"messageSecret"`
			} `json:"messageContextInfo"`
		} `json:"Message"`
	} `json:"CgiRequest"`
}

// 视频消息结构体
type WsDownLoadVideoModel struct {
	CgiRequest struct {
		Message struct {
			VideoMessage struct {
				URL               string `json:"url"`
				Mimetype          string `json:"mimetype"`
				FileSha256        string `json:"fileSha256"`
				FileLength        int64  `json:"fileLength"`
				Seconds           int    `json:"seconds"`
				MediaKey          string `json:"mediaKey"`
				Height            int    `json:"height"`
				Width             int    `json:"width"`
				FileEncSha256     string `json:"fileEncSha256"`
				DirectPath        string `json:"directPath"`
				MediaKeyTimestamp int64  `json:"mediaKeyTimestamp"`
				JpegThumbnail     string `json:"jpegThumbnail"`
				StreamingSidecar  string `json:"streamingSidecar"`
			} `json:"videoMessage"`
			MessageContextInfo struct {
				DeviceListMetadata struct {
					RecipientTimestamp int64 `json:"recipientTimestamp"`
				} `json:"deviceListMetadata"`
				DeviceListMetadataVersion int    `json:"deviceListMetadataVersion"`
				MessageSecret             string `json:"messageSecret"`
			} `json:"messageContextInfo"`
		} `json:"Message"`
	} `json:"CgiRequest"`
}

// 音频消息结构体
type WsDownLoadAudioModel struct {
	CgiRequest struct {
		Message struct {
			AudioMessage struct {
				URL               string `json:"url"`
				Mimetype          string `json:"mimetype"`
				FileSha256        string `json:"fileSha256"`
				FileLength        int64  `json:"fileLength"`
				Seconds           int    `json:"seconds"`
				Ptt               bool   `json:"ptt"`
				MediaKey          string `json:"mediaKey"`
				FileEncSha256     string `json:"fileEncSha256"`
				DirectPath        string `json:"directPath"`
				MediaKeyTimestamp int64  `json:"mediaKeyTimestamp"`
				ContextInfo       struct {
					Expiration                int64 `json:"expiration"`
					EphemeralSettingTimestamp int64 `json:"ephemeralSettingTimestamp"`
					DisappearingMode          struct {
						Initiator     int  `json:"initiator"`
						Trigger       int  `json:"trigger"`
						InitiatedByMe bool `json:"initiatedByMe"`
					} `json:"disappearingMode"`
				} `json:"contextInfo"`
				StreamingSidecar string `json:"streamingSidecar"`
				Waveform         string `json:"waveform"`
			} `json:"audioMessage"`
			MessageContextInfo struct {
				DeviceListMetadata struct {
					RecipientTimestamp int64 `json:"recipientTimestamp"`
				} `json:"deviceListMetadata"`
				DeviceListMetadataVersion int    `json:"deviceListMetadataVersion"`
				MessageSecret             string `json:"messageSecret"`
			} `json:"messageContextInfo"`
		} `json:"Message"`
	} `json:"CgiRequest"`
}

// 添加联系人
type WsAddContactsModel struct {
	CgiRequest struct {
		Contacts []string `json:"Contacts"`
	} `json:"CgiRequest"`
}

// 初始化会话
type Initsession struct {
	CgiRequest struct {
		ToUserId string `json:"ToUserId"`
	} `json:"CgiRequest"`
}

// 初始化会话
type Deletechat struct {
	CgiRequest struct {
		ToUserId string `json:"ToUserId"`
	} `json:"CgiRequest"`
}

// 归档
type Archive struct {
	CgiRequest struct {
		ToUserId string `json:"ToUserId"`
		Achive   bool   `json:"Achive"`
	} `json:"CgiRequest"`
}

// 静音
type MuteChat struct {
	CgiRequest struct {
		ToUserId string `json:"ToUserId"`
		IsMute   bool   `json:"IsMute"`
	} `json:"CgiRequest"`
}
