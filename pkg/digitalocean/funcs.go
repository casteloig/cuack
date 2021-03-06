package digitalocean

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strings"
	"time"

	"github.com/digitalocean/godo"
	"github.com/docker/docker/pkg/namesgenerator"
	"github.com/melbahja/goph"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

// Connects to a droplet via SSH and returns an object *goph.Client which
//	you can use to exec commands. (ip is the IPv4 as a string, Ex: "0.0.0.0")
// It returns an error if is not established correctly
func ConnectSSH(ip string) (*goph.Client, error) {
	// gets private ssh key
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	dir := home + "/.ssh/id_rsa"
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Type password of ssh key:")
	pass, err := reader.ReadString('\n')
	pass = strings.Trim(pass, "\n")
	if err != nil {
		return nil, err
	}

	// gets an auth method goph.Auth for handling the connection request
	auth, err := goph.Key(dir, pass)
	if err != nil {
		return nil, err
	}

	// asks for a new ssh connection returning the client for SSH
	client, err := goph.NewConn(&goph.Config{
		User:     "root",
		Addr:     ip,
		Port:     22,
		Auth:     auth,
		Callback: VerifyHost, //HostCallBack custom (appends host to known_host if not exists)
	})
	if err != nil {
		return nil, err
	}
	return client, nil
}

// Func HostCallBack custom
func VerifyHost(host string, remote net.Addr, key ssh.PublicKey) error {

	hostFound, err := goph.CheckKnownHost(host, remote, key, "")
	// Host in known hosts but key mismatch!
	// Maybe because of MAN IN THE MIDDLE ATTACK!
	if hostFound && err != nil {
		return err
	}
	// handshake because public key already exists.
	if hostFound && err == nil {
		return nil
	}
	// Ask user to check if he trust the host public key.
	if !askIsHostTrusted(host, key) {
		// Make sure to return error on non trusted keys.
		return errors.New("you typed no, aborted")
	}
	// Add the new host to known hosts file.
	return goph.AddKnownHost(host, remote, key, "")
}

// Support func for HostCallBack function custom
func askIsHostTrusted(host string, key ssh.PublicKey) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Would you like to add host to 'known_hosts'? [yes/no] ")
	a, err := reader.ReadString('\n')
	if err != nil {
		logrus.WithFields(logrus.Fields{}).Info(err)
	}
	return strings.ToLower(strings.TrimSpace(a)) == "yes"
}

// Stores in global var (Region string) the region obtained by cuack.config.
// It returns an error if region is not obtained correctly.
func GetRegionFromFile() error {
	home, _ := os.UserHomeDir()
	dir := home + "/.config/cuack.config"
	file, err := os.Open(dir)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		Region = scanner.Text()
		if strings.HasPrefix(Region, "region") {
			Region = strings.Fields(Region)[1]
			return nil
		}
	}
	return errors.New("error getting the region from file")
}

// Stores in global var (Token string) the region obtained by cuack.config.
// It returns an error if token is not obtained correctly.
func GetTokenFromFile() error {
	home, _ := os.UserHomeDir()
	dir := home + "/.config/cuack.config"
	file, err := os.Open(dir)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		Token = scanner.Text()
		if strings.HasPrefix(Token, "key") && strings.Fields(Token)[1] != "" {
			Token = strings.Fields(Token)[1]
			return nil
		}
	}
	return errors.New("error getting the token from file")
}

// Stores in global var (Token string) the region obtained by cuack.config.
// It returns an error if token is not obtained correctly.
func DeleteDropletByName(client *godo.Client, ctx context.Context, name string) error {
	droplets, err := rawListDroplets(client, ctx)
	if err != nil {
		return errors.New("error getting list of existing droplets")
	}

	for _, droplet := range droplets {
		if droplet.Name == name {
			err := rawDeleteDroplet(client, ctx, droplet.ID)
			if err != nil {
				return errors.New("error trying to delete the droplet")
			}
			return nil
		}
	}
	return errors.New("droplet does not exist with that name")
}

