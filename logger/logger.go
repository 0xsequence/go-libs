package logger

import (
	"cmp"
	"log/slog"
	"os"

	"github.com/go-chi/httplog/v3"
	"github.com/go-chi/traceid"
	"github.com/golang-cz/devslog"

	"github.com/0xsequence/go-libs/httpdebug"
)

type Options struct {
	Config Config

	// ServiceName specifies the service name. Defaults to "unknown" if not set
	ServiceName string

	// Version specifies the service version. Defaults to "unknown" if not set
	Version string

	// Use httpdebug header in logging
	HTTPDebug *httpdebug.Header
}

// Config can be used directly in toml config
type Config struct {
	Level   slog.Level `toml:"level"`
	Concise bool       `toml:"concise"`
	Pretty  bool       `toml:"pretty"`
}

var defaultOptions = &Options{
	Config: Config{
		Level:  slog.LevelInfo,
		Pretty: false,
	},
	ServiceName: "unknown",
	Version:     "unknown",
}

func New(o *Options) *slog.Logger {
	o = cmp.Or(o, defaultOptions)

	handlerOptions := &slog.HandlerOptions{
		AddSource:   true,
		Level:       o.Config.Level,
		ReplaceAttr: httplog.SchemaGCP.Concise(o.Config.Concise).ReplaceAttr,
	}

	var slogHandler slog.Handler
	if o.Config.Pretty {
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

	// add traceid handler
	slogHandler = traceid.LogHandler(slogHandler)

	if o.HTTPDebug != nil && o.HTTPDebug.Key != "" && o.HTTPDebug.Value != "" {
		slogHandler = httpdebug.LogHandler(*o.HTTPDebug)(slogHandler)
	}

	logger := slog.New(slogHandler)

	// in JSON mode print service and version
	if !o.Config.Pretty {
		logger = logger.With(
			slog.String("service", cmp.Or(o.ServiceName, "unknown")),
			slog.String("version", cmp.Or(o.Version, "unknown")),
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

type MigrationLogger struct {
	*slog.Logger
}

func (l MigrationLogger) Printf(format string, v ...interface{}) {
	l.Logger.Info(fmt.Sprintf(format, v...))
}

func (l MigrationLogger) Fatalf(format string, v ...interface{}) {
	l.Logger.Error(fmt.Sprintf(format, v...))
}

