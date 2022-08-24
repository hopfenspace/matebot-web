package server

import (
	"errors"
	"fmt"
	mysqlDriver "github.com/go-sql-driver/mysql"
	sdk "github.com/hopfenspace/MateBotSDKGo"
	"github.com/hopfenspace/matebot-web/conf"
	"github.com/hopfenspace/matebot-web/models"
	"github.com/labstack/echo/v4"
	mw "github.com/labstack/echo/v4/middleware"
	"github.com/myOmikron/echotools/color"
	"github.com/myOmikron/echotools/database"
	"github.com/myOmikron/echotools/execution"
	"github.com/myOmikron/echotools/middleware"
	"github.com/myOmikron/echotools/utilitymodels"
	"github.com/myOmikron/echotools/worker"
	"github.com/pelletier/go-toml"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"html/template"
	"io/fs"
	"io/ioutil"
	"net"
	"net/url"
	"os"
	"path"
	"strconv"
	"time"
)

func StartServer(configPath string) {
	config := &conf.Config{}

	if configBytes, err := ioutil.ReadFile(configPath); errors.Is(err, fs.ErrNotExist) {
		color.Printf(color.RED, "Config was not found at %s\n", configPath)
		b, _ := toml.Marshal(config)
		fmt.Print(string(b))
		os.Exit(1)
	} else {
		if err := toml.Unmarshal(configBytes, config); err != nil {
			panic(err)
		}
	}

	// Check for valid config values
	if err := config.CheckConfig(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Database
	var driver gorm.Dialector
	switch config.Database.Driver {
	case "sqlite":
		driver = sqlite.Open(config.Database.Name)
	case "mysql":
		mysqlConf := mysqlDriver.NewConfig()
		mysqlConf.Net = fmt.Sprintf("tcp(%s)", net.JoinHostPort(config.Database.Host, strconv.Itoa(int(config.Database.Port))))
		mysqlConf.DBName = config.Database.Name
		mysqlConf.User = config.Database.User
		mysqlConf.Passwd = config.Database.Password
		mysqlConf.ParseTime = true
		mysqlConf.Params = map[string]string{
			"charset": "utf8mb4",
		}
		driver = mysql.Open(mysqlConf.FormatDSN())
	case "postgresql":
		dsn := url.URL{
			Scheme: "postgres",
			User:   url.UserPassword(config.Database.User, config.Database.Password),
			Host:   net.JoinHostPort(config.Database.Host, strconv.Itoa(int(config.Database.Port))),
			Path:   config.Database.Name,
		}
		driver = postgres.Open(dsn.String())
	}

	db := database.Initialize(
		driver,
		&utilitymodels.Session{},
		&utilitymodels.LocalUser{},
		&models.CoreUser{},
	)

	// Web server
	e := echo.New()
	e.HideBanner = true

	// SDK client
	client, err := sdk.New(&config.MateBot)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Worker pool
	wp := worker.NewPool(&worker.PoolConfig{
		NumWorker: 10,
		QueueSize: 100,
	})
	wp.Start()

	// Template rendering
	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob(path.Join(config.Server.TemplateDir, "*.gohtml"))),
	}
	e.Renderer = renderer

	// Middleware
	e.Use(mw.Logger())
	e.Use(mw.Recover())

	duration := time.Hour * 24
	e.Use(middleware.Session(db, &middleware.SessionConfig{
		CookieName: "session_id",
		CookieAge:  &duration,
	}))
	middleware.RegisterAuthProvider(utilitymodels.GetLocalUser(db))

	allowedHosts := make([]middleware.AllowedHost, 0)
	for _, host := range config.Server.AllowedHosts {
		allowedHosts = append(allowedHosts, middleware.AllowedHost{
			Host:  host.Host,
			Https: host.Https,
		})
	}
	secConfig := &middleware.SecurityConfig{
		AllowedHosts:            allowedHosts,
		UseForwardedProtoHeader: config.Server.UseForwardedProtoHeader,
	}
	e.Use(middleware.Security(secConfig))

	// Router
	defineRoutes(e, db, config, client, wp)

	execution.SignalStart(e, config.Server.Listen, &execution.Config{
		ReloadFunc: func() {
			StartServer(configPath)
		},
		StopFunc: func() {
			if err := client.Shutdown(); err != nil {
				e.Logger.Error(err)
			}
		},
		TerminateFunc: func() {
			if err := client.Shutdown(); err != nil {
				e.Logger.Error(err)
			}
		},
	})
}
