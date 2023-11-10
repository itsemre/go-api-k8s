package server

import (
	"context"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/itsemre/go-api-k8s/pkg/config"
	"github.com/itsemre/go-api-k8s/pkg/controller"
	log "github.com/itsemre/go-api-k8s/pkg/logger"
	"github.com/sirupsen/logrus"
	ginprometheus "github.com/zsais/go-gin-prometheus"
)

type Server struct {
	HTTPServer *http.Server
	Router     *gin.Engine
	Config     *config.Config
	Logger     *logrus.Logger
}

// NewServer returns an instance of the server
func NewServer(conf *config.Config, logger *logrus.Logger) *Server {
	serverAddress := fmt.Sprintf("%s:%s", conf.ServerAddress, conf.ServerPort)
	router := gin.New()
	router.Use(gin.Recovery(), log.JSONLogger(logger))

	// Set up prometheus middleware to expose metrics
	prom := ginprometheus.NewPrometheus("gin")
	prom.Use(router)

	srv := &http.Server{
		Addr:    serverAddress,
		Handler: router,
	}
	return &Server{
		HTTPServer: srv,
		Router:     router,
		Config:     conf,
		Logger:     logger,
	}
}

// Start starts server
func (s *Server) Start() error {
	// Set CORS settings
	s.Router.Use(cors.New(cors.Config{
		AllowOrigins:     s.Config.CORSAllowOrigins,
		AllowMethods:     s.Config.CORSAllowMethods,
		AllowHeaders:     s.Config.CORSAllowHeaders,
		ExposeHeaders:    s.Config.CORSExposeHeaders,
		AllowCredentials: s.Config.CORSAllowCredentials,
		MaxAge:           time.Duration(s.Config.CORSMaxAge) * time.Hour,
	}))

	// Get a new controller instance
	controller := controller.NewController(s.Config)

	// Assign the Gin handlers to their corresponding URL paths and methods
	s.Router.GET("/comics", controller.GetComics)
	s.Router.GET("/ping", controller.Health)

	// Get context for termination signals
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Initializing the server in a goroutine so that it won't block the graceful shutdown handling below
	go func() {
		if err := s.HTTPServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.Logger.Fatalf("listen: %s\n", err)
		}
	}()

	// Listen for the interrupt signal
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown
	stop()
	s.Logger.Println("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has shutDownTimeout to finish the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.Config.ShutDownTimeout))
	defer cancel()
	if err := s.HTTPServer.Shutdown(ctx); err != nil {
		s.Logger.Fatalf("Server forced to shutdown: %s\n", err)
	}

	s.Logger.Println("Server exiting")
	return nil
}
