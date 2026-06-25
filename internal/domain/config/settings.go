package config

// ===
// Configuration file structure for system
// ===
type SystemConfig struct {
	System        SystemSetting                 `mapstructure:"system"`
	HttpServer    HttpServerSetting             `mapstructure:"http_server"`
	GrpcServer    GrpcServerSetting             `mapstructure:"grpc"`
	GrpcClient    map[string]GrpcClientSetting  `mapstructure:"grpc_client"`
	Redis         RedisSetting                  `mapstructure:"redis"`
	ScyllaDb      ScyllaDbSetting               `mapstructure:"scylladb"`
	Postgres      PostgresSetting               `mapstructure:"postgres"`
	Kafka         KafkaSetting                  `mapstructure:"kafka"`
	Logger        LoggerSetting                 `mapstructure:"logger"`
	Observability ObservabilitySetting          `mapstructure:"observability"`
}

// ===
// Setting configuration structure used in the configuration file
// ===

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

// Grpc client structure setting
type GrpcClientSetting struct {
	Connection  ConnectionSetting    `mapstructure:"connection"`
	Timeouts    TimeoutsSetting      `mapstructure:"timeouts"`
	Pool        PoolSetting          `mapstructure:"pool"`
	Keepalive   KeepaliveSetting     `mapstructure:"keepalive"`
	Retry       RetrySetting         `mapstructure:"retry"`
	Tls         GrpcClientTlsSetting `mapstructure:"tls"`
	Discovery   string               `mapstructure:"discovery"`
	Addresses   []string             `mapstructure:"addresses"`
	ServiceName string               `mapstructure:"service_name"`
	Authority   string               `mapstructure:"authority"`
	Service     string               `mapstructure:"service"`
}

type ConnectionSetting struct {
	ConnectTimeoutMs int  `mapstructure:"connect_timeout_ms"`
	InitialBackoffMs int  `mapstructure:"initial_backoff_ms"`
	MaxBackoffMs     int  `mapstructure:"max_backoff_ms"`
	LazyConnect      bool `mapstructure:"lazy_connect"`
	TcpNodelay       bool `mapstructure:"tcp_nodelay"`
}

type TimeoutsSetting struct {
	RequestTimeoutMs  int `mapstructure:"request_timeout_ms"`
	StreamTimeoutMs   int `mapstructure:"stream_timeout_ms"`
	ShutdownTimeoutMs int `mapstructure:"shutdown_timeout_ms"`
}

type PoolSetting struct {
	ConnectionsPerEndpoint             int `mapstructure:"connections_per_endpoint"`
	MaxConcurrentStreamsPerConnection int `mapstructure:"max_concurrent_streams_per_connection"`
	IdleTimeoutMs                      int `mapstructure:"idle_timeout_ms"`
}

type KeepaliveSetting struct {
	Http2KeepaliveIntervalMs int  `mapstructure:"http2_keepalive_interval_ms"`
	Http2KeepaliveTimeoutMs  int  `mapstructure:"http2_keepalive_timeout_ms"`
	KeepaliveWhileIdle       bool `mapstructure:"keepalive_while_idle"`
	TcpKeepaliveMs           int  `mapstructure:"tcp_keepalive_ms"`
}

type RetrySetting struct {
	Enabled              bool     `mapstructure:"enabled"`
	MaxAttempts          int      `mapstructure:"max_attempts"`
	InitialBackoffMs     int      `mapstructure:"initial_backoff_ms"`
	MaxBackoffMs         int      `mapstructure:"max_backoff_ms"`
	BackoffMultiplier    float64  `mapstructure:"backoff_multiplier"`
	RetryableStatusCodes []string `mapstructure:"retryable_status_codes"`
	UseJitter            bool     `mapstructure:"use_jitter"`
}

type GrpcClientTlsSetting struct {
	IsEnabled         bool   `mapstructure:"is_enabled"`
	CertFile          string `mapstructure:"cert_file"`
	KeyFile           string `mapstructure:"key_file"`
	RequireClientCert bool   `mapstructure:"require_client_cert"`
	MinVersion        string `mapstructure:"min_version"`
}

// Redis structure setting
type RedisSetting struct {
	Type           int      `mapstructure:"type"`
	UseTls         bool     `mapstructure:"use_tls"`
	CertPath       string   `mapstructure:"cert_path"`
	KeyPath        string   `mapstructure:"key_path"`
	Password       string   `mapstructure:"password"`
	Db             int      `mapstructure:"db"`
	Host           string   `mapstructure:"host"`
	Port           int      `mapstructure:"port"`
	MasterName     string   `mapstructure:"master_name"`
	SentinelAddrs  []string `mapstructure:"sentinel_addrs"`
	Address        []string `mapstructure:"address"`
	RouteByLatency bool     `mapstructure:"route_by_latency"`
	RouteRandomly  bool     `mapstructure:"route_randomly"`
	PoolSize       int      `mapstructure:"pool_size"`
	MinIdleConns   int      `mapstructure:"min_idle_conns"`
	MaxRetries     int      `mapstructure:"max_retries"`
}

