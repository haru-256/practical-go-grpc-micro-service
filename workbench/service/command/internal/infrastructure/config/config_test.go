package config

import (
	"os"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "config packageのテスト")
}

// testEnvVars はテストで使用される環境変数のリストです。
var testEnvVars = []string{
	"LOG_LEVEL",
	"LOG_FORMAT",
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
func backupEnvVars(envVars []string) map[string]string {
	backup := make(map[string]string)
	for _, key := range envVars {
		backup[key] = os.Getenv(key)
	}
	return backup
}

// restoreEnvVars はバックアップされた環境変数を復元します。
func restoreEnvVars(backup map[string]string) {
	for key, value := range backup {
		if value != "" {
			Expect(os.Setenv(key, value)).To(Succeed())
		} else {
			Expect(os.Unsetenv(key)).To(Succeed())
		}
	}
}

// clearTestEnvVars はすべてのテスト関連環境変数をクリアします。
func clearTestEnvVars() {
	for _, key := range testEnvVars {
		Expect(os.Unsetenv(key)).To(Succeed())
	}
}

// setupTestConfig は一時ディレクトリと設定ファイルを作成します。
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
func cleanupTestConfig(tempDir string) {
	if tempDir != "" {
		_ = os.RemoveAll(tempDir)
	}
}

// defaultTestConfigContent はテストで使用するデフォルトの設定内容です。
const defaultTestConfigContent = `[log]
level = "info"
format = "text"

[mysql]
dbname = "test_db"
host = "test_host"
port = 3307
user = "test_user"
pass = "test_pass"
max_idle_conns = 5
max_open_conns = 50
conn_max_lifetime = "15m"
conn_max_idle_time = "300s"
`

var _ = Describe("NewViper関数", func() {
	var originalEnvVars map[string]string
	var tempDir string

	BeforeEach(func() {
		originalEnvVars = backupEnvVars(testEnvVars)
		tempDir, _ = setupTestConfig(defaultTestConfigContent)
	})

	AfterEach(func() {
		restoreEnvVars(originalEnvVars)
		cleanupTestConfig(tempDir)
	})

	Context("環境変数が設定されていない場合", func() {
		BeforeEach(func() {
			clearTestEnvVars()
		})

		It("TOMLファイルから設定を読み込む", func() {
			v := NewViper(tempDir, "test_config")

			// log設定
			Expect(v.GetString("log.level")).To(Equal("info"))
			Expect(v.GetString("log.format")).To(Equal("text"))

			// mysql設定
			Expect(v.GetString("mysql.dbname")).To(Equal("test_db"))
			Expect(v.GetString("mysql.host")).To(Equal("test_host"))
			Expect(v.GetInt("mysql.port")).To(Equal(3307))
			Expect(v.GetString("mysql.user")).To(Equal("test_user"))
			Expect(v.GetString("mysql.pass")).To(Equal("test_pass"))
			Expect(v.GetInt("mysql.max_idle_conns")).To(Equal(5))
			Expect(v.GetInt("mysql.max_open_conns")).To(Equal(50))
			Expect(v.GetDuration("mysql.conn_max_lifetime")).To(Equal(15 * time.Minute))
			Expect(v.GetDuration("mysql.conn_max_idle_time")).To(Equal(300 * time.Second))
		})
	})

	Context("環境変数が設定されている場合", func() {
		BeforeEach(func() {
			Expect(os.Setenv("DB_DBNAME", "env_db")).To(Succeed())
			Expect(os.Setenv("DB_HOST", "env_host")).To(Succeed())
			Expect(os.Setenv("DB_PORT", "3308")).To(Succeed())
			Expect(os.Setenv("DB_USER", "env_user")).To(Succeed())
			Expect(os.Setenv("DB_PASS", "env_pass")).To(Succeed())
		})

		It("環境変数がTOMLファイルの設定を上書きする", func() {
			v := NewViper(tempDir, "test_config")

			Expect(v.GetString("mysql.dbname")).To(Equal("env_db"))
			Expect(v.GetString("mysql.host")).To(Equal("env_host"))
			Expect(v.GetInt("mysql.port")).To(Equal(3308))
			Expect(v.GetString("mysql.user")).To(Equal("env_user"))
			Expect(v.GetString("mysql.pass")).To(Equal("env_pass"))
		})
	})

	Context("一部の環境変数のみ設定されている場合", func() {
		BeforeEach(func() {
			clearTestEnvVars()
			Expect(os.Setenv("DB_DBNAME", "partial_db")).To(Succeed())
			Expect(os.Setenv("DB_PORT", "3309")).To(Succeed())
		})

		It("環境変数とTOMLファイルの値を組み合わせる", func() {
			v := NewViper(tempDir, "test_config")

			Expect(v.GetString("mysql.dbname")).To(Equal("partial_db"))
			Expect(v.GetString("mysql.host")).To(Equal("test_host")) // TOML
			Expect(v.GetInt("mysql.port")).To(Equal(3309))
			Expect(v.GetString("mysql.user")).To(Equal("test_user")) // TOML
			Expect(v.GetString("mysql.pass")).To(Equal("test_pass")) // TOML
		})
	})

	Context("設定ファイルが存在しない場合", func() {
		It("panicを発生させる", func() {
			Expect(func() {
				NewViper("/nonexistent/path", "missing_config")
			}).To(Panic())
		})
	})

	Context("環境変数のバインディングが正しく動作する場合", func() {
		BeforeEach(func() {
			clearTestEnvVars()
			Expect(os.Setenv("DB_MAX_IDLE_CONNS", "20")).To(Succeed())
			Expect(os.Setenv("DB_MAX_OPEN_CONNS", "200")).To(Succeed())
			Expect(os.Setenv("DB_CONN_MAX_LIFETIME", "60m")).To(Succeed())
			Expect(os.Setenv("DB_CONN_MAX_IDLE_TIME", "10m")).To(Succeed())
		})

		It("DB_プレフィックスの環境変数が正しくバインドされる", func() {
			v := NewViper(tempDir, "test_config")

			Expect(v.GetInt("mysql.max_idle_conns")).To(Equal(20))
			Expect(v.GetInt("mysql.max_open_conns")).To(Equal(200))
			Expect(v.GetDuration("mysql.conn_max_lifetime")).To(Equal(60 * time.Minute))
			Expect(v.GetDuration("mysql.conn_max_idle_time")).To(Equal(10 * time.Minute))
		})
	})

	Context("LOG_LEVEL/LOG_FORMAT環境変数が設定されている場合", func() {
		BeforeEach(func() {
			clearTestEnvVars()
			Expect(os.Setenv("LOG_LEVEL", "debug")).To(Succeed())
			Expect(os.Setenv("LOG_FORMAT", "json")).To(Succeed())
		})

		It("LOG_プレフィックスの環境変数がTOMLファイルの設定を上書きする", func() {
			v := NewViper(tempDir, "test_config")

			Expect(v.GetString("log.level")).To(Equal("debug"))
			Expect(v.GetString("log.format")).To(Equal("json"))
		})
	})

	Context("LOG_LEVEL環境変数のみ設定されている場合", func() {
		BeforeEach(func() {
			clearTestEnvVars()
			Expect(os.Setenv("LOG_LEVEL", "warn")).To(Succeed())
		})

		It("LOG_LEVEL環境変数が上書きされ、LOG_FORMATはTOMLファイルの値を使用する", func() {
			v := NewViper(tempDir, "test_config")

			Expect(v.GetString("log.level")).To(Equal("warn"))
			Expect(v.GetString("log.format")).To(Equal("text")) // TOML
		})
	})
})
