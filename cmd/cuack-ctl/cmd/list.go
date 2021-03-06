/*
Copyright © 2021 Ignacio Castelo <casteloig@outlook.es>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"

	do "cuack/pkg/digitalocean"

	"github.com/digitalocean/godo"
	"github.com/pterm/pterm"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all cuack droplets of DigitalOcean",
	Long:  `List all cuack droplets of DigitalOcean indicating their name and IPv4 address`,
	Run: func(cmd *cobra.Command, args []string) {

		err := do.GetTokenFromFile()
		if err != nil {
			pterm.Error.Println(err)
			logrus.Exit(1)
		}

		client := godo.NewFromToken(do.Token)
		ctx := context.TODO()

		list, err := do.ListCuackDroplets(client, ctx)
		if err != nil {
			pterm.Error.Println(err)
			logrus.Exit(1)
		}

		pterm.DefaultSection.Println("Listsing droplets")
		logrus.WithFields(logrus.Fields{
			"command":  "list",
			"droplets": len(list),
		}).Info("Listing droplets")
		for ip, name := range list {
			pterm.Info.Println("Name: " + name + "\nIPv4: " + ip)
			println()
		}

	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
