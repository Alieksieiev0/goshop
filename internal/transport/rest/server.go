package rest

import "github.com/gin-gonic/gin"

type Server struct {
	r *gin.Engine
}

func NewServer(r *gin.Engine) *Server {
	return &Server{
		r: r,
	}
}

func (s *Server) Start(addr ...string) error {
	s.r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	return s.r.Run(addr...)
}
