package basis

// 请求类型
type FrontendRequest struct {
	Id     string      `json:"id"`
	Type   string      `json:"type"`
	Action string      `json:"action"`
	Token  string      `json:"token"`
	Appkey string      `json:"appkey"` // 新增
	Data   interface{} `json:"data"`
}

// 前端期待的响应类型
type FrontendResponse struct {
	RequestId     string      `json:"requestId"`
	RequestKeep   bool        `json:"requestKeep"`
	ResponseError string      `json:"responseError"`
	ResponseTip   string      `json:"responseTip"`
	ResponseData  interface{} `json:"responseData"`
	SoftwareId    uint32      `json:"softwareId"` // 新增。
}

// 服务请求类型
type SoftwareRequest struct {
	Id   int64       `json:"id"`
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// 服务响应类型
type SoftwareResponse struct {
	Id   int64       `json:"id"`
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}
