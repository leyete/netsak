package models

import (
	//"strings"
	"errors"
	"github.com/mattrbeam/netsak/nsconfig/validator"
)

type ServerConfiguration struct {
	//Type    string `yaml:"type"`
	Url     string `yaml:"url"`
	Port    string `yaml:"port"`
	//Swagger string `yaml:"swagger"`
	User string `yaml:"user"`
	Password string `yaml:"password"`
}

func (server *ServerConfiguration) ServerValidator() error {
	if server.Url == "" {
		return errors.New("Url undefined")
	}

	if server.Port == "" {
		return errors.New("Server port undefined")
	}

	url := validator.IsURL(server.Url)
	if url == false {
		return errors.New("Url not valid")
	}

	if server.User == "" {
		return errors.New("User undefined")
	}

	if server.Password == "" {
		return errors.New("Password undefined")
	}

	return nil

}
