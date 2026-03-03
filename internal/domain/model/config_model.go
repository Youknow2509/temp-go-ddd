package model

// ==============================================================
// Configuration file structure for system
// ==============================================================
type SystemConfig struct {
	System        SystemSetting        `mapstructure:"system"`
	HttpServer    HttpServerSetting    `mapstructure:"http_server"`
	GrpcServer    GrpcServerSetting    `mapstructure:"grpc_server"`
	Logger        LoggerSetting        `mapstructure:"logger"`
	Observability ObservabilitySetting `mapstructure:"observability"`
}

// ==============================================================
// Setting configuration structure used in the configuration file
// ==============================================================

// Observability structure setting
type ObservabilitySetting struct {
	Enabled        bool   `mapstructure:"enabled"`
	MetricsPath    string `mapstructure:"metrics_path"`
	MetricsPort    int    `mapstructure:"metrics_port"`
	TracingEnabled bool   `mapstructure:"tracing_enabled"`
	OtlpEndpoint   string `mapstructure:"otlp_endpoint"`
}

// Logger structure setting
type LoggerSetting struct {
	FolderStore    string `mapstructure:"folder_store"`
	FileMaxSize    int    `mapstructure:"file_max_size"`
	FileMaxBackups int    `mapstructure:"file_max_backups"`
	FileMaxAge     int    `mapstructure:"file_max_age"`
	Compress       bool   `mapstructure:"compress"`
}

// Grpc server structure setting
type GrpcServerSetting struct {
	Network                     string               `mapstructure:"network"`
	Port                        int                  `mapstructure:"port"`
	KeepaliveTimeMs             int                  `mapstructure:"keepalive_time_ms"`
	KeepaliveTimeoutMs          int                  `mapstructure:"keepalive_timeout_ms"`
	Http2MinTimeBetweenPingsMs  int                  `mapstructure:"http2_min_time_between_pings_ms"`
	KeepalivePermitWithoutCalls bool                 `mapstructure:"keepalive_permit_without_calls"`
	Tls                         GrpcServerTlsSetting `mapstructure:"tls"`
}

// Grpc server tsl sub structure setting
type GrpcServerTlsSetting struct {
	Enabled  bool   `mapstructure:"enabled"`
	CertFile string `mapstructure:"cert_file"`
	KeyFile  string `mapstructure:"key_file"`
}

// Http server structure setting
type HttpServerSetting struct {
	Port int `mapstructure:"port"`
}

// System setting structure setting
type SystemSetting struct {
	Name     string `mapstructure:"name"`
	Version  string `mapstructure:"version"`
	Region   string `mapstructure:"region"`
	ShardID  string `mapstructure:"shard_id"`
	Timezone string `mapstructure:"timezone"`
	Mode     string `mapstructure:"mode"`
}
