// Copyright (c) 2025 Half_nothing
// SPDX-License-Identifier: MIT

// Package content
package content

import (
	c "audit-service/src/interfaces/config"
	"audit-service/src/interfaces/repository"

	"half-nothing.cn/service-core/interfaces/cleaner"
	"half-nothing.cn/service-core/interfaces/config"
	"half-nothing.cn/service-core/interfaces/http/jwt"
	"half-nothing.cn/service-core/interfaces/logger"
)

// ApplicationContent 应用程序上下文结构体，包含所有核心组件的接口
type ApplicationContent struct {
	configManager config.ManagerInterface[*c.Config] // 配置管理器
	cleaner       cleaner.Interface                  // 清理器
	logger        logger.Interface                   // 日志
	claimFactory  jwt.ClaimFactoryInterface          // JWT 令牌工厂
	auditLogRepo  repository.AuditLogInterface       // 审计日志数据库操作
}

func (app *ApplicationContent) ConfigManager() config.ManagerInterface[*c.Config] {
	return app.configManager
}

func (app *ApplicationContent) Cleaner() cleaner.Interface { return app.cleaner }

func (app *ApplicationContent) Logger() logger.Interface { return app.logger }

func (app *ApplicationContent) ClaimFactory() jwt.ClaimFactoryInterface {
	return app.claimFactory
}

func (app *ApplicationContent) AuditLog() repository.AuditLogInterface {
	return app.auditLogRepo
}
