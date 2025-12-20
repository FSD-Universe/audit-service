// Copyright (c) 2025 Half_nothing
// SPDX-License-Identifier: MIT

// Package repository
package repository

import (
	"errors"

	"half-nothing.cn/service-core/interfaces/database/entity"
	"half-nothing.cn/service-core/interfaces/database/repository"
)

var (
	ErrAuditLogNotFound = errors.New("audit log not found")
)

type AuditLogInterface interface {
	repository.Base[*entity.AuditLog]
	GetPage(pageNum int, pageSize int, logType []entity.AuditEvent, search string) ([]*entity.AuditLog, int64, error)
}
