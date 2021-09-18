/*
Copyright © 2021 NAME HERE <casteloig@outlook.es>

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
	web "cuack"
	do "cuack/pkg/digitalocean"
	"fmt"
	"log"

	"github.com/digitalocean/godo"
	"github.com/spf13/cobra"
)

// inspectCmd represents the inspect command
var inspectCmd = &cobra.Command{
	Use:   "inspect",
	Short: "Get some usefull information of a droplet",
	Long:  `Get some usefull information of a droplet`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) > 0 {
			serverName := args[0]

			err := do.GetTokenFromFile()
			if err != nil {
				log.Fatal(err)
			}

			client := godo.NewFromToken(do.Token)
			ctx := context.TODO()

			inspected, err := do.InspectDroplet(client, ctx, serverName)
			if err != nil {
				fmt.Println(err)
			}

			err = web.CreateWebsite(inspected)
			if err != nil {
				log.Println(err)
			}

		}
	},
}

func init() {
	rootCmd.AddCommand(inspectCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// inspectCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// inspectCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
