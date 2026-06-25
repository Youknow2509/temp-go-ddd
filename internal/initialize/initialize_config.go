package initialize

import (
	"context"

	"github.com/youknow2509/temp-go-ddd/internal/constant"
	domain_config "github.com/youknow2509/temp-go-ddd/internal/domain/config"
	infra_config "github.com/youknow2509/temp-go-ddd/internal/infrastructure/config"
)

// ===
// private function to initialize the configurations system
// ===

func initializeConfig(ctx context.Context) (*domain_config.SystemConfig, error) {
	config_instance := infra_config.NewViperConfig()
	config_data, err := config_instance.LoadConfig(ctx, constant.ConfigFileName)
	if err != nil {
		return nil, err
	}

	return config_data, nil
}
