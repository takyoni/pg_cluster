package logger

import (
	"os"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Setup() {
	runLogFile, _ := os.OpenFile(
		"myapp.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0664,
	)
	multi := zerolog.MultiLevelWriter(os.Stdout, runLogFile)
	log.Logger = zerolog.New(multi).With().Timestamp().Logger()
}
