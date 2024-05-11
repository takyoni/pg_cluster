package main

import (
	arb "agent/internal/arbiter"
	ct "agent/internal/cluster"
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

	cluster := ct.Init(config)
	defer cluster.Close()

	switch strings.ToLower(config.ROLE) {
	case "arbiter":
		arb.RunArbiter(cluster)
	case "master":
		mr.RunMaster(cluster)
	case "slave":
		sl.RunSlave(cluster)
	}
}
