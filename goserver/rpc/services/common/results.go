package common

var (
	ResultOK = &CodeAndMessage{
		Code:    1,
		Message: "ok",
	}
)

type CodeAndMessage struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
