// Copyright (c) 2025 Half_nothing
// SPDX-License-Identifier: MIT

// Package dto
package dto

import (
	"time"

	"half-nothing.cn/service-core/interfaces/database/entity"
	"half-nothing.cn/service-core/interfaces/http/jwt"
)

type AuditLog struct {
	Subject   string               `json:"subject"`
	Object    string               `json:"object"`
	Event     string               `json:"event"`
	Ip        string               `json:"ip"`
	UserAgent string               `json:"user_agent"`
	Detail    *entity.ChangeDetail `json:"detail"`
	LogTime   time.Time            `json:"log_time"`
}

func (log *AuditLog) FromAuditLogEntity(entity *entity.AuditLog) *AuditLog {
	log.Subject = entity.Subject
	log.Object = entity.Object
	log.Event = entity.Event
	log.Ip = entity.Ip
	log.UserAgent = entity.UserAgent
	log.Detail = entity.ChangeDetails
	log.LogTime = entity.CreatedAt
	return log
}

type GetAuditLogs struct {
	jwt.Content
	PageNum  int      `query:"page_num" valid:"required,min=0"`
	PageSize int      `query:"page_size" valid:"required,min=0"`
	Search   string   `query:"search"`
	Event    []string `query:"event"`
}

type GetAuditLogsResponse struct {
	Data     []*AuditLog `json:"page_data"`
	PageNum  int         `json:"page_num"`
	PageSize int         `json:"page_size"`
	Total    int64       `json:"total"`
}
