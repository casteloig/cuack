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
	"errors"
	"log"

	do "cuack/pkg/digitalocean"

	"github.com/digitalocean/godo"
	"github.com/spf13/cobra"
)

var name string

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Deletes a droplet",
	Long:  `Deletes a droplet taking as a flag the name of the droplet`,
	Run: func(cmd *cobra.Command, args []string) {

		err := do.GetTokenFromFile()
		if err != nil {
			log.Println(err)
		}

		if len(args) > 0 {
			client := godo.NewFromToken(do.Token)
			ctx := context.TODO()

			serverName := args[0]
			log.Println("Deleting " + serverName)
			err = do.DeleteDropletByName(client, ctx, serverName)
			if err != nil {
				log.Println(err)
			}
		} else {
			err = errors.New("not enough arguments")
			log.Println(err)
		}

	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
