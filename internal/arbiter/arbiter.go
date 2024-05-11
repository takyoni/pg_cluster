package arbiter

import (
	"agent/internal/config"
	"agent/internal/database"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

type Server struct {
	config *config.Config
}

func RunArbiter(cfg *config.Config) {
	log.Info().Msg("Run as Arbiter")
	handler := &Server{config: cfg}

	server := gin.Default()
	server.GET("/master", handler.MasterStatus)
	server.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, "pong")
	})

	server.Run(":8080")
}

func (s *Server) MasterStatus(c *gin.Context) {
	result := database.CheckDB(s.config, database.Master)
	log.Info().Bool("result", result).Msg("Check Master")

	if !result {
		c.JSON(http.StatusBadGateway, gin.H{"master": result})
	} else {
		c.JSON(http.StatusOK, gin.H{"master": result})
	}
}
