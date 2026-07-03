package response

type Response struct {
	Status  string `json:"status"`
	Data    any    `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
	Code    int    `json:"code,omitempty"`
}

func Success(data any) Response {
	return Response{
		Status: "success",
		Data:   data,
	}
}

func Fail(data any) Response {
	return Response{
		Status: "fail",
		Data:   data,
	}
}

func Error(message string) Response {
	return Response{
		Status:  "error",
		Message: message,
	}
}

func ErrorCode(code int, message string, data any) Response {
	return Response{
		Status:  "error",
		Message: message,
		Code:    code,
		Data:    data,
	}
}
