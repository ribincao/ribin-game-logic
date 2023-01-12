package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ribincao/ribin-game-logic/constant"
	"github.com/ribincao/ribin-game-logic/handler"
	"github.com/ribincao/ribin-game-server/config"
	"github.com/ribincao/ribin-game-server/logger"
	"github.com/ribincao/ribin-game-server/server"
	"github.com/ribincao/ribin-game-server/utils"
	"go.uber.org/zap"
)

const TestPort = 8080

func main() {
	initLogger()

	ctx, cancel := context.WithCancel(context.Background())
	run(ctx)
	handleSignal(ctx, cancel)
}

func initLogger() {
	config.ParseConf(constant.CONFIG_PATH, config.GlobalConfig)
	config.GlobalConfig.LogConfig.LogPath = fmt.Sprintf("%v-%s", os.Getppid(), "server.log")
	logger.InitLogger(config.GlobalConfig.LogConfig)
}

func run(ctx context.Context) {
	var port int32
	// TODO: Match-Server allocate Server
	if config.GlobalConfig.ServiceConfig.Env == "local" {
		port = TestPort
	}
	srv := server.NewServer(server.RoomServer, server.WithAddress(fmt.Sprintf(":%d", port)))
	srv.SetConnCloseCallback(handler.OnClose)
	srv.SetHandler(handler.HandleServerMessage)

	utils.GoWithRecover(srv.Serve)
	logger.Info("Server Start Success.", zap.Any("Port", port))
}

func handleSignal(ctx context.Context, cancel context.CancelFunc) {
	sigC := make(chan os.Signal, 2)
	signal.Notify(sigC,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGABRT,
		syscall.SIGSEGV)

	sig := <-sigC
	switch sig {
	case syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGSEGV:
		logger.Error("Report Crash")
	}

	cancel()
	go func() {
		time.Sleep(1 * time.Second)
		os.Exit(0)
	}()
}
