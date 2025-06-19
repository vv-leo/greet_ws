package resp

// 消息返回（普通）
type WSPacketCommonModel struct {
	EventData   interface{}            `json:"EventData"`
	EventName   string                 `json:"EventName"`
	OtherFields map[string]interface{} `json:"OtherFields,omitempty"`
}

type WSPacketRichText struct {
	WSPacket struct {
		EventData struct {
			FromJID   string `json:"fromJID"`
			MessageID string `json:"messageID"`
			Timestamp int64  `json:"timestamp"`
		} `json:"EventData"`
		EventName string `json:"EventName"`
		WSInfo    struct {
			ApiInfo string `json:"ApiInfo"`
			RegInfo string `json:"RegInfo"`
			UUID    string `json:"UUID"`
			UserJID string `json:"UserJID"`
		} `json:"WSInfo"`
	} `json:"WSPacket"`
}

// 消息返回（普通）
type WSPacketModel struct {
	EventData struct {
		FromJID       string `json:"fromJID"`
		MessageID     string `json:"messageID"`
		MessageSource struct {
			Conversation       string `json:"conversation"`
			MessageContextInfo struct {
				DeviceListMetadata struct {
					SenderTimestamp int64 `json:"senderTimestamp"`
				} `json:"deviceListMetadata"`
				DeviceListMetadataVersion int `json:"deviceListMetadataVersion"`
			} `json:"messageContextInfo"`
		} `json:"messageSource"`
		MessageType int   `json:"messageType"`
		Timestamp   int64 `json:"timestamp"`
	} `json:"EventData"`
	EventName string `json:"EventName"`
	WSInfo    struct {
		ApiInfo string `json:"ApiInfo"`
		UUID    string `json:"UUID"`
		UserJID string `json:"UserJID"`
	} `json:"WSInfo"`
}

// 消息返回（富文本）
type WSPacketRichTextModel struct {
	EventData struct {
		FromJID       string `json:"fromJID"`
		MessageID     string `json:"messageID"`
		MessageSource struct {
			ExtendedTextMessage struct {
				Text        string `json:"text"`
				PreviewType int    `json:"previewType"`
				ContextInfo struct {
					EntryPointConversionSource       string `json:"entryPointConversionSource"`
					EntryPointConversionDelaySeconds int    `json:"entryPointConversionDelaySeconds"`
				} `json:"contextInfo"`
				InviteLinkGroupTypeV2 int `json:"inviteLinkGroupTypeV2"`
			} `json:"extendedTextMessage"`
			MessageContextInfo struct {
				DeviceListMetadata struct {
					SenderTimestamp int64 `json:"senderTimestamp"`
				} `json:"deviceListMetadata"`
				DeviceListMetadataVersion int `json:"deviceListMetadataVersion"`
			} `json:"messageContextInfo"`
		} `json:"messageSource"`
		Timestamp int64 `json:"timestamp"`
	} `json:"EventData"`
	EventName string `json:"EventName"`
	WSInfo    struct {
		ApiInfo string `json:"ApiInfo"`
		UUID    string `json:"UUID"`
		UserJID string `json:"UserJID"`
	} `json:"WSInfo"`
}
