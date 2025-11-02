package handler

import (
	"database/sql"
	"errors"
	"fmt"
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
	LogLevel        string        // ログレベル
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
		DBName:          utils.GetKey[string](v, "mysql.dbname", &configErrors),
		Host:            utils.GetKey[string](v, "mysql.host", &configErrors),
		Port:            utils.GetKey[int](v, "mysql.port", &configErrors),
		User:            utils.GetKey[string](v, "mysql.user", &configErrors),
		Pass:            utils.GetKey[string](v, "mysql.pass", &configErrors),
		MaxIdleConns:    utils.GetKey[int](v, "mysql.max_idle_conns", &configErrors),
		MaxOpenConns:    utils.GetKey[int](v, "mysql.max_open_conns", &configErrors),
		ConnMaxLifetime: utils.GetKey[time.Duration](v, "mysql.conn_max_lifetime", &configErrors),
		ConnMaxIdleTime: utils.GetKey[time.Duration](v, "mysql.conn_max_idle_time", &configErrors),
		LogLevel:        utils.GetKey[string](v, "log.level", &configErrors),
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
