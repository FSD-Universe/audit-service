// Copyright (c) 2025 Half_nothing
// SPDX-License-Identifier: MIT

// Package server
package server

import (
	"audit-service/src/interfaces/content"
	"audit-service/src/server/controller"
	"audit-service/src/server/service"
	"io"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"half-nothing.cn/service-core/http"
	"half-nothing.cn/service-core/interfaces/logger"
)

func StartHttpServer(content *content.ApplicationContent) {
	c := content.ConfigManager().GetConfig()
	lg := logger.NewLoggerAdapter(content.Logger(), "http-server")

	lg.Info("Http server initializing...")
	e := echo.New()
	e.Logger.SetOutput(io.Discard)
	e.Logger.SetLevel(log.OFF)

	http.SetEchoConfig(lg, e, c.ServerConfig.HttpServerConfig, nil)

	jwtMiddleware, requireNoFresh, _ := http.GetJWTMiddleware(content.ClaimFactory())

	if c.TelemetryConfig.HttpServerTrace {
		http.SetTelemetry(e, c.TelemetryConfig)
	}

	auditController := controller.NewAuditLogController(lg, service.NewAuditLogService(lg, content.AuditLog()))

	http.SetHealthPoint(e)

	apiGroup := e.Group("/api/v1")
	auditGroup := apiGroup.Group("/audits")
	auditGroup.GET("", auditController.GetAuditLogPage, jwtMiddleware, requireNoFresh)

	http.SetUnmatchedRoute(e)
	http.SetCleaner(content.Cleaner(), e)

	http.Serve(lg, e, c.ServerConfig.HttpServerConfig)
}
