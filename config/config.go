package config

type Config struct {
	Pull bool
	Push bool

	Sync []*Sync
}

type Sync struct {
	Src string
	Dst string
}
