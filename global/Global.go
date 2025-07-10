package global

import (
	"fmt"
	"git.exclouds.org/ww/greet_ws/config"
	"github.com/go-redis/redis"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"time"
)

var (
	SERVER_CONFIG config.ServerConfig
	LOGGER        *zap.Logger
	REDIS         *redis.Client
	VIPER         *viper.Viper
)

// 外部API URL变量
var (
	WsSendTxtUrl                = ""
	WsSendImageUrl              = ""
	WsSendRichTextUrl           = ""
	WsAddContactsUrl            = ""
	WsInitsession               = ""
	WsDeteleChatUrl             = ""
	WsAchiveUrl                 = ""
	WsMuteChatUrl               = ""
	ChatHistorySaveMessagesUrl  = ""
	ChatHistorySaveContactsUrl  = ""
	ChatHistoryGetContactUrl    = ""
	ChatHistoryUpdateContactUrl = ""

	ChatWorkQueryWaitingTaskUrl      = ""
	ChatWorkExecuteDetailUrl         = ""
	ChatWorkTenantHelloUrl           = ""
	ChatAccountAssignUrl             = ""
	FileMediaDownloadAndUploadApiUrl = ""
	ChatSocketPushMsgUrl             = ""
)

// InitExternalApiUrl 初始化外部接口URL
func InitExternalApiUrl() {
	wsServerUrl := SERVER_CONFIG.EndpointsConfig.WsServer.HttpUrl
	//ws接口
	WsSendTxtUrl = wsServerUrl + "/v1/message/sendtext"        //发送文本消息
	WsSendImageUrl = wsServerUrl + "/v1/message/sendimage"     //发送图片消息
	WsSendRichTextUrl = wsServerUrl + "/v1/message/adreply"    //发送富文本消息（图文）
	WsAddContactsUrl = wsServerUrl + "/v1/profile/addcontacts" //添加联系人
	WsInitsession = wsServerUrl + "/v1/message/initsession"    //初始化会话
	WsDeteleChatUrl = wsServerUrl + "/v1/message/deletechat"   //删除会话
	WsAchiveUrl = wsServerUrl + "/v1/message/archivechat"      //归档
	WsMuteChatUrl = wsServerUrl + "/v1/message/mutechat"       //静音

	// ChatHistory 服务接口路径
	ChatHistorySaveMessagesUrl = "/api/v1/im/saveMessage"         //批量保存消息
	ChatHistorySaveContactsUrl = "/api/v1/im/saveContact"         //保存联系人
	ChatHistoryGetContactUrl = "/api/v1/im/getContactByCondition" //获取联系人信息
	ChatHistoryUpdateContactUrl = "/api/v1/im/updateContact"      //更新联系人信息

	// ChatWork 服务接口路径
	ChatWorkQueryWaitingTaskUrl = "/chat_work/v1/task/query_waiting" //查询等待任务
	ChatWorkExecuteDetailUrl = "/chat_work/v1/task/execute_detail"   //执行任务详情
	ChatWorkTenantHelloUrl = "/chat_work/v1/tenant/hello"            //租户问候配置

	// ChatAccount 服务接口路径
	ChatAccountAssignUrl = "/account/v1/assign" //账号分配

	FileMediaDownloadAndUploadApiUrl = "/api/v1/im/mediaDownloadAndUpload" // 媒体文件上传下载
	ChatSocketPushMsgUrl = "/api/v1/push"                                  // 消息推送
}

func GenerateUniqueId() string {
	charset := "0123456789"
	length := 10

	timestamp := time.Now().UnixMilli() // 毫秒级时间戳
	randomID, err := gonanoid.Generate(charset, length)
	if err != nil {
		LOGGER.Error(err.Error())
		return "-1"
	}

	id := fmt.Sprintf("%d%s", timestamp, randomID)
	return id
}
