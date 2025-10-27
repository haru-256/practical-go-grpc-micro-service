package errs

import "fmt"

// DomainError はドメインモデルのバリデーションやビジネスルール違反を表すエラーです。
type DomainError struct {
	Code    string // エラーコード
	Message string // エラーメッセージ
	Cause   error  // 原因となったエラー（オプション）
}

// Error は error インターフェースを実装します。
func (e *DomainError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap は errors.Unwrap をサポートするためのメソッドです。
func (e *DomainError) Unwrap() error {
	return e.Cause
}

// NewDomainError は新しいドメインエラーを生成します。
func NewDomainError(code, message string) *DomainError {
	return &DomainError{Code: code, Message: message}
}

// NewDomainErrorWithCause は原因となったエラーを含むドメインエラーを生成します。
func NewDomainErrorWithCause(code, message string, cause error) *DomainError {
	return &DomainError{Code: code, Message: message, Cause: cause}
}
