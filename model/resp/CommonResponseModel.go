package resp

// 分页数据响应结构体
type PageResult struct {
	List     interface{} `json:"list" comment:"数据列表"`
	Total    int64       `json:"total" comment:"记录总数"`
	Page     int         `json:"page" comment:"页码"`
	PageSize int         `json:"pageSize" comment:"每页显示条数"`
}
