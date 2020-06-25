//go:generate wire

//+build wireinject

package main

import (
	"dt/config"
	"dt/controller"
	"dt/handlers"
	"dt/logwrap"
	"dt/managers"
	"dt/rpc"
	"dt/stores"
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	_ "github.com/lib/pq"
	"strconv"
	"time"
)

func main() {
	conf, err := ParseConfig()
	if err != nil {
		logwrap.InitializePanic("conf err: " + err.Error())
	}

	r, cleanUp, err := Injector(conf)
	if err != nil {
		logwrap.InitializePanic("injection err: " + err.Error())
	}

	go func() {
		time.Sleep(time.Second * 7)
		logwrap.Debug("service started")
	}()

	r.Run(conf.Host + ":" + strconv.Itoa(int(conf.Port)))
	cleanUp()
}

func ParseConfig() (*config.Config, error) {
	configPath := flag.String("config", "", "path to config json file")
	flag.Parse()

	return config.ParseConfig(config.PathConfig(*configPath))
}

func Injector(c *config.Config) (*gin.Engine, func(), error) {
	wire.Build(
		handlers.InitHandlers,
		managers.ProviderSet,
		initer.Claserize("172.68.52.57"),
		stores.InitDB,
		rpc.ProviderSet,
		controller.NewServer,
	)
	return nil, nil, nil
}
