package resp

// 验证码
type SysCaptchaResponse struct {
	CaptchaId string `json:"captchaId" comment:"验证码id"`
	PicPath   string `json:"picPath" comment:"验证码图片地址"`
}
