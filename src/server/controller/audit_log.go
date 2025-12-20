// Copyright (c) 2025 Half_nothing
// SPDX-License-Identifier: MIT

// Package controller
package controller

import (
	DTO "audit-service/src/interfaces/server/dto"
	"audit-service/src/interfaces/server/service"

	"github.com/labstack/echo/v4"
	"half-nothing.cn/service-core/interfaces/http/dto"
	"half-nothing.cn/service-core/interfaces/http/jwt"
	"half-nothing.cn/service-core/interfaces/logger"
)

type AuditLogController struct {
	logger  logger.Interface
	service service.AuditLogInterface
}

func NewAuditLogController(
	lg logger.Interface,
	service service.AuditLogInterface,
) *AuditLogController {
	return &AuditLogController{
		logger:  logger.NewLoggerAdapter(lg, "audit-log-controller"),
		service: service,
	}
}

func (controller *AuditLogController) GetAuditLogPage(ctx echo.Context) error {
	data := &DTO.GetAuditLogs{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.Errorf("GetAuditLogPage handle fail, parse argument fail, %v", err)
		return dto.ErrorResponse(ctx, dto.ErrErrorParam)
	}
	res, err := dto.ValidStruct(data)
	if err != nil {
		controller.logger.Errorf("GetAuditLogPage handle fail, validate err, %v", err)
		return dto.ErrorResponse(ctx, dto.ErrServerError)
	}
	if res != nil {
		controller.logger.Errorf("GetAuditLogPage handle fail, validate argument fail, %v", res)
		return dto.ErrorResponse(ctx, res)
	}
	if err := jwt.SetJwtContent(data, ctx); err != nil {
		controller.logger.Errorf("GetAuditLogPage handle fail, set jwt content fail, %v", err)
		return dto.ErrorResponse(ctx, dto.ErrUnknownJwtError)
	}
	controller.logger.Debugf("GetAuditLogPage with argument %#v", data)
	return controller.service.GetAuditLogPage(data).Response(ctx)
}
