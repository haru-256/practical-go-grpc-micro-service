package config

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

var _ = Describe("Config Module", func() {
	configOption := fx.Supply(
		fx.Annotate("../../", fx.ResultTags(`name:"configPath"`)),
		fx.Annotate("config", fx.ResultTags(`name:"configName"`)),
	)
	Context("when creating an fx app with the config module", func() {
		It("should successfully initialize without errors", func() {
			var v *viper.Viper

			app := fx.New(
				// configPathとconfigNameを提供
				configOption,
				Module,
				fx.Populate(&v),
				fx.NopLogger, // テスト時はログを抑制
			)

			Expect(app.Err()).ToNot(HaveOccurred(), "fx app should initialize without errors")
			Expect(v).ToNot(BeNil(), "viper instance should be provided")
		})

		It("should provide a properly configured viper instance", func() {
			var v *viper.Viper

			app := fx.New(
				configOption,
				Module,
				fx.Populate(&v),
				fx.NopLogger,
			)

			Expect(app.Err()).ToNot(HaveOccurred())
			Expect(v).ToNot(BeNil())

			// Viperインスタンスが使用可能であることを確認
			// configPathとconfigNameが提供されているため、
			// 何らかの設定が読み込まれているはず
			Expect(v.AllKeys()).ToNot(BeEmpty(), "viper should have loaded configuration")
		})
	})
})
