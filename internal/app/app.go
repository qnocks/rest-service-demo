package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"team-task/internal/config"
	"team-task/internal/service"
	"team-task/internal/storage"
	"team-task/internal/stream"
	"team-task/internal/transport/http"
	"team-task/pkg/httpserver"
	"team-task/pkg/logger"
	"time"
)

const timeout = 5 * time.Second

func Run() {
	cfg := config.GetConfig()

	stanClient, err := stream.NewSTANClient(cfg.STAN.ClusterID, cfg.STAN.ClientID, cfg.STAN.URL, cfg.STAN.Subject)
	if err != nil {
		logger.Errorf("error connecting to stan: %s\n", err.Error())
		return
	}

	store := storage.NewStorage()
	services := service.NewService(store, stanClient)
	handler := http.NewHandler(services)
	setRouter, getRouter := handler.Init(cfg)

	firstServer := httpserver.NewServer(cfg.Server.PortOne, setRouter)
	secondServer := httpserver.NewServer(cfg.Server.PortTwo, getRouter)

	go func() {
		logger.Infof("running server on port: %s\n", cfg.Server.PortOne)
		if err := firstServer.Run(); err != nil {
			logger.Errorf("error running http server: %s\n", err.Error())
			return
		}
	}()

	go func() {
		logger.Infof("running server on port: %s\n", cfg.Server.PortTwo)
		if err := secondServer.Run(); err != nil {
			logger.Errorf("error running http server: %s\n", err.Error())
			return
		}
	}()

	go func() {
		if err := stanClient.Listen(store); err != nil {
			logger.Errorf("error listening to streaming connection: %s\n", err.Error())
			return
		}
	}()

	graceful()
	shutdown(firstServer, secondServer, stanClient)
}

func graceful() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}

func shutdown(firstServer *httpserver.Server, secondServer *httpserver.Server, stanClient *stream.STANClient) {
	ctx, shutdownCancel := context.WithTimeout(context.Background(), timeout)
	defer shutdownCancel()

	if err := firstServer.Shutdown(ctx); err != nil {
		logger.Errorf("failed to stop server: %s\n", err.Error())
		return
	}

	if err := secondServer.Shutdown(ctx); err != nil {
		logger.Errorf("failed to stop server: %s\n", err.Error())
		return
	}

	if err := stanClient.Shutdown(); err != nil {
		logger.Errorf("failed to close stan connection: %s\n", err.Error())
		return
	}
}
