package main

import (
	"message_app/internal/handler"
	"message_app/internal/infra"
	"message_app/internal/logger"
	"message_app/internal/router"
	"message_app/internal/service"

	"go.uber.org/fx"
)

var Modules = fx.Module("server",
	logger.Module,
	infra.Module,
	service.Module,
	handler.Module,
	router.Module,
)
