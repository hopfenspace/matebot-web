package conf

import "github.com/hopfenspace/MateBotSDKGo"

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
	return nil
}
