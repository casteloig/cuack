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
	"bufio"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	do "cuack/digitalocean"

	"github.com/digitalocean/godo"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Set the basic config",
	Long: `Set the key token of Digital Ocean and asks for the preffered region.
	
	This step is needed for creating a server`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("init called")

		fmt.Println("Enter the token of Digital Ocean:")
		reader := bufio.NewReader(os.Stdin)
		tokenDO, _ := reader.ReadString('\n')
		tokenDO = strings.Trim(tokenDO, "\n")

		client := godo.NewFromToken(tokenDO)
		ctx := context.TODO()

		regions := listRegions(client, ctx)

		regionPref := "lon1"
		fmt.Println("Enter the prefered region slug (default lon1)")

		auxRegion, _ := reader.ReadString('\n')
		auxRegion = strings.Trim(auxRegion, "\n")
		if auxRegion != "" {
			regionSlug, err := selectRegion(regions, auxRegion)
			if regionSlug != "" && err == nil {
				regionPref = auxRegion
			} else {
				fmt.Println(err)
			}
		}

		createInitFile(tokenDO, regionPref)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	//initCmd.Flags().StringVarP(&keyDO, "key", "k", "", "Key token of DigitalOcean")
	//initCmd.Flags().StringVarP(&regionPref, "region", "r", "london", "Prefered region")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func createInitFile(key string, region string) {
	home, _ := os.UserHomeDir()
	dir := home + "/.config/cuack.config"

	err := ioutil.WriteFile(dir, []byte("key "+key+"\n"+"region "+region), 0755)
	if err != nil {
		fmt.Println(err)
	}
}

func listRegions(client *godo.Client, ctx context.Context) []godo.Region {
	regions, _ := do.GetAvailableRegions(client, ctx)

	for _, region := range regions {
		fmt.Println(region.Slug + " (" + region.Name + ")")
	}
	return regions
}

func selectRegion(regions []godo.Region, slug string) (string, error) {
	for _, region := range regions {
		if strings.Compare(region.Slug, slug) == 0 {
			return region.Slug, nil
		}
	}
	return "", errors.New("that slug does not exist")
}
