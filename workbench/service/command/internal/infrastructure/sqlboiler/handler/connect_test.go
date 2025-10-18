//go:build integration || !ci

package handler

import (
	"os"
	"testing"
	"time"

	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/config"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestConn(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "infrastructure/sqlboiler/handler packageのテスト")
}

// testEnvVars はテストで使用される環境変数のリストです。
var testEnvVars = []string{
	"LOG_LEVEL",
	"DB_DBNAME",
	"DB_HOST",
	"DB_PORT",
	"DB_USER",
	"DB_PASS",
	"DB_MAX_IDLE_CONNS",
	"DB_MAX_OPEN_CONNS",
	"DB_CONN_MAX_LIFETIME",
	"DB_CONN_MAX_IDLE_TIME",
}

// backupEnvVars は指定された環境変数のバックアップを作成します。
//
// Parameters:
//   - envVars: バックアップする環境変数名のリスト
//
// Returns:
//   - map[string]string: 環境変数名と値のマップ
func backupEnvVars(envVars []string) map[string]string {
	backup := make(map[string]string)
	for _, key := range envVars {
		backup[key] = os.Getenv(key)
	}
	return backup
}

// restoreEnvVars はバックアップされた環境変数を復元します。
//
// Parameters:
//   - backup: 環境変数名と値のマップ
func restoreEnvVars(backup map[string]string) {
	for key, value := range backup {
		if value != "" {
			Expect(os.Setenv(key, value)).To(Succeed())
		} else {
			Expect(os.Unsetenv(key)).To(Succeed())
		}
	}
}

// clearDBEnvVars はすべてのDB関連環境変数をクリアします。
func clearDBEnvVars() {
	for _, key := range testEnvVars {
		Expect(os.Unsetenv(key)).To(Succeed())
	}
}

// setupTestConfig は一時ディレクトリと設定ファイルを作成します。
//
// Parameters:
//   - configContent: 設定ファイルの内容
//
// Returns:
//   - string: 一時ディレクトリのパス
//   - *os.File: 作成された設定ファイル
func setupTestConfig(configContent string) (string, *os.File) {
	tempDir, err := os.MkdirTemp("", "config_test_*")
	Expect(err).NotTo(HaveOccurred())

	configFile, err := os.Create(tempDir + "/test_config.toml")
	Expect(err).NotTo(HaveOccurred())

	_, err = configFile.WriteString(configContent)
	Expect(err).NotTo(HaveOccurred())
	Expect(configFile.Close()).To(Succeed())

	return tempDir, configFile
}

// cleanupTestConfig は一時ディレクトリを削除します。
//
// Parameters:
//   - tempDir: 削除する一時ディレクトリのパス
func cleanupTestConfig(tempDir string) {
	if tempDir != "" {
		_ = os.RemoveAll(tempDir)
	}
}

// productionTestConfigContent は本番環境に近い設定内容です。
const productionTestConfigContent = `[log]
level = "info"

[mysql]
dbname = "sample_db"
host = "localhost"
port = 3306
user = "root"
pass = "password"
max_idle_conns = 10
max_open_conns = 100
conn_max_lifetime = "30m"
conn_max_idle_time = "5s"
`

var _ = Describe("NewDBConfig関数", func() {
	var originalEnvVars map[string]string
	var tempDir string

	BeforeEach(func() {
		originalEnvVars = backupEnvVars(testEnvVars)
		tempDir, _ = setupTestConfig(productionTestConfigContent)
	})

	AfterEach(func() {
		restoreEnvVars(originalEnvVars)
		cleanupTestConfig(tempDir)
	})

	Context("環境変数が設定されていない場合", func() {
		BeforeEach(func() {
			clearDBEnvVars()
		})

		It("TOMLファイルから設定を読み込む", func() {
			v := config.NewViper(tempDir, "test_config")
			config, err := NewDBConfig(v)
			Expect(err).NotTo(HaveOccurred())
			Expect(config).NotTo(BeNil())
			Expect(config.DBName).To(Equal("sample_db"))
			Expect(config.Host).To(Equal("localhost"))
			Expect(config.Port).To(Equal(3306))
		})
	})

	Context("環境変数が設定されている場合", func() {
		BeforeEach(func() {
			Expect(os.Setenv("DB_DBNAME", "env_test_db")).To(Succeed())
			Expect(os.Setenv("DB_HOST", "envhost")).To(Succeed())
			Expect(os.Setenv("DB_PORT", "3309")).To(Succeed())
			Expect(os.Setenv("DB_USER", "envuser")).To(Succeed())
			Expect(os.Setenv("DB_PASS", "envpass")).To(Succeed())
		})

		It("環境変数がtomlファイルの設定を上書きする", func() {
			v := config.NewViper(tempDir, "test_config")
			config, err := NewDBConfig(v)
			Expect(err).NotTo(HaveOccurred())
			Expect(config).NotTo(BeNil())
			Expect(config.DBName).To(Equal("env_test_db"))
			Expect(config.Host).To(Equal("envhost"))
			Expect(config.Port).To(Equal(3309))
			Expect(config.User).To(Equal("envuser"))
			Expect(config.Pass).To(Equal("envpass"))
		})
	})

	Context("一部の環境変数のみ設定されている場合", func() {
		BeforeEach(func() {
			Expect(os.Unsetenv("DB_DBNAME")).To(Succeed())
			Expect(os.Unsetenv("DB_HOST")).To(Succeed())
			Expect(os.Unsetenv("DB_PORT")).To(Succeed())
			Expect(os.Unsetenv("DB_USER")).To(Succeed())
			Expect(os.Unsetenv("DB_PASS")).To(Succeed())

			Expect(os.Setenv("DB_DBNAME", "custom_db")).To(Succeed())
			Expect(os.Setenv("DB_PORT", "3308")).To(Succeed())
		})

		It("設定された環境変数とTOMLファイルのデフォルト値を組み合わせる", func() {
			v := config.NewViper(tempDir, "test_config")
			config, err := NewDBConfig(v)

			Expect(err).NotTo(HaveOccurred())
			Expect(config.DBName).To(Equal("custom_db"))
			Expect(config.Host).To(Equal("localhost")) // TOMLファイルのデフォルト
			Expect(config.Port).To(Equal(3308))
			Expect(config.User).To(Equal("root"))     // TOMLファイルのデフォルト
			Expect(config.Pass).To(Equal("password")) // TOMLファイルのデフォルト
		})
	})
})

