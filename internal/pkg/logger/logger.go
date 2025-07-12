package logger

import (
	"io"
	"os"
	"strconv"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func New() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	zerolog.TimestampFieldName = "t"
	zerolog.LevelFieldName = "l"
	zerolog.MessageFieldName = "m"

	logLevelStr := os.Getenv("LOG_LEVEL")
	logLevel, err := strconv.Atoi(logLevelStr)
	if err != nil || logLevel > 7 || logLevel < 0 {
		logLevel = int(zerolog.InfoLevel)
	}

	var output io.Writer = os.Stdout
	if os.Getenv("APP_ENV") == "development" {
		output = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}
	}

	log.Logger = zerolog.New(output).
		Level(zerolog.Level(logLevel)).
		With().
		Timestamp().
		Logger()

	log.Info().Msg("Logger initialized")
}
