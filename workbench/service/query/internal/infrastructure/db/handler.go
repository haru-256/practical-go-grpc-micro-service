package db

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"strings"
	"time"

	mysql_go "github.com/go-sql-driver/mysql"
	"github.com/haru-256/practical-go-grpc-micro-service/pkg/errs"
	"github.com/haru-256/practical-go-grpc-micro-service/pkg/utils"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gorm_logger "gorm.io/gorm/logger"
)

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

func NewDBConfig(v *viper.Viper) (*DBConfig, error) {
	var configErrors []error
	cfg := &DBConfig{
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
		return cfg, errors.Join(configErrors...)
	}
	return cfg, nil
}

func NewDatabase(config *DBConfig, logger *slog.Logger) (*gorm.DB, error) {
	ctx := context.Background()
	connectStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.User, config.Pass, config.Host, config.Port, config.DBName)
	conn, err := gorm.Open(mysql.Open(connectStr), &gorm.Config{})
	if err != nil {
		return nil, DBErrHandler(ctx, err, logger)
	}
	if db, err := conn.DB(); err != nil {
		return nil, DBErrHandler(ctx, err, logger)
	} else {
		if err := db.Ping(); err != nil {
			return nil, DBErrHandler(ctx, err, logger)
		}
		// 接続プールの設定
		db.SetMaxIdleConns(config.MaxIdleConns)       // 最大アイドル接続数
		db.SetMaxOpenConns(config.MaxOpenConns)       // 最大接続数
		db.SetConnMaxLifetime(config.ConnMaxLifetime) // 接続の最大生存時間
		db.SetConnMaxIdleTime(config.ConnMaxIdleTime) // 接続の最大アイドル時間

		// 生成されたSQLをログに出力する設定
		// 0: Silent, 1: Error, 2: Warn, 3: Info
		var loggerLevel gorm_logger.LogLevel
		switch strings.ToLower(config.LogLevel) {
		case "debug":
			loggerLevel = gorm_logger.Info
		case "info":
			loggerLevel = gorm_logger.Info
		case "warn":
			loggerLevel = gorm_logger.Warn
		case "error":
			loggerLevel = gorm_logger.Error
		default:
			return nil, errs.NewInternalErrorWithCause("INVALID_LOG_LEVEL", fmt.Sprintf("不正なログレベルです: %s", config.LogLevel), nil)
		}
		conn.Logger = conn.Logger.LogMode(loggerLevel)
	}

	return conn, nil
}

// DBErrHandler はデータベースアクセスエラーを適切なドメインエラーに変換します。
//
// この関数は以下のエラータイプを処理します:
//   - *net.OpError: ネットワーク接続エラー（接続タイムアウト等）
//   - *mysql.MySQLError: MySQLドライバ固有のエラー
//   - 1062: 一意制約違反
//   - その他: ドライバエラー
//   - その他: 不明なエラー
//
// Parameters:
//   - err: データベース操作から返されたエラー
//
// Returns:
//   - error: ドメイン層のエラー型に変換されたエラー
func DBErrHandler(ctx context.Context, err error, logger *slog.Logger) error {
	var opErr *net.OpError
	var driverErr *mysql_go.MySQLError
	if errors.As(err, &opErr) { // 接続タイムアウトやネットワーク関連の問題で接続が確立できない場合
		// TODO: slogに変更する
		logger.ErrorContext(ctx, opErr.Error())
		return errs.NewInternalErrorWithCause("DB_CONNECTION_ERROR", opErr.Error(), opErr)
	} else if errors.As(err, &driverErr) { // MySQLドライバエラーの場合
		logger.WarnContext(ctx, "MySQL driver error", slog.Int("code", int(driverErr.Number)), slog.String("message", driverErr.Message))
		return errs.NewInternalErrorWithCause("DB_DRIVER_ERROR", driverErr.Message, driverErr)
	} else { // その他のエラー
		logger.ErrorContext(ctx, err.Error())
		return errs.NewInternalErrorWithCause("UNKNOWN_ERROR", err.Error(), err)
	}
}
