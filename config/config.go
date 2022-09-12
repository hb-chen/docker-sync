package config

type Config struct {
	Pull        bool     `yaml:"pull"`
	PullOptions []string `yaml:"pull_options"`

	Push bool `yaml:"push"`

	Sync []*Sync `yaml:"sync"`
}

type Sync struct {
	Src string `yaml:"src"`
	Dst string `yaml:"dst"`
}
