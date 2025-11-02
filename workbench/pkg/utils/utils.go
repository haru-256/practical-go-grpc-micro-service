package utils

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// getKey はViperから型安全に設定値を取得するヘルパー関数です。
// 指定されたキーが存在しない場合、またはサポートされていない型の場合はエラーを記録します。
//
// サポートされる型:
//   - string
//   - int
//   - bool
//   - time.Duration
//
// Parameters:
//   - v: Viperインスタンス
//   - key: 設定キー（例: "mysql.host"）
//   - errs: エラーを蓄積するスライスへのポインタ
//
// Returns:
//   - T: 設定値（エラーの場合はゼロ値）
func GetKey[T any](v *viper.Viper, key string, errs *[]error) T {
	var zero T
	if !v.IsSet(key) {
		*errs = append(*errs, fmt.Errorf("config key '%s' is not set", key))
		return zero
	}

	switch any(zero).(type) {
	case string:
		return any(v.GetString(key)).(T)
	case int:
		return any(v.GetInt(key)).(T)
	case bool:
		return any(v.GetBool(key)).(T)
	case time.Duration:
		// v.GetDuration() を使うことで、"30m" や "1h" のような文字列を
		// time.Duration型へ安全にパースする処理をViperに任せます。
		return any(v.GetDuration(key)).(T)
	default:
		*errs = append(*errs, fmt.Errorf("unsupported type for key '%s'", key))
		return zero
	}
}
