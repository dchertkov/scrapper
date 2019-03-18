package main

import (
	"context"
	"os"

	"github.com/dchertkov/scrapper/pkg/api"
	"github.com/dchertkov/scrapper/pkg/checker"
	"github.com/dchertkov/scrapper/pkg/config"
	"github.com/dchertkov/scrapper/pkg/file"
	"github.com/dchertkov/scrapper/pkg/server"
	"github.com/dchertkov/scrapper/pkg/signal"
	"github.com/dchertkov/scrapper/pkg/store/service"
	"github.com/dchertkov/scrapper/pkg/store/stat"

	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

func main() {
	conf, err := config.FromEnv()
	if err != nil {
		logrus.WithError(err).Fatal("loading config")
	}

	app, err := initApp(conf)
	if err != nil {
		logrus.WithError(err).Fatal("init application")
	}

	ctx := signal.WithCancel(context.Background(), func(sig os.Signal) {
		logrus.WithField("signal", sig).Info("terminating process")
	})

	g := errgroup.Group{}

	g.Go(func() error {
		return app.checker.Start(ctx)
	})

	g.Go(func() error {
		logrus.WithFields(logrus.Fields{
			"addr": conf.Server.Addr(),
		}).Info("listening server")
		return app.server.Listen(ctx)
	})

	if err = g.Wait(); err != nil {
		logrus.WithError(err).Fatal("exit")
	}
}

type app struct {
	server  *server.Server
	checker *checker.Checker
}

func initApp(conf *config.Config) (*app, error) {
	serviceList, err := file.Load(conf.SourceFile)
	if err != nil {
		return nil, err
	}

	serviceStore := service.NewStore()
	statStore := stat.NewStore()
	handler := api.NewHandler(serviceStore, statStore)

	return &app{
		server:  server.NewServer(conf.Server, handler.Handler()),
		checker: checker.NewChecker(conf.Checker, serviceStore, serviceList),
	}, nil
}