// ScyllaDb structure setting
type ScyllaDbSetting struct {
	Authentication  ScyllaDbAuthSetting `mapstructure:"authentication"`
	Address         []string            `mapstructure:"address"`
	Keyspace        string              `mapstructure:"keyspace"`
	Ssl             ScyllaDbSslSetting  `mapstructure:"ssl"`
	MaxIdleConns    int                 `mapstructure:"maxIdleConns"`
	MaxOpenConns    int                 `mapstructure:"maxOpenConns"`
	ConnMaxLifetime int                 `mapstructure:"connMaxLifetime"`
}

type ScyllaDbAuthSetting struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type ScyllaDbSslSetting struct {
	Enabled      bool   `mapstructure:"enabled"`
	CertfilePath string `mapstructure:"certfile_path"`
	Validate     bool   `mapstructure:"validate"`
	UserkeyPath  string `mapstructure:"userkey_path"`
	UsercertPath string `mapstructure:"usercert_path"`
}

// Postgres structure setting
type PostgresSetting struct {
	Address               []string `mapstructure:"address"`
	Database              string   `mapstructure:"database"`
	Username              string   `mapstructure:"username"`
	Password              string   `mapstructure:"password"`
	SslMode               string   `mapstructure:"sslmode"`
	SslCert               string   `mapstructure:"sslcert"`
	SslKey                string   `mapstructure:"sslkey"`
	SslRootCert           string   `mapstructure:"sslrootcert"`
	SslPassword           string   `mapstructure:"sslpassword"`
	AppName               string   `mapstructure:"appname"`
	ConnectionTimeout     int      `mapstructure:"connectionTimeout"`
	Tz                    string   `mapstructure:"tz"`
	MaxConns              int      `mapstructure:"maxConns"`
	MinConns              int      `mapstructure:"minConns"`
	MinIdleConns          int      `mapstructure:"minIdleConns"`
	MaxConnIdleTime       int      `mapstructure:"maxConnIdleTime"`
	MaxConnLifetimeJitter int      `mapstructure:"maxConnLifetimeJitter"`
	HealthCheckPeriod     int      `mapstructure:"healthCheckPeriod"`
}

// Kafka structure setting
type KafkaSetting struct {
	Brokers   []string           `mapstructure:"brokers"`
	Security  KafkaSecurity      `mapstructure:"security"`
	Producer  KafkaProducer      `mapstructure:"producer"`
	Consumers KafkaConsumerGroup `mapstructure:"consumers"`
}

type KafkaSecurity struct {
	Sasl KafkaSasl `mapstructure:"sasl"`
	Tls  KafkaTls  `mapstructure:"tls"`
}

type KafkaSasl struct {
	Enabled   bool   `mapstructure:"enabled"`
	Mechanism string `mapstructure:"mechanism"`
	Username  string `mapstructure:"username"`
	Password  string `mapstructure:"password"`
}

type KafkaTls struct {
	Enabled    bool   `mapstructure:"enabled"`
	SkipVerify bool   `mapstructure:"skip_verify"`
	CaFile     string `mapstructure:"ca_file"`
}

type KafkaProducer struct {
	CompressionType string `mapstructure:"compression_type"`
	Retries         int    `mapstructure:"retries"`
	RetryBackoffMs  int    `mapstructure:"retry_backoff_ms"`
	LingerMs        int    `mapstructure:"linger_ms"`
	BatchSize       int    `mapstructure:"batch_size"`
	BatchBytes      int    `mapstructure:"batch_bytes"`
	MaxAttempts     int    `mapstructure:"max_attempts"`
	Async           bool   `mapstructure:"async"`
	WriteTimeoutMs  int    `mapstructure:"write_timeout_ms"`
	ReadTimeoutMs   int    `mapstructure:"read_timeout_ms"`
	Balancer        string `mapstructure:"balancer"`
}

type KafkaConsumerGroup struct {
	GroupId             string `mapstructure:"group_id"`
	CommitIntervalMs    int    `mapstructure:"commit_interval_ms"`
	MinBytes            int    `mapstructure:"min_bytes"`
	MaxBytes            int    `mapstructure:"max_bytes"`
	MaxWaitMs           int    `mapstructure:"max_wait_ms"`
	ReadBatchTimeoutMs  int    `mapstructure:"read_batch_timeout_ms"`
	HeartbeatIntervalMs int    `mapstructure:"heartbeat_interval_ms"`
	SessionTimeoutMs    int    `mapstructure:"session_timeout_ms"`
	RebalanceTimeoutMs  int    `mapstructure:"rebalance_timeout_ms"`
	JoinGroupBackoffMs  int    `mapstructure:"join_group_backoff_ms"`
	ReadBackoffMinMs    int    `mapstructure:"read_backoff_min_ms"`
	ReadBackoffMaxMs    int    `mapstructure:"read_backoff_max_ms"`
	ReadLagIntervalMs   int    `mapstructure:"read_lag_interval_ms"`
	MaxAttempts         int    `mapstructure:"max_attempts"`
	QueueCapacity       int    `mapstructure:"queue_capacity"`
	RetentionTimeMs     int    `mapstructure:"retention_time_ms"`
}
