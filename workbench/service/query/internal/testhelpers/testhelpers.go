package testhelpers

import (
	"io"
	"log/slog"

	"buf.build/go/protovalidate"
	"github.com/haru-256/practical-go-grpc-micro-service/service/query/internal/infrastructure/config"
	"github.com/haru-256/practical-go-grpc-micro-service/service/query/internal/infrastructure/db"
	"gorm.io/gorm"
)

// TestLogger はテスト用のロガーです。出力は破棄されます。
var TestLogger = slog.New(slog.NewTextHandler(io.Discard, nil))

// TestValidator はテスト用のProtobufバリデータです。
var TestValidator = func() protovalidate.Validator {
	v, err := protovalidate.New()
	if err != nil {
		panic(err)
	}
	return v
}()

// SetupDB はテスト用のデータベース接続を確立します。
//
// Parameters:
//   - configPath: 設定ファイルのディレクトリパス
//   - configName: 設定ファイル名（拡張子なし）
//
// Returns:
//   - *gorm.DB: データベース接続
//   - error: エラー
func SetupDB(configPath, configName string) (*gorm.DB, error) {
	v := config.NewViper(configPath, configName)
	dbConfig, err := db.NewDBConfig(v)
	if err != nil {
		return nil, err
	}
	dbConn, err := db.NewDatabase(dbConfig, TestLogger)
	if err != nil {
		return nil, err
	}
	return dbConn, nil
}

// TeardownDB はデータベース接続をクローズします。
//
// Parameters:
//   - dbConn: データベース接続
//
// Returns:
//   - error: エラー
func TeardownDB(dbConn *gorm.DB) error {
	if dbConn == nil {
		return nil
	}
	sqlDB, err := dbConn.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// IntegrationTestSetup は統合テスト用の共通セットアップを保持します。
type IntegrationTestSetup struct {
	DBConn *gorm.DB     // データベース接続
	Logger *slog.Logger // ロガー
}
