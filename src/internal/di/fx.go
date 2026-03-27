package di

import (
	"context"
	"net/http"

	"tictactoe/internal/datasource"
	"tictactoe/internal/domain"
	"tictactoe/internal/web"

	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		datasource.NewStorage,
		datasource.NewRepository,
		domain.NewService,
		web.NewHandler,
		http.NewServeMux,
	),
	fx.Invoke(registerHooks),
)

func registerHooks(lc fx.Lifecycle, handler *web.Handler, mux *http.ServeMux) {
	handler.RegisterRoutes(mux)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go server.ListenAndServe()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return server.Shutdown(ctx)
		},
	})
}
