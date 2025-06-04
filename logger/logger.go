package logger

import (
	"cmp"
	"log/slog"
	"os"

	"github.com/go-chi/traceid"
	"github.com/golang-cz/devslog"

	"github.com/0xsequence/go-libs/debugger"
)

type Options struct {
	LoggerConfig Config

	// if the debug client is passed, then the logger would use it
	DebugClient *debugger.Client

	// if the service name is not specified, then it would add "unknown"
	ServiceName string
	// if the version name is not specified, then it would add "unknown"
	Version string

	// when the development is true, then the logger would use devslog
	DevelopmentMode bool
	DevslogOptions  *devslog.Options
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
	if o.DevelopmentMode && !o.LoggerConfig.JSON {
		// Pretty logger for localhost development.
		o.DevslogOptions = cmp.Or(
			o.DevslogOptions,
			&devslog.Options{
				// default options, if user doesn't send options
				MaxSlicePrintSize: 20,
				SortKeys:          true,
				TimeFormat:        "[15:04:05.000]",
				StringerFormatter: true,
			},
		)

		// apply handler options to devslog options
		o.DevslogOptions.HandlerOptions = handlerOptions

		slogHandler = devslog.NewHandler(os.Stdout, o.DevslogOptions)
	} else {
		// JSON logger for production
		slogHandler = slog.NewJSONHandler(os.Stdout, handlerOptions)
	}

	// Log "traceId"
	slogHandler = traceid.LogHandler(slogHandler)

	if o.DebugClient != nil {
		slogHandler = o.DebugClient.LogHandler(slogHandler)
	}

	logger := slog.New(slogHandler)

	if !o.DevelopmentMode {
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
