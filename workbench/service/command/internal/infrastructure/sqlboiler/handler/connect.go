package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/haru-256/practical-go-grpc-micro-service/pkg/utils"
	"github.com/spf13/viper"
)

// DBConfig はデータベース接続設定を保持する構造体です。
// TOMLファイルまたは環境変数から読み込まれます。
type DBConfig struct {
	DBName          string        //	データベース名
	Host            string        //	ホスト名
	Port            int           //	ポート番号
	User            string        //	ユーザー名
	Pass            string        //	パスワード
	MaxIdleConns    int           //	最大アイドル接続数
	MaxOpenConns    int           //	最大接続数
	ConnMaxLifetime time.Duration //	接続の最大生存時間(分)
	ConnMaxIdleTime time.Duration //	接続の最大アイドル時間(分)
}

func setupViper(configPath string, configName string) *viper.Viper {
	v := viper.New()
	v.AddConfigPath(configPath)
	v.SetConfigName(configName)
	v.SetConfigType("toml")

	// 環境変数でtomlの設定を上書き可能にする
	v.AutomaticEnv()

	// 各設定キーを環境変数に明示的にバインド
	bindings := map[string]string{
		"mysql.dbname":             "DB_DBNAME",
		"mysql.host":               "DB_HOST",
		"mysql.port":               "DB_PORT",
		"mysql.user":               "DB_USER",
		"mysql.pass":               "DB_PASS",
		"mysql.max_idle_conns":     "DB_MAX_IDLE_CONNS",
		"mysql.max_open_conns":     "DB_MAX_OPEN_CONNS",
		"mysql.conn_max_lifetime":  "DB_CONN_MAX_LIFETIME",
		"mysql.conn_max_idle_time": "DB_CONN_MAX_IDLE_TIME",
	}
	for key, env := range bindings {
		if err := v.BindEnv(key, env); err != nil {
			log.Printf("failed to bind env for %s: %v", key, err)
		}
	}

	if err := v.ReadInConfig(); err != nil {
		log.Printf("config file not found: %v", err)
	}

	return v
}

func getKey[T any](v *viper.Viper, key string, errs *[]error) T {
	var zero T
	if !v.IsSet(key) {
		*errs = append(*errs, fmt.Errorf("config key '%s' is not set", key))
		return zero
	}

	switch any(zero).(type) {
	case string:
		// durationなど文字列を特定の型にparseするのはviperに任せたいため、
		// any()を使ってGetString()の結果をinterface{}に変換し、
		// さらに.(T)で元の型Tにアサーション（変換）する
		return any(v.GetString(key)).(T)
	case int:
		return any(v.GetInt(key)).(T)
	case bool:
		return any(v.GetBool(key)).(T)
	case time.Duration:
		return any(v.GetDuration(key)).(T)
	default:
		*errs = append(*errs, fmt.Errorf("unsupported type for key '%s'", key))
		return zero
	}
}

// NewDBConfig はデータベース設定を生成します。
// TOMLファイルまたは環境変数から読み込みます。
// デフォルトでは "../config" から "database" という名前の設定ファイルを読み込みます。
//
// 環境変数（DB_プレフィックス）:
//   - DB_DBNAME, DB_HOST, DB_PORT, DB_USER, DB_PASS
//   - DB_MAX_IDLE_CONNS, DB_MAX_OPEN_CONNS
//   - DB_CONN_MAX_LIFETIME, DB_CONN_MAX_IDLE_TIME
//
// Returns:
//   - *DBConfig: データベース設定
//   - error: ファイル読み込みエラー、パースエラー、または環境変数エラー
func NewDBConfig() (*DBConfig, error) {
	v := setupViper("../config", "database")
	var configErrors []error
	cfg := DBConfig{
		DBName:          getKey[string](v, "mysql.dbname", &configErrors),
		Host:            getKey[string](v, "mysql.host", &configErrors),
		Port:            getKey[int](v, "mysql.port", &configErrors),
		User:            getKey[string](v, "mysql.user", &configErrors),
		Pass:            getKey[string](v, "mysql.pass", &configErrors),
		MaxIdleConns:    getKey[int](v, "mysql.max_idle_conns", &configErrors),
		MaxOpenConns:    getKey[int](v, "mysql.max_open_conns", &configErrors),
		ConnMaxLifetime: getKey[time.Duration](v, "mysql.conn_max_lifetime", &configErrors),
		ConnMaxIdleTime: getKey[time.Duration](v, "mysql.conn_max_idle_time", &configErrors),
	}
	// すべての環境変数を読み込んだ後、エラーがあればまとめて返す
	if len(configErrors) > 0 {
		return &cfg, errors.Join(configErrors...)
	}
	return &cfg, nil
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
		_ = db.Close() // エラーハンドリングよりもPingエラーを優先
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
