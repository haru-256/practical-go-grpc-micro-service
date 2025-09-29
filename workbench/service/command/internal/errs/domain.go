package errs

import "fmt"

// ドメインに関連するエラー型
// スタイルガイド (4.2.2) に基づき、構造化されたエラー情報を保持します。
type DomainError struct {
	Code    string
	Message string
	Cause   error
}

// Error は error インターフェースを実装します。
func (e *DomainError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// NewDomainError はドメインエラーを生成します。
// 例: errs.NewDomainError("INVALID_ARGUMENT", "商品IDはUUIDの形式である必要があります")
func NewDomainError(code, message string) *DomainError {
	return &DomainError{Code: code, Message: message}
}

// NewDomainErrorWithCause は原因となったエラーを含むドメインエラーを生成します。
func NewDomainErrorWithCause(code, message string, cause error) *DomainError {
	return &DomainError{Code: code, Message: message, Cause: cause}
}