var _ = Describe("NewDatabase関数", Label("DBConnect"), func() {
	var originalEnvVars map[string]string
	var tempDir string

	BeforeEach(func() {
		allEnvVars := append([]string{"LOG_LEVEL"}, testEnvVars...)
		originalEnvVars = backupEnvVars(allEnvVars)
		tempDir, _ = setupTestConfig(productionTestConfigContent)
	})

	AfterEach(func() {
		restoreEnvVars(originalEnvVars)
		cleanupTestConfig(tempDir)
	})

	DescribeTable("NewDatabaseの動作確認",
		func(setupConfig func() *DBConfig, expectError bool, expectedErrorSubstring string) {
			config := setupConfig()
			db, err := NewDatabase(config)
			if expectError {
				Expect(err).To(HaveOccurred())
				if expectedErrorSubstring != "" {
					Expect(err.Error()).To(ContainSubstring(expectedErrorSubstring))
				}
				Expect(db).To(BeNil())
			} else {
				Expect(err).NotTo(HaveOccurred())
				Expect(db).NotTo(BeNil())
				// 接続をクリーンアップ
				if db != nil {
					Expect(db.Close()).To(Succeed())
				}
			}
		},
		Entry("正常系: 有効な設定で接続", func() *DBConfig {
			return &DBConfig{
				DBName:          "sample_db",
				Host:            "localhost",
				Port:            3306,
				User:            "root",
				Pass:            "password",
				MaxIdleConns:    10,
				MaxOpenConns:    100,
				ConnMaxLifetime: 30 * time.Minute,
				ConnMaxIdleTime: 5 * time.Second,
			}
		}, false, ""),
		Entry("異常系: 不正なホスト名", func() *DBConfig {
			return &DBConfig{
				DBName:          "sample_db",
				Host:            "invalid_host",
				Port:            3306,
				User:            "root",
				Pass:            "password",
				MaxIdleConns:    10,
				MaxOpenConns:    100,
				ConnMaxLifetime: 30 * time.Minute,
				ConnMaxIdleTime: 5 * time.Second,
			}
		}, true, ""),
		Entry("異常系: 不正なポート番号", func() *DBConfig {
			return &DBConfig{
				DBName:          "sample_db",
				Host:            "localhost",
				Port:            99999,
				User:            "root",
				Pass:            "password",
				MaxIdleConns:    10,
				MaxOpenConns:    100,
				ConnMaxLifetime: 30 * time.Minute,
				ConnMaxIdleTime: 5 * time.Second,
			}
		}, true, ""),
	)

	Context("環境変数から設定を取得して接続する場合", func() {
		BeforeEach(func() {
			Expect(os.Setenv("DB_DBNAME", "sample_db")).To(Succeed())
			Expect(os.Setenv("DB_HOST", "localhost")).To(Succeed())
			Expect(os.Setenv("DB_PORT", "3306")).To(Succeed())
			Expect(os.Setenv("DB_USER", "root")).To(Succeed())
			Expect(os.Setenv("DB_PASS", "password")).To(Succeed())
		})

		It("環境変数から取得した設定で接続できる", func() {
			v := config.NewViper(tempDir, "test_config")
			cfg, err := NewDBConfig(v)
			Expect(err).NotTo(HaveOccurred())

			db, err := NewDatabase(cfg)
			Expect(err).NotTo(HaveOccurred())
			Expect(db).NotTo(BeNil())

			// 接続をクリーンアップ
			Expect(db.Close()).To(Succeed())
		})
	})

	Context("TOMLファイルから設定を取得して接続する場合", func() {
		BeforeEach(func() {
			clearDBEnvVars()
		})

		It("TOMLファイルから取得した設定で接続できる", func() {
			v := config.NewViper(tempDir, "test_config")
			cfg, err := NewDBConfig(v)
			Expect(err).NotTo(HaveOccurred())

			db, err := NewDatabase(cfg)
			Expect(err).NotTo(HaveOccurred())
			Expect(db).NotTo(BeNil())

			// 接続をクリーンアップ
			Expect(db.Close()).To(Succeed())
		})
	})
})
