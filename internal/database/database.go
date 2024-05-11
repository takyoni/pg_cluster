package database

import (
	"agent/internal/config"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

type HostRole int32

const (
	Master HostRole = iota
	Slave
)

func CheckDB(cfg *config.Config, role HostRole) bool {
	host := ""
	switch role {
	case Master:
		host = cfg.MASTER_HOST
	case Slave:
		host = cfg.SLAVE_HOST
	}
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, 5432, cfg.POSTGRES_USER, cfg.POSTGRES_PASSWORD, "postgres")
	conn, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Info().Bool("result", false).Msg("Check Master")
		return false
	}

	defer conn.Close()
	err = conn.Ping()
	log.Info().Bool("result", err == nil).Msg("Check Master")
	return err == nil
}
