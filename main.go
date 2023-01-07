package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ribincao/ribin-game-logic/constant"
	"github.com/ribincao/ribin-game-server/config"
	"github.com/ribincao/ribin-game-server/logger"
	"github.com/ribincao/ribin-game-server/server"
	"github.com/ribincao/ribin-game-server/utils"
)

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
	srv := server.NewServer(server.RoomServer)
	// TODO: Set Handler / Set ConnCloseCallBack

	utils.GoWithRecover(srv.Serve)
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

	// TODO: Destruct: Destroy Room

	go func() {
		time.Sleep(1 * time.Second)
		os.Exit(0)
	}()
}
