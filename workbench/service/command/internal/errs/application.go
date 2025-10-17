package errs

import "fmt"

// ApplicationError はアプリケーションサービス層でのビジネスロジックエラーを表します。
// 複数のドメインオブジェクトを跨ぐビジネスルールの違反や、
// アプリケーション固有のエラー条件を表現する際に使用されます。
//
// 主なエラーコード:
//   - CATEGORY_ALREADY_EXISTS: カテゴリが既に存在する
//   - OPERATION_FAILED: 操作が失敗した
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
//
// Parameters:
//   - code: エラーコード（例: "CATEGORY_ALREADY_EXISTS"）
//   - message: エラーメッセージ（例: "Category already exists"）
//
// Returns:
//   - *ApplicationError: 生成されたアプリケーションエラー
func NewApplicationError(code, message string) *ApplicationError {
	return &ApplicationError{Code: code, Message: message}
}

// NewApplicationErrorWithCause は原因となったエラーを含むアプリケーションエラーを生成します。
//
// Parameters:
//   - code: エラーコード
//   - message: エラーメッセージ
//   - cause: 原因となったエラー
//
// Returns:
//   - *ApplicationError: 生成されたアプリケーションエラー
func NewApplicationErrorWithCause(code, message string, cause error) *ApplicationError {
	return &ApplicationError{Code: code, Message: message, Cause: cause}
}
