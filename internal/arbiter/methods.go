package arbiter

import (
	"agent/internal/config"
	"fmt"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

func CheckArbiter(cfg *config.Config) bool {
	psqlInfo := fmt.Sprintf("http://%s:8080/master",
		cfg.ARBITER_HOST)
	result, err := http.Get(psqlInfo)
	if err != nil || result.StatusCode != http.StatusOK {
		log.Info().Bool("result", false).Msg("Check Arbiter")
		return false
	}
	log.Info().Bool("result", true).Msg("Check Arbiter")
	return true
}
