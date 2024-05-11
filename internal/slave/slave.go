package slave

import (
	"agent/internal/arbiter"
	"agent/internal/config"
	"agent/internal/database"
	"os/exec"
	"time"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

func RunSlave(cfg *config.Config) {
	log.Info().Msg("Run as Slave")
	checkCluster(cfg)
	for {
		time.Sleep(1 * time.Second) // Проверяем состояние каждые 10 секунд
		arbiter, err := arbiter.CheckArbiter(cfg)
		if err == nil && !arbiter && !database.CheckDB(cfg, database.Master) {
			log.Info().Msg("Promote to Master")

			cmd := exec.Command("touch", "/tmp/touch_me_to_promote_to_me_master")
			err := cmd.Run()

			if err == nil {
				log.Info().Msg("Success promote to Master")
				break
			}

			log.Info().Err(err).Msg("Error promote to Master")
		}
	}
}
func checkCluster(cfg *config.Config) {
	for {
		_, err := arbiter.CheckArbiter(cfg)
		if err == nil && !database.CheckDB(cfg, database.Master) {
			break
		}
		log.Info().Msg("Waiting for cluster")
		time.Sleep(5 * time.Second)
	}
}
