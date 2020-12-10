package config

type Config struct {
	Sync []*Sync
}

type Sync struct {
	Src string
	Dst string
}
