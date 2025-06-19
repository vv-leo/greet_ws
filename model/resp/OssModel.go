package resp

type OssModel struct {
	Code            int    `json:"code"`              // oss响应代码
	OssFile         string `json:"oss_file"`          // oss外部文件路径
	OssInternalFile string `json:"oss_internal_file"` // oss内部文件路径
}
