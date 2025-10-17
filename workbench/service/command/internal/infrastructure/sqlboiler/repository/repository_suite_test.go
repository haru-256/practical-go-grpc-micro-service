//go:build integration || !ci

package repository_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/infrastructure/sqlboiler/handler"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestRepImplPackage(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Repository Implementation Suite")
}

var _ = BeforeSuite(func() {
	absPath, _ := filepath.Abs("../config/database.toml")
	Expect(os.Setenv("DATABASE_TOML_PATH", absPath)).To(Succeed())
	config, err := handler.NewDBConfig()
	Expect(err).NotTo(HaveOccurred(), "DBConfigの生成に失敗しました。")
	_, err = handler.NewDatabase(config)
	Expect(err).NotTo(HaveOccurred(), "データベース接続が失敗したのでテストを中止します。")
})
