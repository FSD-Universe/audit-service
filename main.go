// Copyright (c) 2025 Half_nothing
// SPDX-License-Identifier: MIT

// Package main
package main

import (
	grpcImpl "audit-service/src/grpc"
	c "audit-service/src/interfaces/config"
	"audit-service/src/interfaces/content"
	g "audit-service/src/interfaces/global"
	pb "audit-service/src/interfaces/grpc"
	"audit-service/src/repository"
	"audit-service/src/server"
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"half-nothing.cn/service-core/cleaner"
	"half-nothing.cn/service-core/config"
	"half-nothing.cn/service-core/database"
	"half-nothing.cn/service-core/discovery"
	grpcUtils "half-nothing.cn/service-core/grpc"
	"half-nothing.cn/service-core/interfaces/global"
	"half-nothing.cn/service-core/jwt"
	"half-nothing.cn/service-core/logger"
	"half-nothing.cn/service-core/telemetry"
)

func main() {
	global.CheckFlags()

	configManager := config.NewManager[*c.Config]()
	if err := configManager.Init(); err != nil {
		fmt.Printf("fail to initialize configuration file: %v", err)
		return
	}

	applicationConfig := configManager.GetConfig()
	lg := logger.NewLogger()
	lg.Init(
		global.LogName,
		applicationConfig.GlobalConfig.LogConfig,
	)

	lg.Info(" _____       _ _ _   _____             _")
	lg.Info("|  _  |_ _ _| |_| |_|   __|___ ___ _ _|_|___ ___")
	lg.Info("|     | | | . | |  _|__   | -_|  _| | | |  _| -_|")
	lg.Info("|__|__|___|___|_|_| |_____|___|_|  \\_/|_|___|___|")
	lg.Info(fmt.Sprintf("%49s", fmt.Sprintf("Copyright © %d-%d Half_nothing", global.BeginYear, time.Now().Year())))
	lg.Info(fmt.Sprintf("%49s", fmt.Sprintf("AuditService v%s", g.AppVersion)))

	cl := cleaner.NewCleaner(lg)
	cl.Init()

	defer cl.Clean()

	if applicationConfig.TelemetryConfig.Enable {
		sdk := telemetry.NewSDK(lg, applicationConfig.TelemetryConfig)
		shutdown, err := sdk.SetupOTelSDK(context.Background())
		if err != nil {
			lg.Fatalf("fail to initialize telemetry: %v", err)
			return
		}
		cl.Add("Telemetry", shutdown)
	}

	closeFunc, db, err := database.InitDatabase(lg, applicationConfig.DatabaseConfig)
	if err != nil {
		lg.Fatalf("fail to initialize database: %v", err)
		return
	}
	cl.Add("Database", closeFunc)

	if applicationConfig.TelemetryConfig.DatabaseTrace {
		if err := database.ApplyDBTracing(db, "mysql"); err != nil {
			lg.Fatalf("fail to apply database tracing: %v", err)
			return
		}
	}

	auditLogRepo := repository.NewAuditLogRepository(lg, db, applicationConfig.DatabaseConfig.QueryTimeoutDuration)

	contentBuilder := content.NewApplicationContentBuilder().
		SetConfigManager(configManager).
		SetCleaner(cl).
		SetLogger(lg).
		SetJwtClaimFactory(jwt.NewClaimFactory(applicationConfig.JwtConfig)).
		SetAuditLog(auditLogRepo)

	started := make(chan bool)
	initFunc := func(s *grpc.Server) {
		grpcServer := grpcImpl.NewAuditLogServer(lg, auditLogRepo)
		pb.RegisterAuditLogServer(s, grpcServer)
	}
	if applicationConfig.TelemetryConfig.Enable && applicationConfig.TelemetryConfig.GrpcServerTrace {
		go grpcUtils.StartGrpcServerWithTrace(lg, cl, applicationConfig.ServerConfig.GrpcServerConfig, started, initFunc)
	} else {
		go grpcUtils.StartGrpcServer(lg, cl, applicationConfig.ServerConfig.GrpcServerConfig, started, initFunc)
	}

	if ok := <-started; !ok {
		lg.Fatal("fail to start grpc server")
		return
	}

	consulClient := discovery.NewConsulClient(lg, applicationConfig.GlobalConfig.Discovery, g.AppVersion)

	if err := consulClient.RegisterServer(); err != nil {
		lg.Fatalf("fail to register server: %v", err)
		return
	}

	cl.Add("Discovery", consulClient.UnregisterServer)

	go func() {
		for {
			<-consulClient.EventChan
		}
	}()

	go server.StartHttpServer(contentBuilder.Build())

	cl.Wait()
}
