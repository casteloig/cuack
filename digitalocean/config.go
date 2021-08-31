package digitalocean

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/digitalocean/godo"
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

type Server struct {
	Name     string `yaml:"name"`
	Image    string `yaml:"image"`
	Provider `yaml:"provider"`
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
	fmt.Println("Connecting via ssh to the droplet...")
	clientSSH, err := ConnectSSH(ip)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Retrying again...")
		time.Sleep(30 * time.Second)
		clientSSH, err = ConnectSSH(ip)
		if err != nil {
			return err
		}
	}

	defer clientSSH.Close() // Remember to close the connection
	fmt.Println("Deploying server in: " + ip)

	_, err = clientSSH.Run("docker pull " + Servers.Image)
	if err != nil {
		return err
	}

	err = clientSSH.Upload("setup_1.sh", "/root/setup_1.sh")
	if err != nil {
		return err
	}
	_, err = clientSSH.Run("chmod +x /root/setup_1.sh && bash /root/setup_1.sh")
	if err != nil {
		return err
	}

	err = clientSSH.Upload("setup_2.py", "/root/setup_2.py")
	if err != nil {
		return err
	}
	execPythonEnv := "python3 setup_2.py"
	execPythonEnv += " name=" + Servers.Name + " image=" + Servers.Image
	for i, env := range Servers.Params {
		execPythonEnv += " "
		switch e := env.(type) {
		case string:
			execPythonEnv += i + "=" + e
		case int:
			execPythonEnv += i + "=" + strconv.Itoa(e)
		case []interface{}:
			execPythonEnv += i + "="
			for _, param := range e {
				execPythonEnv += ","
				execPythonEnv += strconv.Itoa(param.(int))

			}
		default:
			return errors.New("param type not supported")
		}
	}
	fmt.Println(execPythonEnv)
	_, err = clientSSH.Run(execPythonEnv)
	if err != nil {
		return err
	}

	return nil
}
