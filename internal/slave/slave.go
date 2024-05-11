package slave

import (
	"agent/internal/cluster"
	"os/exec"
	"time"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

func RunSlave(ct *cluster.Replicas) {
	log.Info().Msg("Run as Slave")
	for {
		time.Sleep(1 * time.Second)
		arbiter, err := ct.CheckAM()
		if err == nil && !arbiter && !ct.CheckMaster() {
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
