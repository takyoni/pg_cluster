package logger

import (
	"fmt"
	"os"
	"strings"

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
	output := zerolog.ConsoleWriter{Out: runLogFile}
	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
	}
	output.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("***%s****", i)
	}
	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s:", i)
	}
	output.FormatFieldValue = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("%s", i))
	}
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: runLogFile})

	//multi := zerolog.MultiLevelWriter(output, runLogFile)
	// log.Logger = zerolog.New(multi).With().Timestamp().Logger()
}
