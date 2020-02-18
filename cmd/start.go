// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
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
	"fmt"
	"os"
	"time"

	"github.com/llimon/churndr/common"
	"github.com/llimon/churndr/controller"
	"github.com/llimon/churndr/notifier"
	"github.com/llimon/churndr/server"
	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start rest service",
	Long:  `Stuff Super Service`,
	Run: func(cmd *cobra.Command, args []string) {

		defer common.Logger.Sync() // flushes buffer, if any

		if !common.Config.DissableEmailNotifications {
			if common.Config.EmailSMTPServer == "" ||
				common.Config.EmailFrom == "" ||
				common.Config.EmailTo == "" { //||
				//common.Config.EmailLogin == "" ||
				//common.Config.EmailPassword == "" {
				//fmt.Println("Missing parameter(s) one of [email-from, email-to, email-login, email-password, smtp]")
				fmt.Println("Missing parameter(s) one of [email-from, email-to, smtp]")
				os.Exit(1)
			}
		}

		common.Sugar.Infof("Creating notified ticker....")

		ticker := time.NewTicker(time.Duration(common.Config.NotificationFrequency) * time.Second)
		go func() {
			for t := range ticker.C {
				//Call the periodic function here.
				notifier.KubeNotifierFunc(t)
			}
		}()

		quit := make(chan bool, 1)

		if !common.Config.NoAPIServer {
			go server.RESTServer()
		}

		go controller.RunChurnNotifierController()

		controller.KubeGetPods()

		// main will continue to wait until there is an entry in quit
		<-quit
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:

	startCmd.PersistentFlags().IntVarP(&common.Config.Port, "port", "p", 8080, "Port to listen for https requests")

}
