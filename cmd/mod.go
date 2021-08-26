/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
	do "cuack/digitalocean"
	mine "cuack/minecraft"
	"fmt"
	"strconv"
	"strings"

	"github.com/digitalocean/godo"
	"github.com/spf13/cobra"
)

// modCmd represents the mod command
var modCmd = &cobra.Command{
	Use:   "mod",
	Short: "A brief description of your command",
	Long:  `Droplets cannot be modified, just the gameservers`,
	Run: func(cmd *cobra.Command, args []string) {
		do.GetRegionFromFile()
		do.GetTokenFromFile()

		file, _ := cmd.Flags().GetString("file")
		fmt.Println(file)
		if file != "" {
			yamlToStruct(file)
		}

		client := godo.NewFromToken(do.Token)
		ctx := context.TODO()

		// Iterate over yaml configs
		for index, sv := range do.Servers {
			fmt.Println("Server number " + strconv.Itoa(index+1))
			// So far every server provider must be "digitalocean"
			if sv.Sv.NameProv == "digitalocean" {
				fmt.Println("Checking if it exists")
				exists, err := do.CheckDropletExists(client, ctx, sv.Sv.Name)
				if err != nil {
					fmt.Println(err)
				}
				if exists {
					// Create request = MINECRAFT
					if strings.ToLower(sv.Sv.NameGame) == "minecraft" {
						fmt.Println("Updating")
						mine.UpdateServer(client, ctx, sv)
					}
				}

			}

		}
	},
}

func init() {
	rootCmd.AddCommand(modCmd)
	modCmd.Flags().StringVarP(&file, "file", "f", "", "new config file of the server")
	modCmd.MarkFlagRequired("file")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// modCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// modCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
