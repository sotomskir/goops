// Copyright © 2019 Robert Sotomski <sotomskie@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"github.com/fatih/color"
	"github.com/mitchellh/go-homedir"
	"github.com/sotomskir/gitlab-cli/execService"
	"github.com/sotomskir/gitlab-cli/gitService"
	"github.com/sotomskir/gitlab-cli/gitlabApi"
	"github.com/sotomskir/gitlab-cli/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path"
)

var cfgFile string
var noColor bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gitlab-cli",
	Short: "Gitlab REST API command line client.",
	TraverseChildren: true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logger.ErrorLn(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gitlab-cli.yaml)")

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().StringP("server", "s", "", "Gitlab API Url")
	rootCmd.Flags().StringP("token", "a", "", "Gitlab API auth token")
	viper.BindPFlag("api_v4_url", rootCmd.Flags().Lookup("server"))
	viper.BindPFlag("build_token", rootCmd.Flags().Lookup("token"))
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "Disable ANSI color output")
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	//rootCmd.PersistentFlags().StringVarP(&user, "user", "u", "admin", "gitlab username")
	//rootCmd.PersistentFlags().StringVarP(&password, "password", "p", "admin", "gitlab password")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	gitService.Initialize(execService.Service{})
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if noColor {
		color.NoColor = true
	}
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			logger.ErrorLn(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".gitlab-cli" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".gitlab-cli")
	}

	viper.SetEnvPrefix("CI")
	viper.AutomaticEnv() // read in environment variables that match
	gitlabApi.Initialize()
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		//fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		home, err := homedir.Dir()
		if err != nil {
			logger.ErrorLn(err)
			os.Exit(1)
		}
		viper.WriteConfigAs(path.Join(home, "/.gitlab-cli.yaml"))
	}
}
