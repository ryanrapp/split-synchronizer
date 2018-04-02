package v2

// Config values grouped by execution mode
type Config struct {
	Main     Main     `yaml:"main" split-cli-option-group:"true"`
	Log      Log      `yaml:"log" split-cli-option-group:"true"`
	Proxy    Proxy    `yaml:"proxy" split-cli-option-group:"true"`
	Producer Producer `yaml:"producer" split-cli-option-group:"true"`
	Admin    Admin    `yaml:"admin" split-cli-option-group:"true"`
}

//---------------------------------------------------------------------------------
//  Common configuration Section
//---------------------------------------------------------------------------------

// Main configuration common among execution modes
type Main struct {
	APIKey                 string             `yaml:"apikey" split-cli-option:"api-key" split-default-value:"YOUR API KEY" split-cli-description:"Your Split API-KEY"`
	SplitsRefreshRate      int                `yaml:"splitsRefreshRate" split-cli-option:"split-refresh-rate" split-default-value:"60" split-cli-description:"Refresh rate of splits fetcher"`
	SegmentsRefreshRate    int                `yaml:"segmentsRefreshRate" split-default-value:"60" split-cli-option:"segment-refresh-rate" split-cli-description:"Refresh rate of segments fetcher"`
	ImpressionsRefreshRate int                `yaml:"impressionsRefreshRate" split-default-value:"60" split-cli-option:"impressions-post-rate" split-cli-description:"Post rate of impressions recorder"`
	ImpressionsPerPost     int                `yaml:"impressionsPerPost" split-cli-option:"impressions-per-post" split-default-value:"1000" split-cli-description:"Number of impressions to send in a POST request"`
	ImpressionsThreads     int                `yaml:"impressionsThreads" split-default-value:"1" split-cli-option:"impressions-recorder-threads" split-cli-description:"Number of impressions recorder threads"`
	EventsPushRate         int                `yaml:"eventsPushRate" split-default-value:"60" split-cli-option:"events-push-rate" split-cli-description:"Post rate of event recorder (seconds)"`
	EventsConsumerReadSize int                `yaml:"eventsConsumerReadSize" split-default-value:"10000" split-cli-option:"events-consumer-read-size" split-cli-description:"Events queue read size"`
	EventsConsumerThreads  int                `yaml:"eventsConsumerThreads" split-default-value:"1" split-cli-option:"events-consumer-threads" split-cli-description:"Number of events consumer threads"`
	MetricsRefreshRate     int                `yaml:"metricsRefreshRate" split-default-value:"60" split-cli-option:"metrics-post-rate" split-cli-description:"Post rate of metrics recorder"`
	HTTPTimeout            int                `yaml:"httpTimeout" split-default-value:"60" split-cli-option:"http-timeout" split-cli-description:"Timeout specifies a time limit for requests"`
	ImpressionListener     ImpressionListener `yaml:"impressionListener" split-cli-option-group:"true"`
}

// ImpressionListener represents configuration for impression bulk poster
type ImpressionListener struct {
	Endpoint string `yaml:"endpoint" split-default-value:"" split-cli-option:"impression-listener-endpoint" split-cli-description:"HTTP endpoint where impression bulks will be posted"`
}

//---------------------------------------------------------------------------------
//  Log configuration Section
//---------------------------------------------------------------------------------

// Log configuration values
type Log struct {
	Verbose          bool   `yaml:"verbose" split-default-value:"false" split-cli-option:"log-verbose" split-cli-description:"Enable verbose mode"`
	Debug            bool   `yaml:"debug" split-default-value:"false" split-cli-option:"log-debug" split-cli-description:"Enable debug mode"`
	Stdout           bool   `yaml:"stdout" split-default-value:"false" split-cli-option:"log-stdout" split-cli-description:"Enable log standard output"`
	File             string `yaml:"file" split-default-value:"/tmp/split-agent.log" split-cli-option:"log-file" split-cli-description:"Set the log file"`
	FileMaxSizeBytes int64  `yaml:"fileMaxSizeBytes" split-cli-option:"log-file-max-size" split-default-value:"2000000" split-cli-description:"Max file log size in bytes"`
	FileBackupCount  int    `yaml:"fileBackupCount" split-cli-option:"log-file-backup-count" split-default-value:"3" split-cli-description:"Number of last log files to keep in filesystem"`
	SlackChannel     string `yaml:"slackChannel" split-default-value:"" split-cli-option:"log-slack-channel" split-cli-description:"Set the Slack channel or user"`
	SlackWebhookURL  string `yaml:"slackWebhookURL" split-default-value:"" split-cli-option:"log-slack-webhook-url" split-cli-description:"Set the Slack webhook url"`
}

