package config

import (
	"context"
)

// ===
// Interface for configuration settings
// ===
type IConfig interface {
	LoadConfig(ctx context.Context, pathConfig string) (*SystemConfig, error)
}

