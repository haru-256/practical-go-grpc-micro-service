package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/haru-256/practical-go-grpc-micro-service/service/client/internal/infrastructure/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewViper(t *testing.T) {
	t.Run("正常系: 設定ファイルを読み込める", func(t *testing.T) {
		// Arrange
		configContent := `
[server]
host = "localhost"
port = 8080

[command]
url = "http://localhost:8083"

[query]
url = "http://localhost:8085"
`
		tmpDir := t.TempDir()
		configPath := tmpDir
		configName := "test_config"
		configFile := filepath.Join(tmpDir, "test_config.toml")
		err := os.WriteFile(configFile, []byte(configContent), 0644)
		require.NoError(t, err, "設定ファイルの作成に失敗しました")

		// Act
		v := config.NewViper(configPath, configName)

		// Assert
		require.NotNil(t, v, "Viperインスタンスがnilです")
		assert.Equal(t, "localhost", v.GetString("server.host"))
		assert.Equal(t, 8080, v.GetInt("server.port"))
		assert.Equal(t, "http://localhost:8083", v.GetString("command.url"))
		assert.Equal(t, "http://localhost:8085", v.GetString("query.url"))
	})

	t.Run("正常系: 環境変数で設定を上書きできる", func(t *testing.T) {
		// Arrange
		configContent := `
[server]
host = "localhost"
port = 8080

[command]
url = "http://localhost:8083"

[query]
url = "http://localhost:8085"
`
		tmpDir := t.TempDir()
		configPath := tmpDir
		configName := "test_config"
		configFile := filepath.Join(tmpDir, "test_config.toml")
		err := os.WriteFile(configFile, []byte(configContent), 0644)
		require.NoError(t, err, "設定ファイルの作成に失敗しました")

		// 環境変数を設定
		t.Setenv("SERVER_HOST", "0.0.0.0")
		t.Setenv("SERVER_PORT", "9090")
		t.Setenv("COMMAND_URL", "http://command-service:8083")
		t.Setenv("QUERY_URL", "http://query-service:8085")

		// Act
		v := config.NewViper(configPath, configName)

		// Assert
		require.NotNil(t, v, "Viperインスタンスがnilです")
		assert.Equal(t, "0.0.0.0", v.GetString("server.host"))
		assert.Equal(t, "9090", v.GetString("server.port"))
		assert.Equal(t, "http://command-service:8083", v.GetString("command.url"))
		assert.Equal(t, "http://query-service:8085", v.GetString("query.url"))
	})

	t.Run("異常系: 設定ファイルが存在しない場合panicする", func(t *testing.T) {
		// Arrange
		tmpDir := t.TempDir()
		configPath := tmpDir
		configName := "non_existent_config"

		// Act & Assert
		assert.Panics(t, func() {
			config.NewViper(configPath, configName)
		}, "設定ファイルが存在しない場合、panicが発生するべきです")
	})

	t.Run("正常系: TOMLの入れ子構造を読み込める", func(t *testing.T) {
		// Arrange
		configContent := `
[server]
host = "localhost"
port = 8080

[server.tls]
enabled = true
cert_file = "/path/to/cert.pem"
key_file = "/path/to/key.pem"
`
		tmpDir := t.TempDir()
		configPath := tmpDir
		configName := "test_config"
		configFile := filepath.Join(tmpDir, "test_config.toml")
		err := os.WriteFile(configFile, []byte(configContent), 0644)
		require.NoError(t, err, "設定ファイルの作成に失敗しました")

		// Act
		v := config.NewViper(configPath, configName)

		// Assert
		require.NotNil(t, v, "Viperインスタンスがnilです")
		assert.Equal(t, "localhost", v.GetString("server.host"))
		assert.Equal(t, 8080, v.GetInt("server.port"))
		assert.True(t, v.GetBool("server.tls.enabled"))
		assert.Equal(t, "/path/to/cert.pem", v.GetString("server.tls.cert_file"))
		assert.Equal(t, "/path/to/key.pem", v.GetString("server.tls.key_file"))
	})

	t.Run("正常系: デフォルト値を設定できる", func(t *testing.T) {
		// Arrange
		configContent := `
[server]
host = "localhost"
`
		tmpDir := t.TempDir()
		configPath := tmpDir
		configName := "test_config"
		configFile := filepath.Join(tmpDir, "test_config.toml")
		err := os.WriteFile(configFile, []byte(configContent), 0644)
		require.NoError(t, err, "設定ファイルの作成に失敗しました")

		// Act
		v := config.NewViper(configPath, configName)

		// Assert
		require.NotNil(t, v, "Viperインスタンスがnilです")
		assert.Equal(t, "localhost", v.GetString("server.host"))
		// 設定ファイルに存在しないキーは、Viperのデフォルト動作でゼロ値を返す
		assert.Equal(t, 0, v.GetInt("server.port"))
		assert.Empty(t, v.GetString("command.url"))
	})

	t.Run("正常系: 環境変数のプレフィックスが機能する", func(t *testing.T) {
		// Arrange
		configContent := `
[server]
host = "localhost"
port = 8080
`
		tmpDir := t.TempDir()
		configPath := tmpDir
		configName := "test_config"
		configFile := filepath.Join(tmpDir, "test_config.toml")
		err := os.WriteFile(configFile, []byte(configContent), 0644)
		require.NoError(t, err, "設定ファイルの作成に失敗しました")

		// 環境変数を設定（プレフィックスなし - 動作確認用）
		t.Setenv("SERVER_HOST", "env-host")

		// Act
		v := config.NewViper(configPath, configName)

		// Assert
		require.NotNil(t, v, "Viperインスタンスがnilです")
		// 環境変数の値で上書きされる
		assert.Equal(t, "env-host", v.GetString("server.host"))
		assert.Equal(t, 8080, v.GetInt("server.port"))
	})
}
