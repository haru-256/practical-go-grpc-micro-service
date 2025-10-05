package utils

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// EnvType は環境変数から取得可能な型を制約するインターフェースです。
// サポートされる型: string, int, int64, bool, float64, time.Duration
// 注意: time.Durationは内部的にint64のエイリアスのため、型制約には含まれていませんが、
// 実行時に適切に処理されます。
type EnvType interface {
	~string | ~int | ~int64 | ~bool | ~float64
}

// GetEnv は環境変数から値を取得し、指定された型に変換して返します。
//
// 環境変数が設定されていない場合は、defaultValueを返します。
// 環境変数の値が指定された型に変換できない場合は、panicします。
//
// サポートされる型:
//   - string: そのまま文字列として返す
//   - int: strconv.Atoiで整数に変換
//   - bool: strconv.ParseBoolでブール値に変換（true, false, 1, 0など）
//   - float64: strconv.ParseFloatで浮動小数点数に変換
//   - time.Duration: time.ParseDurationで時間に変換（"300ms", "1.5h", "2h45m"など）
//
// 使用例:
//
//	// 文字列型
//	host := GetEnv("DB_HOST", "localhost")
//
//	// 整数型
//	port := GetEnv("DB_PORT", 3306)
//
//	// ブール型
//	debug := GetEnv("DEBUG", false)
//
//	// 浮動小数点型
//	timeoutSec := GetEnv("TIMEOUT_SEC", 30.0)
//
//	// 時間型
//	timeout := GetEnv("TIMEOUT", 30*time.Second)
//
// Parameters:
//   - key: 環境変数のキー名
//   - defaultValue: 環境変数が設定されていない場合のデフォルト値
//
// Returns:
//   - T: 環境変数の値、または環境変数が未設定の場合はdefaultValue
//
// Panics:
//   - 環境変数の値が指定された型に変換できない場合
func GetEnv[T EnvType](key string, defaultValue T) (T, error) {
	valStr, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue, nil
	}

	var result any
	switch any(defaultValue).(type) {
	case string:
		result = valStr
	case int:
		v, err := strconv.Atoi(valStr)
		if err != nil {
			return defaultValue, err
		}
		result = v
	case int64:
		v, err := strconv.ParseInt(valStr, 10, 64)
		if err != nil {
			return defaultValue, err
		}
		result = v
	case bool:
		v, err := strconv.ParseBool(valStr)
		if err != nil {
			return defaultValue, err
		}
		result = v
	case float64:
		v, err := strconv.ParseFloat(valStr, 64)
		if err != nil {
			return defaultValue, err
		}
		result = v
	case time.Duration:
		// ここで文字列をパースして time.Duration に変換
		v, err := time.ParseDuration(valStr)
		if err != nil {
			return defaultValue, err
		}
		result = v
	default:
		return defaultValue, fmt.Errorf("GetEnvでサポートされていない型が指定されました: %T", defaultValue)
	}

	return result.(T), nil
}
