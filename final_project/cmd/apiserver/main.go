package main

import (
	"final_project/internal/app/apiserver"
	"flag"

	"github.com/BurntSushi/toml"
)

var (
	confPath string
)

func init() {
	flag.StringVar(&confPath, "configpath", "configs/apiserver.toml", "path to config")
}

func main() {
	flag.Parse()

	config := apiserver.NewConfig()
	_, err := toml.DecodeFile(confPath, config)
	if err != nil {
		panic(err)
	}

	if err = apiserver.Start(config); err != nil {
		panic(err)
	}
}
