package config

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// NewViper は設定ファイルと環境変数を読み込むViperインスタンスを生成します。
//
// Parameters:
//   - configPath: 設定ファイルのディレクトリパス
//   - configName: 設定ファイル名（拡張子なし）
//
// Returns:
//   - *viper.Viper: 設定が読み込まれたViperインスタンス
//
// Panics:
//   - 設定の読み込みに失敗した場合
func NewViper(configPath string, configName string) *viper.Viper {
	v := viper.New()
	v.AddConfigPath(configPath)
	v.SetConfigName(configName)
	v.SetConfigType("toml")

	absolutePath, err := filepath.Abs(v.ConfigFileUsed())
	if err != nil {
		panic(err)
	}

	v.AutomaticEnv()

	// . -> _ に変換して環境変数を読み込む
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 特殊なバインディング
	bindings := map[string]string{
		"mysql.dbname":             "DB_DBNAME",
		"mysql.host":               "DB_HOST",
		"mysql.port":               "DB_PORT",
		"mysql.user":               "DB_USER",
		"mysql.pass":               "DB_PASS",
		"mysql.max_idle_conns":     "DB_MAX_IDLE_CONNS",
		"mysql.max_open_conns":     "DB_MAX_OPEN_CONNS",
		"mysql.conn_max_lifetime":  "DB_CONN_MAX_LIFETIME",
		"mysql.conn_max_idle_time": "DB_CONN_MAX_IDLE_TIME",
	}
	for key, env := range bindings {
		if bindErr := v.BindEnv(key, env); bindErr != nil {
			panic(fmt.Sprintf("failed to bind env for %s: %v", key, bindErr))
		}
	}

	if readErr := v.ReadInConfig(); readErr != nil {
		panic(fmt.Errorf("failed to read config file %s: %w", absolutePath, readErr))
	}

	return v
}
