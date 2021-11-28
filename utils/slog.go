package utils

import (
	"github.com/gookit/slog"
)

func InitSlog(verbose bool) {
	slog.Configure(func(logger *slog.SugaredLogger) {
		f := logger.Formatter.(*slog.TextFormatter)
		f.Template = "[{{datetime}}] [{{level}}] {{message}} {{data}}\n"
		f.EnableColor = true
		if verbose {
			logger.Level = slog.TraceLevel
		} else {
			logger.Level = slog.InfoLevel
		}
	})
}
