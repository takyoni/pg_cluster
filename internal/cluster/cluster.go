package cluster

import (
	"agent/internal/config"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

type HostRole int32

type Replicas struct {
	SlaveConn   *sql.DB
	MasterConn  *sql.DB
	ArbiterHost string
}

func Init(cfg *config.Config) *Replicas {
	replica := &Replicas{}
	for {
		var err error

		if cfg.MASTER_HOST != "" {
			replica.MasterConn, err = connectToDatabase(cfg.MASTER_HOST, cfg.POSTGRES_USER, cfg.POSTGRES_PASSWORD)
		}

		if cfg.SLAVE_HOST != "" {
			replica.SlaveConn, err = connectToDatabase(cfg.SLAVE_HOST, cfg.POSTGRES_USER, cfg.POSTGRES_PASSWORD)
		}

		if cfg.ARBITER_HOST != "" {
			replica.ArbiterHost = cfg.ARBITER_HOST
		}

		if err == nil {
			log.Info().Msg("Success start")
			break
		}

		log.Info().Msg("Waiting for other hosts")
		time.Sleep(2 * time.Second)
	}
	return replica

}

func connectToDatabase(host, user, password string) (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, 5432, user, password, "postgres")
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (r *Replicas) CheckMaster() bool {
	err := r.MasterConn.Ping()
	if err != nil {
		log.Err(err).Msg("Check Master")
	}
	log.Info().Bool("result", err == nil).Msg("Check Master")
	return err == nil
}

func (r *Replicas) CheckSlave() bool {
	err := r.SlaveConn.Ping()
	if err != nil {
		log.Err(err).Msg("Check Slave")
	}
	log.Info().Bool("result", err == nil).Msg("Check Master")
	return err == nil
}

func (r *Replicas) CheckAM() (bool, error) {
	psqlInfo := fmt.Sprintf("http://%s:8080/master",
		r.ArbiterHost)
	result, err := http.Get(psqlInfo)
	if err != nil {
		log.Info().Bool("result", false).Msg("Check Arbiter")
		return false, err
	}
	if result.StatusCode != http.StatusOK {
		log.Info().Bool("result", false).Msg("Check Arbiter")
		return false, nil
	}
	log.Info().Bool("result", true).Msg("Check Arbiter")
	return true, nil
}

func (r *Replicas) CheckArbiter() (bool, error) {
	psqlInfo := fmt.Sprintf("http://%s:8080/ping",
		r.ArbiterHost)
	result, err := http.Get(psqlInfo)
	if err != nil {
		log.Info().Bool("result", false).Msg("Check Arbiter")
		return false, err
	}
	if result.StatusCode != http.StatusOK {
		log.Info().Bool("result", false).Msg("Check Arbiter")
		return false, nil
	}
	log.Info().Bool("result", true).Msg("Check Arbiter")
	return true, nil
}

func (r *Replicas) Close() {
	r.MasterConn.Close()
	r.SlaveConn.Close()
}
