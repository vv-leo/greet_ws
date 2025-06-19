package resp

type TimeSegmentModel struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

type GreetMessageContent struct {
	Type           string   `json:"type"`
	TextList       []string `json:"textList"`
	ImageUrl       string   `json:"imageUrl"`
	Title          string   `json:"title"`
	SecondaryTitle string   `json:"secondaryTitle"`
	Content        string   `json:"content"`
	UrlList        []string `json:"urlList"`
}

type TenantModel struct {
	ID                    string               `json:"id"`
	CreateBy              string               `json:"createBy"`
	CreateTime            string               `json:"createTime"`
	UpdateBy              string               `json:"updateBy"`
	UpdateTime            string               `json:"updateTime"`
	GroupID               *string              `json:"groupId"`
	MaxSendQuantityPerDay int                  `json:"maxSendQuantityPerDay"`
	TranslateConfig       *string              `json:"translateConfig"`
	TimeSegmentList       []int64              `json:"timeSegmentList"`
	GreetTextList         []string             `json:"greetTextList"`
	MessagesPerMinute     int                  `json:"messagesPerMinute"`
	GreetMessageContent   *GreetMessageContent `json:"greetMessageContent"`
	AllocSession          *string              `json:"allocSession"`
	CountryCode           string               `json:"countryCode"`
}
