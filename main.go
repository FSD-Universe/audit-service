// Copyright (c) 2025 Half_nothing
// SPDX-License-Identifier: MIT

// Package main
package main

import (
	c "audit-service/src/interfaces/config"
	"fmt"

	"half-nothing.cn/service-core/config"
	"half-nothing.cn/service-core/interfaces/global"
)

func main() {
	global.CheckFlags()

	configManager := config.NewManager[*c.Config]()
	if err := configManager.Init(); err != nil {
		fmt.Printf("fail to initialize configuration file: %v", err)
		return
	}
}
