package digitalocean

var Servers []ServerGeneral
var Token string
var Region string

type Provider struct {
	NameProv string `yaml:"name-prov"`
	SshName  string `yaml:"ssh-name"`
	Cpu      int    `yaml:"cpu"`
	Ram      string `yaml:"ram"`
}

type Game struct {
	NameGame   string `yaml:"name-game"`
	Image      string `yaml:"image"`
	WorldName  string `yaml:"world-name"`
	Players    int    `yaml:"players"`
	Difficulty string `yaml:"difficulty"`
}

type Server struct {
	Name     string `yaml:"name"`
	Provider `yaml:"provider"`
	Game     `yaml:"game"`
}

type ServerGeneral struct {
	Sv Server `yaml:"server"`
}
