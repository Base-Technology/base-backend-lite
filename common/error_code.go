package common

const (
	ErrorPanic = 1 + iota
	ErrorInvalidToken
	ErrorInvalidParams
	ErrorPassword
	ErrorInner
	ErrorLimitExceedeed
)

type BaseResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
