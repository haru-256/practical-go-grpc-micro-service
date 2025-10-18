//go:build integration || !ci

package repository

import (
	"testing"

	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/config"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/infrastructure/sqlboiler/handler"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestRepImplPackage(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Repository Implementation Suite")
}

var _ = BeforeSuite(func() {
	v := config.NewViper("../../../../", "config")
	config, err := handler.NewDBConfig(v)
	Expect(err).NotTo(HaveOccurred(), "DBConfigの生成に失敗しました。")
	_, err = handler.NewDatabase(config)
	Expect(err).NotTo(HaveOccurred(), "データベース接続が失敗したのでテストを中止します。")
})
