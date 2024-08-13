package router

import (
	"context"
	"message_app/internal/logger"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

var Module = fx.Module("router",
	fx.Provide(
		NewRouter,
	),
	fx.Invoke(startServer),
)

func startServer(lc fx.Lifecycle, e *echo.Echo, logger *logger.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {

			logger.L.Info("Starting HTTP server")
			go e.Start(":1323")
			logger.L.Info("Started HTTP server")
			return nil
		},
	})
}
