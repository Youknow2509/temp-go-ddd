package config

// Grpc server structure setting
type GrpcServerSetting struct {
	Network   GrpcNetworkSetting   `mapstructure:"network"`
	Timeouts  GrpcTimeoutsSetting  `mapstructure:"timouts"` // Notice spelling in yaml: "timouts"
	Limits    GrpcLimitsSetting    `mapstructure:"limits"`
	Http2     GrpcHttp2Setting     `mapstructure:"http2"`
	Keepalive GrpcKeepaliveSetting `mapstructure:"keepalive"`
	Tcp       GrpcTcpSetting       `mapstructure:"tcp"`
	Protocol  GrpcProtocolSetting  `mapstructure:"protocol"`
	Security  GrpcSecuritySetting  `mapstructure:"security"`
}

type GrpcNetworkSetting struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type GrpcTimeoutsSetting struct {
	IdleTimeoutMs int `mapstructure:"idle_timeout_ms"`
}

type GrpcLimitsSetting struct {
	MaxDecodingMessageSize int `mapstructure:"max_decoding_message_size"`
	MaxEncodingMessageSize int `mapstructure:"max_encoding_message_size"`
}

type GrpcHttp2Setting struct {
	InitialStreamWindowSize     int `mapstructure:"initial_stream_window_size"`
	InitialConnectionWindowSize int `mapstructure:"initial_connection_window_size"`
	MaxConcurrentStreams        int `mapstructure:"max_concurrent_streams"`
}

type GrpcKeepaliveSetting struct {
	Http2KeepaliveIntervalMs int  `mapstructure:"http2_keepalive_interval_ms"`
	Http2KeepaliveTimeoutMs  int  `mapstructure:"http2_keepalive_timeout_ms"`
	KeepaliveWhileIdle       bool `mapstructure:"keepalive_while_idle"`
}

type GrpcTcpSetting struct {
	TcpNodelay         bool `mapstructure:"tcp_nodelay"`
	TcpKeepalive       bool `mapstructure:"tcp_keepalive"`
	TcpKeepaliveTimeMs int  `mapstructure:"tcp_keepalive_time_ms"`
}

type GrpcProtocolSetting struct {
	ReflectionEnabled  bool `mapstructure:"reflection_enabled"`
	HealthCheckEnabled bool `mapstructure:"health_check_enabled"`
}

type GrpcSecuritySetting struct {
	Tls GrpcServerTlsSetting `mapstructure:"tls"`
}

type GrpcServerTlsSetting struct {
	IsEnabled         bool   `mapstructure:"is_enabled"`
	CertFile          string `mapstructure:"cert_file"`
	KeyFile           string `mapstructure:"key_file"`
	RequireClientCert bool   `mapstructure:"require_client_cert"`
	MinVersion        string `mapstructure:"min_version"`
}
