package handler

import (
	"os"
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestConn(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "infrastructure/sqlboiler/handler packageのテスト")
}

var _ = Describe("DBConnect関数", Label("DBConnect"), func() {
	BeforeEach(func() {
		// テスト前のセットアップ処理
		absPath, err := filepath.Abs("../config/sqlboiler.toml")
		if err != nil {
			Fail("Failed to get absolute path")
		}
		Expect(os.Setenv("DATABASE_TOML_PATH", absPath)).To(Succeed())
	})

	DescribeTable("DBConnectの動作確認",
		func(expectError bool, expectedErrorCode string) {
			err := DBConnect()
			if expectError {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(expectedErrorCode))
			} else {
				Expect(err).NotTo(HaveOccurred())
			}
		},
		Entry("正常系", false, ""),
	)
})
