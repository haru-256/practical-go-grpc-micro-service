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

// DBConfig はデータベース接続設定を保持する構造体です。
// TOMLファイルまたは環境変数から読み込まれます。
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

// loadConfigFromEnv は環境変数からデータベース設定を読み込みます。
// 環境変数が設定されていない場合は、デフォルト値が使用されます。
//
// 環境変数:
//   - DB_NAME: データベース名（デフォルト: "sample_db"）
//   - DB_HOST: ホスト名（デフォルト: "localhost"）
//   - DB_PORT: ポート番号（デフォルト: 3306）
//   - DB_USER: ユーザー名（デフォルト: "root"）
//   - DB_PASSWORD: パスワード（デフォルト: "password"）
//   - DB_MAX_IDLE_CONNS: 最大アイドル接続数（デフォルト: 10）
//   - DB_MAX_OPEN_CONNS: 最大接続数（デフォルト: 100）
//   - DB_CONN_MAX_LIFETIME: 接続の最大生存時間（デフォルト: 30分）
//   - DB_CONN_MAX_IDLE_TIME: 接続の最大アイドル時間（デフォルト: 5秒）
//
// Returns:
//   - DBConfig: 読み込まれた設定（エラーがある場合もデフォルト値で設定されます）
//   - error: 環境変数の解析エラーがある場合、すべてのエラーを結合したもの
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

// NewDBConfig はデータベース設定を生成します。
// DATABASE_TOML_PATH環境変数が設定されている場合はTOMLファイルから読み込み、
// 設定されていない場合は環境変数から読み込みます。
//
// 環境変数:
//   - DATABASE_TOML_PATH: TOMLファイルのパス（任意）
//
// Returns:
//   - *DBConfig: データベース設定
//   - error: ファイル読み込みエラー、パースエラー、または環境変数エラー
func NewDBConfig() (*DBConfig, error) {
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

// NewDatabase は指定された設定でMySQLデータベース接続を確立します。
// この関数はSQLBoilerのグローバル状態も設定します。
//
// Note: この関数はboil.SetDB()を通じてグローバルなDB接続を設定します。
// これはSQLBoilerの設計上の要件によるもので、以下の理由から意図的な設計選択です:
//   - SQLBoilerの生成コードはグローバルに設定されたDB接続を期待します
//   - リポジトリメソッドは*sql.Txを引数で受け取り、トランザクション境界を明示的に制御します
//   - グローバルDBはトランザクション生成のみに使用され、実際のクエリは注入された*sql.Txを使用します
//
// Parameters:
//   - config: データベース設定
//
// Returns:
//   - *sql.DB: データベース接続（呼び出し元でクローズする必要があります）
//   - error: 接続エラーまたは設定エラー
func NewDatabase(config *DBConfig) (*sql.DB, error) {
	// 接続文字列を生成する
	rdbms := "mysql"
	connectStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", config.User, config.Pass, config.Host, config.Port, config.DBName)

	// DB接続
	db, err := sql.Open(rdbms, connectStr)
	if err != nil {
		return nil, DBErrHandler(err)
	}
	// 接続確認
	if err = db.Ping(); err != nil {
		defer db.Close()
		return nil, DBErrHandler(err)
	}
	// 接続プールの設定
	db.SetMaxIdleConns(config.MaxIdleConns)       // 最大アイドル接続数
	db.SetMaxOpenConns(config.MaxOpenConns)       // 最大接続数
	db.SetConnMaxLifetime(config.ConnMaxLifetime) // 接続の最大生存時間
	db.SetConnMaxIdleTime(config.ConnMaxIdleTime) // 接続の最大アイドル時間

	// Configure SQLBoiler globally (required by SQLBoiler's design)
	boil.SetDB(db)
	logLevel, err := utils.GetEnv("LOG_LEVEL", "debug")
	if err != nil {
		return nil, err
	}
	boil.DebugMode = strings.ToLower(logLevel) == "debug" // デバッグモードに設定 生成されたSQLを出力する
	return db, nil
}
