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
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	do "cuack/digitalocean"
	mine "cuack/minecraft"

	"github.com/digitalocean/godo"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v3"
)

var file string

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates and deploys a server",
	Long: `Creates a droplet with docker installed.
	
	Then it takes the yaml file as an argument and deploys inside the droplet a 
	running server.
	
	It can create and deploy more than one server. Check the templates and default packages`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("create called")

		do.GetRegionFromFile()
		do.GetTokenFromFile()

		fmt.Println(do.Token)
		fmt.Println(do.Region)

		file, _ := cmd.Flags().GetString("file")
		fmt.Println(file)
		if file != "" {
			yamlToStruct(file)
		}

		client := godo.NewFromToken(do.Token)
		ctx := context.TODO()

		// Iterate over all servers specified in the yaml file
		for index, sv := range do.Servers {
			slugDroplet := "s-"
			cpu := sv.Sv.Provider.Cpu
			ram := sv.Sv.Provider.Ram
			slugDroplet = slugDroplet + strconv.Itoa(cpu) + "vcpu-" + strings.ToLower(ram)

			fmt.Println("Server number " + strconv.Itoa(index))
			// So far every server provider must be "digitalocean"
			if sv.Sv.NameProv == "digitalocean" {
				fmt.Println("Checking if it can be created")
				exists, err := do.CheckDropletExists(client, ctx, sv.Sv.Name)
				if err != nil {
					fmt.Println(err)
				}
				// If not reached max number of droplets and the droplet with that name does not exists yet...
				if !exists && (do.GetMaxDroplets(client, ctx)-do.GetNumberDroplets(client, ctx) >= 1) {
					// Create droplet
					fmt.Println("Creating droplet ...")
					_, err := do.CreateDropletWithSSH(client, ctx, sv.Sv.Name, do.Region, slugDroplet, sv.Sv.Provider.SshName, sv.Sv.NameGame)
					if err != nil {
						fmt.Println(err)
					}
					// Create request = MINECRAFT
					if strings.ToLower(sv.Sv.NameGame) == "minecraft" {
						fmt.Println("Creating minecraft server")
						mine.CreateServer(client, ctx, sv)
					}
				}

			}

		}

	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().StringVarP(&file, "file", "f", "", "config file of the server")
	createCmd.MarkFlagRequired("file")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func yamlToStruct(file string) error {
	fileContent, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	yaml.Unmarshal(fileContent, &do.Servers)
	fmt.Println(do.Servers)
	return nil
}

func GetServers() []do.ServerGeneral {
	return do.Servers
}

func GetToken() string {
	return do.Token
}

func GetRegion() string {
	return do.Region
}
