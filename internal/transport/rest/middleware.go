package rest

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger(logger *log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		c.Next()

		logger.Printf("Request latency: %v", time.Since(t))

		for i, err := range c.Errors {
			logger.Printf("Error #%d: %v", i+1, err)
		}

		logger.Printf("Returned status: %d", c.Writer.Status())
	}
}
