package log

import (
	"context"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"gopkg.in/natefinch/lumberjack.v2"
)

//Log object for logger use.
var Log *ZapLogger

//NewZapLogger Create a new ZapLogger and return.
func NewZapLogger(opts ...Option) *ZapLogger {
	zapLogger := &ZapLogger{
		Options: new(Options),
	}

	configure(zapLogger, opts...)
	// if there is no logger , create one.
	if zapLogger.Options.Logger == nil {
		// create zap logger.
		syncWriter := zapcore.AddSync(&lumberjack.Logger{
			Filename:   zapLogger.Options.Filename,
			MaxSize:    int(zapLogger.Options.MaxSize),
			LocalTime:  zapLogger.Options.LocalTime,
			Compress:   zapLogger.Options.Compress,
			MaxBackups: zapLogger.Options.MaxBackups,
		})
		pEncoder := zap.NewDevelopmentEncoderConfig()
		pEncoder.EncodeTime = zapcore.ISO8601TimeEncoder // time fommat.
		encoder := zap.NewProductionEncoderConfig()
		encoder.EncodeLevel = zapcore.CapitalColorLevelEncoder //output console color.
		encoder.EncodeTime = zapcore.ISO8601TimeEncoder
		core := zapcore.NewTee(
			zapcore.NewCore(zapcore.NewConsoleEncoder(pEncoder), syncWriter, zap.NewAtomicLevelAt(zapLogger.Options.Level)),                //the core of production , insert log to file.
			zapcore.NewCore(zapcore.NewConsoleEncoder(encoder), zapcore.AddSync(os.Stdout), zap.NewAtomicLevelAt(zapLogger.Options.Level)), //the core for both production and development, log to console.
		)

		logger := zap.New(core, zap.AddCaller())
		zapLogger.SugaredLogger = logger.Sugar()
	} else {
		zapLogger.SugaredLogger = zapLogger.Options.Logger
	}
	Log = zapLogger
	grpclog.SetLoggerV2(&GrpcLog{
		SugaredLogger: zapLogger.SugaredLogger,
	})
	return zapLogger
}

//StreamClient Use ZapLogger to implement middleware
func (zap *ZapLogger) StreamClient(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (cs grpc.ClientStream, err error) {
	//todo: filter func
	defer func() {
		if err != nil {
			zap.SugaredLogger.Errorw("Streamer function", "method", method, "errcode", err)
		} else if zap.Options.Level == zapcore.DebugLevel {
			zap.SugaredLogger.Debugw("Streamer function", "method", method, "errcode", err)
		}
	}()
	cs, err = streamer(ctx, desc, cc, method, opts...)
	return
}

//Info implement for zap log info method.
func Info(arg ...interface{}) {
	Log.Info(arg...)
}

//Infof implement for zap log infof method.
func Infof(template string, arg ...interface{}) {
	Log.Infof(template, arg...)
}

//Error implement for zap log error method.
func Error(arg ...interface{}) {
	Log.Error(arg...)
}

//Errorf implement for zap log errorf method.
func Errorf(template string, arg ...interface{}) {
	Log.Errorf(template, arg...)
}

//Fatal implement for zap log fatal method.
func Fatal(arg ...interface{}) {
	Log.Fatal(arg...)
}

//Fatalf implement for zap log fatalf method.
func Fatalf(template string, arg ...interface{}) {
	Log.Fatalf(template, arg...)
}

//Panic implement for zap log panic method.
func Panic(arg ...interface{}) {
	Log.Panic(arg...)
}

//Panicf implement for zap log panicf method.
func Panicf(template string, arg ...interface{}) {
	Log.Panicf(template, arg...)
}

//Warn implement for zap log warn method.
func Warn(arg ...interface{}) {
	Log.Warn(arg...)
}

//Warnf implement for zap log warnf method.
func Warnf(template string, arg ...interface{}) {
	Log.Warnf(template, arg...)
}

//Debug implement for zap log debug method.
func Debug(arg ...interface{}) {
	Log.Debug(arg...)
}

//Debugf implement for zap log debugf method.
func Debugf(template string, arg ...interface{}) {
	Log.Debugf(template, arg...)
}
