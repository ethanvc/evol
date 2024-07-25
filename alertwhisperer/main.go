package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/ethanvc/evol/alertwhisperer/controller"
	"github.com/ethanvc/evol/alertwhisperer/domain"
	"github.com/ethanvc/evol/svrkit"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net"
	"os"
)

func main() {
	app := fx.New(
		fx.Provide(NewConfig),
		fx.Provide(NewDatabase),
		fx.Provide(controller.NewAlertRuleController),
		fx.Provide(NewHttpServer),
		fx.Provide(domain.NewAlertRuleRepository),
		fx.Invoke(func(engine *gin.Engine) {}),
	)
	app.Run()
}

type NewHttpServerParam struct {
	fx.In
	Conf      *Config
	AlertRule *controller.AlertRuleController
}

func NewHttpServer(lc fx.Lifecycle, param NewHttpServerParam) *gin.Engine {
	conf := param.Conf
	engine := gin.New()
	registerControllers(engine, param)
	var ln net.Listener
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			var err error
			ln, err = net.Listen("tcp", conf.Server.ListenAddr)
			if err != nil {
				return err
			}
			go func() {
				err := engine.RunListener(ln)
				panic(err)
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return ln.Close()
		},
	})
	return engine
}

func registerControllers(engine *gin.Engine, param NewHttpServerParam) {
	var interceptors []svrkit.InterceptorFunc
	interceptors = append(interceptors,
		svrkit.NewAccessInterceptor().Intercept,
		svrkit.NewHttpEncoder().Intercept,
		svrkit.NewHttpDecoder().Intercept,
	)
	{
		controller := param.AlertRule
		r := engine.Group("/api/alert-rule/")
		r.POST("/create", svrkit.NewGinChain(interceptors, controller.CreateAlertRule))
	}
}

type Config struct {
	Server   Server   `toml:"server"`
	Database Database `toml:"database"`
}

type Server struct {
	ListenAddr string `toml:"listen_addr"`
}

type Database struct {
	User         string `toml:"user"`
	Password     string `toml:"password"`
	Host         string `toml:"host"`
	Port         int    `toml:"port"`
	DatabaseName string `toml:"database_name"`
}

func NewConfig() (*Config, error) {
	flagSet := flag.NewFlagSet("", flag.ContinueOnError)
	configFile := flagSet.String("config", "config.toml", "config file")
	err := flagSet.Parse(os.Args[1:])
	if err != nil {
		return nil, err
	}
	var conf Config
	_, err = toml.DecodeFile(*configFile, &conf)
	if err != nil {
		return nil, err
	}
	return &conf, nil
}

func NewDatabase(conf *Config) (*gorm.DB, error) {
	dbConf := conf.Database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbConf.User, dbConf.Password, dbConf.Host, dbConf.Port, dbConf.DatabaseName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err = db.AutoMigrate(&domain.AlertRule{}); err != nil {
		return nil, err
	}
	return db, nil
}
