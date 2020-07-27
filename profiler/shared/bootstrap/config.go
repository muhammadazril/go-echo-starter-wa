package bootstrap

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// InitConfig will return initated config
func InitConfig() *viper.Viper {

	var (
		err     error
		prjPath string
	)

	// set the project path correctly, for unit test run (to read the exact file)
	// if not set, set with current dir (unit test will fail since cannot read the config properlly)
	prjPath = os.Getenv("APP_PATH")
	if prjPath == "" {
		prjPath = "."
	}

	config := viper.New()
	config.SetConfigName("config")
	config.AddConfigPath(prjPath)
	config.SetConfigType("yaml")

	err = config.ReadInConfig()
	if err != nil {
		panic(fmt.Sprintf("Failed read config file | %s", err.Error()))
	}

	return config
}
