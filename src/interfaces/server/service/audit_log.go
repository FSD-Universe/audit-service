// Copyright (c) 2025 Half_nothing
// SPDX-License-Identifier: MIT

// Package service
package service

import (
	DTO "audit-service/src/interfaces/server/dto"

	"half-nothing.cn/service-core/interfaces/http/dto"
)

type AuditLogInterface interface {
	GetAuditLogPage(form *DTO.GetAuditLogs) *dto.ApiResponse[*DTO.GetAuditLogsResponse]
}
