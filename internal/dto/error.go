package dto

const (
	ErrCodeNotFound     = 404
	ErrCodeBadRequest   = 400
	ErrCodeUnauthorized = 401
	ErrCodeConflict     = 409
)

type (
	AppError struct {
		Message string
		Code    int
	}
)

func (e AppError) Error() string {
	return e.Message
}
