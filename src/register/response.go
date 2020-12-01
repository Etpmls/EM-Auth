package register

import (
	"encoding/json"
	"github.com/Etpmls/Etpmls-Micro"
	"github.com/Etpmls/Etpmls-Micro/library"
	em_utils "github.com/Etpmls/Etpmls-Micro/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
)

// Return error message in json format
// 返回json格式的错误信息
func MakeRpcError(rcpStatusCode codes.Code, code string, message string, data interface{}, err error) (*em.Response, error) {
	// Data interface => Json string
	var tmp_data []byte
	if data != nil {
		// First judge whether it is a string, if it is a string, then judge whether it belongs to json, if it already belongs to json, no longer convert json again
		// 首先判断是不是字符串，如果是字符串，则判断属不属于json，如果已经属于json，则不再二次转化json
		v, ok := data.(string)
		if ok {
			ok_json := json.Valid([]byte(v))
			if ok_json {
				tmp_data = []byte(v)
			} else {
				tmp_data, _ = json.Marshal(data)
			}
		} else {
			tmp_data, _ = json.Marshal(data)
		}
	}

	// If enabled, use HTTP CODE instead of system default CODE
	// 如果开启使用HTTP CODE 代替系统的默认CODE
	if em_library.Config.App.UseHttpCode == true {
		code = strconv.Itoa(int(rcpStatusCode))
	}

	// If it is a Debug environment, return information with Error
	// 如果是Debug环境，返回带有Error的信息
	if em.IsDebug() && err != nil {
		err := status.Error(rcpStatusCode, em.Response{
			Code:    code,
			Status:  em.ERROR_Status,
			Message: message + "Error: " + err.Error(),
			Data:    string(tmp_data),
		}.String())
		return nil, err
	}

	err3 := status.Error(rcpStatusCode, em.Response{
		Code:    code,
		Status:  em.ERROR_Status,
		Message: message,
		Data:    string(tmp_data),
	}.String())
	return nil, err3
}

// Return success information in json format
// 返回json格式的成功信息
func MakeRpcSuccess(code string, message string, data interface{}) (*em.Response, error) {
	// First judge whether it is a string, if it is a string, then judge whether it belongs to json, if it already belongs to json, no longer convert json again
	// 首先判断是不是字符串，如果是字符串，则判断属不属于json，如果已经属于json，则不再二次转化json
	if v, ok := data.(string); ok {
		ok_json := json.Valid([]byte(v))
		if ok_json {
			return &em.Response{
				Code:    code,
				Status:  em.SUCCESS_Status,
				Message: message,
				Data:    v,
			}, nil
		}
	}

	return &em.Response{
		Code:    code,
		Status:  em.SUCCESS_Status,
		Message: message,
		Data:    em_utils.MustConvertJson(data),
	}, nil
}

