package errs

import "fmt"

// ApplicationError はアプリケーションサービス層でのビジネスロジックエラーを表します。
type ApplicationError struct {
	Code    string // エラーコード
	Message string // エラーメッセージ
	Cause   error  // 原因となったエラー（オプション）
}

// Error は error インターフェースを実装します。
func (e *ApplicationError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap は errors.Unwrap をサポートするためのメソッドです。
func (e *ApplicationError) Unwrap() error {
	return e.Cause
}

// NewApplicationError は新しいアプリケーションエラーを生成します。
func NewApplicationError(code, message string) *ApplicationError {
	return &ApplicationError{Code: code, Message: message}
}

// NewApplicationErrorWithCause は原因となったエラーを含むアプリケーションエラーを生成します。
func NewApplicationErrorWithCause(code, message string, cause error) *ApplicationError {
	return &ApplicationError{Code: code, Message: message, Cause: cause}
}
