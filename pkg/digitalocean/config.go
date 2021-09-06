package digitalocean

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/digitalocean/godo"
	"github.com/melbahja/goph"
)

var Servers Server
var Token string
var Region string

type Provider struct {
	NameProv string `yaml:"name-prov"`
	SshName  string `yaml:"ssh-name"`
	Cpu      int    `yaml:"cpu"`
	Ram      string `yaml:"ram"`
}

type Ports struct {
	Main       int   `yaml:"main"`
	Additional []int `yaml:"additional"`
}

type Server struct {
	Name     string `yaml:"name"`
	Image    string `yaml:"image"`
	Provider `yaml:"provider"`
	Ports    `yaml:"ports"`
	Params   map[string]interface{} `yaml:"params"`
}

// Once the droplet is created, this func connects to it via SSH and deploys the game server
// It returns an error if it is not deployed correctly
func CreateServer(client *godo.Client, ctx context.Context) error {
	// Make time for droplet to be deployed correctly
	time.Sleep(60 * time.Second)
	// Get the IPv4
	ip, err := GetIPv4(client, ctx, Servers.Name)
	if err != nil {
		return err
	}

	// Connect to droplet and exec server
	log.Println("Connecting via ssh to the droplet...")
	clientSSH, err := ConnectSSH(ip)
	if err != nil {
		for i := 0; i < 5; i++ {
			time.Sleep(30 * time.Second)

			fmt.Println("Do you want to try again? [yes/no]")
			reader := bufio.NewReader(os.Stdin)
			a, _ := reader.ReadString('\n')
			retry := strings.ToLower(strings.TrimSpace(a)) == "yes"

			if retry {
				log.Println("Retrying again...")

				clientSSH, err = ConnectSSH(ip)
				if err != nil {
					return err
				}
			} else {
				return errors.New(err.Error() + "Stop retrying connections by user")
			}
		}

	}

	defer clientSSH.Close() // Remember to close the connection
	log.Println("Deploying server in: " + ip)

	name := Servers.Name
	image := Servers.Image
	mainPort := Servers.Ports.Main
	additionalPorts := Servers.Ports.Additional
	params := Servers.Params

	env := make(map[string]string)

	env = iterateParams(env, params)

	initCommands(clientSSH, name, image, mainPort, additionalPorts, env)

	return nil
}

func initCommands(clientSSH *goph.Client, name string, image string, mainPort int, additionalPorts []int, env map[string]string) error {
	_, err := clientSSH.Run("docker pull " + Servers.Image)
	if err != nil {
		return err
	}

	_, err = clientSSH.Run(`mkdir -p /root/logs ;
							chmod 777 /root/logs`)
	if err != nil {
		return err
	}

	portString := " -p " + strconv.Itoa(Servers.Ports.Main) + ":" + strconv.Itoa(Servers.Ports.Main)
	if len(Servers.Ports.Additional) > 0 {
		for _, port := range Servers.Ports.Additional {
			portString += " -p " + strconv.Itoa(port) + ":" + strconv.Itoa(port)
		}
	}

	volString := " -v /root/logs/:/mnt/cuack/"

	var envString string
	for i, e := range env {
		envString += " -e \"" + strings.ToUpper(i) + "=" + e + "\""
	}

	com := "docker run -d --name " + Servers.Name + portString + volString + envString + " " + image
	log.Println(com)

	_, err = clientSSH.Run(com)
	if err != nil {
		return err
	}

	return nil
}

func iterateParams(env map[string]string, params map[string]interface{}) map[string]string {
	for i, param := range params {
		switch paramType := param.(type) {
		case string:
			env[i] = paramType
		case int:
			env[i] = strconv.Itoa(paramType)
		case map[string]interface{}:
			iterateParams(env, paramType)
		}
	}
	return env
}
