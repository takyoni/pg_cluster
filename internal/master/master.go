package master

import (
	"agent/internal/arbiter"
	"agent/internal/config"
	"agent/internal/database"
	"os/exec"
	"time"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

func RunMaster(cfg *config.Config) {
	log.Info().Msg("Run as Master")
	checkCluster(cfg)
	for {
		time.Sleep(1 * time.Second)
		arbiter, err := arbiter.CheckArbiter(cfg)
		if err == nil && !arbiter && !database.CheckDB(cfg, database.Slave) {
			cmd := exec.Command("iptables", "-P", "INPUT", "DROP")
			err := cmd.Run()

			if err == nil {
				log.Info().Msg("Success block connections to Master")
				break
			}

			log.Info().Err(err).Msg("Block input connection to Master")
		}
	}
}

func checkCluster(cfg *config.Config) {
	for {
		_, err := arbiter.CheckArbiter(cfg)
		if err == nil && !database.CheckDB(cfg, database.Slave) {
			break
		}
		log.Info().Msg("Waiting for cluster")
		time.Sleep(5 * time.Second)
	}
}
