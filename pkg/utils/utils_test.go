package utils

import (
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"
)

func TestUtils(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Utils Suite")
}

var _ = Describe("GetKey", func() {
	var (
		v    *viper.Viper
		errs []error
	)

	BeforeEach(func() {
		v = viper.New()
		errs = []error{}
	})

	Describe("String型", func() {
		Context("キーが設定されている場合", func() {
			BeforeEach(func() {
				v.Set("test.string", "test_value")
			})

			It("設定値を返す", func() {
				result := GetKey[string](v, "test.string", &errs)
				Expect(result).To(Equal("test_value"))
				Expect(errs).To(BeEmpty())
			})
		})

		Context("キーが設定されていない場合", func() {
			It("ゼロ値を返し、エラーを記録する", func() {
				result := GetKey[string](v, "test.not_set", &errs)
				Expect(result).To(Equal(""))
				Expect(errs).To(HaveLen(1))
				Expect(errs[0].Error()).To(ContainSubstring("config key 'test.not_set' is not set"))
			})
		})
	})

	Describe("Int型", func() {
		Context("キーが設定されている場合", func() {
			BeforeEach(func() {
				v.Set("test.int", 42)
			})

			It("設定値を返す", func() {
				result := GetKey[int](v, "test.int", &errs)
				Expect(result).To(Equal(42))
				Expect(errs).To(BeEmpty())
			})
		})

		Context("キーが設定されていない場合", func() {
			It("ゼロ値を返し、エラーを記録する", func() {
				result := GetKey[int](v, "test.not_set", &errs)
				Expect(result).To(Equal(0))
				Expect(errs).To(HaveLen(1))
				Expect(errs[0].Error()).To(ContainSubstring("config key 'test.not_set' is not set"))
			})
		})
	})

	Describe("Bool型", func() {
		Context("キーが設定されている場合", func() {
			BeforeEach(func() {
				v.Set("test.bool", true)
			})

			It("設定値を返す", func() {
				result := GetKey[bool](v, "test.bool", &errs)
				Expect(result).To(Equal(true))
				Expect(errs).To(BeEmpty())
			})
		})

		Context("キーが設定されていない場合", func() {
			It("ゼロ値を返し、エラーを記録する", func() {
				result := GetKey[bool](v, "test.not_set", &errs)
				Expect(result).To(Equal(false))
				Expect(errs).To(HaveLen(1))
				Expect(errs[0].Error()).To(ContainSubstring("config key 'test.not_set' is not set"))
			})
		})
	})

	Describe("time.Duration型", func() {
		Context("キーが設定されている場合", func() {
			BeforeEach(func() {
				v.Set("test.duration", "30m")
			})

			It("設定値を返す", func() {
				result := GetKey[time.Duration](v, "test.duration", &errs)
				Expect(result).To(Equal(30 * time.Minute))
				Expect(errs).To(BeEmpty())
			})
		})

		Context("キーが設定されていない場合", func() {
			It("ゼロ値を返し、エラーを記録する", func() {
				result := GetKey[time.Duration](v, "test.not_set", &errs)
				Expect(result).To(Equal(time.Duration(0)))
				Expect(errs).To(HaveLen(1))
				Expect(errs[0].Error()).To(ContainSubstring("config key 'test.not_set' is not set"))
			})
		})
	})

	Describe("サポートされていない型", func() {
		Context("float64型を指定した場合", func() {
			BeforeEach(func() {
				v.Set("test.float", 3.14)
			})

			It("ゼロ値を返し、エラーを記録する", func() {
				result := GetKey[float64](v, "test.float", &errs)
				Expect(result).To(Equal(0.0))
				Expect(errs).To(HaveLen(1))
				Expect(errs[0].Error()).To(ContainSubstring("unsupported type for key 'test.float'"))
			})
		})
	})

	Describe("複数のエラーが蓄積される", func() {
		It("複数のキーでエラーが発生した場合、すべてのエラーが記録される", func() {
			GetKey[string](v, "test.not_set1", &errs)
			GetKey[int](v, "test.not_set2", &errs)
			GetKey[bool](v, "test.not_set3", &errs)

			Expect(errs).To(HaveLen(3))
			Expect(errs[0].Error()).To(ContainSubstring("config key 'test.not_set1' is not set"))
			Expect(errs[1].Error()).To(ContainSubstring("config key 'test.not_set2' is not set"))
			Expect(errs[2].Error()).To(ContainSubstring("config key 'test.not_set3' is not set"))
		})
	})
})
