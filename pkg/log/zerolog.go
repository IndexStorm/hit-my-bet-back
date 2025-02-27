package log

import (
	"github.com/rs/zerolog"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

type LevelOutputWriter struct {
	io.Writer
	m map[zerolog.Level]io.Writer
}

func (w LevelOutputWriter) WriteLevel(level zerolog.Level, p []byte) (int, error) {
	if dst, ok := w.m[level]; ok {
		return dst.Write(p)
	}
	return os.Stdout.Write(p)
}

func NewWithLevel(lvl zerolog.Level) zerolog.Logger {
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.TimestampFieldName = "timestamp"
	zerolog.LevelFieldName = "severity"
	zerolog.MessageFieldName = "message"
	zerolog.DurationFieldInteger = true
	return zerolog.New(
		LevelOutputWriter{
			m: map[zerolog.Level]io.Writer{
				zerolog.WarnLevel:  os.Stderr,
				zerolog.ErrorLevel: os.Stderr,
				zerolog.FatalLevel: os.Stderr,
				zerolog.PanicLevel: os.Stderr,
			},
		},
	).With().Timestamp().Caller().Stack().Logger().Level(lvl)
}

func SetupCallerRootRewrite(roots ...string) {
	if len(roots) == 0 {
		roots = []string{
			"cmd/",
			"internal/",
			"pkg/",
			"testing/",
			"tests/",
		}
	}
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		for _, root := range roots {
			if i := strings.Index(file, root); i > 0 {
				file = file[i:]
				break
			}
		}
		return file + ":" + strconv.Itoa(line)
	}
}
