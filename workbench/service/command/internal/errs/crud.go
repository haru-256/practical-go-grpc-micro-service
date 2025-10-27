package errs

import "fmt"

// CRUDError はデータベースのCRUD操作に関連するエラーを表します。
type CRUDError struct {
	Code    string // エラーコード
	Message string // エラーメッセージ
	Cause   error  // 原因となったエラー
}

// Error は error インターフェースを実装します。
func (e *CRUDError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap は errors.Unwrap をサポートするためのメソッドです。
func (e *CRUDError) Unwrap() error {
	return e.Cause
}

// NewCRUDError は新しいCRUDエラーを生成します。
func NewCRUDError(code, message string) *CRUDError {
	return &CRUDError{Code: code, Message: message}
}

// NewCRUDErrorWithCause は原因となったエラーを含むCRUDエラーを生成します。
func NewCRUDErrorWithCause(code, message string, cause error) *CRUDError {
	return &CRUDError{Code: code, Message: message, Cause: cause}
}
