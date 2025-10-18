package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// NewViper は設定ファイルと環境変数を読み込むViperインスタンスを生成します。
// TOMLファイルから基本設定を読み込み、環境変数で上書き可能にします。
//
// 設定ファイル:
//   - 形式: TOML
//   - 場所: configPathディレクトリ内のconfigName.toml
//
// 環境変数のマッピング:
//   - 自動マッピング: "." を "_" に変換（例: log.level → LOG_LEVEL）
//   - 明示的バインディング: mysql.* → DB_*（例: mysql.dbname → DB_DBNAME）
//
// Parameters:
//   - configPath: 設定ファイルのディレクトリパス
//   - configName: 設定ファイル名（拡張子なし）
//
// Returns:
//   - *viper.Viper: 設定が読み込まれたViperインスタンス
//
// Panics:
//   - 環境変数のバインディングに失敗した場合
//   - 設定ファイルの読み込みに失敗した場合
//
// 使用例:
//
//	v := config.NewViper(".", "config")
//	dbName := v.GetString("mysql.dbname")  // TOMLまたはDB_DBNAME環境変数から取得
func NewViper(configPath string, configName string) *viper.Viper {
	v := viper.New()
	v.AddConfigPath(configPath)
	v.SetConfigName(configName)
	v.SetConfigType("toml")

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
		if err := v.BindEnv(key, env); err != nil {
			panic(fmt.Sprintf("failed to bind env for %s: %v", key, err))
		}
	}

	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	return v
}
