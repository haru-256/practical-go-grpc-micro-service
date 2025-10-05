package errs

import "fmt"

// 内部エラー型(データベース接続エラーなど)
type InternalError struct {
	Code    string
	Message string // エラーメッセージ
	Cause   error
}

// エラーメッセージを返すメソッド
func (e *InternalError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *InternalError) Unwrap() error {
	return e.Cause
}

// コンストラクタ
func NewInternalError(code, message string) *InternalError {
	return &InternalError{Code: code, Message: message}
}

// NewInternalErrorWithCause は原因となったエラーを含む内部エラーを生成します。
func NewInternalErrorWithCause(code, message string, cause error) *InternalError {
	return &InternalError{Code: code, Message: message, Cause: cause}
}
