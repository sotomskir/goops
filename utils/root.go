package utils

import (
	"github.com/sotomskir/gitlab-cli/logger"
	"github.com/spf13/viper"
	"os"
)

func ViperValidate(name string, flagName string, envVarName string) {
	if viper.GetString(name) == "" {
		logger.ErrorF("Error: empty %s. Please use --%s flag or %s env variable\n", name, flagName, envVarName)
		os.Exit(1)
	}
}
