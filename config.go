package log

import (
	"go.uber.org/zap"
)

//ZapLogger zap logs implement.
type ZapLogger struct {
	Options *Options
	*zap.SugaredLogger
}

//Configure the log print.
func configure(zap *ZapLogger, ops ...Option) {
	// default value.
	zap.Options.LocalTime = true
	zap.Options.Compress = true
	// Deal with options which are set.
	for _, o := range ops {
		o(zap.Options)
	}
	// When args is empty, use default.
	if zap.Options.Filename == "" {
		zap.Options.Filename = DefaultFilename
	}
	if zap.Options.MaxSize <= 0 {
		zap.Options.MaxSize = DefaultMaxSize
	}
	if zap.Options.Level < -1 || zap.Options.Level > 5 {
		zap.Options.Level = DefaultLevel
	}
}
