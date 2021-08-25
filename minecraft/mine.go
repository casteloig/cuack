package minecraft

import (
	"context"
	do "cuack/digitalocean"
	"fmt"
	"strconv"
	"time"

	"github.com/digitalocean/godo"
	"github.com/melbahja/goph"
)

// Once the droplet is created, this func connects to it via SSH and deploys the game server
// It returns an error if it is not deployed correctly
func CreateServer(client *godo.Client, ctx context.Context, server do.ServerGeneral) error {
	// Make time for droplet to be deployed correctly
	time.Sleep(60 * time.Second)
	// Get the IPv4
	ip, err := do.GetIPv4(client, ctx, server.Sv.Name)
	if err != nil {
		return err
	}
	// Connect to droplet and exec server
	fmt.Println("Connecting via ssh to the droplet...")
	clientSSH, err := do.ConnectSSH(ip)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Retrying again...")
		time.Sleep(30 * time.Second)
		clientSSH, err = do.ConnectSSH(ip)
		if err != nil {
			return err
		}
	}

	defer clientSSH.Close() // Remember to close the connection
	fmt.Println("Deploying server in: " + ip)
	err = commandsInitSSH(clientSSH, server)
	if err != nil {
		return err
	}
	return nil
}

// Uploads every file needed and makes the changes to the game server config keeping in mind the
//	config specified in the yaml file. It also starts the server.
// It returns an error if the server is not deployed correctly with that conf
func commandsInitSSH(clientSSH *goph.Client, server do.ServerGeneral) error {
	// Upload all
	err := clientSSH.Upload("images/"+server.Sv.Image+"/eula.txt", "/root/eula.txt")
	if err != nil {
		return err
	}
	err = clientSSH.Upload("images/"+server.Sv.Image+"/Dockerfile", "/root/Dockerfile")
	if err != nil {
		return err
	}
	err = clientSSH.Upload("images/"+server.Sv.Image+"/server.properties", "/root/server.properties")
	if err != nil {
		return err
	}
	err = clientSSH.Upload("images/"+server.Sv.Image+"/start.sh", "/root/start.sh")
	if err != nil {
		return err
	}

	_, err = clientSSH.Run("sed -i 's/^max-players=.*/max-players=" + strconv.Itoa(server.Sv.Players) + "/' server.properties")
	if err != nil {
		return err
	}
	_, err = clientSSH.Run("sed -i 's/^difficulty=.*/difficulty=" + (server.Sv.Difficulty) + "/' server.properties")
	if err != nil {
		return err
	}

	_, err = clientSSH.Run("docker build -t mine-server:latest /root/")
	if err != nil {
		return err
	}

	_, err = clientSSH.Run("docker volume create mine")
	if err != nil {
		return err
	}

	_, err = clientSSH.Run("docker run -d --name minecraft --mount source=mine,destination=/root/ -p 25565:25565 mine-server:latest ")
	if err != nil {
		return err
	}

	return nil
}
