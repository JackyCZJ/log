package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

//Option return setted options.
type Option func(*Options)

// Default value.
const (
	// Default Log Level
	DefaultLevel = zapcore.DebugLevel
	// Default File Name
	DefaultFilename = "./app.log"
	// Default Max File Size
	DefaultMaxSize = 100
)

var zapLevel = map[string]zapcore.Level{
	"debug": zapcore.DebugLevel,
	"info":  zapcore.InfoLevel,
	"warn":  zapcore.WarnLevel,
	"error": zapcore.ErrorLevel,
}

//Options stores all options , and the logger which want to use direct.
type Options struct {
	Logger     *zap.SugaredLogger
	Filename   string
	MaxSize    int32
	LocalTime  bool
	Compress   bool
	MaxBackups int
	Level      zapcore.Level
}

//Logger set logger into zaplogger
func Logger(logger *zap.SugaredLogger) Option {
	return func(o *Options) {
		o.Logger = logger
	}
}

//Filename set filename instead of default filename
func Filename(filename string) Option {
	return func(o *Options) {
		o.Filename = filename
	}
}

// MaxSize set maxSize of log file.
func MaxSize(maxSize int32) Option {
	return func(o *Options) {
		o.MaxSize = maxSize
	}
}

// LocalTime Use localtime?
func LocalTime(localTime bool) Option {
	return func(o *Options) {
		o.LocalTime = localTime
	}
}

// Compress Use Compress ?
func Compress(compress bool) Option {
	return func(o *Options) {
		o.Compress = compress
	}
}

// Level Set output log level
func Level(l string) Option {
	return func(o *Options) {
		if level, ok := zapLevel[l]; ok {
			o.Level = level
		} else {
			o.Level = DefaultLevel
		}
	}
}

// FilterOutFunc set filter function.
//func FilterOutFunc(filterOutFunc) Option {
//	return func(o *Options) {
//		o.FilterOutFunc = filterOutFunc
//	}
//}

//MaxBackups max backups number
func MaxBackups(a int) Option {
	return func(options *Options) {
		options.MaxBackups = a
	}
}
