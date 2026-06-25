package initialize

import (
	"context"

	"github.com/youknow2509/temp-go-ddd/internal/constant"
	"github.com/youknow2509/temp-go-ddd/internal/global"
	infra_config "github.com/youknow2509/temp-go-ddd/internal/infrastructure/config"
)

// ===
// private function to initialize the configurations system
// ===

func initializeConfig() error {
	config_instance := infra_config.NewViperConfig()
	config_data, err := config_instance.LoadConfig(context.Background(), constant.ConfigFileName)
	if err != nil {
		return err
	}
	// fmt.Printf("%+v\n", config_data)

	// Store the loaded configuration in the global variable
	global.SystemConfig = config_data
	return nil
}
