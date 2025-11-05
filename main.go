package main

import (
	"context"
	"os/signal"
	"syscall"

	"bitbucket.org/xoduxcrt/ctlog-update-processor/certwatch"
	"bitbucket.org/xoduxcrt/ctlog-update-processor/ct"
	"bitbucket.org/xoduxcrt/ctlog-update-processor/logger"
	"bitbucket.org/xoduxcrt/ctlog-update-processor/msg"
	"bitbucket.org/xoduxcrt/ctlog-update-processor/server"
	"bitbucket.org/xoduxcrt/ctlog-update-processor/utils"
)

func main() {
	// The certwatch database connections, which were opened automatically by the init() function, need to be closed on exit.
	defer certwatch.Close()

	// utils.ScrapeCTLogList()
	// certwatch.InitConfig()

	// Configure graceful shutdown capabilities.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	defer msg.ShutdownWG.Wait()

	// Start the N goroutines.
	msg.ShutdownWG.Add(4)
	go ct.LogConfigSyncer(ctx)
	go ct.GetEntriesLauncher(ctx)
	go certwatch.NewEntriesWriter(ctx)
	go utils.RTUpdate(ctx)

	// Start the Monitoring HTTP server.
	server.Run()
	defer server.Shutdown()

	// Wait to be interrupted.
	<-ctx.Done()

	// Ensure all log messages are flushed before we exit.
	logger.Logger.Sync()
}
