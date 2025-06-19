package service

import (
	"git.exclouds.org/ww/greet_ws/global"
	"git.exclouds.org/ww/greet_ws/utils"
	"git.exclouds.org/ww/greet_ws/utils/encrypt"
	"time"
)

// 工具
var bcryptUtils encrypt.BcryptUtils
var stringUtils utils.StringUtils
var redisUtils utils.RedisUtils
var uuidUtils utils.UUidUtils
var timeUtils utils.TimeUtils

const HttpRequestTimeout = 100

var (
	TaskPullUrl            = ""
	TaskPullAccountListUrl = ""
	TaskUpdateUrl          = ""
	UserLoginUrl           = ""
	ProxyChangeUrl         = ""

	GetTenantList         = ""
	GetTargetList         = ""
	GetWsAccountList      = ""
	WsTaskFinish          = ""
	TargetFinish          = ""
	WsAccountSyncLogin    = ""
	WsAccountForceReLogin = ""

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
	taskServerUrl := global.SERVER_CONFIG.EndpointsConfig.TaskServer.HttpUrl
	if taskServerUrl == "" {
		global.LOGGER.Error("taskServerUrl配置为空没读取到")
		return
	}

	//任务接口
	TaskPullUrl = taskServerUrl + "/ru/task/pull"
	TaskPullAccountListUrl = taskServerUrl + "/ru/task/pullAccountList"
	TaskUpdateUrl = taskServerUrl + "/ru/task/update"
	UserLoginUrl = taskServerUrl + "/user/login"
	ProxyChangeUrl = taskServerUrl + "/ws/account/proxy/change" //代理切换

	GetTenantList = taskServerUrl + "/user/tenant/list"                //获取租户列表
	GetTargetList = taskServerUrl + "/target/data/fetch"               //获取目标联系人
	GetWsAccountList = taskServerUrl + "/ws/task/pullAccountList"      //获取ws账号
	WsTaskFinish = taskServerUrl + "/ws/task/finish"                   //ws账号完成上报
	TargetFinish = taskServerUrl + "/target/data/update"               //目标联系人上报
	WsAccountSyncLogin = taskServerUrl + "/ws/account/syncLogin"       //ws账号登录
	WsAccountForceReLogin = taskServerUrl + "/ws/account/forceReLogin" //ws账号强制登录

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

// 获取分布式锁
func acquireLock(lockKey string, expiration time.Duration) (bool, error) {
	// 使用 SETNX 尝试获取锁
	ok, err := global.REDIS.SetNX(lockKey, "locked", expiration).Result()
	if err != nil {
		return false, err
	}
	return ok, nil
}

// 释放分布式锁
func releaseLock(lockKey string) error {
	return global.REDIS.Del(lockKey).Err()
}
