package config

import (
	"context"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"github.com/youknow2509/temp-go-ddd/internal/constant"
	domain_config "github.com/youknow2509/temp-go-ddd/internal/domain/config"
)

type ViperConfig struct{}

func (v *ViperConfig) LoadConfig(ctx context.Context, pathConfig string) (*domain_config.SystemConfig, error) {
	// Load default config
	defaultViper := viper.New()
	defaultViper.SetConfigType(constant.ConfigFileType)
	pathFileConfigDefault := filepath.Join(constant.FolderConfig, constant.DefaultNameFileConfig)
	defaultViper.SetConfigFile(pathFileConfigDefault)
	if err := defaultViper.ReadInConfig(); err != nil {
		return nil, err
	}

	// Load override config và deep merge if pathConfig is provided
	if pathConfig != "" {
		overrideViper := viper.New()
		overrideViper.SetConfigType(constant.ConfigFileType)
		pathFileConfigOverride := filepath.Join(constant.FolderConfig, pathConfig)
		overrideViper.SetConfigFile(pathFileConfigOverride)
		if err := overrideViper.ReadInConfig(); err != nil {
			return nil, err
		}
		defaultSettings := defaultViper.AllSettings()
		overrideSettings := overrideViper.AllSettings()
		merged := helperDeepMerge(defaultSettings, overrideSettings)

		if err := defaultViper.MergeConfigMap(merged); err != nil {
			return nil, err
		}
	}

	// Unmarshal config data
	configData := &domain_config.SystemConfig{}
	if err := defaultViper.Unmarshal(configData); err != nil {
		return nil, err
	}

	// Override system mode from env variable if exists
	switch systemMode := os.Getenv(constant.SystemModeEnvKey); systemMode {
	case constant.SystemModeDevelopment:
		configData.System.Mode = constant.SystemModeDevelopment
	default:
		configData.System.Mode = constant.SystemModeProduction
	}

	return configData, nil
}

func NewViperConfig() domain_config.IConfig {
	return &ViperConfig{}
}

// ===
// Helper
// ===

// helperDeepMerge: merge src vào dst đệ quy
// - Key có trong src -> ghi đè dst
// - Key KHÔNG có trong src -> giữ nguyên dst
// - Nested map -> đệ quy merge tiếp
func helperDeepMerge(dst, src map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	// Copy tất cả key từ dst (default) trước
	for k, v := range dst {
		result[k] = v
	}

	// Override bằng src, nếu cả hai là map -> đệ quy
	for k, srcVal := range src {
		dstVal, exists := result[k]
		if !exists {
			// Key mới chỉ có trong override -> thêm vào
			result[k] = srcVal
			continue
		}

		// Cả hai đều là map -> deep merge đệ quy
		srcMap, srcIsMap := srcVal.(map[string]interface{})
		dstMap, dstIsMap := dstVal.(map[string]interface{})
		if srcIsMap && dstIsMap {
			result[k] = helperDeepMerge(dstMap, srcMap)
			continue
		}

		// Không phải map -> override trực tiếp
		result[k] = srcVal
	}

	return result
}
