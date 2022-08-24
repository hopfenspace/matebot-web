package conf

import (
	"errors"
	"github.com/hopfenspace/MateBotSDKGo"
)

type Database struct {
	Driver   string
	Host     string
	Port     uint16
	Name     string
	User     string
	Password string
}

type AllowedHost struct {
	Host  string
	Https bool
}

type Server struct {
	Listen                  string
	TemplateDir             string
	StaticDir               string
	AllowedHosts            []AllowedHost
	UseForwardedProtoHeader bool
}

type Config struct {
	Server   Server
	MateBot  MateBotSDKGo.Config
	Database Database
}

func (c *Config) CheckConfig() error {
	if c.MateBot.CallbackSecret == nil || *c.MateBot.CallbackSecret == "" {
		return errors.New("callback secret for MateBot must be set")
	}
	return nil
}
