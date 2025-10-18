//go:build integration || !ci

// Package repository_test provides integration tests for the transaction manager implementation.
// These tests verify the transactionManagerImpl's Begin and Complete methods
// using a real database connection.
package repository

import (
	"context"
	"database/sql"
	"errors"
	"io"
	"log/slog"

	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/application/service"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("transactionManagerImpl", Ordered, Label("TransactionManagerインターフェースメソッドのテスト"), func() {
	var (
		ctx context.Context
		tm  service.TransactionManager
	)

	BeforeAll(func() {
		// テストではログ出力を破棄するloggerを使用
		logger := slog.New(slog.NewTextHandler(io.Discard, nil))
		tm = NewTransactionManagerImpl(logger)
	})

	BeforeEach(func() {
		ctx = context.Background()
	})

	Context("Beginメソッドの動作確認", func() {
		It("新しいトランザクションを開始できること", func() {
			tx, err := tm.Begin(ctx)
			defer func() {
				if tx != nil {
					// テストクリーンアップのため強制ロールバック
					_ = tx.Rollback()
				}
			}()

			Expect(err).NotTo(HaveOccurred(), "トランザクションの開始に失敗しました")
			Expect(tx).NotTo(BeNil(), "開始したトランザクションがnilです")
		})
	})

	Context("Completeメソッドの動作確認", func() {
		var tx *sql.Tx

		BeforeEach(func() {
			var err error
			tx, err = tm.Begin(ctx)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("エラーが渡された場合", func() {
			It("トランザクションをロールバックすること", func() {
				err := tm.Complete(ctx, tx, errors.New("何らかのエラー"))
				Expect(err).NotTo(HaveOccurred(), "ロールバックに失敗しました")
			})
		})

		Context("エラーがnilの場合", func() {
			It("トランザクションをコミットすること", func() {
				err := tm.Complete(ctx, tx, nil)
				Expect(err).NotTo(HaveOccurred(), "コミットに失敗しました")
			})
		})
	})

	Context("トランザクションのライフサイクル", func() {
		It("Begin→ロールバックの一連の流れが正常に動作すること", func() {
			tx, err := tm.Begin(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(tx).NotTo(BeNil())

			err = tm.Complete(ctx, tx, errors.New("テスト用エラー"))
			Expect(err).NotTo(HaveOccurred())
		})

		It("Begin→コミットの一連の流れが正常に動作すること", func() {
			tx, err := tm.Begin(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(tx).NotTo(BeNil())

			err = tm.Complete(ctx, tx, nil)
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
