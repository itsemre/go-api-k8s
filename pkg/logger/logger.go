package logger

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sbecker/gin-api-demo/util"
	log "github.com/sirupsen/logrus"
)

const (
	ISO8601layout = "2006-01-02T15:04:05-0700"
)

// InitLogger creates a new Logrus logger instance
func InitLogger(logLevel, logFormat string) (*log.Logger, error) {
	level, err := log.ParseLevel(logLevel)
	if err != nil {
		return &log.Logger{}, errors.Wrapf(err, "error initialising logger")
	}

	var logger = &log.Logger{
		Out:   os.Stdout,
		Level: level,
	}

	if logFormat == "json" {
		logger.SetFormatter(&log.JSONFormatter{
			TimestampFormat: ISO8601layout,
			FieldMap: log.FieldMap{
				"msg":  "message",
				"time": "timestamp",
			},
		})
	} else {
		logger.SetFormatter(&log.TextFormatter{})
	}

	return logger, nil
}

// JSONLogger is a Gin handler function that defines our JSON logging structure
func JSONLogger(logger *log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		// Process Request
		c.Next()
		// Stop timer
		duration := util.GetDurationInMillseconds(start)

		entry := logger.WithFields(log.Fields{
			"duration":    duration,
			"client_ip":   util.GetClientIP(c),
			"url":         c.Request.URL.Path,
			"status_code": c.Writer.Status(),
			"method":      c.Request.Method,
			"headers": log.Fields{
				"referer":    c.GetHeader("Referer"),
				"user_agent": c.GetHeader("User-Agent"),
			},
		})

		if c.Writer.Status() >= 400 {
			entry.Error(c.Errors.String())
		} else {
			entry.Info()
		}
	}
}
