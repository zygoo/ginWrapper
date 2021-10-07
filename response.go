package ginWrapper

import "net/http"

type ResponseInfo struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Response struct {
	Result ResponseInfo `json:"result"`
	Data   interface{}  `json:"data"`
}

const (
	CodeSuccess = 0
	CodeFailure = 1

	CodeBindError        = 10
	CodeRecordNotFound   = 11
	CodePermissionDenied = 12
)

var ResponseMap = map[int]ResponseInfo{
	CodeSuccess:   {http.StatusOK, http.StatusText(http.StatusOK)},
	CodeFailure:   {http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)},
	CodeBindError: {http.StatusInternalServerError, "params bind error"},
}

// SendResponse code为可变参数，code[0]为business code, code[1]为http code
func (c *GinContext) SendResponse(data interface{}, codes ...int) {
	httpCode := http.StatusOK
	response, exist := ResponseMap[codes[0]]
	if !exist {
		response = ResponseMap[CodeFailure]
	}

	result := ResponseInfo{codes[0], response.Message}
	if len(codes) > 1 {
		httpCode = codes[1]
	} else if result.Code <= 20 {
		httpCode = result.Code
	}
	c.JSON(httpCode, Response{result, data})
}

func (c *GinContext) OK(data interface{}) {
	c.SendResponse(data, CodeSuccess)
}

func (c *GinContext) BindFailure() {
	c.SendResponse(nil, CodeBindError)
}

func (c *GinContext) Failure(err error) {
	c.JSON(http.StatusInternalServerError, Response{
		Result: ResponseInfo{
			Code:    CodeFailure,
			Message: err.Error(),
		},
		Data: nil,
	})
}
