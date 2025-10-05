package errs

import "fmt"

// データベースアクセスエラー型
type CRUDError struct {
	Code    string
	Message string // エラーメッセージ
	Cause   error
}

// エラーメッセージを返すメソッド
func (e *CRUDError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *CRUDError) Unwrap() error {
	return e.Cause
}

// コンストラクタ
func NewCRUDError(code, message string) *CRUDError {
	return &CRUDError{Code: code, Message: message}
}

// NewCRUDErrorWithCause は原因となったエラーを含むCRUDエラーを生成します。
func NewCRUDErrorWithCause(code, message string, cause error) *CRUDError {
	return &CRUDError{Code: code, Message: message, Cause: cause}
}
