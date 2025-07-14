package app

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/1612421/cinema-booking/pkg/go-kit/log"
	"github.com/1612421/cinema-booking/pkg/go-kit/thread"
)

type Application interface {
	Start()
	Shutdown()
}

func Run() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	rg := thread.NewRoutineGroup()

	httpServer := NewHTTPServer()
	// Start run server in goroutines
	rg.Run(httpServer.Start)

	<-stop

	log.Bg().Info("Graceful shutdown")

	// Clean up resources in goroutines
	rg.Run(httpServer.Shutdown)

	rg.Wait()

	log.Bg().Info("Shut down successfully")
}
