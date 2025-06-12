package logger

import (
	"cmp"
	"log/slog"
	"os"

	"github.com/golang-cz/devslog"
)

type Options struct {
	LoggerConfig Config

	LogHandlers []func(slog.Handler) slog.Handler

	// ServiceName specifies the service name. Defaults to "unknown" if not set
	ServiceName string
	// Version specifies the service version. Defaults to "unknown" if not set
	Version string

	// If JSON is true, use JSON logger. Use pretty logger otherwise.
	JSON bool
}

type Config struct {
	Level   slog.Level `toml:"level"`
	JSON    bool       `toml:"json"`
	Concise bool       `toml:"concise"`
}

var defaultOptions = &Options{
	LoggerConfig: Config{
		Level: slog.LevelInfo,
		JSON:  true,
	},
	ServiceName: "undefined",
	Version:     "undefined",
}

func New(o *Options) *slog.Logger {
	o = cmp.Or(o, defaultOptions)

	handlerOptions := &slog.HandlerOptions{
		AddSource: true,
		Level:     o.LoggerConfig.Level,
	}

	var slogHandler slog.Handler
	if o.JSON && !o.LoggerConfig.JSON {
		// Pretty logger for localhost development.
		slogHandler = devslog.NewHandler(os.Stdout, &devslog.Options{
			MaxSlicePrintSize: 20,
			SortKeys:          true,
			TimeFormat:        "[15:04:05.000]",
			StringerFormatter: true,
			HandlerOptions:    handlerOptions,
		})
	} else {
		// JSON logger for production
		slogHandler = slog.NewJSONHandler(os.Stdout, handlerOptions)
	}

	for _, handler := range o.LogHandlers {
		slogHandler = handler(slogHandler)
	}

	logger := slog.New(slogHandler)

	if !o.JSON {
		logger.With(
			slog.String("service", cmp.Or(o.ServiceName, "undefined")),
			slog.String("version", cmp.Or(o.Version, "undefined")),
		)
	}

	// set default log and slog logger, if somebody would use plain slog or log
	// we would at least get proper JSON log format
	slog.SetDefault(logger)

	// set log error level for plain log, we shouldn't use log at all,
	// but if that happen we can see the logs and errors and resolve it
	// only place where we use it is in main for log.Fatalf()
	slog.SetLogLoggerLevel(slog.LevelError)

	return logger
}
