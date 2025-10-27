package errs

import "fmt"

// InternalError は予期しないシステムエラーや内部エラーを表します。
type InternalError struct {
	Code    string // エラーコード
	Message string // エラーメッセージ
	Cause   error  // 原因となったエラー（オプション）
}

// Error は error インターフェースを実装します。
func (e *InternalError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap は errors.Unwrap をサポートするためのメソッドです。
func (e *InternalError) Unwrap() error {
	return e.Cause
}

// NewInternalError は新しい内部エラーを生成します。
func NewInternalError(code, message string) *InternalError {
	return &InternalError{Code: code, Message: message}
}

// NewInternalErrorWithCause は原因となったエラーを含む内部エラーを生成します。
func NewInternalErrorWithCause(code, message string, cause error) *InternalError {
	return &InternalError{Code: code, Message: message, Cause: cause}
}
