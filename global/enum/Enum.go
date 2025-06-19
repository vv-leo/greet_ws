/**
 * @Description $
 * @Author $
 * @Date $ $
 **/
package enum

const (
	//登录类型，短信
	LOGIN_TYPE_SMS = "sms"

	//发送类型账号
	SEND_TYPE_SMS_ACCOUNT = "account"

	//发送类型手机号
	SEND_TYPE_SMS_PHONE = "phone"
)

const (
	MessageTypeText     = "text"  //文本消息
	MessageTypeMedia    = "media" //发送消息时的图片消息是media
	MessageTypeImage    = "image" //收消息时图片保存到表的类型是image
	MessageTypeAudio    = "audio" //语音
	MessageTypeVideo    = "video" //视频
	MessageTypeEvent    = "event" //事件消息
	MessageTypeVcard    = "vcard" //名片
	MessageTypeRichText = "richText"
)

const (
	RichTextTypeQuoted    = "quoted"    //引用
	RichTextTypeImageText = "imageText" //图文
	RichTextTypeHyperlink = "hyperLink" //超链接
)
const (
	PlatformAutoReply = "auto_reply" //AI客服系统自动推进模块项目
)

const (
	WsAccountProxyTypeStatic  = iota + 1 //静态代理
	WsAccountProxyTypeDynamic            //动态代理
)