// Creates a new Droplet and returns an object *godo.Droplet. It also binds an ssh key (based on the name string)
//	to the droplet, so you can connect to it via ssh.
// It returns an error if request is not done correctly.
func CreateDropletWithSSH(client *godo.Client, ctx context.Context, size string) (*godo.Droplet, error) {
	keys, err := rawListSSH(client, ctx)
	if err != nil {
		return nil, errors.New("error getting existing ssh keys")
	}

	newName, err := generateName(client, ctx, Servers.Name)
	if err != nil {
		return nil, errors.New("could not generate new name")
	}
	Servers.Name = newName

	for _, key := range keys {
		if key.Name == Servers.Provider.SshName {
			droplet, err := rawCreateDropletWithSSH(client, ctx, Servers.Name, Region, size, key.ID)
			if err != nil {
				return nil, errors.New("error creating droplet")
			}
			logrus.WithFields(logrus.Fields{
				"command":    "create",
				"final-name": Servers.Name,
			}).Info("Droplet created")
			return droplet, nil
		}
	}
	return nil, errors.New("ssh key does not exist with that name")
}

// Returns a slice of objects []godo.Region with all the AVAILABLE regions and their attributes.
// It returns an error if request is not done correctly.
func GetAvailableRegions(client *godo.Client, ctx context.Context) ([]godo.Region, error) {
	regions, err := rawGetRegions(client, ctx)
	if err != nil {
		return nil, err
	}

	var listRegions []godo.Region
	for _, region := range regions {
		if region.Available {
			listRegions = append(listRegions, region)
		}
	}
	return listRegions, nil
}

// Returns a string with the IPv4 of the server.
// It returns an error if it is not archieved correctly.
func GetIPv4(client *godo.Client, ctx context.Context, name string) (string, error) {
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
			return ip, nil
		}
	}
	return "", errors.New("droplet does not exist with that name")
}

// Returns the maximun number of droplets an account can have.
func GetMaxDroplets(client *godo.Client, ctx context.Context) int {
	acc, _, _ := client.Account.Get(ctx)
	return acc.DropletLimit
}

// Returns the number of droplets exist in an account.
func GetNumberDroplets(client *godo.Client, ctx context.Context) int {
	lo := godo.ListOptions{
		Page:    1,
		PerPage: 20,
	}
	// returns []Droplets
	list, _, _ := client.Droplets.List(ctx, &lo)

	return len(list)
}

func CheckDropletExists(client *godo.Client, ctx context.Context, name string) (bool, error) {
	droplets, err := rawListDroplets(client, ctx)
	if err != nil {
		return false, err
	}

	for _, droplet := range droplets {
		if droplet.Name == name {
			return true, nil
		}
	}
	return false, nil
}

func ListCuackDroplets(client *godo.Client, ctx context.Context) (map[string]string, error) {
	opt := &godo.ListOptions{
		Page:    1,
		PerPage: 200,
	}

	droplets, _, err := client.Droplets.List(ctx, opt)
	if err != nil {
		return nil, errors.New("error listting cuack droplets")
	}

	list := make(map[string]string)
	for _, droplet := range droplets {
		if contains(droplet.Tags, "cuack") {
			ip, err := GetIPv4(client, ctx, droplet.Name)
			if err != nil {
				return nil, errors.New("could not get the IPv4 from the server")
			}
			list[ip] = droplet.Name
		}
	}

	return list, nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func generateName(client *godo.Client, ctx context.Context, oldName string) (string, error) {
	rand.Seed(time.Now().UnixNano())

	suffixes := strings.ReplaceAll(namesgenerator.GetRandomName(0), "_", "-")
	newName := oldName + "-" + suffixes

	exists, err := CheckDropletExists(client, ctx, newName)
	if err != nil {
		return "", err
	}
	if exists {
		newName, err = generateName(client, ctx, oldName)
		if err != nil {
			return "", err
		}
	}

	return newName, nil
}
