package apiserver

type Config struct {
	Addr        string `toml:"addr"`
	DatabaseURL string `toml:"dburl"`
	SessionKey  string `toml:"session_key"`
}

func NewConfig() *Config {
	return &Config{
		Addr: ":8080",
	}
}
