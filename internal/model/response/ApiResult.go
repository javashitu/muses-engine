package response

import "muses-engine/internal/common"

type ApiResult struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func Ok() *ApiResult {
	return constructResp(common.NoneError.Code, common.NoneError.Message)
}

func SuccessApi(data interface{}) *ApiResult {
	result := Ok()
	result.Data = data
	return result
}

func FailureApi(errorDescPair common.ErrorDescPair) *ApiResult {
	return constructResp(errorDescPair.Code, errorDescPair.Message)
}

func constructResp(code string, messaage string) *ApiResult {
	result := new(ApiResult)
	result.Code = code
	result.Message = messaage
	return result
}
