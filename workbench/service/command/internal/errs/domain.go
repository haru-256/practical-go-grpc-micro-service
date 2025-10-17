package errs

import "fmt"

// DomainError はドメインモデルのバリデーションやビジネスルール違反を表すエラーです。
// 値オブジェクトやエンティティの生成時、ドメインルールの検証時に使用されます。
//
// 主なエラーコード:
//   - INVALID_ARGUMENT: 引数の値が不正（例: 文字数制約違反、形式不正）
//   - BUSINESS_RULE_VIOLATION: ビジネスルール違反
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
//
// Parameters:
//   - code: エラーコード（例: "INVALID_ARGUMENT"）
//   - message: エラーメッセージ（例: "商品IDはUUIDの形式である必要があります"）
//
// Returns:
//   - *DomainError: 生成されたドメインエラー
//
// Example:
//
//	errs.NewDomainError("INVALID_ARGUMENT", "カテゴリ名は1文字以上20文字以下で入力してください")
func NewDomainError(code, message string) *DomainError {
	return &DomainError{Code: code, Message: message}
}

// NewDomainErrorWithCause は原因となったエラーを含むドメインエラーを生成します。
//
// Parameters:
//   - code: エラーコード
//   - message: エラーメッセージ
//   - cause: 原因となったエラー
//
// Returns:
//   - *DomainError: 生成されたドメインエラー
func NewDomainErrorWithCause(code, message string, cause error) *DomainError {
	return &DomainError{Code: code, Message: message, Cause: cause}
}
