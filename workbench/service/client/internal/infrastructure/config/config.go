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
	if readErr := v.ReadInConfig(); readErr != nil {
		panic(fmt.Errorf("failed to read config file %s: %w", absolutePath, readErr))
	}
	// 環境変数の読み込み設定
	v.AutomaticEnv()
	// . -> _ に変換して環境変数を読み込む
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	return v
}
