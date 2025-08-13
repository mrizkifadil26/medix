package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mrizkifadil26/medix/logger"
	"github.com/mrizkifadil26/medix/server"
	"github.com/mrizkifadil26/medix/webgen"
)

func main() {
	logger.Info("🔁 Starting dev server with auto-rebuild")

	if err := webgen.GenerateSite("data", "dist"); err != nil {
		logger.Error("❌ Initial site generation failed: " + err.Error())
	}
	logger.Info("Initial site generation complete")

	// Setup graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go server.WatchAndBuild()
	go server.OpenBrowser("http://localhost:8080")

	// Serve in main goroutine
	go func() {
		if err := server.Serve("dist", "8080"); err != nil {
			logger.Error("❌ Server failed: " + err.Error())
			stop() // trigger shutdown
		}
	}()

	// Wait for interrupt
	<-ctx.Done()
	logger.Warn("⚠️  Received shutdown signal. Cleaning up...")

	// Optional: give time to finish writes, close file handles, etc.
	time.Sleep(300 * time.Millisecond)
	logger.Info("👋 Gracefully exited.")
}
