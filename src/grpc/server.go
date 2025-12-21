// Copyright (c) 2025 Half_nothing
// SPDX-License-Identifier: MIT

// Package grpc
package grpc

import (
	pb "audit-service/src/interfaces/grpc"
	"audit-service/src/interfaces/repository"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"half-nothing.cn/service-core/interfaces/database/entity"
	"half-nothing.cn/service-core/interfaces/logger"
)

type AuditLogServer struct {
	pb.UnimplementedAuditLogServer
	logger logger.Interface
	repo   repository.AuditLogInterface
}

func NewAuditLogServer(
	lg logger.Interface,
	auditLogRepo repository.AuditLogInterface,
) *AuditLogServer {
	return &AuditLogServer{
		logger: logger.NewLoggerAdapter(lg, "grpc-server"),
		repo:   auditLogRepo,
	}
}

func (a *AuditLogServer) Log(_ context.Context, request *pb.AuditLogRequest) (*pb.AuditLogResponse, error) {
	if request.Event == "" || !entity.AuditEventManager.IsValidEnum(request.Event) {
		return nil, status.Errorf(codes.InvalidArgument, "invalid event: %s", request.Event)
	}
	if request.Subject == "" || request.Ip == "" {
		return nil, status.Errorf(codes.InvalidArgument, "subject or ip is empty")
	}
	var detail *entity.ChangeDetail
	if request.NewValue == "" && request.OldValue == "" {
		detail = nil
	} else {
		detail = &entity.ChangeDetail{
			NewValue: request.NewValue,
			OldValue: request.OldValue,
		}
	}
	log := &entity.AuditLog{
		Event:         request.Event,
		Subject:       request.Subject,
		Object:        request.Object,
		Ip:            request.Ip,
		UserAgent:     request.UserAgent,
		ChangeDetails: detail,
	}
	a.logger.Infof("grpc call, create audit log %#v", log)
	if err := a.repo.Save(log); err != nil {
		a.logger.Errorf("failed to create audit log, %+v", err)
		return nil, status.Errorf(codes.Internal, "failed to create audit log")
	}
	return &pb.AuditLogResponse{Success: true}, nil
}
