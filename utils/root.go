package utils

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"strings"
)

func ViperValidate(name string, flagName string, envVarName string) {
	if viper.GetString(name) == "" {
		logrus.Fatalf("Error: empty %s. Please use --%s flag or %s env variable\n", name, flagName, envVarName)
	}
}

func ViperValidateEnv(envVarName ...string) {
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

func SaveExportInt(variableName string, value int) {
	saveExport(fmt.Sprintf("export %s=%d", variableName, value))
}

func SaveExportString(variableName string, value string) {
	saveExport(fmt.Sprintf("export %s=%s", variableName, value))
}

func saveExport(s string) {
	f, err := os.OpenFile(".goops.env", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		logrus.Fatalln(err)
	}

	defer f.Close()
	fmt.Println(s)

	str := fmt.Sprintln(s)
	if _, err = f.WriteString(str); err != nil {
		logrus.Fatalln(err)
	}
}
