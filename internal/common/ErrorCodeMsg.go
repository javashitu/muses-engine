package common

// 规则xx-yyy-zzzz
// XX 系统 yyy 模块 zzzz错误码
//还是有可能透出来使用的，比如错误处理
const (
	CodeOK            = "0"
	UploadFileFailure = "301000015"
	ParamWrong        = "301000020"
	BindParamFailure  = "301000025"
	PreSignUrlFailure = "301000030"
	Unauthenticated   = "301000035"
	TokenExpire       = "301000040"
)

var (
	NoneError            = newErrorDescPair(CodeOK, "success")
	ParamError           = newErrorDescPair(ParamWrong, "param wrong")
	BindParamError       = newErrorDescPair(BindParamFailure, "bind param failure")
	UploadFileError      = newErrorDescPair(UploadFileFailure, "upload file failure")
	PreSignUrlError      = newErrorDescPair(PreSignUrlFailure, "pre sign visit url failure")
	UnauthenticatedError = newErrorDescPair(Unauthenticated, "authentication failure")
	TokenExpireError     = newErrorDescPair(TokenExpire, "token expire failure")
)

type ErrorDescPair struct {
	Code    string
	Message string
}

func newErrorDescPair(code string, message string) ErrorDescPair {
	return ErrorDescPair{
		Code:    code,
		Message: message,
	}
}
