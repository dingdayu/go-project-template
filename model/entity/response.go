package entity

var (
	SuccessResponse = Response{Code: 0, Message: "ok"} //success default

	// 系统响应    100xxx
	ErrInternalServer = Response{Code: 10001, Message: "系统错误"}
	ErrMissParams     = Response{Code: 10002, Message: "缺少参数"}
	ErrFailParams     = Response{Code: 10003, Message: "参数格式错误"}
	ErrNotExist       = Response{Code: 10004, Message: "数据不存在"}
	ErrDefault        = Response{Code: 10005, Message: "操作失败"}
	ErrDataPermission = Response{Code: 10006, Message: "没有此数据权限"}
	ErrDuplicatedKey  = Response{Code: 10007, Message: "无法创建重复数据"}

	// Auth 认证登陆 响应 101xx
	ErrAuthForbidden    = Response{Code: 10100, Message: "登录信息不存"}
	ErrAuthUserNotFound = Response{Code: 10101, Message: "登陆用户不存在"}
	// Token 失效，需要重新登录
	ErrAuthTokenInvalid = Response{Code: 10102, Message: "Token 失效，需要重新登录"}
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

func NewErrResponse(err any) Response {
	s := ErrInternalServer
	s.Error = err
	return s
}

func NewSucResponse(data interface{}) Response {
	s := SuccessResponse
	s.Data = data
	return s
}
func NewErrParamsResponse(message string) Response {
	s := ErrMissParams
	s.Error = message
	return s
}
