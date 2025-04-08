package logger

import (
	"io"
	"log"
	"os"

	"github.com/Fl0rencess720/Springbroad/consts"

	"github.com/natefinch/lumberjack"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Init(logFilePath string) {

	if logFilePath == "" {
		logFilePath = consts.DefaultLogFilePath
	}

	zap.ReplaceGlobals(zap.New(
		zapcore.NewCore(
			getEncoder(),
			getWriteSyncer(logFilePath),
			getLogLevel(),
		),
		zap.AddCaller(),
	))

}

func Sync(l *zap.Logger) {
	err := l.Sync()
	if err != nil {
		log.Println("Error when flushing zap logger")
	}
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	return zapcore.NewJSONEncoder(encoderConfig)
}

func getLogLevel() zapcore.Level {
	if viper.GetString("project.mode") == "dev" {
		return zapcore.DebugLevel
	} else {
		return zapcore.ErrorLevel
	}
}

func getWriteSyncer(logFilePath string) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   logFilePath,
		MaxSize:    5,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   true,
	}

	if viper.GetString("project.mode") == "dev" {
		return zapcore.AddSync(io.MultiWriter(os.Stdout, lumberJackLogger))
	}
	return zapcore.AddSync(lumberJackLogger)
}
