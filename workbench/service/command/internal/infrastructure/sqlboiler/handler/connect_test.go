//go:build integration || !ci

package handler

import (
	"os"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestConn(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "infrastructure/sqlboiler/handler packageのテスト")
}

var _ = Describe("setupViper関数", Label("setupViper"), func() {
	var originalEnvVars map[string]string
	var tempDir string
	var configFile *os.File

	BeforeEach(func() {
		// 環境変数のバックアップ
		originalEnvVars = map[string]string{
			"DB_DBNAME":             os.Getenv("DB_DBNAME"),
			"DB_HOST":               os.Getenv("DB_HOST"),
			"DB_PORT":               os.Getenv("DB_PORT"),
			"DB_USER":               os.Getenv("DB_USER"),
			"DB_PASS":               os.Getenv("DB_PASS"),
			"DB_MAX_IDLE_CONNS":     os.Getenv("DB_MAX_IDLE_CONNS"),
			"DB_MAX_OPEN_CONNS":     os.Getenv("DB_MAX_OPEN_CONNS"),
			"DB_CONN_MAX_LIFETIME":  os.Getenv("DB_CONN_MAX_LIFETIME"),
			"DB_CONN_MAX_IDLE_TIME": os.Getenv("DB_CONN_MAX_IDLE_TIME"),
		}

		// 一時ディレクトリを作成
		var err error
		tempDir, err = os.MkdirTemp("", "config_test_*")
		Expect(err).NotTo(HaveOccurred())

		// 一時設定ファイルを作成
		configFile, err = os.Create(tempDir + "/test_config.toml")
		Expect(err).NotTo(HaveOccurred())

		// テスト用の設定を書き込む
		_, err = configFile.WriteString(`[mysql]
dbname = "test_sample_db"
host = "test_localhost"
port = 3306
user = "test_root"
pass = "test_password"
max_idle_conns = 10
max_open_conns = 100
conn_max_lifetime = "30m"
conn_max_idle_time = "500s"
`)
		Expect(err).NotTo(HaveOccurred())
		Expect(configFile.Close()).To(Succeed())
	})

	AfterEach(func() {
		// 環境変数を元に戻す
		for key, value := range originalEnvVars {
			if value != "" {
				Expect(os.Setenv(key, value)).To(Succeed())
			} else {
				Expect(os.Unsetenv(key)).To(Succeed())
			}
		}

		// 一時ディレクトリを削除
		if tempDir != "" {
			os.RemoveAll(tempDir)
		}
	})

	Context("環境変数が設定されていない場合", func() {
		BeforeEach(func() {
			// すべての環境変数をクリア
			Expect(os.Unsetenv("DB_DBNAME")).To(Succeed())
			Expect(os.Unsetenv("DB_HOST")).To(Succeed())
			Expect(os.Unsetenv("DB_PORT")).To(Succeed())
			Expect(os.Unsetenv("DB_USER")).To(Succeed())
			Expect(os.Unsetenv("DB_PASS")).To(Succeed())
			Expect(os.Unsetenv("DB_MAX_IDLE_CONNS")).To(Succeed())
			Expect(os.Unsetenv("DB_MAX_OPEN_CONNS")).To(Succeed())
			Expect(os.Unsetenv("DB_CONN_MAX_LIFETIME")).To(Succeed())
			Expect(os.Unsetenv("DB_CONN_MAX_IDLE_TIME")).To(Succeed())
		})

		It("tomlファイルから設定を読み込む", func() {
			v := setupViper(tempDir, "test_config")

			Expect(v.GetString("mysql.dbname")).To(Equal("test_sample_db"))
			Expect(v.GetString("mysql.host")).To(Equal("test_localhost"))
			Expect(v.GetInt("mysql.port")).To(Equal(3306))
			Expect(v.GetString("mysql.user")).To(Equal("test_root"))
			Expect(v.GetString("mysql.pass")).To(Equal("test_password"))
			Expect(v.GetInt("mysql.max_idle_conns")).To(Equal(10))
			Expect(v.GetInt("mysql.max_open_conns")).To(Equal(100))
			Expect(v.GetDuration("mysql.conn_max_lifetime")).To(Equal(30 * time.Minute))
			Expect(v.GetDuration("mysql.conn_max_idle_time")).To(Equal(500 * time.Second))
		})
	})

	Context("環境変数が設定されている場合", func() {
		BeforeEach(func() {
			Expect(os.Setenv("DB_DBNAME", "env_test_db")).To(Succeed())
			Expect(os.Setenv("DB_HOST", "envhost")).To(Succeed())
			Expect(os.Setenv("DB_PORT", "3307")).To(Succeed())
			Expect(os.Setenv("DB_USER", "envuser")).To(Succeed())
			Expect(os.Setenv("DB_PASS", "envpass")).To(Succeed())
		})

		It("環境変数がtomlファイルの設定を上書きする", func() {
			v := setupViper(tempDir, "test_config")

			Expect(v.GetString("mysql.dbname")).To(Equal("env_test_db"))
			Expect(v.GetString("mysql.host")).To(Equal("envhost"))
			Expect(v.GetInt("mysql.port")).To(Equal(3307))
			Expect(v.GetString("mysql.user")).To(Equal("envuser"))
			Expect(v.GetString("mysql.pass")).To(Equal("envpass"))
		})
	})
})

