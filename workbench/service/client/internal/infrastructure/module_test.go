//go:build integration || !ci

package infrastructure_test

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"testing"
	"time"

	"github.com/haru-256/practical-go-grpc-micro-service/service/client/internal/domain/repository"
	"github.com/haru-256/practical-go-grpc-micro-service/service/client/internal/infrastructure"
	"github.com/haru-256/practical-go-grpc-micro-service/service/client/internal/infrastructure/cqrs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
)

func TestModule(t *testing.T) {
	configOption := fx.Supply(
		fx.Annotate("../../", fx.ResultTags(`name:"configPath"`)),
		fx.Annotate("config", fx.ResultTags(`name:"configName"`)),
	)
	t.Run("依存関係の注入と初期化", func(t *testing.T) {
		var (
			client               *http.Client
			logger               *slog.Logger
			commandServiceClient *cqrs.CommandServiceClient
			queryServiceClient   *cqrs.QueryServiceClient
			repo                 repository.CQRSRepository
		)

		app := fx.New(
			configOption,
			infrastructure.Module,
			// テスト時のみデコレートして差し替える
			fx.Decorate(func() (*slog.Logger, error) {
				return slog.New(slog.NewTextHandler(io.Discard, nil)), nil
			}),
			fx.Populate(&client, &logger, &commandServiceClient, &queryServiceClient, &repo),
			fx.NopLogger, // テスト時はログを抑制
		)

		t.Cleanup(func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			assert.NoError(t, app.Stop(ctx), "fx app should stop without errors")
		})

		require.NoError(t, app.Err(), "fx app should initialize without errors")

		// 各依存関係がnilでないことを確認
		assert.NotNil(t, client, "http client should not be nil")
		assert.NotNil(t, logger, "logger should not be nil")
		assert.NotNil(t, commandServiceClient, "command service client should not be nil")
		assert.NotNil(t, queryServiceClient, "query service client should not be nil")
		assert.NotNil(t, repo, "repository should not be nil")
	})

	t.Run("repositoryの動作確認", func(t *testing.T) {
		var (
			repo repository.CQRSRepository
		)

		// Arrange
		app := fx.New(
			configOption,
			infrastructure.Module,
			fx.Populate(&repo),
			fx.NopLogger, // テスト時はログを抑制
		)
		require.NoError(t, app.Err(), "fx app should initialize without errors")

		// Act
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		t.Cleanup(cancel)
		require.NoError(t, app.Start(ctx), "fx app should start without errors")
		t.Cleanup(func() {
			stopCtx, stopCancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer stopCancel()
			assert.NoError(t, app.Stop(stopCtx), "fx app should stop without errors")
		})

		// Assert
		assert.NotNil(t, repo, "repository should not be nil")
		categories, err := repo.CategoryList(ctx)
		require.NoError(t, err, "should fetch category list without error")
		assert.Greater(t, len(categories), 0, "category list should not be empty")
	})
}
