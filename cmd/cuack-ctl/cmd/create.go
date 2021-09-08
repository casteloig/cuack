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
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	do "cuack/pkg/digitalocean"

	"github.com/digitalocean/godo"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v3"
)

var file string
var sel string

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates and deploys a server",
	Long: `Creates a droplet with docker installed.
	
	Then it takes the yaml file as an argument and deploys inside the droplet a 
	running server.`,

	Run: func(cmd *cobra.Command, args []string) {

		err := do.GetRegionFromFile()
		if err != nil {
			log.Println(err)
		}
		err = do.GetTokenFromFile()
		if err != nil {
			log.Println(err)
		}

		file, _ := cmd.Flags().GetString("file")
		if file != "" {
			sel, _ := cmd.Flags().GetString("select")
			if sel != "" {
				err = yamlToStruct(file, sel)
				if err != nil {
					log.Println(err)
				}
			} else {
				err = yamlToStructFirst(file)
				if err != nil {
					log.Println(err)
				}
			}
		}

		client := godo.NewFromToken(do.Token)
		ctx := context.TODO()

		slugDroplet := "s-"
		cpu := do.Servers.Cpu
		ram := do.Servers.Ram
		slugDroplet = slugDroplet + strconv.Itoa(cpu) + "vcpu-" + strings.ToLower(ram)

		// So far every server provider must be "digitalocean"
		if do.Servers.Provider.NameProv == "digitalocean" {

			// If not reached max number of droplets
			if do.GetMaxDroplets(client, ctx)-do.GetNumberDroplets(client, ctx) >= 1 {
				// Create droplet
				log.Println("Trying to create droplet ...")
				_, err := do.CreateDropletWithSSH(client, ctx, slugDroplet)
				if err != nil {
					log.Println(err)
				}

				// Create server inside the droplet
				log.Println("Creating server ...")
				err = do.CreateServer(client, ctx)
				if err != nil {
					log.Println(err)
				}
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().StringVarP(&file, "file", "f", "", "config file of the server")
	createCmd.MarkFlagRequired("file")

	createCmd.Flags().StringVarP(&sel, "select", "s", "", "select one option from the yaml by his name. If there is more than one.")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func yamlToStruct(file string, name string) error {
	var fileContent []byte

	if strings.HasPrefix(file, "https://") || strings.HasPrefix(file, "http://") {
		file = strings.Replace(file, "blob", "raw", 1)
		resp, err := http.Get(file)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		fileContent = buf.Bytes()

	} else {
		var err error
		fileContent, err = ioutil.ReadFile(file)
		if err != nil {
			return err
		}
	}

	allYamlBytes, err := splitYAML(fileContent)
	if err != nil {
		return err
	}

	var eachYaml do.Server
	for _, y := range allYamlBytes {
		yaml.Unmarshal(y, &eachYaml)
		if eachYaml.Name == name {
			do.Servers = eachYaml
			break
		}
	}

	log.Println(do.Servers)

	return nil
}

func yamlToStructFirst(file string) error {
	var fileContent []byte

	if strings.HasPrefix(file, "https://") || strings.HasPrefix(file, "http://") {
		file = strings.Replace(file, "blob", "raw", 1)
		resp, err := http.Get(file)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		fileContent = buf.Bytes()

	} else {
		var err error
		fileContent, err = ioutil.ReadFile(file)
		if err != nil {
			return err
		}
	}

	allYamlBytes, err := splitYAML(fileContent)
	if err != nil {
		return err
	}

	yaml.Unmarshal(allYamlBytes[0], &do.Servers)

	log.Println(do.Servers)

	return nil
}

func splitYAML(resources []byte) ([][]byte, error) {

	dec := yaml.NewDecoder(bytes.NewReader(resources))

	var res [][]byte
	for {
		var value interface{}
		err := dec.Decode(&value)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		valueBytes, err := yaml.Marshal(value)
		if err != nil {
			return nil, err
		}
		res = append(res, valueBytes)
	}
	return res, nil
}
