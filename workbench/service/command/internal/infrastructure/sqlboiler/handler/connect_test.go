//go:build integration || !ci

package handler

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestConn(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "infrastructure/sqlboiler/handler packageのテスト")
}

var _ = Describe("loadConfigFromEnv関数", func() {
	var originalEnvVars map[string]string

	BeforeEach(func() {
		// 環境変数のバックアップ
		originalEnvVars = map[string]string{
			"DB_NAME":               os.Getenv("DB_NAME"),
			"DB_HOST":               os.Getenv("DB_HOST"),
			"DB_PORT":               os.Getenv("DB_PORT"),
			"DB_USER":               os.Getenv("DB_USER"),
			"DB_PASSWORD":           os.Getenv("DB_PASSWORD"),
			"DB_MAX_IDLE_CONNS":     os.Getenv("DB_MAX_IDLE_CONNS"),
			"DB_MAX_OPEN_CONNS":     os.Getenv("DB_MAX_OPEN_CONNS"),
			"DB_CONN_MAX_LIFETIME":  os.Getenv("DB_CONN_MAX_LIFETIME"),
			"DB_CONN_MAX_IDLE_TIME": os.Getenv("DB_CONN_MAX_IDLE_TIME"),
			"DATABASE_TOML_PATH":    os.Getenv("DATABASE_TOML_PATH"),
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

	Context("環境変数が設定されている場合", func() {
		BeforeEach(func() {
			Expect(os.Setenv("DB_NAME", "test_db")).To(Succeed())
			Expect(os.Setenv("DB_HOST", "testhost")).To(Succeed())
			Expect(os.Setenv("DB_PORT", "3307")).To(Succeed())
			Expect(os.Setenv("DB_USER", "testuser")).To(Succeed())
			Expect(os.Setenv("DB_PASSWORD", "testpass")).To(Succeed())
			Expect(os.Setenv("DB_MAX_IDLE_CONNS", "20")).To(Succeed())
			Expect(os.Setenv("DB_MAX_OPEN_CONNS", "200")).To(Succeed())
			Expect(os.Setenv("DB_CONN_MAX_LIFETIME", "1h")).To(Succeed())
			Expect(os.Setenv("DB_CONN_MAX_IDLE_TIME", "10s")).To(Succeed())
		})

		It("環境変数から正しく設定を読み込む", func() {
			config, err := loadConfigFromEnv()

			Expect(err).NotTo(HaveOccurred())
			Expect(config.DBName).To(Equal("test_db"))
			Expect(config.Host).To(Equal("testhost"))
			Expect(config.Port).To(Equal(3307))
			Expect(config.User).To(Equal("testuser"))
			Expect(config.Pass).To(Equal("testpass"))
			Expect(config.MaxIdleConns).To(Equal(20))
			Expect(config.MaxOpenConns).To(Equal(200))
			Expect(config.ConnMaxLifetime).To(Equal(time.Hour))
			Expect(config.ConnMaxIdleTime).To(Equal(10 * time.Second))
		})
	})

	Context("環境変数が設定されていない場合", func() {
		BeforeEach(func() {
			// すべての環境変数をクリア
			Expect(os.Unsetenv("DB_NAME")).To(Succeed())
			Expect(os.Unsetenv("DB_HOST")).To(Succeed())
			Expect(os.Unsetenv("DB_PORT")).To(Succeed())
			Expect(os.Unsetenv("DB_USER")).To(Succeed())
			Expect(os.Unsetenv("DB_PASSWORD")).To(Succeed())
			Expect(os.Unsetenv("DB_MAX_IDLE_CONNS")).To(Succeed())
			Expect(os.Unsetenv("DB_MAX_OPEN_CONNS")).To(Succeed())
			Expect(os.Unsetenv("DB_CONN_MAX_LIFETIME")).To(Succeed())
			Expect(os.Unsetenv("DB_CONN_MAX_IDLE_TIME")).To(Succeed())
		})

		It("デフォルト値を使用する", func() {
			config, err := loadConfigFromEnv()

			Expect(err).NotTo(HaveOccurred())
			Expect(config.DBName).To(Equal("sample_db"))
			Expect(config.Host).To(Equal("localhost"))
			Expect(config.Port).To(Equal(3306))
			Expect(config.User).To(Equal("root"))
			Expect(config.Pass).To(Equal("password"))
			Expect(config.MaxIdleConns).To(Equal(10))
			Expect(config.MaxOpenConns).To(Equal(100))
			Expect(config.ConnMaxLifetime).To(Equal(30 * time.Minute))
			Expect(config.ConnMaxIdleTime).To(Equal(5 * time.Second))
		})
	})

	Context("一部の環境変数のみ設定されている場合", func() {
		BeforeEach(func() {
			Expect(os.Unsetenv("DB_NAME")).To(Succeed())
			Expect(os.Unsetenv("DB_HOST")).To(Succeed())
			Expect(os.Unsetenv("DB_PORT")).To(Succeed())
			Expect(os.Unsetenv("DB_USER")).To(Succeed())
			Expect(os.Unsetenv("DB_PASSWORD")).To(Succeed())
			Expect(os.Unsetenv("DB_MAX_IDLE_CONNS")).To(Succeed())
			Expect(os.Unsetenv("DB_MAX_OPEN_CONNS")).To(Succeed())
			Expect(os.Unsetenv("DB_CONN_MAX_LIFETIME")).To(Succeed())
			Expect(os.Unsetenv("DB_CONN_MAX_IDLE_TIME")).To(Succeed())

			Expect(os.Setenv("DB_NAME", "custom_db")).To(Succeed())
			Expect(os.Setenv("DB_PORT", "3308")).To(Succeed())
		})

		It("設定された環境変数と未設定のデフォルト値を組み合わせる", func() {
			config, err := loadConfigFromEnv()

			Expect(err).NotTo(HaveOccurred())
			Expect(config.DBName).To(Equal("custom_db"))
			Expect(config.Host).To(Equal("localhost")) // デフォルト
			Expect(config.Port).To(Equal(3308))
			Expect(config.User).To(Equal("root"))     // デフォルト
			Expect(config.Pass).To(Equal("password")) // デフォルト
		})
	})

	Context("不正な環境変数値が設定されている場合", func() {
		BeforeEach(func() {
			// すべての環境変数をクリア
			Expect(os.Unsetenv("DB_NAME")).To(Succeed())
			Expect(os.Unsetenv("DB_HOST")).To(Succeed())
			Expect(os.Unsetenv("DB_PORT")).To(Succeed())
			Expect(os.Unsetenv("DB_USER")).To(Succeed())
			Expect(os.Unsetenv("DB_PASSWORD")).To(Succeed())
			Expect(os.Unsetenv("DB_MAX_IDLE_CONNS")).To(Succeed())
			Expect(os.Unsetenv("DB_MAX_OPEN_CONNS")).To(Succeed())
			Expect(os.Unsetenv("DB_CONN_MAX_LIFETIME")).To(Succeed())
			Expect(os.Unsetenv("DB_CONN_MAX_IDLE_TIME")).To(Succeed())
		})

		It("DB_PORTが不正な値の場合、エラーを返す", func() {
			Expect(os.Setenv("DB_PORT", "invalid_port")).To(Succeed())

			config, err := loadConfigFromEnv()

			Expect(err).To(HaveOccurred())
			// strconv.Atoiのパースエラーが含まれることを確認
			Expect(err.Error()).To(Or(
				ContainSubstring("invalid_port"),
				ContainSubstring("invalid syntax"),
			))
			// エラーが発生してもデフォルト値が設定される
			Expect(config.Port).To(Equal(3306))
		})

		It("DB_MAX_IDLE_CONNSが不正な値の場合、エラーを返す", func() {
			Expect(os.Setenv("DB_MAX_IDLE_CONNS", "not_a_number")).To(Succeed())

			config, err := loadConfigFromEnv()

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Or(
				ContainSubstring("not_a_number"),
				ContainSubstring("invalid syntax"),
			))
			Expect(config.MaxIdleConns).To(Equal(10))
		})

		It("DB_CONN_MAX_LIFETIMEが不正な値の場合、エラーを返す", func() {
			Expect(os.Setenv("DB_CONN_MAX_LIFETIME", "invalid_duration")).To(Succeed())

			config, err := loadConfigFromEnv()

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Or(
				ContainSubstring("invalid_duration"),
				ContainSubstring("time: invalid duration"),
			))
			Expect(config.ConnMaxLifetime).To(Equal(30 * time.Minute))
		})

		It("複数の環境変数が不正な値の場合、すべてのエラーを返す", func() {
			Expect(os.Setenv("DB_PORT", "invalid")).To(Succeed())
			Expect(os.Setenv("DB_MAX_IDLE_CONNS", "invalid")).To(Succeed())
			Expect(os.Setenv("DB_CONN_MAX_LIFETIME", "invalid")).To(Succeed())

			config, err := loadConfigFromEnv()

			Expect(err).To(HaveOccurred())
			// errors.Joinで結合されたエラーを確認
			// 各エラーメッセージには不正な値またはエラータイプが含まれる
			errMsg := err.Error()
			Expect(errMsg).To(ContainSubstring("invalid"))
			// すべてデフォルト値が設定される
			Expect(config.Port).To(Equal(3306))
			Expect(config.MaxIdleConns).To(Equal(10))
			Expect(config.ConnMaxLifetime).To(Equal(30 * time.Minute))
		})
	})
})

var _ = Describe("loadConfig関数", func() {
	var originalEnvVars map[string]string

	BeforeEach(func() {
		// 環境変数のバックアップ
		originalEnvVars = map[string]string{
			"DB_NAME":               os.Getenv("DB_NAME"),
			"DB_HOST":               os.Getenv("DB_HOST"),
			"DB_PORT":               os.Getenv("DB_PORT"),
			"DB_USER":               os.Getenv("DB_USER"),
			"DB_PASSWORD":           os.Getenv("DB_PASSWORD"),
			"DB_MAX_IDLE_CONNS":     os.Getenv("DB_MAX_IDLE_CONNS"),
			"DB_MAX_OPEN_CONNS":     os.Getenv("DB_MAX_OPEN_CONNS"),
			"DB_CONN_MAX_LIFETIME":  os.Getenv("DB_CONN_MAX_LIFETIME"),
			"DB_CONN_MAX_IDLE_TIME": os.Getenv("DB_CONN_MAX_IDLE_TIME"),
			"DATABASE_TOML_PATH":    os.Getenv("DATABASE_TOML_PATH"),
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

	Context("DATABASE_TOML_PATHが設定されている場合", func() {
		BeforeEach(func() {
			absPath, err := filepath.Abs("../config/sqlboiler.toml")
			if err != nil {
				Fail("Failed to get absolute path")
			}
			Expect(os.Setenv("DATABASE_TOML_PATH", absPath)).To(Succeed())
		})

		It("TOMLファイルから設定を読み込む", func() {
			config, err := loadConfig()
			Expect(err).NotTo(HaveOccurred())
			Expect(config).NotTo(BeNil())
			Expect(config.DBName).To(Equal("sample_db"))
			Expect(config.Host).To(Equal("localhost"))
			Expect(config.Port).To(Equal(3306))
		})
	})

	Context("DATABASE_TOML_PATHが設定されていない場合", func() {
		BeforeEach(func() {
			Expect(os.Unsetenv("DATABASE_TOML_PATH")).To(Succeed())
			Expect(os.Setenv("DB_NAME", "env_test_db")).To(Succeed())
			Expect(os.Setenv("DB_HOST", "envhost")).To(Succeed())
			Expect(os.Setenv("DB_PORT", "3309")).To(Succeed())
		})

		It("環境変数から設定を読み込む", func() {
			config, err := loadConfig()
			Expect(err).NotTo(HaveOccurred())
			Expect(config).NotTo(BeNil())
			Expect(config.DBName).To(Equal("env_test_db"))
			Expect(config.Host).To(Equal("envhost"))
			Expect(config.Port).To(Equal(3309))
		})
	})

	Context("存在しないTOMLファイルパスが設定されている場合", func() {
		BeforeEach(func() {
			Expect(os.Setenv("DATABASE_TOML_PATH", "/nonexistent/path/config.toml")).To(Succeed())
		})

		It("エラーを返す", func() {
			config, err := loadConfig()
			Expect(err).To(HaveOccurred())
			Expect(config).To(BeNil())
		})
	})

	Context("TOMLファイルにmysqlキーが存在しない場合", func() {
		var tempFile *os.File

		BeforeEach(func() {
			var err error
			tempFile, err = os.CreateTemp("", "invalid_config_*.toml")
			Expect(err).NotTo(HaveOccurred())

			// mysqlキーが存在しない設定ファイルを作成
			_, err = tempFile.WriteString(`
[postgres]
dbname = "test_db"
host = "localhost"
port = 5432
`)
			Expect(err).NotTo(HaveOccurred())
			Expect(tempFile.Close()).To(Succeed())

			Expect(os.Setenv("DATABASE_TOML_PATH", tempFile.Name())).To(Succeed())
		})

		AfterEach(func() {
			if tempFile != nil {
				os.Remove(tempFile.Name())
			}
		})

		It("エラーを返す", func() {
			config, err := loadConfig()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("key 'mysql' not found"))
			Expect(config).To(BeNil())
		})
	})
})

var _ = Describe("DBConnect関数", Label("DBConnect"), func() {
	var originalEnvVars map[string]string

	BeforeEach(func() {
		// 環境変数のバックアップ
		originalEnvVars = map[string]string{
			"DATABASE_TOML_PATH": os.Getenv("DATABASE_TOML_PATH"),
			"LOG_LEVEL":          os.Getenv("LOG_LEVEL"),
		}

		// テスト前のセットアップ処理
		absPath, err := filepath.Abs("../config/sqlboiler.toml")
		if err != nil {
			Fail("Failed to get absolute path")
		}
		Expect(os.Setenv("DATABASE_TOML_PATH", absPath)).To(Succeed())
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

	DescribeTable("DBConnectの動作確認",
		func(expectError bool, expectedErrorCode string) {
			err := DBConnect()
			if expectError {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(expectedErrorCode))
			} else {
				Expect(err).NotTo(HaveOccurred())
			}
		},
		Entry("正常系: TOMLファイルから接続", false, ""),
	)

	Context("環境変数から接続する場合", func() {
		BeforeEach(func() {
			Expect(os.Unsetenv("DATABASE_TOML_PATH")).To(Succeed())
			Expect(os.Setenv("DB_NAME", "sample_db")).To(Succeed())
			Expect(os.Setenv("DB_HOST", "localhost")).To(Succeed())
			Expect(os.Setenv("DB_PORT", "3306")).To(Succeed())
			Expect(os.Setenv("DB_USER", "root")).To(Succeed())
			Expect(os.Setenv("DB_PASSWORD", "password")).To(Succeed())
		})

		It("環境変数から接続できる", func() {
			err := DBConnect()
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
