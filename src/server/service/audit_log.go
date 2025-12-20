// Copyright (c) 2025 Half_nothing
// SPDX-License-Identifier: MIT

// Package service
package service

import (
	"audit-service/src/interfaces/repository"
	DTO "audit-service/src/interfaces/server/dto"

	"half-nothing.cn/service-core/interfaces/database/entity"
	"half-nothing.cn/service-core/interfaces/http/dto"
	"half-nothing.cn/service-core/interfaces/logger"
	"half-nothing.cn/service-core/permission"
	"half-nothing.cn/service-core/utils"
)

type AuditLogService struct {
	logger logger.Interface
	repo   repository.AuditLogInterface
}

func NewAuditLogService(
	lg logger.Interface,
	repo repository.AuditLogInterface,
) *AuditLogService {
	return &AuditLogService{
		logger: logger.NewLoggerAdapter(lg, "audit-log-service"),
		repo:   repo,
	}
}

func (a *AuditLogService) GetAuditLogPage(form *DTO.GetAuditLogs) *dto.ApiResponse[*DTO.GetAuditLogsResponse] {
	perm := permission.Permission(form.Permission)
	if !perm.HasPermission(permission.AuditLogShow) {
		return dto.NewApiResponse[*DTO.GetAuditLogsResponse](dto.ErrNoPermission, nil)
	}
	events := make([]entity.AuditEvent, len(form.Event))
	for i, event := range form.Event {
		if !entity.AuditEventManager.IsValidEnum(event) {
			return dto.NewApiResponse[*DTO.GetAuditLogsResponse](dto.ErrErrorParam, nil)
		}
		events[i] = entity.AuditEventManager.GetEnum(event)
	}
	auditLogs, total, err := a.repo.GetPage(form.PageNum, form.PageSize, events, form.Search)
	if err != nil {
		return dto.NewApiResponse[*DTO.GetAuditLogsResponse](dto.ErrServerError, nil)
	}
	logs := make([]*DTO.AuditLog, len(auditLogs))
	utils.ForEach(auditLogs, func(i int, auditLog *entity.AuditLog) {
		log := &DTO.AuditLog{}
		logs[i] = log.FromAuditLogEntity(auditLog)
	})
	return dto.NewApiResponse[*DTO.GetAuditLogsResponse](dto.SuccessHandleRequest, &DTO.GetAuditLogsResponse{
		Data:     logs,
		Total:    total,
		PageNum:  form.PageNum,
		PageSize: form.PageSize,
	})
}
