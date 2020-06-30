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

var Log *ZapLogger

func NewZapLogger(opts ...Option) *ZapLogger {
	zapLogger := &ZapLogger{
		Options: new(Options),
	}
	// 配置
	configure(zapLogger, opts...)
	// 未设置日志对象，则创建一个
	if zapLogger.Options.Logger == nil {
		// 创建zap日志对象
		syncWriter := zapcore.AddSync(&lumberjack.Logger{
			Filename:   zapLogger.Options.Filename,
			MaxSize:    int(zapLogger.Options.MaxSize),
			LocalTime:  zapLogger.Options.LocalTime,
			Compress:   zapLogger.Options.Compress,
			MaxBackups: zapLogger.Options.MaxBackups,
		})
		pEncoder := zap.NewDevelopmentEncoderConfig()
		pEncoder.EncodeTime = zapcore.ISO8601TimeEncoder // 时间格式
		encoder := zap.NewProductionEncoderConfig()
		encoder.EncodeLevel = zapcore.CapitalColorLevelEncoder //命令行颜色
		encoder.EncodeTime = zapcore.ISO8601TimeEncoder
		core := zapcore.NewTee(
			zapcore.NewCore(zapcore.NewConsoleEncoder(pEncoder), syncWriter, zap.NewAtomicLevelAt(zapLogger.Options.Level)),                //生产环境核心 ， 输出日志至文件
			zapcore.NewCore(zapcore.NewConsoleEncoder(encoder), zapcore.AddSync(os.Stdout), zap.NewAtomicLevelAt(zapLogger.Options.Level)), //开发环境核心，输出带颜色参数的日志至命令行
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

//use ZapLogger to implement middleware
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


func Info(arg ...interface{}){
	Log.Info(arg...)
}

func Infof(template string,arg ...interface{})  {
	Log.Infof(template,arg...)
}

func Error(arg ...interface{}){
	Log.Error(arg...)
}

func Errorf(template string,arg ...interface{}){
	Log.Errorf(template ,arg...)
}

func Fatal(arg ...interface{}){
	Log.Fatal(arg...)
}

func Fatalf(template string,arg ...interface{}){
	Log.Fatalf(template ,arg...)
}

func Panic(arg ...interface{}){
	Log.Panic(arg...)
}

func Panicf(template string,arg ...interface{}){
	Log.Panicf(template ,arg...)
}

func Warn(arg ...interface{}){
	Log.Warn(arg...)
}

func Warnf(template string,arg ...interface{}){
	Log.Warnf(template ,arg...)
}

func Debug(arg ...interface{}){
	Log.Debug(arg...)
}

func Debugf(template string,arg ...interface{}){
	Log.Debugf(template,arg...)
}