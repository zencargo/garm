// Copyright 2022 Cloudbase Solutions SRL
//
//    Licensed under the Apache License, Version 2.0 (the "License"); you may
//    not use this file except in compliance with the License. You may obtain
//    a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//    WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//    License for the specific language governing permissions and limitations
//    under the License.

package cmd

import (
	"fmt"
	"garm/cmd/garm-cli/common"
	"garm/cmd/garm-cli/config"
	"garm/params"
	"strings"

	"github.com/spf13/cobra"
)

var (
	loginProfileName string
	loginURL         string
	loginPassword    string
	loginUserName    string
	loginFullName    string
	loginEmail       string
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:          "login",
	SilenceUsage: true,
	Short:        "Log into a manager",
	Long: `Performs login for this machine on a remote garm:

garm-cli login --name=dev --url=https://runner.example.com`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if cfg != nil {
			if cfg.HasManager(loginProfileName) {
				return fmt.Errorf("a manager with name %s already exists in your local config", loginProfileName)
			}
		}

		if err := promptUnsetLoginVariables(); err != nil {
			return err
		}

		url := strings.TrimSuffix(loginURL, "/")
		loginParams := params.PasswordLoginParams{
			Username: loginUserName,
			Password: loginPassword,
		}

		resp, err := cli.Login(url, loginParams)
		if err != nil {
			return err
		}

		cfg.Managers = append(cfg.Managers, config.Manager{
			Name:    loginProfileName,
			BaseURL: url,
			Token:   resp,
		})
		cfg.ActiveManager = loginProfileName

		if err := cfg.SaveConfig(); err != nil {
			return err
		}
		return nil
	},
}

func promptUnsetLoginVariables() error {
	var err error
	if loginUserName == "" {
		loginUserName, err = common.PromptString("Username")
		if err != nil {
			return err
		}
	}

	if loginPassword == "" {
		loginPassword, err = common.PromptPassword("Password")
		if err != nil {
			return err
		}
	}
	return nil
}

func init() {
	rootCmd.AddCommand(loginCmd)

	loginCmd.Flags().StringVarP(&loginProfileName, "name", "n", "", "A name for this runner manager")
	loginCmd.Flags().StringVarP(&loginURL, "url", "a", "", "The base URL for the runner manager API")
	loginCmd.Flags().StringVarP(&loginUserName, "username", "u", "", "Username to log in as")
	loginCmd.Flags().StringVarP(&loginPassword, "password", "p", "", "The user passowrd")

	loginCmd.MarkFlagRequired("name")
	loginCmd.MarkFlagRequired("url")
}
