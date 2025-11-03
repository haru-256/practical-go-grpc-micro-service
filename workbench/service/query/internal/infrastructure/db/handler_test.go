//go:build integration || !ci

package db_test

import (
	"errors"
	"net"
	"testing"
	"time"

	mysql_go "github.com/go-sql-driver/mysql"
	"github.com/haru-256/practical-go-grpc-micro-service/service/query/internal/infrastructure/config"
	"github.com/haru-256/practical-go-grpc-micro-service/service/query/internal/infrastructure/db"
	"github.com/haru-256/practical-go-grpc-micro-service/service/query/internal/testhelpers"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDBConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setupViper func() *viper.Viper
		wantErr    bool
		assertions func(t *testing.T, cfg *db.DBConfig, err error)
	}{
		{
			name: "正常系: すべての設定値を読み込める",
			setupViper: func() *viper.Viper {
				v := viper.New()
				v.Set("mysql.dbname", "sample_db")
				v.Set("mysql.host", "localhost")
				v.Set("mysql.port", 3307)
				v.Set("mysql.user", "root")
				v.Set("mysql.pass", "password")
				v.Set("mysql.max_idle_conns", 10)
				v.Set("mysql.max_open_conns", 100)
				v.Set("mysql.conn_max_lifetime", 1800*time.Second)
				v.Set("mysql.conn_max_idle_time", 500*time.Second)
				v.Set("log.level", "info")
				return v
			},
			wantErr: false,
			assertions: func(t *testing.T, cfg *db.DBConfig, err error) {
				require.NoError(t, err)
				assert.Equal(t, "sample_db", cfg.DBName)
				assert.Equal(t, "localhost", cfg.Host)
				assert.Equal(t, 3307, cfg.Port)
				assert.Equal(t, "root", cfg.User)
				assert.Equal(t, "password", cfg.Pass)
				assert.Equal(t, 10, cfg.MaxIdleConns)
				assert.Equal(t, 100, cfg.MaxOpenConns)
				assert.Equal(t, 1800*time.Second, cfg.ConnMaxLifetime)
				assert.Equal(t, 500*time.Second, cfg.ConnMaxIdleTime)
				assert.Equal(t, "info", cfg.LogLevel)
			},
		},
		{
			name: "異常系: 必須設定値が欠けている場合エラーを返す",
			setupViper: func() *viper.Viper {
				v := viper.New()
				// mysql.dbnameとmysql.hostのみ設定
				v.Set("mysql.dbname", "test_db")
				v.Set("mysql.host", "localhost")
				// その他の必須項目は未設定
				return v
			},
			wantErr: true,
			assertions: func(t *testing.T, cfg *db.DBConfig, err error) {
				require.Error(t, err)
				assert.NotNil(t, cfg)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			v := tt.setupViper()
			cfg, err := db.NewDBConfig(v)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			tt.assertions(t, cfg, err)
		})
	}
}

func TestNewDatabase(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setupFunc  func() *db.DBConfig
		wantErr    bool
		assertions func(t *testing.T, conn interface{}, err error)
	}{
		{
			name: "正常系: データベースに接続できる",
			setupFunc: func() *db.DBConfig {
				configPath := "../../../"
				configName := "config"
				v := config.NewViper(configPath, configName)
				cfg, err := db.NewDBConfig(v)
				require.NoError(t, err)
				return cfg
			},
			wantErr: false,
			assertions: func(t *testing.T, conn interface{}, err error) {
				require.NoError(t, err)
				assert.NotNil(t, conn)
			},
		},
		{
			name: "異常系: 不正なログレベルの場合エラーを返す",
			setupFunc: func() *db.DBConfig {
				configPath := "../../../"
				configName := "config"
				v := config.NewViper(configPath, configName)
				// 不正なログレベルに上書き
				v.Set("log.level", "invalid_level")
				cfg, err := db.NewDBConfig(v)
				require.NoError(t, err)
				return cfg
			},
			wantErr: true,
			assertions: func(t *testing.T, conn interface{}, err error) {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "不正なログレベルです")
			},
		},
		{
			name: "異常系: 接続できないホストの場合エラーを返す",
			setupFunc: func() *db.DBConfig {
				configPath := "../../../"
				configName := "config"
				v := config.NewViper(configPath, configName)
				// 無効なホストに上書き
				v.Set("mysql.host", "invalid_host_12345")
				cfg, err := db.NewDBConfig(v)
				require.NoError(t, err)
				return cfg
			},
			wantErr: true,
			assertions: func(t *testing.T, conn interface{}, err error) {
				require.Error(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cfg := tt.setupFunc()
			conn, err := db.NewDatabase(cfg, testhelpers.TestLogger)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if tt.assertions != nil {
				tt.assertions(t, conn, err)
			}
		})
	}
}

func TestDBErrHandler(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		inputErr   error
		assertions func(t *testing.T, err error)
	}{
		{
			name:     "ネットワークエラーの場合DB_CONNECTION_ERRORを返す",
			inputErr: &net.OpError{Op: "dial", Net: "tcp", Err: errors.New("connection timeout")},
			assertions: func(t *testing.T, err error) {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "DB_CONNECTION_ERROR")
			},
		},
		{
			name:     "MySQLドライバエラーの場合DB_DRIVER_ERRORを返す",
			inputErr: &mysql_go.MySQLError{Number: 1062, Message: "Duplicate entry"},
			assertions: func(t *testing.T, err error) {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "DB_DRIVER_ERROR")
			},
		},
		{
			name:     "その他のエラーの場合UNKNOWN_ERRORを返す",
			inputErr: errors.New("unknown error"),
			assertions: func(t *testing.T, err error) {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "UNKNOWN_ERROR")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := db.DBErrHandler(t.Context(), tt.inputErr, testhelpers.TestLogger)
			tt.assertions(t, err)
		})
	}
}
