package main

import (
	"context"
	"github.com/web-page-analysis/bootstrap"
	"github.com/web-page-analysis/container"
	"github.com/web-page-analysis/server"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cancel()
	}()
	// config loader
	conf, err := bootstrap.InitConfig()
	if err != nil {
		return
	}

	// container resolver
	ctr := container.Resolver(ctx, conf)

	// server start
	server.InitRouter(ctx, conf, *ctr)
	<-ctx.Done()
}
