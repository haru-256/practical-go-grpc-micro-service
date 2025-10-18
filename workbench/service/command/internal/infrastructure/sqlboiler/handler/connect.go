package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aarondl/sqlboiler/v4/boil"
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
	LogLevel        string        // ログレベル
}

// getKey はViperから型安全に設定値を取得するヘルパー関数です。
// 指定されたキーが存在しない場合、またはサポートされていない型の場合はエラーを記録します。
//
// サポートされる型:
//   - string
//   - int
//   - bool
//   - time.Duration
//
// Parameters:
//   - v: Viperインスタンス
//   - key: 設定キー（例: "mysql.host"）
//   - errs: エラーを蓄積するスライスへのポインタ
//
// Returns:
//   - T: 設定値（エラーの場合はゼロ値）
func getKey[T any](v *viper.Viper, key string, errs *[]error) T {
	var zero T
	if !v.IsSet(key) {
		*errs = append(*errs, fmt.Errorf("config key '%s' is not set", key))
		return zero
	}

	switch any(zero).(type) {
	case string:
		return any(v.GetString(key)).(T)
	case int:
		return any(v.GetInt(key)).(T)
	case bool:
		return any(v.GetBool(key)).(T)
	case time.Duration:
		// v.GetDuration() を使うことで、"30m" や "1h" のような文字列を
		// time.Duration型へ安全にパースする処理をViperに任せます。
		return any(v.GetDuration(key)).(T)
	default:
		*errs = append(*errs, fmt.Errorf("unsupported type for key '%s'", key))
		return zero
	}
}

// NewDBConfig はデータベース設定を生成します。
// 外部から注入されたviperインスタンスを使用してTOMLファイルまたは環境変数から設定を読み込みます。
//
// 環境変数（DB_プレフィックス）:
//   - DB_DBNAME, DB_HOST, DB_PORT, DB_USER, DB_PASS
//   - DB_MAX_IDLE_CONNS, DB_MAX_OPEN_CONNS
//   - DB_CONN_MAX_LIFETIME, DB_CONN_MAX_IDLE_TIME
//
// Parameters:
//   - v: viper設定インスタンス（通常はconfig.NewViper()から取得）
//
// Returns:
//   - *DBConfig: データベース設定
//   - error: パースエラー、または必須設定キーが存在しない場合
func NewDBConfig(v *viper.Viper) (*DBConfig, error) {
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
		LogLevel:        getKey[string](v, "log.level", &configErrors),
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
	boil.DebugMode = strings.ToLower(config.LogLevel) == "debug" // デバッグモードに設定 生成されたSQLを出力する
	return db, nil
}
