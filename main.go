package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path"
	"runtime"
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ykgk/kgo/config"
	"gitlab.com/ykgk/kgo/controller"
	"gitlab.com/ykgk/kgo/db"
)

const (
	APPName = "kgo"
)

func init() {
	if err := config.Init(APPName); err != nil {
		panic(err)
	}

	initLog()
	initDatabase()
}

func initLog() {
	log.SetReportCaller(true)
	log.SetFormatter(&log.JSONFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)
			return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
		},
	})
}

func initDatabase() {
	type PSQLConfig struct {
		URL   string `json:"url" yaml:"url"`
		Debug bool   `json:"debug" yaml:"debug"`
	}

	// type MongoConfig struct {
	// 	URL string `json:"url" yaml:"url"`
	// 	DB  string `json:"db" yaml:"db"`
	// }

	psqlCfg := PSQLConfig{}
	if err := config.Unmarshal("db.psql", &psqlCfg); err != nil {
		panic(err)
	}

	if err := db.InitPG(psqlCfg.URL, psqlCfg.Debug); err != nil {
		panic(err)
	}

	// mongoCfg := MongoConfig{}
	// if err := config.Unmarshal("db/mongo", &mongoCfg); err != nil {
	// 	panic(err)
	// }

	// if err := db.InitMongo(mongoCfg.URL, mongoCfg.DB); err != nil {
	// 	panic(err)
	// }
}

func main() {
	type HTTPConfig struct {
		Listen  string `json:"listen" yaml:"listen"`
		Session string `json:"session" yaml:"session"`
	}

	httpCfg := HTTPConfig{}
	if err := config.Unmarshal("http", &httpCfg); err != nil {
		panic(err)
	}
	var listen string
	flag.StringVar(&listen, "http.listen", "", "bind address")
	flag.Parse()
	if listen != "" {
		httpCfg.Listen = listen
	}

	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Output: os.Stdout,
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderXRequestedWith},
		AllowCredentials: true,
		AllowMethods:     []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(httpCfg.Session))))
	e.HTTPErrorHandler = errorHandler

	controller.Register(e)

	instanceName := APPName + httpCfg.Listen

	go func() {
		if err := e.Start(httpCfg.Listen); err != nil {
			log.Fatalf("shutting down the server: %v", err)
		}
	}()
	log.Infof("%s started", instanceName)

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, os.Kill)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
	log.Warnf("%s quited", instanceName)
}

func errorHandler(err error, c echo.Context) {
	c.NoContent(500)
}
