package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/haru-256/practical-go-grpc-micro-service/service/query/internal/infrastructure/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewViper(t *testing.T) {
	tests := []struct {
		name          string
		configContent string
		envVars       map[string]string
		wantPanic     bool
		assertions    func(t *testing.T, configPath, configName string)
	}{
		{
			name: "正常系: 設定ファイルを読み込める",
			configContent: `
[server]
port = 50051

[mysql]
dbname = "test_db"
host = "localhost"
port = 3306
user = "test_user"
pass = "test_pass"
max_idle_conns = 10
max_open_conns = 100
conn_max_lifetime = 300
conn_max_idle_time = 60
`,
			envVars:   map[string]string{},
			wantPanic: false,
			assertions: func(t *testing.T, configPath, configName string) {
				v := config.NewViper(configPath, configName)
				assert.NotNil(t, v)
				assert.Equal(t, 50051, v.GetInt("server.port"))
				assert.Equal(t, "test_db", v.GetString("mysql.dbname"))
				assert.Equal(t, "localhost", v.GetString("mysql.host"))
				assert.Equal(t, 3306, v.GetInt("mysql.port"))
				assert.Equal(t, "test_user", v.GetString("mysql.user"))
				assert.Equal(t, "test_pass", v.GetString("mysql.pass"))
				assert.Equal(t, 10, v.GetInt("mysql.max_idle_conns"))
				assert.Equal(t, 100, v.GetInt("mysql.max_open_conns"))
				assert.Equal(t, 300, v.GetInt("mysql.conn_max_lifetime"))
				assert.Equal(t, 60, v.GetInt("mysql.conn_max_idle_time"))
			},
		},
		{
			name: "正常系: 環境変数で設定を上書きできる",
			configContent: `
[mysql]
dbname = "test_db"
host = "localhost"
port = 3306
`,
			envVars: map[string]string{
				"DB_DBNAME": "env_db",
				"DB_HOST":   "env_host",
				"DB_PORT":   "3307",
			},
			wantPanic: false,
			assertions: func(t *testing.T, configPath, configName string) {
				v := config.NewViper(configPath, configName)
				assert.Equal(t, "env_db", v.GetString("mysql.dbname"))
				assert.Equal(t, "env_host", v.GetString("mysql.host"))
				assert.Equal(t, "3307", v.GetString("mysql.port"))
			},
		},
		{
			name:          "異常系: 設定ファイルが存在しない場合panicする",
			configContent: "",
			envVars:       map[string]string{},
			wantPanic:     true,
			assertions: func(t *testing.T, configPath, configName string) {
				assert.Panics(t, func() {
					config.NewViper(configPath, "non_existent_config")
				})
			},
		},
		{
			name: "正常系: ドット区切りのキーを環境変数のアンダースコア区切りに変換できる",
			configContent: `
[mysql]
max_idle_conns = 10
`,
			envVars: map[string]string{
				"DB_MAX_IDLE_CONNS": "20",
			},
			wantPanic: false,
			assertions: func(t *testing.T, configPath, configName string) {
				v := config.NewViper(configPath, configName)
				assert.Equal(t, "20", v.GetString("mysql.max_idle_conns"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// テスト用の設定ファイルを作成
			tmpDir := t.TempDir()
			configPath := tmpDir
			configName := "test_config"

			if tt.configContent != "" {
				configFile := filepath.Join(tmpDir, "test_config.toml")
				err := os.WriteFile(configFile, []byte(tt.configContent), 0644)
				require.NoError(t, err)
			}

			// 環境変数を設定
			for key, value := range tt.envVars {
				t.Setenv(key, value)
			}

			// アサーションを実行
			tt.assertions(t, configPath, configName)
		})
	}
}
