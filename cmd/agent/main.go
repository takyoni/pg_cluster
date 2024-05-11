package main

import (
	arb "agent/internal/arbiter"
	cfg "agent/internal/config"
	"agent/internal/logger"
	mr "agent/internal/master"
	sl "agent/internal/slave"
	"strings"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

// Основная функция агента
func main() {
	config, err := cfg.Load()
	if err != nil {
		return
	}

	logger.Setup()
	log.Info().Msg("Success parsed config")

	switch strings.ToLower(config.ROLE) {
	case "arbiter":
		arb.RunArbiter(config)
	case "master":
		mr.RunMaster(config)
	case "slave":
		sl.RunSlave(config)
	}
}
