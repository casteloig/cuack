<h1 align="center">
  Cuack
</h1>
<p align="center">Deploy automatically your gameservers on cloud providers (only DigitalOcean is supported so far) in a easy and fast way</p>

[![Docker Downloads](https://img.shields.io/docker/pulls/casteloig/cuack?style=flat-square)](https://hub.docker.com/repository/docker/casteloig/cuack)

## Quick start
1. Create the **initialization file** and introduce the Personal Access Token (DigitalOcean's can be found [here](https://docs.digitalocean.com/reference/api/create-personal-access-token/)) and the prefered region slug (these are displayed on screen.)

```bash
cuack-ctl init
```

2. **Create your gameserver** with default options (minecraft and factorio supported so far)

```bash
# Create a default minecraft server
cuack-ctl create -f https://github.com/casteloig/cuack/tree/main/games/minecraft/default.yaml

# Create a default factorio server
cuack-ctl create -f https://github.com/casteloig/cuack/tree/main/games/factorio/default.yaml
```
<!-- 
**List** all of your _cuack servers_:

```bash
cuack-ctl list
```

**Delete** a running server:

```bash
# Delete a DigitalOcean droplet named "minecraft-foo-bar"
cuack-ctl delete minecraft-foo-bar
``` -->
## Commands and options
### `init`
The **init** command is used to create the init file with the basic configuration of Cuack. It is stored by default on `~/.config/cuack.config`.

This file is composed by two parts:

| Parts  | Description                                              | Type   | Default |
| ------ | -------------------------------------------------------- | ------ | ------- |
| token   | Where the server provider's personal token is stored (**it should never be shared!**). | `string` |         |
| region   | It is the slug that represents the prefered region where all droplets will be created on. | `string` | lon1 |

> This file **must** exist and contain all the information above to perform the rest of the commands.

### `create`
The **create** command is used to create a droplet and, then, create a gameserver in that droplet. It does not have any arguments, but it has two flags:

| Flags      | Flags (short) | Description                                              | Type   | Required | Default |
| ---------- | ------------- |--------------------------------------------------------- | ------ | -------- | ------- |
| `--file`   | `-f` | It is the configuration file of the droplet and server | `string` | Yes |  |
| `--select` | `-s` | Select one of the existing configurations in the config file passed on `-f` | `string` | No | Selects always the first configuration from the file |

A configuration file is a `.yaml`/`.yml` file that can either a local file or a remote file (with a URL that starts with _https://_ or _http://_). One example of the file would be this one (as an example of minecraft file):

```yaml
name: minecraft
image: casteloig/mine-server:latest
provider:
  name_prov: digitalocean
  ssh_name: asus-laptop-linux
  cpu: 2
  ram: 2GB
ports:
  main: 25565
  additional: 
    - 25575
params:   
  world_name: world # not used yet
  players: 5
  difficulty: easy
---
name: minecraft-big
image: casteloig/mine-server:latest
provider:
  name_prov: digitalocean
  ssh_name: asus-laptop-linux
  cpu: 2
  ram: 2GB
ports:
  main: 25565
  additional: 
    - 25575
params:   
  world_name: world # not used yet
  players: 10
  difficulty: medium
```

As you can see, there are two different configurations, you can select which one you want by passing the `name` of the configuration on the `-s` flag (in this case `minecraft` or `minecraft-big`).

```bash
# If we want to create a droplet with the 'minecraft-big' configuration we would do
cuack-ctl -f <path_to_yaml_file> -s minecraft-big
```

By default, in this repository, exist some basic configs you can use directly or download and modify to your own belong.

Every config file must be composed by two parts:

#### First part

This part must be **static**, that means that the variable names and indentation cannot be modified, only the values.

| Variables  | Description                                              |
| ------ | -------------------------------------------------------- |
| name | The name that will use the `-s` flag and also will be used to create the droplet name |
| image | The Docker image that will be used to create the server in the droplet. [(more info about how it works)](#how-cuack-uses-the-docker-images-to-start-his-servers) |
| provider | All the information that the server provider needs to create the droplet. [(see the sub-table here)](#provider-sub-table) |
| ports | The ports that the server needs (so docker may use). There must be always one main port, and may be any number of additional ports. |


##### Provider sub-table
| Variables  | Description                                              | Type   |
| ------ | -------------------------------------------------------- | ------ |
| name_prov | So far we only work with DigitalOcean | `string` |
| ssh_name | The name of the ssh key that you must have created in your server provider | `string` |
| cpu | Number of CPUs that the droplet will have | `integer` |
| ram | RAM size that the droplet will have | `string` composed by '_number of GB_' + 'GB' |


#### Second part

This second part it is represented by the **upper node called `params`**, and it can be dynamic. Why...? All this second part represent environment variables that Cuack will pass to the docker container. So all Cuack is going to do is to read all the final nodes and treat them as strings that will pass as environment variables. That's why you can create the yaml structure that you prefer as long as every final node represents a environment variable.

> Keep in mind that all values will be strings when passed as environment variables



### `list`
The **list** command is used to list all droplets created on the server providers by cuack.

It shows the droplet's name and the IP where it is deployed on.

```bash
cuack-ctl list
```

### `inspect`
The **inspect** command is used to see the configuration file (the yaml file) and some other details of the droplet, like CPU and RAM usage.

You must know the name of the droplet, which can be obtained via the `list` command.

```bash
cuack-ctl inspect <name_of_droplet>
```

It prints in command line the URL where is served a tiny website where it's displayed all this information. By detault it's `localhost:8080/`

### `delete`
The **delete** command is used to delete a server and the droplet that it is deployed on.

It only accepts one argument:

| Argument | Description                                              | Type   |
| -------- | -------------------------------------------------------- | ------ |
| `name_of_droplet` | It represents the name of the droplet | `string` |

```bash
cuack-ctl delete <name_of_droplet>
```

Keep in mind that **the name of the droplet is not the same as the name represented in the yaml file** when you created the server. That's because cuack will always take the name in the yaml file and add some extra strings to distinguish two servers created with the same file. So you should check the name with the `list` command before.


## How Cuack uses docker images to start his servers?

Cuack has his own images to deploy his servers because cuack abstracts himself from the way the images are created.

We explained a little bit of this [here](#second-part).

In a nutshell, Cuack just creates the droplet where the server will be deployed and runs the container within the droplet passing all the params as environment variables. This provides a lot of freedom to the way Cuack can be used, because you actually can build your own image of another game/service and run inside of it the application that you want as long as you keep in mind that all the params of the `yaml` file will be passed as environment variables to the running image.