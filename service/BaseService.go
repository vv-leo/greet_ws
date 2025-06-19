package service

import (
	"git.exclouds.org/ww/greet_ws/global"
	"git.exclouds.org/ww/greet_ws/utils"
)

// 工具
var redisUtils utils.RedisUtils

const HttpRequestTimeout = 100

var (
	WsSendTxtUrl      = ""
	WsSendImageUrl    = ""
	WsSendRichTextUrl = ""
	WsAddContactsUrl  = ""
	WsInitsession     = ""
	WsDeteleChatUrl   = ""
	WsAchiveUrl       = ""
	WsMuteChatUrl     = ""
)

// 外部接口配置
func InitExternalApiUrl() {
	wsServerUrl := global.SERVER_CONFIG.EndpointsConfig.WsServer.HttpUrl
	//ws接口
	WsSendTxtUrl = wsServerUrl + "/v1/message/sendtext"        //发送文本消息
	WsSendImageUrl = wsServerUrl + "/v1/message/sendimage"     //发送图片消息
	WsSendRichTextUrl = wsServerUrl + "/v1/message/adreply"    //发送富文本消息（图文）
	WsAddContactsUrl = wsServerUrl + "/v1/profile/addcontacts" //添加联系人
	WsInitsession = wsServerUrl + "/v1/message/initsession"    //初始化会话
	WsDeteleChatUrl = wsServerUrl + "/v1/message/deletechat"   //删除会话
	WsAchiveUrl = wsServerUrl + "/v1/message/archivechat"      //归档
	WsMuteChatUrl = wsServerUrl + "/v1/message/mutechat"       //静音
}
