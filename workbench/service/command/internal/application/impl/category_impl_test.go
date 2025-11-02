package impl

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log/slog"

	"github.com/haru-256/practical-go-grpc-micro-service/pkg/errs"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/application/dto"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/application/service"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/domain/models/categories"
	mock_repository "github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/mock/repository"
	mock_service "github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/mock/service"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("CategoryService", Label("UnitTests"), func() {
	var (
		ctrl         *gomock.Controller
		mockRepo     *mock_repository.MockCategoryRepository
		mockTm       *mock_service.MockTransactionManager
		cs           service.CategoryService
		ctx          context.Context
		mockTx       *sql.Tx
		testCategory *categories.Category
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockRepo = mock_repository.NewMockCategoryRepository(ctrl)
		mockTm = mock_service.NewMockTransactionManager(ctrl)
		// テストではログ出力を破棄するloggerを使用
		logger := slog.New(slog.NewTextHandler(io.Discard, nil))
		cs = NewCategoryServiceImpl(logger, mockRepo, mockTm)
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
				// Arrange
				createDTO := &dto.CreateCategoryDTO{
					Name: "TestCategory",
				}

				// Arrange: モックの期待値を順序付きで設定
				gomock.InOrder(
					mockTm.EXPECT().Begin(ctx).Return(mockTx, nil),
					mockRepo.EXPECT().ExistsByName(ctx, mockTx, testCategory.Name()).Return(false, nil),
					mockRepo.EXPECT().Create(ctx, mockTx, gomock.Any()).Do(
						func(ctx context.Context, tx *sql.Tx, category *categories.Category) {
							Expect(category).NotTo(BeNil())
							Expect(category.Name()).To(Equal(testCategory.Name()))
						}).Return(nil),
					mockTm.EXPECT().Complete(ctx, mockTx, nil).Return(nil),
				)

				// Act
				result, err := cs.Add(ctx, createDTO)

				// Assert
				Expect(err).NotTo(HaveOccurred())
				Expect(result).NotTo(BeNil())
				Expect(result.Name).To(Equal("TestCategory"))
			})
		})

		Context("when category name already exists", func() {
			It("should return ApplicationError with CATEGORY_ALREADY_EXISTS code", func() {
				// Arrange
				createDTO := &dto.CreateCategoryDTO{
					Name: "TestCategory",
				}

				gomock.InOrder(
					mockTm.EXPECT().Begin(ctx).Return(mockTx, nil),
					mockRepo.EXPECT().ExistsByName(ctx, mockTx, testCategory.Name()).Return(true, nil),
					mockTm.EXPECT().Complete(ctx, mockTx, gomock.Any()).
						Do(func(ctx context.Context, tx *sql.Tx, err error) {
							// Completeに渡されるエラーがApplicationErrorであることを検証
							Expect(err).To(BeAssignableToTypeOf(&errs.ApplicationError{}))
							appErr := err.(*errs.ApplicationError)
							Expect(appErr.Code).To(Equal("CATEGORY_ALREADY_EXISTS"))
						}).
						Return(nil),
				)

				// Act
				result, err := cs.Add(ctx, createDTO)

				// Assert
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				Expect(err).To(BeAssignableToTypeOf(&errs.ApplicationError{}))
				appErr := err.(*errs.ApplicationError)
				Expect(appErr.Code).To(Equal("CATEGORY_ALREADY_EXISTS"))
			})
		})

		Context("when Begin fails", func() {
			It("should return the error from Begin", func() {
				// Arrange
				createDTO := &dto.CreateCategoryDTO{
					Name: "TestCategory",
				}
				beginErr := fmt.Errorf("failed to begin transaction")
				mockTm.EXPECT().Begin(ctx).Return(nil, beginErr)

				// Act
				result, err := cs.Add(ctx, createDTO)

				// Assert
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				Expect(err).To(Equal(beginErr))
			})
		})

		Context("when ExistsByName fails", func() {
			It("should return the error and rollback", func() {
				// Arrange
				createDTO := &dto.CreateCategoryDTO{
					Name: "TestCategory",
				}
				existsErr := fmt.Errorf("database error")
				gomock.InOrder(
					mockTm.EXPECT().Begin(ctx).Return(mockTx, nil),
					mockRepo.EXPECT().ExistsByName(ctx, mockTx, testCategory.Name()).Return(false, existsErr),
					mockTm.EXPECT().Complete(ctx, mockTx, existsErr).Return(nil),
				)

				// Act
				result, err := cs.Add(ctx, createDTO)

				// Assert
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				Expect(err).To(Equal(existsErr))
			})
		})

		Context("when Create fails", func() {
			It("should return the error and rollback", func() {
				// Arrange
				createDTO := &dto.CreateCategoryDTO{
					Name: "TestCategory",
				}
				createErr := errs.NewCRUDError("DB_ERROR", "failed to create category")
				gomock.InOrder(
					mockTm.EXPECT().Begin(ctx).Return(mockTx, nil),
					mockRepo.EXPECT().ExistsByName(ctx, mockTx, testCategory.Name()).Return(false, nil),
					mockRepo.EXPECT().Create(ctx, mockTx, gomock.Any()).Do(
						func(ctx context.Context, tx *sql.Tx, category *categories.Category) {
							Expect(category).NotTo(BeNil())
							Expect(category.Name()).To(Equal(testCategory.Name()))
						}).Return(createErr),
					mockTm.EXPECT().Complete(ctx, mockTx, createErr).Return(nil),
				)

				// Act
				result, err := cs.Add(ctx, createDTO)

				// Assert
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				Expect(err).To(Equal(createErr))
			})
		})

		Context("when commit fails", func() {
			It("should return the commit error", func() {
				// Arrange
				createDTO := &dto.CreateCategoryDTO{
					Name: "TestCategory",
				}
				commitErr := fmt.Errorf("commit failed")
				gomock.InOrder(
					mockTm.EXPECT().Begin(ctx).Return(mockTx, nil),
					mockRepo.EXPECT().ExistsByName(ctx, mockTx, testCategory.Name()).Return(false, nil),
					mockRepo.EXPECT().Create(ctx, mockTx, gomock.Any()).Return(nil),
					mockTm.EXPECT().Complete(ctx, mockTx, nil).Return(commitErr),
				)

				// Act
				result, err := cs.Add(ctx, createDTO)

				// Assert
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				Expect(err).To(Equal(commitErr))
			})
		})
	})

	Describe("Update", func() {
		Context("when category exists", func() {
			It("should successfully update the category", func() {
				// Arrange
				updateDTO := &dto.UpdateCategoryDTO{
					Id:   testCategory.Id().Value(),
					Name: "UpdatedCategory",
				}
				gomock.InOrder(
					mockTm.EXPECT().Begin(ctx).Return(mockTx, nil),
					mockRepo.EXPECT().UpdateById(ctx, mockTx, gomock.Any()).Return(nil),
					mockTm.EXPECT().Complete(ctx, mockTx, nil).Return(nil),
				)

				// Act
				result, err := cs.Update(ctx, updateDTO)

				// Assert
				Expect(err).NotTo(HaveOccurred())
				Expect(result).NotTo(BeNil())
				Expect(result.Name).To(Equal("UpdatedCategory"))
			})
		})

		Context("when UpdateById fails", func() {
			It("should return the error and rollback", func() {
				// Arrange
				updateDTO := &dto.UpdateCategoryDTO{
					Id:   testCategory.Id().Value(),
					Name: "UpdatedCategory",
				}
				updateErr := errs.NewCRUDError("NOT_FOUND", "category not found")
				gomock.InOrder(
					mockTm.EXPECT().Begin(ctx).Return(mockTx, nil),
					mockRepo.EXPECT().UpdateById(ctx, mockTx, gomock.Any()).Return(updateErr),
					mockTm.EXPECT().Complete(ctx, mockTx, updateErr).Return(nil),
				)

				// Act
				result, err := cs.Update(ctx, updateDTO)

				// Assert
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				Expect(err).To(Equal(updateErr))
			})
		})

		Context("when Begin fails", func() {
			It("should return the error from Begin", func() {
				// Arrange
				updateDTO := &dto.UpdateCategoryDTO{
					Id:   testCategory.Id().Value(),
					Name: "UpdatedCategory",
				}
				beginErr := fmt.Errorf("failed to begin transaction")
				mockTm.EXPECT().Begin(ctx).Return(nil, beginErr)

				// Act
				result, err := cs.Update(ctx, updateDTO)

				// Assert
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				Expect(err).To(Equal(beginErr))
			})
		})
	})

	Describe("Delete", func() {
		Context("when category exists", func() {
			It("should successfully delete the category", func() {
				// Arrange
				deleteDTO := &dto.DeleteCategoryDTO{
					Id: testCategory.Id().Value(),
				}
				gomock.InOrder(
					mockTm.EXPECT().Begin(ctx).Return(mockTx, nil),
					mockRepo.EXPECT().FindById(ctx, mockTx, testCategory.Id()).Return(testCategory, nil),
					mockRepo.EXPECT().DeleteById(ctx, mockTx, testCategory.Id()).Return(nil),
					mockTm.EXPECT().Complete(ctx, mockTx, nil).Return(nil),
				)

				// Act
				result, err := cs.Delete(ctx, deleteDTO)

				// Assert
				Expect(err).NotTo(HaveOccurred())
				Expect(result).NotTo(BeNil())
				Expect(result.Name).To(Equal("TestCategory"))
			})
		})

		Context("when DeleteById fails", func() {
			It("should return the error and rollback", func() {
				// Arrange
				deleteDTO := &dto.DeleteCategoryDTO{
					Id: testCategory.Id().Value(),
				}
				deleteErr := errs.NewCRUDError("NOT_FOUND", "category not found")
				gomock.InOrder(
					mockTm.EXPECT().Begin(ctx).Return(mockTx, nil),
					mockRepo.EXPECT().FindById(ctx, mockTx, testCategory.Id()).Return(nil, deleteErr),
					mockTm.EXPECT().Complete(ctx, mockTx, deleteErr).Return(nil),
				)

				// Act
				result, err := cs.Delete(ctx, deleteDTO)

				// Assert
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				Expect(err).To(Equal(deleteErr))
			})
		})

		Context("when Begin fails", func() {
			It("should return the error from Begin", func() {
				// Arrange
				deleteDTO := &dto.DeleteCategoryDTO{
					Id: testCategory.Id().Value(),
				}
				beginErr := fmt.Errorf("failed to begin transaction")
				mockTm.EXPECT().Begin(ctx).Return(nil, beginErr)

				// Act
				result, err := cs.Delete(ctx, deleteDTO)

				// Assert
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				Expect(err).To(Equal(beginErr))
			})
		})
	})
})
