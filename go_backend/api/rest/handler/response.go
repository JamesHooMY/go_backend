package handler

var (
	ErrRequestInvalid = 40001
	ErrUnauthorized   = 40101
	ErrForbidden      = 40301
	ErrNotFound       = 40401
	ErrInternalServer = 50001

	ErrInternalServerMsg = "internal server error"
)

type Response struct {
	Code int    `json:"code,omitempty"`
	Msg  string `json:"msg,omitempty"`
	Data any    `json:"data,omitempty"`
}
