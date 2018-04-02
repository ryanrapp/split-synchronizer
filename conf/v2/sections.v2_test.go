package v2

import (
	"testing"

	yaml "gopkg.in/yaml.v2"
)

var yaml4Test = `
# Main: configuration shared between all execution modes
main:
  apiKey: "YOUR API KEY"
  splitsRefreshRate: 60
  segmentsRefreshRate: 60
  impressionsRefreshRate: 60
  impressionsPerPost: 1000
  impressionsThreads: 1
  eventsPushRate: 60
  eventsConsumerReadSize: 10000
  eventsConsumerThreads: 1
  metricsRefreshRate: 60
  httpTimeout: 60
  impressionListener:
    endpoint: ""

# Proxy: settings to apply only in proxy mode execution
proxy:
  port: 3000
  persistInFilePath": ""
  impressionsMaxSize": 10485760
  eventsMaxSize": 10485760
  sdkAPIKeys:
   - "SDK_API_KEY"

# Producer: settings to apply only in sync mode execution
producer:
  storage:
    redis:
      host: "localhost"
      port: 6379
      db: 0
      password: ""
      prefix: ""
      network: "tcp"
      maxRetries: 0
      dialTimeout: 5
      readTimeout: 10
      writeTimeout: 5
      poolSize": 10

# Admin services configuration
admin:
  port: 3010
  dashboardTitle: ""
  basicAuth:
    username: ""
    password: ""

# Log configuration
log:
  verbose: false
  debug: true
  stdout: false
  file: "/tmp/split-sync.log"
  fileMaxSizeBytes: 2000000
  fileBackupCount: 3
  slackChannel: ""
  slackWebhookURL: ""
`

func TestConfigStruct(t *testing.T) {
	config := Config{}

	err := yaml.Unmarshal([]byte(yaml4Test), &config)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	if config.Log.FileMaxSizeBytes != 2000000 {
		t.Error("Error reading config.Log.FileMaxSizeBytes")
	}

	if !config.Log.Debug {
		t.Error("Debug should be enabled")
	}

	if config.Producer.Storage.Redis.Host != "localhost" {
		t.Error("Redis Host should be localhost")
	}
}
