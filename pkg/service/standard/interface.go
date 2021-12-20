package standard

import "github.com/lichuan0620/secret-keeper-backend/pkg/service/standard/common"

// Error defines interface that provide Error elements
type Error interface {
	GetHTTPCode() int32
	GetCode() string
	GetMessage() string
	GetData() map[string]string
}

// NewError returns an interface of base error
func NewError(httpCode int, code, message string, data map[string]string) Error {
	return &common.ErrorBase{
		HTTPCode: int32(httpCode),
		Code:     code,
		Message:  message,
		Data:     data,
	}
}
