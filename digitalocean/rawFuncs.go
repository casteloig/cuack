package digitalocean

import (
	"context"

	"github.com/digitalocean/godo"
)

// Returns a slice of objects []godo.Region with all the regions and their attributes.
// It returns an error if request is not done correctly.
func rawGetRegions(client *godo.Client, ctx context.Context) ([]godo.Region, error) {
	opt := &godo.ListOptions{
		Page:    1,
		PerPage: 200,
	}

	regions, _, err := client.Regions.List(ctx, opt)
	if err != nil {
		return nil, err
	}
	return regions, nil
}

// Deletes a droplet based on his ID (int)
// It returns an error if request is not done correctly.
func rawDeleteDroplet(client *godo.Client, ctx context.Context, id int) error {
	_, err := client.Droplets.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

// Returns a slice of objects []godo.Droplet with a list of all droplets existing in your account.
// It returns an error if request is not done correctly.
func rawListDroplets(client *godo.Client, ctx context.Context) ([]godo.Droplet, error) {
	opt := &godo.ListOptions{
		Page:    1,
		PerPage: 200,
	}

	droplets, _, err := client.Droplets.List(ctx, opt)
	if err != nil {
		return nil, err
	}
	return droplets, nil
}

// Creates a new Droplet and returns an object *godo.Droplet. It also binds an ssh key (based on the ID, not the name)
//	to the droplet, so you can connect to it via ssh.
// It returns an error if request is not done correctly.
func rawCreateDropletWithSSH(client *godo.Client, ctx context.Context, name string, region string, size string, ssh int, game string) (*godo.Droplet, error) {
	request := &godo.DropletCreateRequest{
		Name:   name,
		Region: region,
		Size:   size,
		Image: godo.DropletCreateImage{
			Slug: "docker-20-04",
		},
		SSHKeys: []godo.DropletCreateSSHKey{
			{ID: ssh},
		},
		Tags: []string{"cuack", game},
	}

	droplet, _, err := client.Droplets.Create(ctx, request)
	if err != nil {
		return nil, err
	}
	return droplet, nil
}

// Creates a new Droplet and returns an object *godo.Droplet.
// It returns an error if request is not done correctly.
func rawCreateDroplet(client *godo.Client, ctx context.Context, name string, region string, size string, game string) (*godo.Droplet, error) {
	request := &godo.DropletCreateRequest{
		Name:   name,
		Region: region,
		Size:   size,
		Image: godo.DropletCreateImage{
			Slug: "docker-20-04",
		},
		Tags: []string{"cuack", game},
	}

	droplet, _, err := client.Droplets.Create(ctx, request)
	if err != nil {
		return nil, err
	}
	return droplet, nil
}

// Returns a slice of objects []godo.Key with all the ssh Keys existing in your account.
// It returns an error if request is not done correctly.
func rawListSSH(client *godo.Client, ctx context.Context) ([]godo.Key, error) {
	opt := &godo.ListOptions{
		Page:    1,
		PerPage: 200,
	}

	keys, _, err := client.Keys.List(ctx, opt)
	if err != nil {
		return nil, err
	}
	return keys, nil
}
