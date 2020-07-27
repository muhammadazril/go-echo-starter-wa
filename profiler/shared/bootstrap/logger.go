package bootstrap

import (
	"encoding/json"
	"fmt"

	"go.uber.org/zap"
)

// InitLogger will return initated logger
func InitLogger(appname string) *zap.Logger {
	cfg := zap.NewProductionConfig()
	rawJSON := []byte(`{
		"level": "debug",
		"encoding": "json",
		"outputPaths": ["stdout", "/tmp/log/` + appname + `.log"],
		"errorOutputPaths": ["stderr", "/tmp/log/` + appname + `.log"],
		"initialFields": {"appname": "` + appname + `"},
		"encoderConfig": {
			"messageKey": "msg",
			"levelKey": "level",
			"levelEncoder": "lowercase"
		}
	}`)
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		panic(err)
	}
	logger, err := cfg.Build()
	if err != nil {
		panic(fmt.Sprintf("Failed Initiate Logger | %s", err.Error()))
	}
	return logger
}
