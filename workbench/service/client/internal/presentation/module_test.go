//go:build integration || !ci

package presentation_test

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/haru-256/practical-go-grpc-micro-service/service/client/internal/presentation"
	"github.com/haru-256/practical-go-grpc-micro-service/service/client/internal/presentation/dto"
	"github.com/haru-256/practical-go-grpc-micro-service/service/client/internal/presentation/server"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
)

// setupTestApp は統合テスト用のFxアプリケーションをセットアップします
func setupTestApp(t *testing.T, populate ...interface{}) *fx.App {
	t.Helper()
	configOption := fx.Supply(
		fx.Annotate("../../", fx.ResultTags(`name:"configPath"`)),
		fx.Annotate("config", fx.ResultTags(`name:"configName"`)),
	)

	app := fx.New(
		configOption,
		presentation.Module,
		// テスト時のみデコレートして差し替える
		fx.Decorate(func() (*slog.Logger, error) {
			return slog.New(slog.NewTextHandler(io.Discard, nil)), nil
		}),
		fx.Populate(populate...),
		fx.NopLogger,
	)
	require.NoError(t, app.Err(), "fx app should initialize without errors")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	t.Cleanup(cancel)
	require.NoError(t, app.Start(ctx), "fx app should start without errors")
	t.Cleanup(func() {
		stopCtx, stopCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer stopCancel()
		assert.NoError(t, app.Stop(stopCtx), "fx app should stop without errors")
	})

	return app
}

func TestModule(t *testing.T) {
	t.Run("依存関係の注入と初期化", func(t *testing.T) {
		// Arrange
		var (
			cfg *server.CQRSServiceConfig
			h   *server.CQRSServiceHandler
			srv *server.CQRSServiceServer
		)

		// Act
		setupTestApp(t, &cfg, &h, &srv)

		// Assert
		assert.NotNil(t, cfg, "config should not be nil")
		assert.NotNil(t, h, "handler should not be nil")
		assert.NotNil(t, srv, "server should not be nil")

		// 設定値の確認
		assert.NotEmpty(t, cfg.Port, "server port should not be empty")
	})

	t.Run("CategoryListハンドラーの動作確認", func(t *testing.T) {
		// Arrange
		var h *server.CQRSServiceHandler

		setupTestApp(t, &h)
		require.NotNil(t, h, "handler should not be nil")

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/categories", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Act
		err := h.CategoryList(c)

		// Assert
		require.NoError(t, err, "CategoryList should not return an error")
		assert.Equal(t, http.StatusOK, rec.Code, "response status should be 200 OK")

		var resp dto.CategoryListResponse
		err = json.Unmarshal(rec.Body.Bytes(), &resp)
		require.NoError(t, err, "should unmarshal response body without error")
		assert.NotNil(t, resp.Categories, "categories should not be nil")

		// カテゴリが存在する場合のみ詳細チェック
		if len(resp.Categories) > 0 {
			for _, cat := range resp.Categories {
				assert.NotEmpty(t, cat.Id, "category ID should not be empty")
				assert.NotEmpty(t, cat.Name, "category name should not be empty")
			}
		}
	})
}
