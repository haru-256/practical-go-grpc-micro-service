package impl

import (
	"context"
	"log"

	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/application/service"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/domain/models/categories"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/errs"
)

type categoryServiceImpl struct {
	repo categories.CategoryRepository
	tm   service.TransactionManager
}

func NewCategoryService(repo categories.CategoryRepository, tm service.TransactionManager) service.CategoryService {
	return &categoryServiceImpl{
		repo: repo,
		tm:   tm,
	}
}

func (s *categoryServiceImpl) Add(ctx context.Context, category *categories.Category) error {
	tx, err := s.tm.Begin(ctx)
	if err != nil {
		return err
	}
	// NOTE: defer内でerrを評価するため、クロージャで囲む。defer時点のerrを参照させるため。
	defer func() {
		err = s.tm.Complete(tx, err)
		if err != nil {
			log.Fatalln("トランザクションの完了に失敗しました:", err)
		}
	}()

	exists, err := s.repo.ExistsByName(ctx, tx, category.Name())
	if err != nil {
		return err
	}
	if exists {
		return errs.NewApplicationError("CATEGORY_ALREADY_EXISTS", "Category already exists")
	}

	if err = s.repo.Create(ctx, tx, category); err != nil {
		return err
	}

	return nil
}

func (s *categoryServiceImpl) Update(ctx context.Context, category *categories.Category) error {
	tx, err := s.tm.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		err = s.tm.Complete(tx, err)
		if err != nil {
			log.Fatalln("トランザクションの完了に失敗しました:", err)
		}
	}()

	if err = s.repo.UpdateById(ctx, tx, category); err != nil {
		return err
	}

	return nil
}

func (s *categoryServiceImpl) Delete(ctx context.Context, category *categories.Category) error {
	tx, err := s.tm.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		err = s.tm.Complete(tx, err)
		if err != nil {
			log.Fatalln("トランザクションの完了に失敗しました:", err)
		}
	}()

	if err = s.repo.DeleteById(ctx, tx, category.Id()); err != nil {
		return err
	}

	return nil
}
