package req

// 分页请求结构体
type PageInfoModel struct {
	Page     int `json:"page" form:"page" binding:"required" comment:"页码，从1开始"`
	PageSize int `json:"pageSize" form:"pageSize" binding:"required" comment:"每页记录数"`
}

// Find by id structure
type IdModel struct {
	Id string `json:"id" form:"id" binding:"required" `
}

// Find by id structure
type ContactIdModel struct {
	Id string `json:"id" form:"id" binding:"required" ` //联系人id
}

// Find by id int64 structure
type IdInt64Model struct {
	Id int64 `json:"id" binding:"required"`
}

// Find by id int64 structure
type IdBigIntModel struct {
	Id string `json:"id" binding:"required"`
}

// Find by name structure
type NameModel struct {
	Name string `json:"name" form:"name" binding:"required"`
}

// Find by name structure
type GroupNameModel struct {
	Name string `json:"name" form:"name" binding:"required"` //联系人分组名
}

// Find by ids structure
type IdsModel struct {
	Ids []string `json:"ids" form:"ids"`
}

// Find by fiel structure
type FileModel struct {
	File string `json:"file" form:"file"` //base64编码的文件字符串（前面不带“data:image/png;base64,”）。
}
