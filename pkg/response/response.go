package response

import "time"

type Response struct {
	Status  string    `json:"status"`
	Message string    `json:"mssage"`
	Date    time.Time `json:"time"`
}

func Success(mes string) Response {
	return Response{
		Status:  "Success",
		Message: mes,
		Date:    time.Now(),
	}
}

func Error(mes string) Response {
	return Response{
		Status:  "Error",
		Message: mes,
		Date:    time.Now(),
	}
}
