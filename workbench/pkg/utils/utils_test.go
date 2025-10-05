package utils

import (
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestUtils(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Utils Suite")
}

// ヘルパー関数: 文字列のポインタを返す
func stringPtr(s string) *string {
	return &s
}

var _ = Describe("GetEnv", func() {
	Describe("String型", func() {
		DescribeTable("環境変数のパース",
			func(envKey string, envValue *string, defaultValue string, expected string) {
				if envValue != nil {
					Expect(os.Setenv(envKey, *envValue)).To(Succeed())
					DeferCleanup(func() {
						Expect(os.Unsetenv(envKey)).To(Succeed())
					})
				}

				result, err := GetEnv(envKey, defaultValue)
				Expect(err).NotTo(HaveOccurred())
				Expect(result).To(Equal(expected))
			},
			Entry("環境変数が設定されている場合、環境変数の値を返す",
				"TEST_STRING_VAR", stringPtr("test_value"), "default", "test_value"),
			Entry("環境変数が設定されていない場合、デフォルト値を返す",
				"TEST_STRING_VAR_NOT_SET", nil, "default_value", "default_value"),
			Entry("空文字列が設定されている場合、空文字列を返す",
				"TEST_STRING_EMPTY", stringPtr(""), "default", ""),
		)
	})

	Describe("Int型", func() {
		DescribeTable("環境変数のパース",
			func(envKey string, envValue *string, defaultValue int, expected int) {
				if envValue != nil {
					Expect(os.Setenv(envKey, *envValue)).To(Succeed())
					DeferCleanup(func() {
						Expect(os.Unsetenv(envKey)).To(Succeed())
					})
				}

				result, err := GetEnv(envKey, defaultValue)
				Expect(err).NotTo(HaveOccurred())
				Expect(result).To(Equal(expected))
			},
			Entry("環境変数が設定されている場合、環境変数の値を返す",
				"TEST_INT_VAR", stringPtr("42"), 0, 42),
			Entry("環境変数が設定されていない場合、デフォルト値を返す",
				"TEST_INT_VAR_NOT_SET", nil, 100, 100),
			Entry("負の数が設定されている場合、負の数を返す",
				"TEST_INT_NEGATIVE", stringPtr("-10"), 0, -10),
		)

		Context("不正な値が設定されている場合", func() {
			const key = "TEST_INT_INVALID"

			BeforeEach(func() {
				Expect(os.Setenv(key, "invalid")).To(Succeed())
			})

			AfterEach(func() {
				Expect(os.Unsetenv(key)).To(Succeed())
			})

			It("エラーを返す", func() {
				result, err := GetEnv(key, 0)
				Expect(err).To(HaveOccurred())
				Expect(result).To(Equal(0)) // デフォルト値を返す
			})
		})
	})

	Describe("Bool型", func() {
		DescribeTable("環境変数のパース",
			func(envKey string, envValue *string, defaultValue bool, expected bool) {
				if envValue != nil {
					Expect(os.Setenv(envKey, *envValue)).To(Succeed())
					DeferCleanup(func() {
						Expect(os.Unsetenv(envKey)).To(Succeed())
					})
				}

				result, err := GetEnv(envKey, defaultValue)
				Expect(err).NotTo(HaveOccurred())
				Expect(result).To(Equal(expected))
			},
			Entry("trueが設定されている場合、trueを返す",
				"TEST_BOOL_TRUE", stringPtr("true"), false, true),
			Entry("falseが設定されている場合、falseを返す",
				"TEST_BOOL_FALSE", stringPtr("false"), true, false),
			Entry("1が設定されている場合、trueを返す",
				"TEST_BOOL_ONE", stringPtr("1"), false, true),
			Entry("0が設定されている場合、falseを返す",
				"TEST_BOOL_ZERO", stringPtr("0"), true, false),
			Entry("環境変数が設定されていない場合、デフォルト値を返す",
				"TEST_BOOL_VAR_NOT_SET", nil, true, true),
		)

		Context("不正な値が設定されている場合", func() {
			const key = "TEST_BOOL_INVALID"

			BeforeEach(func() {
				Expect(os.Setenv(key, "invalid")).To(Succeed())
			})

			AfterEach(func() {
				Expect(os.Unsetenv(key)).To(Succeed())
			})

			It("エラーを返す", func() {
				result, err := GetEnv(key, false)
				Expect(err).To(HaveOccurred())
				Expect(result).To(Equal(false)) // デフォルト値を返す
			})
		})
	})

	Describe("Float64型", func() {
		DescribeTable("環境変数のパース",
			func(envKey string, envValue *string, defaultValue float64, expected float64) {
				if envValue != nil {
					Expect(os.Setenv(envKey, *envValue)).To(Succeed())
					DeferCleanup(func() {
						Expect(os.Unsetenv(envKey)).To(Succeed())
					})
				}

				result, err := GetEnv(envKey, defaultValue)
				Expect(err).NotTo(HaveOccurred())
				Expect(result).To(Equal(expected))
			},
			Entry("環境変数が設定されている場合、環境変数の値を返す",
				"TEST_FLOAT_VAR", stringPtr("3.14"), 0.0, 3.14),
			Entry("環境変数が設定されていない場合、デフォルト値を返す",
				"TEST_FLOAT_VAR_NOT_SET", nil, 2.71, 2.71),
			Entry("負の数が設定されている場合、負の数を返す",
				"TEST_FLOAT_NEGATIVE", stringPtr("-1.5"), 0.0, -1.5),
			Entry("整数が設定されている場合、整数を浮動小数点数として返す",
				"TEST_FLOAT_INTEGER", stringPtr("42"), 0.0, 42.0),
		)

		Context("不正な値が設定されている場合", func() {
			const key = "TEST_FLOAT_INVALID"

			BeforeEach(func() {
				Expect(os.Setenv(key, "invalid")).To(Succeed())
			})

			AfterEach(func() {
				Expect(os.Unsetenv(key)).To(Succeed())
			})

			It("エラーを返す", func() {
				result, err := GetEnv(key, 0.0)
				Expect(err).To(HaveOccurred())
				Expect(result).To(Equal(0.0)) // デフォルト値を返す
			})
		})
	})
})
