package impl

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/application/service"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/domain/models/categories"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/errs"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("CategoryService", Label("UnitTests"), func() {
	var (
		ctrl         *gomock.Controller
		mockRepo     *categories.MockCategoryRepository
		mockTm       *service.MockTransactionManager
		cs           service.CategoryService
		ctx          context.Context
		mockTx       *sql.Tx
		testCategory *categories.Category
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockRepo = categories.NewMockCategoryRepository(ctrl)
		mockTm = service.NewMockTransactionManager(ctrl)
		cs = NewCategoryServiceImpl(mockRepo, mockTm)
		ctx = context.Background()

		// モックトランザクション（実際のオブジェクトをシミュレート）
		mockTx = &sql.Tx{}

		// テスト用カテゴリの作成（ユニットテストではユニークである必要はない）
		name, err := categories.NewCategoryName("TestCategory")
		Expect(err).NotTo(HaveOccurred())
		testCategory, err = categories.NewCategory(name)
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Describe("Add", func() {
		Context("when category name does not exist", func() {
			It("should successfully add a new category", func() {
				// Arrange: モックの期待値を順序付きで設定
				gomock.InOrder(
					mockTm.EXPECT().Begin(ctx).Return(mockTx, nil),
					mockRepo.EXPECT().ExistsByName(ctx, mockTx, testCategory.Name()).Return(false, nil),
					mockRepo.EXPECT().Create(ctx, mockTx, testCategory).Return(nil),
					mockTm.EXPECT().Complete(mockTx, nil).Return(nil),
				)

				// Act
				err := cs.Add(ctx, testCategory)

				// Assert
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when category name already exists", func() {
			It("should return ApplicationError with CATEGORY_ALREADY_EXISTS code", func() {
				// Arrange
				gomock.InOrder(
					mockTm.EXPECT().Begin(ctx).Return(mockTx, nil),
					mockRepo.EXPECT().ExistsByName(ctx, mockTx, testCategory.Name()).Return(true, nil),
					mockTm.EXPECT().Complete(mockTx, gomock.Any()).
						Do(func(tx *sql.Tx, err error) {
							// Completeに渡されるエラーがApplicationErrorであることを検証
							Expect(err).To(BeAssignableToTypeOf(&errs.ApplicationError{}))
							appErr := err.(*errs.ApplicationError)
							Expect(appErr.Code).To(Equal("CATEGORY_ALREADY_EXISTS"))
						}).
						Return(nil),
				)

				// Act
				err := cs.Add(ctx, testCategory)

				// Assert
				Expect(err).To(HaveOccurred())
				Expect(err).To(BeAssignableToTypeOf(&errs.ApplicationError{}))
				appErr := err.(*errs.ApplicationError)
				Expect(appErr.Code).To(Equal("CATEGORY_ALREADY_EXISTS"))
			})
		})

		Context("when Begin fails", func() {
			It("should return the error from Begin", func() {
				// Arrange
				beginErr := fmt.Errorf("failed to begin transaction")
				mockTm.EXPECT().Begin(ctx).Return(nil, beginErr)

				// Act
				err := cs.Add(ctx, testCategory)

				// Assert
				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(beginErr))
			})
		})

		Context("when ExistsByName fails", func() {
			It("should return the error and rollback", func() {
				// Arrange
				existsErr := fmt.Errorf("database error")
				gomock.InOrder(
					mockTm.EXPECT().Begin(ctx).Return(mockTx, nil),
					mockRepo.EXPECT().ExistsByName(ctx, mockTx, testCategory.Name()).Return(false, existsErr),
					mockTm.EXPECT().Complete(mockTx, existsErr).Return(nil),
				)

				// Act
				err := cs.Add(ctx, testCategory)

				// Assert
				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(existsErr))
			})
		})

		Context("when Create fails", func() {
			It("should return the error and rollback", func() {
				// Arrange
				createErr := errs.NewCRUDError("DB_ERROR", "failed to create category")
				gomock.InOrder(
					mockTm.EXPECT().Begin(ctx).Return(mockTx, nil),
					mockRepo.EXPECT().ExistsByName(ctx, mockTx, testCategory.Name()).Return(false, nil),
					mockRepo.EXPECT().Create(ctx, mockTx, testCategory).Return(createErr),
					mockTm.EXPECT().Complete(mockTx, createErr).Return(nil),
				)

				// Act
				err := cs.Add(ctx, testCategory)

				// Assert
				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(createErr))
			})
		})
	})

	Describe("Update", func() {
		Context("when category exists", func() {
			It("should successfully update the category", func() {
				// Arrange
				gomock.InOrder(
					mockTm.EXPECT().Begin(ctx).Return(mockTx, nil),
					mockRepo.EXPECT().UpdateById(ctx, mockTx, testCategory).Return(nil),
					mockTm.EXPECT().Complete(mockTx, nil).Return(nil),
				)

				// Act
				err := cs.Update(ctx, testCategory)

				// Assert
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when UpdateById fails", func() {
			It("should return the error and rollback", func() {
				// Arrange
				updateErr := errs.NewCRUDError("NOT_FOUND", "category not found")
				gomock.InOrder(
					mockTm.EXPECT().Begin(ctx).Return(mockTx, nil),
					mockRepo.EXPECT().UpdateById(ctx, mockTx, testCategory).Return(updateErr),
					mockTm.EXPECT().Complete(mockTx, updateErr).Return(nil),
				)

				// Act
				err := cs.Update(ctx, testCategory)

				// Assert
				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(updateErr))
			})
		})

		Context("when Begin fails", func() {
			It("should return the error from Begin", func() {
				// Arrange
				beginErr := fmt.Errorf("failed to begin transaction")
				mockTm.EXPECT().Begin(ctx).Return(nil, beginErr)

				// Act
				err := cs.Update(ctx, testCategory)

				// Assert
				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(beginErr))
			})
		})
	})

	Describe("Delete", func() {
		Context("when category exists", func() {
			It("should successfully delete the category", func() {
				// Arrange
				gomock.InOrder(
					mockTm.EXPECT().Begin(ctx).Return(mockTx, nil),
					mockRepo.EXPECT().DeleteById(ctx, mockTx, testCategory.Id()).Return(nil),
					mockTm.EXPECT().Complete(mockTx, nil).Return(nil),
				)

				// Act
				err := cs.Delete(ctx, testCategory)

				// Assert
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when DeleteById fails", func() {
			It("should return the error and rollback", func() {
				// Arrange
				deleteErr := errs.NewCRUDError("NOT_FOUND", "category not found")
				gomock.InOrder(
					mockTm.EXPECT().Begin(ctx).Return(mockTx, nil),
					mockRepo.EXPECT().DeleteById(ctx, mockTx, testCategory.Id()).Return(deleteErr),
					mockTm.EXPECT().Complete(mockTx, deleteErr).Return(nil),
				)

				// Act
				err := cs.Delete(ctx, testCategory)

				// Assert
				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(deleteErr))
			})
		})

		Context("when Begin fails", func() {
			It("should return the error from Begin", func() {
				// Arrange
				beginErr := fmt.Errorf("failed to begin transaction")
				mockTm.EXPECT().Begin(ctx).Return(nil, beginErr)

				// Act
				err := cs.Delete(ctx, testCategory)

				// Assert
				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(beginErr))
			})
		})
	})
})
