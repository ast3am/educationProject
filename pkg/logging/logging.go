package logging

import (
	"github.com/rs/zerolog"
	"net/http"
	"os"
	"strconv"
	"time"
)

func Get() zerolog.Logger {
	log := (zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).
		Level(zerolog.DebugLevel)).
		With().
		Timestamp().
		Logger()
	return log
}

type Logger struct {
	zerolog.Logger
}

func GetLogger() *Logger {
	return &Logger{Get()}
}

func (l *Logger) HandlerLog(r *http.Request, status int, msg string) {
	code := strconv.Itoa(status)
	l.Info().Str("method", r.Method).
		Str("host", r.Host).
		Str("URL", r.RequestURI).
		Str("from", r.RemoteAddr).
		Str("status", code).
		Msg(msg)
}

func (l *Logger) HandlerErrorLog(r *http.Request, status int, msg string, err error) {
	code := strconv.Itoa(status)
	l.Error().Str("method", r.Method).
		Str("host", r.Host).
		Str("URL", r.RequestURI).
		Str("from", r.RemoteAddr).
		Str("status", code).
		Err(err).
		Msg(msg)
}
