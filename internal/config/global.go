package config

import (
	"sync"
)

// 全局配置实例
var (
	globalConfig *ButterflyConfig
	configMutex  sync.RWMutex
)

// SetGlobalConfig 设置全局配置
func SetGlobalConfig(cfg *ButterflyConfig) {
	configMutex.Lock()
	globalConfig = cfg
	configMutex.Unlock()
}

// GetGlobalConfig 获取全局配置
func GetGlobalConfig() *ButterflyConfig {
	configMutex.RLock()
	defer configMutex.RUnlock()
	if globalConfig == nil {
		return &ButterflyConfig{
			JWT: JWTConfig{
				Secret:        "default-secret-key",
				AccessExpiry:  3600,
				RefreshExpiry: 604800,
			},
		}
	}
	return globalConfig
}

// Get 获取配置的快捷方式 (兼容旧代码)
func Get() *ButterflyConfig {
	return GetGlobalConfig()
}