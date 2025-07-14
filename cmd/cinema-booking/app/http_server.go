package app

import (
	"context"
	"net/http"
	"time"

	"github.com/1612421/cinema-booking/pkg/go-kit/log"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	"github.com/1612421/cinema-booking/config"
	"github.com/1612421/cinema-booking/internal"
)

type httpServer struct {
	server   *http.Server
	cleanups []func()
}

var (
	skipPaths = []string{"/metrics", "/health", "/info"}
)

func (h *httpServer) Start() {
	log.Bg().Info("HTTP server is starting", log.String("address", h.server.Addr))
	err := h.server.ListenAndServe()
	if err != nil {
		return
	}
}

func (h *httpServer) Shutdown() {
	log.Bg().Info("Shutting down HTTP server...")

	// The context is used to inform the server it has 10 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := h.server.Shutdown(ctx)
	if err != nil {
		log.Bg().Error("Shut down http server failed", log.Error(err))
	} else {
		log.Bg().Info("Shutting down http server successfully")
	}

	for _, cleanup := range h.cleanups {
		cleanup()
	}

	log.Bg().Info("Shut down HTTP server successfully")
}

func NewHTTPServer() Application {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.ContextWithFallback = true // enable to use context in request
	// Middlewares
	r.Use(
		gin.Recovery(),
		otelgin.Middleware(config.GetConfig().Service.Name, otelgin.WithFilter(filterPaths)),
	)

	controller, cleanup, err := internal.InitializeHTTPAPIController()
	if err != nil {
		log.Bg().Fatal("Failed to initialize HTTP API controller", log.Error(err))
	}

	controller.SetupRouter(r)

	server := &http.Server{
		Addr:              config.GetConfig().Service.HTTP.Address(),
		Handler:           r,
		ReadHeaderTimeout: 3 * time.Second, // G112: Potential Slowloris Attack https://app.deepsource.com/directory/analyzers/go/issues/GO-S2112
	}

	return &httpServer{
		server:   server,
		cleanups: []func(){cleanup},
	}
}

func filterPaths(request *http.Request) bool {
	for _, skipPath := range skipPaths {
		if request.URL.Path == skipPath {
			return false
		}
	}
	return true
}
