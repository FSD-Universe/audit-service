// Package config
package config

import (
	"audit-service/src/interfaces/global"

	"half-nothing.cn/service-core/interfaces/config"
)

type GlobalConfig struct {
	*config.GlobalConfig `yaml:",inline"`
}

func (g *GlobalConfig) InitDefaults() {
	g.GlobalConfig = &config.GlobalConfig{}
	g.GlobalConfig.InitDefaults()
	g.Name = "audit-service"
	g.Version = global.ConfigVersion
	g.ConfigVersion = global.ConfigVersion
}
