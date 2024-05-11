package database

import (
	"agent/internal/config"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

func CheckMaster(cfg *config.Config) bool {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cfg.MASTER_HOST, 5432, cfg.POSTGRES_USER, cfg.POSTGRES_PASSWORD, "postgres")
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
func CheckSlave(cfg *config.Config) bool {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cfg.SLAVE_HOST, 5432, cfg.POSTGRES_USER, cfg.POSTGRES_PASSWORD, "postgres")
	conn, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return false
	}

	defer conn.Close()
	err = conn.Ping()
	log.Info().Bool("result", err == nil).Msg("Check Slave")
	return err == nil
}
