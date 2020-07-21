/**
@Author: wei-g
@Date:   2020/3/26 2:57 下午
@Description: 提供便捷的返回数据格式
*/

package data

const (
	// FailedCode 处理失败 其他错误
	FailedCode = iota - 1 // -1
	// SucCode 处理成功
	SucCode // 0
	// RequestErr 请求方错误
	RequestErr = 1000
	// 请求参数错误
	RequestParaErr = 1001
	// 响应方错误
	ResponseErr = 2001
)

const (
	FailedMsg = "failed"
	SucMsg    = "success"
)

// ResultJson 前后端分离开发时 约定好的标准的数据结构
type ResultJson struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// NewResultJson 构造一个标准的数据
func NewResultJson(code int, message string, data interface{}) (result *ResultJson) {
	return &ResultJson{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

// SucJson 构建一个处理成功的数据
func SucJson(data interface{}) (result *ResultJson) {
	return &ResultJson{
		Code:    SucCode,
		Message: SucMsg,
		Data:    data,
	}
}

// OtherFailedJson 其他错误的错误
func OtherFailedJson(msg string) (result *ResultJson) {
	return &ResultJson{
		Code:    FailedCode,
		Message: msg,
		Data:    "",
	}
}

// RequestErrJson 请求端错误
func RequestErrJson(msg string) (result *ResultJson) {
	return &ResultJson{
		Code:    RequestErr,
		Message: msg,
		Data:    "",
	}
}
