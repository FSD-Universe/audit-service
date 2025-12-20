// Copyright (c) 2025 Half_nothing
// SPDX-License-Identifier: MIT

// Package repository
package repository

import (
	r "audit-service/src/interfaces/repository"
	"errors"
	"time"

	"gorm.io/gorm"
	"half-nothing.cn/service-core/database"
	"half-nothing.cn/service-core/interfaces/database/entity"
	"half-nothing.cn/service-core/interfaces/database/repository"
	"half-nothing.cn/service-core/interfaces/logger"
	"half-nothing.cn/service-core/utils"
)

type AuditLogRepository struct {
	*database.BaseRepository[*entity.AuditLog]
	pageReq database.PageableInterface[*entity.AuditLog]
}

func NewAuditLogRepository(
	lg logger.Interface,
	db *gorm.DB,
	queryTimeout time.Duration,
) *AuditLogRepository {
	return &AuditLogRepository{
		BaseRepository: database.NewBaseRepository[*entity.AuditLog](lg, "audit-log-repository", db, queryTimeout),
		pageReq:        database.NewPageRequest[*entity.AuditLog](db),
	}
}

func (repo *AuditLogRepository) GetById(id uint) (*entity.AuditLog, error) {
	if id <= 0 {
		return nil, repository.ErrArgument
	}
	auditLog := &entity.AuditLog{ID: id}
	err := repo.Query(func(tx *gorm.DB) error {
		return tx.First(auditLog).Error
	})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = r.ErrAuditLogNotFound
	}
	return auditLog, err
}

func (repo *AuditLogRepository) GetPage(
	pageNum int,
	pageSize int,
	logType []entity.AuditEvent,
	search string,
) (auditLogs []*entity.AuditLog, total int64, err error) {
	auditLogs = make([]*entity.AuditLog, 0, pageSize)
	queryFunc := func(tx *gorm.DB) *gorm.DB {
		if search != "" {
			tx = tx.Where("object like ? or subject like ?", "%"+search+"%", "%"+search+"%")
		}
		if logType != nil && len(logType) > 0 {
			types := make([]string, len(logType))
			utils.ForEach(logType, func(i int, logType entity.AuditEvent) { types[i] = logType.Value })
			tx = tx.Where("event in ?", types)
		}
		return tx.Order("created_at desc")
	}
	total, err = repo.QueryWithPagination(repo.pageReq, database.NewPage(pageNum, pageSize, &auditLogs, &entity.AuditLog{}, queryFunc))
	return
}
