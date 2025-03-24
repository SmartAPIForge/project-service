package main

import (
	"os"
	"os/signal"
	"project-service/internal/app"
	"project-service/internal/config"
	"project-service/internal/lib/logger"
	"syscall"
)

func main() {
	cfg := config.MustLoad()
	log := logger.MustSetupLogger(cfg.Env)

	application := app.NewApp(log, cfg)
	application.GrpcApp.MustRun()

	stopWait(application)
}

func stopWait(application *app.App) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop
	application.GrpcApp.Stop()
}
