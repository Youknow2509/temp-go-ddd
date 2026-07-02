package config

// Http server structure setting
type HttpServerSetting struct {
	Network   HttpNetworkSetting   `mapstructure:"network"`
	Timeouts  HttpTimeoutsSetting  `mapstructure:"timeouts"`
	Limits    HttpLimitsSetting    `mapstructure:"limits"`
	Tcp       HttpTcpSetting       `mapstructure:"tcp"`
	Security  HttpSecuritySetting  `mapstructure:"security"`
	RateLimit HttpRateLimitSetting `mapstructure:"rate_limit"`
}

type HttpNetworkSetting struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type HttpTimeoutsSetting struct {
	ReadTimeoutMs     int `mapstructure:"read_timeout_ms"`
	WriteTimeoutMs    int `mapstructure:"write_timeout_ms"`
	IdleTimeoutMs     int `mapstructure:"idle_timeout_ms"`
	ShutdownTimeoutMs int `mapstructure:"shutdown_timeout_ms"`
}

type HttpLimitsSetting struct {
	MaxBodySize   int64 `mapstructure:"max_body_size"`
	MaxHeaderSize int   `mapstructure:"max_header_size"`
}

type HttpTcpSetting struct {
	TcpNodelay         bool `mapstructure:"tcp_nodelay"`
	TcpKeepalive       bool `mapstructure:"tcp_keepalive"`
	TcpKeepaliveTimeMs int  `mapstructure:"tcp_keepalive_time_ms"`
}

type HttpSecuritySetting struct {
	Cors HttpCorsSetting `mapstructure:"cors"`
	Tls  HttpTlsSetting  `mapstructure:"tls"`
}

type HttpCorsSetting struct {
	Enabled             bool                   `mapstructure:"enabled"`
	AllowCredentials    bool                   `mapstructure:"allow_credentials"`
	MaxAgeSecs          int                    `mapstructure:"max_age_secs"`
	AllowPrivateNetwork bool                   `mapstructure:"allow_private_network"`
	Origin              HttpCorsOriginSetting  `mapstructure:"origin"`
	Methods             HttpCorsMethodsSetting `mapstructure:"methods"`
	Headers             HttpCorsHeadersSetting `mapstructure:"headers"`
}

type HttpCorsOriginSetting struct {
	AllowedOrigins []string `mapstructure:"allowed_origins"`
	AllowAnyOrigin bool     `mapstructure:"allow_any_origin"`
}

type HttpCorsMethodsSetting struct {
	AllowedMethods []string `mapstructure:"allowed_methods"`
	AllowAnyMethod bool     `mapstructure:"allow_any_method"`
}

type HttpCorsHeadersSetting struct {
	AllowedHeaders []string `mapstructure:"allowed_headers"`
	AllowAnyHeader bool     `mapstructure:"allow_any_header"`
	ExposedHeaders []string `mapstructure:"exposed_headers"`
}

type HttpTlsSetting struct {
	IsEnabled         bool   `mapstructure:"is_enabled"`
	CertFile          string `mapstructure:"cert_file"`
	KeyFile           string `mapstructure:"key_file"`
	RequireClientCert bool   `mapstructure:"require_client_cert"`
	MinVersion        string `mapstructure:"min_version"`
}

type HttpRateLimitSetting struct {
	Enabled           bool                                    `mapstructure:"enabled"`
	DefaultLimit      int                                     `mapstructure:"default_limit"`
	DefaultWindowSecs int                                     `mapstructure:"default_window_secs"`
	Endpoints         map[string]HttpEndpointRateLimitSetting `mapstructure:"endpoints"`
}

type HttpEndpointRateLimitSetting struct {
	Limit      int `mapstructure:"limit"`
	WindowSecs int `mapstructure:"window_secs"`
}
