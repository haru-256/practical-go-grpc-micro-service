package errs

import "fmt"

// InternalError は予期しないシステムエラーや内部エラーを表します。
// データベース接続エラー、設定ファイルの読み込みエラーなど、
// アプリケーションの正常な動作を妨げる技術的なエラーに使用されます。
//
// 主なエラーコード:
//   - DB_CONNECTION_ERROR: データベース接続エラー
//   - CONFIG_ERROR: 設定エラー
//   - INTERNAL_ERROR: その他の内部エラー
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
//
// Parameters:
//   - code: エラーコード（例: "DB_CONNECTION_ERROR"）
//   - message: エラーメッセージ
//
// Returns:
//   - *InternalError: 生成された内部エラー
func NewInternalError(code, message string) *InternalError {
	return &InternalError{Code: code, Message: message}
}

// NewInternalErrorWithCause は原因となったエラーを含む内部エラーを生成します。
//
// Parameters:
//   - code: エラーコード
//   - message: エラーメッセージ
//   - cause: 原因となったエラー
//
// Returns:
//   - *InternalError: 生成された内部エラー
func NewInternalErrorWithCause(code, message string, cause error) *InternalError {
	return &InternalError{Code: code, Message: message, Cause: cause}
}
