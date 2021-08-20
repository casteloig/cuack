package main

import (
	"context"

	"github.com/digitalocean/godo"
)

func rawDeleteDroplet(client *godo.Client, ctx context.Context, id int) error {
	_, err := client.Droplets.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

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

func rawCreateDropletWithSSH(client *godo.Client, ctx context.Context, name string, region string, size string, ssh int) (*godo.Droplet, error) {
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
	}

	droplet, _, err := client.Droplets.Create(ctx, request)
	if err != nil {
		return nil, err
	}
	return droplet, nil
}

func rawCreateDroplet(client *godo.Client, ctx context.Context, name string, region string, size string) (*godo.Droplet, error) {
	request := &godo.DropletCreateRequest{
		Name:   name,
		Region: region,
		Size:   size,
		Image: godo.DropletCreateImage{
			Slug: "docker-20-04",
		},
	}

	droplet, _, err := client.Droplets.Create(ctx, request)
	if err != nil {
		return nil, err
	}
	return droplet, nil
}

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
