package conf

type Database struct {
	Driver   string
	Host     string
	Port     uint16
	Name     string
	User     string
	Password string
}

type Consumable struct {
	ID   uint
	Name string
}

type MateBot struct {
	Url         string
	User        string
	Password    string
	ID          uint
	Consumables []Consumable
}

type Generic struct {
	Listen      string
	TemplateDir string
	StaticDir   string
}

type Config struct {
	Generic  Generic
	MateBot  MateBot
	Database Database
}

func (c *Config) CheckConfig() error {
	return nil
}
