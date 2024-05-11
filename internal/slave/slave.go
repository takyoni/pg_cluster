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
	for {
		time.Sleep(5 * time.Second)                                   // Проверяем состояние каждые 10 секунд
		if !database.CheckMaster(cfg) && !arbiter.CheckArbiter(cfg) { // Проверяем доступность мастера и арбитра
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
