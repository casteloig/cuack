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
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	do "cuack/pkg/digitalocean"

	"github.com/digitalocean/godo"
	"github.com/pterm/pterm"
	"github.com/sirupsen/logrus"
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
			pterm.Error.Println(err)
		}
		err = do.GetTokenFromFile()
		if err != nil {
			pterm.Error.Println(err)
		}

		var yamlFile []byte
		file, _ := cmd.Flags().GetString("file")
		if file != "" {
			sel, _ := cmd.Flags().GetString("select")
			if sel != "" {
				yamlFile, err = yamlToStruct(file, sel)
				if err != nil {
					pterm.Error.Println(err)
					logrus.Exit(1)
				}
			} else {
				yamlFile, err = yamlToStructFirst(file)
				if err != nil {
					pterm.Error.Println(err)
					logrus.Exit(1)
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
				p, _ := pterm.DefaultProgressbar.WithTotal(3).Start()
				p.Title = "Creating droplet"
				p.Increment()
				_, err := do.CreateDropletWithSSH(client, ctx, slugDroplet)
				if err != nil {
					pterm.Error.Println(err)
					logrus.Exit(1)
				}
				pterm.Success.Println("Droplet created successfully")

				// Create server inside the droplet
				p.Title = "Creating server"
				p.Increment()
				createTempYaml(yamlFile)
				ip, err := do.CreateServer(client, ctx)
				if err != nil {
					pterm.Error.Println(err)
					logrus.Exit(1)
				}

				p.Title = "Creating droplet"
				p.Increment()
				pterm.Success.Println("Server created successfully on ", ip)
				pterm.Info.Println("Server successfully created!")
				logrus.WithFields(logrus.Fields{
					"command":    "create",
					"final-name": do.Servers.Name,
					"file":       file,
					"ip":         ip,
				}).Info("Sucesfully created droplet and server")
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

func yamlToStruct(file string, name string) ([]byte, error) {
	var fileContent []byte

	if strings.HasPrefix(file, "https://") || strings.HasPrefix(file, "http://") {
		file = strings.Replace(file, "blob", "raw", 1)
		resp, err := http.Get(file)
		if err != nil {
			return nil, errors.New("error downloading yaml file")
		}
		defer resp.Body.Close()
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		fileContent = buf.Bytes()

	} else {
		var err error
		fileContent, err = ioutil.ReadFile(file)
		if err != nil {
			return nil, errors.New("error reading yaml file")
		}
	}

	allYamlBytes, err := splitYAML(fileContent)
	if err != nil {
		return nil, errors.New("sintactical error in yaml file")
	}

	var fileStructSelected []byte
	var eachYaml do.Server
	for _, y := range allYamlBytes {
		fileStructSelected = y
		yaml.Unmarshal(y, &eachYaml)
		if eachYaml.Name == name {
			do.Servers = eachYaml
			break
		}
	}

	return fileStructSelected, nil
}

func createTempYaml(fileContent []byte) error {
	fileDir := "/tmp/" + do.Servers.Name
	err := ioutil.WriteFile(fileDir, fileContent, 0755)
	if err != nil {
		return errors.New("error creating temp yaml file")
	}

	return nil
}

func yamlToStructFirst(file string) ([]byte, error) {
	var fileContent []byte

	if strings.HasPrefix(file, "https://") || strings.HasPrefix(file, "http://") {
		file = strings.Replace(file, "blob", "raw", 1)
		resp, err := http.Get(file)
		if err != nil {
			return nil, errors.New("error downloading yaml file")
		}
		defer resp.Body.Close()
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		fileContent = buf.Bytes()

	} else {
		var err error
		fileContent, err = ioutil.ReadFile(file)
		if err != nil {
			return nil, errors.New("error reading yaml file")
		}
	}

	allYamlBytes, err := splitYAML(fileContent)
	if err != nil {
		return nil, errors.New("sintactical error in yaml file")
	}

	yaml.Unmarshal(allYamlBytes[0], &do.Servers)

	return allYamlBytes[0], nil
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
