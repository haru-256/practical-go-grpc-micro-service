package errs

import "fmt"

// CRUDError はデータベースのCRUD操作に関連するエラーを表します。
// リポジトリ層からアプリケーション層へエラー情報を伝達する際に使用されます。
//
// 主なエラーコード:
//   - NOT_FOUND: レコードが見つからない
//   - DB_UNIQUE_CONSTRAINT_VIOLATION: 一意制約違反
//   - DB_ERROR: その他のデータベースエラー
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
//
// Parameters:
//   - code: エラーコード（例: "NOT_FOUND", "DB_UNIQUE_CONSTRAINT_VIOLATION"）
//   - message: エラーメッセージ
//
// Returns:
//   - *CRUDError: 生成されたCRUDエラー
func NewCRUDError(code, message string) *CRUDError {
	return &CRUDError{Code: code, Message: message}
}

// NewCRUDErrorWithCause は原因となったエラーを含むCRUDエラーを生成します。
//
// Parameters:
//   - code: エラーコード
//   - message: エラーメッセージ
//   - cause: 原因となったエラー
//
// Returns:
//   - *CRUDError: 生成されたCRUDエラー
func NewCRUDErrorWithCause(code, message string, cause error) *CRUDError {
	return &CRUDError{Code: code, Message: message, Cause: cause}
}
