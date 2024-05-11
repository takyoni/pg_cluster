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
	for {
		time.Sleep(5 * time.Second) // Проверяем состояние каждые 10 секунд
		if !arbiter.CheckArbiter(cfg) && !database.CheckSlave(cfg) {
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
