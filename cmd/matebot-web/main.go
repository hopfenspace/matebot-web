package main

import (
	"github.com/hellflame/argparse"
	"github.com/hopfenspace/matebot-web/server"
	"github.com/myOmikron/echotools/color"
)

func main() {
	parser := argparse.NewParser("matebot-web", "", &argparse.ParserConfig{DisableDefaultShowHelp: true})

	configPath := parser.String("", "config-path", &argparse.Option{
		Default: "/etc/matebot-web/config.toml",
		Help:    "Specify an alternative path to the configuration file.",
	})

	if err := parser.Parse(nil); err != nil {
		color.Println(color.RED, err.Error())
	}

	server.StartServer(*configPath)
}
