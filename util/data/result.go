/**
@Author: wei-g
@Date:   2020/3/26 2:57 下午
@Description: 提供便捷的返回数据格式
*/

package data

const (
	// FailedCode 处理失败 其他错误
	FailedCode = -1
	// SucCode 处理成功
	SucCode = 200
	// 响应方内部错误
	ResponseInternalErr = 5000
	// 请求参数错误
	RequestParaErr = 1000
	// RequestErr 请求方错误
	RequestErr = 1001
	// 无权限访问错误
	UnauthorizedErr = 1002
)

const (
	FailedMsg = "failed"
	SucMsg    = "success"
)

var (
	codeMsg = map[int]string{
		RequestParaErr:      "参数错误",
		UnauthorizedErr:     "无权限访问",
		ResponseInternalErr: "内部错误",
	}
)

// ResultJson 前后端分离开发时 约定好的标准的数据结构
type ResultJson struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// CodeToMsg code转换成 msg
func CodeToMsg(code int) (msg string, ok bool) {
	msg, ok = codeMsg[code]
	if ok {
		return
	}
	msg = FailedMsg
	return
}

// CodeJoinMsg code 对应的 信息和msg进行合拼
func CodeJoinMsg(code int, msg ...string) (join string) {
	temp, _ := CodeToMsg(code)
	if len(msg) == 0 {
		return temp
	}
	if len(msg[0]) == 0 {
		return temp
	}
	join = temp + ", " + msg[0]
	return

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
func OtherFailedJson(msg string, data ...interface{}) (result *ResultJson) {
	return commonErrJson(FailedCode, msg, data...)
}

// ResponseInternalErrJson 服务端内部错误
func ResponseInternalErrJson(msg string, data ...interface{}) (result *ResultJson) {
	return commonErrJson(ResponseInternalErr, CodeJoinMsg(ResponseInternalErr, msg), data...)
}

// RequestErrJson 请求端错误
func RequestErrJson(msg string, data ...interface{}) (result *ResultJson) {
	return commonErrJson(RequestErr, CodeJoinMsg(RequestErr, msg), data...)
}

// RequestParaErrJson 请求参数错误
func RequestParaErrJson(msg string, data ...interface{}) (result *ResultJson) {
	return commonErrJson(RequestParaErr, CodeJoinMsg(RequestParaErr, msg), data...)

}

// RequestParaErrJson 请求参数错误
func UnauthorizedErrJson(msg string, data ...interface{}) (result *ResultJson) {
	return commonErrJson(UnauthorizedErr, CodeJoinMsg(UnauthorizedErr, msg), data...)

}

func commonErrJson(code int, msg string, data ...interface{}) (result *ResultJson) {
	if len(data) < 1 {
		return &ResultJson{
			Code:    code,
			Message: msg,
			Data:    nil,
		}
	} else {
		return &ResultJson{
			Code:    code,
			Message: msg,
			Data:    data[0],
		}
	}

}

// CommonErrJson 通用的错误返回结构
func CommonErrJson(code int, msg string, data ...interface{}) (result *ResultJson) {
	return commonErrJson(code, msg, data)

}
