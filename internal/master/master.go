package master

import (
	"agent/internal/cluster"
	"os/exec"
	"time"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

func RunMaster(ct *cluster.Replicas) {
	log.Info().Msg("Run as Master")
	for {
		time.Sleep(1 * time.Second)
		arbiter, err := ct.CheckArbiter()
		if err == nil && !arbiter && !ct.CheckSlave() {
			log.Info().Msg("Start blocking input connection to Master")

			cmd := exec.Command("iptables", "-P", "INPUT", "DROP")
			err := cmd.Run()

			if err == nil {
				log.Info().Msg("Success block connections to Master")
			}

			cmd = exec.Command("iptables-save", ">", "/etc/iptables/rules.v4")
			err = cmd.Run()

			if err == nil {
				log.Info().Msg("Success save iptables changes")
				break
			}

			log.Info().Err(err).Msg("Block input connection to Master")
		}
	}
}
