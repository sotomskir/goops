package utils

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"strings"
)

func ViperValidate(name string, flagName string, envVarName string) {
	if viper.GetString(name) == "" {
		logrus.Fatalf("Error: empty %s. Please use --%s flag or %s env variable\n", name, flagName, envVarName)
	}
}

func ViperValidateEnv(envVarName... string) {
	errors := make([]string, 0)
	for _, v := range envVarName {
		value := viper.GetString(v)
		logrus.Tracef("ViperValidateEnv %s=%s\n", v, value)
		if value == "" {
			errors = append(errors, v)
		}
	}
	if len(errors) > 0 {
		logrus.Fatalf("required variables not set: %s", strings.Join(errors, ", "))
	}
}
