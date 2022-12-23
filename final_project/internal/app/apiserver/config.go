package apiserver

type Config struct {
	Addr        string `toml:"addr"`
	DatabaseURL string `toml:"dburl"`
}

func NewConfig() *Config {
	return &Config{
		Addr: ":8080",
	}
}
