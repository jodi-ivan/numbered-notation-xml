package config

import (
	"fmt"

	"gopkg.in/gcfg.v1"
)

type Config struct {
	Webserver WebServerConfig
}

type WebServerConfig struct {
	Port string
}

func InitConfig(env string) (Config, error) {
	result := Config{}

	path := "/etc/numbered-mutation-xml/"
	if env == "development" {
		path = "files/etc/numbered-mutation-xml/"
	}

	filenames := []string{
		"config",
	}

	var fullpath string
	for _, filename := range filenames {
		fullpath = fmt.Sprintf("%s%s.ini", path, filename)

		err := gcfg.ReadFileInto(&result, fullpath)
		if err != nil {
			return result, err
		}
	}
	return result, nil
}