var _ = Describe("NewDBConfig関数", func() {
	var originalEnvVars map[string]string

	BeforeEach(func() {
		// 環境変数のバックアップ
		originalEnvVars = map[string]string{
			"DB_DBNAME":             os.Getenv("DB_DBNAME"),
			"DB_HOST":               os.Getenv("DB_HOST"),
			"DB_PORT":               os.Getenv("DB_PORT"),
			"DB_USER":               os.Getenv("DB_USER"),
			"DB_PASS":               os.Getenv("DB_PASS"),
			"DB_MAX_IDLE_CONNS":     os.Getenv("DB_MAX_IDLE_CONNS"),
			"DB_MAX_OPEN_CONNS":     os.Getenv("DB_MAX_OPEN_CONNS"),
			"DB_CONN_MAX_LIFETIME":  os.Getenv("DB_CONN_MAX_LIFETIME"),
			"DB_CONN_MAX_IDLE_TIME": os.Getenv("DB_CONN_MAX_IDLE_TIME"),
		}
	})

	AfterEach(func() {
		// 環境変数を元に戻す
		for key, value := range originalEnvVars {
			if value != "" {
				Expect(os.Setenv(key, value)).To(Succeed())
			} else {
				Expect(os.Unsetenv(key)).To(Succeed())
			}
		}
	})

	Context("環境変数が設定されていない場合", func() {
		BeforeEach(func() {
			// すべての環境変数をクリア
			Expect(os.Unsetenv("DB_DBNAME")).To(Succeed())
			Expect(os.Unsetenv("DB_HOST")).To(Succeed())
			Expect(os.Unsetenv("DB_PORT")).To(Succeed())
			Expect(os.Unsetenv("DB_USER")).To(Succeed())
			Expect(os.Unsetenv("DB_PASS")).To(Succeed())
			Expect(os.Unsetenv("DB_MAX_IDLE_CONNS")).To(Succeed())
			Expect(os.Unsetenv("DB_MAX_OPEN_CONNS")).To(Succeed())
			Expect(os.Unsetenv("DB_CONN_MAX_LIFETIME")).To(Succeed())
			Expect(os.Unsetenv("DB_CONN_MAX_IDLE_TIME")).To(Succeed())
		})

		It("TOMLファイルから設定を読み込む", func() {
			config, err := NewDBConfig()
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
			config, err := NewDBConfig()
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
			config, err := NewDBConfig()

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

	BeforeEach(func() {
		// 環境変数のバックアップ
		originalEnvVars = map[string]string{
			"LOG_LEVEL": os.Getenv("LOG_LEVEL"),
		}
	})

	AfterEach(func() {
		// 環境変数を元に戻す
		for key, value := range originalEnvVars {
			if value != "" {
				Expect(os.Setenv(key, value)).To(Succeed())
			} else {
				Expect(os.Unsetenv(key)).To(Succeed())
			}
		}
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
			config, err := NewDBConfig()
			Expect(err).NotTo(HaveOccurred())

			db, err := NewDatabase(config)
			Expect(err).NotTo(HaveOccurred())
			Expect(db).NotTo(BeNil())

			// 接続をクリーンアップ
			Expect(db.Close()).To(Succeed())
		})
	})

	Context("TOMLファイルから設定を取得して接続する場合", func() {
		BeforeEach(func() {
			// すべての環境変数をクリアしてTOMLファイルのみに依存する
			Expect(os.Unsetenv("DB_DBNAME")).To(Succeed())
			Expect(os.Unsetenv("DB_HOST")).To(Succeed())
			Expect(os.Unsetenv("DB_PORT")).To(Succeed())
			Expect(os.Unsetenv("DB_USER")).To(Succeed())
			Expect(os.Unsetenv("DB_PASS")).To(Succeed())
		})

		It("TOMLファイルから取得した設定で接続できる", func() {
			config, err := NewDBConfig()
			Expect(err).NotTo(HaveOccurred())

			db, err := NewDatabase(config)
			Expect(err).NotTo(HaveOccurred())
			Expect(db).NotTo(BeNil())

			// 接続をクリーンアップ
			Expect(db.Close()).To(Succeed())
		})
	})

	Context("TOMLファイルから設定を取得して接続する場合", func() {
		BeforeEach(func() {
			// すべての環境変数をクリアしてTOMLファイルのみに依存する
			Expect(os.Unsetenv("DB_DBNAME")).To(Succeed())
			Expect(os.Unsetenv("DB_HOST")).To(Succeed())
			Expect(os.Unsetenv("DB_PORT")).To(Succeed())
			Expect(os.Unsetenv("DB_USER")).To(Succeed())
			Expect(os.Unsetenv("DB_PASS")).To(Succeed())
		})

		It("TOMLファイルから取得した設定で接続できる", func() {
			config, err := NewDBConfig()
			Expect(err).NotTo(HaveOccurred())

			db, err := NewDatabase(config)
			Expect(err).NotTo(HaveOccurred())
			Expect(db).NotTo(BeNil())

			// 接続をクリーンアップ
			Expect(db.Close()).To(Succeed())
		})
	})
})