//---------------------------------------------------------------------------------
//  Proxy configuration Section
//---------------------------------------------------------------------------------

// Proxy configuration values
type Proxy struct {
	Port               int      `yaml:"port" split-default-value:"3000" split-cli-option:"proxy-port" split-cli-description:"Proxy port to listen connections"`
	PersistInFilePath  string   `yaml:"persistInFilePath" split-default-value:"" split-cli-option:"proxy-mmap-path" split-cli-description:"File path to persist memory in proxy mode"`
	ImpressionsMaxSize int      `yaml:"impressionsMaxSize" split-default-value:"10485760" split-cli-option:"proxy-impressions-max-size" split-cli-description:"Max size, in bytes, to send impressions in proxy mode"`
	EventsMaxSize      int      `yaml:"eventsMaxSize" split-default-value:"10485760" split-cli-option:"proxy-events-max-size" split-cli-description:"Max size, in bytes, to send events in proxy mode"`
	SDKAPIKeys         []string `yaml:"sdkAPIKeys,flow" split-default-value:"SDK_API_KEY" split-cli-option:"proxy-apikeys" split-cli-description:"List of allowed custom API Keys for SDKs"`
}

//---------------------------------------------------------------------------------
//  Production configuration Section
//---------------------------------------------------------------------------------

// Producer configuration section
type Producer struct {
	Storage Storage `yaml:"storage" split-cli-option-group:"true"`
}

// Storage configuration for producer execution mode
type Storage struct {
	Redis Redis `yaml:"redis" split-cli-option-group:"true"`
}

// Redis instance information
type Redis struct {
	Host   string `yaml:"host" split-default-value:"localhost" split-cli-option:"redis-host" split-cli-description:"Redis server hostname"`
	Port   int    `yaml:"port" split-default-value:"6379" split-cli-option:"redis-port" split-cli-description:"Redis Server port"`
	Db     int    `yaml:"db" split-default-value:"0" split-cli-option:"redis-db" split-cli-description:"Redis DB"`
	Pass   string `yaml:"password" split-default-value:"" split-cli-option:"redis-pass" split-cli-description:"Redis password"`
	Prefix string `yaml:"prefix" split-default-value:"" split-cli-option:"redis-prefix" split-cli-description:"Redis key prefix"`

	// The network type, either tcp or unix.
	// Default is tcp.
	Network string `yaml:"network" split-default-value:"tcp" split-cli-option:"redis-network" split-cli-description:"Redis network protocol"`

	// Maximum number of retries before giving up.
	// Default is to not retry failed commands.
	MaxRetries int `yaml:"maxRetries" split-default-value:"0" split-cli-option:"redis-max-retries" split-cli-description:"Redis connection max retries"`

	// Dial timeout for establishing new connections.
	// Default is 5 seconds.
	DialTimeout int `yaml:"dialTimeout" split-default-value:"5" split-cli-option:"redis-dial-timeout" split-cli-description:"Redis connection dial timeout"`

	// Timeout for socket reads. If reached, commands will fail
	// with a timeout instead of blocking.
	// Default is 10 seconds.
	ReadTimeout int `yaml:"readTimeout" split-default-value:"10" split-cli-option:"redis-read-timeout" split-cli-description:"Redis connection read timeout"`

	// Timeout for socket writes. If reached, commands will fail
	// with a timeout instead of blocking.
	// Default is 3 seconds.
	WriteTimeout int `yaml:"writeTimeout" split-default-value:"5" split-cli-option:"redis-write-timeout" split-cli-description:"Redis connection write timeout"`

	// Maximum number of socket connections.
	// Default is 10 connections.
	PoolSize int `yaml:"poolSize" split-default-value:"10" split-cli-option:"redis-pool" split-cli-description:"Redis connection pool size"`
}

//---------------------------------------------------------------------------------
//  Admin tools configuration Section
//---------------------------------------------------------------------------------

// Admin tools configuration values
type Admin struct {
	Port           int       `yaml:"port" split-default-value:"3010" split-cli-option:"admin-port" split-cli-description:"Admin port to listen connections"`
	DashboardTitle string    `yaml:"dashboardTitle" split-default-value:"" split-cli-option:"admin-dashboard-title" split-cli-description:"Descriptive title to be shown in Admin Dashboard"`
	BasicAuth      BasicAuth `yaml:"basicAuth" split-cli-option-group:"true"`
}

// BasicAuth basic HTTP authentication for admin tools
type BasicAuth struct {
	Username string `yaml:"username" split-default-value:"" split-cli-option:"admin-username" split-cli-description:"HTTP basic auth username for admin endpoints"`
	Password string `yaml:"password" split-default-value:"" split-cli-option:"admin-password" split-cli-description:"HTTP basic auth password for admin endpoints"`
}
