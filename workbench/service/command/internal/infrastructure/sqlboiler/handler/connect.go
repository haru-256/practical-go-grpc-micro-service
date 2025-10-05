package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/haru-256/practical-go-grpc-micro-service/pkg/utils"
)

type DBConfig struct {
	DBName          string        `toml:"dbname"`            //	データベース名
	Host            string        `toml:"host"`              //	ホスト名
	Port            int           `toml:"port"`              //	ポート番号
	User            string        `toml:"user"`              //	ユーザー名
	Pass            string        `toml:"pass"`              //	パスワード
	MaxIdleConns    int           `toml:"max_idle_conns"`    //	最大アイドル接続数
	MaxOpenConns    int           `toml:"max_open_conns"`    //	最大接続数
	ConnMaxLifetime time.Duration `toml:"conn_max_lifetime"` //	接続の最大生存時間(分)
	ConnMaxIdleTime time.Duration `toml:"idle_timeout"`      //	接続の最大アイドル時間(分)
}

// getEnvWithError は GetEnv を呼び出し、エラーが発生した場合はエラーリストに追加してデフォルト値を返します。
func getEnvWithError[T utils.EnvType](key string, defaultValue T, errs *[]error) T {
	val, err := utils.GetEnv(key, defaultValue)
	if err != nil {
		*errs = append(*errs, err)
		return defaultValue // エラー発生時はデフォルト値を返す
	}
	return val
}

func loadConfigFromEnv() (DBConfig, error) {
	var configErrors []error

	cfg := DBConfig{
		DBName:          getEnvWithError("DB_NAME", "sample_db", &configErrors),
		Host:            getEnvWithError("DB_HOST", "localhost", &configErrors),
		Port:            getEnvWithError("DB_PORT", 3306, &configErrors),
		User:            getEnvWithError("DB_USER", "root", &configErrors),
		Pass:            getEnvWithError("DB_PASSWORD", "password", &configErrors),
		MaxIdleConns:    getEnvWithError("DB_MAX_IDLE_CONNS", 10, &configErrors),
		MaxOpenConns:    getEnvWithError("DB_MAX_OPEN_CONNS", 100, &configErrors),
		ConnMaxLifetime: getEnvWithError("DB_CONN_MAX_LIFETIME", time.Duration(30)*time.Minute, &configErrors),
		ConnMaxIdleTime: getEnvWithError("DB_CONN_MAX_IDLE_TIME", time.Duration(5)*time.Second, &configErrors),
	}

	// すべての環境変数を読み込んだ後、エラーがあればまとめて返す
	if len(configErrors) > 0 {
		return cfg, errors.Join(configErrors...)
	}
	return cfg, nil
}

// database.tomlから接続情報を取得してDbConfig型で返す
func loadConfig() (*DBConfig, error) {
	// 環境変数からファイルパスを取得する
	path, ok := os.LookupEnv("DATABASE_TOML_PATH")
	// 設定されている場合はそのパスを使用する
	var config *DBConfig
	if ok {
		// database.tomlを読取りDBConfigにマッピングする
		m := map[string]DBConfig{}
		_, err := toml.DecodeFile(path, &m)
		if err != nil {
			return nil, err
		}
		var mysqlConfig DBConfig
		mysqlConfig, ok = m["mysql"]
		if !ok {
			return nil, fmt.Errorf("key 'mysql' not found in config file: %s", path)
		}
		config = &mysqlConfig
	} else {
		// 環境変数が無い場合は環境変数から取得する
		c, err := loadConfigFromEnv()
		if err != nil {
			return nil, err
		}
		config = &c
	}
	return config, nil
}

func DBConnect() error {
	config, err := loadConfig()
	if err != nil {
		return DBErrHandler(err)
	}

	// 接続文字列を生成する
	rdbms := "mysql"
	connectStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", config.User, config.Pass, config.Host, config.Port, config.DBName)

	// DB接続
	conn, err := sql.Open(rdbms, connectStr)
	if err != nil {
		return DBErrHandler(err)
	}
	// 接続確認
	if err = conn.Ping(); err != nil {
		return DBErrHandler(err)
	}
	// 接続プールの設定
	conn.SetMaxIdleConns(config.MaxIdleConns)       // 最大アイドル接続数
	conn.SetMaxOpenConns(config.MaxOpenConns)       // 最大接続数
	conn.SetConnMaxLifetime(config.ConnMaxLifetime) // 接続の最大生存時間
	conn.SetConnMaxIdleTime(config.ConnMaxIdleTime) // 接続の最大アイドル時間

	// FIXME: DBConnect 関数でグローバルな状態を設定する代わりに、生成した *sql.DB コネクションを返し、それを必要とするリポジトリ層などのコンポーネントにコンストラクタ経由で注入（Inject）するようにリファクタリングしてください。
	// boil.SetDBはグローバルにDB接続を設定する
	boil.SetDB(conn)
	logLevel, err := utils.GetEnv("LOG_LEVEL", "debug")
	if err != nil {
		return err
	}
	boil.DebugMode = strings.ToLower(logLevel) == "debug" // デバッグモードに設定 生成されたSQLを出力する
	return nil
}
