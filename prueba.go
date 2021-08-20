package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/digitalocean/godo"
	"github.com/melbahja/goph"
)

func main() {
	client := godo.NewFromToken("46845e1df2c97b20e4ac1fe3ddab145e53c699168cd9d150cad9d40322eb4a27")
	ctx := context.TODO()

	// fmt.Println("Creating droplet")
	// _, err := createDropletWithSSH(client, ctx, "minecraft", "lon1", "s-1vcpu-1gb", "msi-laptop-linux")
	// if err != nil {
	// 	fmt.Println("Problema creando droplet")
	// }

	// err := deleteDropletByName(client, ctx, "test1")
	// if err != nil {
	// 	fmt.Println("cagada en delete")
	// }

	// time.Sleep(20 * time.Second)

	droplets, err := rawListDroplets(client, ctx)
	if err != nil {
		fmt.Println("cagada en list")
	}
	for _, droplet := range droplets {
		fmt.Println(droplet.ID, droplet.Name)
	}

	err = connectSSH(client, ctx, "minecraft")
	if err != nil {
		fmt.Println(err)
	}

}

func connectSSH(client *godo.Client, ctx context.Context, name string) error {
	home, _ := os.UserHomeDir()
	dir := home + "/.ssh/id_rsa"
	auth, err := goph.Key(dir, "2311")
	if err != nil {
		return err
	}

	ip, err := getIPv4(client, ctx, name)
	if err != nil {
		return err
	}

	sshClient, err := goph.New("root", ip, auth)
	if err != nil {
		return err
	}

	defer sshClient.Close()

	out, err := sshClient.Run("docker pull casteloig/mine-server:latest")
	if err != nil {
		return err
	}

	fmt.Println(string(out))
	return nil
}

func deleteDropletByName(client *godo.Client, ctx context.Context, name string) error {
	droplets, err := rawListDroplets(client, ctx)
	if err != nil {
		return err
	}

	for _, droplet := range droplets {
		if droplet.Name == name {
			err := rawDeleteDroplet(client, ctx, droplet.ID)
			if err != nil {
				return err
			}
		}
	}
	return errors.New("droplet does not exist with that name")
}

func createDropletWithSSH(client *godo.Client, ctx context.Context, name string, region string, size string, sshName string) (*godo.Droplet, error) {
	keys, err := rawListSSH(client, ctx)
	if err != nil {
		return nil, err
	}

	for _, key := range keys {
		if key.Name == sshName {
			droplet, err := rawCreateDropletWithSSH(client, ctx, name, region, size, key.ID)
			if err != nil {
				return nil, err
			}
			return droplet, nil
		}
	}
	return nil, errors.New("ssh key does not exist with that name")
}

func getIPv4(client *godo.Client, ctx context.Context, name string) (string, error) {
	droplets, err := rawListDroplets(client, ctx)
	if err != nil {
		return "", err
	}

	for _, droplet := range droplets {
		if droplet.Name == name {
			ip, err := droplet.PublicIPv4()
			if err != nil {
				return "", err
			}
			fmt.Println(ip)
			return ip, nil
		}
	}
	return "", errors.New("droplet does not exist with that name")
}

//////// RAWS

// getMaxDroplets returns the maximun number of droplets an account can have.
func getMaxDroplets(client *godo.Client, ctx context.Context) int {
	acc, _, _ := client.Account.Get(ctx)
	return acc.DropletLimit
}

// getNumberDroplets returns the number of droplets exist in an account.
func getNumberDroplets(client *godo.Client, ctx context.Context, max int) int {
	lo := godo.ListOptions{
		Page:    1,
		PerPage: max,
	}
	// returns []Droplets
	list, _, _ := client.Droplets.List(ctx, &lo)

	return len(list)
}

func checkDropletExists(client *godo.Client, ctx context.Context, name string) (bool, error) {
	droplets, err := rawListDroplets(client, ctx)
	if err != nil {
		return false, err
	}

	for _, droplet := range droplets {
		if droplet.Name == name {
			return true, nil
		}
	}
	return false, errors.New("droplet does not exist with that name")
}
